package node

import (
	Database "blockchain/Database"
)

// * File: Niels, s204503

type NodeFromPostRequest struct {
	Address     string                        `json:"address"`
	PeerSet     PeerSet                       `json:"peer_set"`
	State       Database.StateFromPostRequest `json:"state"`
	ChainHashes []string                      `json:"chain_hashes"`
}

type Node struct {
	Address     string         `json:"address"`
	PeerSet     PeerSet        `json:"peer_set"`
	State       Database.State `json:"state"`
	ChainHashes []string       `json:"chain_hashes"`
}

//Models the balances data recived
type balancesResult struct {
	Hash     [32]byte                         `json:"block_hash"`
	Balances map[Database.AccountAddress]uint `json:"balances"`
}

//Models the data for sending
type TxRequest struct {
	From   string  `json:"From"`
	To     string  `json:"To"`
	Amount float64 `json:"Amount"`
	Type   string  `json:"Type"`
}

//Models the transaction data recived
type TxResult struct {
	Status bool `json:"status"`
}

type PingResponse struct {
	Address string
	Ok      bool
	Latency int64
}

type PingResponseList []PingResponse

func (p PingResponseList) Len() int           { return len(p) }
func (p PingResponseList) Less(i, j int) bool { return p[i].Latency < p[j].Latency }
func (p PingResponseList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
