package node

import (
	Database "blockchain/Database"
	"fmt"
	"time"
)

//Function that essentialy is the implemented version of our peer sync algorithm
func synchronization() {
	nodeState := GetNodeState()

	for {
		updateNodeState(&nodeState) // Updates nodeState to get local changes in case any has been made

		// Creates a copy of the peers to check, to avoid weird looping as we append to the slice in the loop
		peersToCheck := make([]string, len(nodeState.PeerList))
		copy(peersToCheck, nodeState.PeerList)

		for _, peer := range peersToCheck {

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
				if !contains(nodeState.PeerList, peer2) && Ping(peer2) {
					nodeState.PeerList = append(nodeState.PeerList, peer2)
				}
			}
			time.Sleep(40 * time.Second)
		}
	}
}

// Updated the current nodeState.State field to get latest local changes, in case the user has added local transactions.
func updateNodeState(currentNodeState *NodeState) {
	currentNodeState.State = *Database.LoadState()
}
