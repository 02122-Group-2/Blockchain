package database

import (
	"testing"
)

func TestLoadState(t *testing.T) {
	t.Log("Start load state test")
	state := LoadState()
	t.Log(state)
}
