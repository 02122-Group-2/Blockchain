package database

import (
	"fmt"
	"math/rand"
	"testing"
)

var names = []AccountAddress{"Magn", "Niels", "Emilie", "Asger", "Alberto", "Bill", "Andrej"}

func makeDummyTransaction() Transaction {
	return state.CreateTransaction(names[rand.Int()*7], "niels", 6969.0)
}

var state_block, _ = LoadState()
var blockchain_original = LoadBlockchain()

func TestCreateBlock(t *testing.T) {
	tx1 := state_block.CreateTransaction("Niels", "Magn", 10)
	tx2 := state_block.CreateTransaction("Magn", "Emilie", 4)
	state_block.AddTransaction(tx1)
	state_block.AddTransaction(tx2)
	block := state_block.CreateBlock(state_block.TxMempool)
	fmt.Println(block)
}

func TestSaveBlock(t *testing.T) {
	blockchain_original = LoadBlockchain()

	// Create a block
	tx1 := state_block.CreateTransaction("Niels", "Magn", 10)
	tx2 := state_block.CreateTransaction("Magn", "Emilie", 4)
	state_block.AddTransaction(tx1)
	state_block.AddTransaction(tx2)
	block := state_block.CreateBlock(state_block.TxMempool)

	var blockList []Block

	blockList = append(blockList, block)

	SaveBlockchain(blockList)

	SaveBlockchain(blockchain_original)
}

func TestLoadBlockchain(t *testing.T) {
	res := LoadBlockchain()
	fmt.Println(res)
}

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
		t.Errorf("failed - expected zero errors and that the length of the TxMemPool is 1")
	}

	SaveBlockchain(blockchain_original) // Re-safe the original blockchain
}

// This tests makes sure the functionality of sharing the blocks work correctly.
// Two states will be creates, who are orignally identical.
// One state will create some transactions, Then create a block.
// The other will create a few transactions too. The first and last should be invalidated when the block from the first state when it is synced.
func TestSeperateStatesShareBlock(t *testing.T) {
	original_state := state_block.copyState()
	stateOne := state_block.copyState()
	stateTwo := state_block.copyState()

	stateOne.AddTransaction(stateOne.CreateTransaction("Magn", "Niels", 10))
	stateOne.AddTransaction(stateOne.CreateTransaction("Niels", "Magn", 10))
	stateOne.AddTransaction(stateOne.CreateTransaction("Magn", "Emilie", 10))

	blockOne := stateOne.CreateBlock(stateOne.TxMempool)

	stateTwo.AddTransaction(stateTwo.CreateTransaction("Magn", "Niels", 10))  // Should be invalid when merging the other block - Because of Nounces
	stateTwo.AddTransaction(stateTwo.CreateTransaction("Asger", "Emilie", 10)) // Should be valid
	stateTwo.AddTransaction(stateTwo.CreateTransaction("Niels", "Asger", 10)) // Should be invalid when merging the other block

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

	SaveBlockchain(blockchain_original)
}

func TestMarshalBlock(t *testing.T) {
	data, _ := blockchain_original[1].Header.MarshalJSON()
	t.Logf("%s", data)
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