package iputil

import (
	"net"
	"net/netip"
)

// DetectLocalPrivateIP 获取本机内网 IP
func DetectLocalPrivateIP() string {
	loopBack := "127.0.0.1"
	inter, err := getNetInterface()
	if err != nil {
		return loopBack
	}

	addrs, err := inter.Addrs()
	if err != nil {
		return loopBack
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return loopBack
}

func IsIp(addr string) bool {
	ip, err := netip.ParseAddr(addr)
	if err != nil {
		return false
	}

	return ip.Is4() || ip.Is6()
}
