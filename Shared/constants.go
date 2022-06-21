package shared

import (
	"os"
	"path/filepath"
)


var LocalDirToWallets string = getProgLocation() + "\\Wallets\\"
var LocalDirToFileFolder string = getProgLocation() + "\\Persistence\\"
var HttpPort = 8080
var BootstrapNode = "localhost:8080"

func getProgLocation() string {
	ex, err := os.Executable()
	if err != nil {
			panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}