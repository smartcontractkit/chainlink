package rmn

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccip_integration_tests/integrationhelpers"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	assert "github.com/stretchr/testify/assert"

	readerpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
)

func TestRMNHomeReader_GetRMNNodesInfo(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	uni := integrationhelpers.NewTestUniverse(ctx, t, lggr)
	zeroBytes := [32]byte{0}

	const (
		chainID1        = 1
		minObservers1   = 1
		observerBitmap1 = 1

		chainID2        = 2
		minObservers2   = 0
		observerBitmap2 = 1
	)

	//================================Deploy and configure RMNHome===============================
	rmnHomeAddress, _, rmnHome, err := rmn_home.DeployRMNHome(uni.Transactor, uni.Backend)
	require.NoError(t, err)
	uni.Backend.Commit()

	staticConfig, dynamicConfig := integrationhelpers.GenerateRMNHomeConfigs(
		"PeerID1",
		"DummyPublicKey1",
		"This is a sample offchain configuration in the static config",
		chainID1,
		minObservers1,
		big.NewInt(observerBitmap1),
	)

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

	rmnHomeReader := readerpkg.NewRMNHomePoller(uni.HomeContractReader, rmnHomeBoundContract, lggr, 1*time.Millisecond)

	err = rmnHomeReader.Start(testutils.Context(t))
	require.NoError(t, err)

	//================================Test RMNHome Reader===============================
	expectedNodesInfo := integrationhelpers.GenerateExpectedRMNHomeNodesInfo(staticConfig, chainID1)

	testutils.AssertEventually(t, func() bool {
		nodesInfo, err2 := rmnHomeReader.GetRMNNodesInfo(configDigest)
		if err2 != nil {
			t.Logf("Error getting RMN nodes info: %v", err2)
			return false
		}

		if !assert.Equal(t, expectedNodesInfo, nodesInfo) {
			t.Logf("Expected nodes info doesn't match actual nodes info")
			return false
		}

		return true
	})

	// Add a new candidate config
	staticConfig, dynamicConfig = integrationhelpers.GenerateRMNHomeConfigs(
		"PeerID2",
		"DummyPublicKey2",
		"This is a sample offchain configuration in the static config 2",
		chainID2,
		minObservers2,
		big.NewInt(observerBitmap2),
	)

	_, err = rmnHome.SetCandidate(uni.Transactor, staticConfig, dynamicConfig, zeroBytes)
	require.NoError(t, err)
	uni.Backend.Commit()

	candidateConfigDigest, err := rmnHome.GetCandidateDigest(&bind.CallOpts{})
	require.NoError(t, err)

	expectedCandidateNodesInfo := integrationhelpers.GenerateExpectedRMNHomeNodesInfo(staticConfig, chainID2)

	testutils.AssertEventually(t, func() bool {
		nodesInfo, err2 := rmnHomeReader.GetRMNNodesInfo(candidateConfigDigest)
		if err2 != nil {
			t.Logf("Error getting RMN nodes info: %v", err2)
			return false
		}

		if !assert.Equal(t, expectedCandidateNodesInfo, nodesInfo) {
			t.Logf("Expected nodes info doesn't match actual nodes info")
			return false
		}

		return true
	})

	// Promote the candidate config
	_, err = rmnHome.PromoteCandidateAndRevokeActive(uni.Transactor, candidateConfigDigest, configDigest)
	require.NoError(t, err)
	uni.Backend.Commit()

	testutils.AssertEventually(t, func() bool {
		nodesInfo, err2 := rmnHomeReader.GetRMNNodesInfo(candidateConfigDigest)
		if err2 != nil {
			t.Logf("Error getting RMN nodes info: %v", err2)
			return false
		}

		if !assert.Equal(t, expectedCandidateNodesInfo, nodesInfo) {
			t.Logf("Expected nodes info doesn't match actual nodes info")
			return false
		}

		isPrevConfigStillSet := rmnHomeReader.IsRMNHomeConfigDigestSet(configDigest)
		if isPrevConfigStillSet {
			t.Logf("Previous config is still set")
			return false
		}

		return true
	})
}
