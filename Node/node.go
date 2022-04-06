package node

import "fmt"

// Business logic for interacting between nodes
// Uses an implementation of NodeConnectionManager interface
// to perform networking calls (or other types of queries, such as in MockConnectionManager)

// interface describes contracts that should be fulfilled in order for the node to be able to interact
// with peer nodes. Ideally, when this interface (w.i.p.) is implemented, we should be able to create
// the higher level algorithms such as sync, mining, consensus
type INode interface {
	Fetch2ndLevelPeerList()
	SendState()
	GetState()
	CheckNetworkStatus()
}

type Node struct {
	NCM   NodeConnectionManager
	FQDN  string
	Peers *PeerList
}

type NetworkStatusDTO struct {
	ReadyPeers  int
	BusyPeers   int
	FailedPeers int
	OnlinePeers int
}

// (integer) amount of (ready, busy, failed, online) peers
func CreateNetworkStatusDTO(ready int, busy int, failed int, online int) NetworkStatusDTO {
	return NetworkStatusDTO{ReadyPeers: ready, BusyPeers: busy, FailedPeers: failed, OnlinePeers: online}
}

// returns 4-tuple of (ready, busy, failed, online) peers
func (nws NetworkStatusDTO) ExtractTuple() (int, int, int, int) {
	return int(nws.ReadyPeers), int(nws.BusyPeers), int(nws.FailedPeers), int(nws.OnlinePeers)
}

func (node *Node) Fetch2ndLevelPeerList() {

}

func (node *Node) SendState() {

}

func (node *Node) FetchState(fqdn string) {
	stateData, err := node.NCM.FetchStateData(fqdn)
	println("%+v %s", stateData, err)
}

// pass ncm as argument?
// or have singleton declared somewhere?
// or like now, where it is called from node, which uses its own ncm to handle data context
// Check whether node has any registered peers and whether they are online and ready to communicate
func (node *Node) CheckNetworkStatus() (NetworkStatusDTO, error) {
	peers, err := LoadPeers()
	if err != nil {
		panic(err)
	}
	noOfPeers := countPeers(peers)

	if noOfPeers == 0 {
		return CreateNetworkStatusDTO(0, 0, 0, 0), nil
	}

	var readyPeers, busyPeers, failedPeers, onlinePeers int
	for _, peer := range peers.Peers {
		nodeRes := node.NCM.GetHeartBeat(peer.FQDN)
		switch statusCode := nodeRes.Status; statusCode {
		case READY:
			readyPeers++
			onlinePeers++
		case BUSY:
			busyPeers++
			onlinePeers++
		case FAILED:
			failedPeers++
			onlinePeers++
		case NO_CONN:
			{
			}
		default:
			panic(fmt.Sprintf("Status Code %d should not be able to be returned from heartbeat check", statusCode))
		}
	}

	return CreateNetworkStatusDTO(readyPeers, busyPeers, failedPeers, onlinePeers), err
}
