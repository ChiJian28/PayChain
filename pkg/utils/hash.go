package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// HashStrings computes SHA-256 over the '|' joined parts.
func HashStrings(parts ...string) string {
	h := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(h[:])
}
