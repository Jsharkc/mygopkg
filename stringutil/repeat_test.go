package stringutil_test

import (
	"testing"

	"github.com/Jsharkc/mygopkg/stringutil"
)

func TestRepeatConcat(t *testing.T) {
	type args struct {
		s     string
		count int
		sep   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{"?", 2, ","}, "?,?"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringutil.RepeatConcat(tt.args.s, tt.args.count, tt.args.sep); got != tt.want {
				t.Errorf("RepeatConcat() = %v, want %v", got, tt.want)
			}
		})
	}
}
