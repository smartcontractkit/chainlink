package ccip_integration_tests

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

	ccipreader "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ocr3_config_encoder"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/services/ccipcapability/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const chainID = 1337

func NewReader(t *testing.T, logPoller logpoller.LogPoller, client client.Client, address common.Address, chainReaderConfig evmrelaytypes.ChainReaderConfig, contractName string) types.ContractReader {
	cr, err := evm.NewChainReaderService(testutils.Context(t), logger.TestLogger(t), logPoller, client, chainReaderConfig)
	require.NoError(t, err)
	err = cr.Bind(testutils.Context(t), []types.BoundContract{
		{
			Address: address.String(),
			Name:    contractName,
			Pending: false,
		},
	})
	require.NoError(t, err)
	require.NoError(t, cr.Start(testutils.Context(t)))
	for {
		if err := cr.Ready(); err == nil {
			break
		}
	}

	return cr
}

const (
	ChainA  uint64 = 1
	FChainA uint8  = 1

	ChainB  uint64 = 2
	FChainB uint8  = 2

	ChainC  uint64 = 3
	FChainC uint8  = 3

	CcipCapabilityLabelledName = "ccip"
	CcipCapabilityVersion      = "v1.0"
)

type TestUniverse struct {
	Transactor      *bind.TransactOpts
	Backend         *backends.SimulatedBackend
	CapReg          *kcr.CapabilitiesRegistry
	CcipCfg         *ccip_config.CCIPConfig
	TestingT        *testing.T
	LogPoller       logpoller.LogPoller
	SimClient       client.Client
	HomeChainReader ccipreader.HomeChain
}

func NewTestUniverse(ctx context.Context, t *testing.T, lggr logger.Logger) TestUniverse {
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

	db := pgtest.NewSqlxDB(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            0,
		BackfillBatchSize:        10,
		RpcBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	cl := client.NewSimulatedBackendClient(t, backend, big.NewInt(chainID))
	lp := logpoller.NewLogPoller(logpoller.NewORM(big.NewInt(chainID), db, lggr), cl, logger.NullLogger, lpOpts)
	require.NoError(t, lp.Start(ctx))
	t.Cleanup(func() { require.NoError(t, lp.Close()) })

	hcr := NewHomeChainReader(t, lp, cl, ccAddress)
	return TestUniverse{
		Transactor:      transactor,
		Backend:         backend,
		CapReg:          capReg,
		CcipCfg:         cc,
		TestingT:        t,
		SimClient:       cl,
		LogPoller:       lp,
		HomeChainReader: hcr,
	}
}

func (t TestUniverse) NewContractReader(ctx context.Context, cfg []byte) (types.ContractReader, error) {
	var config evmrelaytypes.ChainReaderConfig
	err := json.Unmarshal(cfg, &config)
	require.NoError(t.TestingT, err)
	return evm.NewChainReaderService(ctx, logger.TestLogger(t.TestingT), t.LogPoller, t.SimClient, config)
}

func P2pIDsFromInts(ints []int64) [][32]byte {
	var p2pIDs [][32]byte
	for _, i := range ints {
		p2pID := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(i)).PeerID()
		p2pIDs = append(p2pIDs, p2pID)
	}
	return p2pIDs
}

func (t *TestUniverse) AddCapability(p2pIDs [][32]byte) {
	_, err := t.CapReg.AddCapabilities(t.Transactor, []kcr.CapabilitiesRegistryCapability{
		{
			LabelledName:          CcipCapabilityLabelledName,
			Version:               CcipCapabilityVersion,
			CapabilityType:        0,
			ResponseType:          0,
			ConfigurationContract: t.CcipCfg.Address(),
		},
	})
	require.NoError(t.TestingT, err, "failed to add capability to registry")
	t.Backend.Commit()

	ccipCapabilityID, err := t.CapReg.GetHashedCapabilityId(nil, CcipCapabilityLabelledName, CcipCapabilityVersion)
	require.NoError(t.TestingT, err)

	for i := 0; i < len(p2pIDs); i++ {
		_, err = t.CapReg.AddNodeOperators(t.Transactor, []kcr.CapabilitiesRegistryNodeOperator{
			{
				Admin: t.Transactor.From,
				Name:  fmt.Sprintf("nop-%d", i),
			},
		})
		require.NoError(t.TestingT, err)
		t.Backend.Commit()

		// get the node operator id from the event
		it, err := t.CapReg.FilterNodeOperatorAdded(nil, nil, nil)
		require.NoError(t.TestingT, err)
		var nodeOperatorID uint32
		for it.Next() {
			if it.Event.Name == fmt.Sprintf("nop-%d", i) {
				nodeOperatorID = it.Event.NodeOperatorId
				break
			}
		}
		require.NotZero(t.TestingT, nodeOperatorID)

		_, err = t.CapReg.AddNodes(t.Transactor, []kcr.CapabilitiesRegistryNodeParams{
			{
				NodeOperatorId:      nodeOperatorID,
				Signer:              testutils.Random32Byte(),
				P2pId:               p2pIDs[i],
				HashedCapabilityIds: [][32]byte{ccipCapabilityID},
			},
		})
		require.NoError(t.TestingT, err)
		t.Backend.Commit()

		// verify that the node was added successfully
		nodeInfo, err := t.CapReg.GetNode(nil, p2pIDs[i])
		require.NoError(t.TestingT, err)

		require.Equal(t.TestingT, nodeOperatorID, nodeInfo.NodeOperatorId)
		require.Equal(t.TestingT, p2pIDs[i][:], nodeInfo.P2pId[:])
	}
}

func NewHomeChainReader(t *testing.T, logPoller logpoller.LogPoller, client client.Client, ccAddress common.Address) ccipreader.HomeChain {
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

	cr := NewReader(t, logPoller, client, ccAddress, cfg, "CCIPConfig")

	hcr := ccipreader.NewHomeChainReader(cr, logger.TestLogger(t), 500*time.Millisecond)
	require.NoError(t, hcr.Start(testutils.Context(t)))
	t.Cleanup(func() { require.NoError(t, hcr.Close()) })

	return hcr
}

func (t *TestUniverse) AddDONToRegistry(
	ccipCapabilityID [32]byte,
	chainSelector uint64,
	f uint8,
	bootstrapP2PID [32]byte,
	p2pIDs [][32]byte,
) {
	tabi, err := ocr3_config_encoder.IOCR3ConfigEncoderMetaData.GetAbi()
	require.NoError(t.TestingT, err)

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
	require.NoError(t.TestingT, err)

	// Trim first four bytes to remove function selector.
	encodedConfigs := encodedCall[4:]

	_, err = t.CapReg.AddDON(t.Transactor, p2pIDs, []kcr.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: ccipCapabilityID,
			Config:       encodedConfigs,
		},
	}, false, false, f)
	require.NoError(t.TestingT, err)
	t.Backend.Commit()
}

func SetupConfigInfo(chainSelector uint64, readers [][32]byte, fChain uint8, cfg []byte) ccip_config.CCIPConfigTypesChainConfigInfo {
	return ccip_config.CCIPConfigTypesChainConfigInfo{
		ChainSelector: chainSelector,
		ChainConfig: ccip_config.CCIPConfigTypesChainConfig{
			Readers: readers,
			FChain:  fChain,
			Config:  cfg,
		},
	}
}
