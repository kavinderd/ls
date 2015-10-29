package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ls Suite")
}
