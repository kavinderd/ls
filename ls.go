package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var err error

func SimpleLs(path string, writer *bufio.Writer) {
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			writer.WriteString("Error getting current wd, not my problem")
			return
		}
	}

	contents, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, file := range contents {
		name := file.Name()
		if strings.Index(name, ".") != 0 {
			writer.WriteString(file.Name())
			writer.WriteByte(10)
		}
	}
	writer.Flush()
}

func main() {
	path, _ := os.Getwd()
	contents, _ := ioutil.ReadDir(path)
	for _, file := range contents {
		fmt.Println(file.Name())
	}
}
