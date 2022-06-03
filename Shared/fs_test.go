package shared

import (
	"os"
	"testing"
)

//Tests if CheckingForNeededFiles work for all the needed files
func TestCheckForNeededFiles(t *testing.T) {
	t.Log("begin CheckForNeededFiles test")

	//1. The file is already present and so nothing happens
	err := CheckForNeededFiles()
	if err != nil {
		t.Errorf("failed to do nothing...")
	}

	//Remove the file
	os.Remove(LocalDirToFileFolder + "CurrentState.json")
	os.Remove(LocalDirToFileFolder + "Blockchain.db")
	os.Remove(LocalDirToFileFolder + "state.json")
	os.Remove(LocalDirToFileFolder + "LatestSnapshot.json")
	os.Remove(LocalDirToFileFolder + "Transactions.json")

	//2. The file should not be present and therefore a new empty one is created
	err = CheckForNeededFiles()
	if err != nil {
		t.Errorf("failed to create missing files...")
	}

	//Check if the files are present
	if !fileExist(LocalDirToFileFolder + "CurrentState.json") {
		t.Log("Error CurrentState.json was not created")
		t.Fail()
	} else if !fileExist(LocalDirToFileFolder + "Blockchain.db") {
		t.Log("Error Blockchain.db was not created")
		t.Fail()
	} else if !fileExist(LocalDirToFileFolder + "state.json") {
		t.Log("Error state.json was not created")
		t.Fail()
	} else if !fileExist(LocalDirToFileFolder + "LatestSnapshot.json") {
		t.Log("Error LatestSnapshot.json was not created")
		t.Fail()
	} else if !fileExist(LocalDirToFileFolder + "Transactions.json") {
		t.Log("Error Transactions.json was not created")
		t.Fail()
	}
}
