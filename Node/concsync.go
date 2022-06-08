package node

import (
	db "blockchain/Database"
	shared "blockchain/Shared"
	"sort"
	"time"
)

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
		// get latest node data
		node := GetNode()

		// copy peersToCheck - is this needed anymore?
		peersToCheck := node.PeerSet.DeepCopy()
		noOfPeers := len(peersToCheck)

		// construct subsets for parallel goroutines to iterate through
		peerSubsets := constructSubsets(peersToCheck)

		// make channel for goroutines to write to, main routine to read from
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

		// get all data from the channels
		for i := 0; i < noOfPeers; i++ {
			pingResp := <-pingChannel
			if !pingResp.Ok {
				nodes = nodes[:len(nodes)-1]
				pings = pings[:len(pings)-1]
				continue
			}
			nodes[i] = <-nodeChannel
			pings[i] = pingResp
		}

		// TODO: following can only be done in last iteration. Listen for SIGTERM on main process?
		// close channels, since they will no longer be used
		// close(nodeChannel)
		// close(pingChannel)

		// add own node to nodes argument
		// nodes = append(nodes, node)

		// compute consensus
		consensusNode := computeConsensusNode(nodes) // gets node object that has consensus chain

		// match blockchain with consensus chain, newest blocks
		var deltaIdx int
		if len(node.ChainHashes) < len(consensusNode.ChainHashes) {
			deltaIdx = chainDiffIdx(node.ChainHashes, consensusNode.ChainHashes)
		} else {
			deltaIdx = chainDiffIdx(consensusNode.ChainHashes, node.ChainHashes)
		}

		// fetch peer blocks delta
		var peerBlocks []db.Block
		if deltaIdx != -1 {
			peerBlocks = GetPeerBlocks(consensusNode.Address, deltaIdx)
			clearConflictingSubchain(deltaIdx)
		}

		// apply the fetched blocks
		//    - TODO: should it make sure blocks are clear up until this point in own state?
		if len(peerBlocks) > 0 {
			for _, block := range peerBlocks {
				node.State.AddBlock(block)
			}
		}
		// TODO: or maybe clear them at this point since we do not know yet if the chain is accepted back then?

		// apply states from peers with newest chain
		tryApplyPeerStates(node, nodes)

		// compute new PeerSet based on top XX fastest pings
		pings = add2ndLevelPeers(pings, peersToCheck, nodes)
		newPeers := getNFastestPeers(pings, MAX_PEERS)

		// persist new peerset to file if there are any - otherwise, it might be because of bad connection
		if len(newPeers) > 0 {
			SavePeerSetAsJSON(newPeers, shared.PeerSetFile)
		}

		time.Sleep(20 * time.Second)
	}
}

func clearConflictingSubchain(deltaIdx int) {
	blockchain := db.LoadBlockchain()[:deltaIdx-1]
	// slicedBlockchain := blockchain
	db.SaveBlockchain(blockchain)
}

func tryApplyPeerStates(node Node, nodes []Node) {
	for _, peer := range nodes {
		if !Ping(peer.Address).Ok {
			continue
		}
		if len(peer.ChainHashes) < len(node.ChainHashes) {
			if chainsAgree(peer.ChainHashes, node.ChainHashes) {
				node.State.TryAddTransactions(peer.State.TxMempool)
			}
		} else {
			if chainsAgree(node.ChainHashes, peer.ChainHashes) {
				node.State.TryAddTransactions(peer.State.TxMempool)
			}
		}
	}
}

func add2ndLevelPeers(pings PingResponseList, peersToCheck PeerSet, nodes []Node) PingResponseList {
	for _, n := range nodes {
		for peer2 := range n.PeerSet {
			if !peersToCheck.Exists(peer2) && peer2 != getLocalIP() {
				pingRes := Ping(peer2)
				pings = append(pings, pingRes)
			}
		}
	}
	return pings
}

func getNFastestPeers(pings PingResponseList, amount int) PeerSet {
	sort.Sort(pings)
	ps := PeerSet{}
	for i, pingRes := range pings {
		if i >= amount {
			break
		}
		ps.Add(pingRes.Address)
	}
	return ps
}

// contract: c1 is the shorter, c2 is the longer chain
func chainDiffIdx(c1 []string, c2 []string) int {
	// if chains are identical, return -1
	if len(c1) == len(c2) && chainsAgree(c1, c2) {
		return -1
	}

	// find index where the two chains no longer agree
	for idx, h1 := range c1 {
		if c2[idx] != h1 {
			return idx
		}
	}

	// otherwise they agree, and it will always be from the last index of the shorter chain
	return len(c1)
}

// returns first node that contains the consensus chain (longest chain that most agree upon)
func computeConsensusNode(nodes []Node) Node {
	// make map for storing latest hash mapped to its node
	latestHash2Node := make(map[string]Node)

	// find unique chains (on last hash and serial no.) and store the no. of times they appear in different nodes and serialNo
	latestHashes := make(map[string]cPair)
	seenNodeAddresses := PeerSet{}
	for _, n := range nodes {
		if seenNodeAddresses.Exists(n.Address) || n.Address == "" {
			continue
		}
		seenNodeAddresses.Add(n.Address)
		latestHash := n.ChainHashes[len(n.ChainHashes)-1]
		if val, ok := latestHashes[latestHash]; ok {
			latestHashes[latestHash] = cPair{val.serialNo, val.count + 1}
		} else {
			latestHashes[latestHash] = cPair{n.State.LastBlockSerialNo, 1}
			// store first node with unique hash in map
			latestHash2Node[latestHash] = n
		}
	}
	// iterate unique hashes (pop each time)
	// for each, iterate all other unique hashes
	// if serialNo (block height) is greater on one, if they agree, add the count of the lower block height chain to the longer one

	// store how many nodes agree on chain
	agreeCount := make(map[string]int)
	for h1, cPair1 := range latestHashes {
		// remove to avoid duplicates
		if _, ok := agreeCount[h1]; !ok {
			agreeCount[h1] = latestHashes[h1].count
		}
		delete(latestHashes, h1)
		for h2, cPair2 := range latestHashes {
			if _, ok := agreeCount[h2]; !ok {
				agreeCount[h2] = latestHashes[h2].count
			}

			if cPair1.serialNo < cPair2.serialNo {
				c1 := latestHash2Node[h1].ChainHashes
				c2 := latestHash2Node[h2].ChainHashes
				if chainsAgree(c1, c2) {
					agreeCount[h2] = agreeCount[h1] + agreeCount[h2]
				}
			} else if cPair1.serialNo > cPair2.serialNo {
				c1 := latestHash2Node[h1].ChainHashes
				c2 := latestHash2Node[h2].ChainHashes
				if chainsAgree(c2, c1) {
					agreeCount[h1] = agreeCount[h1] + agreeCount[h2]
				}
			}
		}
	}
	maxAgreeHash := getMaxAgreeHash(agreeCount)

	return latestHash2Node[maxAgreeHash]
}

func getMaxAgreeHash(agreeCount map[string]int) string {
	var max int = 0
	var maxHash string
	for hash, count := range agreeCount {
		if count > max {
			max = count
			maxHash = hash
		}
	}
	return maxHash
}

// contract: c1 is the shorter chain, c2 is the longer
func chainsAgree(c1 []string, c2 []string) bool {
	// compare the chains at the latest hash in c1
	compIdx := len(c1) - 1
	return c1[compIdx] == c2[compIdx]
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
			nch <- Node{}
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
