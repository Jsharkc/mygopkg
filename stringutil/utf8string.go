package stringutil

import "unicode/utf8"

// Slice 对 utf8 的 str 进行 reslice 操作，i、j 表示起始位置
func Slice(str string, i, j int) string {
	if utf8.RuneCountInString(str) > j {
		return string([]rune(str)[i:j])
	}

	if i > 0 {
		return string([]rune(str)[i:])
	}

	return str
}
