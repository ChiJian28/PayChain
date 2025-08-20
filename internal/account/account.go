package account

import (
	"paychain/internal/blockchain"
	"sync"
)

type Store struct {
	mu      sync.RWMutex
	balance map[string]int
}

func NewStore() *Store {
	return &Store{balance: make(map[string]int)}
}

func (s *Store) GetBalance(user string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.balance[user]
}

// ApplyTransaction updates balances if the sender has enough funds.
// Returns true if applied; false otherwise.
func (s *Store) ApplyTransaction(tx blockchain.Transaction) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if tx.Amount <= 0 {
		return false
	}
	if tx.From != "" {
		if s.balance[tx.From] < tx.Amount {
			return false
		}
		s.balance[tx.From] -= tx.Amount
	}
	if tx.To != "" {
		s.balance[tx.To] += tx.Amount
	}
	return true
}
