package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Jsharkc/mygopkg/logger"
)

type Builder struct {
	strings.Builder
	f Callback
}

func (b *Builder) Write(p []byte) (n int, err error) {
	err = b.f(p)
	if err != nil {
		return 0, err
	}
	return b.Builder.Write(p)
}

type sseHttpClient struct {
	retryCount     int
	client         *http.Client
	retryCondition func(resp *http.Response) bool
}

var DefaultSSEClient = NewSSE()

func NewSSE() *sseHttpClient {
	client := &http.Client{Timeout: 10 * time.Minute}
	// 默认在返回 502 Bad Gateway 时进行一次重试
	retrycondition := func(resp *http.Response) bool {
		return resp.StatusCode == http.StatusBadGateway
	}
	return &sseHttpClient{
		retryCount:     1,
		client:         client,
		retryCondition: retrycondition,
	}
}
func (c *sseHttpClient) getCallback(callbacks ...Callback) Callback {
	var succCallback Callback

	if len(callbacks) > 0 {
		succCallback = callbacks[0]
	}

	return succCallback
}
func (c *sseHttpClient) doRetry(req *http.Request, seekfunc func(offset int64, whence int) (int64, error)) (*http.Response, error) {
	if c.retryCount == 0 {
		return c.client.Do(req)
	} else {
		var (
			resp *http.Response
			err  error
		)
		for i := 0; i <= c.retryCount; i++ {
			resp, err = c.client.Do(req)
			if err != nil {
				return resp, err
			}
			if c.retryCondition(resp) {
				seekfunc(0, 0)
				continue
			}
			return resp, nil
		}
		return resp, err
	}

}
func (c *sseHttpClient) dealResp(resp *http.Response, succCallback Callback) (string, error) {
	res := Builder{f: succCallback}
	_, err := io.Copy(&res, resp.Body)
	if err != nil {
		return "", err
	}
	return res.String(), nil
}

func (c *sseHttpClient) PostJSONWithContext(ctx context.Context, url string, body any, callbacks ...Callback) (string, error) {
	succCallback := c.getCallback(callbacks...)
	bytebody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	reader := bytes.NewReader(bytebody)
	req, err := http.NewRequestWithContext(ctx, "POST", url, reader)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	traceId, ok := ctx.Value(logger.TraceIDKey).(string)
	if ok {
		req.Header.Set("x-trace-id", traceId)
	}

	resp, err := c.doRetry(req, reader.Seek)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return c.dealResp(resp, succCallback)
}

func (c *sseHttpClient) PostJSON(url string, body any, callbacks ...Callback) (string, error) {
	return c.PostJSONWithContext(context.Background(), url, body, callbacks...)
}

func SSEPostJSON(url string, body any, callbacks ...Callback) (string, error) {
	return DefaultSSEClient.PostJSON(url, body, callbacks...)
}
func SSEPostJSONWithContext(ctx context.Context, url string, body any, callbacks ...Callback) (string, error) {
	return DefaultSSEClient.PostJSONWithContext(ctx, url, body, callbacks...)
}
