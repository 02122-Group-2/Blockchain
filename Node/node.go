package node

import (
	Database "blockchain/Database"
	shared "blockchain/Shared"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

func Run() error {
	err := shared.EnsureNeededFilesExist()
	if err != nil {
		return err
	}

	shared.Log(fmt.Sprintf("Listening on port %d", httpPort))

	// go concSynchronization()
	startNode()

	return nil
}

//Sets up all the http connections and their handlers
func startNode() error {
	state := Database.LoadState()

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		balancesHandler(w, r, state)
	})
	shared.Log("/balances/list setup complete")

	http.HandleFunc("/transaction/create", func(w http.ResponseWriter, r *http.Request) {
		transactionHandler(w, r, state)
	})
	shared.Log("/transaction/create setup complete")

	http.HandleFunc("/getState", func(w http.ResponseWriter, r *http.Request) {
		getStateHandler(w, r, state)
	})
	shared.Log("/getState setup complete")

	http.HandleFunc("/blockDelta", func(w http.ResponseWriter, r *http.Request) {
		blockDeltaHandler(w, r, state)
	})
	shared.Log("/blockDelta setup complete")

	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

// Get the initial node state
func GetNode() Node {
	node := Node{}
	node.Address = getLocalIP()
	node.State = *Database.LoadState()
	node.PeerSet = GetPeerSet()
	node.ChainHashes = Database.GetLocalChainHashes(node.State, 0)
	return node
}

// Get the stored set of nodes
// If this hasn't been created before, create it using the bootstrap node
func GetPeerSet() PeerSet {
	ps := LoadPeerSetFromJSON(shared.PeerSetFile)
	if ps == nil {
		ps = PeerSet{bootstrapNode: true}
	}
	if len(ps) == 0 {
		ps.Add(bootstrapNode)
	}
	return ps
}

func Ping(peerAddr string) PingResponse {
	if !shared.LegalIpAddress(peerAddr) {
		return PingResponse{"nil", false, -1}
	}

	startTime := shared.MakeTimestamp()
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", peerAddr, timeout)
	if err != nil {
		fmt.Println("Site unreachable, error: ", err)
		return PingResponse{"nil", false, -1}
	}
	endTime := shared.MakeTimestamp()
	latency := endTime - startTime
	conn.Close()
	return PingResponse{peerAddr, true, latency}
}

func getLocalIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return fmt.Sprintf("%v:%d", ipv4, httpPort)
		}
	}
	return fmt.Sprintf("localhost:%d", httpPort)
}
