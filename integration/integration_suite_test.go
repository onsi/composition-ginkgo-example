package integration_test

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/onsi/composition-ginkgo-example/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	//tests to run
	"github.com/onsi/composition-ginkgo-example/tests/key_value_tests"
	"github.com/onsi/composition-ginkgo-example/tests/prefix_tests"

	"testing"
)

var context helpers.SharedContext
var keyValueStoreSession *gexec.Session

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	keyValueBinary, err := gexec.Build("github.com/onsi/composition-ginkgo-example/key_value_store")
	Ω(err).ShouldNot(HaveOccurred())

	cmd := exec.Command(keyValueBinary)
	keyValueStoreSession, err = gexec.Start(cmd, os.Stdout, os.Stdout)
	Ω(err).ShouldNot(HaveOccurred())

	address := "http://localhost:9999"

	//wait for the store to come up
	Eventually(helpers.KeyValueStorePinger(address)).Should(Succeed())

	return []byte(address)
}, func(address []byte) {
	context = helpers.NewSharedContext(
		string(address),
		fmt.Sprintf("prefix-%d", GinkgoParallelNode()),
	)

	key_value_tests.SharedContext = context
	prefix_tests.SharedContext = context
})

var _ = AfterEach(func() {
	//clean up data under this Ginkgo node's prefix
	Ω(context.Client.DeletePrefix(context.Prefix)).Should(Succeed())
})

var _ = SynchronizedAfterSuite(func() {
}, func() {
	keyValueStoreSession.Kill().Wait()
	gexec.CleanupBuildArtifacts()
})
