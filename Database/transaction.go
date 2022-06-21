package database

import (
	Crypto "blockchain/Cryptography"
	shared "blockchain/Shared"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// * Magnus, s204509 #clutch
type AccountAddress string

// * Asger, s204435
type Transaction_Old struct {
	From         AccountAddress
	To           AccountAddress
	Amount       float64
	SenderNounce uint
	Timestamp    int64 // UNIX time
	Type         string
}

// * Emilie, s204471
type TransactionList []Transaction_Old

// * Magnus, s204509
type SignedTransaction struct {
	Signature []byte
	Tx        Transaction_Old
}

// * Magnus, s204509
type SignedTransactionList []SignedTransaction

// * Emilie, s204471
type LoadedTransactions struct {
	Transactions SignedTransactionList `json:"transactions"`
}

// * Asger, s204435
func (transaction Transaction_Old) toJsonString() string {
	json, err := json.Marshal(transaction)
	if err != nil {
		panic(err)
	}
	return string(json)
}

// * Asger, s204435
func (transaction *Transaction_Old) hash() [32]byte {
	return Crypto.HashTransaction(transaction.toJsonString())
}

// * Magnus, s204509
func (state *State) newAccountAddr(value string) AccountAddress {
	return AccountAddress(value)
}

// * Niels, s204503
// Create a custom transaction. Used as a helper function.
func (state *State) CreateCustomTransaction(from AccountAddress, to AccountAddress, amount float64, _type string) Transaction_Old {
	accountNounce := state.AccountNounces[from] + 1
	t := Transaction_Old{
		from,
		to,
		amount,
		accountNounce,
		shared.MakeTimestamp(),
		_type,
	}

	// fmt.Println(t)
	return t
}

// * Emilie, s204471
// Creates an ordinary transaction between two users.
// Takes two addresses (strings) and the amount sent (float)
func (state *State) CreateTransaction(from AccountAddress, to AccountAddress, amount float64) Transaction_Old {
	return state.CreateCustomTransaction(from, to, amount, "transaction")
}

// * Magnus, s204509
// Takes a wallet, password, receiver, amount and returns a signed transaction
func (state *State) CreateSignedTransaction(wallet Crypto.Account, password string, receiver AccountAddress, amount float64) (SignedTransaction, error) {
	tx := state.CreateTransaction(AccountAddress(wallet.Address), receiver, amount)
	return state.SignTransaction(wallet, password, tx)
}

// * Emilie, s204471
// Creates a genesis type transaction from the system to a certain user.
// Takes the receiver address (string) and the amount sent (float)
func (state *State) CreateGenesisTransaction(accountAddress AccountAddress, amount float64) SignedTransaction {
	return SignedTransaction{
		Signature: []byte{},
		Tx:        state.CreateCustomTransaction("system", accountAddress, amount, "genesis"),
	}
}

// * Emilie, s204471
// Creates a reward type transaction from the system to a certain user.
// Takes the receiver address (string) and the amount sent (float)
// Is automatically created as signed transaction
func (state *State) CreateReward(accountAddress AccountAddress, amount float64) SignedTransaction {
	return SignedTransaction{
		Signature: []byte{},
		Tx:        state.CreateCustomTransaction("system", accountAddress, amount, "reward"),
	}
}

// * Magnus, s204509
// Given the password for the wallet and a regular transaction, sign the transaction, if the sender is equal to the address of the wallet
// Returns the signed transaction or an error
func (state *State) SignTransaction(wallet Crypto.Account, password string, transaction Transaction_Old) (SignedTransaction, error) {
	if transaction.From != AccountAddress(wallet.Address) {
		return SignedTransaction{}, fmt.Errorf("this transaction is not able to be signed by you!")
	}
	txHash := transaction.hash()
	signature, err := wallet.SignTransaction(password, txHash)
	if err != nil {
		return SignedTransaction{}, err
	}
	return SignedTransaction{Signature: signature, Tx: transaction}, nil
}

// * Asger, s204435
func ClearTransactions() {
	err := os.Truncate(shared.LocatePersistenceFile("Transactions.json", ""), 0)
	if err != nil {
		panic(err)
	}
}

// * Asger, s204435
// Given a list of transactions, it saves these transactions as a JSON string in a local text file.
// Returns a boolean value indicating whether or not it was saved succesfully.
func SaveTransaction(transactionList SignedTransactionList) bool {
	toSave := LoadedTransactions{transactionList}
	txFile, _ := json.MarshalIndent(toSave, "", "  ")

	err := ioutil.WriteFile(shared.LocatePersistenceFile("Transactions.json", ""), txFile, 0644)
	if err != nil {
		panic(err)
	}

	return true
}

// * Asger, s204435
// Loads the local transactions, saved in the transactions.json file. This is deprecated and only used in early versions of the blockchain.
// It returns a list of transactions.
func LoadTransactions() SignedTransactionList {
	data, err := os.ReadFile(shared.LocatePersistenceFile("Transactions.json", ""))
	if err != nil {
		panic(err)
	}

	var loadedTransactions LoadedTransactions
	json.Unmarshal(data, &loadedTransactions)

	return loadedTransactions.Transactions
}

// * Asger, s204435
// Given a list of transactions, save the list of transactions as the local transactions.
func (transaction_list *TransactionList) SaveTransactions() error {
	return saveTransactionsAsJSON(transaction_list, "Transactions.json")
}

// * Asger, s204435
// Function that saves list of transactions as a json file
func saveTransactionsAsJSON(transaction_list *TransactionList, filename string) error {
	txFile, _ := json.MarshalIndent(transaction_list, "", "  ")

	err := ioutil.WriteFile(shared.LocatePersistenceFile(filename, ""), txFile, 0644)
	if err != nil {
		panic(err)
	}

	return nil
}

// * Magnus, s204509
// Formats a given transaction to text format.
func TxToString(transaction Transaction_Old) string {
	return "From: " + string(transaction.From) + "\n To: " + string(transaction.To) + "\n Amount: " + fmt.Sprintf("%v", transaction.Amount)
}

// * Magnus, s204509
// Formats a given signed transaction to text format.
func SignedTxToString(transaction SignedTransaction) string {
	return "From: " + string(transaction.Tx.From) + "\n To: " + string(transaction.Tx.To) + "\n Amount: " + fmt.Sprintf("%v", transaction.Tx.Amount)
}
