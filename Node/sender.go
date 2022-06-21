package node

import (
	Database "blockchain/Database"
	shared "blockchain/Shared"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// * Niels, s204503
func GetPeerBlocks(peerAddr string, deltaIdx int) []Database.Block {
	lastBlockSerialNo := deltaIdx - 1
	URI := fmt.Sprintf("http://"+peerAddr+"/blockDelta?lastLocalBlockSerialNo=%d", lastBlockSerialNo)
	resp, err := http.Get(URI)

	if err != nil {
		log.Fatalln(err)
		fmt.Printf("something went wrong when sending GET req to %s\n", URI)
	}

	var blockDelta []Database.Block

	bytes, _ := readResp(resp)

	json.Unmarshal(bytes, &blockDelta)

	return blockDelta
}

// * Emilie, s204471
//The following is done using POST,
//The header contain the address of the peer that is currently being accessed
//The body should contain the current state of the requesting node
func GetPeerState(peerAddr string) Node {
	httpposturl := "http://" + peerAddr + "/getState"

	currNode := GetNode()
	jsonData, err := json.Marshal(currNode)
	if err != nil {
		fmt.Println("Could not marshal current node state")
		panic(err)
	}

	resp, err := http.Post(httpposturl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Unable to get state")
		panic(err)
	}

	var peerNodeFromRequest NodeFromPostRequest
	var peerNode Node
	bytes, _ := readResp(resp)
	shared.Log(fmt.Sprintf("Get State response from %s", peerAddr))
	// str := string(bytes)
	// fmt.Println(str)
	json.Unmarshal(bytes, &peerNodeFromRequest)
	json.Unmarshal(bytes, &peerNode)
	//At this point the data recived should have been saved into peerNode

	var lh32 [32]byte
	for i := 0; i < 32; i++ {
		lh32[i] = peerNodeFromRequest.State.LatestHash[i]
	}
	peerNode.State.LatestHash = lh32

	peerNode.ChainHashes = peerNodeFromRequest.ChainHashes

	return peerNode
}

// * Asger, s204435
//Reading the server response
func readResp(r *http.Response) ([]byte, error) {
	reqJson, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to read reqsponse body. %s", err.Error())
	}

	return reqJson, nil
}
