package stringutil_test

import (
	"testing"

	"github.com/Jsharkc/mygopkg/stringutil"
)

func TestSimilarity(t *testing.T) {
	type args struct {
		str1 string
		str2 string
	}
	tests := []struct {
		name  string
		input args
		want  float32
	}{
		{
			"a",
			args{
				str1: "那你工作忙不忙",
				str2: "工作是否忙碌",
			},
			0.8,
		},
		{
			"b",
			args{
				str1: "个人对自学通过的信心或判断",
				str2: "自学效果怎么样",
			},
			0.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := stringutil.Similarity(tt.input.str1, tt.input.str2); got != tt.want {
				t.Errorf("CamelCaseToUdnderscore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimilarity2(t *testing.T) {
	type args struct {
		str1 string
		str2 []string
	}
	tests := []struct {
		name  string
		input args
		want  int
	}{
		{
			"a",
			args{
				str1: "那你工作忙不忙",
				str2: []string{"工作是否忙碌", "你工作忙不忙", "自学效果怎么样", "个人对自学通过的信心或判断"},
			},
			1,
		},
		{
			"b",
			args{
				str1: "个人对自学通过的信心或判断",
				str2: []string{"自学效果怎么样", "那你工作忙不忙", "你对自学通过的信心或判断", "工作是否忙碌"},
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, f, _ := stringutil.FindSimilarityIndex(tt.input.str1, tt.input.str2...); got != tt.want {
				t.Errorf("CamelCaseToUdnderscore() = %v, want %v, f %v", got, tt.want, f)
			} else {
				t.Logf("CamelCaseToUdnderscore() = %v, want %v, f %v", got, tt.want, f)
			}
		})
	}
}
