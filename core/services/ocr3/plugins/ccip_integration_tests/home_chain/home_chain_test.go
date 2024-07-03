package home_chain

import (
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/onsi/gomega"

	libocrtypes "github.com/smartcontractkit/libocr/ragep2p/types"

	ccipreader "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	capcfg "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	it "github.com/smartcontractkit/chainlink/v2/core/services/ocr3/plugins/ccip_integration_tests"

	"github.com/stretchr/testify/require"
)

func TestHomeChainReader(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	uni := it.NewTestUniverse(ctx, t, lggr)
	// We need 3*f + 1 p2pIDs to have enough nodes to bootstrap
	var arr []int64
	n := int(it.FChainA*3 + 1)
	for i := 0; i <= n; i++ {
		arr = append(arr, int64(i))
	}
	p2pIDs := it.P2pIDsFromInts(arr)
	uni.AddCapability(p2pIDs)
	//==============================Apply configs to Capability Contract=================================
	chainAConf := it.SetupConfigInfo(it.ChainA, p2pIDs, it.FChainA, []byte("ChainA"))
	chainBConf := it.SetupConfigInfo(it.ChainB, p2pIDs[1:], it.FChainB, []byte("ChainB"))
	chainCConf := it.SetupConfigInfo(it.ChainC, p2pIDs[2:], it.FChainC, []byte("ChainC"))
	inputConfig := []capcfg.CCIPConfigTypesChainConfigInfo{
		chainAConf,
		chainBConf,
		chainCConf,
	}
	_, err := uni.CcipCfg.ApplyChainConfigUpdates(uni.Transactor, nil, inputConfig)
	require.NoError(t, err)
	uni.Backend.Commit()
	//================================Setup HomeChainReader===============================

	pollDuration := time.Second
	homeChain := uni.HomeChainReader

	gomega.NewWithT(t).Eventually(func() bool {
		configs, _ := homeChain.GetAllChainConfigs()
		return configs != nil
	}, testutils.WaitTimeout(t), pollDuration*5).Should(gomega.BeTrue())

	t.Logf("homchain reader is ready")
	//================================Test HomeChain Reader===============================
	expectedChainConfigs := map[cciptypes.ChainSelector]ccipreader.ChainConfig{}
	for _, c := range inputConfig {
		expectedChainConfigs[cciptypes.ChainSelector(c.ChainSelector)] = ccipreader.ChainConfig{
			FChain:         int(c.ChainConfig.FChain),
			SupportedNodes: toPeerIDs(c.ChainConfig.Readers),
		}
	}
	configs, err := homeChain.GetAllChainConfigs()
	require.NoError(t, err)
	require.Equal(t, expectedChainConfigs, configs)
	//=================================Remove ChainC from OnChainConfig=========================================
	_, err = uni.CcipCfg.ApplyChainConfigUpdates(uni.Transactor, []uint64{it.ChainC}, nil)
	require.NoError(t, err)
	uni.Backend.Commit()
	time.Sleep(pollDuration * 5) // Wait for the chain reader to update
	configs, err = homeChain.GetAllChainConfigs()
	require.NoError(t, err)
	delete(expectedChainConfigs, cciptypes.ChainSelector(it.ChainC))
	require.Equal(t, expectedChainConfigs, configs)
	//================================Close HomeChain Reader===============================
	//require.NoError(t, homeChain.Close())
	//require.NoError(t, uni.LogPoller.Close())
	t.Logf("homchain reader successfully closed")
}

func toPeerIDs(readers [][32]byte) mapset.Set[libocrtypes.PeerID] {
	peerIDs := mapset.NewSet[libocrtypes.PeerID]()
	for _, r := range readers {
		peerIDs.Add(r)
	}
	return peerIDs
}
