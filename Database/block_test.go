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
	res := LoadBlockchain();
	fmt.Println(res)
}


func TestAddBlockToBlockchain(t *testing.T) {
	block := state_block.CreateBlock(state_block.TxMempool)
	err := state_block.AddBlock(block)
	fmt.Println(err)

	tx1 := state_block.CreateTransaction("Niels", "Magn", 10)
	tx2 := state_block.CreateTransaction("Magn", "Emilie", 4)
	// state_block.AddTransaction(tx1)
	// state_block.AddTransaction(tx2)
	block2 := state_block.CreateBlock(TransactionList{tx1,tx2})

	err = state_block.AddBlock(block2)
	fmt.Println(err)
}
