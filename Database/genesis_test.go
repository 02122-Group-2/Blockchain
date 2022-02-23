package database

import (
	"testing"
)

func TestLoadGenesis(t *testing.T) {
	t.Log("Begin load genesis test")
	loadedGen := LoadGenesis()
	if &loadedGen == nil {
		t.Error()
	}
}

func TestLoadTransactions(t *testing.T) {
	t.Log("Begin load transactions test")
	loadedTrans := LoadTransactions()
	if &loadedTrans == nil {
		t.Error()
	}
}
