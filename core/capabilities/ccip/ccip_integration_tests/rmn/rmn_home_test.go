package rmn

import (
	"bytes"
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	readerpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccip_integration_tests/integrationhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestRMNHomeReader_GetRMNNodesInfo(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
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
