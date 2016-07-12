package gotalk

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"golang.org/x/net/context"
)

var _ = Describe("StartContainer", func() {
	It("Should start a docker container", func() {
		status, err := StartContainer(context.TODO(), "github.com/kevin-cantwell/gotalks", "localhost", "3999")
		Expect(err).To(BeNil())
		Expect(status.Repo).To(Equal("github.com/kevin-cantwell/gotalks"))
		Expect(status.HostPort).ToNot(Equal(""))
		Expect(status.Name).To(Equal("github.com_kevin-cantwell_gotalks"))
		StopContainer(status)
	})
})
