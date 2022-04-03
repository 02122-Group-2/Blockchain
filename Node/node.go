package node

import (
	Database "blockchain/database"
	"encoding/json"
	"fmt"
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
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	Timestamp int64   `json:"timestamp"`
	Type      string  `json:"type"`
}

//Models the transaction data recived
type TxResult struct {
	Hash [32]byte `json:"block_hash"`
}

func Run(dataDir string) error {
	fmt.Println(fmt.Sprintf("Listening on port %d", httpPort))

	state, err := Database.LoadState() //TODO: Load from the dataDir path
	if err != nil {
		return err
	}

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		balancesHandler(w, r, state)
	})

	/*http.HandleFunc("/transaction/create", func(w http.ResponseWriter, r *http.Request) {
		transactionHandler(w,r,state)
	})
	*/
	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)

}

func balancesHandler(w http.ResponseWriter, r *http.Request, state *Database.State) {
	writeResult(w, balancesResult{state.LatestHash, state.Balances})
}

/*
func transactionHandler(w http.ResponseWriter, r http.Request, state *database.State) {
	req := TxRequest{}
	err := readRequest(r, &req)
	if err != null {
		writeErrResult(w, err)
		return
	}

	tx := database.CreateTransaction(database.newAccountAddr(req.From), database.newAccountAddr(req.To), req.Amount, req.Type)

	err := database.AddTransaction(tx)
	if err != nil {
		writeErrResult(w, err)
		return
	}

	hash, err := state.Persist()
	if err != nil {
		WriteErrRes(w, err)
		return
	}

	writeResult(w, TxResult{hash})
}
*/

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

func writeErrResult(w http.ResponseWriter)
