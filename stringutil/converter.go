package stringutil

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	commonInitialisms         = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	commonInitialismsReplacer *strings.Replacer
)

func init() {
	caser := cases.Title(language.Und)

	commonInitialismsForReplacer := make([]string, 0, len(commonInitialisms))
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, caser.String(strings.ToLower(initialism)))
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
}

// UderscoreToUpperCamelCase 下划线单词转首字母大写驼峰单词
func UderscoreToUpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = cases.Title(language.Und).String(s)
	return strings.Replace(s, " ", "", -1)
}

// UderscoreToLowerCamelCase 下划线单词转首字母小写驼峰单词
func UderscoreToLowerCamelCase(s string) string {
	s = UderscoreToUpperCamelCase(s)
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

// CamelCaseToUdnderscore 驼峰单词转下划线单词
func CamelCaseToUdnderscore(s string) string {
	s = commonInitialismsReplacer.Replace(s)

	output := make([]rune, 0, utf8.RuneCountInString(s)+2)
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
		} else {
			if unicode.IsUpper(r) {
				output = append(output, '_')
			}

			output = append(output, unicode.ToLower(r))
		}
	}

	return string(output)
}
