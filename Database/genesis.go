package database

type Genesis struct {
	Balances map[AccountAddress]int `json:"balances"`
}
