package database

type Block struct {
	ParentHash, Hash string
	CreatedAt        int64
	SerialNo         int
	Transactions     []Transaction
}

func createBlock(parentHash string, createdAt int64, tx []Transaction) Block {
	block := Block{
		parentHash,
		"implement me",
		makeTimestamp(),
		69,
		tx,
	}
	return block
}
