package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
)

// RsaGenerate 创建一组RSA密钥
// bits参数为密钥位长
func RsaGenerate(bits int) ([]byte, []byte, error) {
	private, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	privateKey, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := x509.MarshalPKIXPublicKey(&private.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}

// RSAOAEPEncrypt OAEP模式加密
// label 是盐可以为任意值，用于验证解密结果
func RSAOAEPEncrypt(target, publicKey, label []byte) ([]byte, error) {
	pubKey, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	rsaPubKey := pubKey.(*rsa.PublicKey)
	data, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, rsaPubKey, target, label)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// RSAOAEPDecrypt OAEP模式解密
// label需要使用与加密相同的值 是盐可以为任意值，用于验证解密结果
func RSAOAEPDecrypt(target, privateKey, label []byte) ([]byte, error) {
	privKey, err := x509.ParsePKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	rsaPrivKey := privKey.(*rsa.PrivateKey)
	return rsa.DecryptOAEP(sha512.New(), rand.Reader, rsaPrivKey, target, label)
}

// RSAEncrypt 普通模式加密
func RSAEncrypt(target, publicKey []byte) ([]byte, error) {
	pubKey, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	rsaPubKey := pubKey.(*rsa.PublicKey)
	data, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, target)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// RSADecrypt 普通模式解密
func RSADecrypt(target, privateKey []byte) ([]byte, error) {
	privKey, err := x509.ParsePKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	rsaPrivKey := privKey.(*rsa.PrivateKey)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, rsaPrivKey, target)
}

// RSAOAEPEncrypt OAEP模式加密
// label 是盐可以为任意值，用于验证解密结果
// 返回值为Base64编码密文
func RSAOAEPEncryptBase64(target, publicKey, label []byte) (string, error) {
	data, err := RSAOAEPEncrypt(target, publicKey, label)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// RSAOAEPDecrypt OAEP模式解密
// label需要使用与加密相同的值 是盐可以为任意值，用于验证解密结果
// 输入target为Base64编码密文
func RSAOAEPDecryptBase64(target string, privateKey, label []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(target)
	if err != nil {
		return nil, err
	}

	return RSAOAEPDecrypt(data, privateKey, label)
}

// RSAEncrypt 普通模式加密
// 返回值为Base64编码密文
func RSAEncryptBase64(target, publicKey []byte) (string, error) {
	data, err := RSAEncrypt(target, publicKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// RSADecrypt 普通模式解密
// 输入target为Base64编码密文
func RSADecryptBase64(target string, privateKey []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(target)
	if err != nil {
		return nil, err
	}

	return RSADecrypt(data, privateKey)
}

// RSASign RSA签名
// 输入target待签名数据
func RSASign(target, privateKey []byte) ([]byte, error) {
	privKey, err := x509.ParsePKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	msgHash := sha512.New()
	_, err = msgHash.Write(target)
	if err != nil {
		return nil, err
	}
	msgHashSum := msgHash.Sum(nil)
	rsaPrivKey := privKey.(*rsa.PrivateKey)

	return rsa.SignPSS(rand.Reader, rsaPrivKey, crypto.SHA512, msgHashSum, &rsa.PSSOptions{SaltLength: 16, Hash: crypto.SHA512})
}

// RSAVerify RSA签名验证
// 输入target待验证数据
// sig为签名
func RSAVerify(target, publicKey, sig []byte) error {
	pubKey, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	rsaPubKey := pubKey.(*rsa.PublicKey)
	msgHash := sha512.New()
	_, err = msgHash.Write(target)
	if err != nil {
		return err
	}
	msgHashSum := msgHash.Sum(nil)

	return rsa.VerifyPSS(rsaPubKey, crypto.SHA512, msgHashSum, sig, &rsa.PSSOptions{SaltLength: 16, Hash: crypto.SHA512})
}
