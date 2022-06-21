package node

import (
	shared "blockchain/Shared"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// * file: Niels, s204503

// simple set type for storing IP addresses, with only the needed set operations; Add, Remove, Exists, Union

type PeerSet map[string]bool

func (p PeerSet) Add(k string) {
	// don't add itself
	if k == fmt.Sprintf("localhost:%d", shared.HttpPort) {
		return
	}
	if shared.LegalIpAddress(k) {
		p[k] = true
	}
}

func (p PeerSet) Remove(k string) {
	if p.Exists(k) {
		delete(p, k)
	}
}

func (p PeerSet) Exists(k string) bool {
	_, exists := p[k]
	return exists
}

func (p PeerSet) DeepCopy() PeerSet {
	new := PeerSet{}
	for k := range p {
		new[k] = true
	}
	return new
}

func (p PeerSet) UnionWith(p2 PeerSet) {
	for k := range p2 {
		p[k] = true
	}
}

func Union(s1 PeerSet, s2 PeerSet) PeerSet {
	s_union := PeerSet{}
	for k := range s1 {
		s_union[k] = true
	}
	for k := range s2 {
		s_union[k] = true
	}
	return s_union
}

// Get Peer List from JSON file
func LoadPeerSetFromJSON(filename string) PeerSet {
	// Create the file if it doesnt exist
	shared.InitDataDirIfNotExists(filename)

	data, err := os.ReadFile(shared.LocatePersistenceFile(filename, ""))
	if err != nil {
		return PeerSet{}
	}

	var ps PeerSet
	json.Unmarshal(data, &ps)

	if ps == nil {
		return PeerSet{}
	}

	return ps
}

func PersistPeerSet(ps PeerSet) {
	SavePeerSetAsJSON(ps, shared.PeerSetFile)
}

// Save the peer list in a JSON file
func SavePeerSetAsJSON(ps PeerSet, filename string) error {
	psJSON, _ := json.MarshalIndent(ps, "", "  ")

	err := ioutil.WriteFile(shared.LocatePersistenceFile(filename, ""), psJSON, 0644)
	if err != nil {
		panic(err)
	}

	return nil
}
