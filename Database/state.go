package database

import (
	Crypto "blockchain/Cryptography"
	shared "blockchain/Shared"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	// "path/filepath"
)

// * Magnus, s204509
type StateFromPostRequest struct {
	AccountBalances    map[AccountAddress]uint `json:"AccountBalances"`
	AccountNounces     map[AccountAddress]uint `json:"AccountNounces"`
	TxMempool          TransactionList         `json:"TxMempool"`
	DbFile             *os.File                `json:"DbFile"`
	LastBlockSerialNo  int                     `json:"LastBlockSerialNo"`
	LastBlockTimestamp int64                   `json:"LastBlockTimestamp"`
	LatestHash         []byte                  `json:"LatestHash"`
	LatestTimestamp    int64                   `json:"LatestTimestamp"`
}

// * Emilie, s204471
type State struct {
	AccountBalances    map[AccountAddress]uint `json: "AccountBalances"`
	AccountNounces     map[AccountAddress]uint `json: "AccountNounces"`
	TxMempool          SignedTransactionList   `json: "TxMempool"`
	DbFile             *os.File                `json: "DbFile"`
	LastBlockSerialNo  int                     `json: "LastBlockSerialNo"`
	LastBlockTimestamp int64                   `json: "LastBlockTimestamp"`
	LatestHash         [32]byte                `json: "LatestHash"`
	LatestTimestamp    int64                   `json: "LatestTimestamp"`
}

// * Niels, s204503
func (s *State) getNextBlockSerialNo() int {
	return s.LastBlockSerialNo + 1
}

// * Niels, s204503
func (s *State) getLatestHash() [32]byte {
	return s.LatestHash
}

// * Niels, s204503
func (s *State) MarshalJSON() ([]byte, error) {
	type Alias State
	return json.Marshal(&struct {
		LatestHash string `json:"LatestHash"`
		*Alias
	}{
		LatestHash: fmt.Sprintf("%x", s.LatestHash),
		Alias:      (*Alias)(s),
	})
}

// * Niels, s204503
func (s *State) UnmarshalJSON(data []byte) error {
	type Alias State
	aux := &struct {
		LatestHash string `json:"LatestHash"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	lh, _ := hex.DecodeString(aux.LatestHash)
	var lh32 [32]byte
	for i := 0; i < 32; i++ {
		lh32[i] = lh[i]
	}
	s.LatestHash = lh32

	return nil
}

// * Magnus, s204509
// Creates a state based from the data in the local blockchain.db file.
func LoadState() *State {
	state := loadStateFromJSON("CurrentState.json")
	return &state
}

// * Asger, s204435
func (state *State) ClearState() {
	state.LastBlockSerialNo = 0
	err := os.Truncate(shared.LocatePersistenceFile("CurrentState.json", ""), 0)
	if err != nil {
		panic(err)
	}

	err = os.Truncate(shared.LocatePersistenceFile("CurrentState.json", ""), 0)
	if err != nil {
		panic(err)
	}

}

// * Magnus, s204509
// Adds a transaction to the state. It will validate the transaction, then apply the transaction to the state,
// then add the transaction to its MemPool and update its latest timestamp field.
func (state *State) AddTransaction(transaction SignedTransaction) error {
	if err := state.ValidateTransaction(transaction); err != nil {
		return err
	}

	state.TxMempool = append(state.TxMempool, transaction)

	state.ApplyTransaction(transaction)

	state.LatestTimestamp = transaction.Tx.Timestamp

	state.SaveState()
	return nil
}

// * Magnus, s204509
// Apply the transaction by updating the balances of the affected users.
func (state *State) ApplyTransaction(transaction SignedTransaction) {
	if transaction.Tx.Type != "genesis" && transaction.Tx.Type != "reward" {
		state.AccountBalances[transaction.Tx.From] -= uint(transaction.Tx.Amount)
	}
	state.AccountNounces[transaction.Tx.From]++
	state.AccountBalances[transaction.Tx.To] += uint(transaction.Tx.Amount)
}

// * Magnus, s204509
// Validates a given signed transaction against the state. It validate the sender and receiver and amount and timestamp and the balance of the sender.
func (state *State) ValidateTransaction(signedTx SignedTransaction) error {
	if state.AccountNounces[signedTx.Tx.From]+1 != signedTx.Tx.SenderNounce {
		return fmt.Errorf("Transaction Nounce doesn't match account nounce")
	}

	if signedTx.Tx.Amount <= 0 {
		return fmt.Errorf("illegal to make a transaction with 0 or less coins")
	}

	if (state.LastBlockSerialNo == 0 && signedTx.Tx.Type == "genesis") || signedTx.Tx.Type == "reward" {
		return nil
	}

	addrOfTransaction, err := Crypto.GetAddressFromSignedTransaction(signedTx.Signature, signedTx.Tx.hash())
	if err != nil {
		return err
	}
	if AccountAddress(addrOfTransaction) != signedTx.Tx.From {
		return fmt.Errorf("sender of the transaction did not create the transaction!")
	}

	if state.LastBlockSerialNo != 0 && signedTx.Tx.Type == "genesis" {
		return fmt.Errorf("a genesis transaction can only be applied to the genesis block (serial 0)")
	}

	if signedTx.Tx.From == signedTx.Tx.To {
		return fmt.Errorf("a normal transaction is not allowed to same account")
	}

	if _, ok := state.AccountBalances[signedTx.Tx.From]; !ok {
		return fmt.Errorf("sending from Undefined Account \"%s\"", signedTx.Tx.From)
	}

	if state.AccountBalances[signedTx.Tx.From] < uint(signedTx.Tx.Amount) {
		return fmt.Errorf("Sender ain't that liquid right now")
	}

	return nil
}

// * Magnus, s204509
// Validates a list of transactions against the state.
func (state *State) ValidateTransactionList(transactionList SignedTransactionList) error {
	for i, t := range transactionList {
		err := state.ValidateTransaction(t)
		if err != nil {
			return fmt.Errorf("transaction nr. %d is not valid. Received Error: %s", i, err.Error())
		}
	}
	return nil
}

// * Magnus, s204509
// Adds a list of transaction to the state.
func (state *State) AddTransactionList(transactionList SignedTransactionList) error {
	for i, t := range transactionList {
		err := state.AddTransaction(t)
		if err != nil {
			return fmt.Errorf("Transaction idx[%d] is not able to be added. Received Error: %s", i, err.Error())
		}
	}
	return nil
}

// * Magnus, s204509
// Tries to add all transactions to a state
// This assumes that all the transactions that tries to be added have been validated before.
// This function is meant to be used to add the remaining transactions in the local memory pool after receiving a block
// Any transaction that has been validated before (is in the mempool) but is no more, must be invalidated (already applied) by the new block
// This removes duplicates
func (state *State) TryAddTransactions(transactionList SignedTransactionList) error {
	for _, t := range transactionList {
		state.AddTransaction(t) // It won't add the transaction if validation fails but will simply continue.
	}
	return nil
}

// * Niels, s204503
// recomputes state snapshot corresponding to a given index (serial no.) on the blockchain
// mutates state of {state} and the persisted snapshot
func (state *State) RecomputeState(deltaIdx int) {
	newState := BlankState()
	bc := LoadBlockchain()[:deltaIdx-1]

	for _, b := range bc {
		newState.ApplyBlock(b)
	}

	// Add pending transactions to the new state
	newState.TryAddTransactions(state.TxMempool)

	*state = newState
	state.SaveSnapshot()
}

// * Magnus, s204509 & Niels, s204503
// Given a block delta, try and add the new blocks to the current blockchain from the point where the fork happens.
func (state *State) TryMergeBlockDelta(deltaIdx int, newBlocks []Block) error {
	originalBlockchain := LoadBlockchain()
	originalSnapshot := LoadSnapshot()
	originalState := LoadState()
	agreedBlockchain := originalBlockchain[:deltaIdx-1]

	newState := BlankState()
	newState.SaveSnapshot()   // Save the snapshot to use for adding blocks
	newState.SaveState()      // Save the current state to handle transaction logic
	SaveBlockchain([]Block{}) // Save the empty block to handle block append

	// Add the old blocks from the local blockchain
	for _, block := range agreedBlockchain {
		blockErr := newState.AddBlock(block)
		if blockErr != nil {
			originalSnapshot.SaveSnapshot() // If it fails to add the new blocks, revert to the original snapshot and blockchain
			originalState.SaveState()       // It it fails to add the blocks revert to original state
			SaveBlockchain(originalBlockchain)
			return fmt.Errorf("One of the local blocks has been tampered with or otherwise corrupted. Got error: " + blockErr.Error())
		}
	}
	//Add the new blocks
	for _, block := range newBlocks {
		blockErr := newState.AddBlock(block)
		if blockErr != nil {
			originalSnapshot.SaveSnapshot() // If it fails to add the new blocks, revert to the original snapshot and blockchain
			originalState.SaveState()       // It it fails to add the blocks revert to original state
			SaveBlockchain(originalBlockchain)
			return fmt.Errorf("One of the new blocks are invalid. Got error: " + blockErr.Error())
		}
	}

	// Add pending transactions to the new state
	newState.TryAddTransactions(state.TxMempool)

	// Save the new current State
	newState.SaveState()

	// Save the new blockchain

	// Overwrite old state with NEW state
	*state = newState

	return nil
}

// * Niels, s204503
// Returns a blank state object
func BlankState() State {
	newState := State{}
	newState.AccountBalances = map[AccountAddress]uint{}
	newState.AccountBalances = map[AccountAddress]uint{}
	newState.AccountNounces = map[AccountAddress]uint{}
	newState.TxMempool = SignedTransactionList{}
	newState.DbFile = &os.File{}
	newState.LastBlockSerialNo = 0
	newState.LastBlockTimestamp = 0
	newState.LatestHash = [32]byte{}
	newState.LatestTimestamp = 0
	return newState
}

// * Emilie, s204471
// Loads the latest snapshot of the state. Each snapshot is meant as the state right after a block has been added.
func LoadSnapshot() State {
	return loadStateFromJSON("LatestSnapshot.json")
}

// * Emilie, s204471
// Given a state, save the state as the Current State, including local changes.
// This is different from a snapshot, as the current state also saves local changes, aka. transactions.
func (state *State) SaveState() error {
	return saveStateAsJSON(state, "CurrentState.json")
}

// * Asger, s204435
// Given a state, save the state as the local state snapshot.
// I.e. the state at the moment a new block is added. Any local Tx's are therefore not included.
func (state *State) SaveSnapshot() error {
	if len(state.TxMempool) > 0 { // Local transactions are not allowed
		return fmt.Errorf("cannot save snapshot of state with local changes")
	}

	return saveStateAsJSON(state, "LatestSnapshot.json")
}

// * Asger, s204435
// Function that saves a state as a json file
func saveStateAsJSON(state *State, filename string) error {
	txFile, _ := json.MarshalIndent(state, "", "  ")

	err := ioutil.WriteFile(shared.LocatePersistenceFile(filename, ""), txFile, 0644)
	if err != nil {
		panic(err)
	}

	return nil
}

// * Asger, s204435
// Function that loads a state from a JSON file
func loadStateFromJSON(filename string) State {
	data, err := os.ReadFile(shared.LocatePersistenceFile(filename, ""))
	if err != nil {
		return BlankState()
	}

	var state State
	err = json.Unmarshal(data, &state)
	if err != nil {
		return BlankState()
	}
	return state
}

// * Emilie, s204471
// Given a state, make a deep copy of the state and return the copy.
func (currState *State) copyState() State {
	copy := State{}

	copy.TxMempool = make([]SignedTransaction, 0)
	copy.AccountBalances = make(map[AccountAddress]uint)
	copy.AccountNounces = make(map[AccountAddress]uint)

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

	copy.TxMempool = append(copy.TxMempool, currState.TxMempool...)

	return copy
}
