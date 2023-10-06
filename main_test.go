package main

import (
	"net/url"
	"testing"
)

func Test_start(t *testing.T) {
	type args struct {
		urlParam    *url.URL
		concurrency int
	}
	url_, _ := url.Parse(
		"http://localhost:8080/videos/64ba3cc3ca1e274ca6425a2b/index.m3u8?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3MiOiJ2aWV3IiwiaWF0IjoxNjk2Mjk0NTY5LCJleHAiOjE2OTYyOTQ2Njl9.cpCZ5b20adlYy5D-vgYGR57nEF5mEZ-C3mxP5z3ql0w",
	)
	tests := []struct {
		name string
		args args
	}{
		{
			name: "m3u8",
			args: args{
				urlParam:    url_,
				concurrency: 8,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start(tt.args.urlParam, tt.args.concurrency, "test")
		})
	}
}
