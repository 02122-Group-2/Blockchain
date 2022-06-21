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

// * Emilie, s204471
func Run() error {
	err := shared.EnsureNeededFilesExist()
	if err != nil {
		return err
	}

	shared.Log(fmt.Sprintf("Listening on port %d", shared.HttpPort))

	go concSynchronization()
	startNode()

	return nil
}

// * Emilie, s204471
//Sets up all the http connections and their handlers
func startNode() error {
	state := Database.LoadState()

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		balancesHandler(w, r, state)
	})
	shared.Log("/balances/list setup complete")

	http.HandleFunc("/getState", func(w http.ResponseWriter, r *http.Request) {
		getStateHandler(w, r, state)
	})
	shared.Log("/getState setup complete")

	http.HandleFunc("/blockDelta", func(w http.ResponseWriter, r *http.Request) {
		blockDeltaHandler(w, r, state)
	})
	shared.Log("/blockDelta setup complete")

	return http.ListenAndServe(fmt.Sprintf(":%d", shared.HttpPort), nil)
}

// * Niels, s204503
// Get the initial node state
func GetNode() Node {
	node := Node{}
	node.Address = getLocalIP()
	node.State = *Database.LoadState()
	node.PeerSet = GetPeerSet()
	node.ChainHashes = Database.GetLocalChainHashes(node.State, 0)
	return node
}

// * Niels, s204503
// Get the stored set of nodes
// If this hasn't been created before, create it using the bootstrap node
func GetPeerSet() PeerSet {
	ps := LoadPeerSetFromJSON(shared.PeerSetFile)
	if ps == nil {
		ps = PeerSet{shared.BootstrapNode: true}
	}
	if len(ps) == 0 {
		ps = PeerSet{shared.BootstrapNode: true}
	}
	ps.Add(shared.BootstrapNode)
	return ps
}

// * Niels, s204503
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

// * Niels, s204503
func getLocalIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return fmt.Sprintf("%v:%d", ipv4, shared.HttpPort)
		}
	}
	return fmt.Sprintf("localhost:%d", shared.HttpPort)
}
