package utils

import (
	"net/url"
	"path"
	"testing"
)

func Test_url(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "good http url",
			url:  "http://127.0.0.1:8080/test.m3u8?a=1&b=2",
		},
		{
			name: "good http url",
			url:  "127.0.0.1/test.m3u8?a=1&b=2",
		},
		{
			name: "bad http url",
			url:  "://127.0.0.1:8080",
		},
		{
			name: "bad http url",
			url:  "/path1/path2/file.m3u8",
		},
		{
			name: "good file url",
			url:  "test.m3u8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url_, err := url.Parse(tt.url)
			if err != nil {
				t.Logf("%s: Parse() error = %v", tt.name, err)
			} else {
				t.Logf("%s: scheme: %s, host: %s, path: %s, dir: %s, basename: %s, query: %s",
					tt.name, url_.Scheme, url_.Host, url_.Path, path.Dir(url_.Path), path.Base(url_.Path), url_.Query())
			}
		})
	}
}
