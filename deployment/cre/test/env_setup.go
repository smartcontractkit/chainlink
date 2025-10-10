package test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"maps"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	changeset2 "github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	envtest "github.com/smartcontractkit/chainlink/deployment/environment/test"
)

const (
	DONName           = "test-don"
	RegistryQualifier = "test-registry"
)

type EnvWrapperV2 struct {
	t *testing.T

	TestJD *envtest.JDNodeService

	Env              *cldf.Environment
	RegistrySelector uint64
	RegistryAddress  common.Address
	AptosSelector    uint64
}

type donConfig struct {
	Name             string
	N                int
	F                int
	CapabilityConfig map[string]*pb.CapabilityConfig
	Labels           map[string]string
	RegistryChainSel uint64
}

// TODO CRE-999; aptos can be made optional
func initEnv(t *testing.T, lggr logger.Logger) (registryChainSel, aptosChainSel uint64, env *cldf.Environment) {
	evmChains := memory.NewMemoryChainsEVM(t, 1, 1)
	chains := cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{
		evmChains[0],
	})
	registryChainSel = evmChains[0].ChainSelector()

	ds := datastore.NewMemoryDataStore()
	localEnv := cldf.Environment{
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

	localEnv, _, err := changeset.ApplyChangesets(t, localEnv, changes)
	require.NoError(t, err)

	env = &localEnv
	require.NotNil(t, env)
	require.Len(t, env.BlockChains.EVMChains(), 1)

	// by inspection, the only chain that is needed is evm, but some callers expect aptos keys and therefore an aptos selector to use for generating the keys
	return registryChainSel, chain_selectors.APTOS_LOCALNET.Selector, env
}

// SetupEnvV2 starts an environment with a single DON, 4 nodes and a capabilities registry v2 deployed and configured.
func SetupEnvV2(t *testing.T, useMCMS bool) *EnvWrapperV2 {
	t.Helper()

	lggr := logger.Test(t)

	registryChainSel, aptosChainSel, envInitiated := initEnv(t, lggr)
	lggr.Debug("Initialized environment", "registryChainSel", registryChainSel)

	n := 4
	donCfg := donConfig{
		Name:             DONName,
		N:                n,
		F:                (n-1)/3 + 1,
		RegistryChainSel: registryChainSel,
	}

	// Only need one DON
	don, env, jd := setupViewOnlyNodeTest(t, registryChainSel, aptosChainSel, envInitiated.BlockChains, donCfg)

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

	var mcmsConfig *ocr3.MCMSConfig
	if useMCMS {
		mcmsConfig = &ocr3.MCMSConfig{
			MinDuration: 10 * time.Second,
		}
	}

	configCapRegChangeset := changeset2.ConfigureCapabilitiesRegistry{}
	changes := []changeset.ConfiguredChangeSet{
		changeset.Configure(
			cldf.CreateChangeSet(configCapRegChangeset.Apply, configCapRegChangeset.VerifyPreconditions),
			changeset2.ConfigureCapabilitiesRegistryInput{
				ChainSelector:               registryChainSel,
				CapabilitiesRegistryAddress: registryAddrs[0].Address,
				MCMSConfig:                  mcmsConfig,
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
						Metadata:     map[string]any{"capabilityType": 2},
					},
				},
				DONs: []changeset2.CapabilitiesRegistryNewDONParams{
					{
						Name:        donCfg.Name,
						F:           uint8(donCfg.F), //nolint:gosec // disable G115
						Nodes:       nodesP2PIDs,
						DonFamilies: []string{"test-family"},
						Config:      map[string]any{"consensus": "basic", "timeout": "30s"},
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
	require.Len(t, gotNodes, donCfg.N+1) // +1 for bootstrap
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
		TestJD:           jd,
		Env:              &env,
		AptosSelector:    aptosChainSel,
		RegistrySelector: registryChainSel,
		RegistryAddress:  common.HexToAddress(registryAddrs[0].Address),
	}
}

func setupViewOnlyNodeTest(t *testing.T, registryChainSel, aptosChainSel uint64, chains cldf_chain.BlockChains, donCfg donConfig) (*viewOnlyDon, cldf.Environment, *envtest.JDNodeService) {
	var (
		don      *viewOnlyDon
		nodesCfg []envtest.NodeConfig
	)

	for i := 0; i < donCfg.N; i++ {
		labels := map[string]string{
			"don-" + donCfg.Name: donCfg.Name,
			"environment":        "test",
			"product":            "cre",
			"type":               "plugin",
		}
		if donCfg.Labels != nil {
			maps.Copy(labels, donCfg.Labels)
		}

		nCfg := envtest.NodeConfig{
			ChainSelectors: []uint64{registryChainSel, aptosChainSel},
			Name:           fmt.Sprintf("%s-%d", donCfg.Name, i),
			Labels:         labels,
		}
		nodesCfg = append(nodesCfg, nCfg)
	}

	btLabels := map[string]string{
		"don-" + donCfg.Name: donCfg.Name,
		"environment":        "test",
		"product":            "cre",
		"type":               "bootstrap",
	}
	if donCfg.Labels != nil {
		maps.Copy(btLabels, donCfg.Labels)
	}
	nodesCfg = append(nodesCfg, envtest.NodeConfig{
		ChainSelectors: []uint64{registryChainSel, aptosChainSel},
		Name:           donCfg.Name + "-bootstrap",
		Labels:         btLabels,
	})

	n := envtest.NewNodes(t, nodesCfg)
	require.Len(t, n, donCfg.N+1) // +1 for bootstrap

	don = newViewOnlyDon(donCfg.Name, n)

	nodes := make(deployment.Nodes, 0, don.N())
	for _, v := range don.m {
		nodes = append(nodes, *v)
	}

	blockChains := map[uint64]cldf_chain.BlockChain{}
	for sel, c := range chains.EVMChains() {
		blockChains[sel] = c
	}
	for sel, c := range chains.AptosChains() {
		blockChains[sel] = c
	}

	jd := envtest.NewJDService(nodes)
	env := cldf.NewEnvironment(
		"test",
		logger.Test(t),
		cldf.NewMemoryAddressBook(),
		datastore.NewMemoryDataStore().Seal(),
		nodes.IDs(),
		jd,
		t.Context,
		focr.XXXGenerateTestOCRSecrets(),
		cldf_chain.NewBlockChains(blockChains),
	)

	return don, *env, jd
}
