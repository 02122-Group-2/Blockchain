package node

import (
	shared "blockchain/Shared"
	"fmt"
	"testing"
)

func TestIpRegex(t *testing.T) {
	shouldBeTrue := "192.168.0.1:8080"
	shouldBeFalse := "192.168.0.1:808022"
	shouldBeFalse2 := "asdf"

	if !shared.LegalIpAddress(shouldBeTrue) {
		panic(fmt.Sprintf("%s should be true", shouldBeTrue))
	}

	if shared.LegalIpAddress(shouldBeFalse) {
		panic(fmt.Sprintf("%s should be false", shouldBeFalse2))
	}

	if shared.LegalIpAddress(shouldBeFalse2) {
		panic(fmt.Sprintf("%s should be false", shouldBeFalse2))
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

