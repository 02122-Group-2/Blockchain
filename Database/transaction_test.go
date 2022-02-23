package database

import (
	"testing"
)

// type State struct {
// 	Balances  map[AccountAddress]uint
// 	txMempool []Transaction
// 	dbFile    *os.File

// 	lastTxSerialNo    int
// 	lastBlockSerialNo int
// 	latestHash        string
// }

var state, _ = LoadState()

func TestCreate(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("magn", "niels", 6969.0)

	if tr.Amount != 6969.0 {
		t.Errorf("Amount is set wrong")
	}
	if tr.From != "magn" {
		t.Errorf("From is set wrong")
	}
	if tr.To != "niels" {
		t.Errorf("To is set wrong")
	}
	if tr.Type != "transaction" {
		t.Errorf("Type is wrong")
	}
}

func TestReward(t *testing.T) {
	r := state.CreateReward("niels", 1337.420)

	if r.Amount != 1337.420 {
		t.Errorf("Amount is set wrong")
	}
	if r.From != "system" {
		t.Errorf("From is set wrong")
	}
	if r.Type != "reward" {
		t.Errorf("Type is wrong")
	}
}
