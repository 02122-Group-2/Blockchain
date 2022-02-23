package database

import (
	"math/rand"
	"testing"
)

var names = []string{"magn", "niels", "emilie", "asger", "alberto", "bill", "andrej"}

func makeDummyTransaction() Transaction {
	return dbInfo.CreateTransaction(names[rand.Int()*7], "niels", 6969.0)
}

func TestCreateBlock(t *testing.T) {

}
