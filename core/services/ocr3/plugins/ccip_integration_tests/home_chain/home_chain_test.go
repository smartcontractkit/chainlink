package home_chain

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/onsi/gomega"

	libocrtypes "github.com/smartcontractkit/libocr/ragep2p/types"

	ccipreader "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	capcfg "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	helpers "github.com/smartcontractkit/chainlink/v2/core/services/ocr3/plugins/ccip_integration_tests"

	"github.com/stretchr/testify/require"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const (
	chainA  uint64 = 1
	fChainA uint8  = 1

	chainB  uint64 = 2
	fChainB uint8  = 2

	chainC  uint64 = 3
	fChainC uint8  = 3
)

func TestHomeChainReader(t *testing.T) {
	// Initialize chainReader
	cfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			"CCIPConfig": {
				ContractABI: capcfg.CCIPConfigMetaData.ABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					"getAllChainConfigs": {
						ChainSpecificName: "getAllChainConfigs",
					},
				},
			},
		},
	}
	//============================Setup Backend===================================
	transactor := testutils.MustNewSimTransactor(t)
	backend := backends.NewSimulatedBackend(core.GenesisAlloc{
		transactor.From: {Balance: assets.Ether(1000).ToInt()},
	}, 30e6)
	//==============================Setup Contracts - Add capabilities=================================
	capRegAddress, capRegContract, err := prepareCapabilityRegistry(t, backend, transactor)
	require.NoError(t, err)
	capConfAddress, capConfContract, err := prepareCCIPCapabilityConfig(t, backend, transactor, capRegAddress)
	require.NoError(t, err)
	p2pIDS := addCapabilities(t, backend, transactor, capRegContract, capConfAddress)
	//==============================Apply configs to Capability Contract=================================
	chainAConf := setupConfigInfo(chainA, p2pIDS, fChainA, []byte("chainA"))
	chainBConf := setupConfigInfo(chainB, p2pIDS[1:], fChainB, []byte("chainB"))
	chainCConf := setupConfigInfo(chainC, p2pIDS[2:], fChainC, []byte("chainC"))
	inputConfig := []capcfg.CCIPConfigChainConfigInfo{
		chainAConf,
		chainBConf,
		chainCConf,
	}
	_, err = capConfContract.ApplyChainConfigUpdates(transactor, nil, inputConfig)
	require.NoError(t, err)
	backend.Commit()
	//================================Setup HomeChainReader===============================
	ctx := testutils.Context(t)
	testData := helpers.SetupReaderTestData(ctx, t, backend, capConfAddress, cfg, "CCIPConfig")
	chainReader := testData.ChainReader
	logPoller := testData.LogPoller
	require.NoError(t, err)
	pollDuration := 5 * time.Millisecond
	homeChain := ccipreader.NewHomeChainReader(chainReader, logger.TestLogger(t), pollDuration)
	require.NoError(t, homeChain.Start(ctx))

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
	_, err = capConfContract.ApplyChainConfigUpdates(transactor, []uint64{chainC}, nil)
	require.NoError(t, err)
	backend.Commit()
	time.Sleep(pollDuration * 5) // Wait for the chain reader to update
	configs, err = homeChain.GetAllChainConfigs()
	require.NoError(t, err)
	delete(expectedChainConfigs, cciptypes.ChainSelector(chainC))
	require.Equal(t, expectedChainConfigs, configs)
	//================================Close HomeChain Reader===============================
	require.NoError(t, homeChain.Close())
	require.NoError(t, logPoller.Close())
	require.NoError(t, chainReader.Close())
	t.Logf("homchain reader successfully closed")
}

func toPeerIDs(readers [][32]byte) mapset.Set[libocrtypes.PeerID] {
	peerIDs := mapset.NewSet[libocrtypes.PeerID]()
	for _, r := range readers {
		peerIDs.Add(r)
	}
	return peerIDs
}

func setupConfigInfo(chainSelector uint64, readers [][32]byte, fChain uint8, cfg []byte) capcfg.CCIPConfigChainConfigInfo {
	return capcfg.CCIPConfigChainConfigInfo{
		ChainSelector: chainSelector,
		ChainConfig: capcfg.CCIPConfigChainConfig{
			Readers: readers,
			FChain:  fChain,
			Config:  cfg,
		},
	}
}

func prepareCCIPCapabilityConfig(t *testing.T, backend *backends.SimulatedBackend, transactor *bind.TransactOpts, capRegAddress common.Address) (common.Address, *capcfg.CCIPConfig, error) {
	ccAddress, _, _, err := capcfg.DeployCCIPConfig(transactor, backend, capRegAddress)
	require.NoError(t, err)
	backend.Commit()

	contract, err := capcfg.NewCCIPConfig(ccAddress, backend)
	require.NoError(t, err)
	backend.Commit()

	return ccAddress, contract, nil
}

func prepareCapabilityRegistry(t *testing.T, backend *backends.SimulatedBackend, transactor *bind.TransactOpts) (common.Address, *capabilities_registry.CapabilitiesRegistry, error) {
	crAddress, _, _, err := capabilities_registry.DeployCapabilitiesRegistry(transactor, backend)
	require.NoError(t, err)
	backend.Commit()

	capReg, err := capabilities_registry.NewCapabilitiesRegistry(crAddress, backend)
	require.NoError(t, err)
	backend.Commit()

	return crAddress, capReg, nil
}

func addCapabilities(
	t *testing.T,
	backend *backends.SimulatedBackend,
	transactor *bind.TransactOpts,
	capReg *capabilities_registry.CapabilitiesRegistry,
	capConfAddress common.Address) [][32]byte {
	// add the CCIP capability to the registry
	_, err := capReg.AddCapabilities(transactor, []capabilities_registry.CapabilitiesRegistryCapability{
		{
			LabelledName:          "ccip",
			Version:               "v1.0",
			CapabilityType:        0,
			ResponseType:          0,
			ConfigurationContract: capConfAddress,
		},
	})
	require.NoError(t, err, "failed to add capability to registry")
	backend.Commit()

	ccipCapabilityID, err := capReg.GetHashedCapabilityId(nil, "ccip", "v1.0")
	require.NoError(t, err)

	// Add the p2p ids of the ccip nodes
	var p2pIDs [][32]byte
	for i := 0; i < 4; i++ {
		p2pID := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(i + 1))).PeerID()
		p2pIDs = append(p2pIDs, p2pID)
		_, err = capReg.AddNodeOperators(transactor, []capabilities_registry.CapabilitiesRegistryNodeOperator{
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

		_, err = capReg.AddNodes(transactor, []capabilities_registry.CapabilitiesRegistryNodeParams{
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
