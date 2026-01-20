package stringutil

import (
	"regexp"
)

const (
	PatternPhone     = `^(0|\+?86)?[1-9]\d{10}$`
	PatternPhonePart = `(0|\+?86)?[1-9]\d{10}`
	PatternEmail     = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	PatternAccount    = `^[a-zA-Z0-9_]{4,20}$`
)

var (
	regexpPhone     = regexp.MustCompile(PatternPhone)
	regexpPhonePart = regexp.MustCompile(PatternPhonePart)
	regexpEmail     = regexp.MustCompile(PatternEmail)
	regexpAccount   = regexp.MustCompile(PatternAccount)
)

// 获取参数的字符串中的所有手机号, max为最大的匹配结果数量, -1为所有
func FindAllPhone(input string, max int) []string {
	return regexpPhonePart.FindAllString(input, max)
}

// 校验手机号是否有效
func CheckPhone(phone string) bool {
	found := regexpPhone.MatchString(phone)
	return found
}

// 校验邮箱是否有效
func CheckEmail(email string) bool {
	found := regexpEmail.MatchString(email)
	return found
}

// 校验用户账号
func CheckUserAccount(account string) bool {
	found := regexpAccount.MatchString(account)
	return found
}
