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
	"strconv"
	"time"
)

type NodeStateFromPostRequest struct {
	PeerList []string                      `json:"peer_list"`
	State    Database.StateFromPostRequest `json:"state"`
}

type NodeState struct {
	PeerList []string       `json:"peer_list"`
	State    Database.State `json:"state"`
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

func Run() error {
	fmt.Printf("Listening on port %d\n", httpPort)
	go synchronization()
	startNode()
	return nil
}

func synchronization() {
	nodeState := GetNodeState()

	for {
		updateNodeState(&nodeState) // Updates nodeState to get local changes in case any has been made

		// Creates a copy of the peers to check, to avoid weird looping as we append to the slice in the loop
		peersToCheck := make([]string, len(nodeState.PeerList))
		copy(peersToCheck, nodeState.PeerList)

		for _, peer := range peersToCheck {

			peerState := GetPeerState(peer)

			fmt.Println("Got peer state")
			fmt.Println(peerState)

			if peerState.State.LastBlockSerialNo > nodeState.State.LastBlockSerialNo {
				peerBlocks := GetPeerBlocks(peer, nodeState.State.LastBlockSerialNo)
				nodeState.State.ApplyBlocks(peerBlocks)
			}

			nodeState.State.TryAddTransactions(peerState.State.TxMempool)

			for _, peer2 := range peerState.PeerList {
				if !contains(nodeState.PeerList, peer2) && Ping(peer2) {
					nodeState.PeerList = append(nodeState.PeerList, peer2)
				}
			}
			time.Sleep(40 * time.Second)
		}
	}
}

func Ping(peerAddr string) bool {
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", peerAddr, timeout)
	if err != nil {
		fmt.Println("Site unreachable, error: ", err)
		return false
	}
	conn.Close()
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

func GetPeerBlocks(peerAddr string, lastLocalBlockSerialNo int) []Database.Block {
	URI := fmt.Sprintf("http://"+peerAddr+"/blockDelta?lastLocalBlockSerialNo=%d", lastLocalBlockSerialNo)
	resp, err := http.Get(URI)

	if err != nil {
		log.Fatalln(err)
		fmt.Printf("something went wrong when sending GET req to %s\n", URI)
	}

	var blockDelta []Database.Block

	bytes, err := readResp(resp)

	json.Unmarshal(bytes, &blockDelta)

	return blockDelta
}

//The following is done using POST,
//The header contain the address of the peer that is currently being accessed
//The body should contain the current state of the requesting node
func GetPeerState(peerAddr string) NodeState {
	httpposturl := "http://" + peerAddr + "/getState"

	currNodeState := GetNodeState()
	jsonData, err := json.Marshal(currNodeState)
	if err != nil {
		fmt.Println("Could not marshal current node state")
		panic(err)
	}

	resp, err := http.Post(httpposturl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Unable to get state")
		panic(err)
	}

	var peerNodeStateFromRequest NodeStateFromPostRequest
	var peerNodeState NodeState
	bytes, _ := readResp(resp)
	fmt.Println("Get State response")
	str := string(bytes)
	fmt.Println(str)
	json.Unmarshal(bytes, &peerNodeStateFromRequest)
	json.Unmarshal(bytes, &peerNodeState)
	//At this point the data recived should have been saved into peerNodeState

	var lh32 [32]byte
	for i := 0; i < 32; i++ {
		lh32[i] = peerNodeStateFromRequest.State.LatestHash[i]
	}
	peerNodeState.State.LatestHash = lh32

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

	http.HandleFunc("/blockDelta", func(w http.ResponseWriter, r *http.Request) {
		blockDeltaHandler(w, r, state)
	})
	fmt.Println("/blockDelta setup complete")

	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

func blockDeltaHandler(w http.ResponseWriter, r *http.Request, state *Database.State) {
	localBlockChain := Database.LoadBlockchain()
	serialNoParam := r.URL.Query().Get("lastLocalBlockSerialNo")
	var fromSerial int
	if serialNoParam == "" {
		fmt.Errorf("no serial number was provided in GET request\n")
		return
	}

	fromSerial, _ = strconv.Atoi(serialNoParam)
	delta := Database.GetBlockChainDelta(localBlockChain, fromSerial)

	writeResult(w, delta)
}

func getStateHandler(w http.ResponseWriter, r *http.Request, state *Database.State) {
	// Read the body containing the state of the node requesting
	getStateRequest := NodeState{}
	bytes, err := readReq(r)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bytes, &getStateRequest)

	//TODO: Do something with the peer state

	//Response: Send your current state
	nodeState := GetNodeState()

	fmt.Println(nodeState.PeerList)
	writeResult(w, nodeState)
}

func balancesHandler(w http.ResponseWriter, r *http.Request, state *Database.State) {
	writeResult(w, balancesResult{state.LatestHash, state.AccountBalances})
}

func transactionHandler(w http.ResponseWriter, r *http.Request, state *Database.State) {
	req := TxRequest{}
	bytes, err := readReq(r)
	if err != nil {
		return
	}
	json.Unmarshal(bytes, &req)

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
func readReq(r *http.Request) ([]byte, error) {
	reqJson, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to read request body. %s", err.Error())
	}

	return reqJson, nil
}

func readResp(r *http.Response) ([]byte, error) {
	reqJson, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to read reqsponse body. %s", err.Error())
	}

	return reqJson, nil
}
