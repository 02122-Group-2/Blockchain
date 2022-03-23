package cryptography

import (
	"crypto/sha256"
)

// Takes a Block in JSON string format and calculates the 32-byte hash of this block and returns it.
func HashBlock(blockString string) [32]byte {
	return sha256.Sum256([]byte(blockString)) 
}
