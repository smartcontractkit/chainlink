package test

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

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

type EnvWrapperV2 struct {
	t *testing.T

	Env              *cldf.Environment
	RegistrySelector uint64
}

type DonConfig struct {
	Name             string // required, must be unique across all dons
	N                int
	F                *int                            // if nil, defaults to floor(N-1/3) + 1
	CapabilityConfig map[string]*pb.CapabilityConfig // optional DON specific configuration for the given capability (capabilityID => config)
	Labels           map[string]string               // optional
	RegistryChainSel uint64                          // require, must be the same for all dons
	ChainSelectors   []uint64                        // optional chains
}

func (c DonConfig) Validate() error {
	if c.N < 4 {
		return errors.New("N must be at least 4")
	}
	return nil
}

type EnvWrapperConfigV2 struct {
	DonConfig
	NumChains int

	UseMCMS bool

	// if true, use in-memory nodes for testing
	// if false, view only nodes will be used
	useInMemoryNodes bool
}

func (c EnvWrapperConfigV2) Validate() error {
	if err := c.DonConfig.Validate(); err != nil {
		return err
	}

	if c.NumChains < 1 {
		return errors.New("NumChains must be at least 1")
	}

	return nil
}

func initEnv(t *testing.T, lggr logger.Logger, nChains int) (uint64, *cldf.Environment) {
	chains := cldf_chain.NewBlockChainsFromSlice(memory.NewMemoryChainsEVM(t, nChains, 1))
	registryChainSel := registryChain(chains.EVMChains())

	ds := datastore.NewMemoryDataStore()
	env := cldf.Environment{
		Logger:            lggr,
		GetContext:        t.Context,
		ExistingAddresses: cldf.NewMemoryAddressBook(), // this hasn't been wiped out yet, so we have to at least instantiate it
		DataStore:         ds.Seal(),
		BlockChains:       chains,
		OperationsBundle:  operations.NewBundle(t.Context, lggr, operations.NewMemoryReporter()),
	}

	deployCapRegChangeset := changeset2.DeployCapabilitiesRegistry{}
	changes := []changeset.ConfiguredChangeSet{
		changeset.Configure(
			cldf.CreateChangeSet(deployCapRegChangeset.Apply, deployCapRegChangeset.VerifyPreconditions),
			changeset2.DeployCapabilitiesRegistryInput{
				ChainSelector: registryChainSel,
				Qualifier:     "test-registry",
			},
		),
	}

	env, _, err := changeset.ApplyChangesets(t, env, changes)
	require.NoError(t, err)
	require.NotNil(t, env)
	require.Len(t, env.BlockChains.EVMChains(), nChains)

	validateInitialChainState(t, env, registryChainSel)

	return registryChainSel, &env
}

func SetupEnvV2(t *testing.T, cfg EnvWrapperConfigV2) *EnvWrapperV2 {
	t.Helper()

	require.NoError(t, cfg.Validate())

	lggr := logger.Test(t)

	registryChainSel, envInitiated := initEnv(t, lggr, cfg.NumChains)
	lggr.Debug("Initialized environment", "registryChainSel", registryChainSel)

	var (
		dons testDons
		env  cldf.Environment
	)

	if cfg.useInMemoryNodes {
		dons, env = setupMemoryNodeTest(t, registryChainSel, envInitiated.BlockChains, cfg)
	} else {
		dons, env = setupViewOnlyNodeTest(t, registryChainSel, envInitiated.BlockChains.EVMChains(), cfg)
	}

	err := env.ExistingAddresses.Merge(envInitiated.ExistingAddresses)
	require.NoError(t, err)
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

	donNodes, err := dons.AllNodes()
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
			Signer:              string(ocrConfig.OnchainPublicKey),
			CapabilityIDs: []string{
				"offchain_reporting@1.0.0",
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
				UseMCMS:                     cfg.UseMCMS,
				Nops: []changeset2.CapabilitiesRegistryNodeOperator{
					{
						Name:  "Operator 1",
						Admin: common.HexToAddress("0x01"),
					},
				},
				Nodes: nodes,
				Capabilities: []changeset2.CapabilitiesRegistryCapability{
					{
						CapabilityID: "offchain_reporting@1.0.0",
						Metadata:     map[string]interface{}{"capabilityType": 2},
					},
				},
				DONs: []changeset2.CapabilitiesRegistryNewDONParams{
					{
						Name:        cfg.DonConfig.Name,
						F:           uint8(*cfg.DonConfig.F),
						Nodes:       nodesP2PIDs,
						DonFamilies: []string{"test-family"},
						Config:      map[string]interface{}{"consensus": "basic", "timeout": "30s"},
						CapabilityConfigurations: []changeset2.CapabilitiesRegistryCapabilityConfiguration{
							{
								CapabilityID: "offchain_reporting@1.0.0",
							},
						},
						IsPublic:         false,
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
	require.Len(t, gotNodes, len(dons.P2PIDs()))
	require.Equal(t, gotNodes, cfg.DonConfig.N)
	for _, n := range gotNodes {
		require.Equal(t, "offchain_reporting@1.0.0", n.CapabilityIds[0])
	}

	gotDON, err := capReg.GetDONByName(nil, cfg.DonConfig.Name)
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
		Env:              nil,
		RegistrySelector: registryChainSel,
	}
}

func setupViewOnlyNodeTest(t *testing.T, registryChainSel uint64, chains map[uint64]cldf_evm.Chain, c EnvWrapperConfigV2) (testDons, cldf.Environment) {
	// now that we have the initial contracts deployed, we can configure the nodes with the addresses
	dons := newViewOnlyDons()
	for _, donCfg := range []DonConfig{c.DonConfig} {
		require.NoError(t, donCfg.Validate())

		ncfg := make([]envtest.NodeConfig, 0, len(donCfg.ChainSelectors))
		for i := 0; i < donCfg.N; i++ {
			labels := map[string]string{
				"don": donCfg.Name,
			}
			if donCfg.Labels != nil {
				for k, v := range donCfg.Labels {
					labels[k] = v
				}
			}
			cfg := envtest.NodeConfig{
				ChainSelectors: []uint64{registryChainSel},
				Name:           fmt.Sprintf("%s-%d", donCfg.Name, i),
				Labels:         labels,
			}
			cfg.ChainSelectors = append(cfg.ChainSelectors, donCfg.ChainSelectors...)
			ncfg = append(ncfg, cfg)
		}
		n := envtest.NewNodes(t, ncfg)
		require.Len(t, n, donCfg.N)
		dons.Put(newViewOnlyDon(donCfg.Name, n))
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
		dons.NodeList().IDs(),
		envtest.NewJDService(dons.NodeList()),
		t.Context,
		cldf.XXXGenerateTestOCRSecrets(),
		cldf_chain.NewBlockChains(blockChains),
	)

	return dons, *env
}

func setupMemoryNodeTest(
	t *testing.T, registryChainSel uint64, blockchains cldf_chain.BlockChains, c EnvWrapperConfigV2,
) (testDons, cldf.Environment) {
	lggr := logger.Test(t)
	crConfig := deployment.CapabilityRegistryConfig{
		EVMChainID: registryChainSel,
		Contract:   [20]byte{},
	}

	wfChains := map[uint64]cldf_evm.Chain{}
	wfChains[registryChainSel] = blockchains.EVMChains()[registryChainSel]
	wfConf := memory.NewNodesConfig{
		LogLevel:       zapcore.InfoLevel,
		BlockChains:    blockchains,
		NumNodes:       c.DonConfig.N,
		NumBootstraps:  0,
		RegistryConfig: crConfig,
		CustomDBSetup:  nil,
	}
	wfNodes := memory.NewNodes(t, wfConf)
	require.Len(t, wfNodes, c.DonConfig.N)

	dons := newMemoryDons()
	dons.Put(newMemoryDon(c.DonConfig.Name, wfNodes))

	env := memory.NewMemoryEnvironmentFromChainsNodes(
		t.Context, lggr, blockchains, dons.AllNodesForJD(),
	)
	return dons, env
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

// validateInitialChainState checks that the initial chain state
// has the expected contracts deployed
func validateInitialChainState(t *testing.T, env cldf.Environment, registryChainSel uint64) {
	t.Helper()

	ds := env.DataStore
	ds.Addresses().Filter()
	// all contracts on registry chain
	registryChainAddrs := ds.Addresses().Filter(datastore.AddressRefByChainSelector(registryChainSel))
	require.Len(t, registryChainAddrs, 1) // registry
	// only forwarder on non-home chain
	for sel := range env.BlockChains.EVMChains() {
		chainAddrs := ds.Addresses().Filter(datastore.AddressRefByChainSelector(sel))
		if sel != registryChainSel {
			require.Len(t, chainAddrs, 0)
		} else {
			require.Len(t, chainAddrs, 1)
		}
	}
}
