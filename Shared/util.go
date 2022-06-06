package shared

import "time"

func MakeTimestamp() int64 {
	return time.Now().UnixNano()
}
