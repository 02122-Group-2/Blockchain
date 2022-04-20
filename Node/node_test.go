package node

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	t.Log("begin init test")

	Run() //change to yout own path when testing

	//Database.ResetTest()
}

type MockNodeConnectionManager struct {
	peers PeerList

	node2status map[string]NodeResponse // maps node to response
}

// onlinePercentage 0-100% -> {0.0..1.0}
func (ncm *MockNodeConnectionManager) seedNodeData(onlinePercentage float32) {
	if onlinePercentage < 0.0 || onlinePercentage > 1.0 {
		panic("onlinePercentage should be a number from 0.0 to 1.0")
	}

	peerList, err := LoadPeers()
	if err != nil {
		panic(err.Error())
	}
	ncm.peers = peerList

	if ncm.node2status == nil {
		ncm.node2status = map[string]NodeResponse{}
	}

	peerCount := len(ncm.peers.Peers)
	for i, peer := range ncm.peers.Peers {
		i++
		percentile := float32(i) / float32(peerCount)
		var nodeStatus Status
		if percentile > onlinePercentage {
			nodeStatus = NO_CONN
		} else {
			nodeStatus = READY
		}

		res := fmt.Sprintf("test@%d", i)
		ncm.node2status[peer.FQDN] = NodeResponse{
			Status:   nodeStatus,
			Response: res,
		}
	}
}

func (ncm *MockNodeConnectionManager) ensureStatusOnMockDB(amount int, status Status) {
	peers := ncm.peers.Peers
	if amount > len(peers) {
		panic("Cannot set status on more nodes than are in the DB")
	}
	for i := 0; i < amount; i++ {
		cur := ncm.node2status[peers[i].FQDN]
		cur.Status = status
		ncm.node2status[peers[i].FQDN] = cur
	}
}

func (ncm MockNodeConnectionManager) GetHeartBeat(fqdn string) NodeResponse {
	return ncm.node2status[fqdn]
}

func (ncm MockNodeConnectionManager) FetchStateData(fqdn string) ([]byte, error) {
	return []byte{}, nil
}

var mockNode Node
var mockNcm MockNodeConnectionManager

func (ncm MockNodeConnectionManager) printMockDB() string {
	var s string
	for _, obj := range ncm.peers.Peers {
		s += "\n" + fmt.Sprintf("%s: %+v", obj.Alias, ncm.node2status[obj.FQDN])
	}
	return s
}

func networkStatusToString(online int, ready int) string {
	return fmt.Sprintf("\nPeer Status\nOnline: %d Ready: %d", online, ready)
}

func TestMockNodeConnectionManager(t *testing.T) {
	mockNcm.seedNodeData(0.0)
	mockNode = Node{NCM: mockNcm, FQDN: "nil", Peers: nil}
	t.Log(mockNcm.printMockDB())
	networkStatus, err := mockNode.CheckNetworkStatus()
	if err != nil {
		panic(err.Error())
	}
	ready, busy, failed, online := networkStatus.ExtractTuple()
	t.Log(networkStatusToString(online, ready))

	if online != 0 && ready != 0 {
		panic("mockNcm should be seeded with 0 online peers")
	}

	mockNcm.seedNodeData(1.0)
	mockNode.NCM = mockNcm
	t.Log(mockNcm.printMockDB())
	networkStatus, err = mockNode.CheckNetworkStatus()
	if err != nil {
		panic(err.Error())
	}
	//nolint:staticcheck
	ready, busy, failed, online = networkStatus.ExtractTuple()
	t.Log(networkStatusToString(online, ready))

	noOfPeers := len(mockNcm.peers.Peers)
	if online != noOfPeers && ready != noOfPeers {
		panic(fmt.Sprintf("mockNcm should be seeded with all online and ready peers. Actual: (online) [%d/%d] (ready) [%d/%d]", online, noOfPeers, ready, noOfPeers))
	}

	amountFailed, amountBusy := 2, 1
	mockNcm.ensureStatusOnMockDB(amountFailed, FAILED)
	mockNcm.ensureStatusOnMockDB(amountBusy, BUSY)
	t.Log(mockNcm.printMockDB())
	networkStatus, err = mockNode.CheckNetworkStatus()
	if err != nil {
		panic(err.Error())
	}
	ready, busy, failed, online = networkStatus.ExtractTuple()
	t.Log(networkStatusToString(online, ready))
	if online != noOfPeers || ready != (noOfPeers-amountFailed) {
		panic(fmt.Sprintf("mockNcm should currently have %d ready peers at beginning of peer list", (noOfPeers - 2)))
	}
	if busy != amountBusy {
		panic(fmt.Sprintf("mockNcm should have %d busy peers", amountBusy))
	}
	if failed != (amountFailed - amountBusy) {
		panic(fmt.Sprintf("mockNcm should have %d failed peers (due to one of the failed being set as busy in previous assignments)", (amountFailed - amountBusy)))
	}
}
