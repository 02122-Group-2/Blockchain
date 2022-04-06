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
