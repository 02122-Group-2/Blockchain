package node

import (
	"fmt"
	"testing"
)

func TestIpRegex(t *testing.T) {
	shouldBeTrue := "192.168.0.1:8080"
	shouldBeFalse := "192.168.0.1:808022"
	shouldBeFalse2 := "asdf"

	if !legalIpAddress(shouldBeTrue) {
		panic(fmt.Sprintf("%s should be true", shouldBeTrue))
	}

	if legalIpAddress(shouldBeFalse) {
		panic(fmt.Sprintf("%s should be false", shouldBeFalse2))
	}

	if legalIpAddress(shouldBeFalse2) {
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
}
