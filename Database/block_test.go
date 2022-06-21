package database

import (
	Crypto "blockchain/Cryptography"
	shared "blockchain/Shared"

	// "encoding/hex"
	// "encoding/json"
	"fmt"
	"reflect"
	"testing"
)

// * file: Emilie, s204471

var blockchain_original = LoadBlockchain()
var state_original = LoadState()
var snapshot_orignal = LoadSnapshot()
var transactions_original = LoadTransactions()

func TestCreateBlock(t *testing.T) {
	// Prepare Files
	shared.ResetPersistenceFilesForTest()

	var state_block = LoadState()

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

	if len(block.SignedTx) != 2 {
		t.Logf("Expected number of transactions in block to be 2, but was %v", len(block.SignedTx))
		t.Fail()
	}

	trans1 := block.SignedTx[0]
	trans2 := block.SignedTx[1]

	if trans1.Tx.From != signedTx1.Tx.From && trans1.Tx.To != signedTx1.Tx.To && trans1.Tx.Amount != signedTx1.Tx.Amount {
		t.Log("Expected the first transaction in the block to be equal to the first transaction but wasn't")
		t.Fail()
	}

	if trans2.Tx.From != signedTx2.Tx.From && trans2.Tx.To != signedTx2.Tx.To && trans2.Tx.Amount != signedTx2.Tx.Amount {
		t.Log("Expected the second transaction in the block to be equal to the second transaction but wasn't")
		t.Fail()
	}

	ResetTest()
	testWallet.HardDelete()
}

func TestSaveBlock(t *testing.T) {
	t.Log("begin save block")
	shared.ResetPersistenceFilesForTest()

	var state_block = LoadState()

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

	//Now this blocklist should be equal to the one in the file
	loadedBlockchain := LoadBlockchain()

	if !reflect.DeepEqual(blockList, loadedBlockchain) {
		t.Log("Expected blockchains to be equal but weren't")
		t.Fail()
	}

	ResetTest()
	testWallet.HardDelete()
}

func TestLoadBlockchain(t *testing.T) {

	t.Log("begin load blockchain test")
	shared.ResetPersistenceFilesForTest()

	res := LoadBlockchain()
	fmt.Println(res)
	ResetTest()
}

func TestAddLegalBlockToBlockchain(t *testing.T) {

	t.Log("begin add legal block to blockchain")
	shared.ResetPersistenceFilesForTest()

	var state_block = LoadState()

	// Create both wallets
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet1, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state_block.AccountBalances[AccountAddress(testWallet1.Address)] = 1000

	Crypto.CreateNewWallet(walletUsername2, pswd)
	testWallet2, _ := Crypto.AccessWallet(walletUsername2, pswd)
	state_block.AccountBalances[AccountAddress(testWallet2.Address)] = 1000

	// Save Snapshot so it will accept the blockchain with the new balance
	state_block.SaveSnapshot()

	//Create the transactions first Transaction
	tx1, _ := state_block.CreateSignedTransaction(testWallet1, pswd, "Magn", 10)

	//Create the transactions Second Transaction
	tx2, _ := state_block.CreateSignedTransaction(testWallet2, pswd, "Emilie", 4)

	//Add transactions to the state
	err := state_block.AddTransaction(tx1)
	if err != nil {
		t.Error("Failed to add transaction 1")
	}

	err = state_block.AddTransaction(tx2)

	if err != nil {
		t.Error("Failed to add transaction 2")
	}

	//Create the block
	block1 := state_block.CreateBlock(state_block.TxMempool)

	//Add the block
	err = state_block.AddBlock(block1)
	if err != nil {
		t.Errorf("Expected block to be legal, but wasn't")
	}

	testWallet1.HardDelete()
	testWallet2.HardDelete()
	ResetTest()
}

func TestAddMultipleLegalBlockToBlockchain(t *testing.T) {

	t.Log("begin add legal block to blockchain")
	shared.ResetPersistenceFilesForTest()

	var state_block = LoadState()

	// Create both wallets
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet1, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state_block.AccountBalances[AccountAddress(testWallet1.Address)] = 1000

	Crypto.CreateNewWallet(walletUsername2, pswd)
	testWallet2, _ := Crypto.AccessWallet(walletUsername2, pswd)
	state_block.AccountBalances[AccountAddress(testWallet2.Address)] = 1000

	// Save Snapshot so it will accept the blockchain with the new balance
	state_block.SaveSnapshot()

	//Create the transactions first Transaction
	tx1, _ := state_block.CreateSignedTransaction(testWallet1, pswd, "Magn", 10)

	//Create the transactions Second Transaction
	tx2, _ := state_block.CreateSignedTransaction(testWallet2, pswd, "Emilie", 4)

	//Create the block 1
	block1 := state_block.CreateBlock(SignedTransactionList{tx1})

	//Add block 1
	err := state_block.AddBlock(block1)
	if err != nil {
		t.Errorf("Expected block 1 to be legal, but wasn't. Got Error: " + err.Error())
	}

	//Create the block 2
	block2 := state_block.CreateBlock(SignedTransactionList{tx2})

	//Add block 2
	err = state_block.AddBlock(block2)
	if err != nil {
		t.Errorf("Expected block 2 to be legal, but wasn't. Got Error: " + err.Error())
	}

	ResetTest()
}

func TestAddIllegalBlockWrongParentHash(t *testing.T) {

	t.Log("begin add illegal block to blockchain: Wrong parent hash")
	shared.ResetPersistenceFilesForTest()

	var state_block = LoadState()

	//Create the transactions first Transaction
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet1, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state_block.AccountBalances[AccountAddress(testWallet1.Address)] = 1000
	tx1, _ := state_block.CreateSignedTransaction(testWallet1, pswd, "Magn", 10)

	//Create the transactions Second Transaction
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet2, _ := Crypto.AccessWallet(walletUsername2, pswd)
	state_block.AccountBalances[AccountAddress(testWallet2.Address)] = 1000
	tx2, _ := state_block.CreateSignedTransaction(testWallet2, pswd, "Emilie", 4)

	//Add transactions to the state
	state_block.AddTransaction(tx1)
	state_block.AddTransaction(tx2)

	//Create the block
	block1 := state_block.CreateBlock(state_block.TxMempool)

	//Mess up parent hash
	block1.Header.ParentHash = [32]byte{}

	//Add the block
	err := state_block.AddBlock(block1)
	if err == nil {
		t.Errorf("Expected block to be illegal due to wrong parent hash, but wasn't")
	}

	ResetTest()

}

func TestAddIllegalBlockIllegalTransaction(t *testing.T) {

	t.Log("begin add illegal block to blockchain: Illegal transaction")
	shared.ResetPersistenceFilesForTest()

	var state_block = LoadState()

	//Create the transactions first Transaction
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet1, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state_block.AccountBalances[AccountAddress(testWallet1.Address)] = 1000
	tx1, _ := state_block.CreateSignedTransaction(testWallet1, pswd, "Magn", 10)

	//Create the transactions Second Transaction
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet2, _ := Crypto.AccessWallet(walletUsername2, pswd)
	state_block.AccountBalances[AccountAddress(testWallet2.Address)] = 1000
	tx2, _ := state_block.CreateSignedTransaction(testWallet2, pswd, "Emilie", 2000)

	//Create the block from a manually created transaction list
	block1 := state_block.CreateBlock(SignedTransactionList{tx1, tx2})

	//Add the block
	err := state_block.AddBlock(block1)
	if err == nil {
		t.Log("Expected block to be illegal due to illegal transaction, but wasn't", err)
		t.Fail()
	}

	ResetTest()

}

func TestAddIllegalBlockWrongTimestamp(t *testing.T) {

	t.Log("begin add illegal block to blockchain: Wrong timestamp")
	shared.ResetPersistenceFilesForTest()

	var state_block = LoadState()

	//Create the transactions first Transaction
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet1, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state_block.AccountBalances[AccountAddress(testWallet1.Address)] = 1000
	tx1, _ := state_block.CreateSignedTransaction(testWallet1, pswd, "Magn", 10)

	//Create the transactions Second Transaction
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet2, _ := Crypto.AccessWallet(walletUsername2, pswd)
	state_block.AccountBalances[AccountAddress(testWallet2.Address)] = 1000
	tx2, _ := state_block.CreateSignedTransaction(testWallet2, pswd, "Emilie", 4)

	//Add transactions to the state
	state_block.AddTransaction(tx1)
	state_block.AddTransaction(tx2)

	//Create the block
	block1 := state_block.CreateBlock(state_block.TxMempool)

	//Mess up the timestamp
	block1.Header.CreatedAt = state_block.LastBlockTimestamp - 1

	//Add the block
	err := state_block.AddBlock(block1)
	if err == nil {
		t.Errorf("Expected block to be illegal due to wrong timestamp, but wasn't")

	}

	ResetTest()

}

func TestAddIllegalBlockWrongBlockHeigh(t *testing.T) {

	t.Log("begin add illegal block to blockchain: Wrong block height")
	shared.ResetPersistenceFilesForTest()

	var state_block = LoadState()

	//Create the transactions first Transaction
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet1, _ := Crypto.AccessWallet(walletUsername1, pswd)
	state_block.AccountBalances[AccountAddress(testWallet1.Address)] = 1000
	tx1, _ := state_block.CreateSignedTransaction(testWallet1, pswd, "Magn", 10)

	//Create the transactions Second Transaction
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet2, _ := Crypto.AccessWallet(walletUsername2, pswd)
	state_block.AccountBalances[AccountAddress(testWallet2.Address)] = 1000
	tx2, _ := state_block.CreateSignedTransaction(testWallet2, pswd, "Emilie", 4)

	//Add transactions to the state
	state_block.AddTransaction(tx1)
	state_block.AddTransaction(tx2)

	//Create the block
	block1 := state_block.CreateBlock(state_block.TxMempool)

	//Mess up the height
	block1.Header.SerialNo -= 1

	//Add the block
	err := state_block.AddBlock(block1)
	if err == nil {
		t.Errorf("Expected block to be illegal due to wrong blockheight, but wasn't")

	}

	ResetTest()

}

func TestAddIllegalBlockNoTransactions(t *testing.T) {

	t.Log("begin add illegal block to blockchain: No transactions")
	shared.ResetPersistenceFilesForTest()

	var state_block = LoadState()

	//Create the block with no transactions
	block1 := state_block.CreateBlock(SignedTransactionList{})

	//Add the block
	err := state_block.AddBlock(block1)
	if err == nil {
		t.Errorf("Expected block to be illegal due to no transactions, but wasn't")
	}

	ResetTest()

}

func TestAddBlockWhereSomeTransactionsFromStateAreInvalidatedAfterBlock(t *testing.T) {
	t.Log("begin Add Block Where Some Transactions From State Are Invalidated After Block has been added test")

	shared.ResetPersistenceFilesForTest()

	original_state := LoadSnapshot()

	//State one will be the local state with two transactions
	stateOne := original_state.copyState()

	// Create the two wallets and save the state snapshot
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet1, _ := Crypto.AccessWallet(walletUsername1, pswd)
	stateOne.AccountBalances[AccountAddress(testWallet1.Address)] = 1000

	Crypto.CreateNewWallet(walletUsername2, pswd)
	testWallet2, _ := Crypto.AccessWallet(walletUsername2, pswd)
	stateOne.AccountBalances[AccountAddress(testWallet2.Address)] = 1000

	stateOne.SaveSnapshot()

	//Create the transactions
	tx1, _ := stateOne.CreateSignedTransaction(testWallet1, pswd, "Niels", 10)
	tx2, _ := stateOne.CreateSignedTransaction(testWallet2, pswd, "Magn", 67) //this will be invalid after block has been added
	stateOne.AddTransaction(tx1)
	stateOne.AddTransaction(tx2)

	//State two will be the "peer" state that contains the block
	stateTwo := original_state.copyState()
	stateTwo.AccountBalances[AccountAddress(testWallet2.Address)] = 1000
	tx3, _ := stateTwo.CreateSignedTransaction(testWallet2, pswd, "Magn", 10) //Will be valid and cause transaction in state one to be invalid
	stateTwo.AddTransaction(tx3)

	// Creates the block
	block := stateTwo.CreateBlock(stateTwo.TxMempool)

	//Now this block will be added to state one causing one of the transactions to become invalid
	err := stateOne.AddBlock(block)
	if err != nil {
		SaveBlockchain(blockchain_original)
		t.Errorf("failed to add block to first state...")
	}

	//At this point state one should only have one legal transaction and the other should have been invalidated
	if len(stateOne.TxMempool) != 1 {
		t.Logf("Expected number of transactions to be one, but was %v", len(stateOne.TxMempool))
		t.Fail()
	}
	ResetTest()
}

//The local state is kept empty. After a block has been feteched it will be added to this local state. The changes from this block should be applied to the local state
func TestAddBlockWhereSomeTransactionsAreNotInCurrentState(t *testing.T) {
	t.Log("begin add block where some transactions are not in current state test")

	shared.ResetPersistenceFilesForTest()

	original_state := LoadSnapshot()
	//State one will be the local state
	stateOne := original_state.copyState()

	//State two will be the "peer" state that contains the block
	stateTwo := original_state.copyState()

	// Create and add the transactions
	Crypto.CreateNewWallet(walletUsername1, pswd)
	testWallet1, _ := Crypto.AccessWallet(walletUsername1, pswd)
	stateOne.AccountBalances[AccountAddress(testWallet1.Address)] = 1000
	stateTwo.AccountBalances[AccountAddress(testWallet1.Address)] = 1000

	// Create the transactions
	tx1, _ := stateTwo.CreateSignedTransaction(testWallet1, pswd, "Niels", 10)
	tx2, _ := stateTwo.CreateSignedTransaction(testWallet1, pswd, "Magn", 10)
	stateTwo.AddTransaction(tx1)
	stateTwo.AddTransaction(tx2)

	// Create the block
	block := stateTwo.CreateBlock(stateTwo.TxMempool)

	// Prepare Snapshot
	original_state.AccountBalances[AccountAddress(testWallet1.Address)] = 1000
	original_state.SaveSnapshot()

	//Now this block will be added to state one
	err := stateOne.AddBlock(block)
	if err != nil {
		SaveBlockchain(blockchain_original)
		t.Errorf("failed to add block to first state...")
	}

	// Saves the snapshot, since the snapshot is still "outdated" for the other account. This error is due to the fact that we run the software on the same pc.
	original_state.SaveSnapshot()

	//Now this block will be added to state two
	err = stateTwo.AddBlock(block)
	if err != nil {
		SaveBlockchain(blockchain_original)
		t.Errorf("failed to add block to second state...")
	}
	//At this point state one and state two should be identical in therms of balances and block height + latest hash
	if stateOne.getLatestHash() != stateTwo.getLatestHash() && stateOne.LastBlockSerialNo != stateTwo.LastBlockSerialNo && reflect.DeepEqual(stateOne.AccountBalances, stateTwo.AccountBalances) {
		t.Log("Expected the hashes, nounces and balances to be equal, but they aren't")
		t.Fail()
	}

	ResetTest()

}

/*
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
	testWallet.HardDelete()
}
*/

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

	testWallet1.HardDelete()
	testWallet2.HardDelete()
	testWallet3.HardDelete()
	ResetTest()
}

/*
func TestMarshalUnmarshalBlock(t *testing.T) {
	t.Log("begin marshal unmarshal block test")

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
// 	// Genesis transaction to Emilie
// 	tx1 := state_block.CreateGenesisTransaction("0xC98f5180FF9836CC2EF67158EfEd9AA5ddeC54F6", 100000000)
// 	_    = state_block.AddTransaction(tx1)

// 	// Genesis transaction to Asger
// 	tx2 := state_block.CreateGenesisTransaction("0x5b355Cd0C7fB6aD65b2e9342Fb6FBf0146585D7b", 100000000)
// 	_    = state_block.AddTransaction(tx2)

// 	// Genesis transaction to Magn
// 	tx3 := state_block.CreateGenesisTransaction("0x5D34001173D5d05fA3AC865fb2b30131478a13d7", 100000000)
// 	_    = state_block.AddTransaction(tx3)

// 	block := state_block.CreateBlock(state_block.TxMempool)
// 	err   := state_block.AddBlock(block)

// 	if err != nil {
// 		fmt.Println("d")
// 	}
// 	fmt.Print("Uo")
// }

// // func TestByteSliceToHexString (t *testing.T)
