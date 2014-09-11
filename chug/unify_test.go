package chug_test

import (
	"bytes"
	"io"
	"regexp"
	"time"

	. "github.com/cloudfoundry-incubator/veritas/chug"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Unify", func() {
	var loggerA, loggerB lager.Logger
	var bufferA, bufferB *bytes.Buffer
	var t0, t1 time.Time

	BeforeEach(func() {
		loggerA = lager.NewLogger("A")
		loggerB = lager.NewLogger("B")

		bufferA = &bytes.Buffer{}
		bufferB = &bytes.Buffer{}

		loggerA.RegisterSink(lager.NewWriterSink(bufferA, lager.DEBUG))
		loggerB.RegisterSink(lager.NewWriterSink(bufferB, lager.DEBUG))

		bufferA.Write([]byte("non-lager-A\n"))
		bufferB.Write([]byte("non-lager-B-1\n"))
		loggerA.Info("A-1")

		time.Sleep(10 * time.Millisecond)
		t0 = time.Now()
		time.Sleep(10 * time.Millisecond)

		loggerB.Info("B-1")
		loggerA.Info("A-2")
		loggerB.Info("B-2")
		loggerB.Info("B-3")

		time.Sleep(10 * time.Millisecond)
		t1 = time.Now()
		time.Sleep(10 * time.Millisecond)

		loggerA.Info("A-3")
		bufferB.Write([]byte("non-lager-B-2\n"))
		loggerB.Info("B-4")
	})

	It("should unify the independent streams", func(done Done) {
		out := gbytes.NewBuffer()
		err := Unify([]io.Reader{bufferA, bufferB}, out, time.Time{}, time.Time{}, nil, nil)
		Ω(err).ShouldNot(HaveOccurred())

		Ω(out).Should(gbytes.Say("non-lager-A"))
		Ω(out).Should(gbytes.Say("non-lager-B-1"))
		Ω(out).Should(gbytes.Say("A-1"))
		Ω(out).Should(gbytes.Say("B-1"))
		Ω(out).Should(gbytes.Say("A-2"))
		Ω(out).Should(gbytes.Say("B-2"))
		Ω(out).Should(gbytes.Say("B-3"))
		Ω(out).Should(gbytes.Say("non-lager-B-2"), "non-lager lines always beat lager lines")
		Ω(out).Should(gbytes.Say("A-3"))
		Ω(out).Should(gbytes.Say("B-4"))

		close(done)
	})

	Context("with a minimum time", func() {
		It("should only show log lines after that time", func(done Done) {
			out := gbytes.NewBuffer()
			err := Unify([]io.Reader{bufferA, bufferB}, out, t0, time.Time{}, nil, nil)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(out).ShouldNot(gbytes.Say("non-lager-A"))
			Ω(out).ShouldNot(gbytes.Say("non-lager-B-1"))
			Ω(out).ShouldNot(gbytes.Say("A-1"))
			Ω(out).Should(gbytes.Say("B-1"))
			Ω(out).Should(gbytes.Say("A-2"))
			Ω(out).Should(gbytes.Say("B-2"))
			Ω(out).Should(gbytes.Say("B-3"))
			Ω(out).Should(gbytes.Say("non-lager-B-2"), "non-lager lines always beat lager lines")
			Ω(out).Should(gbytes.Say("A-3"))
			Ω(out).Should(gbytes.Say("B-4"))

			close(done)
		})
	})

	Context("with a maximum time", func() {
		It("should only show log lines before that time", func(done Done) {
			out := gbytes.NewBuffer()
			err := Unify([]io.Reader{bufferA, bufferB}, out, time.Time{}, t1, nil, nil)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(out).Should(gbytes.Say("non-lager-A"))
			Ω(out).Should(gbytes.Say("non-lager-B-1"))
			Ω(out).Should(gbytes.Say("A-1"))
			Ω(out).Should(gbytes.Say("B-1"))
			Ω(out).Should(gbytes.Say("A-2"))
			Ω(out).Should(gbytes.Say("B-2"))
			Ω(out).Should(gbytes.Say("B-3"))
			Ω(out).Should(gbytes.Say("non-lager-B-2"), "non-lager lines always beat lager lines")
			Ω(out).ShouldNot(gbytes.Say("A-3"))
			Ω(out).ShouldNot(gbytes.Say("B-4"))

			close(done)
		})
	})

	Context("with both a minimum and maximum time", func() {
		It("should apply both limits", func(done Done) {
			out := gbytes.NewBuffer()
			err := Unify([]io.Reader{bufferA, bufferB}, out, t0, t1, nil, nil)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(out).ShouldNot(gbytes.Say("non-lager-A"))
			Ω(out).ShouldNot(gbytes.Say("non-lager-B-1"))
			Ω(out).ShouldNot(gbytes.Say("A-1"))
			Ω(out).Should(gbytes.Say("B-1"))
			Ω(out).Should(gbytes.Say("A-2"))
			Ω(out).Should(gbytes.Say("B-2"))
			Ω(out).Should(gbytes.Say("B-3"))
			Ω(out).Should(gbytes.Say("non-lager-B-2"), "non-lager lines always beat lager lines")
			Ω(out).ShouldNot(gbytes.Say("A-3"))
			Ω(out).ShouldNot(gbytes.Say("B-4"))

			close(done)
		})
	})

	Context("with a time limits and filters", func() {
		It("should apply the time limits and the filters", func(done Done) {
			out := gbytes.NewBuffer()
			err := Unify([]io.Reader{bufferA, bufferB}, out, t0, t1, regexp.MustCompile("B"), regexp.MustCompile("3"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(out).ShouldNot(gbytes.Say("non-lager-A"))
			Ω(out).ShouldNot(gbytes.Say("non-lager-B-1"))
			Ω(out).ShouldNot(gbytes.Say("A-1"))
			Ω(out).Should(gbytes.Say("B-1"))
			Ω(out).ShouldNot(gbytes.Say("A-2"))
			Ω(out).Should(gbytes.Say("B-2"))
			Ω(out).ShouldNot(gbytes.Say("B-3"))
			Ω(out).Should(gbytes.Say("non-lager-B-2"), "non-lager lines always beat lager lines")
			Ω(out).ShouldNot(gbytes.Say("A-3"))
			Ω(out).ShouldNot(gbytes.Say("B-4"))

			close(done)
		})
	})
})
