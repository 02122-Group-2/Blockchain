package node

import (
	Database "blockchain/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const httpPort = 8080

//Models the balances data recived
type balancesResult struct {
	Hash     [32]byte                         `json:"block_hash`
	Balances map[Database.AccountAddress]uint `json:"balances`
}

//Models the data for sending
type TxRequest struct {
	From   string  `json:"From"`
	To     string  `json:"To"`
	Amount float64 `json:"Amount"`
	Type   string  `json:"Type"`
}

//Models the transaction data recived
type TxResult struct {
	Status bool `json:"status"`
}

func Run(dataPath string) error {
	fmt.Println(fmt.Sprintf("Listening on port %d", httpPort))

	state := Database.LoadState() //TODO: Load from the dataDir path

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		balancesHandler(w, r, state)
	})

	http.HandleFunc("/transaction/create", func(w http.ResponseWriter, r *http.Request) {
		transactionHandler(w, r, state)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)

}

func balancesHandler(w http.ResponseWriter, r *http.Request, state *Database.State) {
	writeResult(w, balancesResult{state.LatestHash, state.AccountBalances})
}

func transactionHandler(w http.ResponseWriter, r *http.Request, state *Database.State) {
	req := TxRequest{}
	err := readReq(r, &req)

	if err != nil {
		return
	}

	var transaction Database.Transaction
	println("TYPE OF REQUEST " + req.Type)
	switch req.Type {
	case "genesis":
		transaction = state.CreateGenesisTransaction(Database.AccountAddress(req.From), float64(req.Amount))

		fmt.Println("Genesis created" + Database.TxToString(transaction))

	case "reward":
		transaction = state.CreateReward(Database.AccountAddress(req.From), float64(req.Amount))

		fmt.Println("Reward created" + Database.TxToString(transaction))

	case "transaction":
		if req.To != "" {
			transaction = state.CreateTransaction(Database.AccountAddress(req.From), Database.AccountAddress(req.To), float64(req.Amount))

			fmt.Println("Transaction created" + Database.TxToString(transaction))
		}
	}

	fmt.Println(transaction)

	err = state.AddTransaction(transaction)
	if err != nil {
		return
	}

	status := Database.SaveTransaction(state.TxMempool)

	writeResult(w, TxResult{status})
}

//When using get method
func writeResult(w http.ResponseWriter, content interface{}) {
	contentJson, err := json.Marshal(content)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(contentJson)
}

//When using post method
func readReq(r *http.Request, reqBody interface{}) error {
	reqJson, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("ERROR WAS NOT NIL READ")
		return fmt.Errorf("unable to read request body. %s", err.Error())

	}
	defer r.Body.Close()

	err = json.Unmarshal(reqJson, reqBody)
	if err != nil {
		fmt.Printf("unable to unmarshal request body. %s", err.Error())
		return fmt.Errorf("unable to unmarshal request body. %s", err.Error())
	}

	return nil

}
