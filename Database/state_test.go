package database

import (
	shared "blockchain/Shared"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
)

// * file: Niels, s204503

func TestLoadState(t *testing.T) {
	shared.ResetPersistenceFilesForTest()
	t.Log("Start load state test")
	hash, _ := hex.DecodeString("811a21a6ad322ab9e5f68cbcb47bf20a094ba55612a404f00a83ccb93e57c063")
	hash32 := [32]byte{}
	for i := 0; i < 32; i++ {
		hash32[i] = hash[i]
	}
	state_comp := State{LatestHash: hash32, AccountBalances: map[AccountAddress]uint{"Alberto": 15100, "Asger": 423, "Emilie": 22, "Magn": 421011, "Niels": 579039, "gggg": 3}, AccountNounces: map[AccountAddress]uint{"Emilie": 3, "Magn": 10, "Niels": 9, "system": 10}, TxMempool: nil, DbFile: nil, LastBlockSerialNo: 3, LastBlockTimestamp: 1648641894935865700, LatestTimestamp: 1648641894935865700}
	state := loadStateFromJSON("test_data/state_test.json")

	state_comp_json, _ := json.Marshal(state_comp)
	state_comp_json_string := fmt.Sprintf("%v", state_comp_json)
	state_json, _ := json.Marshal(state)
	state_json_string := fmt.Sprintf("%v", state_json)

	if state_json_string != state_comp_json_string {
		panic(fmt.Sprintf("%s\n%s\n", state_comp_json, state_comp_json_string))
	}
}

func TestRecomputeState(t *testing.T) {
	s := LoadState()
	s.RecomputeState(4)
	t.Logf("%x\n", s.LatestHash)
}
