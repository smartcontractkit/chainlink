package test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	envtest "github.com/smartcontractkit/chainlink/deployment/environment/test"
	kschangeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
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

type P2PIDs []p2pkey.PeerID

func (ps P2PIDs) Strings() []string {
	out := make([]string, len(ps))
	for i, p := range ps {
		out[i] = p.String()
	}
	return out
}

func (ps P2PIDs) Bytes32() [][32]byte {
	out := make([][32]byte, len(ps))
	for i, p := range ps {
		out[i] = p
	}
	return out
}

func (ps P2PIDs) Unique() P2PIDs {
	dedup := make(map[p2pkey.PeerID]struct{})
	var out []p2pkey.PeerID
	for _, p := range ps {
		if _, exists := dedup[p]; !exists {
			out = append(out, p)

			dedup[p] = struct{}{}
		}
	}
	return out
}

type TestEnvI interface {
	ContractSets() map[uint64]internal.ContractSet
	CapabilitiesRegistry() *kcr.CapabilitiesRegistry
	CapabilityInfos() []kcr.CapabilitiesRegistryCapabilityInfo
	Nops() []kcr.CapabilitiesRegistryNodeOperatorAdded

	GetP2PIDs(donName string) P2PIDs
}

type TestDon interface {
	GetP2PIDs() P2PIDs
	N() int
	F() int
	Name() string
}

var _ TestDon = (*memoryDon)(nil)

type memoryDon struct {
	name string
	m    map[string]memory.Node
}

func newMemoryDon(name string, m map[string]memory.Node) *memoryDon {
	return &memoryDon{name: name, m: m}
}

func (d *memoryDon) GetP2PIDs() P2PIDs {
	var out []p2pkey.PeerID
	for _, n := range d.m {
		out = append(out, n.Keys.PeerID)
	}
	return out
}

func (d *memoryDon) N() int {
	return len(d.m)
}

func (d *memoryDon) F() int {
	return (d.N() - 1) / 3
}

func (d *memoryDon) Name() string {
	return d.name
}

type viewOnlyDon struct {
	name string
	m    map[string]*deployment.Node
}

func newViewOnlyDon(name string, nodes []*deployment.Node) *viewOnlyDon {
	m := make(map[string]*deployment.Node)
	for _, n := range nodes {
		m[n.PeerID.String()] = n
	}
	return &viewOnlyDon{name: name, m: m}
}

func (d *viewOnlyDon) GetP2PIDs() P2PIDs {
	var out []p2pkey.PeerID
	for _, n := range d.m {
		out = append(out, n.PeerID)
	}
	return out
}

func (d *viewOnlyDon) N() int {
	return len(d.m)
}

func (d *viewOnlyDon) F() int {
	return (d.N() - 1) / 3
}

func (d *viewOnlyDon) Name() string {
	return d.name
}

type TestDons interface {
	Get(name string) TestDon
	Put(d TestDon)
	List() []TestDon
	// Unique list of p2pIDs across all dons
	P2PIDs() P2PIDs
}

var _ TestDons = (*memoryDons)(nil)

type memoryDons struct {
	dons map[string]*memoryDon
}

func newMemoryDons() *memoryDons {
	return &memoryDons{dons: make(map[string]*memoryDon)}
}

func (d *memoryDons) Get(name string) TestDon {
	x := d.dons[name]
	return x
}

func (d *memoryDons) Put(d2 TestDon) {
	d.dons[d2.Name()] = d2.(*memoryDon)
}

func (d *memoryDons) List() []TestDon {
	out := make([]TestDon, 0, len(d.dons))
	donNames := make([]string, 0, len(d.dons))
	for k := range d.dons {
		donNames = append(donNames, k)
	}
	sort.Strings(donNames)
	for _, name := range donNames {
		out = append(out, d.dons[name])
	}
	return out
}

func (d *memoryDons) P2PIDs() P2PIDs {
	var out P2PIDs
	for _, d := range d.dons {
		out = append(out, d.GetP2PIDs()...)
	}
	return out.Unique()
}

func (d *memoryDons) AllNodes() map[string]memory.Node {
	out := make(map[string]memory.Node)
	for _, d := range d.dons {
		for k, v := range d.m {
			out[k] = v
		}
	}
	return out
}

type viewOnlyDons struct {
	dons map[string]*viewOnlyDon
}

func newViewOnlyDons() *viewOnlyDons {
	return &viewOnlyDons{dons: make(map[string]*viewOnlyDon)}
}

func (d *viewOnlyDons) Get(name string) TestDon {
	x := d.dons[name]
	return x
}

func (d *viewOnlyDons) Put(d2 TestDon) {
	d.dons[d2.Name()] = d2.(*viewOnlyDon)
}

func (d *viewOnlyDons) List() []TestDon {
	out := make([]TestDon, 0, len(d.dons))
	donNames := make([]string, 0, len(d.dons))
	for k := range d.dons {
		donNames = append(donNames, k)
	}
	sort.Strings(donNames)
	for _, name := range donNames {
		out = append(out, d.dons[name])
	}
	return out
}

func (d *viewOnlyDons) P2PIDs() P2PIDs {
	var out P2PIDs
	for _, d := range d.dons {
		out = append(out, d.GetP2PIDs()...)
	}
	return out.Unique()
}

func (d *viewOnlyDons) AllNodes() map[string]*deployment.Node {
	out := make(map[string]*deployment.Node)
	for _, d := range d.dons {
		for k, v := range d.m {
			out[k] = v
		}
	}
	return out
}

func (d *viewOnlyDons) NodeList() deployment.Nodes {
	tmp := d.AllNodes()
	nodes := make([]deployment.Node, 0, len(tmp))
	for _, v := range tmp {
		nodes = append(nodes, *v)
	}
	return nodes
}

// TODO: separate the config into different types; wf should expand to types of ocr keybundles; writer to target chains; ...
type WFDonConfig = DonConfig
type AssetDonConfig = DonConfig
type WriterDonConfig = DonConfig

type TestConfig struct {
	WFDonConfig
	AssetDonConfig
	WriterDonConfig
	NumChains int

	contractOnly bool
	UseMCMS      bool
}

func (c TestConfig) Validate() error {
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

var _ TestEnvI = (*TestEnv)(nil)

type TestEnv struct {
	t                *testing.T
	Env              deployment.Environment
	RegistrySelector uint64

	dons TestDons
}

func (te TestEnv) ContractSets() map[uint64]internal.ContractSet {
	r, err := internal.GetContractSets(te.Env.Logger, &internal.GetContractSetsRequest{
		Chains:      te.Env.Chains,
		AddressBook: te.Env.ExistingAddresses,
	})
	require.NoError(te.t, err)
	return r.ContractSets
}

func (te TestEnv) CapabilitiesRegistry() *kcr.CapabilitiesRegistry {
	r, err := internal.GetContractSets(te.Env.Logger, &internal.GetContractSetsRequest{
		Chains:      te.Env.Chains,
		AddressBook: te.Env.ExistingAddresses,
	})
	require.NoError(te.t, err)
	return r.ContractSets[te.RegistrySelector].CapabilitiesRegistry
}

func (te TestEnv) CapabilityInfos() []kcr.CapabilitiesRegistryCapabilityInfo {
	te.t.Helper()
	caps, err := te.CapabilitiesRegistry().GetCapabilities(nil)
	require.NoError(te.t, err)
	return caps
}

func (te TestEnv) Nops() []kcr.CapabilitiesRegistryNodeOperatorAdded {
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

func (te TestEnv) GetP2PIDs(donName string) P2PIDs {
	return te.dons.Get(donName).GetP2PIDs()
}

func memoryNodesP2pIDs(t *testing.T, m map[string]memory.Node) []p2pkey.PeerID {
	var out []p2pkey.PeerID
	for _, n := range m {
		out = append(out, n.Keys.PeerID)
	}
	return out
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
	var err error
	env, err = commonchangeset.ApplyChangesets(t, env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(kschangeset.DeployCapabilityRegistry),
			Config:    registryChainSel,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(kschangeset.DeployOCR3),
			Config:    registryChainSel,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(kschangeset.DeployForwarder),
			Config:    kschangeset.DeployForwarderRequest{},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(workflowregistry.Deploy),
			Config:    registryChainSel,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, env)
	require.Len(t, env.Chains, nChains)
	validateInitialChainState(t, env, registryChainSel)
	return registryChainSel, env
}

func SetupContractOnlyTestEnv(t *testing.T, c TestConfig) TestEnv {
	c.contractOnly = true
	return SetupTestEnv(t, c)
}

// SetupTestEnv sets up a keystone test environment with the given configuration
// TODO: make more configurable; eg many tests don't need all the nodes (like when testing a registry change)
func SetupTestEnv(t *testing.T, c TestConfig) TestEnv {
	require.NoError(t, c.Validate())
	lggr := logger.Test(t)

	registryChainSel, envWithContracts := initEnv(t, c.NumChains)
	lggr.Debug("done init env")
	var (
		dons TestDons
		env  deployment.Environment
	)
	if c.contractOnly {
		dons, env = setupViewOnlyNodeTest(t, registryChainSel, envWithContracts.Chains, c)
	} else {
		dons, env = setupMemoryNodeTest(t, registryChainSel, envWithContracts.Chains, c)
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

	csOut, err := kschangeset.ConfigureInitialContractsChangeset(env, kschangeset.InitialContractsCfg{
		RegistryChainSel: registryChainSel,
		Dons:             allDons,
		OCR3Config:       &ocr3Config,
	})
	require.NoError(t, err)
	require.Nil(t, csOut.AddressBook, "no new addresses should be created in configure initial contracts")

	req := &internal.GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: env.ExistingAddresses,
	}

	contractSetsResp, err := internal.GetContractSets(lggr, req)
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
		timelockCfgs := make(map[uint64]commontypes.MCMSWithTimelockConfig)
		for sel := range env.Chains {
			t.Logf("Enabling MCMS on chain %d", sel)
			timelockCfgs[sel] = proposalutils.SingleGroupTimelockConfig(t)
		}
		env, err = commonchangeset.ApplyChangesets(t, env, nil, []commonchangeset.ChangesetApplication{
			{
				Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
				Config:    timelockCfgs,
			},
		})
		require.NoError(t, err)
		// extract the MCMS address
		r, err := internal.GetContractSets(lggr, &internal.GetContractSetsRequest{
			Chains:      env.Chains,
			AddressBook: env.ExistingAddresses,
		})
		require.NoError(t, err)
		for sel := range env.Chains {
			mcms := r.ContractSets[sel].MCMSWithTimelockState
			require.NotNil(t, mcms, "MCMS not found on chain %d", sel)
			require.NoError(t, mcms.Validate())

			// transfer ownership of all contracts to the MCMS
			env, err = commonchangeset.ApplyChangesets(t, env, map[uint64]*proposalutils.TimelockExecutionContracts{sel: {Timelock: mcms.Timelock, CallProxy: mcms.CallProxy}}, []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(kschangeset.AcceptAllOwnershipsProposal),
					Config: &kschangeset.AcceptAllOwnershipRequest{
						ChainSelector: sel,
						MinDelay:      0,
					},
				},
			})
			require.NoError(t, err)
		}
	}
	return TestEnv{
		t:                t,
		Env:              env,
		RegistrySelector: registryChainSel,
		dons:             dons,
	}
}

func setupViewOnlyNodeTest(t *testing.T, registryChainSel uint64, chains map[uint64]deployment.Chain, c TestConfig) (TestDons, deployment.Environment) {
	// now that we have the initial contracts deployed, we can configure the nodes with the addresses
	wfConfig := make([]envtest.TestNodeConfig, 0, len(c.WFDonConfig.ChainSelectors))
	for i := 0; i < c.WFDonConfig.N; i++ {
		wfConfig = append(wfConfig, envtest.TestNodeConfig{
			ChainSelectors: []uint64{registryChainSel},
			Name:           fmt.Sprintf("%s-%d", c.WFDonConfig.Name, i),
		})
	}
	wfNodes := envtest.NewNodes(t, wfConfig)
	require.Len(t, wfNodes, c.WFDonConfig.N)

	assetConfig := make([]envtest.TestNodeConfig, 0, len(c.AssetDonConfig.ChainSelectors))
	for i := 0; i < c.AssetDonConfig.N; i++ {
		assetConfig = append(assetConfig, envtest.TestNodeConfig{
			ChainSelectors: maps.Keys(chains),
			Name:           fmt.Sprintf("%s-%d", c.AssetDonConfig.Name, i),
		})
	}
	assetNodes := envtest.NewNodes(t, assetConfig)
	require.Len(t, assetNodes, c.AssetDonConfig.N)

	writerConfig := make([]envtest.TestNodeConfig, 0, len(c.WriterDonConfig.ChainSelectors))
	for i := 0; i < c.WriterDonConfig.N; i++ {
		writerConfig = append(writerConfig, envtest.TestNodeConfig{
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

func setupMemoryNodeTest(t *testing.T, registryChainSel uint64, chains map[uint64]deployment.Chain, c TestConfig) (TestDons, deployment.Environment) {
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
func validateNodes(t *testing.T, gotRegistry *kcr.CapabilitiesRegistry, nodes TestDon, expectedHashedCaps [][32]byte) {
	gotNodes, err := gotRegistry.GetNodesByP2PIds(nil, p2p32Bytes(t, nodes.GetP2PIDs()))
	require.NoError(t, err)
	require.Len(t, gotNodes, nodes.N())
	for _, n := range gotNodes {
		require.Equal(t, expectedHashedCaps, n.HashedCapabilityIds)
	}
}

// validateDon checks that the don exists and has the expected capabilities
func validateDon(t *testing.T, gotRegistry *kcr.CapabilitiesRegistry, nodes TestDon, don internal.DonCapabilities) {
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

func capIDs(t *testing.T, cfgs []kcr.CapabilitiesRegistryCapabilityConfiguration) [][32]byte {
	var out [][32]byte
	for _, cfg := range cfgs {
		out = append(out, cfg.CapabilityId)
	}
	return out
}

func expectedHashedCapabilities(t *testing.T, registry *kcr.CapabilitiesRegistry, don internal.DonCapabilities) [][32]byte {
	out := make([][32]byte, len(don.Capabilities))
	var err error
	for i, capWithCfg := range don.Capabilities {
		out[i], err = registry.GetHashedCapabilityId(nil, capWithCfg.Capability.LabelledName, capWithCfg.Capability.Version)
		require.NoError(t, err)
	}
	return out
}

func sortedHash(p2pids [][32]byte) string {
	sha256Hash := sha256.New()
	sort.Slice(p2pids, func(i, j int) bool {
		return bytes.Compare(p2pids[i][:], p2pids[j][:]) < 0
	})
	for _, id := range p2pids {
		sha256Hash.Write(id[:])
	}
	return hex.EncodeToString(sha256Hash.Sum(nil))
}

func p2p32Bytes(t *testing.T, p2pIDs []p2pkey.PeerID) [][32]byte {
	bs := make([][32]byte, len(p2pIDs))
	for i, p := range p2pIDs {
		bs[i] = p
	}
	return bs
}

func p2pStrings(t *testing.T, p2pIDs []p2pkey.PeerID) []string {
	bs := make([]string, len(p2pIDs))
	for i, p := range p2pIDs {
		bs[i] = p.String()
	}
	return bs
}
