package ffmpeg

import (
	"os"
	"testing"
)

func TestMergeM3u8Ts(t *testing.T) {
	type args struct {
		m3u8 string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "m3u8",
			args: args{
				m3u8: "local.m3u8",
			},
		},
	}
	_ = os.Chdir("../test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MergeM3u8Ts(tt.args.m3u8)
		})
	}
}
