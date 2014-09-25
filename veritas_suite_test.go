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

var _ = SynchronizedBeforeSuite(func() []byte {
	veritasPath, err := gexec.Build("github.com/pivotal-cf-experimental/veritas")
	Ω(err).ShouldNot(HaveOccurred())
	return []byte(veritasPath)
}, func(veritasPath []byte) {
	veritas = string(veritasPath)
	etcdRunner = etcdstorerunner.NewETCDClusterRunner(GinkgoParallelNode()+4001, 1)
})

var _ = BeforeEach(func() {
	etcdRunner.Start()
})

var _ = AfterEach(func() {
	etcdRunner.Stop()
})

var _ = SynchronizedAfterSuite(func() {
	if etcdRunner != nil {
		etcdRunner.Stop()
	}
}, func() {
	gexec.CleanupBuildArtifacts()
})
