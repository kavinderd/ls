package main_test

import (
	"bufio"
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	. "ls"
	"os"
	"os/exec"
)

var _ = Describe("ls", func() {

	var stdout *os.File
	var readFile *os.File
	var writeFile *os.File
	var err error
	var buf bytes.Buffer
	var routine func(reader io.ReadCloser, channel chan string)

	BeforeEach(func() {
		stdout = os.Stdout
		readFile, writeFile, err = os.Pipe()
		if err != nil {
			Fail("Couldn't create Pipe")
		}

		routine = func(reader io.ReadCloser, channel chan string) {
			var b bytes.Buffer
			_, err := io.Copy(&b, reader)
			reader.Close()
			if err != nil {
				Fail("Error in Channel")
			}
			channel <- b.String()
		}
		os.Stdout = writeFile
	})

	AfterEach(func() {
		buf.Reset()
	})

	cleanup := func() {
		writeFile.Close()
		os.Stdout = stdout
	}

	var _ = Describe("without any args", func() {
		It("lists the contents of the current working dir", func() {
			path := ""

			writer := bufio.NewWriter(os.Stdout)
			SimpleLs(path, writer)

			outC := make(chan string)
			go routine(readFile, outC)

			ls := exec.Command("ls")
			b, err := ls.Output()
			if err != nil {
				Fail("Error line 70")
			}

			buf.Write(b)
			cleanup()

			out := <-outC

			fmt.Print(out)
			Expect(out).To(Equal(buf.String()))
		})
	})
})
