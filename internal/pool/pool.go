package pool

import (
	"paychain/internal/blockchain"
	"sync"
)

type Pool struct {
	mu  sync.Mutex
	txs []blockchain.Transaction
}

func NewPool() *Pool {
	return &Pool{txs: make([]blockchain.Transaction, 0, 128)}
}

func (p *Pool) AddTransaction(tx blockchain.Transaction) {
	p.mu.Lock()
	p.txs = append(p.txs, tx)
	p.mu.Unlock()
}

func (p *Pool) GetBatch(n int) []blockchain.Transaction {
	p.mu.Lock()
	defer p.mu.Unlock()
	if n <= 0 || len(p.txs) == 0 {
		return nil
	}
	if n > len(p.txs) {
		n = len(p.txs)
	}
	batch := make([]blockchain.Transaction, n)
	copy(batch, p.txs[:n])
	p.txs = append([]blockchain.Transaction(nil), p.txs[n:]...)
	return batch
}

func (p *Pool) List() []blockchain.Transaction {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]blockchain.Transaction, len(p.txs))
	copy(out, p.txs)
	return out
}

func (p *Pool) Size() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.txs)
}
