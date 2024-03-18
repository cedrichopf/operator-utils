package network_test

import (
	"github.com/cedrichopf/operator-utils/network"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test HostIsResolvable", func() {
	Context("when host is resolvable", func() {
		It("returns true", func() {
			host := "localhost"
			result := network.HostIsResolvable(host)
			Expect(result).To(Equal(true))
		})
	})

	Context("when host is not resolvable", func() {
		It("returns false", func() {
			host := "some-non-existing-host.default.tmp.domain"
			result := network.HostIsResolvable(host)
			Expect(result).To(Equal(false))
		})
	})
})
