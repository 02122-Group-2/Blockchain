package database

import (
	"math/rand"
	"testing"
)

var names = []AccountAddress{"Magn", "Niels", "Emilie", "Asger", "Alberto", "Bill", "Andrej"}

func makeDummyTransaction() Transaction {
	return state.CreateTransaction(names[rand.Int()*7], "niels", 6969.0)
}

func TestCreateBlock(t *testing.T) {

}
