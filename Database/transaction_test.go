package database

import (
	shared "blockchain/Shared"
	"testing"
)

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

	ResetTest()
}

func TestReward(t *testing.T) {
	r := state.CreateReward("niels", 1337.420)

	if r.Amount != 1337.420 {
		t.Errorf("Amount is set wrong")
	}
	if r.From != "system" {
		t.Errorf("From is set wrong")
	}
	if r.To != "niels" {
		t.Errorf("To is set wrong")
	}
	if r.Type != "reward" {
		t.Errorf("Type is wrong")
	}

	ResetTest()
}

func EmptyBlockchain() {
	ClearBlockchain()
	ClearTransactions()
	state.ClearState()
}

func TestGenesis(t *testing.T) {
	EmptyBlockchain()
	g := state.CreateGenesisTransaction("asger", 42.42)

	if g.Amount != 42.42 {
		t.Errorf("Amount is set wrong")
	}
	if g.From != "system" {
		t.Errorf("From is set wrong")
	}
	if g.To != "asger" {
		t.Errorf("To is set wrong")
	}
	if g.Type != "genesis" {
		t.Errorf("Type is wrong")
	}
	ResetTest()
}

func TestAddLegalGenesisTransaction(t *testing.T) {
	t.Log("begin add legal genesis transaction test")

	EmptyBlockchain()
	g := state.CreateGenesisTransaction("asger", 42.42)

	err := state.AddTransaction(g)
	if err != nil {
		t.Error("Failed to add transaction. Error: " + err.Error())
	}

	ResetTest()
}

func TestAddIllegalGenesisTransaction(t *testing.T) {
	// for this specific test case, the following reset is needed
	shared.ResetPersistenceFilesForTest()
	state = LoadState()
	t.Log("begin create illegal genesis transaction test")

	g := state.CreateGenesisTransaction("asger", 666.66)
	err := state.AddTransaction(g)
	if err == nil {
		t.Error("Adding genesis transactions later in blockchain is not allowed")
	}
	ResetTest()
}

func TestAddLegalTransaction(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("Magn", "Niels", 42.0)
	err := state.AddTransaction(tr)
	if err != nil {
		t.Error("Failed to add transaction. Error: " + err.Error())
	}

	ResetTest()
}

func TestAddIllegalTransaction(t *testing.T) {
	t.Log("begin create too large transaction test")

	tr := state.CreateTransaction("Magn", "Niels", 898989.0)
	err := state.AddTransaction(tr)
	if err == nil {
		t.Error("Succesfully added transaction but expected to fail.")
	}

	ResetTest()
}

func TestSendMoneyToSameUser(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("Magn", "Magn", 100.0)
	err := state.AddTransaction(tr)
	if err == nil {
		t.Error("Normal transaction from account to itself is not allowed")
	}

	ResetTest()
}

func TestAddTransactionWithNegativeAmount(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("Magn", "Niels", -10.0)
	err := state.AddTransaction(tr)
	if err == nil {
		t.Error("Succesfully added transaction but expected to fail.")
	}

	ResetTest()
}

func TestAddTransactionFromAnUnknownAccount(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("llll", "Niels", 1.0)
	err := state.AddTransaction(tr)
	if err == nil {
		t.Error("Shouldnt be able to make a transaction from an unknown account")
	}

	ResetTest()
}

func TestAddTransactionToAnUnknownAccount(t *testing.T) {
	t.Log("begin create transaction test")

	tr := state.CreateTransaction("Niels", "gggg", 1.0)
	err := state.AddTransaction(tr)
	if err != nil {
		t.Error("Should be able to send to unknown account")
	}

	ResetTest()
}

func TestAddTransactionWithWrongNounce(t *testing.T) {
	t.Log("begin create transaction with wrong account nounce")

	shared.ResetPersistenceFilesForTest()
	tr := state.CreateTransaction("Emilie", "Niels", 2.0)
	tr.SenderNounce = 2
	err := state.AddTransaction(tr)
	if err == nil {
		t.Error("Should not be able to add transactions with older nounces")
	}
	ResetTest()
}

func TestAddRewardToAccount(t *testing.T) {
	t.Log("Begin test reward system")
	tr := state.CreateReward("Alberto", 5000)
	err := state.AddTransaction(tr)

	if err != nil {
		t.Error("Unable to add reward to user")
	}

	ResetTest()
}

func TestCreateLegalTransactionAndPersist(t *testing.T) {
	t.Log("Begin test persisting to transaction.JSON")
	tr := state.CreateTransaction("Niels", "Magn", 200000.0)
	err := state.AddTransaction(tr)
	if err != nil {
		t.Error("Failed to add transaction. Error: " + err.Error())
	}

	SaveTransaction(state.TxMempool)
	ResetTest()
}

func TestAddTransactionAndCheckTheyAreSaved(t *testing.T) {
	state1 := LoadState()
	state1.AddTransaction(state.CreateTransaction("Magn", "Niels", 10))
	state1.AddTransaction(state.CreateTransaction("Niels", "Magn", 10))

	state2 := LoadState()

	state1Json, _ := state1.MarshalJSON()
	state2Json, _ := state2.MarshalJSON()
	if string(state1Json) != string(state2Json) {
		panic("the local changes should be saved but are not")
	}

	ResetTest()
}
