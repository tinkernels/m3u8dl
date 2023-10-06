package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"

	mu "github.com/tinkernels/m3u8dl/m3u8"
)

func MergeM3u8Ts(m3u8 string) {
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
	if err != nil {
		fmt.Println("ffmpeg command error: ", err)
	}
}
