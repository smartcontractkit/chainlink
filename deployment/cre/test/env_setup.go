package test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	changeset2 "github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	envtest "github.com/smartcontractkit/chainlink/deployment/environment/test"
)

const (
	DONName           = "test-don"
	RegistryQualifier = "test-registry"
)

type EnvWrapperV2 struct {
	t *testing.T

	Env              *cldf.Environment
	RegistrySelector uint64
	RegistryAddress  common.Address
}

type donConfig struct {
	Name             string
	N                int
	F                int
	CapabilityConfig map[string]*pb.CapabilityConfig
	Labels           map[string]string
	RegistryChainSel uint64
}

func initEnv(t *testing.T, lggr logger.Logger) (uint64, *cldf.Environment) {
	chains := cldf_chain.NewBlockChainsFromSlice(memory.NewMemoryChainsEVM(t, 1, 1))
	registryChainSel := registryChain(chains.EVMChains())

	ds := datastore.NewMemoryDataStore()
	env := cldf.Environment{
		Logger:           lggr,
		GetContext:       t.Context,
		DataStore:        ds.Seal(),
		BlockChains:      chains,
		OperationsBundle: operations.NewBundle(t.Context, lggr, operations.NewMemoryReporter()),
	}

	deployCapRegChangeset := changeset2.DeployCapabilitiesRegistry{}
	changes := []changeset.ConfiguredChangeSet{
		changeset.Configure(
			cldf.CreateChangeSet(deployCapRegChangeset.Apply, deployCapRegChangeset.VerifyPreconditions),
			changeset2.DeployCapabilitiesRegistryInput{
				ChainSelector: registryChainSel,
				Qualifier:     RegistryQualifier,
			},
		),
	}

	env, _, err := changeset.ApplyChangesets(t, env, changes)
	require.NoError(t, err)
	require.NotNil(t, env)
	require.Len(t, env.BlockChains.EVMChains(), 1)

	return registryChainSel, &env
}

// SetupEnvV2 starts an environment with a single DON, 4 nodes and a capabilities registry v2 deployed and configured.
func SetupEnvV2(t *testing.T, useMCMS bool) *EnvWrapperV2 {
	t.Helper()

	lggr := logger.Test(t)

	registryChainSel, envInitiated := initEnv(t, lggr)
	lggr.Debug("Initialized environment", "registryChainSel", registryChainSel)

	n := 4
	donCfg := donConfig{
		Name:             DONName,
		N:                n,
		F:                (n-1)/3 + 1,
		RegistryChainSel: registryChainSel,
	}

	// Only need one DON
	don, env := setupViewOnlyNodeTest(t, registryChainSel, envInitiated.BlockChains.EVMChains(), donCfg)

	env.DataStore = envInitiated.DataStore

	registryAddrs := env.DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(registryChainSel),
		datastore.AddressRefByType("CapabilitiesRegistry"),
	)
	require.Len(t, registryAddrs, 1)

	chainID, err := chain_selectors.GetChainIDFromSelector(registryChainSel)
	require.NoError(t, err)

	registryChainDetails, err := chain_selectors.GetChainDetailsByChainIDAndFamily(chainID, chain_selectors.FamilyEVM)
	require.NoError(t, err)

	donNodes, err := don.AllNodes()
	require.NoError(t, err)

	nodesP2PIDs := make([]string, 0, len(donNodes))
	nodesP2PIDsBytes := make([][32]byte, 0, len(donNodes))

	var nodes []changeset2.CapabilitiesRegistryNodeParams
	for _, n := range donNodes {
		p2pID := n.PeerID.String()
		ocrConfig, ok := n.OCRConfigs[registryChainDetails]
		require.True(t, ok, "node %s does not have OCR config for registry chain %d", n.Name, registryChainSel)

		nodesP2PIDs = append(nodesP2PIDs, p2pID)
		nodesP2PIDsBytes = append(nodesP2PIDsBytes, n.PeerID)

		nodes = append(nodes, changeset2.CapabilitiesRegistryNodeParams{
			NodeOperatorID:      1,
			P2pID:               p2pID,
			CsaKey:              n.CSA,
			EncryptionPublicKey: n.WorkflowKey,
			Signer:              hex.EncodeToString(ocrConfig.OnchainPublicKey),
			CapabilityIDs: []string{
				"test-capability@1.0.0",
			},
		})
	}

	configCapRegChangeset := changeset2.ConfigureCapabilitiesRegistry{}
	changes := []changeset.ConfiguredChangeSet{
		changeset.Configure(
			cldf.CreateChangeSet(configCapRegChangeset.Apply, configCapRegChangeset.VerifyPreconditions),
			changeset2.ConfigureCapabilitiesRegistryInput{
				ChainSelector:               registryChainSel,
				CapabilitiesRegistryAddress: registryAddrs[0].Address,
				UseMCMS:                     useMCMS,
				Nops: []changeset2.CapabilitiesRegistryNodeOperator{
					{
						Name:  "Operator 1",
						Admin: common.HexToAddress("0x01"),
					},
				},
				Nodes: nodes,
				Capabilities: []changeset2.CapabilitiesRegistryCapability{
					{
						CapabilityID: "test-capability@1.0.0",
						Metadata:     map[string]interface{}{"capabilityType": 2},
					},
				},
				DONs: []changeset2.CapabilitiesRegistryNewDONParams{
					{
						Name:        donCfg.Name,
						F:           uint8(donCfg.F), //nolint:gosec // disable G115
						Nodes:       nodesP2PIDs,
						DonFamilies: []string{"test-family"},
						Config:      map[string]interface{}{"consensus": "basic", "timeout": "30s"},
						CapabilityConfigurations: []changeset2.CapabilitiesRegistryCapabilityConfiguration{
							{
								CapabilityID: "test-capability@1.0.0",
							},
						},
						IsPublic:         true,
						AcceptsWorkflows: true,
					},
				},
			},
		),
	}

	env, _, err = changeset.ApplyChangesets(t, env, changes)
	require.NoError(t, err)
	require.NotNil(t, env)

	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(common.HexToAddress(registryAddrs[0].Address), env.BlockChains.EVMChains()[registryChainSel].Client)
	require.NoError(t, err)
	require.NotNil(t, capReg)

	gotNodes, err := capReg.GetNodesByP2PIds(nil, nodesP2PIDsBytes)
	require.NoError(t, err)
	require.Len(t, gotNodes, len(don.GetP2PIDs()))
	require.Len(t, gotNodes, donCfg.N)
	for _, n := range gotNodes {
		require.Equal(t, "test-capability@1.0.0", n.CapabilityIds[0])
	}

	gotDON, err := capReg.GetDONByName(nil, donCfg.Name)
	require.NoError(t, err)
	require.Len(t, gotDON.NodeP2PIds, len(nodesP2PIDsBytes))

	// Sort both slices before comparison
	sort.Slice(gotDON.NodeP2PIds, func(i, j int) bool {
		return bytes.Compare(gotDON.NodeP2PIds[i][:], gotDON.NodeP2PIds[j][:]) < 0
	})
	sortedNodesP2PIDsBytes := make([][32]byte, len(nodesP2PIDsBytes))
	copy(sortedNodesP2PIDsBytes, nodesP2PIDsBytes)
	sort.Slice(sortedNodesP2PIDsBytes, func(i, j int) bool {
		return bytes.Compare(sortedNodesP2PIDsBytes[i][:], sortedNodesP2PIDsBytes[j][:]) < 0
	})
	for i, id := range gotDON.NodeP2PIds {
		require.Equal(t, sortedNodesP2PIDsBytes[i], id)
	}

	return &EnvWrapperV2{
		t:                t,
		Env:              &env,
		RegistrySelector: registryChainSel,
		RegistryAddress:  common.HexToAddress(registryAddrs[0].Address),
	}
}

func setupViewOnlyNodeTest(t *testing.T, registryChainSel uint64, chains map[uint64]cldf_evm.Chain, donCfg donConfig) (*viewOnlyDon, cldf.Environment) {
	var (
		don      *viewOnlyDon
		nodesCfg []envtest.NodeConfig
	)

	for i := 0; i < donCfg.N; i++ {
		labels := map[string]string{
			"don": donCfg.Name,
		}
		if donCfg.Labels != nil {
			for k, v := range donCfg.Labels {
				labels[k] = v
			}
		}
		nCfg := envtest.NodeConfig{
			ChainSelectors: []uint64{registryChainSel},
			Name:           fmt.Sprintf("%s-%d", donCfg.Name, i),
			Labels:         labels,
		}
		nodesCfg = append(nodesCfg, nCfg)
	}

	n := envtest.NewNodes(t, nodesCfg)
	require.Len(t, n, donCfg.N)

	don = newViewOnlyDon(donCfg.Name, n)

	nodes := make(deployment.Nodes, 0, don.N())
	for _, v := range don.m {
		nodes = append(nodes, *v)
	}

	blockChains := map[uint64]cldf_chain.BlockChain{}
	for sel, c := range chains {
		blockChains[sel] = c
	}

	env := cldf.NewEnvironment(
		"view only nodes",
		logger.Test(t),
		cldf.NewMemoryAddressBook(),
		datastore.NewMemoryDataStore().Seal(),
		nodes.IDs(),
		envtest.NewJDService(nodes),
		t.Context,
		cldf.XXXGenerateTestOCRSecrets(),
		cldf_chain.NewBlockChains(blockChains),
	)

	return don, *env
}

func registryChain(chains map[uint64]cldf_evm.Chain) uint64 {
	var registryChainSel uint64 = math.MaxUint64
	for sel := range chains {
		if sel < registryChainSel {
			registryChainSel = sel
		}
	}
	return registryChainSel
}
