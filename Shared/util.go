package shared

import (
	"fmt"
	"regexp"
	"time"
)

func MakeTimestamp() int64 {
	return time.Now().UnixNano()
}

func PrettyTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func LegalIpAddress(addr string) bool {
	regexIPv4withPort := "^(localhost|((([0-1]{0,1}[0-9]{1,2})|2([0-4][0-9]|5[0-5])).){3}(([0-1]{0,1}[0-9]{1,2})|2([0-4][0-9]|5[0-5]))):[0-9]{4,5}$"
	match, _ := regexp.MatchString(regexIPv4withPort, addr)
	return match
}

func Log(msg string) {
	fmt.Printf("%s: %s\n", PrettyTimestamp(), msg)
}
