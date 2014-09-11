package chug_test

import (
	"bytes"
	"time"

	. "github.com/cloudfoundry-incubator/veritas/chug"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/chug"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ChugFilter", func() {
	var buffer *bytes.Buffer
	var t0, t1, t2, t3 time.Time
	var out chan chug.Entry
	var minTime, maxTime time.Time

	BeforeEach(func() {
		minTime = time.Time{}
		maxTime = time.Time{}

		logger := lager.NewLogger("logger")
		buffer = &bytes.Buffer{}
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.DEBUG))

		t0 = time.Now()
		time.Sleep(10 * time.Millisecond)

		buffer.WriteString("none-lager-1\n")
		logger.Info("lager-1")

		time.Sleep(10 * time.Millisecond)
		t1 = time.Now()
		time.Sleep(10 * time.Millisecond)

		logger.Info("lager-2")
		buffer.WriteString("none-lager-2\n")
		logger.Info("lager-3")

		time.Sleep(10 * time.Millisecond)
		t2 = time.Now()
		time.Sleep(10 * time.Millisecond)

		logger.Info("lager-4")

		time.Sleep(10 * time.Millisecond)
		t3 = time.Now()
		time.Sleep(10 * time.Millisecond)

		buffer.WriteString("none-lager-3\n")
	})

	JustBeforeEach(func() {
		out = ChugWithFilter(buffer, minTime, maxTime)
	})

	nextEntry := func() chug.Entry {
		var entry chug.Entry
		Eventually(out).Should(Receive(&entry))
		return entry
	}

	Context("with no time constraints", func() {
		BeforeEach(func() {
			minTime = time.Time{}
			maxTime = time.Time{}
		})

		It("should return all the entries", func() {
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-1"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-1"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-2"))
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-2"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-3"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-4"))
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-3"))
			Eventually(out).Should(BeClosed())
		})
	})

	Context("with a minimum time constraint that precedes a non-lager line", func() {
		BeforeEach(func() {
			minTime = t0
		})

		It("should only return entries after the minimum time, ignoring leading non-lager lines", func() {
			Ω(nextEntry().Log.Message).Should(Equal("lager-1"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-2"))
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-2"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-3"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-4"))
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-3"))
			Eventually(out).Should(BeClosed())
		})
	})

	Context("with a minimum time constraint that's after a lager line", func() {
		BeforeEach(func() {
			minTime = t1
		})

		It("should only return entries after the minimum time", func() {
			Ω(nextEntry().Log.Message).Should(Equal("lager-2"))
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-2"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-3"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-4"))
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-3"))
			Eventually(out).Should(BeClosed())
		})
	})

	Context("with a maximum time constraint that's before a lager line", func() {
		BeforeEach(func() {
			maxTime = t2
		})

		It("should only return entries before the maximum time", func() {
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-1"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-1"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-2"))
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-2"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-3"))
			Consistently(out).ShouldNot(Receive())
			Eventually(out).Should(BeClosed())
		})
	})

	Context("with a maximum time constraint that's before a none-lager line", func() {
		BeforeEach(func() {
			maxTime = t3
		})

		It("should return the lager lines", func() {
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-1"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-1"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-2"))
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-2"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-3"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-4"))
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-3"))
			Eventually(out).Should(BeClosed())
		})
	})
})
