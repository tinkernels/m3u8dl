package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func ExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func GetCWD() string {
	// using the function
	cwd_, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd_
}

func PathExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func AbsPath(path string) (absPath string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("path error: ", err)
		return path
	}
	return
}
