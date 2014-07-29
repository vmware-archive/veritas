package main_test

import (
	"github.com/cloudfoundry/storeadapter/storerunner/etcdstorerunner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var veritas string
var etcdRunner *etcdstorerunner.ETCDClusterRunner

func TestVeritas(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Veritas Suite")
}

var _ = BeforeSuite(func() {
	var err error
	veritas, err = gexec.Build(".")
	Ω(err).ShouldNot(HaveOccurred())

	etcdRunner = etcdstorerunner.NewETCDClusterRunner(4001+GinkgoParallelNode(), 1)
})

var _ = BeforeEach(func() {
	etcdRunner.Start()
})

var _ = AfterEach(func() {
	etcdRunner.Stop()
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
	etcdRunner.Stop()
})
