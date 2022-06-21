package shared

import (
	"fmt"
	"os"
	"testing"
)

// * Emilie, s204471
//Tests if CheckingForNeededFiles work for all the needed files
func TestEnsureNeededFilesExist(t *testing.T) {
	t.Log("begin EnsureNeededFilesExist test")

	//1. The file is already present and so nothing happens
	err := EnsureNeededFilesExist()
	if err != nil {
		t.Errorf("failed to do nothing...")
	}

	//Remove the files
	for _, fileName := range runtimeFiles {
		os.Remove(LocatePersistenceFile(fileName, ""))
	}

	//2. The file should not be present and therefore a new empty one is created
	err = EnsureNeededFilesExist()
	if err != nil {
		t.Errorf("failed to create missing files...")
	}

	for _, fileName := range runtimeFiles {
		if !fileExist(LocatePersistenceFile(fileName, "")) {
			panic(fmt.Sprintf("Error %s was not created\n", fileName))
		}
	}
}

// * Niels, s204503
// test result is cached and as a result the reset is not run more than once when running the test
func TestResetPersistenceFiles(t *testing.T) {
	ResetPersistenceFilesForTest()

	for _, fileMapping := range persistenceFileMappings {
		replacedFile, checkFile := LocatePersistenceFile(fileMapping.from, "test_data"), LocatePersistenceFile(fileMapping.to, "")
		replacedFileSum, checkFileSum := GetChecksum(replacedFile), GetChecksum(checkFile)
		if replacedFileSum != checkFileSum {
			panic(fmt.Sprintf("Checksums do not match for files %s and %s\n%x\n%x\n", replacedFile, checkFile, replacedFileSum, checkFileSum))
		}
		fmt.Printf("Checksums for files %s and %s:\n%x\n%x\n\n", fileMapping.from, fileMapping.to, replacedFileSum, checkFileSum)
	}
}
