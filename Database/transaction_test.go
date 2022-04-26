package database

import (
	Crypto "blockchain/Cryptography"
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
var walletUsername1 = "testingWallet1"
var walletUsername2 = "testingWallet265"
var walletUsername3 = "testingWaller333"
var pswd = "password"

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

	if r.Tx.Amount != 1337.420 {
		t.Errorf("Amount is set wrong")
	}
	if r.Tx.From != "system" {
		t.Errorf("From is set wrong")
	}
	if r.Tx.Type != "reward" {
		t.Errorf("Type is wrong")
	}

	ResetTest()
}

func TestApplyLegalTransaction(t *testing.T) {
	t.Log("begin create transaction test")

	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state.AccountBalances[AccountAddress(testWallet.Address)] = 1000

	senderBalanceBefore := state.AccountBalances[AccountAddress(testWallet.Address)]
	receiverBalanceBefore := state.AccountBalances["Niels"]

	tr := state.CreateTransaction(AccountAddress(testWallet.Address), "Niels", 42.0)
	signedTx, err := state.SignTransaction(testWallet, pswd, tr)

	if (err != nil) {
		t.Error("failed to sign transaction. Get error: " + err.Error())
	}

	err = state.AddTransaction(signedTx)
	if err != nil {
		t.Error("Failed to add transaction. Error: " + err.Error())
	}

	if senderBalanceBefore - 42 != state.AccountBalances[AccountAddress(testWallet.Address)] {
		t.Error("Sender should have lost 42 tokens")
	}

	if receiverBalanceBefore + 42 != state.AccountBalances["Niels"] {
		t.Error("Receiver should have received 42 tokens")
	}

	testWallet.Delete()
	ResetTest()
}

func TestApplyIllegalTransaction(t *testing.T) {
	t.Log("begin create transaction test")

	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state.AccountBalances[AccountAddress(testWallet.Address)] = 1000

	tr := state.CreateTransaction(AccountAddress(testWallet.Address), "Niels", 898989.0)
	signedTx, _ := state.SignTransaction(testWallet, pswd, tr)
	err := state.AddTransaction(signedTx)
	if err == nil {
		t.Error("Succesfully added transaction but expected to fail.")
	}

	testWallet.Delete()
	ResetTest()
}

func TestSendMoneyToSameUser(t *testing.T) {
	t.Log("begin create transaction test")

	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state.AccountBalances[AccountAddress(testWallet.Address)] = 1000

	tr := state.CreateTransaction(AccountAddress(testWallet.Address), AccountAddress(testWallet.Address), 100.0)
	signedTx, _ := state.SignTransaction(testWallet, pswd, tr)
	err := state.AddTransaction(signedTx)
	if err == nil {
		t.Error("Normal transaction from account to itself is not allowed")
	}

	testWallet.Delete()
	ResetTest()
}

func TestApplyTransactionWithNegativeAmount(t *testing.T) {
	t.Log("begin create transaction test")

	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state.AccountBalances[AccountAddress(testWallet.Address)] = 1000

	tr := state.CreateTransaction(AccountAddress(testWallet.Address), "Niels", -10.0)
	signedTx, _ := state.SignTransaction(testWallet, pswd, tr)
	err := state.AddTransaction(signedTx)
	if err == nil {
		t.Error("Succesfully added transaction but expected to fail.")
	}

	testWallet.Delete()
	ResetTest()
}

func TestAddTransactionFromAnUnknownAccount(t *testing.T) {
	t.Log("begin create transaction test")

	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet, _ := Crypto.AccessWallet(walletUsername1, pswd)

	tr := state.CreateTransaction("llll", "Niels", 1.0)
	signedTx, _ := state.SignTransaction(testWallet, pswd, tr)
	err := state.AddTransaction(signedTx)
	if err == nil {
		t.Error("Shouldnt be able to make a transaction from an unknown account")
	}

	testWallet.Delete()
	ResetTest()
}

func TestAddTransactionToAnUnknownAccount(t *testing.T) {
	t.Log("begin create transaction test")

	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state.AccountBalances[AccountAddress(testWallet.Address)] = 1000

	tr := state.CreateTransaction(AccountAddress(testWallet.Address), "gggg", 1.0)
	signedTx, _ := state.SignTransaction(testWallet, pswd, tr)
	err := state.AddTransaction(signedTx)
	if err != nil {
		t.Error("Should be able to send to unknown account")
	}

	testWallet.Delete()
	ResetTest()
}

func TestAddRewardToAccount(t *testing.T) {
	t.Log("Begin test reward system")

	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state.AccountBalances[AccountAddress(testWallet.Address)] = 1000	

	rewardTx := state.CreateReward("Alberto", 5000)
	err := state.AddTransaction(rewardTx)

	if err != nil {
		t.Error("Unable to add reward to user. Gor error: " + err.Error())
	}

	testWallet.Delete()
	ResetTest()
}

func TestCreateLegalTransactionAndPersist(t *testing.T) {
	t.Log("Begin test persisting to transaction.JSON")

	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state.AccountBalances[AccountAddress(testWallet.Address)] = 1000	

	tr := state.CreateTransaction(AccountAddress(testWallet.Address), "Magn", 12.0)
	signedTx, _ := state.SignTransaction(testWallet, pswd, tr)
	err := state.AddTransaction(signedTx)
	if err != nil {
		t.Error("Failed to add transaction. Error: " + err.Error())
	}

	testWallet.Delete()
	SaveTransaction(state.TxMempool)
	ResetTest()
}

func TestAddTransactionAndCheckTheyAreSaved(t *testing.T) {
	state1 := LoadState()

	Crypto.CreateNewWallet(walletUsername3, pswd)
	testWallet1, _ := Crypto.AccessWallet(walletUsername3, pswd)
	state.AccountBalances[AccountAddress(testWallet1.Address)] = 1000

	Crypto.CreateNewWallet(walletUsername2, pswd)
	testWallet2, _ := Crypto.AccessWallet(walletUsername2, pswd)
	state.AccountBalances[AccountAddress(testWallet2.Address)] = 1000

	signedTx1, _ := state.SignTransaction(testWallet1, pswd, state.CreateTransaction(AccountAddress(testWallet1.Address), AccountAddress(testWallet2.Address), 10))
	signedTx2, _ := state.SignTransaction(testWallet2, pswd, state.CreateTransaction(AccountAddress(testWallet2.Address), AccountAddress(testWallet1.Address), 10))

	state1.AddTransaction(signedTx1)
	state1.AddTransaction(signedTx2)

	state2 := LoadState()

	state1Json, _ := state1.MarshalJSON()
	state2Json, _ := state2.MarshalJSON()
	if string(state1Json) != string(state2Json) {
		t.Error("the local changes should be saved but are not")
	}

	testWallet1.Delete()
	testWallet2.Delete()
	ResetTest()
}
