package node

import (
	Database "blockchain/Database"
	"fmt"
	"time"
)

var peerListFile = "PeerList.json"

//Function that essentialy is the implemented version of our peer sync algorithm
func synchronization() {
	for {
		nodeState := GetNodeState() // Updates nodeState to get local changes in case any has been made

		// Creates a copy of the peers to check, to avoid weird looping as we append to the slice in the loop
		peersToCheck := make([]string, len(nodeState.PeerList))
		copy(peersToCheck, nodeState.PeerList)

		// Create a new slice of Addresses - Add those who are new and works
		var newAddresses []string

		for _, peer := range peersToCheck {
			if !Ping(peer) { // If the connection is lost, dont send to it
				continue
			} else {
				newAddresses = append(newAddresses, peer)
			}

			peerState := GetPeerState(peer)

			fmt.Println("Got peer state")
			fmt.Println(peerState)

			if peerState.State.LastBlockSerialNo > nodeState.State.LastBlockSerialNo {
				peerBlocks := GetPeerBlocks(peer, nodeState.State.LastBlockSerialNo)
				for _, block := range peerBlocks {
					nodeState.State.AddBlock(block)
				}
			}

			nodeState.State.TryAddTransactions(peerState.State.TxMempool)

			for _, peer2 := range peerState.PeerList {
				if !contains(nodeState.PeerList, peer2) && Ping(peer2) { // If the incomming address wasn't in the original list, add it to the new list of addresses
					newAddresses = append(newAddresses, peer2)
				}
			}
		}
		// Save any potential changes in the peer list
		Database.SavePeerListAsJSON(newAddresses, peerListFile)
		// Wait 20 seconds to continue
		time.Sleep(20 * time.Second)
	}
}

// Get the initial node state
func GetNodeState() NodeState {
	nodeState := NodeState{}
	nodeState.State = *Database.LoadState()
	nodeState.PeerList = GetPeerList()
	return nodeState
}

// Get the stored list of nodes
// If this hasn't been created before, create it using the bootstrap node
func GetPeerList() []string {
	peerList := Database.LoadPeerListFromJSON(peerListFile)
	if len(peerList) == 0 {
		peerList = append(peerList, bootstrapNode)
		Database.SavePeerListAsJSON(peerList, peerListFile)
	}
	return peerList
}
