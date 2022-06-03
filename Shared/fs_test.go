package shared

import (
	"os"
	"testing"
)

//Test currently only tests function for one file.
//It is assumed it works for all other files in the same directory
func TestCheckForNeededFiles(t *testing.T) {
	t.Log("begin CheckForNeededFiles test")

	//1. The file is already present and so nothing happens
	err := CheckForNeededFiles()
	if err != nil {
		t.Errorf("failed to do nothing...")
	}

	//Remove the file
	os.Remove(LocalDirToFileFolder + "CurrentState.json")

	//2. The file should not be present and therefore a new empty one is created
	err = CheckForNeededFiles()
	if err != nil {
		t.Errorf("failed to create missing file...")
	}
}
