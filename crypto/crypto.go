package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

func Md5(text string) string {
	hashMd5 := md5.New()
	_, _ = io.WriteString(hashMd5, text)
	return hex.EncodeToString(hashMd5.Sum(nil))
}

func Md5Buf(buf []byte) string {
	hashMd5 := md5.New()
	hashMd5.Write(buf)
	return hex.EncodeToString(hashMd5.Sum(nil))
}

// Md5File 计算文件内容的 md5
// 请确保文件能正常打开，如果不能，返回文件路径的 md5
func Md5File(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		return Md5(filename)
	}
	defer f.Close()

	return Md5Stream(f)
}

func Md5Stream(r io.Reader) string {
	var buf = make([]byte, 4096)
	hashMd5 := md5.New()
	for {
		n, err := r.Read(buf)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			break
		}

		hashMd5.Write(buf[:n])
	}

	return hex.EncodeToString(hashMd5.Sum(nil))
}

func Sha1(text string) string {
	hashSha1 := sha1.New()
	_, _ = io.WriteString(hashSha1, text)
	return hex.EncodeToString(hashSha1.Sum(nil))
}

func Sha1Buf(buf []byte) string {
	hashSha1 := sha1.New()
	hashSha1.Write(buf)
	return hex.EncodeToString(hashSha1.Sum(nil))
}

func Sha256(text string) string {
	hashSha256 := sha256.New()
	_, _ = io.WriteString(hashSha256, text)
	return hex.EncodeToString(hashSha256.Sum(nil))
}

func Sha256Base64(text string) string {
	hashSha256 := sha256.New()
	_, _ = io.WriteString(hashSha256, text)
	return base64.StdEncoding.EncodeToString(hashSha256.Sum(nil))
}

func Sha256Buf(buf []byte) string {
	hashSha256 := sha256.New()
	hashSha256.Write(buf)
	return hex.EncodeToString(hashSha256.Sum(nil))
}

func Sha512(text string) string {
	hashSha512 := sha512.New()
	_, _ = io.WriteString(hashSha512, text)
	return hex.EncodeToString(hashSha512.Sum(nil))
}

func Sha512Buf(buf []byte) string {
	hashSha512 := sha512.New()
	hashSha512.Write(buf)
	return hex.EncodeToString(hashSha512.Sum(nil))
}

// HmacSha1 是 HMAC Sha1 编码，返回原始的字节数组
func HmacSha1(text string, key []byte) []byte {
	hmacHash := hmac.New(sha1.New, key)
	_, _ = io.WriteString(hmacHash, text)
	return hmacHash.Sum(nil)
}

// HmacSha1 是 HMAC Sha1 编码，返回值做了十六进制处理
func HmacSha1Hex(text string, key []byte) string {
	return hex.EncodeToString(HmacSha1(text, key))
}

// HmacSha1 是 HMAC Sha1 编码，返回值做了十六进制处理
func HmacSha1Buf(text, key []byte) string {
	return hex.EncodeToString(hmacSha1Buf(text, key))
}

// HmacSha1 是 HMAC Sha1 编码，返回原始的字节数组
func hmacSha1Buf(text, key []byte) []byte {
	hmacHash := hmac.New(sha1.New, key)
	hmacHash.Write(text)
	return hmacHash.Sum(nil)
}

// HmacSha1 是 HMAC Sha1 编码，返回值使用 base64 URLEncoding 处理
func HmacSha1URLEncode(text, key []byte) string {
	return base64.URLEncoding.EncodeToString(hmacSha1Buf(text, key))
}

// HmacSha1 是 HMAC Sha1 编码，返回值使用 base64 Encoding 处理
func HmacSha1StdEncode(text, key []byte) string {
	return base64.StdEncoding.EncodeToString(hmacSha1Buf(text, key))
}

// HmacSha256 是 HMAC Sha256 编码，返回原始的字节数组
func HmacSha256(text string, key []byte) []byte {
	hmacHash := hmac.New(sha256.New, key)
	_, _ = io.WriteString(hmacHash, text)
	return hmacHash.Sum(nil)
}

// HmacSha256Hex 是 HMAC Sha256 编码，返回值做了十六进制处理
func HmacSha256Hex(text string, key []byte) string {
	return hex.EncodeToString(HmacSha256(text, key))
}

// AesGenerate 自动生成一个aes密钥
// bits参数必须为128，192，256之一
func AesGenerate(bits int) ([]byte, error) {
	if bits != 128 && bits != 192 && bits != 256 {
		return nil, errors.New("密钥长度必须为128，192，256中的某一个")
	}
	key, err := rand.Prime(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return key.Bytes(), nil
}

// AesIvGenerate 自动创建一个iv值
func AesIvGenerate() ([]byte, error) {
	key, err := rand.Prime(rand.Reader, 128)
	if err != nil {
		return nil, err
	}
	return key.Bytes(), nil
}

// AesEncrypt 使用binary string传递密钥和iv的方式（兼容旧实现）
func AesEncrypt(buf []byte, encryptKey string, iv string) (string, error) {
	tmp, err := AesByteKeyEncrypt(buf, []byte(encryptKey), []byte(iv))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(tmp), nil
}

// AesEncrypt 使用Base64 string传递密钥和iv的方式
func AesBase64Encrypt(buf []byte, base64EncryptKey string, base64iv string) ([]byte, error) {
	encryptKey, err := base64.StdEncoding.DecodeString(base64EncryptKey)
	if err != nil {
		return nil, err
	}
	iv, err := base64.StdEncoding.DecodeString(base64iv)
	if err != nil {
		return nil, err
	}
	return AesByteKeyEncrypt(buf, encryptKey, iv)
}

// AesByteKeyEncrypt AES加密 填充方式：PKCS7Padding
func AesByteKeyEncrypt(buf, encryptKey, iv []byte) ([]byte, error) {
	aesBlockEncrypter, err := aes.NewCipher(encryptKey)
	if err != nil {
		return nil, err
	}
	content := pkcs7Padding(buf, aesBlockEncrypter.BlockSize())
	encrypted := make([]byte, len(content))
	if err != nil {
		return nil, err
	}
	aesEncrypter := cipher.NewCBCEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.CryptBlocks(encrypted, content)
	return encrypted, nil
}

// AesDecrypt 使用binary string传递密钥和iv的方式（兼容旧实现）
func AesDecrypt(src string, encryptKey string, iv string) (data []byte, err error) {
	byt, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}
	return AesByteKeyDecrypt(byt, []byte(encryptKey), []byte(iv))
}

// AesDecrypt 使用Base64 string传递密钥和iv的方式
func AesBase64Decrypt(src []byte, base64EncryptKey string, base64iv string) (data []byte, err error) {
	encryptKey, err := base64.StdEncoding.DecodeString(base64EncryptKey)
	if err != nil {
		return nil, err
	}
	iv, err := base64.StdEncoding.DecodeString(base64iv)
	if err != nil {
		return nil, err
	}
	return AesByteKeyDecrypt(src, encryptKey, iv)
}

// AesByteKeyDecrypt AES加密 buf：加密buf，encryptKey: 加密秘钥，iv：偏移量
func AesByteKeyDecrypt(src []byte, encryptKey, iv []byte) (data []byte, err error) {
	decrypted := make([]byte, len(src))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher(encryptKey)
	if err != nil {
		return nil, err
	}
	aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.CryptBlocks(decrypted, src)
	return pkcs7Trimming(decrypted), nil
}

func pkcs7Padding(buf []byte, blockSize int) []byte {
	padding := blockSize - len(buf)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(buf, padText...)
}

func pkcs7Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
