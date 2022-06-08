package node

import (
	db "blockchain/Database"
	shared "blockchain/Shared"
	"encoding/json"
	"fmt"
	"testing"
)

func TestIpRegex(t *testing.T) {
	shouldBeTrue := "192.168.0.1:8080"
	shouldBeTrue2 := "localhost:8080"
	shouldBeFalse := "192.168.0.1:808022"
	shouldBeFalse2 := "asdf"
	shouldBeFalse3 := "256.168.0.1:8080"

	if !shared.LegalIpAddress(shouldBeTrue) {
		panic(fmt.Sprintf("%s should be true", shouldBeTrue))
	}

	if !shared.LegalIpAddress(shouldBeTrue2) {
		panic(fmt.Sprintf("%s should be true", shouldBeTrue2))
	}

	if shared.LegalIpAddress(shouldBeFalse) {
		panic(fmt.Sprintf("%s should be false", shouldBeFalse2))
	}

	if shared.LegalIpAddress(shouldBeFalse2) {
		panic(fmt.Sprintf("%s should be false", shouldBeFalse2))
	}
	if legalIpAddress(shouldBeFalse3) {
		panic(fmt.Sprintf("%s should be false", shouldBeFalse3))
	}
}

func TestPeerSet(t *testing.T) {
	ps := PeerSet{}
	legalIP := "192.168.0.1:8080"
	ps.Add(legalIP)
	if !ps.Exists(legalIP) {
		panic(fmt.Sprintf("%s should be added, as it is a legal IP", legalIP))
	}
	illegalIP := "192.168.0.1:808022"
	ps.Add(illegalIP)
	if ps.Exists(illegalIP) {
		panic(fmt.Sprintf("%s should not be added, as it is an illegal IP", illegalIP))
	}

	ps2 := PeerSet{}
	legalIP2 := "localhost:8080"
	ps2.Add(legalIP)
	ps2.Add(legalIP2)

	ps_copy := ps.DeepCopy()

	ps.UnionWith(ps2)
	if !ps.Exists(legalIP2) {
		panic(fmt.Sprintf("ps should contain %s after union operation with ps2", legalIP2))
	}

	ps_union := Union(ps_copy, ps2)
	if !(ps_union.Exists(legalIP) && ps_union.Exists(legalIP2)) {
		panic("union of ps's should contain all elements of both ps's")
	}
}

func TestConstructSubsets(t *testing.T) {
	ps := PeerSet{}
	ps.Add("192.168.0.1:8080")
	ps.Add("192.168.0.2:8080")
	ps.Add("192.168.0.3:8080")
	ps.Add("192.168.0.4:8080")
	ps.Add("192.168.0.5:8080")

	subsets := constructSubsets(ps)
	t.Log(subsets)
	if len(subsets) > MAX_SUBSETS {
		panic(fmt.Sprintf("# of subsets (%d) must not be greater than %d", len(subsets), MAX_SUBSETS))
	}
}

func TestGetLocalChainHashes(t *testing.T) {
	shared.ResetPersistenceFilesForTest()
	state := db.LoadState()
	cHashes := db.GetLocalChainHashes(*state, 0)
	t.Log(cHashes)
	chainComp := []string{"0000000000000000000000000000000000000000000000000000000000000000", "c352edf51ac6fdf40de39d11a85c1f1a90620028905a6ead5fa78da04eee75cc", "0116adebf51528def8fdb441daa7620c017d1fe288fdb8071e24717aea74f81c", "811a21a6ad322ab9e5f68cbcb47bf20a094ba55612a404f00a83ccb93e57c063"}
	if db.CompareChainHashes(cHashes, chainComp) != -1 {
		panic("should be equal")
	}
}

func TestLocalIP(t *testing.T) {
	localIP := getLocalIP()
	t.Log(localIP)
	if !legalIpAddress(localIP) {
		panic(fmt.Sprintf("%s should be legal!", localIP))
	}
}

func TestSortByLatency(t *testing.T) {
	pings := make([]PingResponse, 4)
	pings[1] = PingResponse{"localhost:8081", true, 1}
	pings[0] = PingResponse{"localhost:8082", true, 2}
	pings[2] = PingResponse{"localhost:8083", true, 3}
	pings[3] = PingResponse{"localhost:8084", true, 4}
	fastest := getNFastestPeers(pings, 3)
	t.Log(fastest)
	if !(fastest.Exists(pings[1].Address) || fastest.Exists(pings[2].Address) || fastest.Exists(pings[0].Address)) {
		panic(fmt.Sprintf("PeerSet should include fastest 3 pings. [0]: %s [1]: %s", pings[0].Address, pings[1].Address))
	}
	if fastest.Exists(pings[3].Address) {
		panic(fmt.Sprintf("PeerSet should not include slowest 1 ping. [3]: %s", pings[3].Address))
	}
}

func TestChainDiffIdx(t *testing.T) {
	c1 := []string{"a", "b", "c", "d", "e"}
	c2 := []string{"a", "b", "d", "e", "f"}

	idx := chainDiffIdx(c1, c2)
	t.Log(idx)
	if idx != 2 {
		panic("Found index is wrong!")
	}

	c1 = []string{"a", "b"}
	c2 = []string{"a", "b", "d", "e", "f"}
	idx = chainDiffIdx(c1, c2)
	t.Log(idx)
	if idx != 2 {
		panic("Found index is wrong!")
	}
}

func testConsensusUtil(chains [][]string) []Node {
	s := make([]db.State, len(chains))
	nodes := make([]Node, len(chains))
	for i := 0; i < len(chains); i++ {
		s[i] = db.State{}
		s[i].LastBlockSerialNo = len(chains[i]) - 1
		nodes[i] = Node{fmt.Sprintf("localhost:808%d", i+1), nil, s[i], chains[i]}
	}
	return nodes
}

// test fork with equal no. of agreeing nodes
func TestComputeConsensusNode(t *testing.T) {
	c1 := []string{"a", "b", "c", "d", "e"}
	c2 := []string{"a", "b", "d", "e", "f"}
	c3 := []string{"a", "b"}
	c4 := []string{"a", "b"}
	c5 := []string{"a", "b", "c", "d", "e"}

	chains := [][]string{c1, c2, c3, c4, c5}

	nodes := testConsensusUtil(chains)

	cons := computeConsensusNode(nodes)

	t.Log(cons.ChainHashes)
	if chainDiffIdx(c1, cons.ChainHashes) != -1 {
		panic("Consensus algo aint work")
	}
	fmt.Println("succeeded")
}

// test with total separate (longer) chain than consensus chain
func TestComputeConsensusNode2(t *testing.T) {
	c1 := []string{"a", "b", "c", "d", "e"}
	c2 := []string{"a", "b", "d", "e", "f"}
	c3 := []string{"a", "b", "d", "e", "f"}
	c4 := []string{"q", "x", "y", "z", "w", "æ", "ø", "å"}

	chains := [][]string{c1, c2, c3, c4}

	nodes := testConsensusUtil(chains)

	cons := computeConsensusNode(nodes)

	t.Log(cons.ChainHashes)
	if chainDiffIdx(c2, cons.ChainHashes) != -1 {
		panic("Consensus algo aint work")
	}
}

func TestMarshalUnmarshalNode(t *testing.T) {
	node := GetNode()
	node_json, err := json.Marshal(node)
	var node_new NodeFromPostRequest
	err2 := json.Unmarshal(node_json, &node_new)

	if err != nil {
		panic(fmt.Sprintf("Error marshaling Node, %v", err.Error()))
	}

	if err2 != nil {
		panic(fmt.Sprintf("Error unmarshaling Node, %v", err2.Error()))
	}

	t.Log(node_new)
}
