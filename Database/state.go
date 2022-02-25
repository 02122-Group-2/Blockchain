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
	curNo := s.lastTxSerialNo + 1
	s.lastTxSerialNo = curNo
	return curNo
}

func (s *State) getNextBlockSerialNo() int {
	curNo := s.lastBlockSerialNo + 1
	s.lastBlockSerialNo = curNo
	return curNo
}

func (s *State) getLatestHash() string {
	return s.latestHash
}

func LoadState() (*State, error) {
	genesis := LoadGenesis()

	balances := make(map[AccountAddress]uint)
	for account, balance := range genesis.Balances {
		balances[account] = uint(balance)
	}

	var file *os.File
	state := &State{balances, make([]Transaction, 0), file, 0, 0, ""} //TODO fix missing hash

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

	return nil

}

func (state *State) ApplyTransaction(transaction Transaction) {
	state.Balances[transaction.From] -= uint(transaction.Amount)
	state.Balances[transaction.To] += uint(transaction.Amount)
}

func (state *State) ValidateTransaction(transaction Transaction) error {
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
