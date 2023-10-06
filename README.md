# m3u8dl &middot; [![License](https://img.shields.io/hexpm/l/plug?logo=Github&style=flat)](https://github.com/tinkernels/m3u8dl/blob/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/tinkernels/m3u8dl)](https://goreportcard.com/report/github.com/tinkernels/m3u8dl)

m3u8dl is a simple tool to download m3u8 media files and merge them into mp4.

- concurrently download media files in m3u8.
- merge into mp4 with ffmpeg.

## Requirements

- ffmpeg executables in PATH (for merging media files).

## Install

```
go install github.com/tinkernels/m3u8dl@latest
```


## Usage 

```
Usage of m3u8dl:
  -c int
    	concurrency to fetch media segments (default 8)
  -url string
    	url of m3u8 file
  -w string
    	work dir
```