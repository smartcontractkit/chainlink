package changeset

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	verifier_proxy "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/verifier-proxy/0.5"
	verifier "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
)

func TestCallSetConfig(t *testing.T) {
	t.Skip("WIP")
	e := testutil.NewMemoryEnv(t, true)
	ctx := testcontext.Get(t)

	cc := verifier_proxy.DeployVerifierProxyConfig{
		ChainsToDeploy: map[uint64]verifier_proxy.DeployVerifierProxy{
			testutil.TestChain.Selector: {VerifierProxyAddress: common.Address{}},
		},
	}

	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			verifier_proxy.DeployVerifierProxyChangeset,
			cc,
		),
	)

	require.NoError(t, err)

	verifierProxyAddrHex, err := deployment.SearchAddressBook(e.ExistingAddresses, testutil.TestChain.Selector, types.VerifierProxy)
	require.NoError(t, err)
	verifierProxyAddr := common.HexToAddress(verifierProxyAddrHex)

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployVerifierChangeset,
			DeployVerifierConfig{
				ChainsToDeploy: map[uint64]DeployVerifier{
					testutil.TestChain.Selector: {VerifierProxyAddress: verifierProxyAddr},
				},
			},
		),
	)

	require.NoError(t, err)

	verifierAddrHex, err := deployment.SearchAddressBook(e.ExistingAddresses, testutil.TestChain.Selector, types.Verifier)
	require.NoError(t, err)
	verifierAddr := common.HexToAddress(verifierAddrHex)

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

	recipients := []verifier.CommonAddressAndWeight{
		{
			Addr:   common.HexToAddress("0x9999999999999999999999999999999999999999"),
			Weight: 10,
		},
		{
			Addr:   common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
			Weight: 20,
		},
	}

	setConfigPayload := SetConfig{
		VerifierAddress:            verifierAddr,
		configDigest:               configDigest,
		signers:                    signers,
		f:                          f,
		recipientAddressesAndProps: recipients,
	}

	callCfg := SetConfigConfig{
		ConfigsByChain: map[uint64][]SetConfig{
			testutil.TestChain.Selector: {setConfigPayload},
		},
		MCMSConfig: nil,
	}

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			SetConfigChangeset,
			callCfg,
		),
	)
	require.NoError(t, err)

	verifierABI, err := abi.JSON(strings.NewReader(verifier.VerifierMetaData.ABI))
	require.NoError(t, err)

	filter := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   nil,
		Addresses: []common.Address{verifierAddr},
		Topics:    [][]common.Hash{{verifierABI.Events["ConfigSet"].ID}},
	}
	logs, err := e.Chains[testutil.TestChain.Selector].Client.FilterLogs(ctx, filter)
	require.NoError(t, err)
	require.NotEmpty(t, logs, "Expected at least one ConfigSet event")

	for _, lg := range logs {
		decoded, err := verifierABI.Unpack("ConfigSet", lg.Data)
		require.NoError(t, err)

		configDigestFromTopics := lg.Topics[1]

		require.Len(t, decoded, 2, "decoded data should be [signers, f]")

		emittedSigners := decoded[0].([]common.Address)
		emittedF := decoded[1].(uint8)

		t.Logf("ConfigSet Event => configDigest: %s, signers: %v, f: %d",
			configDigestFromTopics.Hex(), emittedSigners, emittedF,
		)

		require.Equal(t, signers, emittedSigners, "Signers mismatch")
		require.Equal(t, f, emittedF, "Fault tolerance mismatch")
	}
}
