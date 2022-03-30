package database

import (
	"fmt"
	"testing"
)

func TestLoadState(t *testing.T) {
	t.Log("Start load state test")
	state := LoadState()
	fmt.Println(state)
}
