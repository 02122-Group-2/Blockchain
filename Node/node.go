package node

import (
	Database "blockchain/Database"
	shared "blockchain/Shared"
	"fmt"
	"net/http"
)

func Run() error {
	//Assesing if all JSON files are present: CurrentState, LatestSnapshot, state,
	//Transactions, Blockchain.db
	err := shared.CheckForNeededFiles()
	if err != nil {
		return err
	}

	fmt.Printf("Listening on port %d\n", httpPort)

	go concSynchronization()
	startNode()

	return nil
}

//Sets up all the http connections and their handlers
func startNode() error {
	state := Database.LoadState()

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		balancesHandler(w, r, state)
	})
	fmt.Println("/balances/list setup complete")

	http.HandleFunc("/transaction/create", func(w http.ResponseWriter, r *http.Request) {
		transactionHandler(w, r, state)
	})
	fmt.Println("/transaction/create setup complete")

	http.HandleFunc("/getState", func(w http.ResponseWriter, r *http.Request) {
		getStateHandler(w, r, state)
	})
	fmt.Println("/getState setup complete")

	http.HandleFunc("/blockDelta", func(w http.ResponseWriter, r *http.Request) {
		blockDeltaHandler(w, r, state)
	})
	fmt.Println("/blockDelta setup complete")

	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}
