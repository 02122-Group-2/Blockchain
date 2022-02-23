package database

import (
	"encoding/json"
	"os"
)

type Genesis struct {
	Balances map[AccountAddress]int `json:"balances"`
}

func LoadGenesis() *Genesis {
	data, err := os.ReadFile("./Genesis.json")
	if err != nil {
		panic(err)
	}

	var loadedGenesis Genesis
	json.Unmarshal(data, &loadedGenesis)

	return &loadedGenesis
}
