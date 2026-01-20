package crypto_test

import (
	"testing"

	"github.com/Jsharkc/mygopkg/crypto"
)

func TestRSA(t *testing.T) {
	priv, pub, err := crypto.RsaGenerate(2048)

	data := "abc"
	if err != nil {
		t.Error()
	}
	a, err := crypto.RSAEncrypt([]byte(data), pub)
	if err != nil {
		t.Error()
	}
	s, err := crypto.RSADecrypt(a, priv)
	if err != nil {
		t.Error()
	}
	if string(s) != data {
		t.Error()
	}
}
func TestRSABase64(t *testing.T) {
	priv, pub, err := crypto.RsaGenerate(2048)

	data := "abc"
	if err != nil {
		t.Error()
	}
	a, err := crypto.RSAEncryptBase64([]byte(data), pub)
	if err != nil {
		t.Error()
	}
	s, err := crypto.RSADecryptBase64(a, priv)
	if err != nil {
		t.Error()
	}
	if string(s) != data {
		t.Error()
	}
}
func TestRSAOAEP(t *testing.T) {
	priv, pub, err := crypto.RsaGenerate(2048)
	data := "abc"
	if err != nil {
		t.Error()
	}
	label := []byte("123456")
	a, err := crypto.RSAOAEPEncrypt([]byte(data), pub, label)
	if err != nil {
		t.Error()
	}
	s, err := crypto.RSAOAEPDecrypt(a, priv, label)
	if err != nil {
		t.Error()
	}
	if string(s) != data {
		t.Error()
	}
}
func TestRSAOAEPBase64(t *testing.T) {
	priv, pub, err := crypto.RsaGenerate(2048)
	data := "abc"
	if err != nil {
		t.Error()
	}
	label := []byte("123456")
	a, err := crypto.RSAOAEPEncryptBase64([]byte(data), pub, label)
	if err != nil {
		t.Error()
	}
	s, err := crypto.RSAOAEPDecryptBase64(a, priv, label)
	if err != nil {
		t.Error()
	}
	if string(s) != data {
		t.Error()
	}
}
func TestRSASign(t *testing.T) {
	priv, pub, err := crypto.RsaGenerate(2048)
	data1, data2 := "abc", "abc"
	if err != nil {
		t.Error()
	}
	a, err := crypto.RSASign([]byte(data1), priv)
	if err != nil {
		t.Error()
	}
	err = crypto.RSAVerify([]byte(data2), pub, a)
	if (err == nil) != (data1 == data2) {
		t.Error()
	}
}
func FuzzXxx(f *testing.F) {
	f.Add("abcsd", "asssas")
	f.Add("asd", "asd")
	f.Fuzz(func(t *testing.T, data1, data2 string) {
		priv, pub, err := crypto.RsaGenerate(2048)
		if err != nil {
			t.Error()
		}
		a, err := crypto.RSASign([]byte(data1), priv)
		if err != nil {
			t.Error()
		}
		err = crypto.RSAVerify([]byte(data2), pub, a)
		if (err == nil) != (data1 == data2) {
			t.Error()
		}
	})

}
