package database

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

type State struct {
	Balances  map[Account]uint
	txMempool []Transaction
	dbFile    *os.File

	latestHash string
}

func LoadState() (*State, error) {
	currWD, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	genesis, err := loadGenesis()
	if err != nil {
		return nil, err
	}

	balances := make(map[Account]uint)
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

		blockFsJson := scanner.Bytes()
		var blockFs BlockFS
		err = json.Unmarshal(blockFsJson, &blockFs)
		if err != nil {
			return nil, err
		}

		err = state.applyBlock(blockFs.Value)
		if err != nil {
			return nil, err
		}

		state.latestHash = blockFs.Key
	}

	return state, nil

}

func AddTransaction() {

}

func Persist() {

}
