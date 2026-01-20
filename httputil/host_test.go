package httputil_test

import (
	"testing"

	"github.com/Jsharkc/mygopkg/httputil"
)

func TestGetTopDomain(t *testing.T) {
	tests := []struct {
		name string
		host string
		want string
	}{
		{"a", "192.168.1.187:2023", "192.168.1.187:2023"},
		{"b", "192.168.1.187", "192.168.1.187"},
		{"c", "megaview.com:2023", "megaview.com:2023"},
		{"d", "megaview.com", "megaview.com"},
		{"e", "app.megaview.com", "megaview.com"},
		{"f", "app.test.megaview.com", "megaview.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := httputil.GetTopDomain(tt.host); got != tt.want {
				t.Errorf("GetTopDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
