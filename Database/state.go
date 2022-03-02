package database

import (
	"fmt"
	"os"

	// "path/filepath"
	"time"
)

type State struct {
	Balances  map[AccountAddress]uint
	txMempool TransactionList
	dbFile    *os.File

	lastTxSerialNo    int
	lastBlockSerialNo int
	latestHash        string
}

func makeTimestamp() int64 {
	return time.Now().UnixNano()
}

func (s *State) getNextTxSerialNo() int {
	return s.lastTxSerialNo + 1
}

func (s *State) getNextBlockSerialNo() int {
	return s.lastBlockSerialNo + 1
}

func (s *State) getLatestHash() string {
	return s.latestHash
}

func LoadState() (*State, error) {
	var file *os.File
	state := &State{make(map[AccountAddress]uint), make([]Transaction, 0), file, 0, 0, ""} //TODO fix missing hash

	genesis := LoadGenesis()

	for account, balance := range genesis.Balances {
		t := state.CreateGenesisTransaction(account, (float64(balance)))
		state.AddTransaction(t)
	}

	loadedTransactions := LoadTransactions()

	for _, t := range loadedTransactions {
		if state.AddTransaction(t) != nil {
			panic("Transaction not allowed")
		}
	}

	return state, nil
}

func (state *State) AddTransaction(transaction Transaction) error {
	if err := state.ValidateTransaction(transaction); err != nil {
		return err
	}

	state.txMempool = append(state.txMempool, transaction)

	state.ApplyTransaction(transaction)

	state.lastTxSerialNo++
	return nil
}

func (state *State) ApplyTransaction(transaction Transaction) {
	if transaction.Type != "genesis" && transaction.Type != "reward" {
		state.Balances[transaction.From] -= uint(transaction.Amount)
	}
	state.Balances[transaction.To] += uint(transaction.Amount)
}

func (state *State) ValidateTransaction(transaction Transaction) error {
	if (state.lastBlockSerialNo == 0 && transaction.Type == "genesis") || transaction.Type == "reward" {
		return nil
	}

	if transaction.SerialNo != (state.lastTxSerialNo + 1) {
		return fmt.Errorf("SerialNo. violates transaction order")
	}

	if transaction.From == transaction.To {
		return fmt.Errorf("A normal transaction is not allowed to same account")
	}

	if _, err := state.Balances[transaction.From]; !err {
		return fmt.Errorf("Sending from Undefined Account")
	}
	if transaction.Amount <= 0 {
		return fmt.Errorf("Illegal to make a transaction with 0 or less coins.")
	}
	if state.Balances[transaction.From] < uint(transaction.Amount) {
		return fmt.Errorf("u broke")
	}

	return nil
}


