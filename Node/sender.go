package node

import (
	Database "blockchain/Database"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

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

//Reading the server response
func readResp(r *http.Response) ([]byte, error) {
	reqJson, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to read reqsponse body. %s", err.Error())
	}

	return reqJson, nil
}
