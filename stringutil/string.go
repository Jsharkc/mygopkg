package stringutil

import (
	"strconv"
	"strings"
	"unicode"
)

// 反转字符串
func Reverse(s string) string {
	rs := []rune(s)
	if len(rs) < 2 {
		return s
	}

	for l, r := 0, len(rs)-1; l < r; l, r = l+1, r-1 {
		rs[l], rs[r] = rs[r], rs[l]
	}
	return string(rs)
}

const (
	_A        = '\u0041'
	_Z        = '\u005a'
	_a        = '\u0061'
	_z        = '\u007a'
	underline = '\u005f'
)

// 字符是否为中文
func IsChinese(ch rune) bool {
	return unicode.Is(unicode.Han, ch)
}

// 字符是否为英文字母（a-z, A-Z）
func IsEnglish(ch rune) bool {
	return (_A <= ch && ch <= _Z) || (_a <= ch && ch <= _z)
}

// 字符是否合法
//
// 合法字符为：中文，英文(a-z, A-Z)，数字(0-9)，下划线(_)
func IsCharValidate(ch rune) bool {
	if IsChinese(ch) {
		return true
	}
	if IsEnglish(ch) {
		return true
	}
	if unicode.IsDigit(ch) {
		return true
	}
	if ch == underline {
		return true
	}
	return false
}

// 名称是否合法
func IsNameValidate(name string) bool {
	if name == "" {
		return false
	}
	for _, ch := range name {
		if !IsCharValidate(ch) {
			return false
		}
	}
	return true
}

// UnicodeDecode 解码 Unicode 编码的字符串
func UnicodeDecode(text string) string {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(text), `\\`, `\`, -1))
	if err != nil {
		return text
	}
	return str
}
