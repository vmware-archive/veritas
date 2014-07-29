package main_test

import (
	"github.com/cloudfoundry/gunk/timeprovider"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("Veritas", func() {
	var store *bbs.BBS
	BeforeEach(func() {
		store = bbs.New(etcdRunner.Adapter(), timeprovider.NewTimeProvider(), lagertest.NewTestLogger("veritas"))
	})

	It("should be able to fetch the contents of the store", func() {

	})
})
