package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

var (
	err               error
	path              string
	fileList          []os.FileInfo
	fileModes         = make([]string, 0)
	fileUsers         = make([]*user.User, 0)
	fileUsernames     = make([]string, 0)
	fileGroups        = make([]string, 0)
	fileSizes         = make([]string, 0)
	fileModTimes      = make([]string, 0)
	detailsMode       *bool
	maxUsernameLength = 0
	maxSizeLength     = 0
	maxFilenameLength = 0
	maxGroupLength    = 0
	totalFileBlocks   = 0
)

func SetPath() {
	path, err = os.Getwd()
	if err != nil {
		os.Exit(1)
	}
}

func CollectData() {
	fileList, err = ioutil.ReadDir(path)
	if err != nil {
		os.Exit(2)
	}

	if *detailsMode {
		getFilePermissionsDone := make(chan bool)
		getFileUsersDone := make(chan bool)
		getFileUsernamesDone := make(chan bool)
		getFileGroupsDone := make(chan bool)
		getFileSizesDone := make(chan bool)
		getFileModTimesDone := make(chan bool)
		getTotalFileBlocks := make(chan bool)

		go GetFilePermissions(getFilePermissionsDone)
		go GetFileUsers(getFileUsersDone)
		go GetTotalFileBlocks(getTotalFileBlocks)
		<-getFilePermissionsDone
		<-getFileUsersDone
		<-getTotalFileBlocks

		go GetFileUsernames(getFileUsernamesDone)
		go GetFileGroups(getFileGroupsDone)
		go GetFileSizes(getFileSizesDone)
		go GetFileModTimes(getFileModTimesDone)
		<-getFileUsernamesDone
		<-getFileGroupsDone
		<-getFileSizesDone
		<-getFileModTimesDone
	}
}

func GetTotalFileBlocks(done chan bool) {
	for _, file := range fileList {
		cmd := exec.Command("stat", "-f", "\"%b\"", file.Name())
		output, err := cmd.Output()
		if err != nil {
			return
		}
		num := strings.TrimSpace(string(output))
		num = strings.Replace(num, "\"", "", -1)
		blockSize, _ := strconv.Atoi(num)
		totalFileBlocks += blockSize
	}
	done <- true
}

func GetFileSizes(done chan bool) {
	for _, file := range fileList {
		size := fmt.Sprintf("%d", file.Size())
		if len(size) > maxSizeLength {
			maxSizeLength = len(size)
		}
		fileSizes = append(fileSizes, size)
	}
	done <- true
}

func GetFilePermissions(done chan bool) {
	for _, file := range fileList {
		fileModes = append(fileModes, file.Mode().String())
	}
	done <- true
}

func GetFileUsers(done chan bool) {
	for _, file := range fileList {
		uid := fmt.Sprintf("%d", file.Sys().(*syscall.Stat_t).Uid)
		user, err := user.LookupId(uid)
		if err != nil {
			os.Exit(1)
		}
		fileUsers = append(fileUsers, user)
	}
	done <- true
}

func GetFileModTimes(done chan bool) {
	for _, file := range fileList {
		fileModTimes = append(fileModTimes, file.ModTime().Format("Jan _2 15:04"))
	}
	done <- true
}
func GetFileGroups(done chan bool) {
	for _, user := range fileUsers {
		command := exec.Command("dscl", ".", "-search", "/Groups", "gid", user.Gid)
		output, err := command.Output()
		if err != nil {
			os.Exit(1)
		}
		outputParts := strings.Fields(string(output))
		outputLine := strings.Fields(string(outputParts[0]))
		group := outputLine[0]
		if len(group) > maxGroupLength {
			maxGroupLength = len(group)
		}
		fileGroups = append(fileGroups, group)
	}
	done <- true
}

func GetFileUsernames(done chan bool) {
	for _, user := range fileUsers {
		if len(user.Username) > maxUsernameLength {
			maxUsernameLength = len(user.Username)
		}
		fileUsernames = append(fileUsernames, user.Username)
	}
	done <- true
}

func PrintFiles() {
	if *detailsMode {
		fmt.Printf("total %d\n", totalFileBlocks)
	}
	for i, file := range fileList {
		if strings.Index(file.Name(), ".") != 0 {
			if *detailsMode {
				usernameSize := fmt.Sprintf("%d", maxUsernameLength+1)
				maxFileSize := fmt.Sprintf("%d", maxSizeLength+1)
				fmt.Printf("%-11s %d %-"+usernameSize+"s %-s %"+maxFileSize+"s %12s %s\n", fileModes[i], 1, fileUsernames[i], fileGroups[i], fileSizes[i], fileModTimes[i], file.Name())
			} else {
				fmt.Println(file.Name())
			}
		}
	}
}

func main() {
	detailsMode = flag.Bool("l", false, "Display details")
	flag.Parse()

	SetPath()
	CollectData()
	PrintFiles()
}
