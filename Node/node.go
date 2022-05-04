package node

import (
	Database "blockchain/database"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type NodeState struct {
	PeerList []string       `json:"peer_list"`
	State    Database.State `json:"state"`
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

	nodeState := GetNodeState()

	for true {

		updateNodeState(&nodeState) // Updates nodeState to get local changes in case any has been made

		// Creates a copy of the peers to check, to avoid weird looping as we append to the slice in the loop
		peersToCheck := make([]string, len(nodeState.PeerList))
		copy(peersToCheck, nodeState.PeerList)

		for _, peer := range peersToCheck {

			peerState := GetPeerState(peer)

			if peerState.State.LastBlockSerialNo > nodeState.State.LastBlockSerialNo {
				peerBlocks := getPeerBlocks(peer)
				nodeState.State.ApplyBlocks(peerBlocks)
			}

			nodeState.State.TryAddTransactions(peerState.State.TxMempool)

			for _, peer2 := range peerState.PeerList {
				if !contains(nodeState.PeerList, peer2) && ping(peer2) {
					nodeState.PeerList = append(nodeState.PeerList, peer2)
				}
			}
		}
		time.Sleep(40 * time.Second)
	}

	return nil
}

func ping(peerAddr string) bool {
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", peerAddr, timeout)
	conn.Close()
	if err != nil {
		fmt.Println("Site unreachable, error: ", err)
		return false
	}
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getPeerBlocks(peerAddr string) []Database.Block {
	return nil
}

//The following is done using POST,
//The header contain the address of the peer that is currently being accessed
//The body should contain the current state of the requesting node
func GetPeerState(peerAddr string) NodeState {

	currNodeState := GetNodeState()
	jsonData, err := json.Marshal(currNodeState)

	resp, err := http.Post("http://"+peerAddr+"/getState", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln(err)
		fmt.Println("something went wrong when posting")
	}
	fmt.Println("Successfully posted current state")

	peerNodeState := NodeState{}

	readResp(resp, &peerNodeState)
	//At this point the data recived should have been saved into peerNodeState

	return peerNodeState
}

func GetNodeState() NodeState {
	nodeState := NodeState{}
	nodeState.State = *Database.LoadState()
	nodeState.PeerList = []string{bootstrapNode}
	return nodeState
}

// Updated the current nodeState.State field to get latest local changes, in case the user has added local transactions.
func updateNodeState(currentNodeState *NodeState) {
	currentNodeState.State = *Database.LoadState()
}

func startNode() error {
	state := Database.LoadState()

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		balancesHandler(w, r, state)
	})
	fmt.Println("/balances/list setup complete")

	http.HandleFunc("/transaction/create", func(w http.ResponseWriter, r *http.Request) {
		transactionHandler(w, r, state)
	})
	fmt.Println("/transaction/create setup complete")

	http.HandleFunc("/getState", func(w http.ResponseWriter, r *http.Request) {
		getStateHandler(w, r, state)
	})
	fmt.Println("/getState setup complete")

	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

func getStateHandler(w http.ResponseWriter, r *http.Request, state *Database.State) {
	// Read the body containing the state of the node requesting
	getStateRequest := NodeState{}

	err := readReq(r, &getStateRequest)
	if err != nil {
		return
	}

	//TODO: Do something with the peer state

	//Response: Send your current state
	nodeState := GetNodeState()

	fmt.Println(nodeState.PeerList)
	writeResult(w, NodeState{PeerList: nodeState.PeerList, State: nodeState.State})

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

//Writing the result from the server
func writeResult(w http.ResponseWriter, content interface{}) {
	contentJson, err := json.Marshal(content)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(contentJson)
	fmt.Println("Server response sent")
}

//Reading the request when using POST method
func readReq(r *http.Request, reqBody interface{}) error {
	reqJson, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return fmt.Errorf("unable to read request body. %s", err.Error())

	}
	defer r.Body.Close()

	err = json.Unmarshal(reqJson, reqBody)
	if err != nil {
		return fmt.Errorf("unable to unmarshal request body. %s", err.Error())
	}

	return nil
}

func readResp(r *http.Response, reqBody interface{}) error {
	reqJson, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return fmt.Errorf("unable to read request body. %s", err.Error())

	}
	defer r.Body.Close()

	err = json.Unmarshal(reqJson, reqBody)
	if err != nil {
		return fmt.Errorf("unable to unmarshal request body. %s", err.Error())
	}
	return nil
}
