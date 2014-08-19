package main_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/cloudfoundry/gunk/timeprovider"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("Veritas", func() {
	var (
		store  *bbs.BBS
		tmpDir string
		err    error
	)

	BeforeEach(func() {
		store = bbs.NewBBS(etcdRunner.Adapter(), timeprovider.NewTimeProvider(), lagertest.NewTestLogger("veritas"))

		err = store.DesireTask(models.Task{
			Guid:   "Task-Guid",
			Stack:  "pancakes",
			Domain: "veritas",
			Actions: []models.ExecutorAction{
				{models.RunAction{Path: "foo"}},
			},
		})
		Ω(err).ShouldNot(HaveOccurred())

		err = store.DesireLRP(models.DesiredLRP{
			ProcessGuid: "Desired-Process-Guid",
			Stack:       "pancakes",
			Domain:      "veritas",
			Actions: []models.ExecutorAction{
				{models.RunAction{Path: "foo"}},
			},
		})
		Ω(err).ShouldNot(HaveOccurred())

		err = store.ReportActualLRPAsRunning(models.ActualLRP{
			ProcessGuid:  "Actual-Process-Guid",
			InstanceGuid: "Instance-Guid",
			Index:        0,
		}, "Executor-ID")
		Ω(err).ShouldNot(HaveOccurred())

		err = store.ReportActualLRPAsRunning(models.ActualLRP{
			ProcessGuid:  "Actual-Process-Guid",
			InstanceGuid: "Instance-Guid-200",
			Index:        200,
		}, "Executor-ID")
		Ω(err).ShouldNot(HaveOccurred())

		err = store.RequestLRPStartAuction(models.LRPStartAuction{
			InstanceGuid: "InstanceGuid",
			DesiredLRP: models.DesiredLRP{
				ProcessGuid: "StartAuction-Process-Guid",
			},
			Index: 1,
		})
		Ω(err).ShouldNot(HaveOccurred())

		err = store.RequestLRPStopAuction(models.LRPStopAuction{
			ProcessGuid: "StopAuction-Process-Guid",
			Index:       2,
		})
		Ω(err).ShouldNot(HaveOccurred())

		err = store.RequestStopLRPInstance(models.StopLRPInstance{
			ProcessGuid:  "StopLRP-Process-Guid",
			Index:        3,
			InstanceGuid: "Instance-Guid",
		})
		Ω(err).ShouldNot(HaveOccurred())

		tmpDir, err = ioutil.TempDir("", "veritas")
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})

	It("should be able to print out the contents of the bbs", func() {
		dumpFile := filepath.Join(tmpDir, "dump")

		session, err := gexec.Start(exec.Command(veritas, "fetch-store", "-etcdCluster="+strings.Join(etcdRunner.NodeURLS(), ","), dumpFile), GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		dump, err := ioutil.ReadFile(dumpFile)
		Ω(err).ShouldNot(HaveOccurred())

		Ω(dump).Should(ContainSubstring("Desired-Process-Guid"))
		Ω(dump).Should(ContainSubstring("Actual-Process-Guid"))
		Ω(dump).Should(ContainSubstring("StartAuction-Process-Guid"))
		Ω(dump).Should(ContainSubstring("StopAuction-Process-Guid"))
		Ω(dump).Should(ContainSubstring("StopLRP-Process-Guid"))

		session, err = gexec.Start(exec.Command(veritas, "print-store", dumpFile), GinkgoWriter, GinkgoWriter)
		Eventually(session).Should(gexec.Exit(0))

		Ω(session.Out.Contents()).Should(ContainSubstring("Desired-Process-Guid"))
		Ω(session.Out.Contents()).Should(ContainSubstring("Actual-Process-Guid"))
		Ω(session.Out.Contents()).Should(ContainSubstring("StartAuction-Process-Guid"))
		Ω(session.Out.Contents()).Should(ContainSubstring("StopAuction-Process-Guid"))
		Ω(session.Out.Contents()).Should(ContainSubstring("StopLRP-Process-Guid"))
	})
})
