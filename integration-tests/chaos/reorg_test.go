package chaos

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/suite/testcommon"
)

var _ = Describe("Reorg example test @reorg", func() {
	i := &testcommon.RunlogSetupInputs{}
	It("Performs reorg and verifies it", func() {
		testcommon.SetupRunlogEnv(i)

		reorgConfirmer, err := NewReorgConfirmer(
			i.SuiteSetup.DefaultNetwork().Client,
			i.SuiteSetup.Environment(),
			5,
			10,
			time.Second*600,
		)
		Expect(err).ShouldNot(HaveOccurred())
		i.SuiteSetup.DefaultNetwork().Client.AddHeaderEventSubscription("reorg", reorgConfirmer)
		err = i.SuiteSetup.DefaultNetwork().Client.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())
		err = reorgConfirmer.Verify()
		Expect(err).ShouldNot(HaveOccurred())
	})
	AfterEach(func() {
		By("Restoring chaos", func() {
			err := i.SuiteSetup.Environment().StopAllChaos()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Tearing down the environment", i.SuiteSetup.TearDown())
	})
})
