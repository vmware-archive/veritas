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
	etcdRunner = etcdstorerunner.NewETCDClusterRunner(GinkgoParallelNode()+4001, 1)
	veritas, err = gexec.Build("github.com/cloudfoundry-incubator/veritas")
	Ω(err).ShouldNot(HaveOccurred())
})

var _ = BeforeEach(func() {
	etcdRunner.Start()
})

var _ = AfterEach(func() {
	etcdRunner.Stop()
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
