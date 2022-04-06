package node

import (
	Database "blockchain/database"
	"testing"
)

func TestRun(t *testing.T) {
	t.Log("begin init test")

	Run(Database.EmRootPath) //change to yout own path when testing

	//Database.ResetTest()
}
