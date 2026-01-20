package stringutil_test

import (
	"strings"
	"testing"

	"github.com/Jsharkc/mygopkg/stringutil"
)

func TestSplit(t *testing.T) {
	type args struct {
		s                string
		sep              string
		includeSeparator bool
	}
	tests := []struct {
		name  string
		input args
		want  []string
	}{
		{
			"a",
			args{
				s:                "这是一句话",
				sep:              "",
				includeSeparator: false,
			},
			[]string{"这是一句话"},
		},
		{
			"b",
			args{
				s:                "",
				sep:              "",
				includeSeparator: false,
			},
			[]string{""},
		},
		{
			"b",
			args{
				s:                "",
				sep:              ".。",
				includeSeparator: false,
			},
			[]string{""},
		},
		{
			"b",
			args{
				s:                "这句话以中文逗号结尾，这句话以中文句号结尾。这句话以中文分号结尾；这句话以英文逗号结尾,这句话以英文句号结尾.这句话以英文分号结尾;",
				sep:              ".。",
				includeSeparator: false,
			},
			[]string{"这句话以中文逗号结尾，这句话以中文句号结尾", "这句话以中文分号结尾；这句话以英文逗号结尾,这句话以英文句号结尾", "这句话以英文分号结尾;"},
		},
		{
			"b",
			args{
				s:                "这句话以中文逗号结尾，这句话以中文句号结尾。这句话以中文分号结尾；这句话以英文逗号结尾,这句话以英文句号结尾.这句话以英文分号结尾;",
				sep:              ".。",
				includeSeparator: true,
			},
			[]string{"这句话以中文逗号结尾，这句话以中文句号结尾。", "这句话以中文分号结尾；这句话以英文逗号结尾,这句话以英文句号结尾.", "这句话以英文分号结尾;"},
		},
		{
			"b",
			args{
				s:                "这句话以中文逗号结尾，这句话以中文句号结尾。这句话以中文分号结尾；这句话以英文逗号结尾,这句话以英文句号结尾.这句话以英文分号结尾;",
				sep:              ".。，,；;",
				includeSeparator: false,
			},
			[]string{"这句话以中文逗号结尾", "这句话以中文句号结尾", "这句话以中文分号结尾", "这句话以英文逗号结尾", "这句话以英文句号结尾", "这句话以英文分号结尾"},
		},
		{
			"b",
			args{
				s:                "这句话以中文逗号结尾，这句话以中文句号结尾。这句话以中文分号结尾；这句话以英文逗号结尾,这句话以英文句号结尾.这句话以英文分号结尾;",
				sep:              ".。",
				includeSeparator: true,
			},
			[]string{"这句话以中文逗号结尾，", "这句话以中文句号结尾。", "这句话以中文分号结尾；", "这句话以英文逗号结尾,", "这句话以英文句号结尾.", "这句话以英文分号结尾;"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringutil.Split(tt.input.s, tt.input.sep, tt.input.includeSeparator); strings.Join(got, "") != strings.Join(tt.want, "") {
				t.Errorf("stringutil.Split() = %v, want %v", got, tt.want)
			}
		})
	}
}
