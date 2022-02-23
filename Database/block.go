package database

type Block struct {
	Header       BlockHeader
	Transactions []Transaction
}

type BlockHeader struct {
	ParentHash string
	CreatedAt  int64
	SerialNo   int
}

func (state *State) CreateBlock(txs []Transaction) Block {
	return Block{
		BlockHeader{
			state.getLastHash(),
			makeTimestamp(),
			state.getNextBlockSerialNo(),
		},
		txs,
	}
}
