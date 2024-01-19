package types

import (
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

type VRFv2TestConfig interface {
	tc.CommonTestConfig
	tc.GlobalTestConfig
	tc.VRFv2TestConfig
}

type VRFv2PlusTestConfig interface {
	tc.CommonTestConfig
	tc.GlobalTestConfig
	tc.VRFv2PlusTestConfig
}

type FunctionsTestConfig interface {
	tc.CommonTestConfig
	tc.GlobalTestConfig
	tc.FunctionsTestConfig
}

type AutomationTestConfig interface {
	tc.GlobalTestConfig
	tc.CommonTestConfig
	tc.UpgradeableChainlinkTestConfig
}

type KeeperBenchmarkTestConfig interface {
	tc.GlobalTestConfig
	tc.CommonTestConfig
	tc.KeeperTestConfig
	tc.NamedConfiguration
	testreporters.GrafanaURLProvider
}

type OcrTestConfig interface {
	tc.GlobalTestConfig
	tc.CommonTestConfig
	tc.OcrTestConfig
}
