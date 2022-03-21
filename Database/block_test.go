package database

import (
	Crypto "blockchain/Cryptography"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
)

// var state_block, _ = LoadState()
var state_block = LoadSnapshot2()
var blockchain_original = LoadBlockchain()

func TestCreateBlock(t *testing.T) {
	tx1 := state_block.CreateTransaction("Niels", "Asger", 10)
	tx2 := state_block.CreateTransaction("Asger", "Emilie", 4)
	state_block.AddTransaction(tx1)
	state_block.AddTransaction(tx2)
	block := state_block.CreateBlock(state_block.TxMempool)
	fmt.Println(block)
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

	SaveBlockchain(blockchain_original)
}

func TestLoadBlockchain(t *testing.T) {
	res := LoadBlockchain()
	fmt.Println(res)
}

func TestAddBlocksToBlockchain(t *testing.T) {
	block := state_block.CreateBlock(state_block.TxMempool)

	blJson, _ := BlockToJsonString(blockchain_original[0])
	ogHash := Crypto.HashBlock(blJson)
	fmt.Printf("%x\n", ogHash)

	err := state_block.AddBlock(block)
	if err != nil {
		panic(err)
	}

	tx1 := state_block.CreateTransaction("Niels", "Asger", 10)
	tx2 := state_block.CreateTransaction("Asger", "Emilie", 4)
	state_block.AddTransaction(tx1)
	state_block.AddTransaction(tx2)
	block2 := state_block.CreateBlock(TransactionList{tx1, tx2})

	err = state_block.AddBlock(block2)
	if err != nil {
		panic(err)
	}
	fmt.Println(err)
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
}
