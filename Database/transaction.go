package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	currWD, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(filepath.Join(currWD, "Transactions.json"))
	if err != nil {
		panic(err)
	}

	var loadedTransactions LoadedTransactions
	json.Unmarshal(data, &loadedTransactions)

	return loadedTransactions.Transactions
}

func SaveTransaction(transactionList TransactionList) bool {
	toSave := LoadedTransactions{transactionList}
	txFile, _ := json.MarshalIndent(toSave, "", "  ")

	err := ioutil.WriteFile("./Transactions.json", txFile, 0644)
	if err != nil {
		panic(err)
	}

	return true
}

func Filter(vs TransactionList, f func(Transaction) bool) []Transaction {
	filtered := make([]Transaction, 0)
	for _, v := range vs {
		if f(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func TxToString(transaction Transaction) string {
	return "From: " + string(transaction.From) + "\n To: " + string(transaction.To) + "\n Amount: " + fmt.Sprintf("%v", transaction.Amount)
}
