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

// * Niels, s204503
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

	writeResult(w, r, delta)
}

// * Emilie, s204471
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
	// ! Should this happen? I think we should remove this functionality, since it interferes with consensus
	currentPeerSet := LoadPeerSetFromJSON(shared.PeerSetFile)
	currentPeerSet.UnionWith(getStateRequest.PeerSet)
	SavePeerSetAsJSON(currentPeerSet, shared.PeerSetFile)

	// fmt.Println(node.PeerSet)
	writeResult(w, r, node)
}

// * Emilie, s204471
func balancesHandler(w http.ResponseWriter, r *http.Request, state *Database.State) {
	writeResult(w, r, balancesResult{state.LatestHash, state.AccountBalances})
}

// * Asger, s204435
//Writing the result from the server
func writeResult(w http.ResponseWriter, r *http.Request, content interface{}) {
	contentJson, err := json.Marshal(content)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(contentJson)
	shared.Log(fmt.Sprintf("Server response sent to %s", r.RemoteAddr))
}

// * Asger, s204435
//Reading the request from client
func readReq(r *http.Request) ([]byte, error) {
	reqJson, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to read request body. %s", err.Error())
	}

	return reqJson, nil
}
