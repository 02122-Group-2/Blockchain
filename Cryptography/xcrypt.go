package cryptography

import (
	"crypto/sha256"
)


func HashBlock(blockString string) [32]byte {
	return sha256.Sum256([]byte(blockString)) 
}
