package node

import (
	"fmt"
	"testing"
	shared "blockchain/Shared"
)

func TestCreatePeerSet(t *testing.T) {
	ps := PeerSet{}
	// legalIP := "192.168.0.1:8080"
	localhost := "localhost:8080"
	ps.Add(localhost)
	SavePeerSetAsJSON(ps, shared.PeerSetFile)
}

func TestRun(t *testing.T) {
	t.Log("begin run test")

	err := Run()
	if err != nil {
		t.Log("Could not run node")
		t.Fail()
	}
	//Check for files

}
func TestGetPeerState(t *testing.T) {
	t.Log("begin get peer state test")

	// err := startNode()
	// if err != nil {
	// 	t.Errorf("Node could not start")
	// }
	// t.Log("sucessfully loaded the current state")

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
	// addr := "10.209.222.2:8080"
	pingRes := Ping(addr)
	t.Logf("Latency: %d", pingRes.Latency)
	if !pingRes.Ok {
		t.Error("Connection not active")
	}
}
