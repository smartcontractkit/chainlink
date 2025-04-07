package verification

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	verifier "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func TestCallActivateConfig(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true)
	chainSelector := testutil.TestChain.Selector

	e, _, verifierAddr := DeployVerifierProxyAndVerifier(t, e)

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
		RecipientAddressesAndProps: []verifier.CommonAddressAndWeight{},
	}

	callSetCfg := SetConfigConfig{
		ConfigsByChain: map[uint64][]SetConfig{
			chainSelector: {setConfigPayload},
		},
		MCMSConfig: nil,
	}

	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			SetConfigChangeset,
			callSetCfg,
		),
	)
	require.NoError(t, err)

	activateConfigPayload := ActivateConfig{
		VerifierAddress: verifierAddr,
		ConfigDigest:    configDigest,
	}

	callCfg := ActivateConfigConfig{
		ConfigsByChain: map[uint64][]ActivateConfig{
			chainSelector: {activateConfigPayload},
		},
		MCMSConfig: nil,
	}

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			ActivateConfigChangeset,
			callCfg,
		),
	)
	require.NoError(t, err)

	chain := e.Chains[chainSelector]
	vp, err := verifier.NewVerifier(verifierAddr, chain.Client)
	require.NoError(t, err)
	logIterator, err := vp.FilterConfigActivated(nil, [][32]byte{configDigest})
	require.NoError(t, err)
	defer logIterator.Close()
	require.NoError(t, err)
	foundExpected := false

	for logIterator.Next() {
		event := logIterator.Event
		if configDigest == event.ConfigDigest {
			foundExpected = true
			break
		}
	}
	require.True(t, foundExpected)
}
