package crypto_test

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/Jsharkc/mygopkg/crypto"
)

func TestMd5(t *testing.T) {
	m := crypto.Md5("123456")
	if m != "e10adc3949ba59abbe56e057f20f883e" {
		t.Error("fail")
	}
}

func TestSha1(t *testing.T) {
	m := crypto.Sha1("123456")
	if m != "7c4a8d09ca3762af61e59520943dc26494f8941b" {
		t.Error("fail")
	}
}

func TestSha256(t *testing.T) {
	m := crypto.Sha256("123456")
	if m != "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92" {
		t.Error("fail")
	}
}

// 入参的顺序不一样
func TestShaBase64(t *testing.T) {
	secret := "0PHiMxhjOGw_e1egA202ccClUpW9UZ6V"
	stringToSign := "POST&2023-01-28T16:17:06Z&9e643236-6220-40f0-ab1f-0d08047e9e39"
	m := crypto.HmacSha1URLEncode([]byte(stringToSign), []byte(secret))
	fmt.Println(m)
	n := base64.URLEncoding.EncodeToString(Hmac_sha1(secret, stringToSign))
	fmt.Println(n)
}

func Hmac_sha1(key, code string) []byte {
	signature := hmac.New(sha1.New, []byte(key))
	signature.Write([]byte(code))
	return signature.Sum(nil)
}

func TestAesEncryptDecrypt(t *testing.T) {
	encryptKey := "abcdefghigklmnop"
	iv := "1234567890abcdef"
	text := "123456"
	e, err := crypto.AesEncrypt([]byte(text), encryptKey, iv)
	if err != nil {
		t.Error("fail")
	}
	if e != "DCiBIFfTTfNSLvQ+nZrQKA==" {
		t.Error("fail")
	}
	decode, err := crypto.AesDecrypt(e, encryptKey, iv)
	if err != nil {
		t.Error("fail")
	}
	if string(decode) != text {
		t.Error()
	}
}

func TestAesGenerate(t *testing.T) {
	s, _ := crypto.AesGenerate(256)
	a, _ := crypto.AesIvGenerate()
	println(base64.StdEncoding.EncodeToString(s), "\n", base64.StdEncoding.EncodeToString(a))
}

func FuzzAesByteKeyEncryptDecrypt(f *testing.F) {
	f.Add("ad")
	f.Add("ccc")
	f.Fuzz(func(t *testing.T, text string) {
		encryptKey, _ := crypto.AesGenerate(256)
		iv, _ := crypto.AesIvGenerate()
		e, err := crypto.AesByteKeyEncrypt([]byte(text), encryptKey, iv)
		if err != nil {
			t.Error("fail")
		}
		decode, err := crypto.AesByteKeyDecrypt(e, encryptKey, iv)
		if err != nil {
			t.Error("fail")
		}
		if string(decode) != text {
			t.Error()
		}
	})
}
