package launcher

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/services/ccipcapability/types"
	it "github.com/smartcontractkit/chainlink/v2/core/services/ocr3/plugins/ccip_integration_tests"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

func TestIntegration_Launcher(t *testing.T) {
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

	regSyncer, err := registrysyncer.New(lggr, uni, uni.CapReg.Address().String())
	require.NoError(t, err)

	hcr := uni.HomeChainReader
	launcher := New(
		it.CcipCapabilityVersion,
		it.CcipCapabilityLabelledName,
		p2pIDs[0],
		logger.TestLogger(t),
		hcr,
		&oracleCreatorPrints{
			t: t,
		},
		1*time.Second,
	)
	regSyncer.AddLauncher(launcher)

	require.NoError(t, launcher.Start(ctx))
	require.NoError(t, regSyncer.Start(ctx))
	t.Cleanup(func() { require.NoError(t, regSyncer.Close()) })
	t.Cleanup(func() { require.NoError(t, launcher.Close()) })

	chainAConf := it.SetupConfigInfo(it.ChainA, p2pIDs, it.FChainA, []byte("ChainA"))
	chainBConf := it.SetupConfigInfo(it.ChainB, p2pIDs[1:], it.FChainB, []byte("ChainB"))
	chainCConf := it.SetupConfigInfo(it.ChainC, p2pIDs[2:], it.FChainC, []byte("ChainC"))
	inputConfig := []ccip_config.CCIPConfigTypesChainConfigInfo{
		chainAConf,
		chainBConf,
		chainCConf,
	}
	_, err = uni.CcipCfg.ApplyChainConfigUpdates(uni.Transactor, nil, inputConfig)
	require.NoError(t, err)
	uni.Backend.Commit()

	ccipCapabilityID, err := uni.CapReg.GetHashedCapabilityId(nil, it.CcipCapabilityLabelledName, it.CcipCapabilityVersion)
	require.NoError(t, err)

	uni.AddDONToRegistry(
		ccipCapabilityID,
		it.ChainA,
		it.FChainA,
		p2pIDs[1],
		p2pIDs)

	gomega.NewWithT(t).Eventually(func() bool {
		return len(launcher.runningDONIDs()) == 1
	}, testutils.WaitTimeout(t), testutils.TestInterval).Should(gomega.BeTrue())
}

type oraclePrints struct {
	t           *testing.T
	pluginType  cctypes.PluginType
	config      cctypes.OCR3ConfigWithMeta
	isBootstrap bool
}

func (o *oraclePrints) Start() error {
	o.t.Logf("Starting oracle (pluginType: %s, isBootstrap: %t) with config %+v\n", o.pluginType, o.isBootstrap, o.config)
	return nil
}

func (o *oraclePrints) Close() error {
	o.t.Logf("Closing oracle (pluginType: %s, isBootstrap: %t) with config %+v\n", o.pluginType, o.isBootstrap, o.config)
	return nil
}

type oracleCreatorPrints struct {
	t *testing.T
}

func (o *oracleCreatorPrints) CreatePluginOracle(pluginType cctypes.PluginType, config cctypes.OCR3ConfigWithMeta) (cctypes.CCIPOracle, error) {
	o.t.Logf("Creating plugin oracle (pluginType: %s) with config %+v\n", pluginType, config)
	return &oraclePrints{pluginType: pluginType, config: config, t: o.t}, nil
}

func (o *oracleCreatorPrints) CreateBootstrapOracle(config cctypes.OCR3ConfigWithMeta) (cctypes.CCIPOracle, error) {
	o.t.Logf("Creating bootstrap oracle with config %+v\n", config)
	return &oraclePrints{pluginType: cctypes.PluginTypeCCIPCommit, config: config, isBootstrap: true, t: o.t}, nil
}

var _ cctypes.OracleCreator = &oracleCreatorPrints{}
var _ cctypes.CCIPOracle = &oraclePrints{}
