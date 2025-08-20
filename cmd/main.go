package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"paychain/internal/account"
	"paychain/internal/api"
	"paychain/internal/blockchain"
	"paychain/internal/kafka"
	txpool "paychain/internal/pool"
	"paychain/pkg/logger"
	"paychain/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// basic wiring
	pool := txpool.NewPool()
	acct := account.NewStore()

	// genesis block
	genesis := blockchain.Block{Index: 0, Timestamp: utils.NowUnix(), PrevHash: "", Nonce: 0}
	genesis.Hash = blockchain.ComputeBlockHash(genesis)
	chain := blockchain.NewChain(genesis)

	// preload faucet balance for demo
	acct.ApplyTransaction(blockchain.Transaction{From: "", To: "alice", Amount: 1000, Time: utils.NowUnix()})
	acct.ApplyTransaction(blockchain.Transaction{From: "", To: "bob", Amount: 1000, Time: utils.NowUnix()})

	// Kafka
	brokers := []string{"localhost:9092"}
	topic := "paychain-transactions"
	prod, err := kafka.NewProducer(brokers, topic)
	if err != nil {
		logger.Errorf("producer init: %v", err)
		return
	}
	defer prod.Close()

	groupID := "paychain-consumers"
	cons, err := kafka.NewConsumer(brokers, groupID, topic, pool)
	if err != nil {
		logger.Errorf("consumer init: %v", err)
		return
	}
	defer cons.Close()

	ctx, cancel := signalContext()
	defer cancel()

	// start consumer in background
	go func() {
		if err := cons.Start(ctx); err != nil && ctx.Err() == nil {
			logger.Errorf("consumer err: %v", err)
		}
	}()

	// block packer goroutine
	go func() {
		const batchSize = 3
		const difficulty = 3
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			batch := pool.GetBatch(batchSize)
			if len(batch) == 0 {
				time.Sleep(200 * time.Millisecond)
				continue
			}

			last := chain.LastBlock()
			candidate := blockchain.Block{
				Index:        last.Index + 1,
				Timestamp:    utils.NowUnix(),
				Transactions: batch,
				PrevHash:     last.Hash,
			}

			// Mine with cancelable context
			mineCtx, mineCancel := context.WithCancel(ctx)
			mined, ok := blockchain.MineBlock(mineCtx, candidate, difficulty)
			mineCancel()
			if !ok {
				continue
			}

			// apply transactions to accounts
			applied := make([]blockchain.Transaction, 0, len(mined.Transactions))
			for _, tx := range mined.Transactions {
				if acct.ApplyTransaction(tx) {
					applied = append(applied, tx)
				}
			}
			mined.Transactions = applied
			mined.Hash = blockchain.ComputeBlockHash(mined)
			chain.Append(mined)
			logger.Infof("new block %d, tx=%d, hash=%s", mined.Index, len(applied), mined.Hash)
		}
	}()

	// Gin API
	r := gin.Default()
	api.RegisterRoutes(r, prod, acct, chain, pool)
	if err := r.Run(":8080"); err != nil {
		logger.Errorf("gin run: %v", err)
	}
}

func signalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()
	return ctx, cancel
}
