package cryptography

import (
	"crypto/sha256"
)


func HashBlock(blockString string) string {
	_hash := sha256.Sum256([]byte(blockString)) 
	return string(_hash[:])
}
