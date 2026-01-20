package stringutil

import "github.com/mozillazg/go-pinyin"

// Pinyin 汉字转为拼音，其他字符原样保留
func Pinyin(str string) string {
	arg := pinyin.NewArgs()
	arg.Separator = ""
	arg.Fallback = func(r rune, args pinyin.Args) []string {
		return []string{string(r)}
	}
	py := pinyin.Slug(str, arg)

	if py == "" {
		py = str
	}

	return py
}
