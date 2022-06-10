package shared

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
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
	} else if !fileExist(Locate("PeerSet.json")) {
		t.Log("Error PeerSet.json was not created")
		t.Fail()
	}
}

func getChecksum(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

// test result is cached and as a result the reset is not run more than once when running the test
func TestResetPersistenceFiles(t *testing.T) {
	// fmt.Println("bruh4202")
	ResetPersistenceFilesForTest()

	for _, fileMapping := range persistenceFileMappings {
		replacedFile, checkFile := Locate(fileMapping.from), Locate(fileMapping.to)
		replacedFileSum, checkFileSum := getChecksum(replacedFile), getChecksum(checkFile)
		if replacedFileSum != checkFileSum {
			panic(fmt.Sprintf("Checksums do not match for files %s and %s\n%x\n%x\n", replacedFile, checkFile, replacedFileSum, checkFileSum))
		}
		fmt.Printf("Checksums for files %s and %s:\n%x\n%x\n\n", fileMapping.from, fileMapping.to, replacedFileSum, checkFileSum)
	}
}
