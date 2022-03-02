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

func TestApplyLegalTransaction(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("Magn", "Niels", 42.0)
	err := state.AddTransaction(tr)
	if err != nil {
		t.Error("Failed to add transaction. Error: " + err.Error())
	}
}

func TestApplyIllegalTransaction(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("Magn", "Niels", 89898.0)
	err := state.AddTransaction(tr)
	if err == nil {
		t.Error("Succesfully added transaction but expected to fail.")
	}
}

func TestSendMoneyToSameUser(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("Magn", "Magn", 100.0)
	err := state.AddTransaction(tr)
	if err == nil {
		t.Error("Normal transaction from account to itself is not allowed")
	}
}

func TestApplyTransactionWithNegativeAmount(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("Magn", "Niels", -10.0)
	err := state.AddTransaction(tr)
	if err == nil {
		t.Error("Succesfully added transaction but expected to fail.")
	}
}

func TestAddTransactionFromAnUnknownAccount(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("llll", "Niels", 1.0)
	err := state.AddTransaction(tr)
	if err == nil {
		t.Error("Shouldnt be able to make a transaction from an unknown account")
	}
}

func TestAddTransactionToAnUnknownAccount(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("Niels", "gggg", 1.0)
	err := state.AddTransaction(tr)
	if err != nil {
		t.Error("Should be able to send to unknown account")
	}
}

func TestAddRewardToAccount(t *testing.T) {
	t.Log("Begin test reward system")
	tr := state.CreateReward("Alberto", 5000)
	err := state.AddTransaction(tr)

	if err != nil {
		t.Error("Unable to add reward to user")
	}
}
