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

	//Remove the files
	for _, fileName := range runtimeFiles {
		os.Remove(Locate(fileName))
	}

	//2. The file should not be present and therefore a new empty one is created
	err = CheckForNeededFiles()
	if err != nil {
		t.Errorf("failed to create missing files...")
	}

	for _, fileName := range runtimeFiles {
		if !fileExist(Locate(fileName)) {
			panic(fmt.Sprintf("Error %s was not created\n", fileName))
		}
	}
}

// test result is cached and as a result the reset is not run more than once when running the test
func TestResetPersistenceFiles(t *testing.T) {
	ResetPersistenceFilesForTest()

	for _, fileMapping := range persistenceFileMappings {
		replacedFile, checkFile := Locate(fileMapping.from), Locate(fileMapping.to)
		replacedFileSum, checkFileSum := GetChecksum(replacedFile), GetChecksum(checkFile)
		if replacedFileSum != checkFileSum {
			panic(fmt.Sprintf("Checksums do not match for files %s and %s\n%x\n%x\n", replacedFile, checkFile, replacedFileSum, checkFileSum))
		}
		fmt.Printf("Checksums for files %s and %s:\n%x\n%x\n\n", fileMapping.from, fileMapping.to, replacedFileSum, checkFileSum)
	}
}
