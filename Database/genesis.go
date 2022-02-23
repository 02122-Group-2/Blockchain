package database

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	fmt.Println(loadGenesis())
}

type Genesis struct {
	Balances map[AccountAddress]int `json:"balances"`
}

func loadGenesis() *Genesis {
	data, err := os.ReadFile("./Genesis.json")
	if err != nil {
		panic(err)
	}

	var loadedGenesis *Genesis
	json.Unmarshal([]byte(data), loadedGenesis)

	return loadedGenesis
}
