package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	Balances  map[AccountAddress]uint
	txMempool []Transaction
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
	currWD, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	genesis := LoadGenesis()

	balances := make(map[AccountAddress]uint)
	for account, balance := range genesis.Balances {
		balances[account] = uint(balance)
	}

	file, err := os.OpenFile(filepath.Join(currWD, "database", "block.db"), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	state := &State{balances, make([]Transaction, 0), file, ""} //TODO fix missing hash

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var transaction Transaction
		err = json.Unmarshal(scanner.Bytes(), &transaction)
		if err != nil {
			return nil, err
		}

		if err := state.ValidateTransaction(transaction); err != nil {
			return nil, err
		}

	}

	return state, nil
}

func (state *State) AddTransaction(transaction Transaction) error {
	if err := state.ValidateTransaction(transaction); err != nil {
		return err
	}

	state.txMempool = append(state.txMempool, transaction)

	return nil

}

func (state *State) ValidateTransaction(transaction Transaction) error {
	if state.Balances[AccountAddress(transaction.From)] < uint(transaction.Amount) {
		return fmt.Errorf("u broke")
	}

	state.Balances[AccountAddress(transaction.From)] -= uint(transaction.Amount)
	state.Balances[AccountAddress(transaction.To)] += uint(transaction.Amount)

	return nil

}
