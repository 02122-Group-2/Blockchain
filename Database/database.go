package database

import "time"

type DatabaseInfo struct {
	CurrentSerialNo int
}

func makeTimestamp() int64 {
	return time.Now().UnixNano()
}

func (dbInfo *DatabaseInfo) getNextSerialNo() int {
	curNo := dbInfo.CurrentSerialNo
	dbInfo.CurrentSerialNo = dbInfo.CurrentSerialNo + 1
	return curNo
}
