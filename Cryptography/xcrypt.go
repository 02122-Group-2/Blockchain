package cryptography

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	db "blockchain/database"
)

type Hash [32]byte

func HashBlock(block db.Block) Hash {
	bJson, err := json.Marshal(block)
	if err != nil {
		fmt.Errorf("Unable to convert to json string")
	}
	return sha256.Sum256(bJson)
}
