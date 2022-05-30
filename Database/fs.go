package database

import (
	"os"
	//"io/ioutil"
)

//Function that ensures that all files needed to run a node are present on the current system
//If not they are created
func CheckForNeededFiles() error {

	err := initDataDirIfNotExists("CurrentState.json")
	if err != nil {
		return err
	}

	err = initDataDirIfNotExists("LatestSnapshot.json")
	if err != nil {
		return err
	}

	err = initDataDirIfNotExists("state.json")
	if err != nil {
		return err
	}

	err = initDataDirIfNotExists("Transactions.json")
	if err != nil {
		return err
	}

	err = initDataDirIfNotExists("Blockchain.db")
	if err != nil {
		return err
	}

	return nil
}

func initDataDirIfNotExists(dataDir string) error {
	path := localDirToFileFolder + dataDir

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
