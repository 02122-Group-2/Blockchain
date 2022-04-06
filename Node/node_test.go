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

func (ncm MockNodeConnectionManager) FetchStateData(fqdn string) (int, []byte) {
	return -1, []byte{}
}

var mockNcm MockNodeConnectionManager

func (ncm MockNodeConnectionManager) printMockDB() string {
	var s string
	for _, obj := range ncm.peers.Peers {
		s += "\n" + fmt.Sprintf("%s: %+v", obj.Alias, ncm.node2status[obj.FQDN])
	}
	return s
}

func networkStatusToString(online int, healthy int) string {
	return fmt.Sprintf("\nPeer Status\nOnline: %d Healthy: %d", online, healthy)
}

func TestMockNodeConnectionManager(t *testing.T) {
	mockNcm.seedNodeData(0.0)
	t.Log(mockNcm.printMockDB())
	online, healthy, err := CheckNetworkStatus(mockNcm)
	if err != nil {
		panic(err.Error())
	}
	t.Log(networkStatusToString(online, healthy))

	if online != 0 && healthy != 0 {
		panic("mockNcm should be seeded with 0 online peers")
	}

	mockNcm.seedNodeData(1.0)
	t.Log(mockNcm.printMockDB())
	online, healthy, err = CheckNetworkStatus(mockNcm)
	if err != nil {
		panic(err.Error())
	}
	t.Log(networkStatusToString(online, healthy))

	noOfPeers := len(mockNcm.peers.Peers)
	if online != noOfPeers && healthy != noOfPeers {
		panic(fmt.Sprintf("mockNcm should be seeded with all online and healthy peers. Actual: (online) [%d/%d] (healthy) [%d/%d]", online, noOfPeers, healthy, noOfPeers))
	}

	mockNcm.ensureStatusOnMockDB(2, FAILED)
	mockNcm.ensureStatusOnMockDB(1, BUSY)
	t.Log(mockNcm.printMockDB())
	online, healthy, err = CheckNetworkStatus(mockNcm)
	if err != nil {
		panic(err.Error())
	}
	t.Log(networkStatusToString(online, healthy))

	if online != noOfPeers && healthy != (noOfPeers-2) {
		panic("mockNcm should currently have 2 unhealthy peers at beginning of peer list")
	}
}
