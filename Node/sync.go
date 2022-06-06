package node

import (
	Database "blockchain/Database"
	shared "blockchain/Shared"
	"fmt"
	"net"
	"os"
	"regexp"
	"time"
)

var peerSetFile = "PeerSet.json"

//Function that essentialy is the implemented version of our peer sync algorithm
func synchronization() {
	for {
		// updates Node to get local changes in case any has been made
		node := GetNode()

		// make copy of initial peers to avoid weird behavior when modifying the slice
		peersToCheck := node.PeerSet.DeepCopy()

		// slice for storing all active connections from current sync iteration
		newPeers := PeerSet{}

		// sync with each peer, first blocks then states
		for peer := range peersToCheck {
			newPeers = node.syncPeer(peer, newPeers)
		}

		// Persist the updated peer Set
		SavePeerSetAsJSON(newPeers, peerSetFile)

		// Wait 20 seconds before running next sync iteration
		time.Sleep(20 * time.Second)
	}
}

func (node Node) syncPeer(peer string, newPeers PeerSet) PeerSet {
	// if we cannot connect to a peer, skip it and don't append it
	if !Ping(peer).ok {
		return nil
	} else {
		newPeers.Add(peer)
	}

	peerState := GetPeerState(peer)

	fmt.Println("Got peer state")
	fmt.Println(peerState)

	peerHasNewerBlock := peerState.State.LastBlockSerialNo > node.State.LastBlockSerialNo
	if peerHasNewerBlock {
		peerBlocks := GetPeerBlocks(peer, node.State.LastBlockSerialNo)
		for _, block := range peerBlocks {
			node.State.AddBlock(block)
		}
	}

	node.State.TryAddTransactions(peerState.State.TxMempool)

	reachableIPs := PeerSet{}
	for peer2 := range peerState.PeerSet {
		if Ping(peer2).ok { // If the incoming address wasn't in the original list, add it to the new list of addresses
			reachableIPs.Add(peer2)
		}
	}

	return Union(newPeers, reachableIPs)
}

// Get the initial node state
func GetNode() Node {
	node := Node{}
	node.Address = getLocalIP()
	node.State = *Database.LoadState()
	node.PeerSet = GetPeerSet()
	node.ChainHashes = Database.GetLocalChainHashes(node.State, 0)
	return node
}

// Get the stored set of nodes
// If this hasn't been created before, create it using the bootstrap node
func GetPeerSet() PeerSet {
	ps := LoadPeerSetFromJSON(peerSetFile)
	if len(ps) == 0 {
		ps.Add(bootstrapNode)
	}
	return ps
}

func Ping(peerAddr string) PingResponse {
	if !legalIpAddress(peerAddr) {
		return PingResponse{"nil", false, -1}
	}

	startTime := shared.MakeTimestamp()
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", peerAddr, timeout)
	if err != nil {
		fmt.Println("Site unreachable, error: ", err)
		return PingResponse{"nil", false, -1}
	}
	endTime := shared.MakeTimestamp()
	latency := endTime - startTime
	conn.Close()
	return PingResponse{peerAddr, true, latency}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func legalIpAddress(addr string) bool {
	regexIPwithPort := "^(localhost|([0-9]{1,3}.){1,3}([0-9]{1,3})):([0-9]{4,5})$"
	match, _ := regexp.MatchString(regexIPwithPort, addr)
	return match
}

func getLocalIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return fmt.Sprint(ipv4) + ":8080"
		}
	}
	return "localhost:8080"
}
