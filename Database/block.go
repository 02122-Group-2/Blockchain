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
	Header       BlockHeader            `json:"Header"`
	SignedTx    SignedTransactionList   `json:"Transactions"`
}

type BlockHeader struct {
	ParentHash [32]byte `json:"ParentHash"`
	CreatedAt  int64    `json:"CreatedAt"`
	SerialNo   int      `json:"SerialNo"`
}

type BhDTO struct {
	ParentHash string `json:"ParentHash"`
	CreatedAt  int64  `json:"CreatedAt"` // make date
	SerialNo   int    `json:"SerialNo"`
}

type Blockchain struct {
	Blockchain []Block `json:"Blockchain"`
}

// Create a block object that matches the current state, given a list of transactions
func (state *State) CreateBlock(txs SignedTransactionList) Block {
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
			return fmt.Errorf("the first block must have serial of 1")
		}
	}

	if block.Header.ParentHash != state.LatestHash {
		return fmt.Errorf("the parent hash doesn't match the hash of the Latest block \nBlock.Parent: %x\nState.Latest: %x", block.Header.ParentHash, state.LatestHash)
	}

	if block.Header.SerialNo != state.getNextBlockSerialNo() {
		return fmt.Errorf("block violates serial no. order")
	}

	if block.Header.CreatedAt <= state.LastBlockTimestamp {
		return fmt.Errorf("the new block must have a newer creation date than the Latest block")
	}

	err := state.ValidateTransactionList(block.SignedTx)
	if err != nil {
		return err
	}

	return nil
}

// Applies a single block to the current state.
// It validates all the transactions within the block.
// It applies all the transactions within the block to the state as well.
func (state *State) ApplyBlock(block Block) error {
	err := state.AddTransactionList(block.SignedTx)
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

	err = PersistBlockToDB(block)
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
func PersistBlockToDB(block Block) error {
	oldBlocks := LoadBlockchain()
	oldBlocks = append(oldBlocks, block)

	if !SaveBlockchain(oldBlocks) {
		return fmt.Errorf("failed to save Blockchain locally")
	}

	return nil
}

// Load the local blockchain and return it as a list of blocks
func LoadBlockchain() []Block {
	currWD, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(filepath.Join(currWD, "./Persistence/Blockchain.db"))
	if err != nil {
		panic(err)
	}

	var loadedBlockchain Blockchain
	unm_err := json.Unmarshal(data, &loadedBlockchain)
	if unm_err != nil {
		panic(unm_err)
	}

	return loadedBlockchain.Blockchain
}

// Given a list of blocks, save the list as the local blockchain.
func SaveBlockchain(blockchain []Block) bool {
	toSave := Blockchain{blockchain}
	txFile, _ := json.MarshalIndent(toSave, "", "  ")

	err := ioutil.WriteFile("./Persistence/Blockchain.db", txFile, 0644)
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

func (bh *BlockHeader) encodeBH() BhDTO {
	dto := BhDTO{}
	dto.ParentHash = fmt.Sprintf("%x", bh.ParentHash)
	dto.CreatedAt = bh.CreatedAt
	dto.SerialNo = bh.SerialNo

	return dto
}

func (dto *BhDTO) decodeBH() BlockHeader {
	bh := BlockHeader{}
	ph, _ := hex.DecodeString(dto.ParentHash)
	var ph32 [32]byte
	for i := 0; i < 32; i++ {
		ph32[i] = ph[i]
	}
	bh.ParentHash = ph32
	bh.CreatedAt = dto.CreatedAt
	bh.SerialNo = dto.SerialNo

	return bh
}

func (block *Block) MarshalJSON() ([]byte, error) {
	type Alias Block
	return json.Marshal(&struct {
		Header BhDTO `json:"Header"`
		*Alias
	}{
		Header: block.Header.encodeBH(),
		Alias:  (*Alias)(block),
	})
}

func (block *Block) UnmarshalJSON(data []byte) error {
	type Alias Block
	aux := &struct {
		Header BhDTO `json:"Header"`
		*Alias
	}{
		Alias: (*Alias)(block),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	block.Header = aux.Header.decodeBH()

	return nil
}
