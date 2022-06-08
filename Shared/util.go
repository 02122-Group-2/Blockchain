package shared

import "time"

func MakeTimestamp() int64 {
	return time.Now().UnixNano()
}

func PrettyTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
