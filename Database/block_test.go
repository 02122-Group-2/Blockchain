package database

import (
	Crypto "blockchain/Cryptography"
	// "encoding/hex"
	// "encoding/json"
	"fmt"
	"testing"
)

var state_block = LoadState()
var blockchain_original = LoadBlockchain()
var state_original = LoadState()
var snapshot_orignal = LoadSnapshot()
var transactions_original = LoadTransactions()

func TestCreateBlock(t *testing.T) {
	// Create a wallet to test the functionality
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state.AccountBalances[AccountAddress(testWallet.Address)] = 1000

	// Dynamically set balance of wallet account for testing purposes.
	// This exploit would easily be detected in other systems. As the money would come from nowhere.
	state_block.AccountBalances[AccountAddress(testWallet.Address)] = 1000

	// Create the transactions and sign them with the created wallet
	signedTx1, _ := state_block.CreateSignedTransaction(testWallet, pswd, "Asger", 10)
	err1 := state_block.AddTransaction(signedTx1)

	signedTx2, _ := state_block.CreateSignedTransaction(testWallet, pswd, "Emilie", 4)
	err2 := state_block.AddTransaction(signedTx2)

	block := state_block.CreateBlock(state_block.TxMempool)
	fmt.Println(block)
	fmt.Println(err1)
	fmt.Println(err2)

	ResetTest()
	testWallet.Delete()
}

func TestSaveBlock(t *testing.T) {
	blockchain_original = LoadBlockchain()

	// Creates a wallet to test the functionality
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state_block.AccountBalances[AccountAddress(testWallet.Address)] = 1000

	// Create a block by creating the transactions and signing them with the new wallet
	signedTx1, _ := state_block.CreateSignedTransaction(testWallet, pswd, "Asger", 10)
	err1 := state_block.AddTransaction(signedTx1)

	signedTx2, _ := state_block.CreateSignedTransaction(testWallet, pswd, "Emilie", 4)
	err2 := state_block.AddTransaction(signedTx2)

	if err1 != nil || err2 != nil {
		t.Error("Adding one of the two transactions failed")
	}

	block := state_block.CreateBlock(state_block.TxMempool)

	// var blockList []Block

	blockList := append(blockchain_original, block)

	SaveBlockchain(blockList)

	ResetTest()
	testWallet.Delete()
}

func TestLoadBlockchain(t *testing.T) {
	res := LoadBlockchain()
	fmt.Println(res)
	ResetTest()
}

func TestAddBlockToBlockchain(t *testing.T) {
	// Start by ensuring the setup is correct
	ResetTest()
	state_block = LoadState()
	// Creates a wallet to test the functionality
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state_block.AccountBalances[AccountAddress(testWallet.Address)] = 1000

	// Save the block as the latest snapshot, as it would otherwise fail as the account is undefined.
	err := state_block.SaveSnapshot()

	// Create a block by creating the transactions and signing them with the new wallet
	signedTx1, _ := state_block.CreateSignedTransaction(testWallet, pswd, "Asger", 10)
	state_block.AddTransaction(signedTx1)

	signedTx2, _ := state_block.CreateSignedTransaction(testWallet, pswd, "Emilie", 4)
	state_block.AddTransaction(signedTx2)

	block2 := state_block.CreateBlock(state_block.TxMempool)

	signedTx3, _ := state_block.CreateSignedTransaction(testWallet, pswd, "Niels", 89)
	state_block.AddTransaction(signedTx3)

	err = state_block.AddBlock(block2)
	if err != nil || len(state_block.TxMempool) != 1 {
		t.Log(state_block.TxMempool)
		t.Log(len(state_block.TxMempool))
		t.Errorf("failed - expected no errors and that the length of the TxMemPool is 1")
	}

	ResetTest()
	testWallet.Delete()
}

// This tests makes sure the functionality of sharing the blocks work correctly.
// Two states will be created, who are orignally identical.
// One state will create some transactions, Then create a block.
// The other will create a few transactions too. The first and last should be invalidated when the block from the first state when it is synced.
func TestSeperateStatesShareBlock(t *testing.T) {
	stateOne := LoadSnapshot()
	stateTwo := stateOne.copyState()

	// Creates three wallets
	// Creates the first wallet to test the functionality
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet1, _ := Crypto.AccessWallet(walletUsername1, pswd)
	stateOne.AccountBalances[AccountAddress(testWallet1.Address)] = 1000
	stateTwo.AccountBalances[AccountAddress(testWallet1.Address)] = 1000

	// Creates the second wallet to test the functionality
	Crypto.CreateNewWallet(walletUsername2, pswd)
	testWallet2, _ := Crypto.AccessWallet(walletUsername2, pswd)
	stateOne.AccountBalances[AccountAddress(testWallet2.Address)] = 1000
	stateTwo.AccountBalances[AccountAddress(testWallet2.Address)] = 1000

	// Creates the third wallet to test the functionality
	Crypto.CreateNewWallet(walletUsername3, pswd)
	testWallet3, _ := Crypto.AccessWallet(walletUsername3, pswd)
	stateOne.AccountBalances[AccountAddress(testWallet3.Address)] = 1000
	stateTwo.AccountBalances[AccountAddress(testWallet3.Address)] = 1000

	// Saves the snapshot as the latest snapshot, otherwise the test would fail
	// Makes a copy of stateOne to get the original state.
	original_state := stateOne.copyState()
	stateOne.SaveSnapshot()

	signedTx1, _ := stateOne.CreateSignedTransaction(testWallet1, pswd, "Niels", 10)
	stateOne.AddTransaction(signedTx1)

	signedTx2, _ := stateOne.CreateSignedTransaction(testWallet2, pswd, "Magn", 10)
	stateOne.AddTransaction(signedTx2)

	signedTx3, _ := stateOne.CreateSignedTransaction(testWallet1, pswd, "Emilie", 10)
	stateOne.AddTransaction(signedTx3)

	blockOne := stateOne.CreateBlock(stateOne.TxMempool)

	signedTx1_2, _ := stateTwo.CreateSignedTransaction(testWallet1, pswd, "Niels", 10) // Should be invalid when merging the other block - Because of Nounces
	stateTwo.AddTransaction(signedTx1_2)

	signedTx2_2, _ := stateTwo.CreateSignedTransaction(testWallet3, pswd, "Emilie", 10) // Should be valid
	stateTwo.AddTransaction(signedTx2_2)

	signedTx3_2, _ := stateTwo.CreateSignedTransaction(testWallet2, pswd, "Asger", 10) // Should be invalid when merging the other block
	stateTwo.AddTransaction(signedTx3_2)

	err := stateOne.AddBlock(blockOne)
	if err != nil {
		SaveBlockchain(blockchain_original)
		t.Errorf("failed to add block to first state...")
	}

	// Saves the snapshot, since the snapshot is still "outdated" for the other account. This error is due to the fact that we run the software on the same pc.
	original_state.SaveSnapshot()

	err = stateTwo.AddBlock(blockOne)
	if err != nil {
		SaveBlockchain(blockchain_original)
		t.Errorf("failed to add block to second state...")
	}

	if len(stateOne.TxMempool) != 0 || len(stateTwo.TxMempool) != 1 {
		t.Errorf("failed - all transactions should be removed from the first state and one should remain in the last")
	}

	testWallet1.Delete()
	testWallet2.Delete()
	testWallet3.Delete()
	ResetTest()
}

/*
func TestMarshalUnmarshalBlock(t *testing.T) {
	txList := []Transaction{
		{
			From:      "Niels",
			To:        "Magn",
			Amount:    10,
			Timestamp: 1647079670026215000,
			Type:      "transaction",
		},
		{
			From:      "Magn",
			To:        "Emilie",
			Amount:    4,
			Timestamp: 1647079670578703300,
			Type:      "transaction",
		}}
	phStr := "d4b054173a82144cd6a7f4d7f2157f1504744626b6fe80eb0702cd688429ba43"
	ph, _ := hex.DecodeString(phStr)
	var ph32 [32]byte
	for i := 0; i < 32; i++ {
		ph32[i] = ph[i]
	}
	testBlock := Block{
		Header: BlockHeader{
			ParentHash: ph32,
			CreatedAt:  1647079671155969900,
			SerialNo:   4,
		},
		Transactions: txList,
	}
	jsonData, _ := json.Marshal(&testBlock)
	// t.Logf("%s", jsonData)
	data := Block{}
	unm_err := json.Unmarshal(jsonData, &data)
	if unm_err != nil {
		t.Errorf("Unmarshal failed\n%s\n", unm_err.Error())
	}
	t.Logf(fmt.Sprintln("{Unmarshalled Block}", data))

	if phStr != fmt.Sprintf("%x", data.Header.ParentHash) {
		t.Errorf("ParentHash has been altered by (un)marshaling process")
	}

	ResetTest()
}
*/
func ResetTest() {
	SaveBlockchain(blockchain_original)
	state_original.SaveState()
	snapshot_orignal.SaveSnapshot()
}

// Only run this to remake the local blockchain
// func TestCreateTestDatabase(t *testing.T) {
// 	state_block.SaveSnapshot()
// 	tx1 := state_block.CreateGenesisTransaction("Alberto", 100)
// 	err := state_block.AddTransaction(tx1)
// 	tx2 := state_block.CreateGenesisTransaction("Emilie", 5)
// 	err  = state_block.AddTransaction(tx2)
// 	tx3 := state_block.CreateGenesisTransaction("Niels", 1000000)
// 	err  = state_block.AddTransaction(tx3)
// 	tx4 := state_block.CreateGenesisTransaction("Asger", 420)
// 	err  = state_block.AddTransaction(tx4)
// 	tx5 := state_block.CreateGenesisTransaction("Magn", 69)
// 	err  = state_block.AddTransaction(tx5)
// 	tx6 := state_block.CreateTransaction("Niels", "Magn", 1000)
// 	err  = state_block.AddTransaction(tx6)
// 	tx7 := state_block.CreateTransaction("Magn", "Emilie", 12)
// 	err  = state_block.AddTransaction(tx7)
// 	tx8 := state_block.CreateTransaction("Emilie", "Asger", 3)
// 	err  = state_block.AddTransaction(tx8)
// 	tx9 := state_block.CreateTransaction("Emilie", "Magn", 2)
// 	err  = state_block.AddTransaction(tx9)
// 	tx10 := state_block.CreateTransaction("Emilie", "Niels", 2)
// 	err  = state_block.AddTransaction(tx10)
// 	tx11 := state_block.CreateReward("Emilie", 2)
// 	err  = state_block.AddTransaction(tx11)
// 	tx12 := state_block.CreateReward("Emilie", 2)
// 	err  = state_block.AddTransaction(tx12)
// 	tx13 := state_block.CreateTransaction("Magn", "Niels", 69)
// 	err  = state_block.AddTransaction(tx13)
// 	tx14 := state_block.CreateTransaction("Magn", "Niels", 69)
// 	err  = state_block.AddTransaction(tx14)
// 	tx15 := state_block.CreateTransaction("Magn", "Niels", 42)
// 	err  = state_block.AddTransaction(tx15)
// 	tx16 := state_block.CreateTransaction("Niels", "gggg", 1)
// 	err  = state_block.AddTransaction(tx16)
// 	tx17 := state_block.CreateReward("Alberto", 5000)
// 	err  = state_block.AddTransaction(tx17)
// 	tx18 := state_block.CreateTransaction("Niels", "Magn", 200000)
// 	err  = state_block.AddTransaction(tx18)
// 	tx19 := state_block.CreateTransaction("Magn", "Niels", 42)
// 	err  = state_block.AddTransaction(tx19)
// 	tx20 := state_block.CreateTransaction("Magn", "Niels", 89898)
// 	err  = state_block.AddTransaction(tx20)
// 	tx21 := state_block.CreateTransaction("Niels", "gggg", 1)
// 	err  = state_block.AddTransaction(tx21)
// 	tx22 := state_block.CreateReward("Alberto", 5000)
// 	err  = state_block.AddTransaction(tx22)
// 	tx23 := state_block.CreateTransaction("Niels", "Magn", 200000)
// 	err  = state_block.AddTransaction(tx23)
// 	tx24 := state_block.CreateTransaction("Magn", "Niels", 42)
// 	err  = state_block.AddTransaction(tx24)
// 	tx25 := state_block.CreateTransaction("Magn", "Niels", 89898)
// 	err  = state_block.AddTransaction(tx25)
// 	tx26 := state_block.CreateTransaction("Niels", "gggg", 1)
// 	err  = state_block.AddTransaction(tx26)
// 	tx27 := state_block.CreateReward("Alberto", 5000)
// 	err  = state_block.AddTransaction(tx27)
// 	tx28 := state_block.CreateTransaction("Niels", "Magn", 200000)
// 	err  = state_block.AddTransaction(tx28)

// 	block := state_block.CreateBlock(state_block.TxMempool)
// 	err   = state_block.AddBlock(block)

// 	tx29 := state_block.CreateTransaction("Niels", "Magn", 10)
// 	err  = state_block.AddTransaction(tx29)
// 	tx30 := state_block.CreateTransaction("Magn", "Emilie", 4)
// 	err  = state_block.AddTransaction(tx30)

// 	block = state_block.CreateBlock(state_block.TxMempool)
// 	err = state_block.AddBlock(block)

// 	tx31 := state_block.CreateTransaction("Niels", "Magn", 10)
// 	err  = state_block.AddTransaction(tx31)
// 	tx32 := state_block.CreateTransaction("Magn", "Emilie", 4)
// 	err  = state_block.AddTransaction(tx32)

// 	block = state_block.CreateBlock(state_block.TxMempool)
// 	err = state_block.AddBlock(block)

// 	if err != nil {
// 		fmt.Println("d")
// 	}
// 	fmt.Print("Uo")
// }

// // func TestByteSliceToHexString (t *testing.T)
