package test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	envtest "github.com/smartcontractkit/chainlink/deployment/environment/test"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"
)

var (
	registryQualifier         = "test registry"          // qualifier for the registry chain in the env datastore
	ocr3Qualifier             = "test ocr3"              // qualifier for the ocr3 chain in the env datastore
	forwarderQualifier        = "test forwarder"         // qualifier for the forwarder chain in the env datastore
	workflowRegistryQualifier = "test workflow registry" // qualifier for the workflow registry chain in the env datastore
)

type DonConfig struct {
	Name             string // required, must be unique across all dons
	N                int
	F                *int                                          // if nil, defaults to floor(N-1/3) + 1
	CapabilityConfig map[CapabilityNaturalKey]*pb.CapabilityConfig // optional DON specific configuration for the given capability
	Labels           map[string]string                             // optional
	RegistryChainSel uint64                                        // require, must be the same for all dons
	ChainSelectors   []uint64                                      // optional chains
}

type CapabilityNaturalKey struct {
	LabelledName string
	Version      string
}

func (c DonConfig) Validate() error {
	if c.N < 4 {
		return errors.New("N must be at least 4")
	}
	return nil
}

type testEnvIface interface {
	CapabilitiesRegistry() *kcr.CapabilitiesRegistry
	CapabilityInfos() []kcr.CapabilitiesRegistryCapabilityInfo
	Nops() []kcr.CapabilitiesRegistryNodeOperatorAdded

	GetP2PIDs(donName string) P2PIDs
}

// TODO: separate the config into different types; wf should expand to types of ocr keybundles; writer to target chains; ...
type WFDonConfig = DonConfig
type AssetDonConfig = DonConfig
type WriterDonConfig = DonConfig

type EnvWrapperConfig struct {
	WFDonConfig
	AssetDonConfig
	WriterDonConfig
	NumChains int

	UseMCMS bool
	// if true, use in-memory nodes for testing
	// if false, view only nodes will be used
	useInMemoryNodes bool
}

func (c EnvWrapperConfig) Validate() error {
	if err := c.WFDonConfig.Validate(); err != nil {
		return err
	}

	if err := c.AssetDonConfig.Validate(); err != nil {
		return err
	}

	if err := c.WriterDonConfig.Validate(); err != nil {
		return err
	}

	if c.NumChains < 1 {
		return errors.New("NumChains must be at least 1")
	}
	return nil
}

var _ testEnvIface = (*EnvWrapper)(nil)

type EnvWrapper struct {
	t                *testing.T
	Env              deployment.Environment
	RegistrySelector uint64

	dons testDons
}

func (te EnvWrapper) CapabilitiesRegistry() *kcr.CapabilitiesRegistry {
	return te.OwnedCapabilityRegistry().Contract
}

func (te EnvWrapper) CapabilityInfos() []kcr.CapabilitiesRegistryCapabilityInfo {
	te.t.Helper()
	caps, err := te.CapabilitiesRegistry().GetCapabilities(nil)
	require.NoError(te.t, err)
	return caps
}

func (te EnvWrapper) OwnedCapabilityRegistry() *changeset.OwnedContract[*kcr.CapabilitiesRegistry] {
	return loadOneContract[*kcr.CapabilitiesRegistry](te.t, te.Env, te.Env.Chains[te.RegistrySelector], registryQualifier)
}

func loadOneContract[T changeset.Ownable](t *testing.T, env deployment.Environment, chain deployment.Chain, qualifier string) *changeset.OwnedContract[T] {
	t.Helper()
	addrs := env.DataStore.Addresses().Filter(datastore.AddressRefByQualifier(qualifier))
	require.Len(t, addrs, 1)
	c, err := changeset.GetOwnedContractV2[T](env.DataStore.Addresses(), chain, addrs[0].Address)
	require.NoError(t, err)
	require.NotNil(t, c)
	return c
}

func (te EnvWrapper) CapabilityRegistryAddressRef() datastore.AddressRefKey {
	addrs := te.Env.DataStore.Addresses().Filter(datastore.AddressRefByQualifier(registryQualifier))
	require.Len(te.t, addrs, 1)
	return addrs[0].Key()
}

func (te EnvWrapper) ForwarderAddressRefs() []datastore.AddressRefKey {
	addrs := te.Env.DataStore.Addresses().Filter(datastore.AddressRefByQualifier(forwarderQualifier))
	require.NotEmpty(te.t, addrs)
	out := make([]datastore.AddressRefKey, len(addrs))
	for i, addr := range addrs {
		out[i] = addr.Key()
	}
	return out
}

func (te EnvWrapper) OwnedForwarders() map[uint64][]*changeset.OwnedContract[*forwarder.KeystoneForwarder] { // chain selector -> forwarders
	addrs := te.Env.DataStore.Addresses().Filter(datastore.AddressRefByQualifier(forwarderQualifier))
	require.NotEmpty(te.t, addrs)
	out := make(map[uint64][]*changeset.OwnedContract[*forwarder.KeystoneForwarder])
	for _, addr := range addrs {
		c, err := changeset.GetOwnedContractV2[*forwarder.KeystoneForwarder](te.Env.DataStore.Addresses(), te.Env.Chains[addr.ChainSelector], addr.Address)
		require.NoError(te.t, err)
		require.NotNil(te.t, c)
		out[addr.ChainSelector] = append(out[addr.ChainSelector], c)
	}
	return out
}

func (te EnvWrapper) Nops() []kcr.CapabilitiesRegistryNodeOperatorAdded {
	te.t.Helper()
	nops, err := te.CapabilitiesRegistry().GetNodeOperators(nil)
	require.NoError(te.t, err)
	out := make([]kcr.CapabilitiesRegistryNodeOperatorAdded, len(nops))
	id := uint32(0)
	for i, n := range nops {
		out[i] = kcr.CapabilitiesRegistryNodeOperatorAdded{
			NodeOperatorId: id + 1, // 1-indexed
			Admin:          n.Admin,
			Name:           n.Name,
		}
	}
	return out
}

func (te EnvWrapper) GetP2PIDs(donName string) P2PIDs {
	return te.dons.Get(donName).GetP2PIDs()
}

func initEnv(t *testing.T, nChains int) (registryChainSel uint64, env deployment.Environment) {
	chains, _ := memory.NewMemoryChains(t, nChains, 1)
	registryChainSel = registryChain(t, chains)
	// note that all the nodes require TOML configuration of the cap registry address
	// and writers need forwarder address as TOML config
	// we choose to use changesets to deploy the initial contracts because that's how it's done in the real world
	// this requires a initial environment to house the address book
	env = deployment.Environment{
		GetContext:        t.Context,
		Logger:            logger.Test(t),
		Chains:            chains,
		ExistingAddresses: deployment.NewMemoryAddressBook(),
		DataStore:         datastore.NewMemoryDataStore[datastore.DefaultMetadata, datastore.DefaultMetadata]().Seal(),
	}

	forwarderChangesets := make([]commonchangeset.ConfiguredChangeSet, nChains)
	i := 0
	for _, c := range chains {
		forwarderChangesets[i] = commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.DeployForwarderV2),
			&changeset.DeployRequestV2{
				ChainSel:  c.Selector,
				Qualifier: forwarderQualifier,
			},
		)
		i++
	}

	changes := []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.DeployCapabilityRegistryV2),
			&changeset.DeployRequestV2{
				ChainSel:  registryChainSel,
				Qualifier: registryQualifier,
				Labels:    nil,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.DeployOCR3V2),
			&changeset.DeployRequestV2{
				ChainSel:  registryChainSel,
				Qualifier: ocr3Qualifier,
				Labels:    nil,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(workflowregistry.DeployV2),
			&changeset.DeployRequestV2{
				ChainSel:  registryChainSel,
				Qualifier: workflowRegistryQualifier,
				Labels:    nil,
			},
		),
	}
	changes = append(changes, forwarderChangesets...)
	env, err := commonchangeset.ApplyChangesets(t, env, nil,
		changes,
	)
	require.NoError(t, err)
	require.NotNil(t, env)
	require.Len(t, env.Chains, nChains)
	validateInitialChainState(t, env, registryChainSel)
	return registryChainSel, env
}

// SetupContractTestEnv sets up a keystone test environment for contract testing with the given configuration
// The resulting environment will have the following:
//
// - all the initial contracts deployed (capability registry, ocr3, forwarder, workflow registry) on the registry chain
//
// - the forwarder deployed on all chains
//
// - the capability registry configured with the initial dons (WFDon => ocr capability, AssetDon => stream trigger capability, WriterDon => writer capability for all the chains)
//
// - a view-only Offchain client that supports all the read api operations of the Offchain client
func SetupContractTestEnv(t *testing.T, c EnvWrapperConfig) EnvWrapper {
	c.useInMemoryNodes = false
	return setupTestEnv(t, c)
}

// SetupContractTestEnv sets up a keystone test environment with the given configuration
// TODO: make more configurable; eg many tests don't need all the nodes (like when testing a registry change)
func setupTestEnv(t *testing.T, c EnvWrapperConfig) EnvWrapper {
	require.NoError(t, c.Validate())
	lggr := logger.Test(t)

	registryChainSel, envWithContracts := initEnv(t, c.NumChains)
	lggr.Debug("done init env")
	var (
		dons testDons
		env  deployment.Environment
	)
	if c.useInMemoryNodes {
		dons, env = setupMemoryNodeTest(t, registryChainSel, envWithContracts.Chains, c)
	} else {
		dons, env = setupViewOnlyNodeTest(t, registryChainSel, envWithContracts.Chains, c)
	}
	err := env.ExistingAddresses.Merge(envWithContracts.ExistingAddresses)
	require.NoError(t, err)
	env.DataStore = envWithContracts.DataStore

	ocr3CapCfg := GetDefaultCapConfig(t, internal.OCR3Cap)
	writerChainCapCfg := GetDefaultCapConfig(t, internal.WriteChainCap)
	streamTriggerChainCapCfg := GetDefaultCapConfig(t, internal.StreamTriggerCap)

	// TODO: partition nodes into multiple nops

	wfDonCapabilities := internal.DonCapabilities{
		Name: c.WFDonConfig.Name,
		Nops: []internal.NOP{
			{
				Name:  "nop 1",
				Nodes: dons.Get(c.WFDonConfig.Name).GetP2PIDs().Strings(),
			},
		},
		Capabilities: []internal.DONCapabilityWithConfig{
			{Capability: internal.OCR3Cap, Config: ocr3CapCfg},
		},
	}
	cwDonCapabilities := internal.DonCapabilities{
		Name: c.WriterDonConfig.Name,
		Nops: []internal.NOP{
			{
				Name:  "nop 2",
				Nodes: dons.Get(c.WriterDonConfig.Name).GetP2PIDs().Strings(),
			},
		},
		Capabilities: []internal.DONCapabilityWithConfig{
			{Capability: internal.WriteChainCap, Config: writerChainCapCfg},
		},
	}
	assetDonCapabilities := internal.DonCapabilities{
		Name: c.AssetDonConfig.Name,
		Nops: []internal.NOP{
			{
				Name:  "nop 3",
				Nodes: dons.Get(c.AssetDonConfig.Name).GetP2PIDs().Strings(),
			},
		},
		Capabilities: []internal.DONCapabilityWithConfig{
			{Capability: internal.StreamTriggerCap, Config: streamTriggerChainCapCfg},
		},
	}

	var ocr3Config = internal.OracleConfig{
		MaxFaultyOracles:     dons.Get(c.WFDonConfig.Name).F(),
		TransmissionSchedule: []int{dons.Get(c.WFDonConfig.Name).N()},
	}
	var allDons = []internal.DonCapabilities{wfDonCapabilities, cwDonCapabilities, assetDonCapabilities}

	csOut, err := changeset.ConfigureInitialContractsChangeset(env, changeset.InitialContractsCfg{
		RegistryChainSel: registryChainSel,
		Dons:             allDons,
		OCR3Config:       &ocr3Config,
	})
	require.NoError(t, err)
	require.Nil(t, csOut.AddressBook, "no new addresses should be created in configure initial contracts")

	// check the registry
	gotOwnedRegistry := loadOneContract[*kcr.CapabilitiesRegistry](t, env, env.Chains[registryChainSel], registryQualifier)
	require.NotNil(t, gotOwnedRegistry)
	// validate the registry
	// check the nodes
	gotRegistry := gotOwnedRegistry.Contract
	gotNodes, err := gotRegistry.GetNodes(nil)
	require.NoError(t, err)
	require.Len(t, gotNodes, len(dons.P2PIDs()))
	validateNodes(t, gotRegistry, dons.Get(c.WFDonConfig.Name), expectedHashedCapabilities(t, gotRegistry, wfDonCapabilities))
	validateNodes(t, gotRegistry, dons.Get(c.WriterDonConfig.Name), expectedHashedCapabilities(t, gotRegistry, cwDonCapabilities))
	validateNodes(t, gotRegistry, dons.Get(c.AssetDonConfig.Name), expectedHashedCapabilities(t, gotRegistry, assetDonCapabilities))

	// check the dons
	validateDon(t, gotRegistry, dons.Get(c.WFDonConfig.Name), wfDonCapabilities)
	validateDon(t, gotRegistry, dons.Get(c.WriterDonConfig.Name), cwDonCapabilities)
	validateDon(t, gotRegistry, dons.Get(c.AssetDonConfig.Name), assetDonCapabilities)

	if c.UseMCMS {
		// deploy, configure and xfer ownership of MCMS on all chains
		timelockCfgs := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
		for sel := range env.Chains {
			t.Logf("Enabling MCMS on chain %d", sel)
			timelockCfgs[sel] = proposalutils.SingleGroupTimelockConfigV2(t)
		}
		env, err = commonchangeset.Apply(t, env, nil,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
				timelockCfgs,
			),
		)
		require.NoError(t, err)
		// extract the MCMS address using `GetContractSets` instead of `GetContractSetsV2` because the latter
		// expects contracts to already be owned by MCMS
		r, err := changeset.GetContractSets(lggr, &changeset.GetContractSetsRequest{
			Chains:      env.Chains,
			AddressBook: env.ExistingAddresses,
		})
		require.NoError(t, err)
		for sel := range env.Chains {
			mcms := r.ContractSets[sel].MCMSWithTimelockState
			require.NotNil(t, mcms, "MCMS not found on chain %d", sel)
			require.NoError(t, mcms.Validate())

			// transfer ownership of all contracts to the MCMS
			env, err = commonchangeset.Apply(t, env,
				map[uint64]*proposalutils.TimelockExecutionContracts{
					sel: {Timelock: mcms.Timelock, CallProxy: mcms.CallProxy},
				},
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(changeset.AcceptAllOwnershipsProposal),
					&changeset.AcceptAllOwnershipRequest{
						ChainSelector: sel,
						MinDelay:      0,
					},
				),
			)
			require.NoError(t, err)
		}
	}
	return EnvWrapper{
		t:                t,
		Env:              env,
		RegistrySelector: registryChainSel,
		dons:             dons,
	}
}

func setupViewOnlyNodeTest(t *testing.T, registryChainSel uint64, chains map[uint64]deployment.Chain, c EnvWrapperConfig) (testDons, deployment.Environment) {
	// now that we have the initial contracts deployed, we can configure the nodes with the addresses
	dons := newViewOnlyDons()
	for _, donCfg := range []DonConfig{c.WFDonConfig, c.AssetDonConfig, c.WriterDonConfig} {
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
			ncfg = append(ncfg, envtest.NodeConfig{
				ChainSelectors: []uint64{registryChainSel},
				Name:           fmt.Sprintf("%s-%d", donCfg.Name, i),
				Labels:         labels,
			})
		}
		n := envtest.NewNodes(t, ncfg)
		require.Len(t, n, donCfg.N)
		dons.Put(newViewOnlyDon(donCfg.Name, n))
	}

	env := deployment.NewEnvironment(
		"view only nodes",
		logger.Test(t),
		deployment.NewMemoryAddressBook(),
		datastore.NewMemoryDataStore[
			datastore.DefaultMetadata,
			datastore.DefaultMetadata,
		]().Seal(),
		chains,
		nil,
		nil,
		dons.NodeList().IDs(),
		envtest.NewJDService(dons.NodeList()),
		t.Context,
		cldf.XXXGenerateTestOCRSecrets(),
	)

	return dons, *env
}

func setupMemoryNodeTest(t *testing.T, registryChainSel uint64, chains map[uint64]deployment.Chain, c EnvWrapperConfig) (testDons, deployment.Environment) {
	// now that we have the initial contracts deployed, we can configure the nodes with the addresses
	// TODO: configure the nodes with the correct override functions
	lggr := logger.Test(t)
	crConfig := deployment.CapabilityRegistryConfig{
		EVMChainID: registryChainSel,
		Contract:   [20]byte{},
	}

	wfChains := map[uint64]deployment.Chain{}
	wfChains[registryChainSel] = chains[registryChainSel]
	wfNodes := memory.NewNodes(t, zapcore.InfoLevel, wfChains, nil, nil, c.WFDonConfig.N, 0, crConfig, nil)
	require.Len(t, wfNodes, c.WFDonConfig.N)

	writerChains := map[uint64]deployment.Chain{}
	maps.Copy(writerChains, chains)
	cwNodes := memory.NewNodes(t, zapcore.InfoLevel, writerChains, nil, nil, c.WriterDonConfig.N, 0, crConfig, nil)
	require.Len(t, cwNodes, c.WriterDonConfig.N)

	assetChains := map[uint64]deployment.Chain{}
	assetChains[registryChainSel] = chains[registryChainSel]
	assetNodes := memory.NewNodes(t, zapcore.InfoLevel, assetChains, nil, nil, c.AssetDonConfig.N, 0, crConfig, nil)
	require.Len(t, assetNodes, c.AssetDonConfig.N)

	dons := newMemoryDons()
	dons.Put(newMemoryDon(c.WFDonConfig.Name, wfNodes))
	dons.Put(newMemoryDon(c.AssetDonConfig.Name, assetNodes))
	dons.Put(newMemoryDon(c.WriterDonConfig.Name, cwNodes))

	env := memory.NewMemoryEnvironmentFromChainsNodes(t.Context, lggr, chains, nil, nil, dons.AllNodes())
	return dons, env
}

func registryChain(t *testing.T, chains map[uint64]deployment.Chain) uint64 {
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
func validateInitialChainState(t *testing.T, env deployment.Environment, registryChainSel uint64) {
	ad := env.ExistingAddresses
	// all contracts on registry chain
	registryChainAddrs, err := ad.AddressesForChain(registryChainSel)
	require.NoError(t, err)
	require.Len(t, registryChainAddrs, 4) // registry, ocr3, forwarder, workflowRegistry
	// only forwarder on non-home chain
	for sel := range env.Chains {
		chainAddrs, err := ad.AddressesForChain(sel)
		require.NoError(t, err)
		if sel != registryChainSel {
			require.Len(t, chainAddrs, 1)
		} else {
			require.Len(t, chainAddrs, 4)
		}
		containsForwarder := false
		for _, tv := range chainAddrs {
			if tv.Type == internal.KeystoneForwarder {
				containsForwarder = true
				break
			}
		}
		require.True(t, containsForwarder, "no forwarder found in %v on chain %d for target don", chainAddrs, sel)
	}
}

// validateNodes checks that the nodes exist and have the expected capabilities
func validateNodes(t *testing.T, gotRegistry *kcr.CapabilitiesRegistry, nodes testDon, expectedHashedCaps [][32]byte) {
	gotNodes, err := gotRegistry.GetNodesByP2PIds(nil, p2p32Bytes(t, nodes.GetP2PIDs()))
	require.NoError(t, err)
	require.Len(t, gotNodes, nodes.N())
	for _, n := range gotNodes {
		require.Equal(t, expectedHashedCaps, n.HashedCapabilityIds)
	}
}

// validateDon checks that the don exists and has the expected capabilities
func validateDon(t *testing.T, gotRegistry *kcr.CapabilitiesRegistry, nodes testDon, don internal.DonCapabilities) {
	gotDons, err := gotRegistry.GetDONs(nil)
	require.NoError(t, err)
	wantP2PID := sortedHash(p2p32Bytes(t, nodes.GetP2PIDs()))
	found := false
	for _, have := range gotDons {
		gotP2PID := sortedHash(have.NodeP2PIds)
		if gotP2PID == wantP2PID {
			found = true
			gotCapIDs := capIDs(t, have.CapabilityConfigurations)
			require.Equal(t, expectedHashedCapabilities(t, gotRegistry, don), gotCapIDs)
			break
		}
	}
	require.True(t, found, "don not found in registry")
}
