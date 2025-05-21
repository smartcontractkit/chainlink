package verification

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func TestSetAccessController(t *testing.T) {
	t.Parallel()
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		ShouldDeployMCMS: true,
	})

	e := testEnv.Environment
	acAddress := common.HexToAddress("0x0000000000000000000000000000000000000123")
	testChain := e.AllChainSelectors()[0]

	e, verifierProxyAddr, _ := DeployVerifierProxyAndVerifier(t, e)

	cfg := VerifierProxySetAccessControllerConfig{
		ConfigPerChain: map[uint64][]SetAccessControllerConfig{
			testChain: {
				{AccessControllerAddress: acAddress, ContractAddress: verifierProxyAddr},
			},
		},
	}

	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			SetAccessControllerChangeset,
			cfg,
		),
	)
	require.NoError(t, err)

	client := e.Chains[testChain].Client
	verifierProxy, err := verifier_proxy_v0_5_0.NewVerifierProxy(verifierProxyAddr, client)
	require.NoError(t, err)

	accessController, err := verifierProxy.SAccessController(nil)
	require.NoError(t, err)
	require.Equal(t, acAddress, accessController)
}
