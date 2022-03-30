package database

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	// "path/filepath"
	"time"
)

type State struct {
	AccountBalances    map[AccountAddress]uint `json: "AccountBalances"`
  AccountNounces     map[AccountAddress]uint `json: "AccountNounces"`
	TxMempool          TransactionList         `json: "TxMempool"`
	DbFile             *os.File                `json: "DbFile"`
	LastBlockSerialNo  int                     `json: "LastBlockSerialNo"`
	LastBlockTimestamp int64                   `json: "LastBlockTimestamp"`
	LatestHash         [32]byte                `json: "LatestHash"`
	LatestTimestamp    int64                   `json: "LatestTimestamp"`
}

func makeTimestamp() int64 {
	return time.Now().UnixNano()
}

func (s *State) getNextBlockSerialNo() int {
	return s.LastBlockSerialNo + 1
}

func (s *State) getLatestHash() [32]byte {
	return s.LatestHash
}

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

// Creates a state based from the data in the local blockchain.db file.
func LoadState() (*State, error) {
	var file *os.File
	state := &State{make(map[AccountAddress]uint), make(map[AccountAddress]uint), make([]Transaction, 0), file, 0, 0, [32]byte{}, 0}

	localBlockchain := LoadBlockchain()
	err := state.ApplyBlocks(localBlockchain)

	// set LatestHash property to hash of latest inserted block,
	// since this is the hash that should be used to validate next block
	// state.LatestHash = localBlockchain[len(localBlockchain)-1]

	fmt.Printf("state.LatestHash:%x\n", state.LatestHash)
	if err != nil {
		panic(err)
	}

	(*state).SaveSnapshot()
	return state, nil
}

// Adds a transaction to the state. It will validate the transaction, then apply the transaction to the state,
// then add the transaction to its MemPool and update its latest timestamp field.
func (state *State) AddTransaction(transaction Transaction) error {
	if err := state.ValidateTransaction(transaction); err != nil {
		return err
	}

	state.TxMempool = append(state.TxMempool, transaction)

	state.ApplyTransaction(transaction)

	state.LatestTimestamp = transaction.Timestamp
	return nil
}

// Apply the transaction by updating the balances of the affected users.
func (state *State) ApplyTransaction(transaction Transaction) {
	if transaction.Type != "genesis" && transaction.Type != "reward" {
		state.AccountBalances[transaction.From] -= uint(transaction.Amount)
	}
	state.AccountNounces[transaction.From]++;
	state.AccountBalances[transaction.To] += uint(transaction.Amount)
}

// Validates a given transaction against the state. It validate the sender and receiver and amount and timestamp and the balance of the sender.
func (state *State) ValidateTransaction(transaction Transaction) error {
	if state.AccountNounces[transaction.From]+1 != transaction.SenderNounce  {
		return fmt.Errorf("Transaction Nounce doesn't match account nounce")
	}

	if (state.LastBlockSerialNo == 0 && transaction.Type == "genesis") || transaction.Type == "reward" {
		return nil
	}

	if transaction.From == transaction.To {
		return fmt.Errorf("a normal transaction is not allowed to same account")
	}

	if _, ok := state.AccountBalances[transaction.From]; !ok {
		return fmt.Errorf("sending from Undefined Account \"%s\"", transaction.From)
	}
	if transaction.Amount <= 0 {
		return fmt.Errorf("illegal to make a transaction with 0 or less coins")
	}

	if state.AccountBalances[transaction.From] < uint(transaction.Amount) {
		return fmt.Errorf("Sender ain't that liquid right now")
	}

	return nil
}

// Validates a list of transactions against the state.
func (state *State) ValidateTransactionList(transactionList TransactionList) error {
	for i, t := range transactionList {
		err := state.ValidateTransaction(t)
		if err != nil {
			return fmt.Errorf("Transaction nr. %d is not valid. Received Error: %s", i, err.Error())
		}
	}
	return nil
}

// Adds a list of transaction to the state.
func (state *State) AddTransactionList(transactionList TransactionList) error {
	for i, t := range transactionList {
		err := state.AddTransaction(t)
		if err != nil {
			return fmt.Errorf("Transaction idx[%d] is not able to be added. Received Error: %s", i, err.Error())
		}
	}
	return nil
}

// Tries to add all transactions to a state
// This assumes that all the transactions that tries to be added have been validated before.
// This function is meant to be used to add the remaining transactions in the local memory pool after receiving a block
// Any transaction that has been validated before (is in the mempool) but is no more, must be invalidated (already applied) by the new block
// This removes duplicates 
func (state *State) TryAddTransactions(transactionList TransactionList) error {
	for _, t := range transactionList {
		state.AddTransaction(t) // It won't add the transaction if validation fails but will simply continue.
	}
	return nil
}
