package shared

// * file: Niels, s204503

// name of file to save PeerSet
const PeerSetFile = "PeerSet.json"

// {file to read from => file to write to}
var persistenceFileMappings []struct {
	from string
	to   string
} = []struct {
	from string
	to   string
}{
	{from: "Blockchain_for_testing.db", to: "Blockchain.db"},
	{from: "CurrentState_for_testing.json", to: "CurrentState.json"},
	{from: "LatestSnapshot_for_testing.json", to: "LatestSnapshot.json"},
	{from: "PeerSet_for_testing.json", to: "PeerSet.json"}}

var runtimeFiles = []string{"CurrentState.json", "LatestSnapshot.json", "state.json", "Transactions.json", "Blockchain.db", "PeerSet.json"}
