package httputil_test

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/Jsharkc/mygopkg/httputil"
)

func TestGet(t *testing.T) {
	type args struct {
		url       string
		callbacks []httputil.Callback
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"first", args{"https://megaview.com", nil}, false},
		{"first", args{"https://megaview.com", nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := httputil.Get(tt.args.url, tt.args.callbacks...); (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		time.Sleep(1e6)
		go func() {
			defer wg.Done()
			if err := httputil.Get("http://localhost:8080/"); err != nil {
				b.Error(err)
			}
		}()
	}
	wg.Wait()
}

func TestDownload(t *testing.T) {
	err := httputil.Download("https://github.com/persepolisdm/persepolis/releases/download/3.2.0/persepolis_3.2.0.2_all.deb", "go1.19.5.src.targ.gz")
	if err != nil {
		t.Error(err)
	}
}

func Test_httpClient_Get(t *testing.T) {
	type fields struct {
		Header http.Header
	}
	type args struct {
		url       string
		callbacks []httputil.Callback
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := httputil.New()
			if err := c.Get(tt.args.url, tt.args.callbacks...); (err != nil) != tt.wantErr {
				t.Errorf("httpClient.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestGzPost(t *testing.T) {
	type args struct {
		url       string
		data      map[string]any
		callbacks []httputil.Callback
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"first", args{"http://localhost:8000", map[string]any{"1": "asd", "sadasd": "asdasdasdasdasd"}, nil}, false},
		{"first", args{"http://localhost:8000", map[string]any{"2": "sssq", "sadasd": "asdasdasdasdasd"}, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := httputil.PostGzJSONWithContext(context.Background(), tt.args.url, tt.args.data, tt.args.callbacks...); (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
