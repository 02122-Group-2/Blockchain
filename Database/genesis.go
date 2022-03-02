package database

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Genesis struct {
	Balances map[AccountAddress]int `json:"balances"`
}

func LoadGenesis() *Genesis {
	currWD, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(filepath.Join(currWD, "Genesis.json"))
	if err != nil {
		panic(err)
	}

	var loadedGenesis Genesis
	json.Unmarshal(data, &loadedGenesis)

	return &loadedGenesis
}
