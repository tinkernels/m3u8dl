package ffmpeg

import (
	"fmt"
	mu "github.com/tinkernels/m3u8dl/m3u8"
	"os"
	"os/exec"
)

func MergeM3u8Ts(m3u8 string) {
	//ffmpeg -protocol_whitelist file,crypto -allowed_extensions ALL -i input.m3u8 -bsf:a aac_adtstoasc -c copy output.mp4
	cmd_ := exec.Command("ffmpeg",
		"-y",
		"-protocol_whitelist",
		"file,crypto",
		"-allowed_extensions",
		"ALL",
		"-i",
		mu.FfmpegInputM3u8File,
		"-bsf:a",
		"aac_adtstoasc",
		"-c",
		"copy",
		"output.mp4")
	cmd_.Stdout = os.Stdout
	cmd_.Stderr = os.Stderr
	err := cmd_.Start()
	if err != nil {
		fmt.Println("ffmpeg command error: ", err)
	}
	fmt.Println("Waiting for command to finish...")
	err = cmd_.Wait()
}
