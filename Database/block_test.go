package database

import (
	shared "blockchain/Shared"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
)

var state_block = LoadState()
var blockchain_original = LoadBlockchain()
var state_original = LoadState()
var snapshot_orignal = LoadSnapshot()
var transactions_original = LoadTransactions()

func TestCreateBlock(t *testing.T) {
	shared.ResetPersistenceFilesForTest()

	tx1 := state_block.CreateTransaction("Niels", "Asger", 10)
	tx2 := state_block.CreateTransaction("Asger", "Emilie", 4)
	state_block.AddTransaction(tx1)
	state_block.AddTransaction(tx2)
	block := state_block.CreateBlock(state_block.TxMempool)
	fmt.Println(block)

	ResetTest()
}

func TestSaveBlock(t *testing.T) {
	blockchain_original = LoadBlockchain()

	// Create a block
	tx1 := state_block.CreateTransaction("Niels", "Asger", 10)
	tx2 := state_block.CreateTransaction("Asger", "Emilie", 4)
	state_block.AddTransaction(tx1)
	state_block.AddTransaction(tx2)
	block := state_block.CreateBlock(state_block.TxMempool)

	// var blockList []Block

	blockList := append(blockchain_original, block)

	SaveBlockchain(blockList)

	ResetTest()
}

func TestLoadBlockchain(t *testing.T) {
	res := LoadBlockchain()
	fmt.Println(res)
	ResetTest()
}

func TestAddLegalBlockToBlockchain(t *testing.T) {
	shared.ResetPersistenceFilesForTest()

	//Create the transactions
	tx1 := state_block.CreateTransaction("Niels", "Magn", 10)
	tx2 := state_block.CreateTransaction("Magn", "Emilie", 4)

	//Add transactions to the state
	state_block.AddTransaction(tx1)
	state_block.AddTransaction(tx2)

	//Create the block
	block1 := state_block.CreateBlock(state_block.TxMempool)

	//Add the block
	err := state_block.AddBlock(block1)
	if err != nil {
		t.Errorf("Expected block to be legal, but wasn't")
	}

	ResetTest()
}

func TestAddIllegalBlockWrongParentHash(t *testing.T) {
	shared.ResetPersistenceFilesForTest()

	//Create the transactions
	tx1 := state_block.CreateTransaction("Niels", "Magn", 10)
	tx2 := state_block.CreateTransaction("Magn", "Emilie", 4)

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
	shared.ResetPersistenceFilesForTest()

	//Create the transactions
	tx1 := state_block.CreateTransaction("Niels", "Magn", 10)

	//Illegal transaction
	tx2 := state_block.CreateTransaction("Magn", "Emilie", 10000000)

	//Create the block from a manually created transaction list
	block1 := state_block.CreateBlock(TransactionList{tx1, tx2})

	//Add the block
	err := state_block.AddBlock(block1)
	if err == nil {
		t.Log("Expected block to be illegal due to illegal transaction, but wasn't", err)
		t.Fail()
	}

	ResetTest()

}

func TestAddIllegalBlockWrongTimestamp(t *testing.T) {
	shared.ResetPersistenceFilesForTest()

	//Create the transactions
	tx1 := state_block.CreateTransaction("Niels", "Magn", 10)
	tx2 := state_block.CreateTransaction("Magn", "Emilie", 4)

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
	shared.ResetPersistenceFilesForTest()

	//Create the transactions
	tx1 := state_block.CreateTransaction("Niels", "Magn", 10)
	tx2 := state_block.CreateTransaction("Magn", "Emilie", 4)

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
	shared.ResetPersistenceFilesForTest()

	//Create the block with no transactions
	block1 := state_block.CreateBlock(TransactionList{})

	//Add the block
	err := state_block.AddBlock(block1)
	if err == nil {
		t.Errorf("Expected block to be illegal due to no transactions, but wasn't")
	}

	ResetTest()

}

func TestAddBlockWhereSomeTransactionsFromStateAreInvalidatedAfterBlock(t *testing.T) {
	shared.ResetPersistenceFilesForTest()

	//Create the transaction
	tx1 := state_block.CreateTransaction("Magn", "Emilie", 4)

	//Add transactions to the state
	state_block.AddTransaction(tx1)

	// transaction not in the local mempool
	tx2 := Transaction{"Niels", "Magn", 10, 10, makeTimestamp(), "transaction"}

	//Create the block with the transactions (one local and one unknown)
	block1 := state_block.CreateBlock(TransactionList{tx1, tx2})

	//Create local tranaction which would be legal before block, but is illegal now
	tx3 := state_block.CreateTransaction("Niels", "Asger", 579030) //Current balance for Niels is 579029
	state_block.AddTransaction(tx3)

	//Add the block
	//At this point tx3 will be invalidated and not added to the state and will be lost
	state_block.AddBlock(block1)

	if len(state_block.TxMempool) != 0 {
		t.Errorf("Expected TxMempool to be empty, but isn't")
	}

	ResetTest()

}

/*
func TestAddBlockToBlockchain(t *testing.T) {
	tx1 := state_block.CreateTransaction("Niels", "Magn", 10)
	tx2 := state_block.CreateTransaction("Magn", "Emilie", 4)
	state_block.AddTransaction(tx1)
	state_block.AddTransaction(tx2)
	block2 := state_block.CreateBlock(state_block.TxMempool)

	tx3 := state_block.CreateTransaction("Magn", "Emilie", 22)
	state_block.AddTransaction(tx3)

	err := state_block.AddBlock(block2)
	if err != nil || len(state_block.TxMempool) != 1 {
		t.Errorf("failed - expected no errors and that the length of the TxMemPool is 1")
	}

	ResetTest()
}
*/

// This tests makes sure the functionality of sharing the blocks work correctly.
// Two states will be created, who are orignally identical.
// One state will create some transactions, Then create a block.
// The other will create a few transactions too. The first and last should be invalidated when the block from the first state when it is synced.
func TestSeperateStatesShareBlock(t *testing.T) {
	original_state := LoadSnapshot()
	stateOne := original_state.copyState()
	stateTwo := original_state.copyState()

	stateOne.AddTransaction(stateOne.CreateTransaction("Magn", "Niels", 10))
	stateOne.AddTransaction(stateOne.CreateTransaction("Niels", "Magn", 10))
	stateOne.AddTransaction(stateOne.CreateTransaction("Magn", "Emilie", 10))

	blockOne := stateOne.CreateBlock(stateOne.TxMempool)

	stateTwo.AddTransaction(stateTwo.CreateTransaction("Magn", "Niels", 10))   // Should be invalid when merging the other block - Because of Nounces
	stateTwo.AddTransaction(stateTwo.CreateTransaction("Asger", "Emilie", 10)) // Should be valid
	stateTwo.AddTransaction(stateTwo.CreateTransaction("Niels", "Asger", 10))  // Should be invalid when merging the other block

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

	ResetTest()
}

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

func ResetTest() {
	SaveBlockchain(blockchain_original)
	state_original.SaveState()
	snapshot_orignal.SaveSnapshot()
	transactions_original.SaveTransactions()
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
