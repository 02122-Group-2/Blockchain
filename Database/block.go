package database

import (
	shared "blockchain/Shared"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// * Magnus, s204509
type Block struct {
	Header   BlockHeader           `json:"Header"`
	SignedTx SignedTransactionList `json:"Transactions"`
}

// * Niels, s204503
type BlockHeader struct {
	ParentHash [32]byte `json:"ParentHash"`
	CreatedAt  int64    `json:"CreatedAt"`
	SerialNo   int      `json:"SerialNo"`
}

// * Niels, s204503
type BhDTO struct {
	ParentHash string `json:"ParentHash"`
	CreatedAt  int64  `json:"CreatedAt"` // make date
	SerialNo   int    `json:"SerialNo"`
}

// * Magnus, s204509
type Blockchain struct {
	Blockchain []Block `json:"Blockchain"`
}

// * Magnus, s204509
type Genesis struct {
	Balances map[AccountAddress]int `json:"balances"`
}

// * Asger, s204435
// Create a block object that matches the current state, given a list of transactions
func (state *State) CreateBlock(txs SignedTransactionList) Block {
	return Block{
		BlockHeader{
			state.getLatestHash(),
			shared.MakeTimestamp(),
			state.getNextBlockSerialNo(),
		},
		txs,
	}
}

// * Magnus, s204509
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

	if len(block.SignedTx) == 0 {
		return fmt.Errorf("the number of transactions must be greater than 0")
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

	return nil
}

// * Niels, s204503
// Takes a Block in JSON string format and calculates the 32-byte hash of this block and returns it.
func HashBlock(blockString string) [32]byte {
	return sha256.Sum256([]byte(blockString))
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

	state.LatestHash = HashBlock(jsonString)
	state.LastBlockSerialNo = block.Header.SerialNo
	state.LastBlockTimestamp = block.Header.CreatedAt
	state.TxMempool = nil
	return nil
}

// * Emilie, s204471
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

// * Magnus, s204509
// This functions takes a block and validates it against the state, then saves the block to the local blackchain.db file.
// It then applies the block to the state and saves a snapshot of the last "block"-state.
func (state *State) AddBlock(block Block) error {
	prevState := LoadSnapshot()

	err := prevState.ValidateBlock(block)
	if err != nil {
		return err
	}

	err = prevState.ApplyBlock(block)
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

	// Update the current state with the updates from the blockchain
	prevState.SaveState()

	// Updates the current state
	*state = prevState.copyState()
	return nil
}

// * Niels, s204503
// This updates the local blockchain.db file, by receiving a block and appending it to the list of blocks.
func PersistBlockToDB(block Block) error {
	oldBlocks := LoadBlockchain()
	oldBlocks = append(oldBlocks, block)

	if !SaveBlockchain(oldBlocks) {
		return fmt.Errorf("failed to save Blockchain locally")
	}

	return nil
}

// * Emilie, s204471
// Load the local blockchain and return it as a list of blocks
func LoadBlockchain() []Block {
	data, err := os.ReadFile(shared.LocatePersistenceFile("Blockchain.db", ""))
	if err != nil {
		return []Block{}
	}

	var loadedBlockchain Blockchain
	unm_err := json.Unmarshal(data, &loadedBlockchain)
	if unm_err != nil {
		return []Block{}
	}

	return loadedBlockchain.Blockchain
}

// * Asger, s204435
func ClearBlockchain() {
	err := os.Truncate(shared.LocatePersistenceFile("Blockchain.db", ""), 0)
	if err != nil {
		panic(err)
	}
}

// * Asger, s204435
// Given a list of blocks, save the list as the local blockchain.
func SaveBlockchain(blockchain []Block) bool {
	toSave := Blockchain{blockchain}
	txFile, _ := json.MarshalIndent(toSave, "", "  ")

	err := ioutil.WriteFile(shared.LocatePersistenceFile("Blockchain.db", ""), txFile, 0644)
	if err != nil {
		panic(err)
	}

	return true
}

// * Niels, s204503
// Load difference in contents of blockchain (for sending deltas to peer)
// TODO: data structure of blockchain not the most scalable... 🤔
func GetBlockChainDelta(blockchain []Block, fromBlockSerialNo int) []Block {
	// if 0, special case -> send entire blockchain
	if fromBlockSerialNo == 0 {
		return blockchain
	}

	startIdx := -1
	for i, b := range blockchain {
		if b.Header.SerialNo == fromBlockSerialNo {
			startIdx = i + 1
			break
		}
	}

	if startIdx == -1 || (len(blockchain)-1 < startIdx) {
		return nil
	}

	return blockchain[(startIdx):]
}

// * Niels, s204503
// Given a block, convert it to a JSON string
func BlockToJsonString(block Block) (string, error) {
	json, err := json.Marshal(block)
	if err != nil {
		return "", fmt.Errorf("Unable to convert block to a json string")
	}
	return string(json), nil
}

// * Niels, s204503
func (bh *BlockHeader) encodeBH() BhDTO {
	dto := BhDTO{}
	dto.ParentHash = fmt.Sprintf("%x", bh.ParentHash)
	dto.CreatedAt = bh.CreatedAt
	dto.SerialNo = bh.SerialNo

	return dto
}

// * Niels, s204503
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

// * Niels, s204503
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

// * Niels, s204503
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

// * Emilie, s204471
func (block *Block) BlockToString() string {

	listOfTransactions := ""
	for _, currTransaction := range block.SignedTx {
		listOfTransactions += TxToString(currTransaction.Tx) + "\n"
	}
	return "Header: \n " + "-Parent Hash: " + fmt.Sprintf("%v \n", block.Header.ParentHash) + "-Created at: " + fmt.Sprintf("%v \n", block.Header.CreatedAt) + "-Serial No.: " + fmt.Sprintf("%v \n", block.Header.SerialNo) + "List of Transactions: \n" + listOfTransactions
}

// * Niels, s204503
func GetLocalChainHashes(state State, fromSerialNo int) []string {
	blocks := LoadBlockchain()
	persistedChainHashes := getChainHashes(blocks, fromSerialNo)
	latestHash := fmt.Sprintf("%x", state.getLatestHash())
	return append(persistedChainHashes, latestHash)
}

// * Niels, s204503
func getChainHashes(blockchain []Block, fromSerialNo int) []string {
	var chainHashes []string
	for _, b := range blockchain {
		if b.Header.SerialNo > fromSerialNo {
			chainHashes = append(chainHashes, fmt.Sprintf("%x", b.Header.ParentHash))
		}
	}
	return chainHashes
}

// * Niels, s204503
// returns index of (first) mismatch, -1 if succesful
func CompareChainHashes(cHashes1 []string, cHashes2 []string) int {
	for i, hash := range cHashes1 {
		if i >= len(cHashes2) || hash != cHashes2[i] {
			return i
		}
	}
	return -1
}
