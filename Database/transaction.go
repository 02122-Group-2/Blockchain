package database

import (
	"encoding/json"
	"fmt"
	"os"
)

type AccountAddress string

type Transaction struct {
	From      AccountAddress
	To        AccountAddress
	Amount    float64
	Timestamp int64 // UNIX time
	Type      string
}

type TransactionList []Transaction

type LoadedTransactions struct {
	Transactions TransactionList `json:"transactions"`
}

func (state *State) CreateCustomTransaction(from AccountAddress, to AccountAddress, amount float64, _type string) Transaction {
	fmt.Println("CreateTransaction() called")
	t := Transaction{
		from,
		to,
		amount,
		makeTimestamp(),
		_type,
	}

	fmt.Println(t)
	return t
}

func (state *State) CreateTransaction(from AccountAddress, to AccountAddress, amount float64) Transaction {
	return state.CreateCustomTransaction(from, to, amount, "transaction")
}

func (state *State) CreateGenesisTransaction(accountAddress AccountAddress, amount float64) Transaction {
	return state.CreateCustomTransaction("system", accountAddress, amount, "genesis")
}

func (state *State) CreateReward(accountAddress AccountAddress, amount float64) Transaction {
	return state.CreateCustomTransaction("system", accountAddress, amount, "reward")
}

func LoadTransactions() TransactionList {
	data, err := os.ReadFile("./Transactions.json")
	if err != nil {
		panic(err)
	}

	var loadedTransactions LoadedTransactions
	json.Unmarshal(data, &loadedTransactions)

	return loadedTransactions.Transactions
}
