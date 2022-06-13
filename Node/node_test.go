package node

import (
	shared "blockchain/Shared"
	"fmt"
	"os"
	"testing"
)

func TestCreatePeerSet(t *testing.T) {
	addr := getFirstPeerInPeerset()
	ps := PeerSet{}
	localhost := addr
	ps.Add(localhost)
	peerSetTestFile := fmt.Sprintf("%s_test.json", shared.PeerSetFile[:(len(shared.PeerSetFile)-len(".json"))])
	SavePeerSetAsJSON(ps, peerSetTestFile)

	correctChecksum := "31f74234299b520c9a680a8022b7b0b24f707668c30f34be152a334d015dec0e"

	realChecksum := shared.GetChecksum(shared.LocatePersistenceFile(peerSetTestFile, "test_data"))

	t.Log(realChecksum)

	// t.Logf("Checksums:\n%x\n%x\n\n", correctChecksum, realChecksum)
	if realChecksum != correctChecksum {
		panic(fmt.Sprintf("%s was not created correctly", peerSetTestFile))
	}

	os.Remove(shared.LocatePersistenceFile(peerSetTestFile, ""))
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

	addr := getFirstPeerInPeerset()
	nodeState := GetPeerState(addr)

	if nodeState.PeerSet == nil {
		t.Errorf("Peer list is nil")
	}
	if nodeState.State.AccountBalances == nil {
		t.Errorf("State balances is nil")
	}
	// fmt.Println(nodeState.State)
}

func TestGetPeerBlocks(t *testing.T) {
	addr := getFirstPeerInPeerset()
	res := GetPeerBlocks(addr, 0)

	fmt.Println(res)
}

func TestPingActiveConnection(t *testing.T) {
	addr := getFirstPeerInPeerset()
	pingRes := Ping(addr)
	t.Logf("Latency: %d", pingRes.Latency)
	if !pingRes.Ok {
		t.Error("Connection not active")
	}
}

func getFirstPeerInPeerset() string {
	peers := LoadPeerSetFromJSON(shared.PeerSetFile)
	var addr string
	for p := range peers {
		addr = p
		break
	}
	return addr
}
