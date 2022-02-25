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
	SerialNo  int
}

type TransactionList []Transaction

type LoadedTransactions struct {
	Transactions TransactionList `json:"transactions"`
}

func (state *State) CreateTransaction(from AccountAddress, to AccountAddress, amount float64) Transaction {
	fmt.Println("CreateTransaction() called")
	t := Transaction{
		from,
		to,
		amount,
		makeTimestamp(),
		"transaction",
		state.getNextTxSerialNo(),
	}

	fmt.Println(t)
	return t
}

func (state *State) CreateReward(to AccountAddress, amount float64) Transaction {
	fmt.Println("CreateReward() called")
	r := Transaction{
		"system",
		to,
		amount,
		makeTimestamp(),
		"reward",
		state.getNextTxSerialNo(),
	}

	fmt.Println(r)
	return r
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
