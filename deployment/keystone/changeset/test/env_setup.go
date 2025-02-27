package test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	envtest "github.com/smartcontractkit/chainlink/deployment/environment/test"
	changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

type DonConfig struct {
	Name             string // required, must be unique across all dons
	N                int
	F                *int                                          // if nil, defaults to floor(N-1/3) + 1
	CapabilityConfig map[CapabilityNaturalKey]*pb.CapabilityConfig // optional DON specific configuration for the given capability
	Labels           map[string]string                             // optional
	RegistryChainSel uint64                                        // require, must be the same for all dons
	ChainSelectors   []uint64                                      // optional chains

	generatedKeys []importableKeys
}

type BootstrapConfig struct {
	Name   string
	N      int
	Labels map[string]string

	generatedKeys []importableKeys
	bootstrappers []bootstrapperMetadata
}

func (b *BootstrapConfig) Locations() []ocrcommontypes.BootstrapperLocator {
	locations := make([]ocrcommontypes.BootstrapperLocator, len(b.bootstrappers))
	for i, bs := range b.bootstrappers {
		locations[i] = bs.location()
	}
	return locations
}

type bootstrapperMetadata struct {
	port      int
	importKey keystore.ImportableKey
}

func (b bootstrapperMetadata) mustPeerID() p2pkey.PeerID {
	var x p2pkey.EncryptedP2PKeyExport
	err := json.Unmarshal([]byte(b.importKey.JSON), &x)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal bootstrapper key: %v", err))
	}
	return x.PeerID
}

func (b bootstrapperMetadata) location() ocrcommontypes.BootstrapperLocator {
	return ocrcommontypes.BootstrapperLocator{
		PeerID: b.mustPeerID().String(),
		Addrs:  []string{}, // TODO
	}
}

func (b *BootstrapConfig) generateKeys(t *testing.T, ks *keystore.TestKeystore) {
	if b.generatedKeys != nil {
		return
	}
	b.generatedKeys = generateKeys(t, ks, generateKeysCfg{
		N: b.N,
	})
}

type importableKeys struct {
	P2P     keystore.ImportableKey            // required
	EthKeys map[int]keystore.ImportableEthKey // optional
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

func (c *DonConfig) generateKeys(t *testing.T, ks *keystore.TestKeystore) {
	if c.generatedKeys != nil {
		return
	}
	c.generatedKeys = generateKeys(t, ks, generateKeysCfg{
		N:                 c.N,
		EVMChainSelectors: c.ChainSelectors,
	})
}

func (c *DonConfig) GenerateOpts() map[string][]func(c *chainlink.Config, s *chainlink.Secrets) {
	return nil
}

type capabilitiesTOMLConfigurer struct {
	d2dListener string
	don2don     []ocrcommontypes.BootstrapperLocator
	capCfg      deployment.CapabilityRegistryConfig
	wfCfg       *deployment.CapabilityRegistryConfig
}

func (c *capabilitiesTOMLConfigurer) generate() *toml.Capabilities {
	capabilities := &toml.Capabilities{}
	capabilities.Peering.PeerID = nil
	capabilities.Peering.V2.Enabled = ptr(true)
	capabilities.Peering.V2.ListenAddresses = ptr([]string{c.d2dListener})
	capabilities.Peering.V2.DefaultBootstrappers = ptr(c.don2don)
	capabilities.ExternalRegistry.NetworkID = ptr(relay.NetworkEVM)
	capabilities.ExternalRegistry.ChainID = ptr(strconv.FormatUint(uint64(c.capCfg.EVMChainID), 10))
	capabilities.ExternalRegistry.Address = ptr(c.capCfg.Contract.String())
	if c.wfCfg != nil {
		capabilities.WorkflowRegistry.NetworkID = ptr(relay.NetworkEVM)
		capabilities.WorkflowRegistry.ChainID = ptr(strconv.FormatUint(uint64(c.wfCfg.EVMChainID), 10))
		capabilities.WorkflowRegistry.Address = ptr(c.wfCfg.Contract.String())
	}
	// todo gateway
	return capabilities
}

func ptr[T any](v T) *T {
	return &v
}

type testEnvIface interface {
	ContractSets() map[uint64]changeset.ContractSet
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

func (te EnvWrapper) ContractSets() map[uint64]changeset.ContractSet {
	r, err := changeset.GetContractSets(te.Env.Logger, &changeset.GetContractSetsRequest{
		Chains:      te.Env.Chains,
		AddressBook: te.Env.ExistingAddresses,
	})
	require.NoError(te.t, err)
	return r.ContractSets
}

func (te EnvWrapper) CapabilitiesRegistry() *kcr.CapabilitiesRegistry {
	r, err := changeset.GetContractSets(te.Env.Logger, &changeset.GetContractSetsRequest{
		Chains:      te.Env.Chains,
		AddressBook: te.Env.ExistingAddresses,
	})
	require.NoError(te.t, err)
	return r.ContractSets[te.RegistrySelector].CapabilitiesRegistry
}

func (te EnvWrapper) CapabilityInfos() []kcr.CapabilitiesRegistryCapabilityInfo {
	te.t.Helper()
	caps, err := te.CapabilitiesRegistry().GetCapabilities(nil)
	require.NoError(te.t, err)
	return caps
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
		Logger:            logger.Test(t),
		Chains:            chains,
		ExistingAddresses: deployment.NewMemoryAddressBook(),
	}
	env, err := commonchangeset.Apply(t, env, nil,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.DeployCapabilityRegistry),
			registryChainSel,
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.DeployOCR3),
			registryChainSel,
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.DeployForwarder),
			changeset.DeployForwarderRequest{},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(workflowregistry.Deploy),
			registryChainSel,
		),
	)
	require.NoError(t, err)
	require.NotNil(t, env)
	require.Len(t, env.Chains, nChains)
	validateInitialChainState(t, env, registryChainSel)
	return registryChainSel, env
}

func SetupContractTestEnv(t *testing.T, c EnvWrapperConfig) EnvWrapper {
	c.useInMemoryNodes = false
	return setupTestEnv(t, c)
}

func SetupDevTestEnv(t *testing.T, c EnvWrapperConfig) EnvWrapper {
	c.useInMemoryNodes = true
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

	req := &changeset.GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: env.ExistingAddresses,
	}

	contractSetsResp, err := changeset.GetContractSets(lggr, req)
	require.NoError(t, err)
	require.Len(t, contractSetsResp.ContractSets, len(env.Chains))
	// check the registry
	gotRegistry := contractSetsResp.ContractSets[registryChainSel].CapabilitiesRegistry
	require.NotNil(t, gotRegistry)
	// validate the registry
	// check the nodes
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
				deployment.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
				timelockCfgs,
			),
		)
		require.NoError(t, err)
		// extract the MCMS address
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
					deployment.CreateLegacyChangeSet(changeset.AcceptAllOwnershipsProposal),
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
	wfConfig := make([]envtest.NodeConfig, 0, len(c.WFDonConfig.ChainSelectors))
	for i := 0; i < c.WFDonConfig.N; i++ {
		wfConfig = append(wfConfig, envtest.NodeConfig{
			ChainSelectors: []uint64{registryChainSel},
			Name:           fmt.Sprintf("%s-%d", c.WFDonConfig.Name, i),
		})
	}
	wfNodes := envtest.NewNodes(t, wfConfig)
	require.Len(t, wfNodes, c.WFDonConfig.N)

	assetConfig := make([]envtest.NodeConfig, 0, len(c.AssetDonConfig.ChainSelectors))
	for i := 0; i < c.AssetDonConfig.N; i++ {
		assetConfig = append(assetConfig, envtest.NodeConfig{
			ChainSelectors: maps.Keys(chains),
			Name:           fmt.Sprintf("%s-%d", c.AssetDonConfig.Name, i),
		})
	}
	assetNodes := envtest.NewNodes(t, assetConfig)
	require.Len(t, assetNodes, c.AssetDonConfig.N)

	writerConfig := make([]envtest.NodeConfig, 0, len(c.WriterDonConfig.ChainSelectors))
	for i := 0; i < c.WriterDonConfig.N; i++ {
		writerConfig = append(writerConfig, envtest.NodeConfig{
			ChainSelectors: maps.Keys(chains),
			Name:           fmt.Sprintf("%s-%d", c.WriterDonConfig.Name, i),
		})
	}
	writerNodes := envtest.NewNodes(t, writerConfig)
	require.Len(t, writerNodes, c.WriterDonConfig.N)

	dons := newViewOnlyDons()
	dons.Put(newViewOnlyDon(c.WFDonConfig.Name, wfNodes))
	dons.Put(newViewOnlyDon(c.AssetDonConfig.Name, assetNodes))
	dons.Put(newViewOnlyDon(c.WriterDonConfig.Name, writerNodes))

	env := deployment.NewEnvironment(
		"view only nodes",
		logger.Test(t),
		deployment.NewMemoryAddressBook(),
		chains,
		nil,
		dons.NodeList().IDs(),
		envtest.NewJDService(dons.NodeList()),
		func() context.Context { return tests.Context(t) },
		deployment.XXXGenerateTestOCRSecrets(),
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
	wfNodes := memory.NewNodes(t, zapcore.InfoLevel, wfChains, nil, c.WFDonConfig.N, 0, crConfig)
	require.Len(t, wfNodes, c.WFDonConfig.N)

	writerChains := map[uint64]deployment.Chain{}
	maps.Copy(writerChains, chains)
	cwNodes := memory.NewNodes(t, zapcore.InfoLevel, writerChains, nil, c.WriterDonConfig.N, 0, crConfig)
	require.Len(t, cwNodes, c.WriterDonConfig.N)

	assetChains := map[uint64]deployment.Chain{}
	assetChains[registryChainSel] = chains[registryChainSel]
	assetNodes := memory.NewNodes(t, zapcore.InfoLevel, assetChains, nil, c.AssetDonConfig.N, 0, crConfig)
	require.Len(t, assetNodes, c.AssetDonConfig.N)

	dons := newMemoryDons()
	dons.Put(newMemoryDon(c.WFDonConfig.Name, wfNodes))
	dons.Put(newMemoryDon(c.AssetDonConfig.Name, assetNodes))
	dons.Put(newMemoryDon(c.WriterDonConfig.Name, cwNodes))

	env := memory.NewMemoryEnvironmentFromChainsNodes(func() context.Context { return tests.Context(t) }, lggr, chains, nil, dons.AllNodes())
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

type generateKeysCfg struct {
	N                 int      // number of nodes to generate keys for
	EVMChainSelectors []uint64 // only evm supported in the core node secrets today
}

func generateKeys(t *testing.T, ks *keystore.TestKeystore, c generateKeysCfg) []importableKeys {
	keys := make([]importableKeys, c.N)
	for i := 0; i < c.N; i++ {
		keys[i] = importableKeys{
			P2P:     ks.GenerateP2PKey(),
			EthKeys: make(map[int]keystore.ImportableEthKey, len(c.EVMChainSelectors)),
		}
		evmChainIDs := make([]*big.Int, len(c.EVMChainSelectors))
		for j, sel := range c.EVMChainSelectors {
			cid, err := chain_selectors.GetChainIDFromSelector(sel)
			require.NoError(t, err)
			id, ok := big.NewInt(0).SetString(cid, 10)
			require.True(t, ok)
			evmChainIDs[j] = id
		}

		// under the hood, the keystore adds the same key to all the chains, so we only need to add it once
		k, err := ks.Eth().Create(tests.Context(t), evmChainIDs...)
		require.NoError(t, err)
		json, err := ks.Eth().Export(tests.Context(t), k.ID(), "password")
		require.NoError(t, err)
		for j, chainID := range evmChainIDs {
			keys[i].EthKeys[j] = keystore.ImportableEthKey{
				EVMChainID: chainID.Uint64(),
				ImportableKey: keystore.ImportableKey{
					JSON:     string(json),
					Password: "password",
				},
			}
		}
	}
	return keys
}
