package shared

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
	{from: "LatestSnapshot_for_testing.json", to: "LatestSnapshot.json"}}
