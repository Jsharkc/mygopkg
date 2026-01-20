package idutil

import (
	"log"

	"github.com/jaevor/go-nanoid"
)

const defaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const alphabetWithSpecChar = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

var (
	NanoID            func() string
	NanoIDWithLen     func(int) string
	SpecNanoIDWithLen func(int) string // 生成含有特殊ascii字符的id
)

func init() {
	var err error
	NanoID, err = nanoid.CustomASCII(defaultAlphabet, 22)
	if err != nil {
		log.Fatal(err)
	}

	NanoIDWithLen = func(length int) string {
		f, _ := nanoid.CustomASCII(defaultAlphabet, length)
		return f()
	}

	SpecNanoIDWithLen = func(length int) string {
		f, _ := nanoid.CustomASCII(alphabetWithSpecChar, length)
		return f()
	}
}
