package fetcher

import (
	"fmt"
	"github.com/gojek/heimdall/v7/hystrix"
	"github.com/tinkernels/m3u8dl/utils"
	"io"
	"net/http"
	"os"
	"time"
)

const FetchRetry = 5
const HTTPTimeout = 30 * time.Minute
const HystrixTimeout = FetchRetry * HTTPTimeout

func NewHttpClient() *hystrix.Client {
	return hystrix.NewClient(
		hystrix.WithCommandName("http request for m3u8dl"),
		hystrix.WithRetryCount(FetchRetry),
		hystrix.WithHTTPTimeout(HTTPTimeout),
		hystrix.WithHystrixTimeout(HystrixTimeout),
	)
}

// Fetch http file to local
func Fetch(urlStr string, dst string) (err error) {
	c_ := NewHttpClient()
	rsp_, err := c_.Get(urlStr, nil)
	defer func(resp *http.Response) { _ = resp.Body.Close() }(rsp_)
	if err != nil {
		fmt.Println("fetch error: ", err)
		return
	}
	out, err := os.Create(utils.AbsPath(dst))
	if err != nil {
		fmt.Println("create file error: ", err)
		return
	}
	defer func(out *os.File) { _ = out.Close() }(out)
	wBytes, err := io.Copy(out, rsp_.Body)
	if err != nil {
		fmt.Println("write file error: ", err)
		return
	}
	fmt.Printf("downloaded %s to %s: %dB\n", urlStr, dst, wBytes)
	return
}
