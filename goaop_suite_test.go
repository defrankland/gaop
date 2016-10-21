package gaop_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGaop(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gaop Suite")
}
