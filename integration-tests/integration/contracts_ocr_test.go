//go:build integration

package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/suite/testcommon"
)

var _ = Describe("OCR Feed @ocr", func() {

	DescribeTable("Deploys and watches an OCR feed @ocr", func(
		envInit environment.K8sEnvSpecInit,
	) {
		i := &testcommon.OCRSetupInputs{}
		testcommon.DeployOCRForEnv(i, envInit)
		testcommon.FundNodes(i)
		testcommon.DeployOCRContracts(i, 1)
		testcommon.SendOCRJobs(i)
		testcommon.CheckRound(i)
		By("Printing gas stats", func() {
			i.SuiteSetup.DefaultNetwork().Client.GasStats().PrintStats()
		})
		By("Tearing down the environment", i.SuiteSetup.TearDown())
	},
		Entry("all the same version", environment.NewChainlinkCluster(6)),
		Entry("different versions", environment.NewMixedVersionChainlinkCluster(6, 2)),
	)
})
