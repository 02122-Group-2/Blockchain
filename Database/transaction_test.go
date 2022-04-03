package database

import (
	"fmt"
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

var state = LoadState()

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

	resetTest()
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

	resetTest()
}

func TestApplyLegalTransaction(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("Asger", "Niels", 42.0)
	err := state.AddTransaction(tr)
	if err != nil {
		t.Error("Failed to add transaction. Error: " + err.Error())
	}

	resetTest()
}

func TestApplyIllegalTransaction(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("Magn", "Niels", 898989.0)
	err := state.AddTransaction(tr)
	if err == nil {
		t.Error("Succesfully added transaction but expected to fail.")
	}

	resetTest()
}

func TestSendMoneyToSameUser(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("Magn", "Magn", 100.0)
	err := state.AddTransaction(tr)
	if err == nil {
		t.Error("Normal transaction from account to itself is not allowed")
	}

	resetTest()
}

func TestApplyTransactionWithNegativeAmount(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("Magn", "Niels", -10.0)
	err := state.AddTransaction(tr)
	if err == nil {
		t.Error("Succesfully added transaction but expected to fail.")
	}

	resetTest()
}

func TestAddTransactionFromAnUnknownAccount(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("llll", "Niels", 1.0)
	err := state.AddTransaction(tr)
	if err == nil {
		t.Error("Shouldnt be able to make a transaction from an unknown account")
	}

	resetTest()
}

func TestAddTransactionToAnUnknownAccount(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("Niels", "gggg", 1.0)
	err := state.AddTransaction(tr)
	if err != nil {
		t.Error("Should be able to send to unknown account")
	}

	resetTest()
}

func TestAddRewardToAccount(t *testing.T) {
	t.Log("Begin test reward system")
	tr := state.CreateReward("Alberto", 5000)
	err := state.AddTransaction(tr)

	if err != nil {
		t.Error("Unable to add reward to user")
	}

	resetTest()
}

func TestCreateLegalTransactionAndPersist(t *testing.T) {
	t.Log("Begin test persisting to transaction.JSON")
	tr := state.CreateTransaction("Niels", "Magn", 200000.0)
	err := state.AddTransaction(tr)
	if err != nil {
		t.Error("Failed to add transaction. Error: " + err.Error())
	}

	SaveTransaction(state.TxMempool)
	resetTest()
}

func TestAddTransactionAndCheckTheyAreSaved(t *testing.T) {
	state1 := LoadState()
	state1.AddTransaction(state.CreateTransaction("Magn", "Niels", 10))
	state1.AddTransaction(state.CreateTransaction("Niels", "Magn", 10))

	state2 := LoadState()

	state1Json, _ := state1.MarshalJSON()
	state2Json, _ := state2.MarshalJSON()
	if string(state1Json) != string(state2Json) {
		fmt.Errorf("the local changes should be saved but are not")
	}

	resetTest()
}



