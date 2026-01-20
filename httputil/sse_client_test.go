package httputil_test

import (
	"net/url"
	"testing"

	"github.com/Jsharkc/mygopkg/httputil"
)

func TestSSEPost(t *testing.T) {
	gpturl, err := url.Parse("http://192.168.1.222")
	if err != nil {
		t.Error(err)
	}
	gpturl.Path = "/api/completion"
	gpturl.RawQuery = "model=gpt&force=true&stream=true"
	data := map[string]any{"messages": []map[string]string{{"role": "user", "content": "你好！"}}, "temperature": 0.7}
	_, err = httputil.SSEPostJSON(gpturl.String(), data, func(p []byte) error {
		print(string(p))
		return nil
	})
	print(err)
}
