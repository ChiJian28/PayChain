package kafka

import (
	"context"
	"encoding/json"
	"paychain/internal/blockchain"
	txpool "paychain/internal/pool"

	"github.com/IBM/sarama"
)

type Consumer struct {
	group sarama.ConsumerGroup
	topic string
	pool  *txpool.Pool
}

func NewConsumer(brokers []string, groupID string, topic string, pool *txpool.Pool) (*Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_1_0_0
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	g, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, err
	}
	return &Consumer{group: g, topic: topic, pool: pool}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	handler := &groupHandler{pool: c.pool}
	for {
		if err := c.group.Consume(ctx, []string{c.topic}, handler); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (c *Consumer) Close() error { return c.group.Close() }

type groupHandler struct {
	pool *txpool.Pool
}

func (h *groupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *groupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *groupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var tx blockchain.Transaction
		if err := json.Unmarshal(msg.Value, &tx); err == nil {
			h.pool.AddTransaction(tx)
			sess.MarkMessage(msg, "")
		} else {
			// malformed message, skip but still mark
			sess.MarkMessage(msg, "bad")
		}
	}
	return nil
}
