package fetcher

import "testing"

func TestFetch(t *testing.T) {
	type args struct {
		urlStr string
		dst    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "m3u8",
			args: args{
				urlStr: "http://localhost:8080/videos/64ba3cc3ca1e274ca6425a2b/index.m3u8",
				dst:    "test.m3u8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = Fetch(tt.args.urlStr, tt.args.dst)
		})
	}
}
