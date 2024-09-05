package ccip_integration_tests

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccip_integration_tests/integrationhelpers"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/onsi/gomega"

	libocrtypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	ccipreader "github.com/smartcontractkit/chainlink-ccip/pkg/reader"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/stretchr/testify/require"

	capcfg "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestHomeChainReader(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	uni := integrationhelpers.NewTestUniverse(ctx, t, lggr)
	// We need 3*f + 1 p2pIDs to have enough nodes to bootstrap
	var arr []int64
	n := int(integrationhelpers.FChainA*3 + 1)
	for i := 0; i <= n; i++ {
		arr = append(arr, int64(i))
	}
	p2pIDs := integrationhelpers.P2pIDsFromInts(arr)
	uni.AddCapability(p2pIDs)

	//==============================Apply configs to Capability Contract=================================
	encodedChainConfig, err := chainconfig.EncodeChainConfig(chainconfig.ChainConfig{
		GasPriceDeviationPPB:    cciptypes.NewBigIntFromInt64(1000),
		DAGasPriceDeviationPPB:  cciptypes.NewBigIntFromInt64(1_000_000),
		FinalityDepth:           -1,
		OptimisticConfirmations: 1,
	})
	require.NoError(t, err)
	inputConfig := []capcfg.CCIPConfigTypesChainConfigInfo{
		integrationhelpers.SetupConfigInfo(integrationhelpers.ChainA, p2pIDs, integrationhelpers.FChainA, encodedChainConfig),
		integrationhelpers.SetupConfigInfo(integrationhelpers.ChainB, p2pIDs[1:], integrationhelpers.FChainB, encodedChainConfig),
		integrationhelpers.SetupConfigInfo(integrationhelpers.ChainC, p2pIDs[2:], integrationhelpers.FChainC, encodedChainConfig),
	}
	_, err = uni.CcipCfg.ApplyChainConfigUpdates(uni.Transactor, nil, inputConfig)
	require.NoError(t, err)
	uni.Backend.Commit()
	chainConfigInfos, err := uni.CcipCfg.GetAllChainConfigs(nil, big.NewInt(0), big.NewInt(100))
	require.NoError(t, err)
	require.Len(t, chainConfigInfos, len(inputConfig))

	// Wait for the home chain reader to read the expected amount of chain configs.
	gomega.NewWithT(t).Eventually(func() bool {
		configs, _ := uni.HomeChainReader.GetAllChainConfigs()
		return len(configs) == len(inputConfig)
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	t.Logf("homchain reader is ready")

	//================================Test HomeChain Reader===============================
	expectedChainConfigs := map[cciptypes.ChainSelector]ccipreader.ChainConfig{}
	for _, c := range inputConfig {
		expectedChainConfigs[cciptypes.ChainSelector(c.ChainSelector)] = ccipreader.ChainConfig{
			FChain:         int(c.ChainConfig.FChain),
			SupportedNodes: toPeerIDs(c.ChainConfig.Readers),
			Config:         mustDecodeChainConfig(t, c.ChainConfig.Config),
		}
	}

	configs, err := uni.HomeChainReader.GetAllChainConfigs()
	require.NoError(t, err)

	require.Equal(t, expectedChainConfigs, configs)

	// Remove chain C from the chain configs and expect the home chain reader to
	// update its state accordingly.
	_, err = uni.CcipCfg.ApplyChainConfigUpdates(uni.Transactor, []uint64{integrationhelpers.ChainC}, nil)
	require.NoError(t, err)
	uni.Backend.Commit()

	// Wait for the home chain reader to read the expected amount of chain configs.
	gomega.NewWithT(t).Eventually(func() bool {
		chainConfigs, _ := uni.HomeChainReader.GetAllChainConfigs()
		return len(chainConfigs) == len(inputConfig)-1
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())
	configs, err = uni.HomeChainReader.GetAllChainConfigs()
	require.NoError(t, err)

	delete(expectedChainConfigs, cciptypes.ChainSelector(integrationhelpers.ChainC))
	require.Equal(t, expectedChainConfigs, configs)
}

func toPeerIDs(readers [][32]byte) mapset.Set[libocrtypes.PeerID] {
	peerIDs := mapset.NewSet[libocrtypes.PeerID]()
	for _, r := range readers {
		peerIDs.Add(r)
	}
	return peerIDs
}

func mustDecodeChainConfig(t *testing.T, encodedChainConfig []byte) chainconfig.ChainConfig {
	chainConfig, err := chainconfig.DecodeChainConfig(encodedChainConfig)
	require.NoError(t, err)
	return chainConfig
}
