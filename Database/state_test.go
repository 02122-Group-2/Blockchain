package database

import (
	"fmt"
	"testing"
)

func TestLoadState(t *testing.T) {
	t.Log("Start load state test")
	state, err := LoadState()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(state)
}
