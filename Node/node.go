package node

// Business logic for interacting between nodes
// Uses an implementation of NodeConnectionManager interface
// to perform networking calls (or other types of queries, such as in MockConnectionManager)

func Fetch2ndLevelPeerList() {

}

func SendState() {

}

func FetchState() {

}

// TODO: pass ncm as argument? Or have singleton declared somewhere?
// Check whether node has any registered peers and whether they are online and ready to communicate
func CheckNetworkStatus(ncm NodeConnectionManager) (int, int, error) {
	peers, err := LoadPeers()
	if err != nil {
		panic(err)
	}
	noOfPeers := countPeers(peers)

	if noOfPeers == 0 {
		return 0, 0, nil
	}

	var healthyPeers, onlinePeers int
	for _, peer := range peers.Peers {
		nodeRes := ncm.GetHeartBeat(peer.FQDN)
		statusCode := nodeRes.Status
		if statusCode == READY {
			healthyPeers++
			onlinePeers++
		} else if statusCode != NO_CONN {
			onlinePeers++
		}
	}

	return onlinePeers, healthyPeers, err
}
