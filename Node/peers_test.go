package node

import (
	"encoding/json"
	"testing"
)

func createPopulatedPeerList() PeerList {
	var peers PeerList
	data := []byte(`{"Peers": [{"Alias": "bootstrap-0", "FQDN": "127.0.0.1:13669"},{"Alias": "bootstrap-1","FQDN": "127.0.0.1:13670"},{"Alias": "niller-node1","FQDN": "10.209.184.241:8080"}]}`)
	err := json.Unmarshal(data, &peers)
	if err != nil {
		panic(err)
	}
	return peers
}

func createEmptyPeerList() PeerList {
	var peers PeerList
	return peers
}

func TestPrintPeerLists(t *testing.T) {
	t.Logf("Populated:%s\n", createPopulatedPeerList())
	t.Logf("Empty:%s\n", createEmptyPeerList())
}

func TestLoadPeers(t *testing.T) {
	loaded_peers, err := LoadPeers()

	if err != nil {
		t.Error(err)
	}

	prettyJson, _ := json.MarshalIndent(loaded_peers, "", "  ")
	t.Logf("Loaded Peers:\n%s", prettyJson)
}
