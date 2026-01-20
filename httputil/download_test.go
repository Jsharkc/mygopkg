package httputil_test

import (
	"context"
	"io"
	"testing"

	"github.com/Jsharkc/mygopkg/httputil"
)

func TestDownloadToReader(t *testing.T) {
	reader, err := httputil.DefaultDownloadClient.DownloadToReader(context.Background(), "https://sale-test.oss-cn-zhangjiakou.aliyuncs.com/20240712/c73db08d8f659f86485e2a0d7157392c.amr?Expires=1721131252&OSSAccessKeyId=TMP.3KidEay5mmPYqcJvEZgqfLQabh4uWHvz6MQJKHr7u29pW9UXJ8T5NWWGapW5szHmZkTCDJMWhDbPJkWiqT6yGa7sAKdrKV&Signature=94Z7kHxKOch47hFrM8fKmk4afOg%3D")
	if err != nil {
		t.Error(err)
	}
	bs, err := io.ReadAll(reader)
	if err != nil {
		t.Error(err)
	}
	t.Log(len(bs))
}
func TestDownloadToFile(t *testing.T) {
	err := httputil.DefaultDownloadClient.DownloadToFile(context.Background(), "https://sale-test.oss-cn-zhangjiakou.aliyuncs.com/20240712/c73db08d8f659f86485e2a0d7157392c.amr?Expires=1721131252&OSSAccessKeyId=TMP.3KidEay5mmPYqcJvEZgqfLQabh4uWHvz6MQJKHr7u29pW9UXJ8T5NWWGapW5szHmZkTCDJMWhDbPJkWiqT6yGa7sAKdrKV&Signature=94Z7kHxKOch47hFrM8fKmk4afOg%3D", "c73db08d8f659f86485e2a0d7157392c.amr")
	if err != nil {
		t.Error(err)
	}
}
