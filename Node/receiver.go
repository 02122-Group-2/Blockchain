package node

import (
	Database "blockchain/Database"
	shared "blockchain/Shared"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func blockDeltaHandler(w http.ResponseWriter, r *http.Request, state *Database.State) {
	localBlockChain := Database.LoadBlockchain()
	serialNoParam := r.URL.Query().Get("lastLocalBlockSerialNo")
	var fromSerial int
	if serialNoParam == "" {
		fmt.Println(fmt.Errorf("no serial number was provided in GET request"))
		return
	}

	fromSerial, _ = strconv.Atoi(serialNoParam)
	delta := Database.GetBlockChainDelta(localBlockChain, fromSerial)

	writeResult(w, delta)
}

//Function used to get the state of a peer node
func getStateHandler(w http.ResponseWriter, r *http.Request, state *Database.State) {
	//Response: Get your own state to send to the one requesting it
	node := GetNode()

	// Read the body containing the state of the node requesting
	getStateRequest := Node{}
	bytes, err := readReq(r)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bytes, &getStateRequest)

	// load peer set from file and union it with the peer set from the incoming request
	currentPeerSet := LoadPeerSetFromJSON(shared.PeerSetFile)
	currentPeerSet.UnionWith(getStateRequest.PeerSet)
	SavePeerSetAsJSON(currentPeerSet, shared.PeerSetFile)

	// fmt.Println(node.PeerSet)
	writeResult(w, node)
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

	// fmt.Println(transaction)

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

//Reading the request from client
func readReq(r *http.Request) ([]byte, error) {
	reqJson, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to read request body. %s", err.Error())
	}

	return reqJson, nil
}
