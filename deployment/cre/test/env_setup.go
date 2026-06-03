package test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"maps"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmschangesets "github.com/smartcontractkit/cld-changesets/legacy/mcms/changesets"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	cldftesthelpers "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils/testhelpers"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/onchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	changeset2 "github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	envtest "github.com/smartcontractkit/chainlink/deployment/environment/test"
	changeset3 "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

const (
	DONName           = "test-don"
	RegistryQualifier = "test-registry"
	TestCapabilityID  = "test-capability@1.0.0"
	Zone              = "test-zone-1"
	TotalNodes        = 4
	// EnvironmentName is the name of the environment used in the test
	//
	// This is set by the runtime loader, but the constant is not exposed so we define it here.
	//
	// This will be fixed in a future release of the chainlink-deployments-framework.
	EnvironmentName = "test_environment"
)

var (
	DefaultRegistrySelector uint64 = chain_selectors.TEST_90000001.Selector
	DefaultAptosSelector    uint64 = chain_selectors.APTOS_LOCALNET.Selector
)

type Harness struct {
	t *testing.T

	Runtime *runtime.Runtime
	Don     *viewOnlyDon

	TestJD *envtest.JDNodeService

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

// HarnessConfig is used to configure the initialization of the test harness.
type HarnessConfig struct {
	WithMCMS      bool
	DatastoreSeed *datastore.MemoryDataStore
}

// GetDatastoreSeed returns the datastore that will be seeded into the runtime
// environment.
//
// Defaults to an empty memory data store.
func (cfg *HarnessConfig) GetDatastoreSeed() datastore.DataStore {
	if cfg.DatastoreSeed == nil {
		return datastore.NewMemoryDataStore().Seal()
	}

	return cfg.DatastoreSeed.Seal()
}

// HarnessOpt is used to configure the initialization of the test harness.
type HarnessOpt func(cfg *HarnessConfig)

// WithMCMS configures the test harness to use MCMS.
func WithMCMS() HarnessOpt {
	return func(cfg *HarnessConfig) {
		cfg.WithMCMS = true
	}
}

// WithDatastore configures a custom datastore for the test harness.
func WithDatastore(ds *datastore.MemoryDataStore) HarnessOpt {
	return func(cfg *HarnessConfig) {
		cfg.DatastoreSeed = ds
	}
}

// initHarness sets up the runtime and environment for the test harness.
//
// TODO CRE-999; aptos can be made optional
func initHarness(t *testing.T, lggr logger.Logger, cfg *HarnessConfig) *Harness {
	var (
		registryChainSel = chain_selectors.TEST_90000001.Selector
		// by inspection, the only chain that is needed is evm, but some callers
		// expect aptos keys and therefore an aptos selector to use for generating
		// the keys
		aptosChainSel = chain_selectors.APTOS_LOCALNET.Selector

		donCfg = donConfig{
			Name:             DONName,
			N:                TotalNodes,
			F:                (TotalNodes-1)/3 + 1,
			RegistryChainSel: registryChainSel,
		}
	)

	// Setup the view only DON. Only need one DON
	don := newViewOnlyNodes(t, registryChainSel, aptosChainSel, donCfg)

	// Setup JD service
	jd := envtest.NewJDService(don.Nodes())

	// Setup the runtime environment
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulatedWithConfig(t, []uint64{registryChainSel}, onchain.EVMSimLoaderConfig{
			NumAdditionalAccounts: 1,
		}),
		environment.WithLogger(lggr),
		environment.WithOffchainClient(jd),
		environment.WithNodeIDs(don.Nodes().IDs()),
		environment.WithDatastore(cfg.GetDatastoreSeed()),
	))
	require.NoError(t, err)

	return &Harness{
		t:                t,
		Runtime:          rt,
		TestJD:           jd,
		Don:              don,
		AptosSelector:    aptosChainSel,
		RegistrySelector: registryChainSel,
	}
}

// NewTestHarness starts a runtime with a single DON, 4 nodes and a capabilities registry v2
// deployed and configured.
func NewTestHarness(t *testing.T, opts ...HarnessOpt) *Harness {
	t.Helper()

	cfg := &HarnessConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	lggr := logger.Test(t)
	h := initHarness(t, lggr, cfg)

	t.Log("Initialized runtime", "registryChainSel", h.RegistrySelector)

	// Deploy the capabilities registry
	deployCapabilitiesRegistry(t, h)

	registryAddrs := h.Runtime.Environment().DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(h.RegistrySelector),
		datastore.AddressRefByType("CapabilitiesRegistry"),
	)
	require.Len(t, registryAddrs, 1)
	registryAddress := registryAddrs[0].Address
	h.RegistryAddress = common.HexToAddress(registryAddress) // Inject the registry address into the harness for convenience

	// Configure the capabilities registry
	configureCapabilitiesRegistry(t, h)
	assertCapabilitiesRegistryConfigured(t, h)

	if cfg.WithMCMS {
		setupMCMSInfrastructure(t, h)
	}

	return h
}

// configureCapabilitiesRegistry configures the capabilities registry in the runtime environment
func configureCapabilitiesRegistry(t *testing.T, h *Harness) {
	t.Helper()

	chainID, err := chain_selectors.GetChainIDFromSelector(h.RegistrySelector)
	require.NoError(t, err)

	registryChainDetails, err := chain_selectors.GetChainDetailsByChainIDAndFamily(chainID, chain_selectors.FamilyEVM)
	require.NoError(t, err)

	donNodes, err := h.Don.AllNodes()
	require.NoError(t, err)

	var nodes []changeset2.CapabilitiesRegistryNodeParams
	for _, n := range donNodes {
		p2pID := n.PeerID.String()
		ocrConfig, ok := n.OCRConfigs[registryChainDetails]
		require.True(t, ok, "node %s does not have OCR config for registry chain %d", n.Name, h.RegistrySelector)

		nodes = append(nodes, changeset2.CapabilitiesRegistryNodeParams{
			NOP:                 "Operator 1",
			P2pID:               p2pID,
			CsaKey:              n.CSA,
			EncryptionPublicKey: n.WorkflowKey,
			Signer:              hex.EncodeToString(ocrConfig.OnchainPublicKey),
			CapabilityIDs: []string{
				TestCapabilityID,
			},
		})
	}

	err = h.Runtime.Exec(
		runtime.ChangesetTask(changeset2.ConfigureCapabilitiesRegistry{},
			changeset2.ConfigureCapabilitiesRegistryInput{
				ChainSelector:               h.RegistrySelector,
				CapabilitiesRegistryAddress: h.RegistryAddress.Hex(),
				Nops: []changeset2.CapabilitiesRegistryNodeOperator{
					{
						Name:  "Operator 1",
						Admin: common.HexToAddress("0x01"),
					},
				},
				Nodes: nodes,
				Capabilities: []changeset2.CapabilitiesRegistryCapability{
					{
						CapabilityID: TestCapabilityID,
						Metadata:     map[string]any{"capabilityType": 2},
					},
				},
				DONs: []changeset2.CapabilitiesRegistryNewDONParams{
					{
						Name:        h.Don.Name(),
						F:           uint8(h.Don.F()), //nolint:gosec // disable G115
						Nodes:       h.Don.GetP2PIDs().Strings(),
						DonFamilies: []string{"test-family"},
						Config:      map[string]any{"defaultConfig": map[string]any{}},
						CapabilityConfigurations: []changeset2.CapabilitiesRegistryCapabilityConfiguration{
							{
								CapabilityID: TestCapabilityID,
							},
						},
						IsPublic:         true,
						AcceptsWorkflows: true,
					},
				},
			},
		),
	)
	require.NoError(t, err)
}

// capabilitiesRegistryClient returns a new capabilities registry client
func capabilitiesRegistryClient(t *testing.T, h *Harness) *capabilities_registry_v2.CapabilitiesRegistry {
	t.Helper()

	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		h.RegistryAddress,
		h.Runtime.Environment().BlockChains.EVMChains()[h.RegistrySelector].Client,
	)
	require.NoError(t, err)
	require.NotNil(t, capReg)

	return capReg
}

// assertCapabilitiesRegistryConfigured asserts that the capabilities registry is configured correctly
func assertCapabilitiesRegistryConfigured(t *testing.T, h *Harness) {
	t.Helper()

	capReg := capabilitiesRegistryClient(t, h)
	expectedP2PIDs := h.Don.GetP2PIDs().Bytes32()

	gotNodes, err := capReg.GetNodesByP2PIds(nil, expectedP2PIDs)
	require.NoError(t, err)
	require.Len(t, gotNodes, h.Don.N())
	for _, n := range gotNodes {
		require.Equal(t, TestCapabilityID, n.CapabilityIds[0])
	}

	gotDON, err := capReg.GetDONByName(nil, h.Don.Name())
	require.NoError(t, err)
	require.ElementsMatch(t, expectedP2PIDs, gotDON.NodeP2PIds)
}

// deployCapabilitiesRegistry deploys the capabilities registry in the runtime environment
func deployCapabilitiesRegistry(t *testing.T, h *Harness) {
	t.Helper()

	err := h.Runtime.Exec(
		runtime.ChangesetTask(changeset2.DeployCapabilitiesRegistry{},
			changeset2.DeployCapabilitiesRegistryInput{
				ChainSelector: h.RegistrySelector,
				Qualifier:     RegistryQualifier,
			},
		),
	)
	require.NoError(t, err)
}

// setupMCMSInfrastructure sets up the MCMS infrastructure in the runtime environment
func setupMCMSInfrastructure(t *testing.T, h *Harness) {
	t.Helper()

	t.Log("Setting up MCMS infrastructure...")

	timelockCfgs := map[uint64]cldfproposalutils.MCMSWithTimelockConfig{
		h.RegistrySelector: cldftesthelpers.SingleGroupTimelockConfig(t),
	}

	err := h.Runtime.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(mcmschangesets.DeployMCMSWithTimelockV2), timelockCfgs),
	)
	require.NoError(t, err, "failed to deploy MCMS infrastructure")

	t.Log("MCMS infrastructure deployed successfully")

	t.Log("Transferring ownership to MCMS...")
	err = h.Runtime.Exec(
		runtime.ChangesetTask(
			cldf.CreateLegacyChangeSet(changeset3.AcceptAllOwnershipsProposal),
			&changeset3.AcceptAllOwnershipRequest{
				ChainSelector: h.RegistrySelector,
				MinDelay:      0,
			},
		),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}),
	)
	require.NoError(t, err, "failed to transfer ownership to MCMS")
	t.Log("Ownership transferred to MCMS successfully")
}

// newViewOnlyNodes creates a view only DON with the given configuration.
func newViewOnlyNodes(t *testing.T, registryChainSel, aptosChainSel uint64, donCfg donConfig) *viewOnlyDon {
	var nodesCfg []envtest.NodeConfig

	for i := 0; i < donCfg.N; i++ {
		labels := map[string]string{
			"don-" + donCfg.Name: donCfg.Name,
			"environment":        EnvironmentName,
			"product":            "cre",
			"type":               "plugin",
			"zone":               Zone,
		}
		if donCfg.Labels != nil {
			maps.Copy(labels, donCfg.Labels)
		}

		nCfg := envtest.NodeConfig{
			ChainSelectors: []uint64{registryChainSel, aptosChainSel, chain_selectors.SOLANA_DEVNET.Selector},
			Name:           fmt.Sprintf("%s-%d", donCfg.Name, i),
			Labels:         labels,
		}
		nodesCfg = append(nodesCfg, nCfg)
	}

	btLabels := map[string]string{
		"don-" + donCfg.Name: donCfg.Name,
		"environment":        EnvironmentName,
		"product":            "cre",
		"type":               "bootstrap",
		"zone":               Zone,
	}
	if donCfg.Labels != nil {
		maps.Copy(btLabels, donCfg.Labels)
	}
	nodesCfg = append(nodesCfg, envtest.NodeConfig{
		ChainSelectors: []uint64{registryChainSel, aptosChainSel, chain_selectors.SOLANA_DEVNET.Selector},
		Name:           donCfg.Name + "-bootstrap",
		Labels:         btLabels,
	})

	n := envtest.NewNodes(t, nodesCfg)
	require.Len(t, n, donCfg.N+1) // +1 for bootstrap

	return newViewOnlyDon(donCfg.Name, n)
}
