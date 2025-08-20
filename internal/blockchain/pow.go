package blockchain

import (
	"context"
	"runtime"
)

// MineBlock performs a simple PoW: find Nonce s.t. hash has `difficulty` leading zeros.
func MineBlock(ctx context.Context, base Block, difficulty int) (Block, bool) {
	if difficulty < 1 {
		difficulty = 1
	}
	targetPrefix := make([]byte, difficulty)
	for i := range targetPrefix {
		targetPrefix[i] = '0'
	}
	workers := runtime.NumCPU()
	type result struct {
		block Block
		ok    bool
	}
	resultCh := make(chan result, 1)
	// Launch workers trying disjoint nonce ranges
	for w := 0; w < workers; w++ {
		start := w
		go func(start int) {
			blk := base
			for nonce := start; ; nonce += workers {
				select {
				case <-ctx.Done():
					return
				default:
				}
				blk.Nonce = nonce
				blk.Hash = ComputeBlockHash(blk)
				if hasPrefix(blk.Hash, string(targetPrefix)) {
					select {
					case resultCh <- result{block: blk, ok: true}:
					default:
					}
					return
				}
			}
		}(start)
	}

	select {
	case <-ctx.Done():
		return Block{}, false
	case r := <-resultCh:
		return r.block, r.ok
	}
}

func hasPrefix(s, prefix string) bool {
	if len(prefix) == 0 {
		return true
	}
	if len(s) < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if s[i] != prefix[i] {
			return false
		}
	}
	return true
}
