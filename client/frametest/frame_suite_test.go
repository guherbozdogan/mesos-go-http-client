package frametest

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFrame(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Frame Suite")
}