package node

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	t.Log("begin init test")

	Run()

	//Database.ResetTest()
}

func TestGetPeerState(t *testing.T) {
	t.Log("begin get peer state test")

	/*
		err := startNode()
		if err != nil {
			t.Errorf("Node could not start")
		}
		t.Log("sucessfully loaded the current state")
	*/

	nodeState := GetPeerState("localhost:8080")

	if nodeState.PeerList == nil {
		t.Errorf("Peer list is nil")
	}
	if nodeState.State.AccountBalances == nil {
		t.Errorf("State balances is nil")
	}
	fmt.Println(nodeState.State)

}
