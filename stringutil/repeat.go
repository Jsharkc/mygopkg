package stringutil

import "strings"

// RepeatConcat 将 count 个 s 通过 sep 拼接起来
func RepeatConcat(s string, count int, sep string) string {
	s = strings.Repeat(s+sep, count)
	return strings.TrimSuffix(s, sep)
}
