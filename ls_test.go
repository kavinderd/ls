package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// TODO: Read about go imports
	"github.com/onsi/gomega/gexec"
	"os/exec"
)

var _ = Describe("ls", func() {

	var err error
	var pathToLs string

	BeforeSuite(func() {
		pathToLs, err = gexec.Build("ls")
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	It("Only displays the non-hidden files without any flags", func() {
		command := exec.Command(pathToLs)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		commandReal := exec.Command("ls")
		sessionReal, err := gexec.Start(commandReal, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Eventually(sessionReal).Should(gexec.Exit(0))
		Ω(session.Out.Contents()).Should(Equal(sessionReal.Out.Contents()))
	})
})
