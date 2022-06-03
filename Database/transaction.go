package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type AccountAddress string

type Transaction struct {
	From         AccountAddress
	To           AccountAddress
	Amount       float64
	SenderNounce uint
	Timestamp    int64 // UNIX time
	Type         string
}

type TransactionList []Transaction

type LoadedTransactions struct {
	Transactions TransactionList `json:"transactions"`
}

func (state *State) newAccountAddr(value string) AccountAddress {
	return AccountAddress(value)
}

// Create a custom transaction. Used as a helper function.
func (state *State) CreateCustomTransaction(from AccountAddress, to AccountAddress, amount float64, _type string) Transaction {
	accountNounce := state.AccountNounces[from] + 1
	t := Transaction{
		from,
		to,
		amount,
		accountNounce,
		makeTimestamp(),
		_type,
	}

	fmt.Println(t)
	return t
}

// Creates an ordinary transaction between two users.
// Takes two addresses (strings) and the amount sent (float)
func (state *State) CreateTransaction(from AccountAddress, to AccountAddress, amount float64) Transaction {
	return state.CreateCustomTransaction(from, to, amount, "transaction")
}

// Creates a genesis type transaction from the system to a certain user.
// Takes the receiver address (string) and the amount sent (float)
func (state *State) CreateGenesisTransaction(accountAddress AccountAddress, amount float64) Transaction {
	return state.CreateCustomTransaction("system", accountAddress, amount, "genesis")
}

// Creates a reward type transaction from the system to a certain user.
// Takes the receiver address (string) and the amount sent (float)
func (state *State) CreateReward(accountAddress AccountAddress, amount float64) Transaction {
	return state.CreateCustomTransaction("system", accountAddress, amount, "reward")
}

// Given a list of transactions, it saves these transactions as a JSON string in a local text file.
// Returns a boolean value indicating whether or not it was saved succesfully.
// This is not used in older version of the blockchain.
func SaveTransaction(transactionList TransactionList) bool {
	toSave := LoadedTransactions{transactionList}
	txFile, _ := json.MarshalIndent(toSave, "", "  ")

	err := ioutil.WriteFile(localDirToFileFolder+"Transactions.json", txFile, 0644)
	if err != nil {
		panic(err)
	}

	return true
}

// Loads the local transactions, saved in the transactions.json file. This is deprecated and only used in early versions of the blockchain.
// It returns a list of transactions.
func LoadTransactions() TransactionList {
	data, err := os.ReadFile(localDirToFileFolder + "Transactions.json")
	if err != nil {
		panic(err)
	}

	var loadedTransactions LoadedTransactions
	json.Unmarshal(data, &loadedTransactions)

	return loadedTransactions.Transactions
}

// Given a list of transactions, save the list of transactions as the local transactions.
func (transaction_list *TransactionList) SaveTransactions() error {
	return saveTransactionsAsJSON(transaction_list, "Transactions.json")
}

// Function that saves list of transactions as a json file
func saveTransactionsAsJSON(transaction_list *TransactionList, filename string) error {
	txFile, _ := json.MarshalIndent(transaction_list, "", "  ")

	err := ioutil.WriteFile(localDirToFileFolder+filename, txFile, 0644)
	if err != nil {
		panic(err)
	}

	return nil
}

// Formats a given transaction to text format.
func TxToString(transaction Transaction) string {
	return "From: " + string(transaction.From) + "\n To: " + string(transaction.To) + "\n Amount: " + fmt.Sprintf("%v", transaction.Amount)
}
