package m3u8

import (
	"net/url"
	"os"
	"reflect"
	"testing"
)

func Test_writeTasks(t *testing.T) {
	type args struct {
		tasks []FetchTask
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "m3u8",
			args: args{
				tasks: []FetchTask{
					{
						UrlStr:  "http://localhost:8080/videos/64ba3cc3ca1e274ca6425a2b/index.m3u8",
						DstPath: "test.m3u8",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeTasks(tt.args.tasks); (err != nil) != tt.wantErr {
				t.Errorf("writeTasks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	_ = os.Remove(FetchTaskFilePath)
}

func Test_handleM3u8KeyLine(t *testing.T) {
	type args struct {
		keyLine string
	}
	tests := []struct {
		name               string
		args               args
		wantKeyLine4Ffmpeg string
		wantKeyUri         string
	}{
		{
			name: "m3u8-urlmid",
			args: args{
				keyLine: `#EXT-X-KEY:METHOD=AES-128,URI="http://localhost//videos/64ba3cc3ca1e274ca6425a2b/ts.key",IV=0x00000000000000000000000000000000`,
			},
		},
		{
			name: "m3u8-urlbegin",
			args: args{
				keyLine: `#EXT-X-KEY:URI="http://localhost//videos/64ba3cc3ca1e274ca6425a2b/ts.key",METHOD=AES-128,IV=0x00000000000000000000000000000000`,
			},
		},
		{
			name: "m3u8-urlend",
			args: args{
				keyLine: `#EXT-X-KEY:METHOD=AES-128,IV=0x00000000000000000000000000000000,URI="http://localhost//videos/64ba3cc3ca1e274ca6425a2b/ts.key"`,
			},
		},
		{
			name: "m3u8-noquote",
			args: args{
				keyLine: `#EXT-X-KEY:METHOD=AES-128,IV=0x00000000000000000000000000000000,URI=http://localhost//videos/64ba3cc3ca1e274ca6425a2b/ts.key`,
			},
		},
		{
			name: "m3u8-noquote-mid",
			args: args{
				keyLine: `#EXT-X-KEY:METHOD=AES-128,URI=http://localhost//videos/64ba3cc3ca1e274ca6425a2b/ts.key,IV=0x00000000000000000000000000000000`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeyLine4Ffmpeg, gotKeyUri := handleM3u8KeyLine(tt.args.keyLine)
			if gotKeyLine4Ffmpeg != tt.wantKeyLine4Ffmpeg {
				t.Logf("handleM3u8KeyLine() gotKeyLine4Ffmpeg = %v, want %v", gotKeyLine4Ffmpeg, tt.wantKeyLine4Ffmpeg)
			}
			if gotKeyUri != tt.wantKeyUri {
				t.Logf("handleM3u8KeyLine() gotKeyUri = %v, want %v", gotKeyUri, tt.wantKeyUri)
			}
		})
	}
}

func Test_genFetchTask4Url(t *testing.T) {
	type args struct {
		fileUrlStr string
		urlP       *url.URL
	}
	tests := []struct {
		name          string
		args          args
		wantFetchTask *FetchTask
	}{
		// Add test cases.
		{
			name: "m3u8",
			args: args{
				fileUrlStr: "http://localhost:8080/videos/64ba3cc3ca1e274ca6425a2b/index.m3u8",
				urlP: &url.URL{
					Scheme: "https",
					Host:   "localhost:8080",
					Path:   "/videos/64ba3cc3ca1e274ca6425a2b/index.m3u8",
				},
			},
		},
		{
			name: "m3u8",
			args: args{
				fileUrlStr: "/videos/64ba3cc3ca1e274ca6425a2b/index.m3u8",
				urlP: &url.URL{
					Scheme: "https",
					Host:   "localhost:8080",
					Path:   "/videos/64ba3cc3ca1e274ca6425a2b/index.m3u8",
				},
			},
		},
		{
			name: "m3u8",
			args: args{
				fileUrlStr: "index.m3u8",
				urlP: &url.URL{
					Scheme: "https",
					Host:   "localhost:8080",
					Path:   "/videos/64ba3cc3ca1e274ca6425a2b/index.m3u8",
				},
			},
		},
		{
			name: "m3u8",
			args: args{
				fileUrlStr: "index.m3u8",
				urlP: &url.URL{
					Scheme: "",
					Host:   "",
					Path:   "/tmp/64ba3cc3ca1e274ca6425a2b/index.m3u8",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFetchTask := genFetchTask4Url(tt.args.fileUrlStr, tt.args.urlP); !reflect.DeepEqual(gotFetchTask, tt.wantFetchTask) {
				t.Logf("genFetchTask4Url() = %v, want %v", gotFetchTask, tt.wantFetchTask)
			}
		})
	}
}

func TestHandleM3u8Url(t *testing.T) {
	type args struct {
		urlP *url.URL
	}
	url_, _ := url.Parse("http://localhost:8080/videos/64ba3cc3ca1e274ca6425a2b/index.m3u8?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3MiOiJ2aWV3IiwiaWF0IjoxNjk2Mjg1OTUyLCJleHAiOjE2OTYyODYwNTJ9.VUw-I-BUGNZqq67btHL1IunbDEfYOWuenmYhYS8Wtpc")
	tests := []struct {
		name                  string
		args                  args
		wantM3u8FileForFfmpeg string
		wantFetchTasks        []FetchTask
	}{
		// Add test cases.
		{
			name: "m3u8",
			args: args{
				urlP: url_,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotM3u8FileForFfmpeg, gotFetchTasks := HandleM3u8Url(tt.args.urlP)
			if gotM3u8FileForFfmpeg != tt.wantM3u8FileForFfmpeg {
				t.Errorf("HandleM3u8Url() gotM3u8FileForFfmpeg = %v, want %v", gotM3u8FileForFfmpeg, tt.wantM3u8FileForFfmpeg)
			}
			if !reflect.DeepEqual(gotFetchTasks, tt.wantFetchTasks) {
				t.Errorf("HandleM3u8Url() gotFetchTasks = %v, want %v", gotFetchTasks, tt.wantFetchTasks)
			}
		})
	}
}
