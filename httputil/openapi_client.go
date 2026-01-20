package httputil

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Jsharkc/mygopkg/crypto"
	uuid "github.com/satori/go.uuid"
)

type OpenAPIClient struct {
	*http.Client
	APPkey string
	Secert string
}

func NewOpenAPIClient(APPkey string, Secert string) *OpenAPIClient {
	return &OpenAPIClient{&http.Client{
		// ---- start 2021.12.8 ml 超时时间
		Timeout: 10 * time.Minute,
		// ---- end
	}, APPkey, Secert}
}

func megaAuthentication(HTTPMethod, secret, timeStamp, signatureNonce string) string {
	stringToSign := strings.Join([]string{HTTPMethod, timeStamp, signatureNonce}, "&")
	return crypto.HmacSha1URLEncode([]byte(stringToSign), []byte(secret))
}

func (o *OpenAPIClient) Get(requestUrl string, params url.Values) ([]byte, error) {
	// 生成 uuid
	if params == nil {
		params = make(url.Values)
	}
	signature_nonce := uuid.NewV4().String()
	time_stamp := time.Now().Format("2006-01-02T15:04:05Z")
	signature := megaAuthentication("GET", o.Secert, time_stamp, signature_nonce)

	Url, err := url.Parse(requestUrl)
	if err != nil {
		return nil, err
	}

	// 拼接数据
	params.Add("Signature", signature)
	params.Add("SignatureNonce", signature_nonce)
	params.Add("Timestamp", time_stamp)
	params.Add("AccessKeyId", o.APPkey)

	Url.RawQuery = params.Encode()
	var body []byte
	err = Get(
		Url.String(),
		func(b []byte) error {
			body = b
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return body, nil

}

func (o *OpenAPIClient) PostForm(requestUrl string, params url.Values, form map[string]string) ([]byte, error) {
	// 生成 uuid
	if params == nil {
		params = make(url.Values)
	}
	signature_nonce := uuid.NewV4().String()
	time_stamp := time.Now().Format("2006-01-02T15:04:05Z")
	signature := megaAuthentication("POST", o.Secert, time_stamp, signature_nonce)

	// 拼接数据
	params.Add("Signature", signature)
	params.Add("SignatureNonce", signature_nonce)
	params.Add("Timestamp", time_stamp)
	params.Add("AccessKeyId", o.APPkey)

	Url, err := url.Parse(requestUrl)
	if err != nil {
		return nil, err
	}
	Url.RawQuery = params.Encode()
	var body []byte
	err = PostForm(
		Url.String(),
		form,
		func(b []byte) error {
			body = b
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (o *OpenAPIClient) PostJson(requestUrl string, params url.Values, jsonData map[string]interface{}) ([]byte, error) {
	// 生成 uuid
	if params == nil {
		params = make(url.Values)
	}
	signature_nonce := uuid.NewV4().String()
	time_stamp := time.Now().Format("2006-01-02T15:04:05Z")
	signature := megaAuthentication("POST", o.Secert, time_stamp, signature_nonce)

	// 拼接数据
	params.Add("Signature", signature)
	params.Add("SignatureNonce", signature_nonce)
	params.Add("Timestamp", time_stamp)
	params.Add("AccessKeyId", o.APPkey)

	Url, err := url.Parse(requestUrl)
	if err != nil {
		return nil, err
	}
	Url.RawQuery = params.Encode()
	var body []byte
	err = PostJSON(
		Url.String(),
		jsonData,
		func(b []byte) error {
			body = b
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (o *OpenAPIClient) PostFile(requestUrl string, params url.Values, path string) ([]byte, error) {
	// 生成 uuid
	if params == nil {
		params = make(url.Values)
	}
	signature_nonce := uuid.NewV4().String()
	time_stamp := time.Now().Format("2006-01-02T15:04:05Z")
	signature := megaAuthentication("POST", o.Secert, time_stamp, signature_nonce)

	// 拼接数据
	params.Add("Signature", signature)
	params.Add("SignatureNonce", signature_nonce)
	params.Add("Timestamp", time_stamp)
	params.Add("AccessKeyId", o.APPkey)

	Url, err := url.Parse(requestUrl)
	if err != nil {
		return nil, err
	}
	fp, err := os.Open(path) // 打开文件句柄
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	body := &bytes.Buffer{}             // 初始化body参数
	writer := multipart.NewWriter(body) // 实例化multipart
	if err != nil {
		return nil, err
	}
	part, err := writer.CreateFormFile("file", filepath.Base(path)) // 创建multipart 文件字段
	if err != nil {
		return nil, err
	}
	md5hash := md5.New()
	_, err = io.Copy(md5hash, fp)
	if err != nil {
		return nil, err
	}
	params.Add("md5", hex.EncodeToString(md5hash.Sum(nil)))

	Url.RawQuery = params.Encode()
	fp.Seek(0, 0)
	_, err = io.Copy(part, fp) // 写入文件数据到multipart
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", Url.String(), body) // 新建请求
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType()) // 设置请求头,!!!非常重要，否则远端无法识别请求
	if err != nil {
		return []byte(""), err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
