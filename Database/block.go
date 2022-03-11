package database

import (
	Crypto "blockchain/Cryptography"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)


type Block struct {
	Header       BlockHeader   `json: "Header"`
	Transactions []Transaction `json: "Transactions"`
}

type BlockHeader struct {
	ParentHash string `json: "ParentHash"`
	CreatedAt  int64  `json: "CreatedAt"`
	SerialNo   int    `json: "SerialNo"`
}

type Blockchain struct {
	Blockchain []Block `json: "Blockchain"`
}

func (state *State) CreateBlock(txs []Transaction) Block {
	return Block{
		BlockHeader{
			state.getLatestHash(),
			makeTimestamp(),
			state.getNextBlockSerialNo(),
		},
		txs,
	}
}

func (state *State) ValidateBlock(block Block) error {
	if state.lastBlockSerialNo == 0 { // If no other block is added, add the block if the block has serialNo. 1
		if block.Header.SerialNo == 1 {
			return nil
		} else {
			return fmt.Errorf("The first block must have serial of 1")
		}
	}

	if block.Header.ParentHash != state.latestHash {
		return fmt.Errorf("The parent hash doesn't match the hash of the latest block")
	}

	if block.Header.SerialNo != state.getNextBlockSerialNo() {
		return fmt.Errorf("Block violates serial no. order")
	}

	if block.Header.CreatedAt <= state.latestTimestamp {
		return fmt.Errorf("The new block must have a newer creation date than the latest block")
	}

	err := state.ValidateTransactionList(block.Transactions) 
	if err != nil {
		return err
	}

	return nil
}

func (state *State) AddBlock(block Block) error {
	err := state.ValidateBlock(block) 
	if err != nil {
		return err
	}

	err = state.PersistBlockToDB(block)
	if err != nil {
		return err
	}

	jsonString, jsonErr := BlockToJsonString(block)
	if jsonErr != nil {
		return jsonErr
	}

	state.latestHash = Crypto.HashBlock(jsonString)
	state.lastBlockSerialNo = block.Header.SerialNo
	state.latestTimestamp = block.Header.CreatedAt
	state.TxMempool = nil


	return nil
}

func (state *State) PersistBlockToDB(block Block) error {
	err := state.ValidateBlock(block) 
	if err != nil {
		return err
	}

	oldBlocks := LoadBlockchain()
	oldBlocks = append(oldBlocks, block)
	
	if !SaveBlockchain(oldBlocks) {
		return fmt.Errorf("Failed to save Blockchain locally")
	}

	return nil
}

func BlockToJsonString(block Block) (string, error) {
	json, err  := json.Marshal(block)
	if err != nil {
		return "", fmt.Errorf("Unable to convert block to a json string")
	}
	return string(json), nil
}

func LoadBlockchain() []Block {
	currWD, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(filepath.Join(currWD, "Blockchain.db"))
	if err != nil {
		panic(err)
	}

	var loadedBlockchain Blockchain
	json.Unmarshal(data, &loadedBlockchain)

	return loadedBlockchain.Blockchain
}

func SaveBlockchain(blockchain []Block) bool {
	toSave := Blockchain{blockchain}
	txFile, _ := json.MarshalIndent(toSave, "", "  ")

	err := ioutil.WriteFile("./Blockchain.db", txFile, 0644)
	if err != nil {
		panic(err)
	}

	return true
}