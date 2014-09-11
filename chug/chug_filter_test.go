package chug_test

import (
	"bytes"
	"regexp"
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
	var match, exclude *regexp.Regexp

	BeforeEach(func() {
		minTime = time.Time{}
		maxTime = time.Time{}
		match = nil
		exclude = nil

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
		out = ChugWithFilter(buffer, minTime, maxTime, match, exclude)
	})

	nextEntry := func() chug.Entry {
		var entry chug.Entry
		Eventually(out).Should(Receive(&entry))
		return entry
	}

	Context("with no time constraints or filters", func() {
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

	Context("with a match filter", func() {
		BeforeEach(func() {
			match = regexp.MustCompile(`lager-[14]`)
		})

		It("should return entries with `raw`s that match the filter", func() {
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-1"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-1"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-4"))
			Consistently(out).ShouldNot(Receive())
			Eventually(out).Should(BeClosed())
		})
	})

	Context("with an exclude filter", func() {
		BeforeEach(func() {
			exclude = regexp.MustCompile(`lager-[14]`)
		})

		It("should only return entries with `raw`s that do not match the exclude filter", func() {
			Ω(nextEntry().Log.Message).Should(Equal("lager-2"))
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-2"))
			Ω(nextEntry().Log.Message).Should(Equal("lager-3"))
			Ω(nextEntry().Raw).Should(ContainSubstring("none-lager-3"))
			Eventually(out).Should(BeClosed())
		})
	})

	Context("with both a match and exclude filter", func() {
		BeforeEach(func() {
			match = regexp.MustCompile(`lager-[13]`)
			exclude = regexp.MustCompile(`lager-1`)
		})

		It("should only return entries with `raw`s that match the match filter but not the exclude filter", func() {
			Ω(nextEntry().Log.Message).Should(Equal("lager-3"))
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
