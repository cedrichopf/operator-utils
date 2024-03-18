package network_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cedrichopf/operator-utils/network"
)

var _ = Describe("Test HostFromURL", func() {
	Context("when url is valid", func() {
		It("returns the host", func() {
			url := "http://localhost:9000/some/api/endpoint"
			result, err := network.HostFromURL(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("localhost"))
		})
	})

	Context("when url is invalid", func() {
		It("returns an error", func() {
			url := "localhost:9000/some/api/endpoint"
			result, err := network.HostFromURL(url)
			Expect(err).To(HaveOccurred())
			Expect(result).To(Equal(""))
		})
	})
})
