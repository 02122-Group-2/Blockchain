package cryptography

import (
	"crypto/sha256"
)

// * File: Niels, s204503

// Takes a Block in JSON string format and calculates the 32-byte hash of this block and returns it.
func HashBlock(blockString string) [32]byte {
	return sha256.Sum256([]byte(blockString))
}

// Takes a transationc is Json string format and hashes it - Yes it does exactly the same as above, just named differently.
func HashTransaction(transactionString string) [32]byte {
	return sha256.Sum256([]byte(transactionString))
}
