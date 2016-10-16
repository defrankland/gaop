package goaop_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoAop(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoAop Suite")
}
