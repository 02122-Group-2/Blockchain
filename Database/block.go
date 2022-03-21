package database

import (
	Crypto "blockchain/Cryptography"
	"encoding/hex"
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
	ParentHash [32]byte `json: "ParentHash"`
	CreatedAt  int64    `json: "CreatedAt"`
	SerialNo   int      `json: "SerialNo"`
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
	if state.LastBlockSerialNo == 0 { // If no other block is added, add the block if the block has serialNo. 1
		if block.Header.SerialNo == 1 {
			return nil
		} else {
			return fmt.Errorf("The first block must have serial of 1")
		}
	}

	if block.Header.ParentHash != state.LatestHash {
		return fmt.Errorf("The parent hash doesn't match the hash of the Latest block")
	}

	if block.Header.SerialNo != state.getNextBlockSerialNo() {
		return fmt.Errorf("Block violates serial no. order")
	}

	if block.Header.CreatedAt <= state.LastBlockTimestamp {
		return fmt.Errorf("The new block must have a newer creation date than the Latest block")
	}

	err := state.ValidateTransactionList(block.Transactions)
	if err != nil {
		return err
	}

	return nil
}

func (state *State) ApplyBlocks(blocks []Block) error {
	for _, t := range blocks {
		validation_err := state.ValidateBlock(t)
		if validation_err != nil {
			return validation_err
		}
		if err := state.ApplyBlock(t); err != nil {
			return fmt.Errorf("Block failed: " + err.Error())
		}
	}
	return nil
}

func (state *State) ApplyBlock(block Block) error {
	err := state.AddTransactionList(block.Transactions)
	if err != nil {
		return err
	}

	jsonString, jsonErr := BlockToJsonString(block)
	if jsonErr != nil {
		return jsonErr
	}

	state.LatestHash = Crypto.HashBlock(jsonString)
	state.LastBlockSerialNo = block.Header.SerialNo
	state.LastBlockTimestamp = block.Header.CreatedAt
	state.TxMempool = nil
	return nil
}

func (state *State) AddBlock(block Block) error {
	prevState := LoadSnapshot()
	err := prevState.ValidateBlock(block)
	if err != nil {
		return err
	}

	err = state.PersistBlockToDB(block)
	if err != nil {
		return err
	}

	// This functionality is not working properly yet. Need a better system of applying eiher blocks or transactions. Both will result in applying transactions twice.
	err = state.ApplyBlock(block)
	if err != nil {
		return err
	}

	state.SaveSnapshot()

	return nil
}

func (state *State) PersistBlockToDB(block Block) error {
	oldBlocks := LoadBlockchain()
	oldBlocks = append(oldBlocks, block)

	if !SaveBlockchain(oldBlocks) {
		return fmt.Errorf("Failed to save Blockchain locally")
	}

	return nil
}

func BlockToJsonString(block Block) (string, error) {
	json, err := json.Marshal(block)
	if err != nil {
		return "", fmt.Errorf("Unable to convert block to a json string")
	}
	return string(json), nil
}

func (bh *BlockHeader) MarshalJSON() ([]byte, error) {
	fmt.Println("123bh.MarshalJSON() called!!")
	type BhAlias BlockHeader
	return json.Marshal(&struct {
		ParentHash string `json: "ParentHash"`
		*BhAlias
	}{
		ParentHash: fmt.Sprintf("%x", bh.ParentHash),
		BhAlias:    (*BhAlias)(bh),
	})
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

func (bh *BlockHeader) UnmarshalJSON(data []byte) error {
	fmt.Println("bh.UnmarshalJSON() called!!")
	type BhAlias BlockHeader
	aux := &struct {
		ParentHash string `json: "ParentHash"`
		*BhAlias
	}{
		BhAlias: (*BhAlias)(bh),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	byte_arr, decode_err := hex.DecodeString(aux.ParentHash)
	if decode_err != nil {
		panic(decode_err)
	}

	for i := 0; i < 32; i++ {
		bh.ParentHash[i] = byte_arr[i]
	}

	return nil
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

func LoadSnapshot() State {
	currWD, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(filepath.Join(currWD, "LatestSnapshot.json"))
	if err != nil {
		panic(err)
	}

	var state State
	json.Unmarshal(data, &state)

	return state
}

func (state *State) SaveSnapshot() bool {
	txFile, _ := json.MarshalIndent(state, "", "  ")

	err := ioutil.WriteFile("./LatestSnapshot.json", txFile, 0644)
	if err != nil {
		panic(err)
	}

	return true
}

func (currState *State) copyState() State {
	copy := State{}

	copy.TxMempool = make([]Transaction, len(currState.TxMempool))
	copy.Balances = make(map[AccountAddress]uint)

	copy.LastBlockSerialNo = currState.LastBlockSerialNo
	copy.LastBlockTimestamp = currState.LastBlockTimestamp
	copy.LatestHash = currState.LatestHash
	copy.LatestTimestamp = currState.LatestTimestamp

	for accountA, balance := range currState.Balances {
		copy.Balances[accountA] = balance
	}

	for _, tx := range currState.TxMempool {
		copy.TxMempool = append(copy.TxMempool, tx)
	}

	return copy
}
