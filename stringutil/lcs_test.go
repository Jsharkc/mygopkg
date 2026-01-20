package stringutil_test

import (
	"testing"

	"github.com/Jsharkc/mygopkg/stringutil"
)

func TestLCS(t *testing.T) {
	type args struct {
		str1 string
		str2 string
	}
	tests := []struct {
		name  string
		input args
		want  int
	}{
		{
			"a",
			args{
				str1: "我",
				str2: "我的",
			},
			1,
		},
		{
			"b",
			args{
				str1: "abcd",
				str2: "ab",
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringutil.LongestCommonSubsequence(tt.input.str1, tt.input.str2); got != tt.want {
				t.Errorf("CamelCaseToUdnderscore() = %v, want %v", got, tt.want)
			}
		})
	}
}
