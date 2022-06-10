package shared

import (
	"fmt"
	"os"
	//"io/ioutil"
)

//Function that ensures that all files needed to run a node are present on the current system
//If not they are created
func CheckForNeededFiles() error {
	files := []string{"CurrentState.json", "LatestSnapshot.json", "state.json", "Transactions.json", "Blockchain.db", "PeerList.json", "PeerSet.json"}

	for _, file := range files {
		err := InitDataDirIfNotExists(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitDataDirIfNotExists(dataDir string) error {
	path := LocalDirToFileFolder + dataDir

	if fileExist(path) {
		return nil
	}

	_, err := os.Create(path)
	if err != nil {
		return err
	}

	return nil
}

func fileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// function for creating a system state that we know is a legal blockchain, for testing further functionality
func ResetPersistenceFilesForTest() {
	for _, m := range persistenceFileMappings {
		fromFile := Locate(m.from)
		toFile := Locate(m.to)
		replaceFileContents(fromFile, toFile)
	}
}

func Locate(filename string) string {
	return LocalDirToFileFolder + filename
}

func replaceFileContents(replaceWith string, fileName string) error {
	fmt.Printf("Replacing contents of %s with %s\n", fileName, replaceWith)
	data, err := os.ReadFile(replaceWith)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		panic(err)
	}

	return nil
}
