package node

import (
	shared "blockchain/Shared"
	"fmt"
	"os"
	"testing"
)

func TestCreatePeerSet(t *testing.T) {
	ps := PeerSet{}
	localhost := "localhost:8080"
	ps.Add(localhost)
	peerSetTestFile := fmt.Sprintf("%s_test.json", shared.PeerSetFile[:(len(shared.PeerSetFile)-len(".json"))])
	SavePeerSetAsJSON(ps, peerSetTestFile)

	correctChecksum := "a279eb83cba534c181c993c1a89989982b9c9862ea36a1548e0a78e0be851f69"

	realChecksum := shared.GetChecksum(shared.Locate(peerSetTestFile))

	t.Logf("Checksums:\n%x\n%x\n\n", correctChecksum, realChecksum)
	if realChecksum != correctChecksum {
		panic(fmt.Sprintf("%s was not created correctly", peerSetTestFile))
	}

	os.Remove(shared.Locate(peerSetTestFile))
}

// this test is not an actual unit test.
// it is merely for exploratory tests, e.g. starting a node and debugging
func TestRun(t *testing.T) {
	t.Log("begin run test")

	err := Run()
	if err != nil {
		t.Log("Could not run node")
		t.Fail()
	}
}
func TestGetPeerState(t *testing.T) {
	t.Log("begin get peer state test")

	nodeState := GetPeerState("192.168.0.106:8080")

	if nodeState.PeerSet == nil {
		t.Errorf("Peer list is nil")
	}
	if nodeState.State.AccountBalances == nil {
		t.Errorf("State balances is nil")
	}
	fmt.Println(nodeState.State)
}

func TestGetPeerBlocks(t *testing.T) {
	res := GetPeerBlocks("192.168.0.106:8090", 0)

	fmt.Println(res)
}

func TestPingActiveConnection(t *testing.T) {
	addr := "localhost:8080"
	pingRes := Ping(addr)
	t.Logf("Latency: %d", pingRes.Latency)
	if !pingRes.Ok {
		t.Error("Connection not active")
	}
}
