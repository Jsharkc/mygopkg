package httputil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Jsharkc/mygopkg/fileutil"
)

type downloadHttpClient struct {
	retryCount     int
	client         *http.Client
	retryCondition func(resp *http.Response) bool
}

var DefaultDownloadClient = NewDownload()

func NewDownload() *downloadHttpClient {
	client := &http.Client{Timeout: 30 * time.Minute}
	// 默认在返回 502 Bad Gateway 时进行一次重试
	retrycondition := func(resp *http.Response) bool {
		return resp.StatusCode == http.StatusBadGateway
	}
	return &downloadHttpClient{
		retryCount:     1,
		client:         client,
		retryCondition: retrycondition,
	}
}
func (c *downloadHttpClient) doRetry(req *http.Request) (*http.Response, error) {
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
				continue
			}
			return resp, nil
		}
		return resp, err
	}

}

func (c *downloadHttpClient) DownloadToFile(ctx context.Context, url, output string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.doRetry(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed, status code: %d", resp.StatusCode)
	}

	file, err := fileutil.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *downloadHttpClient) DownloadToReader(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRetry(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed, status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
