package types

import (
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

type VRFv2TestConfig interface {
	tc.CommonTestConfig
	ctf_config.GlobalTestConfig
	tc.VRFv2TestConfig
}

type VRFv2PlusTestConfig interface {
	tc.CommonTestConfig
	ctf_config.GlobalTestConfig
	tc.VRFv2PlusTestConfig
}

type FunctionsTestConfig interface {
	tc.CommonTestConfig
	ctf_config.GlobalTestConfig
	tc.FunctionsTestConfig
}

type AutomationTestConfig interface {
	ctf_config.GlobalTestConfig
	tc.CommonTestConfig
	tc.UpgradeableChainlinkTestConfig
	tc.AutomationTestConfig
}

type KeeperBenchmarkTestConfig interface {
	ctf_config.GlobalTestConfig
	tc.CommonTestConfig
	tc.KeeperTestConfig
	ctf_config.NamedConfiguration
	testreporters.GrafanaURLProvider
}

type OcrTestConfig interface {
	ctf_config.GlobalTestConfig
	tc.CommonTestConfig
	tc.OcrTestConfig
	ctf_config.SethConfig
}

type Ocr2TestConfig interface {
	ctf_config.GlobalTestConfig
	tc.CommonTestConfig
	tc.Ocr2TestConfig
}
