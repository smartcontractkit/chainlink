package launcher

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	it "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccip_integration_tests/integrationhelpers"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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

	db := pgtest.NewSqlxDB(t)
	regSyncer, err := registrysyncer.New(lggr,
		func() (p2ptypes.PeerID, error) {
			return p2pIDs[0], nil
		},
		uni,
		uni.CapReg.Address().String(),
		registrysyncer.NewORM(db, lggr),
	)
	require.NoError(t, err)

	oracleCreator := &oracleCreatorPrints{
		t: t,
	}
	launcher := New(
		it.CapabilityID,
		p2pIDs[0],
		logger.TestLogger(t),
		uni.HomeChainReader,
		1*time.Second,
		oracleCreator,
	)
	regSyncer.AddLauncher(launcher)

	require.NoError(t, launcher.Start(ctx))
	require.NoError(t, regSyncer.Start(ctx))
	t.Cleanup(func() { require.NoError(t, regSyncer.Close()) })
	t.Cleanup(func() { require.NoError(t, launcher.Close()) })

	encodedChainConfig, err := chainconfig.EncodeChainConfig(chainconfig.ChainConfig{
		GasPriceDeviationPPB:    cciptypes.NewBigIntFromInt64(1000),
		DAGasPriceDeviationPPB:  cciptypes.NewBigIntFromInt64(1_000_000),
		OptimisticConfirmations: 1,
	})
	require.NoError(t, err)

	chainAConf := it.SetupConfigInfo(it.ChainA, p2pIDs, it.FChainA, encodedChainConfig)
	chainBConf := it.SetupConfigInfo(it.ChainB, p2pIDs[1:], it.FChainB, encodedChainConfig)
	chainCConf := it.SetupConfigInfo(it.ChainC, p2pIDs[2:], it.FChainC, encodedChainConfig)
	inputConfig := []ccip_home.CCIPHomeChainConfigArgs{
		chainAConf,
		chainBConf,
		chainCConf,
	}
	_, err = uni.CCIPHome.ApplyChainConfigUpdates(uni.Transactor, nil, inputConfig)
	require.NoError(t, err)
	uni.Backend.Commit()

	ccipCapabilityID, err := uni.CapReg.GetHashedCapabilityId(nil, it.CcipCapabilityLabelledName, it.CcipCapabilityVersion)
	require.NoError(t, err)

	uni.AddDONToRegistry(
		ccipCapabilityID,
		it.ChainA,
		it.FChainA,
		p2pIDs)

	require.Eventually(t, func() bool {
		return len(launcher.runningDONIDs()) == 1
	}, testutils.WaitTimeout(t), testutils.TestInterval)
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

func (o *oracleCreatorPrints) Create(ctx context.Context, _ uint32, config cctypes.OCR3ConfigWithMeta) (cctypes.CCIPOracle, error) {
	pluginType := cctypes.PluginType(config.Config.PluginType)
	o.t.Logf("Creating plugin oracle (pluginType: %s) with config %+v\n", pluginType, config)
	return &oraclePrints{pluginType: pluginType, config: config, t: o.t}, nil
}

func (o *oracleCreatorPrints) Type() cctypes.OracleType {
	return cctypes.OracleTypePlugin
}

var _ cctypes.OracleCreator = &oracleCreatorPrints{}
var _ cctypes.CCIPOracle = &oraclePrints{}
