package node

import (
	db "blockchain/Database"
	shared "blockchain/Shared"
	"encoding/json"
	"fmt"
	"testing"
)

// * File: Niels, s204503

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
	if shared.LegalIpAddress(shouldBeFalse3) {
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
	chainComp := []string{"0000000000000000000000000000000000000000000000000000000000000000", "2eff0509e19f074f48791f06d230f0c05bbe09e0a76d6e0dcf9703cff6fc17e0", "b46b7af31379d5b360cc2d0d8e30f2897761746a25421159b5d7698f383cc50a", "1fe834bffdbb95aa634aefca9b808244acb2378128ea5d22a34f8a623e9383fb", "415e4522a692fc4d9234285347f92ff0edec6a5d246a1462c9157cdb909a4f3c", "d1e6d8de29e702ffe8d0b8342c0b4396ac4c2875c313026cb987de48fe37409f"}
	if db.CompareChainHashes(cHashes, chainComp) != -1 {
		panic("should be equal")
	}
}

func TestLocalIP(t *testing.T) {
	localIP := getLocalIP()
	t.Log(localIP)
	if !shared.LegalIpAddress(localIP) {
		panic(fmt.Sprintf("%s should be legal!", localIP))
	}
}

func TestSortByLatency(t *testing.T) {
	pings := make([]PingResponse, 4)
	pings[1] = PingResponse{"localhost:1000", true, 1}
	pings[0] = PingResponse{"localhost:2000", true, 2}
	pings[2] = PingResponse{"localhost:3000", true, 3}
	pings[3] = PingResponse{"localhost:4000", true, 4}
	fastest := getNFastestPeers(pings, 3)
	// t.Log(fastest)
	if !(fastest.Exists(pings[1].Address) || fastest.Exists(pings[2].Address) || fastest.Exists(pings[0].Address)) {
		panic(fmt.Sprintf("PeerSet should include fastest 3 pings. [0]: %s [1]: %s", pings[0].Address, pings[1].Address))
	}
	if fastest.Exists(pings[3].Address) {
		panic(fmt.Sprintf("PeerSet should not include slowest 1 ping. [3]: %s", pings[3].Address))
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

	// t.Log(node_new)
}
