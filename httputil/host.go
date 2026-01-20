package httputil

import (
	"strings"

	"github.com/Jsharkc/mygopkg/iputil"
)

func GetTopDomain(host string) string {
	addrs := strings.Split(host, ":")
	if len(addrs) > 0 {
		if iputil.IsIp(addrs[0]) {
			return host
		}
	}

	domains := strings.Split(host, ".")
	l := len(domains)
	if l <= 2 {
		return host
	}

	return domains[l-2] + "." + domains[l-1]
}
