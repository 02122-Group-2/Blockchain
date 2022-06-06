package shared

import (
	"fmt"
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
	os.Remove(LocalDirToFileFolder + "PeerSet.json")
	os.Remove(LocalDirToFileFolder + "PeerList.json")

	//2. The file should not be present and therefore a new empty one is created
	err = CheckForNeededFiles()
	if err != nil {
		t.Errorf("failed to create missing files...")
	}

	//Check if the files are present
	if !fileExist(Locate("CurrentState.json")) {
		t.Log("Error CurrentState.json was not created")
		t.Fail()
	} else if !fileExist(Locate("Blockchain.db")) {
		t.Log("Error Blockchain.db was not created")
		t.Fail()
	} else if !fileExist(Locate("state.json")) {
		t.Log("Error state.json was not created")
		t.Fail()
	} else if !fileExist(Locate("LatestSnapshot.json")) {
		t.Log("Error LatestSnapshot.json was not created")
		t.Fail()
	} else if !fileExist(Locate("Transactions.json")) {
		t.Log("Error Transactions.json was not created")
		t.Fail()
	} else if !fileExist(Locate("PeerLis.json")) {
		t.Log("Error PeerList.json was not created")
		t.Fail()
	} else if !fileExist(Locate("PeerSet.json")) {
		t.Log("Error PeerSet.json was not created")
		t.Fail()
	}
}

// test result is cached and as a result the reset is not run more than once when running the test
func TestResetPersistenceFiles(t *testing.T) {
	fmt.Println("bruh")
	ResetPersistenceFilesForTest()
}
