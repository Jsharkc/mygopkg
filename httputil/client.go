// Package httputil 是一个 http client 辅助类，基于 Resty 库实现
package httputil

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Jsharkc/mygopkg/logger"
	"github.com/go-resty/resty/v2"
)

const userAgent = "Megaview (https://megaview.com)"

// Get 发送 HTTP/GET 请求，成功时回调 callbacks
func Get(url string, callbacks ...Callback) error {
	return DefaultClient.Get(url, callbacks...)
}

// GetQuery 发送 HTTP/GET 请求，query 参数通过 map 传递进来，成功时回调 callbacks
func GetQuery(url string, query map[string]string, callbacks ...Callback) error {
	return DefaultClient.GetQuery(url, query, callbacks...)
}

// GetWithContext 发送 HTTP/GET 请求，成功时回调 callbacks
// 支持 context
func GetWithContext(ctx context.Context, url string, callbacks ...Callback) error {
	return DefaultClient.GetWithContext(ctx, url, callbacks...)
}

// GetQueryWithContext 发送 HTTP/GET 请求，query 参数通过 map 传递进来，成功时回调 callbacks
// 支持 context
func GetQueryWithContext(ctx context.Context, url string, query map[string]string, callbacks ...Callback) error {
	return DefaultClient.GetQueryWithContext(ctx, url, query, callbacks...)
}

// PostJSON 发送 HTTP/POST 请求，数据格式 JSON ，成功时回调 callbacks
// Supported request body data types is `string`, `[]byte`, `struct`, `map`, `slice` and `io.Reader`.
// Body value can be pointer or non-pointer. Automatic marshalling for JSON and XML content type,
// if it is `struct`, `map`, or `slice`.
func PostJSON(url string, body any, callbacks ...Callback) error {
	return DefaultClient.PostJSON(url, body, callbacks...)
}

// PostJSONWithContext 发送 HTTP/POST 请求，数据格式 JSON ，成功时回调 callbacks
// 支持 context
func PostJSONWithContext(ctx context.Context, url string, body any, callbacks ...Callback) error {
	return DefaultClient.PostJSONWithContext(ctx, url, body, callbacks...)
}
func PostGzJSON(url string, body any, callbacks ...Callback) error {
	return DefaultClient.PostGzJSON(url, body, callbacks...)
}

// PostJSONWithContext 发送 HTTP/POST 请求，数据格式 JSON ，成功时回调 callbacks
// 支持 context
func PostGzJSONWithContext(ctx context.Context, url string, body any, callbacks ...Callback) error {
	return DefaultClient.PostGzJSONWithContext(ctx, url, body, callbacks...)
}

// PostForm 发送 HTTP/POST 请求，数据格式 x-www-form-urlencoded ，成功时回调 callbacks
func PostForm(url string, data map[string]string, callbacks ...Callback) error {
	return DefaultClient.PostForm(url, data, callbacks...)
}

// PostFormWithContext 发送 HTTP/POST 请求，数据格式 x-www-form-urlencoded ，成功时回调 callbacks
// 支持 context
func PostFormWithContext(ctx context.Context, url string, data map[string]string, callbacks ...Callback) error {
	return DefaultClient.PostFormWithContext(ctx, url, data, callbacks...)
}

// Download 下载文件，存入 output 表示的文件中
func Download(url, output string, callbacks ...Callback) error {
	return DefaultClient.Download(url, output, callbacks...)
}

type httpClient struct {
	header http.Header

	client *resty.Client
}

var DefaultClient = New()

func New() *httpClient {
	client := resty.New()
	// 默认在返回 502 Bad Gateway 时进行一次重试
	if client.RetryCount == 0 {
		client.SetRetryCount(1)
	}
	client.SetHeader("User-Agent", userAgent)
	client.SetTimeout(1 * time.Minute)
	return &httpClient{
		client: client,
	}
}

type Callback func([]byte) error

// SetRestyClient 自定义的 RestyClient，只有 httputil 中的 API 无法满足是才使用
func (c *httpClient) SetRestyClient(client *resty.Client) *httpClient {
	if c == DefaultClient {
		panic("don't change DefaultClient")
	}
	c.client = client
	c.client.SetHeader("User-Agent", userAgent)
	return c
}

func (c *httpClient) SetHeader(header http.Header) *httpClient {
	if c == DefaultClient {
		panic("don't change DefaultClient")
	}
	c.header = header
	return c
}

func (c *httpClient) SetTimeout(d time.Duration) *httpClient {
	if c == DefaultClient {
		panic("don't change DefaultClient")
	}
	c.client.SetTimeout(d)
	return c
}

func (c *httpClient) SetRetryCount(count int) *httpClient {
	if c == DefaultClient {
		panic("don't change DefaultClient")
	}
	c.client.SetRetryCount(count)
	return c
}

func (c *httpClient) Get(url string, callbacks ...Callback) error {
	succCallback := c.getCallback(callbacks...)

	resp, err := c.getReq().Get(url)
	if err != nil {
		return err
	}

	return c.dealResp(resp, succCallback)
}

func (c *httpClient) GetQuery(url string, query map[string]string, callbacks ...Callback) error {
	succCallback := c.getCallback(callbacks...)

	resp, err := c.getReq().SetQueryParams(query).Get(url)
	if err != nil {
		return err
	}

	return c.dealResp(resp, succCallback)
}

func (c *httpClient) GetWithContext(ctx context.Context, url string, callbacks ...Callback) error {
	succCallback := c.getCallback(callbacks...)

	resp, err := c.getReq().SetContext(ctx).Get(url)
	if err != nil {
		return err
	}

	return c.dealResp(resp, succCallback)
}

func (c *httpClient) GetQueryWithContext(ctx context.Context, url string, query map[string]string, callbacks ...Callback) error {
	succCallback := c.getCallback(callbacks...)

	resp, err := c.getReq().SetContext(ctx).SetQueryParams(query).Get(url)
	if err != nil {
		return err
	}

	return c.dealResp(resp, succCallback)
}

// Do 执行 http 请求。method 是标准的 HTTP Method
func (c *httpClient) Do(method, url string) (*resty.Response, error) {
	return c.getReq().Execute(method, url)
}

// 获取 Request
func (c *httpClient) Request() *resty.Request {
	return c.getReq()
}

func (c *httpClient) PostJSON(url string, body any, callbacks ...Callback) error {
	succCallback := c.getCallback(callbacks...)

	resp, err := c.getReq().SetHeader("Content-Type", "application/json").
		SetBody(body).Post(url)
	if err != nil {
		return err
	}

	return c.dealResp(resp, succCallback)
}

func (c *httpClient) PostJSONWithContext(ctx context.Context, url string, body any, callbacks ...Callback) error {
	succCallback := c.getCallback(callbacks...)

	req := c.getReq().SetContext(ctx).SetHeader("Content-Type", "application/json")
	traceId, ok := ctx.Value(logger.TraceIDKey).(string)
	if ok {
		req.SetHeader("x-trace-id", traceId)
	}

	resp, err := req.SetBody(body).Post(url)
	if err != nil {
		return err
	}

	return c.dealResp(resp, succCallback)
}

func (*httpClient) gzJsonBody(body any) ([]byte, error) {
	byteBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	var zBuf bytes.Buffer
	zw := gzip.NewWriter(&zBuf)
	if _, err := zw.Write(byteBody); err != nil {
		return nil, err
	}
	zw.Close()
	return zBuf.Bytes(), nil
}

func (c *httpClient) PostGzJSON(url string, body any, callbacks ...Callback) error {
	succCallback := c.getCallback(callbacks...)
	gzBody, err := c.gzJsonBody(body)
	if err != nil {
		return err
	}
	resp, err := c.getReq().SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(gzBody).Post(url)
	if err != nil {
		return err
	}

	return c.dealResp(resp, succCallback)
}

func (c *httpClient) PostGzJSONWithContext(ctx context.Context, url string, body any, callbacks ...Callback) error {
	succCallback := c.getCallback(callbacks...)
	gzBody, err := c.gzJsonBody(body)
	if err != nil {
		return err
	}
	resp, err := c.getReq().SetContext(ctx).SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(gzBody).Post(url)
	if err != nil {
		return err
	}

	return c.dealResp(resp, succCallback)
}

func (c *httpClient) PostForm(url string, data map[string]string, callbacks ...Callback) error {
	succCallback := c.getCallback(callbacks...)

	resp, err := c.getReq().SetFormData(data).Post(url)
	if err != nil {
		return err
	}

	return c.dealResp(resp, succCallback)
}

func (c *httpClient) PostFormWithContext(ctx context.Context, url string, data map[string]string, callbacks ...Callback) error {
	succCallback := c.getCallback(callbacks...)

	resp, err := c.getReq().SetContext(ctx).SetFormData(data).Post(url)
	if err != nil {
		return err
	}

	return c.dealResp(resp, succCallback)
}

func (c *httpClient) Download(url, output string, callbacks ...Callback) error {
	succCallback := c.getCallback(callbacks...)

	resp, err := c.getReq().SetOutput(output).Get(url)
	if err != nil {
		return err
	}

	return c.dealResp(resp, succCallback)
}

func (c *httpClient) getReq() *resty.Request {
	return c.client.R().SetHeaderMultiValues(c.header).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			// 如果发生错误或响应为 nil，不进行重试
			if err != nil || r == nil {
				return false
			}
			// 502 Bad Gateway 重试
			return r.StatusCode() == http.StatusBadGateway
		})
}

func (c *httpClient) getCallback(callbacks ...Callback) Callback {
	var succCallback Callback

	if len(callbacks) > 0 {
		succCallback = callbacks[0]
	}

	return succCallback
}

func (c *httpClient) dealResp(resp *resty.Response, succCallback Callback) error {
	if resp.IsSuccess() {
		if succCallback != nil {
			return succCallback(resp.Body())
		}

		return nil
	}

	originURL := resp.Request.URL
	// return errors.New("http status:" + resp.Status() + "; url:" + originURL)
	return fmt.Errorf("http status %s ; body=%s ; url %s", resp.Status(), resp.Body(), originURL)
}

func (c *httpClient) PutJSONWithContext(ctx context.Context, url string, body any, callbacks ...Callback) error {
	succCallback := c.getCallback(callbacks...)

	req := c.getReq().SetContext(ctx).SetHeader("Content-Type", "application/json")
	traceId, ok := ctx.Value(logger.TraceIDKey).(string)
	if ok {
		req.SetHeader("x-trace-id", traceId)
	}

	resp, err := req.SetBody(body).Put(url)
	if err != nil {
		return err
	}

	return c.dealResp(resp, succCallback)
}
