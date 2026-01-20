package iputil

import (
	"net"
)

func getNetInterface() (*net.Interface, error) {
	inter, err := net.InterfaceByName("en0")
	if err != nil {
		return nil, err
	}

	return inter, nil
}
