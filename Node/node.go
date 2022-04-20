package node

import (
	Database "blockchain/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type NodeState struct {
	peerList []string
	state    Database.State
}

const httpPort = 8080
const bootstrapNode = "bootstrapNode"

//Models the balances data recived
type balancesResult struct {
	Hash     [32]byte                         `json:"block_hash`
	Balances map[Database.AccountAddress]uint `json:"balances`
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

func Run() error {
	fmt.Println(fmt.Sprintf("Listening on port %d", httpPort))
	startNode()

	nodeState := getNodeState()

	for true {
		for _, peer := range nodeState.peerList {
			peerState := getPeerState(peer)

			if peerState.state.LastBlockSerialNo > nodeState.state.LastBlockSerialNo {
				peerBlocks := getPeerBlocks(peer)
				nodeState.state.ApplyBlocks(peerBlocks)
			}

			nodeState.state.TryAddTransactions(peerState.state.TxMempool)

		}
	}

	return nil
}

func getPeerBlocks(peerAddr string) []Database.Block {
	return nil
}

func getPeerState(peerAddr string) NodeState {
	resp, err := http.Get("localhost:")

}

func getNodeState() NodeState {
	nodeState := NodeState{}
	nodeState.state = *Database.LoadState()
	nodeState.peerList = []string{bootstrapNode}
	return nodeState
}

func startNode() error {
	state := Database.LoadState()

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		balancesHandler(w, r, state)
	})

	http.HandleFunc("/transaction/create", func(w http.ResponseWriter, r *http.Request) {
		transactionHandler(w, r, state)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

func balancesHandler(w http.ResponseWriter, r *http.Request, state *Database.State) {
	writeResult(w, balancesResult{state.LatestHash, state.AccountBalances})
}

func transactionHandler(w http.ResponseWriter, r *http.Request, state *Database.State) {
	req := TxRequest{}
	err := readReq(r, &req)

	if err != nil {
		return
	}

	var transaction Database.Transaction
	println("TYPE OF REQUEST " + req.Type)
	switch req.Type {
	case "genesis":
		transaction = state.CreateGenesisTransaction(Database.AccountAddress(req.From), float64(req.Amount))

		fmt.Println("Genesis created" + Database.TxToString(transaction))

	case "reward":
		transaction = state.CreateReward(Database.AccountAddress(req.From), float64(req.Amount))

		fmt.Println("Reward created" + Database.TxToString(transaction))

	case "transaction":
		if req.To != "" {
			transaction = state.CreateTransaction(Database.AccountAddress(req.From), Database.AccountAddress(req.To), float64(req.Amount))

			fmt.Println("Transaction created" + Database.TxToString(transaction))
		}
	}

	fmt.Println(transaction)

	err = state.AddTransaction(transaction)
	if err != nil {
		return
	}

	status := Database.SaveTransaction(state.TxMempool)

	writeResult(w, TxResult{status})
}

//When using get method
func writeResult(w http.ResponseWriter, content interface{}) {
	contentJson, err := json.Marshal(content)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(contentJson)
}

//When using post method
func readReq(r *http.Request, reqBody interface{}) error {
	reqJson, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("ERROR WAS NOT NIL READ")
		return fmt.Errorf("unable to read request body. %s", err.Error())

	}
	defer r.Body.Close()

	err = json.Unmarshal(reqJson, reqBody)
	if err != nil {
		fmt.Printf("unable to unmarshal request body. %s", err.Error())
		return fmt.Errorf("unable to unmarshal request body. %s", err.Error())
	}

	return nil
}

// Business logic for interacting between nodes
// Uses an implementation of NodeConnectionManager interface
// to perform networking calls (or other types of queries, such as in MockConnectionManager)

// interface describes contracts that should be fulfilled in order for the node to be able to interact
// with peer nodes. Ideally, when this interface (w.i.p.) is implemented, we should be able to create
// the higher level algorithms such as sync, mining, consensus
type INode interface {
	Fetch2ndLevelPeerList()
	SendState()
	GetState()
	CheckNetworkStatus()
}

type Node struct {
	NCM   NodeConnectionManager
	FQDN  string
	Peers *PeerList
}

type NetworkStatusDTO struct {
	ReadyPeers  int
	BusyPeers   int
	FailedPeers int
	OnlinePeers int
}

// (integer) amount of (ready, busy, failed, online) peers
func CreateNetworkStatusDTO(ready int, busy int, failed int, online int) NetworkStatusDTO {
	return NetworkStatusDTO{ReadyPeers: ready, BusyPeers: busy, FailedPeers: failed, OnlinePeers: online}
}

// returns 4-tuple of (ready, busy, failed, online) peers
func (nws NetworkStatusDTO) ExtractTuple() (int, int, int, int) {
	return int(nws.ReadyPeers), int(nws.BusyPeers), int(nws.FailedPeers), int(nws.OnlinePeers)
}

func (node *Node) Fetch2ndLevelPeerList() {

}

func (node *Node) SendState() {

}

func (node *Node) FetchState(fqdn string) {
	stateData, err := node.NCM.FetchStateData(fqdn)
	println("%+v %s", stateData, err)
}

// pass ncm as argument?
// or have singleton declared somewhere?
// or like now, where it is called from node, which uses its own ncm to handle data context
// Check whether node has any registered peers and whether they are online and ready to communicate
func (node *Node) CheckNetworkStatus() (NetworkStatusDTO, error) {
	peers, err := LoadPeers()
	if err != nil {
		panic(err)
	}
	noOfPeers := countPeers(peers)

	if noOfPeers == 0 {
		return CreateNetworkStatusDTO(0, 0, 0, 0), nil
	}

	var readyPeers, busyPeers, failedPeers, onlinePeers int
	for _, peer := range peers.Peers {
		nodeRes := node.NCM.GetHeartBeat(peer.FQDN)
		switch statusCode := nodeRes.Status; statusCode {
		case READY:
			readyPeers++
			onlinePeers++
		case BUSY:
			busyPeers++
			onlinePeers++
		case FAILED:
			failedPeers++
			onlinePeers++
		case NO_CONN:
			{
			}
		default:
			panic(fmt.Sprintf("Status Code %d should not be able to be returned from heartbeat check", statusCode))
		}
	}

	return CreateNetworkStatusDTO(readyPeers, busyPeers, failedPeers, onlinePeers), err
}
