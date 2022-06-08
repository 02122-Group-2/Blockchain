package node

import (
	Database "blockchain/Database"
	shared "blockchain/Shared"
	"fmt"
	"net"
	"time"
)

//Function that essentialy is the implemented version of our peer sync algorithm
func synchronization() {
	for {
		// updates Node to get local changes in case any has been made
		node := GetNode()

		// make copy of initial peers to avoid weird behavior when modifying the slice
		peersToCheck := node.PeerSet.DeepCopy()

		// slice for storing all active connections from current sync iteration
		newPeers := PeerSet{}

		for peer, _ := range peersToCheck {
			newPeers = node.syncPeer(peer, newPeers)
		}
		// Persist the updated peer Set
		SavePeerSetAsJSON(newPeers, shared.PeerSetFile)

		// Wait 20 seconds before running next sync iteration
		time.Sleep(20 * time.Second)
	}
}

func (node Node) syncPeer(peer string, newPeers PeerSet) PeerSet {
	// if we cannot connect to a peer, skip it and don't append it
	if !Ping(peer) {
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
	for peer2, _ := range peerState.PeerSet {
		if Ping(peer2) { // If the incoming address wasn't in the original list, add it to the new list of addresses
			reachableIPs.Add(peer2)
		}
	}

	return Union(newPeers, reachableIPs)
}

// Get the initial node state
func GetNode() Node {
	node := Node{}
	node.State = *Database.LoadState()
	node.PeerSet = GetPeerSet()
	return node
}

// Get the stored set of nodes
// If this hasn't been created before, create it using the bootstrap node
func GetPeerSet() PeerSet {
	ps := LoadPeerSetFromJSON(shared.PeerSetFile)
	if len(ps) == 0 {
		ps.Add(bootstrapNode)
	}
	return ps
}

func Ping(peerAddr string) bool {
	if !shared.LegalIpAddress(peerAddr) {
		return false
	}

	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", peerAddr, timeout)
	if err != nil {
		//fmt.Println("Site unreachable, error: ", err)
		return false
	}
	conn.Close()
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}


