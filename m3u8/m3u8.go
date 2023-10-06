package m3u8

import (
	"bufio"
	"fmt"
	"github.com/tinkernels/m3u8dl/fetcher"
	"github.com/tinkernels/m3u8dl/utils"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

const FetchTaskFilePath = "fetch_task.yml"
const OriginalM3u8File = "remote.m3u8"
const FfmpegInputM3u8File = "local.m3u8"

type FetchTask struct {
	UrlStr  string `yaml:"url"`
	DstPath string `yaml:"dst"`
}

func writeTasks(tasks []FetchTask) error {
	yml_, err := yaml.Marshal(tasks)
	if err != nil {
		fmt.Println("marshal yml error: \n", err)
		return err
	}
	err = os.WriteFile(FetchTaskFilePath, yml_, 0644)
	if err != nil {
		fmt.Println("write yml error: \n", err)
		return err
	}
	return nil
}

// HandleM3u8Url handle m3u8 url
func HandleM3u8Url(urlParam *url.URL) (m3u8FileForFfmpeg string, fetchTasks []FetchTask) {
	err := fetcher.Fetch(urlParam.String(), OriginalM3u8File)
	if err != nil {
		fmt.Println("fetch error: \n", err)
		return
	}
	originalM3u8Content_, err := os.ReadFile(OriginalM3u8File)
	if err != nil {
		fmt.Printf("read file %s error: %v \n", OriginalM3u8File, err)
		return
	}
	return handleM3u8File4Ffmpeg(urlParam, originalM3u8Content_)
}

func handleM3u8File4Ffmpeg(urlParam *url.URL, m3u8Content []byte) (
	m3u8FileForFfmpeg string, fetchTasks []FetchTask) {
	// handle m3u8 file to generate fetching task and prepare ffmpeg input m3u8
	m3u8_ := string(m3u8Content)
	scanner_ := bufio.NewScanner(strings.NewReader(m3u8_))
	m3u8Ffmpeg_ := ""
	for scanner_.Scan() {
		line_ := strings.TrimSpace(scanner_.Text())
		if strings.HasPrefix(line_, "#") {
			if match, err := regexp.Match(`^#\s*EXT-X-KEY`, []byte(line_)); match && err == nil {
				keyLine4Ffmpeg_, keyUri_ := handleM3u8KeyLine(line_)
				m3u8Ffmpeg_ += fmt.Sprintf("%s\n", keyLine4Ffmpeg_)
				if task := genFetchTask4Url(keyUri_, urlParam); task != nil {
					fetchTasks = append(fetchTasks, *task)
				}
			} else {
				m3u8Ffmpeg_ += fmt.Sprintf("%s\n", line_)
			}
		} else {
			if task := genFetchTask4Url(line_, urlParam); task != nil {
				m3u8Ffmpeg_ += fmt.Sprintf("%s\n", task.DstPath)
				fetchTasks = append(fetchTasks, *task)
			} else {
				m3u8Ffmpeg_ += fmt.Sprintf("%s\n", line_)
			}
		}
	}
	if err := scanner_.Err(); err != nil {
		fmt.Printf("error occurred: %v\n", err)
	}
	err := os.WriteFile(FfmpegInputM3u8File, []byte(m3u8Ffmpeg_), 0644)
	if err != nil {
		return "", nil
	}
	_ = writeTasks(fetchTasks)
	m3u8FileForFfmpeg = FfmpegInputM3u8File
	return
}

func genFetchTask4Url(fileUrlStr string, m3u8Url *url.URL) (fetchTask *FetchTask) {
	fileUrl_, err := url.Parse(fileUrlStr)
	if err != nil {
		fmt.Printf("raw parsing error: %v \n", err)
		return nil
	}
	if fileUrl_.Scheme == "" {
		fileUrl_.Scheme = m3u8Url.Scheme
	}
	if fileUrl_.Host == "" {
		fileUrl_.Host = m3u8Url.Host
	}
	if path.Dir(fileUrl_.Path) == "." {
		fileUrl_.Path = path.Join(path.Dir(m3u8Url.Path), path.Base(fileUrl_.Path))
	}
	if !utils.IsUrlSchemeHttp(fileUrl_) {
		return nil
	}
	return &FetchTask{
		UrlStr:  fileUrl_.String(),
		DstPath: path.Base(fileUrl_.Path),
	}
}

func handleM3u8KeyLine(keyLine string) (keyLine4Ffmpeg, keyUri string) {
	regex_ := regexp.MustCompile(
		`^(?P<BEFORE_URL>#\s*EXT-X-KEY\s*:.*URI\s*=\s*)(?P<URL>['"]?[^,'"]+['"]?)(?P<AFTER_URL>.*)$`,
	)
	matchResult_ := regex_.FindStringSubmatch(keyLine)
	namedGroups_ := make(map[string]string)
	for i, name := range regex_.SubexpNames() {
		if i != 0 && name != "" {
			namedGroups_[name] = matchResult_[i]
		}
	}
	if namedGroups_["URL"] == "" {
		return keyLine, ""
	} else {
		keyUri = strings.Trim(strings.Trim(namedGroups_["URL"], `"`), `'`)
		keyLine4Ffmpeg = fmt.Sprintf("%s%s%s", namedGroups_["BEFORE_URL"], `"ts.key"`, namedGroups_["AFTER_URL"])
	}
	return
}
