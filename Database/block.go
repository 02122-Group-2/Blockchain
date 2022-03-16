package database

import "fmt"

type Block struct {
	Header       BlockHeader   `json: "Header"`
	Transactions []Transaction `json: "Transactions"`
}

type BlockHeader struct {
	ParentHash string `json: "ParentHash"`
	CreatedAt  int64  `json: "CreatedAt"`
	SerialNo   int    `json: "SerialNo"`
}

func (state *State) CreateBlock(txs []Transaction) Block {
	return Block{
		BlockHeader{
			state.getLatestHash(),
			makeTimestamp(),
			state.getNextBlockSerialNo(),
		},
		txs,
	}
}

func (state *State) AddBlock(block Block) bool {
	if err := state.ValidateBlock(block); err != nil {
		return false
	}

	//TODO
	// if err := state.PersistBlock(block); err != nil {
	// 	return false
	// }

	state.lastBlockSerialNo = block.Header.SerialNo
	state.TxMempool = nil

	return true
}

func (state *State) ValidateBlock(block Block) error {
	if block.Header.ParentHash != state.latestHash {
		return fmt.Errorf("latest hash in state must be parent hash")
	}

	if block.Header.SerialNo != state.lastBlockSerialNo+1 {
		return fmt.Errorf("serial number must be 1 larger than last block's serial number")
	}

	if block.Header.CreatedAt <= state.latestTimestamp {
		return fmt.Errorf("time stamp must be newer than latest time stamp in state")
	}

	if err := state.ValidateTransactionList(block.Transactions); err != nil {
		return fmt.Errorf("transactions in block does not match transactions in state")
	}

	return nil
}
