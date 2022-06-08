package shared

import (
	"regexp"
)

func LegalIpAddress(addr string) bool {
	regexIPwithPort := "^(localhost|([0-9]{1,3}.){1,3}([0-9]{1,3})):([0-9]{4,5})$"
	match, _ := regexp.MatchString(regexIPwithPort, addr)
	return match
}