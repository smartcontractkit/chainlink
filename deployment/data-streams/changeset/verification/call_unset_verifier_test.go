package verification

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func TestUnsetVerifier(t *testing.T) {
	t.Parallel()
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		ShouldDeployMCMS: true,
	})
	e := testEnv.Environment
	e, verifierProxyAddr, verifierAddr := DeployVerifierProxyAndVerifier(t, e)

	chainSelector := e.AllChainSelectors()[0]
	chain := e.Chains[chainSelector]

	verifier, err := verifier_v0_5_0.NewVerifier(verifierAddr, e.Chains[chainSelector].Client)
	require.NoError(t, err)
	require.NotNil(t, verifier)

	// SetConfig on the verifier
	var configDigest [32]byte
	cdBytes, _ := hex.DecodeString("1234567890abcdef1234567890abcdef")
	copy(configDigest[:], cdBytes)

	signers := []common.Address{
		common.HexToAddress("0x1111111111111111111111111111111111111111"),
		common.HexToAddress("0x2222222222222222222222222222222222222222"),
		common.HexToAddress("0x3333333333333333333333333333333333333333"),
		common.HexToAddress("0x4444444444444444444444444444444444444444"),
	}
	f := uint8(1)

	setConfigPayload := SetConfig{
		VerifierAddress:            verifierAddr,
		ConfigDigest:               configDigest,
		Signers:                    signers,
		F:                          f,
		RecipientAddressesAndProps: []verifier_v0_5_0.CommonAddressAndWeight{},
	}

	callSetCfg := SetConfigConfig{
		ConfigsByChain: map[uint64][]SetConfig{
			testutil.TestChain.Selector: {setConfigPayload},
		},
		MCMSConfig: nil,
	}

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			SetConfigChangeset,
			callSetCfg,
		),
	)
	require.NoError(t, err)

	// Unset the verifier
	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			UnsetVerifierChangeset,
			VerifierProxyUnsetVerifierConfig{
				ConfigPerChain: map[uint64][]UnsetVerifierConfig{
					chainSelector: {
						{ContractAddress: verifierProxyAddr, ConfigDigest: configDigest},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	vp, err := verifier_proxy_v0_5_0.NewVerifierProxy(verifierProxyAddr, chain.Client)
	require.NoError(t, err)
	logIterator, err := vp.FilterVerifierUnset(nil)
	require.NoError(t, err)

	foundExpected := false
	for logIterator.Next() {
		if verifierAddr == logIterator.Event.VerifierAddress && configDigest == logIterator.Event.ConfigDigest {
			foundExpected = true
			break
		}
	}
	require.True(t, foundExpected)
	err = logIterator.Close()
	if err != nil {
		t.Errorf("Error closing log iterator: %v", err)
	}
}
