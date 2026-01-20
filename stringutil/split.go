package stringutil

// Split 根据指定的分隔符将字符串 s 分割成多个子字符串,如果分隔符包括多个字符，则按单个字符拆分字符串，并返回分割后的结果。
// 参数 s 为待分割的字符串。
// 参数 sep 为分隔符，可以为多个字符组成的字符串。
// 参数 includeSeparator 为可选参数，表示是否保留分隔符。
// 返回值为分割后的字符串切片。
func Split(s string, sep string, includeSeparator ...bool) []string {
	splits := make([]string, 0, 8)
	if len(sep) == 0 || len(s) == 0 {
		return []string{s}
	}

	includeSeps := len(includeSeparator) > 0 && includeSeparator[0]

	//去重
	sepSet := make(map[rune]struct{})
	for _, char := range sep {
		sepSet[char] = struct{}{}
	}

	str := ""
	for _, r := range s {
		if _, ok := sepSet[r]; ok {
			if includeSeps {
				str += string(r)
			}
			splits = append(splits, str)
			str = ""
		} else {
			str += string(r)
		}
	}
	if len(str) > 0 {
		splits = append(splits, str)
	}
	return splits
}
