package main

import (
	"flag"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/tinkernels/m3u8dl/fetcher"
	"github.com/tinkernels/m3u8dl/ffmpeg"
	"github.com/tinkernels/m3u8dl/m3u8"
	"github.com/tinkernels/m3u8dl/utils"
	"net/url"
	"os"
	"sync"
)

const DefaultFetcherConcurrency = 8
const AntPoolMaxBlocking = 1024 * 1024

var (
	urlFlag         = flag.String("url", "", "url of m3u8 file")
	concurrencyFlag = flag.Int("c", DefaultFetcherConcurrency, "concurrency to fetch media segments")
	workDirFlag     = flag.String("w", "", "work dir")
)

func main() {
	flag.Parse()

	// fetcher concurrency
	concurrency_ := *concurrencyFlag
	if concurrency_ <= 0 {
		concurrency_ = DefaultFetcherConcurrency
	}
	// url
	urlStr_ := *urlFlag
	url_, err := url.Parse(urlStr_)
	if err != nil {
		fmt.Println("url not valid: ", err)
		return
	}

	workDir_ := utils.GetCWD()
	if *workDirFlag != "" {
		workDir_ = utils.AbsPath(*workDirFlag)
	}
	fmt.Printf("work dir: %s\n", workDir_)

	start(url_, concurrency_, workDir_)
}

func start(urlParam *url.URL, concurrency int, workDir string) {
	// Ensure work dir exists
	wd_ := utils.AbsPath(workDir)
	if !utils.PathExists(wd_) {
		err := os.MkdirAll(wd_, os.ModePerm)
		if err != nil {
			fmt.Println("create directory failed: ", err)
			panic(err)
		}
	}
	err := os.Chdir(wd_)
	if err != nil {
		fmt.Println("chdir failed: ", err)
		panic(err)
	}
	// using ants pool to limit concurrency
	antPool_, _ := ants.NewPool(concurrency,
		ants.WithPanicHandler(func(i interface{}) {
			fmt.Println("ants panic", i)
		}),
		ants.WithMaxBlockingTasks(AntPoolMaxBlocking),
		ants.WithNonblocking(false),
		ants.WithPreAlloc(true),
	)
	defer antPool_.Release()

	ffmpegInputM3u8, fetchTasks := m3u8.HandleM3u8Url(urlParam)

	var wg sync.WaitGroup
	for _, task := range fetchTasks {
		task_ := task
		wg.Add(1)
		err := antPool_.Submit(func() {
			_ = fetcher.Fetch(task_.UrlStr, task_.DstPath)
			wg.Done()
		})
		if err != nil {
			wg.Done()
			fmt.Printf("ants pool error: %v\n", err)
			return
		}
	}
	wg.Wait()

	fmt.Printf("ffmpeg will merge %s\n", ffmpegInputM3u8)
	ffmpeg.MergeM3u8Ts(ffmpegInputM3u8)
}
