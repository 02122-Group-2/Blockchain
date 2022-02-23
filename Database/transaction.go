package database

import "fmt"

type Transaction struct {
	From      string
	To        string
	Amount    float64
	Timestamp int64 // UNIX time
	Type      string
	SerialNo  int
}

func (dbInfo *DatabaseInfo) CreateTransaction(from string, to string, amount float64) Transaction {
	fmt.Println("CreateTransaction() called")
	t := Transaction{
		from,
		to,
		amount,
		makeTimestamp(),
		"transaction",
		dbInfo.getNextSerialNo(),
	}

	fmt.Println(t)
	return t
}

func (dbInfo *DatabaseInfo) CreateReward(to string, amount float64) Transaction {
	fmt.Println("CreateReward() called")
	r := Transaction{
		"system",
		to,
		amount,
		makeTimestamp(),
		"reward",
		dbInfo.getNextSerialNo(),
	}

	fmt.Println(r)
	return r
}
