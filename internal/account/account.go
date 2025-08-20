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

// FilterApplicableTransactions returns a subsequence of txs that are valid
// against a snapshot of current balances, applied in order without mutating
// the real store. This is used for pre-validation before mining.
func (s *Store) FilterApplicableTransactions(txs []blockchain.Transaction) []blockchain.Transaction {
	s.mu.RLock()
	// make a snapshot copy
	snapshot := make(map[string]int, len(s.balance)+len(txs))
	for k, v := range s.balance {
		snapshot[k] = v
	}
	s.mu.RUnlock()

	applied := make([]blockchain.Transaction, 0, len(txs))
	for _, tx := range txs {
		if tx.Amount <= 0 {
			continue
		}
		if tx.From != "" {
			if snapshot[tx.From] < tx.Amount {
				continue
			}
			snapshot[tx.From] -= tx.Amount
		}
		if tx.To != "" {
			snapshot[tx.To] += tx.Amount
		}
		applied = append(applied, tx)
	}
	return applied
}

// ApplyBatchIfValid re-validates the given txs against the latest balances
// and applies them atomically if all are valid. Returns true on success.
func (s *Store) ApplyBatchIfValid(txs []blockchain.Transaction) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// simulate on a temp map first
	temp := make(map[string]int, len(s.balance)+len(txs))
	for k, v := range s.balance {
		temp[k] = v
	}
	for _, tx := range txs {
		if tx.Amount <= 0 {
			return false
		}
		if tx.From != "" {
			if temp[tx.From] < tx.Amount {
				return false
			}
			temp[tx.From] -= tx.Amount
		}
		if tx.To != "" {
			temp[tx.To] += tx.Amount
		}
	}
	// commit atomically by replacing the map reference
	s.balance = temp
	return true
}