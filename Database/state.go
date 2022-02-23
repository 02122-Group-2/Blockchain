package database

import (
	"os"
	"time"
)

type State struct {
	txMemPool []Transaction // incoming Txs, not yet added to chain
	dbFile    *os.File

	lastTxSerialNo    int
	lastBlockSerialNo int
	lastHash          string
}

func makeTimestamp() int64 {
	return time.Now().UnixNano()
}

func (s *State) getNextTxSerialNo() int {
	curNo := s.lastTxSerialNo + 1
	s.lastTxSerialNo = curNo
	return curNo
}

func (s *State) getNextBlockSerialNo() int {
	curNo := s.lastBlockSerialNo + 1
	s.lastBlockSerialNo = curNo
	return curNo
}

func (s *State) getLastHash() string {
	return s.lastHash
}
