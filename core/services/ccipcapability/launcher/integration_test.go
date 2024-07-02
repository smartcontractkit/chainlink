package launcher

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	ccipreader "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ocr3_config_encoder"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/services/ccipcapability/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const (
	chainA  uint64 = 1
	fChainA uint8  = 1

	chainB  uint64 = 2
	fChainB uint8  = 2

	chainC  uint64 = 3
	fChainC uint8  = 3

	ccipCapabilityLabelledName = "ccip"
	ccipCapabilityVersion      = "v1.0"
)

type testUniverse struct {
	transactor *bind.TransactOpts
	backend    *backends.SimulatedBackend
	capReg     *kcr.CapabilitiesRegistry
	cc         *ccip_config.CCIPConfig
	testingT   *testing.T
	lp         logpoller.LogPoller
	simClient  client.Client
}

func newTestUniverse(t *testing.T) testUniverse {
	transactor := testutils.MustNewSimTransactor(t)
	backend := backends.NewSimulatedBackend(core.GenesisAlloc{
		transactor.From: {Balance: assets.Ether(1000).ToInt()},
	}, 30e6)

	crAddress, _, _, err := kcr.DeployCapabilitiesRegistry(transactor, backend)
	require.NoError(t, err)
	backend.Commit()

	capReg, err := kcr.NewCapabilitiesRegistry(crAddress, backend)
	require.NoError(t, err)

	ccAddress, _, _, err := ccip_config.DeployCCIPConfig(transactor, backend, crAddress)
	require.NoError(t, err)
	backend.Commit()

	cc, err := ccip_config.NewCCIPConfig(ccAddress, backend)
	require.NoError(t, err)

	return testUniverse{
		transactor: transactor,
		backend:    backend,
		capReg:     capReg,
		cc:         cc,
		testingT:   t,
	}
}

func (t testUniverse) NewContractReader(ctx context.Context, cfg []byte) (types.ContractReader, error) {
	var config evmrelaytypes.ChainReaderConfig
	err := json.Unmarshal(cfg, &config)
	require.NoError(t.testingT, err)
	return evm.NewChainReaderService(ctx, logger.TestLogger(t.testingT), t.lp, t.simClient, config)
}

func addCapabilities(
	t *testing.T,
	backend *backends.SimulatedBackend,
	transactor *bind.TransactOpts,
	capReg *kcr.CapabilitiesRegistry,
	capConfAddress common.Address) [][32]byte {
	// add the CCIP capability to the registry
	_, err := capReg.AddCapabilities(transactor, []kcr.CapabilitiesRegistryCapability{
		{
			LabelledName:          ccipCapabilityLabelledName,
			Version:               ccipCapabilityVersion,
			CapabilityType:        0,
			ResponseType:          0,
			ConfigurationContract: capConfAddress,
		},
	})
	require.NoError(t, err, "failed to add capability to registry")
	backend.Commit()

	ccipCapabilityID, err := capReg.GetHashedCapabilityId(nil, ccipCapabilityLabelledName, ccipCapabilityVersion)
	require.NoError(t, err)

	// Add the p2p ids of the ccip nodes
	var p2pIDs [][32]byte
	for i := 0; i < 4; i++ {
		p2pID := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(i + 1))).PeerID()
		p2pIDs = append(p2pIDs, p2pID)
		_, err = capReg.AddNodeOperators(transactor, []kcr.CapabilitiesRegistryNodeOperator{
			{
				Admin: transactor.From,
				Name:  fmt.Sprintf("nop-%d", i),
			},
		})
		require.NoError(t, err)
		backend.Commit()

		// get the node operator id from the event
		it, err := capReg.FilterNodeOperatorAdded(nil, nil, nil)
		require.NoError(t, err)
		var nodeOperatorID uint32
		for it.Next() {
			if it.Event.Name == fmt.Sprintf("nop-%d", i) {
				nodeOperatorID = it.Event.NodeOperatorId
				break
			}
		}
		require.NotZero(t, nodeOperatorID)

		_, err = capReg.AddNodes(transactor, []kcr.CapabilitiesRegistryNodeParams{
			{
				NodeOperatorId:      nodeOperatorID,
				Signer:              testutils.Random32Byte(),
				P2pId:               p2pID,
				HashedCapabilityIds: [][32]byte{ccipCapabilityID},
			},
		})
		require.NoError(t, err)
		backend.Commit()

		// verify that the node was added successfully
		nodeInfo, err := capReg.GetNode(nil, p2pID)
		require.NoError(t, err)

		require.Equal(t, nodeOperatorID, nodeInfo.NodeOperatorId)
		require.Equal(t, p2pID[:], nodeInfo.P2pId[:])
	}
	return p2pIDs
}

func setupConfigInfo(chainSelector uint64, readers [][32]byte, fChain uint8, cfg []byte) ccip_config.CCIPConfigTypesChainConfigInfo {
	return ccip_config.CCIPConfigTypesChainConfigInfo{
		ChainSelector: chainSelector,
		ChainConfig: ccip_config.CCIPConfigTypesChainConfig{
			Readers: readers,
			FChain:  fChain,
			Config:  cfg,
		},
	}
}

func newHomeChainReader(t *testing.T, logPoller logpoller.LogPoller, client client.Client, ccAddress common.Address) cctypes.HomeChainReader {
	cfg := evmrelaytypes.ChainReaderConfig{
		Contracts: map[string]evmrelaytypes.ChainContractReader{
			"CCIPConfig": {
				ContractABI: ccip_config.CCIPConfigMetaData.ABI,
				Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
					"getAllChainConfigs": {
						ChainSpecificName: "getAllChainConfigs",
					},
					"getOCRConfig": {
						ChainSpecificName: "getOCRConfig",
					},
				},
			},
		},
	}
	cr, err := evm.NewChainReaderService(testutils.Context(t), logger.TestLogger(t), logPoller, client, cfg)
	require.NoError(t, err)

	err = cr.Bind(testutils.Context(t), []types.BoundContract{
		{
			Address: ccAddress.String(),
			Name:    "CCIPConfig",
		},
	})
	require.NoError(t, err)
	require.NoError(t, cr.Start(testutils.Context(t)))

	hcr := ccipreader.NewHomeChainReader(cr, logger.TestLogger(t), time.Second)
	require.NoError(t, hcr.Start(testutils.Context(t)))

	return hcr
}

func addDONToRegistry(t *testing.T,
	transactor *bind.TransactOpts,
	ccipCapabilityID [32]byte,
	chainSelector uint64,
	f uint8,
	capReg *kcr.CapabilitiesRegistry,
	backend *backends.SimulatedBackend,
	bootstrapP2PID [32]byte,
	p2pIDs [][32]byte,
) {
	tabi, err := ocr3_config_encoder.IOCR3ConfigEncoderMetaData.GetAbi()
	require.NoError(t, err)

	var (
		signers      [][]byte
		transmitters [][]byte
	)
	for range p2pIDs {
		signers = append(signers, testutils.NewAddress().Bytes())
		transmitters = append(transmitters, testutils.NewAddress().Bytes())
	}

	var ocr3Configs []ocr3_config_encoder.CCIPConfigTypesOCR3Config
	for _, pluginType := range []cctypes.PluginType{cctypes.PluginTypeCCIPCommit, cctypes.PluginTypeCCIPExec} {
		ocr3Configs = append(ocr3Configs, ocr3_config_encoder.CCIPConfigTypesOCR3Config{
			PluginType:            uint8(pluginType),
			ChainSelector:         chainSelector,
			F:                     f,
			OffchainConfigVersion: 30,
			OfframpAddress:        testutils.NewAddress().Bytes(),
			BootstrapP2PIds:       [][32]byte{bootstrapP2PID},
			P2pIds:                p2pIDs,
			Signers:               signers,
			Transmitters:          transmitters,
			OffchainConfig:        []byte("offchain config"),
		})
	}

	encodedCall, err := tabi.Pack("exposeOCR3Config", ocr3Configs)
	require.NoError(t, err)

	// Trim first four bytes to remove function selector.
	encodedConfigs := encodedCall[4:]

	_, err = capReg.AddDON(transactor, p2pIDs, []kcr.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: ccipCapabilityID,
			Config:       encodedConfigs,
		},
	}, false, false, f)
	require.NoError(t, err)
	backend.Commit()
}

func TestIntegration_Launcher(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	uni := newTestUniverse(t)

	db := pgtest.NewSqlxDB(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            0,
		BackfillBatchSize:        10,
		RpcBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	cl := client.NewSimulatedBackendClient(t, uni.backend, big.NewInt(1337))
	lp := logpoller.NewLogPoller(logpoller.NewORM(big.NewInt(1337), db, lggr), cl, logger.NullLogger, lpOpts)
	require.NoError(t, lp.Start(ctx))
	t.Cleanup(func() { require.NoError(t, lp.Close()) })

	uni.lp = lp
	uni.simClient = cl

	p2pIDs := addCapabilities(t, uni.backend, uni.transactor, uni.capReg, uni.cc.Address())

	regSyncer, err := registrysyncer.New(lggr, uni, uni.capReg.Address().String())
	require.NoError(t, err)

	hcr := newHomeChainReader(t, lp, cl, uni.cc.Address())
	launcher := New(
		ccipCapabilityVersion,
		ccipCapabilityLabelledName,
		p2pIDs[0],
		logger.TestLogger(t),
		hcr,
		&oracleCreatorPrints{
			t: t,
		},
		3*time.Second,
	)
	regSyncer.AddLauncher(launcher)

	require.NoError(t, launcher.Start(ctx))
	require.NoError(t, regSyncer.Start(ctx))
	t.Cleanup(func() { require.NoError(t, regSyncer.Close()) })
	t.Cleanup(func() { require.NoError(t, launcher.Close()) })

	chainAConf := setupConfigInfo(chainA, p2pIDs, fChainA, []byte("chainA"))
	chainBConf := setupConfigInfo(chainB, p2pIDs[1:], fChainB, []byte("chainB"))
	chainCConf := setupConfigInfo(chainC, p2pIDs[2:], fChainC, []byte("chainC"))
	inputConfig := []ccip_config.CCIPConfigTypesChainConfigInfo{
		chainAConf,
		chainBConf,
		chainCConf,
	}
	_, err = uni.cc.ApplyChainConfigUpdates(uni.transactor, nil, inputConfig)
	require.NoError(t, err)
	uni.backend.Commit()

	ccipCapabilityID, err := uni.capReg.GetHashedCapabilityId(nil, ccipCapabilityLabelledName, ccipCapabilityVersion)
	require.NoError(t, err)

	addDONToRegistry(
		t,
		uni.transactor,
		ccipCapabilityID,
		chainA,
		fChainA,
		uni.capReg,
		uni.backend,
		p2pIDs[1], // we're not bootstrapping
		p2pIDs,
	)

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
