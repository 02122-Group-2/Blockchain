package shared

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"
)

// * file: Niels, s204503

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

func GetChecksum(filepath string) string {
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
