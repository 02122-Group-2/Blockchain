package node

import (
	db "blockchain/Database"
	shared "blockchain/Shared"
	"fmt"
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

	// add own node to collection and run consensus algorithm
	nodes = append(nodes, node)
	handleConsensus(node, nodes)

	// apply all possible states from peers with newest chain
	tryApplyPeerStates(node, nodes)

	// compute new PeerSet based on top XX fastest pings
	newPeers := computeNewPeerSet(pings, node.PeerSet, nodes)

	// persist new peerset to file if there are any - otherwise, it might be because of bad connection
	if len(newPeers) > 0 {
		PersistPeerSet(newPeers)
	}
}

// Simple majority consensus algorithm
func handleConsensus(node Node, nodes []Node) {
	// gets node object that has consensus chain, i.e. longest chain that most nodes agree on
	consensusNode := computeConsensusNode(nodes)

	// compute index where chains no longer agree
	var deltaIdx int
	if len(node.ChainHashes) < len(consensusNode.ChainHashes) {
		deltaIdx = chainDiffIdx(node.ChainHashes, consensusNode.ChainHashes)
	} else {
		deltaIdx = chainDiffIdx(consensusNode.ChainHashes, node.ChainHashes)
	}

	// match blockchain with consensus chain, newest blocks
	peerBlocks := fetchConsensusChainDelta(consensusNode, deltaIdx)
	if len(peerBlocks) > 0 {
		// TODO: validate all received blocks before clearing and applying
		// if local chain has blocks that are conflicting at some point with the consensus chain, these must be cleared
		clearConflictingSubchain(deltaIdx) // Is this necessary since the same is performed in the Recomputestate function call?

		// state must match snapshot from before applying the last block before deltaIdx
		node.State.RecomputeState(deltaIdx)
		for _, block := range peerBlocks {
			blockErr := node.State.AddBlock(block)
			if blockErr != nil {
				fmt.Println(blockErr.Error())
			}
		}
	}
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
	for h1 := range latestHashes {
		// remove to avoid duplicates
		if _, ok := agreeCount[h1]; !ok {
			agreeCount[h1] = latestHashes[h1].count
		}
		delete(latestHashes, h1)
		for h2 := range latestHashes {
			if _, ok := agreeCount[h2]; !ok {
				agreeCount[h2] = latestHashes[h2].count
			}

			if chainsAgree(latestHash2Node[h1].ChainHashes, latestHash2Node[h2].ChainHashes) {
				agreeCount[h2] = agreeCount[h1] + agreeCount[h2]
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

// Given two lists of hashes, check that the last element for the shortest list is equal to the hash at the same location for the second list
func chainsAgree(c1 []string, c2 []string) bool {
	// Get the location of the last hash in the shortest list
	compIdx := min(len(c1), len(c2)) - 1

	// compare the chains at the latest hash for the shortest list
	return c1[compIdx] == c2[compIdx]
}

// select which peers to keep for next cycle of sync, ranked on ping latency
func computeNewPeerSet(pings []PingResponse, ps PeerSet, nodes []Node) PeerSet {
	pings = add2ndLevelPeers(pings, ps, nodes)
	newPeers := getNFastestPeers(pings, MAX_PEERS)
	return newPeers
}

// fetch difference in blocks between own chain and the one agreed upon by consensus algorithm
func fetchConsensusChainDelta(consensusNode Node, deltaIdx int) []db.Block {
	// fetch peer blocks delta
	var peerBlocks []db.Block
	if deltaIdx != -1 {
		peerBlocks = GetPeerBlocks(consensusNode.Address, deltaIdx)
	}
	return peerBlocks
}

// reads blockchain from file, slices the conflicting part of chain, and writes it back to the file
func clearConflictingSubchain(deltaIdx int) {
	slicedBlockchain := db.LoadBlockchain()[:deltaIdx-1]
	db.SaveBlockchain(slicedBlockchain)
}

// apply states from nodes with up-to-date chains
func tryApplyPeerStates(node Node, nodes []Node) {
	for _, peer := range nodes {
		if chainsAgree(peer.ChainHashes, node.ChainHashes) {
			node.State.TryAddTransactions(peer.State.TxMempool)
		}
	}
}

// ping peer-of-peers to potentially expand own peer set
func add2ndLevelPeers(pings PingResponseList, peersToCheck PeerSet, nodes []Node) PingResponseList {
	localIp := getLocalIP()
	for _, n := range nodes {
		for peer2 := range n.PeerSet {
			if !peersToCheck.Exists(peer2) && peer2 != localIp {
				pingRes := Ping(peer2)
				pings = append(pings, pingRes)
			}
		}
	}
	return pings
}

// get fixed sized peer set ranked by ping latency (low to high)
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
