package node

import Database "blockchain/Database"

type NodeFromPostRequest struct {
	PeerSet PeerSet                       `json:"peer_set"`
	State   Database.StateFromPostRequest `json:"state"`
}

type Node struct {
	PeerSet PeerSet        `json:"peer_set"`
	State   Database.State `json:"state"`
}

const httpPort = 8080
const bootstrapNode = "localhost:8080"

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