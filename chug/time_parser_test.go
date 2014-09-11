package chug_test

import (
	"time"

	. "github.com/cloudfoundry-incubator/veritas/chug"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TimeParser", func() {
	Context("when passed an empty string", func() {
		It("should return zero", func() {
			Ω(ParseTimeFlag("")).Should(BeZero())
		})
	})

	Context("when passed a unix timestamp", func() {
		It("should return that time", func() {
			Ω(ParseTimeFlag("1123491.1238")).Should(Equal(time.Unix(1123491, 123800000)))
		})
	})

	Context("when passed a duration", func() {
		It("should return a time relative to the current time", func() {
			duration := -(time.Hour + 15*time.Minute + 3*time.Second)
			Ω(ParseTimeFlag("-1h15m3s")).Should(BeTemporally("~", time.Now().Add(duration), time.Second))
		})
	})

	Context("when passed a chug-formatted timestamp", func() {
		It("should return that time", func() {
			expectedTime, err := time.Parse("01/02 15:04:05.00", "09/08 22:45:06.79")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(ParseTimeFlag("09/08 22:45:06.79")).Should(Equal(expectedTime))
		})
	})

	Context("when passed anything else", func() {
		It("should error", func() {
			t, err := ParseTimeFlag("abc")
			Ω(err).Should(HaveOccurred())
			Ω(t).Should(BeZero())
		})
	})
})
