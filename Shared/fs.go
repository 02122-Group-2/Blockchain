package shared

import (
	"fmt"
	"os"
	//"io/ioutil"
)

//Function that ensures that all files needed to run a node are present on the current system
//If not they are created
func CheckForNeededFiles() error {

	err := InitDataDirIfNotExists("CurrentState.json")
	if err != nil {
		return err
	}

	err = InitDataDirIfNotExists("LatestSnapshot.json")
	if err != nil {
		return err
	}

	err = InitDataDirIfNotExists("state.json")
	if err != nil {
		return err
	}

	err = InitDataDirIfNotExists("Transactions.json")
	if err != nil {
		return err
	}

	err = InitDataDirIfNotExists("Blockchain.db")
	if err != nil {
		return err
	}

	err = InitDataDirIfNotExists("PeerList.json")
	if err != nil {
		return err
	}

	err = InitDataDirIfNotExists("PeerSet.json")
	if err != nil {
		return err
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
	// files to read from
	bcTestFile := Locate("Blockchain_for_testing.db")
	csTestFile := Locate("CurrentState_for_testing.json")
	lsTestFile := Locate("LatestSnapshot_for_testing.json")

	// files to write to
	bcFile := Locate("Blockchain.db")
	csFile := Locate("CurrentState.json")
	lsFile := Locate("LatestSnapshot.json")

	replaceFileContents(bcFile, bcTestFile)
	replaceFileContents(csFile, csTestFile)
	replaceFileContents(lsFile, lsTestFile)
}

func Locate(filename string) string {
	return LocalDirToFileFolder + filename
}

func replaceFileContents(fileName string, replaceWith string) error {
	fmt.Printf("Replacing contents of %s with %s", fileName, replaceWith)
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
