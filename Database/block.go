package database

import (
	Crypto "blockchain/Cryptography"
)

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

func (state *State) validateBlock(block Block) bool {
	if block.Header.ParentHash != state.latestHash {
		return false
	}

	if block.Header.SerialNo != state.getNextBlockSerialNo() {
		return false
	}

	if block.Header.CreatedAt <= state.getLatestBlock().Header.CreatedAt {
		return false
	}

	if 

	return true
}