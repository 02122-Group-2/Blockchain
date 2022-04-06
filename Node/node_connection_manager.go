package node

// Service interface for performing networking calls

type Status int

// enum
const (
	READY Status = iota
	BUSY
	FAILED
	NO_CONN
)

type NodeResponse struct {
	Status   Status  `json:"Status"`
	Response string  `json:"Response"`
	Data     *[]byte `json:"Data"` // nullable, for e.g. HeartBeat checks
}

type NodeConnectionManager interface {
	GetHeartBeat(fqdn string) NodeResponse
	FetchStateData(fqdn string) (int, []byte)
	// SendStateData(fqdn string)
	// Fetch2ndLevelPeerList(fqdn string)
}
