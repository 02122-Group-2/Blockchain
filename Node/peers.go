package node

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PeerList struct {
	Peers []Peer `json:"Peers"`
}

type Peer struct {
	Alias string `json:"Alias"`
	FQDN  string `json:"FQDN"`
}

func LoadPeers() (PeerList, error) {
	currWD, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(filepath.Join(currWD, "./Peers.json"))
	if err != nil {
		panic(fmt.Sprintf("Error loading \"Peers\"-file: %s", err.Error()))
	}

	var loaded_peers PeerList
	err = json.Unmarshal(data, &loaded_peers)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling JSON: %s", err.Error()))
	}

	return loaded_peers, err
}

// Check whether node has any registered peers and whether they are online and ready to communicate
func CheckNetworkStatus() (int, error) {
	peers, err := LoadPeers()
	if err != nil {
		panic(err)
	}
	noOfPeers := countPeers(peers)

	if noOfPeers == 0 {
		return 0, nil
	}

	var healthyPeers, onlinePeers int
	for _, peer := range peers.Peers {
		res, _ := Ping(peer.FQDN)
		if res == 200 { // HTTP response code for OK
			healthyPeers++
			onlinePeers++
		}
		if res >= 300 { // HTTP response code for various errors
			onlinePeers++
		}
	}

	if onlinePeers == 0 {
		return 0, nil
	}

	statusCode := (int)(healthyPeers / onlinePeers * 5)

	return statusCode, err
}

func countPeers(peers PeerList) int {
	c := 0
	for _, peer := range peers.Peers {
		if !isBootstrapNode(peer) {
			c++
		}
	}
	return c
}

func isBootstrapNode(p Peer) bool {
	return strings.Contains(p.Alias, "bootstrap")
}

// Ping peers to check if online and ready to communicate
func Ping(fqdn string) (int, string) {
	return 200, "Ready"
}
