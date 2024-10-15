package integrationhelpers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	ccipreader "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/require"

	configsevm "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const chainID = 1337

var CapabilityID = fmt.Sprintf("%s@%s", CcipCapabilityLabelledName, CcipCapabilityVersion)

func NewReader(
	t *testing.T,
	logPoller logpoller.LogPoller,
	headTracker logpoller.HeadTracker,
	client client.Client,
	address common.Address,
	chainReaderConfig evmrelaytypes.ChainReaderConfig,
) types.ContractReader {
	cr, err := evm.NewChainReaderService(testutils.Context(t), logger.TestLogger(t), logPoller, headTracker, client, chainReaderConfig)
	require.NoError(t, err)
	err = cr.Bind(testutils.Context(t), []types.BoundContract{
		{
			Address: address.String(),
			Name:    consts.ContractNameCCIPConfig,
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
	CCIPHome        *ccip_home.CCIPHome
	TestingT        *testing.T
	LogPoller       logpoller.LogPoller
	HeadTracker     logpoller.HeadTracker
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

	ccAddress, _, _, err := ccip_home.DeployCCIPHome(transactor, backend, crAddress)
	require.NoError(t, err)
	backend.Commit()

	cc, err := ccip_home.NewCCIPHome(ccAddress, backend)
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
	headTracker := headtracker.NewSimulatedHeadTracker(cl, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	if lpOpts.PollPeriod == 0 {
		lpOpts.PollPeriod = 1 * time.Hour
	}
	lp := logpoller.NewLogPoller(logpoller.NewORM(big.NewInt(chainID), db, lggr), cl, logger.NullLogger, headTracker, lpOpts)
	require.NoError(t, lp.Start(ctx))
	t.Cleanup(func() { require.NoError(t, lp.Close()) })

	hcr := NewHomeChainReader(t, lp, headTracker, cl, ccAddress)
	return TestUniverse{
		Transactor:      transactor,
		Backend:         backend,
		CapReg:          capReg,
		CCIPHome:        cc,
		TestingT:        t,
		SimClient:       cl,
		LogPoller:       lp,
		HeadTracker:     headTracker,
		HomeChainReader: hcr,
	}
}

func (t TestUniverse) NewContractReader(ctx context.Context, cfg []byte) (types.ContractReader, error) {
	var config evmrelaytypes.ChainReaderConfig
	err := json.Unmarshal(cfg, &config)
	require.NoError(t.TestingT, err)
	return evm.NewChainReaderService(ctx, logger.TestLogger(t.TestingT), t.LogPoller, t.HeadTracker, t.SimClient, config)
}

func P2pIDsFromInts(ints []int64) [][32]byte {
	var p2pIDs [][32]byte
	for _, i := range ints {
		p2pID := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(i)).PeerID()
		p2pIDs = append(p2pIDs, p2pID)
	}
	sort.Slice(p2pIDs, func(i, j int) bool {
		for k := 0; k < 32; k++ {
			if p2pIDs[i][k] < p2pIDs[j][k] {
				return true
			} else if p2pIDs[i][k] > p2pIDs[j][k] {
				return false
			}
		}
		return false
	})
	return p2pIDs
}

func (t *TestUniverse) AddCapability(p2pIDs [][32]byte) {
	_, err := t.CapReg.AddCapabilities(t.Transactor, []kcr.CapabilitiesRegistryCapability{
		{
			LabelledName:          CcipCapabilityLabelledName,
			Version:               CcipCapabilityVersion,
			CapabilityType:        0,
			ResponseType:          0,
			ConfigurationContract: t.CCIPHome.Address(),
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
				EncryptionPublicKey: testutils.Random32Byte(),
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

func NewHomeChainReader(t *testing.T, logPoller logpoller.LogPoller, headTracker logpoller.HeadTracker, client client.Client, ccAddress common.Address) ccipreader.HomeChain {
	cr := NewReader(t, logPoller, headTracker, client, ccAddress, configsevm.HomeChainReaderConfigRaw)

	hcr := ccipreader.NewHomeChainReader(cr, logger.TestLogger(t), 50*time.Millisecond, types.BoundContract{
		Address: ccAddress.String(),
		Name:    consts.ContractNameCCIPConfig,
	})
	require.NoError(t, hcr.Start(testutils.Context(t)))
	t.Cleanup(func() { require.NoError(t, hcr.Close()) })

	return hcr
}

func (t *TestUniverse) AddDONToRegistry(
	ccipCapabilityID [32]byte,
	chainSelector uint64,
	f uint8,
	p2pIDs [][32]byte,
) {
	tabi, err := ccip_home.CCIPHomeMetaData.GetAbi()
	require.NoError(t.TestingT, err)

	var nodes []ccip_home.CCIPHomeOCR3Node

	for i := range p2pIDs {
		nodes = append(nodes, ccip_home.CCIPHomeOCR3Node{
			P2pId:          p2pIDs[i],
			SignerKey:      testutils.NewAddress().Bytes(),
			TransmitterKey: testutils.NewAddress().Bytes(),
		})
	}

	// find the max don id, the next DON id will be max + 1.
	iter, err := t.CapReg.FilterConfigSet(nil, nil)
	require.NoError(t.TestingT, err)
	var maxDonID uint32
	for iter.Next() {
		if iter.Event.DonId > maxDonID {
			maxDonID = iter.Event.DonId
		}
	}

	donID := maxDonID + 1

	for _, pluginType := range []cctypes.PluginType{cctypes.PluginTypeCCIPCommit, cctypes.PluginTypeCCIPExec} {
		ocr3Config := ccip_home.CCIPHomeOCR3Config{
			PluginType:            uint8(pluginType),
			ChainSelector:         chainSelector,
			FRoleDON:              f,
			OffchainConfigVersion: 30,
			OfframpAddress:        testutils.NewAddress().Bytes(),
			RmnHomeAddress:        testutils.NewAddress().Bytes(),
			Nodes:                 nodes,
			OffchainConfig:        []byte("offchain config"),
		}
		encodedSetCandidateCall, err := tabi.Pack(
			"setCandidate",
			donID,
			ocr3Config.PluginType,
			ocr3Config,
			[32]byte{},
		)
		require.NoError(t.TestingT, err)
		// Create DON should be called only once, any subsequent calls should be updating DON
		if pluginType == cctypes.PluginTypeCCIPCommit {
			_, err = t.CapReg.AddDON(
				t.Transactor, p2pIDs, []kcr.CapabilitiesRegistryCapabilityConfiguration{
					{
						CapabilityId: ccipCapabilityID,
						Config:       encodedSetCandidateCall,
					},
				},
				false,
				false,
				f,
			)
		} else {
			_, err = t.CapReg.UpdateDON(
				t.Transactor, donID, p2pIDs, []kcr.CapabilitiesRegistryCapabilityConfiguration{
					{
						CapabilityId: ccipCapabilityID,
						Config:       encodedSetCandidateCall,
					},
				},
				false,
				f,
			)
		}

		require.NoError(t.TestingT, err)
		t.Backend.Commit()

		configs, err := t.CCIPHome.GetAllConfigs(nil, donID, uint8(pluginType))
		require.NoError(t.TestingT, err)
		require.Equal(t.TestingT, ocr3Config, configs.CandidateConfig.Config)

		// get the config digest of the candidate
		candidateDigest, err := t.CCIPHome.GetCandidateDigest(nil, donID, ocr3Config.PluginType)
		require.NoError(t.TestingT, err)
		encodedPromotionCall, err := tabi.Pack(
			"promoteCandidateAndRevokeActive",
			donID,
			ocr3Config.PluginType,
			candidateDigest,
			[32]byte{},
		)
		require.NoError(t.TestingT, err)

		_, err = t.CapReg.UpdateDON(
			t.Transactor, donID, p2pIDs, []kcr.CapabilitiesRegistryCapabilityConfiguration{
				{
					CapabilityId: ccipCapabilityID,
					Config:       encodedPromotionCall,
				},
			},
			false,
			f,
		)

		require.NoError(t.TestingT, err)
		t.Backend.Commit()

		configs, err = t.CCIPHome.GetAllConfigs(nil, donID, uint8(pluginType))
		require.NoError(t.TestingT, err)
		require.Equal(t.TestingT, ocr3Config, configs.ActiveConfig.Config)
	}
}

func SetupConfigInfo(chainSelector uint64, readers [][32]byte, fChain uint8, cfg []byte) ccip_home.CCIPHomeChainConfigArgs {
	return ccip_home.CCIPHomeChainConfigArgs{
		ChainSelector: chainSelector,
		ChainConfig: ccip_home.CCIPHomeChainConfig{
			Readers: readers,
			FChain:  fChain,
			Config:  cfg,
		},
	}
}
