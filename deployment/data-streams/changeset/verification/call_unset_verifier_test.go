package verification

import (
	"encoding/hex"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestUnsetVerifier(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true, 0)

	chainSelector := e.AllChainSelectors()[0]
	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployVerifierProxyChangeset,
			DeployVerifierProxyConfig{
				ChainsToDeploy: map[uint64]DeployVerifierProxy{
					chainSelector: {AccessControllerAddress: common.Address{}},
				},
				Version: *semver.MustParse("0.5.0"),
			},
		),
	)
	require.NoError(t, err)

	// Ensure the VerifierProxy was deployed
	verifierProxyAddrHex, err := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, types.VerifierProxy)
	require.NoError(t, err)
	verifierProxyAddr := common.HexToAddress(verifierProxyAddrHex)

	// Deploy Verifier
	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployVerifierChangeset,
			DeployVerifierConfig{
				ChainsToDeploy: map[uint64]DeployVerifier{
					chainSelector: {VerifierProxyAddress: verifierProxyAddr}},
			},
		),
	)
	require.NoError(t, err)

	// Ensure the Verifier was deployed
	verifierAddrHex, err := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, types.Verifier)
	require.NoError(t, err)
	verifierAddr := common.HexToAddress(verifierAddrHex)

	chain := e.Chains[chainSelector]

	verifier, err := verifier_v0_5_0.NewVerifier(verifierAddr, e.Chains[chainSelector].Client)
	require.NoError(t, err)
	require.NotNil(t, verifier)

	// Initialize the verifier
	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			InitializeVerifierChangeset,
			VerifierProxyInitializeVerifierConfig{
				ConfigPerChain: map[uint64][]InitializeVerifierConfig{
					testutil.TestChain.Selector: {
						{VerifierAddress: verifierAddr, ContractAddress: verifierProxyAddr},
					},
				},
			},
		),
	)
	require.NoError(t, err)

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
