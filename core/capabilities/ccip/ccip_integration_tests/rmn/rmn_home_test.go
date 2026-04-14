package rmn

import (
	"bytes"
	"encoding/json"
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"
	readerpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip/consts"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	configsevm "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccip_integration_tests/integrationhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestRMNRemote_ChainRead(t *testing.T) {
	ctx := t.Context()
	lggr := logger.Test(t)
	uni := integrationhelpers.NewTestUniverse(ctx, t, lggr)

	// Deploy RMNRemote
	rmnRemoteAddr, _, rmnRemote, err := rmn_remote.DeployRMNRemote(uni.Transactor, uni.Backend.Client(), uint64(integrationhelpers.ChainA), common.HexToAddress("0x1"))
	require.NoError(t, err)
	uni.Backend.Commit()

	destReaderConfig, err := json.Marshal(configsevm.DestReaderConfig)
	require.NoError(t, err)
	contractReader, err := uni.NewContractReader(ctx, destReaderConfig)
	require.NoError(t, err)

	// Curse some chain on RMNRemote
	subj := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	_, err = rmnRemote.Curse(uni.Transactor, subj)
	require.NoError(t, err)
	uni.Backend.Commit()

	// Get the cursed subjects
	cursedSubjects, err := rmnRemote.GetCursedSubjects(&bind.CallOpts{})
	require.NoError(t, err)
	require.Len(t, cursedSubjects, 1)
	require.Equal(t, subj, cursedSubjects[0])

	// Use chain reader to get the cursed subjects
	// Bind the contracts first.
	err = contractReader.Bind(ctx, []types.BoundContract{
		{
			Address: rmnRemoteAddr.String(),
			Name:    consts.ContractNameRMNRemote,
		},
	})
	require.NoError(t, err)

	// Call the method via the contract reader.
	// cciptypes.RMNCurseResponse is used by the plugin to unmarshal the response into.
	// Due to it having an incorrectly named field relative to the RMNRemote ABI,
	// which uses a named return, the unmarshal into it will silently fail.
	var cursedSubjectsResponse cciptypes.RMNCurseResponse
	err = contractReader.GetLatestValue(ctx,
		types.BoundContract{
			Address: rmnRemoteAddr.String(),
			Name:    consts.ContractNameRMNRemote,
		}.ReadIdentifier(consts.MethodNameGetCursedSubjects),
		primitives.Unconfirmed,
		map[string]any{}, // args
		&cursedSubjectsResponse,
	)
	require.NoError(t, err)
	require.Empty(t, cursedSubjectsResponse.CursedSubjects)

	// This works because the type has a field matching the name of the named return in the RMNRemote ABI.
	// However, it is a breaking change to the struct.
	type resp struct {
		Subjects [][16]byte
	}
	var r resp
	err = contractReader.GetLatestValue(ctx,
		types.BoundContract{
			Address: rmnRemoteAddr.String(),
			Name:    consts.ContractNameRMNRemote,
		}.ReadIdentifier(consts.MethodNameGetCursedSubjects),
		primitives.Unconfirmed,
		map[string]any{}, // args
		&r,
	)
	require.NoError(t, err)
	// Assert that the cursed subject is in the response.
	require.Len(t, r.Subjects, 1)
	require.Equal(t, subj, r.Subjects[0])

	// This doesn't work, despite the json tag matching the named return in the RMNRemote ABI.
	type respDiffJSONTag struct {
		CursedSubjects [][16]byte `json:"subjects"`
	}
	var r2 respDiffJSONTag
	err = contractReader.GetLatestValue(ctx,
		types.BoundContract{
			Address: rmnRemoteAddr.String(),
			Name:    consts.ContractNameRMNRemote,
		}.ReadIdentifier(consts.MethodNameGetCursedSubjects),
		primitives.Unconfirmed,
		map[string]any{}, // args
		&r2,
	)
	require.NoError(t, err)
	require.Empty(t, r2.CursedSubjects)

	// This works because the codec seems to use mapstructure to decode the response
	// into the struct.
	// See
	// 1. https://github.com/smartcontractkit/chainlink-evm/blob/2dca02f24e983863e72dfce3d715c716c23f6016/pkg/codec/decoder.go#L48
	// 2. https://github.com/smartcontractkit/chainlink-evm/blob/2dca02f24e983863e72dfce3d715c716c23f6016/pkg/codec/decoder.go#L95-L110
	type respDiffMapstructureTag struct {
		CursedSubjects [][16]byte `mapstructure:"subjects"`
	}
	var r3 respDiffMapstructureTag
	err = contractReader.GetLatestValue(ctx,
		types.BoundContract{
			Address: rmnRemoteAddr.String(),
			Name:    consts.ContractNameRMNRemote,
		}.ReadIdentifier(consts.MethodNameGetCursedSubjects),
		primitives.Unconfirmed,
		map[string]any{}, // args
		&r3,
	)
	require.NoError(t, err)
	require.Len(t, r3.CursedSubjects, 1)
	require.Equal(t, subj, r3.CursedSubjects[0])
}

func TestRMNHomeReader_GetRMNNodesInfo(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.Test(t)
	uni := integrationhelpers.NewTestUniverse(ctx, t, lggr)
	zeroBytes := [32]byte{0}

	const (
		chainID1        = 1
		f1              = 0
		observerBitmap1 = 1

		chainID2        = 2
		f2              = 0
		observerBitmap2 = 1
	)

	// ================================Deploy and configure RMNHome===============================
	rmnHomeAddress, _, rmnHome, err := rmn_home.DeployRMNHome(uni.Transactor, uni.Backend.Client())
	require.NoError(t, err)
	uni.Backend.Commit()

	staticConfig, dynamicConfig, err := integrationhelpers.GenerateRMNHomeConfigs(
		"PeerID1",
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"This is a sample offchain configuration in the static config",
		chainID1,
		f1,
		big.NewInt(observerBitmap1),
	)
	require.NoError(t, err)

	_, err = rmnHome.SetCandidate(uni.Transactor, staticConfig, dynamicConfig, zeroBytes)
	require.NoError(t, err)
	uni.Backend.Commit()

	configDigest, err := rmnHome.GetCandidateDigest(&bind.CallOpts{})
	require.NoError(t, err)

	_, err = rmnHome.PromoteCandidateAndRevokeActive(uni.Transactor, configDigest, zeroBytes)
	require.NoError(t, err)
	uni.Backend.Commit()

	rmnHomeBoundContract := types.BoundContract{
		Address: rmnHomeAddress.String(),
		Name:    consts.ContractNameRMNHome,
	}

	err = uni.HomeContractReader.Bind(testutils.Context(t), []types.BoundContract{rmnHomeBoundContract})
	require.NoError(t, err)

	rmnHomeReader, err := readerpkg.NewRMNHomeChainReader(
		ctx,
		lggr,
		100*time.Millisecond,
		cciptypes.ChainSelector(1),
		rmnHomeAddress.Bytes(),
		uni.HomeContractReader,
	)
	require.NoError(t, err)

	err = rmnHomeReader.Start(testutils.Context(t))
	require.NoError(t, err)

	t.Cleanup(func() {
		err1 := rmnHomeReader.Close()
		require.NoError(t, err1)
	})

	// ================================Test RMNHome Reader===============================
	expectedNodesInfo := integrationhelpers.GenerateExpectedRMNHomeNodesInfo(staticConfig, chainID1)

	require.Eventually(
		t,
		assertRMNHomeNodesInfo(t, rmnHomeReader, configDigest, expectedNodesInfo, nil),
		5*time.Second,
		100*time.Millisecond,
	)

	// Add a new candidate config
	staticConfig2, dynamicConfig2, err := integrationhelpers.GenerateRMNHomeConfigs(
		"PeerID2",
		"1123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"This is a sample offchain configuration in the static config 2",
		chainID2,
		f2,
		big.NewInt(observerBitmap2),
	)
	require.NoError(t, err)

	_, err = rmnHome.SetCandidate(uni.Transactor, staticConfig2, dynamicConfig2, zeroBytes)
	require.NoError(t, err)
	uni.Backend.Commit()

	candidateConfigDigest, err := rmnHome.GetCandidateDigest(&bind.CallOpts{})
	require.NoError(t, err)

	expectedCandidateNodesInfo := integrationhelpers.GenerateExpectedRMNHomeNodesInfo(staticConfig2, chainID2)

	require.Eventually(
		t,
		assertRMNHomeNodesInfo(t, rmnHomeReader, candidateConfigDigest, expectedCandidateNodesInfo, nil),
		5*time.Second,
		100*time.Millisecond,
	)

	// Promote the candidate config
	_, err = rmnHome.PromoteCandidateAndRevokeActive(uni.Transactor, candidateConfigDigest, configDigest)
	require.NoError(t, err)
	uni.Backend.Commit()

	require.Eventually(
		t,
		assertRMNHomeNodesInfo(t, rmnHomeReader, candidateConfigDigest, expectedCandidateNodesInfo, &configDigest),
		5*time.Second,
		100*time.Millisecond,
	)
}

func assertRMNHomeNodesInfo(
	t *testing.T,
	rmnHomeReader readerpkg.RMNHome,
	configDigest [32]byte,
	expectedNodesInfo []readerpkg.HomeNodeInfo,
	prevConfigDigest *[32]byte,
) func() bool {
	return func() bool {
		nodesInfo, err := rmnHomeReader.GetRMNNodesInfo(configDigest)
		if err != nil {
			t.Logf("Error getting RMN nodes info: %v", err)
			return false
		}

		equal := slices.EqualFunc(expectedNodesInfo, nodesInfo, func(a, b readerpkg.HomeNodeInfo) bool {
			return a.ID == b.ID &&
				a.PeerID == b.PeerID &&
				bytes.Equal(*a.OffchainPublicKey, *b.OffchainPublicKey) &&
				a.SupportedSourceChains.Equal(b.SupportedSourceChains)
		})

		if !equal {
			t.Logf("Expected nodes info doesn't match actual nodes info")
			t.Logf("Expected: %+v", expectedNodesInfo)
			t.Logf("Actual: %+v", nodesInfo)
			return false
		}

		if prevConfigDigest != nil {
			isPrevConfigStillSet := rmnHomeReader.IsRMNHomeConfigDigestSet(*prevConfigDigest)
			if isPrevConfigStillSet {
				t.Logf("Previous config is still set")
				return false
			}
		}

		return true
	}
}
