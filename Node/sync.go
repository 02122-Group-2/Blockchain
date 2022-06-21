package node

import (
	shared "blockchain/Shared"
	"sort"
	"time"
)

// * file: Niels, s204503

var MAX_SUBSETS int = 4
var MAX_PEERS int = 40

type cPair struct {
	serialNo int
	count    int
}

// concurrent implementation of our synchronization algorithm, with a simple proof-of-work consensus algorithm
func concSynchronization() {
	for {
		shared.Log("Running sync")

		syncLoop()

		time.Sleep(20 * time.Second)
	}
}

func syncLoop() {
	// get latest node data
	node := GetNode()

	noOfPeers := len(node.PeerSet)

	// construct subsets for parallel goroutines to iterate through
	peerSubsets := constructSubsets(node.PeerSet)

	// make channels for goroutines to write to, main routine to read from
	nodeChannel := make(chan Node)
	pingChannel := make(chan PingResponse)

	// create goroutines for concurrent sync and assign channels
	for _, subset := range peerSubsets {
		go getNodesInPeerSet(subset, nodeChannel, pingChannel)
	}

	// read data out of goroutines through channel and store in Node slice
	nodes := make([]Node, noOfPeers)
	// store addresses mapped to their ping response latencies
	pings := make([]PingResponse, len(nodes))

	// read data from the subroutines' channels
	offset := 0
	for i := 0; i < noOfPeers; i++ {
		pingResp := <-pingChannel
		if !pingResp.Ok {
			nodes = nodes[:len(nodes)-1]
			pings = pings[:len(pings)-1]
			offset++
			continue
		}
		nodes[i-offset] = <-nodeChannel
		pings[i-offset] = pingResp
	}

	// TODO: following can only be done in last iteration. Listen for SIGTERM on main process?
	// close channels, since they will no longer be used
	// close(nodeChannel)
	// close(pingChannel)

	// add own node to collection and run consensus algorithm
	nodes = append(nodes, node)
	succesfullyApplied := handleConsensus(node, nodes)

	// apply all possible states from peers with newest chain
	tryApplyPeerStates(node, nodes)

	// compute new PeerSet based on top XX fastest pings
	newPeers := computeNewPeerSet(pings, node.PeerSet, nodes, succesfullyApplied)

	// persist new peerset to file if there are any - otherwise, it might be because of bad connection
	if len(newPeers) > 0 {
		PersistPeerSet(newPeers)
	}
}

// select which peers to keep for next cycle of sync, ranked on ping latency
func computeNewPeerSet(pings []PingResponse, ps PeerSet, nodes []Node, shouldExpandPeerset bool) PeerSet {
	// if last consensus chain was illegal, then the consensus node was removed. This flag is to prevent it from being immediately added again
	if shouldExpandPeerset {
		pings = add2ndLevelPeers(pings, ps, nodes)
	}
	newPeers := getNFastestPeers(pings, MAX_PEERS)
	return newPeers
}

// ping peer-of-peers to potentially expand own peer set
func add2ndLevelPeers(pings PingResponseList, peersToCheck PeerSet, nodes []Node) PingResponseList {
	localIp := getLocalIP()
	for _, n := range nodes {
		for peer2 := range n.PeerSet {
			if !peersToCheck.Exists(peer2) && peer2 != localIp {
				pingRes := Ping(peer2)
				if pingRes.Ok {
					pings = append(pings, pingRes)
				}
			}
		}
	}
	return pings
}

// get fixed sized peer set ranked by ping latency (low to high)
func getNFastestPeers(pings PingResponseList, amount int) PeerSet {
	sort.Sort(pings)
	ps := PeerSet{}
	for _, pingRes := range pings {
		if len(ps) >= amount {
			break
		}
		ps.Add(pingRes.Address)
	}
	return ps
}

func getNodesInPeerSet(ps PeerSet, nch chan Node, pch chan PingResponse) {
	for peer := range ps {
		pingRes := Ping(peer)

		// TODO: configure timeout for GetPeerState?
		if pingRes.Ok {
			pch <- pingRes
			nch <- GetPeerState(peer)
		} else {
			pch <- PingResponse{"nil", false, -1}
			// nch <- Node{}
		}
	}
}

func constructSubsets(peersToCheck PeerSet) []PeerSet {
	psLen := len(peersToCheck)
	noSubsets := min(psLen, MAX_SUBSETS)

	subsets := make([]PeerSet, noSubsets)

	for peer := range peersToCheck {
		idx := len(peersToCheck) % noSubsets
		if subsets[idx] == nil {
			subsets[idx] = PeerSet{}
		}
		subsets[idx].Add(peer)
		peersToCheck.Remove(peer)
	}

	return subsets
}

func min(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}
