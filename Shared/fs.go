package shared

import (
	"fmt"
	"os"
	//"io/ioutil"
)

//Function that ensures that all files needed to run a node are present on the current system
//If not they are created
func EnsureNeededFilesExist() error {
	for _, file := range runtimeFiles {
		err := InitDataDirIfNotExists(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitDataDirIfNotExists(dataDir string) error {
	path := LocatePersistenceFile(dataDir, "")

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
		fromFile := LocatePersistenceFile(m.from, "test_data")
		toFile := LocatePersistenceFile(m.to, "")
		replaceFileContents(fromFile, toFile)
	}
}

func LocatePersistenceFile(filename string, subfolder string) string {
	if subfolder != "" {
		subfolder += "/"
	}
	return fmt.Sprintf("%s/%s%s", LocalDirToFileFolder, subfolder, filename)
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
