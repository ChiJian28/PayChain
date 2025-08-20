package blockchain

import (
	"sync"
)

// Chain is a thread-safe blockchain storage.
type Chain struct {
	mu     sync.RWMutex
	blocks []Block
}

func NewChain(genesis Block) *Chain {
	return &Chain{blocks: []Block{genesis}}
}

func (c *Chain) LastBlock() Block {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.blocks[len(c.blocks)-1]
}

func (c *Chain) Append(block Block) {
	c.mu.Lock()
	c.blocks = append(c.blocks, block)
	c.mu.Unlock()
}

func (c *Chain) All() []Block {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Block, len(c.blocks))
	copy(out, c.blocks)
	return out
}
