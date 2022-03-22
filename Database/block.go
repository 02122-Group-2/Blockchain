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

// Create a block object that matches the current state, given a list of transactions
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

// Validates a given block against the current state
// It checks: The parent hash, Serial No., Timestamp, and the validity of the transactions within the block.
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


// Applies a single block to the current state.
// It validates the block and all the transactions within.
// It applies all the transactions within the block to the state as well.
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


// Applies a list of blocks to the current state. Given a list of block (blockchain) it will apply each block to the state. 
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

// This functions takes a block and validates it against the state, then saves the block to the local blackchain.db file.
// It then applies the block to the state and saves a snapshot of the last "block"-state.
func (state *State) AddBlock(block Block) error {
	prevState := LoadSnapshot()
	err := prevState.ApplyBlock(block)
	if err != nil {
		return err
	}

	err = prevState.PersistBlockToDB(block)
	if err != nil {
		return err
	}

	// Save the new newest block state
	prevState.SaveSnapshot()

	// Apply all the remaining transactions from the current memory pool
	prevState.TryAddTransactions(state.TxMempool)

	// Updates the current state
	*state = prevState.copyState()
	return nil
}

// This updates the local blockchain.db file, by receiving a block and appending it to the list of blocks.
func (state *State) PersistBlockToDB(block Block) error {
	oldBlocks := LoadBlockchain()
	oldBlocks = append(oldBlocks, block)

	if !SaveBlockchain(oldBlocks) {
		return fmt.Errorf("Failed to save Blockchain locally")
	}

	return nil
}

// Load the local blockchain and return it as a list of blocks 
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

// Given a list of blocks, save the list as the local blockchain.
func SaveBlockchain(blockchain []Block) bool {
	toSave := Blockchain{blockchain}
	txFile, _ := json.MarshalIndent(toSave, "", "  ")

	err := ioutil.WriteFile("./Blockchain.db", txFile, 0644)
	if err != nil {
		panic(err)
	}

	return true
}

// Given a block, convert it to a JSON string
func BlockToJsonString(block Block) (string, error) {
	json, err := json.Marshal(block)
	if err != nil {
		return "", fmt.Errorf("Unable to convert block to a json string")
	}
	return string(json), nil
}

// Given a blockheader, convert it into a JSON string object. Performs sepcial formatting on the parent hash.
func (bh *BlockHeader) MarshalJSON() ([]byte, error) {
	type BhAlias BlockHeader
	return json.Marshal(&struct {
		ParentHash string `json: "ParentHash"`
		*BhAlias
	}{
		ParentHash: fmt.Sprintf("%x", bh.ParentHash),
		BhAlias:    (*BhAlias)(bh),
	})
}

// Given a blockheader and an array of bytes, use the bytes to create a header with a parent hash of the correct formatting.
func (bh *BlockHeader) UnmarshalJSON(data []byte) error {
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




// Loads the latest snapchat of the state. Each snapshat is meant as the state right after a block has been added.
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

// Given a state, save the state as the local state snapshot.
func (state *State) SaveSnapshot() bool {
	txFile, _ := json.MarshalIndent(state, "", "  ")

	err := ioutil.WriteFile("./LatestSnapshot.json", txFile, 0644)
	if err != nil {
		panic(err)
	}

	return true
}

// Given a state, make a deep copy of the state and return the copy.
func (currState *State) copyState() State {
	copy := State{}

	copy.TxMempool = make([]Transaction, 0)
	copy.AccountBalances = make(map[AccountAddress]uint)
	copy.AccountNounces  = make(map[AccountAddress]uint)

	copy.LastBlockSerialNo = currState.LastBlockSerialNo
	copy.LastBlockTimestamp = currState.LastBlockTimestamp
	copy.LatestHash = currState.LatestHash
	copy.LatestTimestamp = currState.LatestTimestamp

	for accountA, balance := range currState.AccountBalances {
		copy.AccountBalances[accountA] = balance
	}

	for accountA, nounce := range currState.AccountNounces {
		copy.AccountNounces[accountA] = nounce
	}

	for _, tx := range currState.TxMempool {
		copy.TxMempool = append(copy.TxMempool, tx)
	}

	return copy
}
