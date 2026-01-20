package stringutil_test

import (
	"testing"

	"github.com/Jsharkc/mygopkg/stringutil"
)

func TestCamelCaseToUdnderscore(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"a", "AtomicPerm", "atomic_perm"},
		{"b", "OrgID", "org_id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringutil.CamelCaseToUdnderscore(tt.input); got != tt.want {
				t.Errorf("CamelCaseToUdnderscore() = %v, want %v", got, tt.want)
			}
		})
	}
}
