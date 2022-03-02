package cryptography

import (
	db "blockchain/database"
	"crypto/sha256"
	"fmt"
	"testing"
)

var b = db.Block{
	db.BlockHeader{
		"hash420",
		1234,
		69,
	},
	[]db.Transaction{
		{"bruh", "bruh2", 100, 1235, "transaction"},
		{"bruh2", "bruh3", 2000, 1236, "transaction"},
		{"bruh3", "bruh", 1000, 1237, "transaction"},
	},
}

var hash = HashBlock(b)

func TestHashBlockSameJson(t *testing.T) {
	jsonString := "{\"Header\":{\"ParentHash\":\"hash420\",\"CreatedAt\":1234,\"SerialNo\":69},\"Transactions\":[{\"From\":\"bruh\",\"To\":\"bruh2\",\"Amount\":100,\"Timestamp\":1235,\"Type\":\"transaction\"},{\"From\":\"bruh2\",\"To\":\"bruh3\",\"Amount\":2000,\"Timestamp\":1236,\"Type\":\"transaction\"},{\"From\":\"bruh3\",\"To\":\"bruh\",\"Amount\":1000,\"Timestamp\":1237,\"Type\":\"transaction\"}]}"

	if fmt.Sprintf("%x", hash) != fmt.Sprintf("%x", sha256.Sum256([]byte(jsonString))) {
		t.Errorf("hashing did not work according to protocol")
	}

}

func TestHashBlockDifferentJson(t *testing.T) {
	jsonString := "2{\"Header\":{\"ParentHash\":\"hash420\",\"CreatedAt\":1234,\"SerialNo\":69},\"Transactions\":[{\"From\":\"bruh\",\"To\":\"bruh2\",\"Amount\":100,\"Timestamp\":1235,\"Type\":\"transaction\"},{\"From\":\"bruh2\",\"To\":\"bruh3\",\"Amount\":2000,\"Timestamp\":1236,\"Type\":\"transaction\"},{\"From\":\"bruh3\",\"To\":\"bruh\",\"Amount\":1000,\"Timestamp\":1237,\"Type\":\"transaction\"}]}"

	if fmt.Sprintf("%x", hash) == fmt.Sprintf("%x", sha256.Sum256([]byte(jsonString))) {
		t.Errorf("hashing did not work according to protocol")
	}

}
