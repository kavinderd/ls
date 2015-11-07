package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	err      error
	path     string
	fileList []os.FileInfo
)

func SetPath() {
	path, err = os.Getwd()
	if err != nil {
		os.Exit(1)
	}
}

func CollectFiles() {
	fileList, err = ioutil.ReadDir(path)
	if err != nil {
		os.Exit(2)
	}
}

func PrintFiles() {
	for _, file := range fileList {
		if strings.Index(file.Name(), ".") != 0 {
			fmt.Println(file.Name())
		}
	}
}

func main() {
	SetPath()
	CollectFiles()
	PrintFiles()
}
