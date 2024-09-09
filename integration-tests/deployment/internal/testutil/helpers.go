package testutil

import (
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	chainlink_test_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func GetDefaultNewChainlinkClusterConfig() *testconfig.ChainlinkDeployment {
	return &testconfig.ChainlinkDeployment{
		Common: &testconfig.Node{
			ChainlinkImage: &ctfconfig.ChainlinkImageConfig{
				Image:   ptr.Ptr("public.ecr.aws/chainlink/chainlink"),
				Version: ptr.Ptr("2.13.0"),
			},
			DBImage: "795953128386.dkr.ecr.us-west-2.amazonaws.com/postgres",
			DBTag:   "15.6",
			BaseConfigTOML: `
[Feature]
LogPoller = true

[Log]
Level = 'debug'
JSONConsole = true

[Log.File]
MaxSize = '0b'

[WebServer]
AllowOrigins = '*'
HTTPPort = 6688
SecureCookies = false
HTTPWriteTimeout = '1m'

[WebServer.RateLimit]
Authenticated = 2000
Unauthenticated = 1000

[WebServer.TLS]
HTTPSPort = 0

[Database]
MaxIdleConns = 10
MaxOpenConns = 20
MigrateOnStartup = true

[OCR2]
Enabled = true
DefaultTransactionQueueDepth = 0

[OCR]
Enabled = false
DefaultTransactionQueueDepth = 0

[P2P]
[P2P.V2]
Enabled = true
ListenAddresses = ['0.0.0.0:6690']
AnnounceAddresses = ['0.0.0.0:6690']
DeltaDial = '500ms'
DeltaReconcile = '5s'
`},
		NoOfNodes: ptr.Ptr(5),
	}
}

func GetDefaultNewClusterConfigFromChainlinkTestConfig(config chainlink_test_config.TestConfig) *testconfig.ChainlinkDeployment {
	return &testconfig.ChainlinkDeployment{
		Common: &testconfig.Node{
			ChainlinkImage: &ctfconfig.ChainlinkImageConfig{
				Image:   config.ChainlinkImage.Image,
				Version: config.ChainlinkImage.Version,
			},
			DBImage:        "795953128386.dkr.ecr.us-west-2.amazonaws.com/postgres",
			DBTag:          "15.6",
			BaseConfigTOML: config.NodeConfig.BaseConfigTOML,
		},
		NoOfNodes: ptr.Ptr(5),
	}
}
