package database

type Block struct {
	ParentHash, Hash    string
	CreatedAt, SerialNo int
	Transactions        []Transaction
}
