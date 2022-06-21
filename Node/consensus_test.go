package node

import (
	db "blockchain/Database"
	"fmt"
	"testing"
)

// * file: Niels, s204503

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
