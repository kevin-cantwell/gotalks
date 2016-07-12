package gotalk_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGotalk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gotalk Suite")
}
