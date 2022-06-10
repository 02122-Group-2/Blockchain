package database

import (
	shared "blockchain/Shared"
	"testing"
)

func TestLoadState(t *testing.T) {
	shared.ResetPersistenceFilesForTest()
	t.Log("Start load state test")
	state := LoadState()
	t.Log(state)
}
