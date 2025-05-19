package testconfig

import (
	"encoding/base64"
	"math/big"
	"os"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	a_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/automation"
)

func TestBase64ConfigRead(t *testing.T) {
	os.Setenv("E2E_TEST_ARBITRUM_GOERLI_RPC_HTTP_URL", "https://devnet-1.mt/ABC/rpc/")
	os.Setenv("E2E_TEST_ARBITRUM_GOERLI_RPC_WS_URL", "wss://devnet-1.mt/ABC/rpc/")
	os.Setenv("E2E_TEST_ARBITRUM_GOERLI_WALLET_KEY", "0x3333333333333333333333333333333333333333")
	defer os.Unsetenv("E2E_TEST_ARBITRUM_GOERLI_RPC_HTTP_URL")
	defer os.Unsetenv("E2E_TEST_ARBITRUM_GOERLI_RPC_WS_URL")
	defer os.Unsetenv("E2E_TEST_ARBITRUM_GOERLI_WALLET_KEY")
	os.Setenv("E2E_TEST_OPTIMISM_GOERLI_RPC_HTTP_URL", "https://devnet-3.mt/ABC/rpc/")
	os.Setenv("E2E_TEST_OPTIMISM_GOERLI_RPC_WS_URL", "wss://devnet-3.mt/ABC/rpc/")
	os.Setenv("E2E_TEST_OPTIMISM_GOERLI_WALLET_KEY", "0x3333333333333333333333333333333333333333")
	defer os.Unsetenv("E2E_TEST_OPTIMISM_GOERLI_RPC_HTTP_URL")
	defer os.Unsetenv("E2E_TEST_OPTIMISM_GOERLI_RPC_WS_URL")
	defer os.Unsetenv("E2E_TEST_OPTIMISM_GOERLI_WALLET_KEY")

	testConfig := TestConfig{
		Automation: &a_config.Config{
			General: &a_config.General{
				NumberOfNodes:         ptr.Ptr(7),
				Duration:              ptr.Ptr(9),
				BlockTime:             ptr.Ptr(10),
				SpecType:              ptr.Ptr("minimum"),
				ChainlinkNodeLogLevel: ptr.Ptr("debug"),
				UsePrometheus:         ptr.Ptr(true),
				RemoveNamespace:       ptr.Ptr(true),
			},
			Load: []a_config.Load{
				{
					NumberOfUpkeeps:               ptr.Ptr(1),
					NumberOfEvents:                ptr.Ptr(2),
					NumberOfSpamMatchingEvents:    ptr.Ptr(3),
					NumberOfSpamNonMatchingEvents: ptr.Ptr(4),
					CheckBurnAmount:               big.NewInt(5),
					PerformBurnAmount:             big.NewInt(6),
					SharedTrigger:                 ptr.Ptr(true),
				},
				{
					NumberOfUpkeeps:               ptr.Ptr(3),
					NumberOfEvents:                ptr.Ptr(2),
					NumberOfSpamMatchingEvents:    ptr.Ptr(3),
					NumberOfSpamNonMatchingEvents: ptr.Ptr(7),
					CheckBurnAmount:               big.NewInt(5),
					PerformBurnAmount:             big.NewInt(6),
					SharedTrigger:                 ptr.Ptr(false),
				},
			},
		},
		TestConfig: ctf_config.TestConfig{
			Network: &ctf_config.NetworkConfig{
				SelectedNetworks: []string{"OPTIMISM_GOERLI"},
				RpcHttpUrls: map[string][]string{
					"OPTIMISM_GOERLI": {"http://localhost:8545"},
				},
				WalletKeys: map[string][]string{
					"OPTIMISM_GOERLI": {"0x3333333333333333333333333333333333333333"},
				},
			},
		},
	}

	configMarshalled, err := toml.Marshal(testConfig)
	require.NoError(t, err, "Error marshalling test config")

	testConfigEncoded := base64.StdEncoding.EncodeToString(configMarshalled)
	os.Setenv(Base64OverrideEnvVarName, testConfigEncoded)

	readConfig, err := GetConfig([]string{"test"}, Automation)
	require.NoError(t, err, "Error reading config")

	require.NotNil(t, readConfig.Automation, "Automation config read from base64 is nil")
	require.Equal(t, testConfig.Automation.General, readConfig.Automation.General, "General automation config does not match expected")
	require.EqualValues(t, testConfig.Automation.Load, readConfig.Automation.Load, "Load automation config does not match expected")
	require.NotNil(t, readConfig.Network, "Network config read from base64 is nil")
	require.Equal(t, testConfig.Network.SelectedNetworks, readConfig.Network.SelectedNetworks, "SelectedNetwork config entry read from base64 does not match expected")
	require.Equal(t, []string{"http://localhost:8545"}, readConfig.Network.RpcHttpUrls["OPTIMISM_GOERLI"], "RpcHttpUrls config entry read from base64 does not match expected")
	require.Equal(t, []string{"wss://devnet-2.mt/ABC/rpc/"}, readConfig.Network.RpcWsUrls["OPTIMISM_GOERLI"], "RpcWsUrls config entry read from base64 network defaults does not match expected")
	require.Equal(t, testConfig.Network.WalletKeys, readConfig.Network.WalletKeys, "WalletKeys config entry read from base64 does not match expected")
}
