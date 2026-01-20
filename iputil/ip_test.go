package iputil_test

import (
	"testing"

	"github.com/Jsharkc/mygopkg/iputil"
)

func TestDetectLocalPrivateIP(t *testing.T) {
	localIP := "192.168.50.203"
	ip := iputil.DetectLocalPrivateIP()
	if ip != localIP {
		t.Fatalf("expect:%s, result:%s", localIP, ip)
	}
}
