package test

import (
	"context"
	"fmt"
	"maps"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type Don struct {
	Name              string
	P2PIDs            []p2pkey.PeerID
	CapabilityConfigs []internal.CapabilityConfig
}

type SetupTestRegistryRequest struct {
	// P2pToCapabilities maps a node's p2pID to the capabilities it has
	P2pToCapabilities map[p2pkey.PeerID][]capabilities_registry.CapabilitiesRegistryCapability

	// NopToNodes maps a node operator to the nodes they operate
	NopToNodes map[capabilities_registry.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc

	// Dons groups the p2pIDs of the nodes that comprise it and the capabilities they have
	Dons []Don
	// TODO maybe add support for MCMS at this level
}

type SetupTestRegistryResponse struct {
	CapabilitiesRegistry *capabilities_registry.CapabilitiesRegistry
	Chain                deployment.Chain
	RegistrySelector     uint64
	CapabilityCache      *CapabilityCache
	Nops                 []*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded
}

// SetupTestRegistry deploys a capabilities registry to the given chain
// and adds the given capabilities and node operators
// It can be used in tests that mutate the registry without any other setup such as actual nodes, dons, jobs, etc.
func SetupTestRegistry(t *testing.T, lggr logger.Logger, req *SetupTestRegistryRequest) *SetupTestRegistryResponse {
	chain := testChain(t)

	// deploy the registry
	registry := deployCapReg(t, chain)

	// convert req to nodeoperators
	nops := ToNodeOps(t, req.NopToNodes)
	addNopsResp := addNops(t, lggr, chain, registry, nops)
	require.Len(t, addNopsResp.Nops, len(nops))

	// add capabilities to registry
	registeredCapabilities, capCache := MustAddCapabilities(t, lggr, req.P2pToCapabilities, chain, registry)

	// make the nodes and register node
	nodeParams := ToNodeParams(t,
		req.NopToNodes,
		ToP2PToCapabilities(t, req.P2pToCapabilities, registry, registeredCapabilities),
	)

	AddNodes(t, lggr, chain, registry, nodeParams)

	// add the Dons
	addDons(t, lggr, chain, registry, capCache, req.Dons)

	return &SetupTestRegistryResponse{
		CapabilitiesRegistry: registry,
		Chain:                chain,
		RegistrySelector:     chain.Selector,
		CapabilityCache:      capCache,
		Nops:                 addNopsResp.Nops,
	}
}

// ToNodeParams transforms a map of node operators to nops and a map of node p2pID to capabilities
// must match the number of nodes.
// The order of returned nodeParams is deterministic and sorted by node operator name.
func ToNodeParams(t *testing.T,
	nop2Nodes map[capabilities_registry.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc,
	p2pToCapabilities map[p2pkey.PeerID][][32]byte,
) []capabilities_registry.CapabilitiesRegistryNodeParams {
	t.Helper()

	var nodeParams []capabilities_registry.CapabilitiesRegistryNodeParams
	var i uint32
	// deterministic order
	// get the keys of the map and sort them
	var keys []capabilities_registry.CapabilitiesRegistryNodeOperator
	for k := range maps.Keys(nop2Nodes) {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Name < keys[j].Name
	})
	for _, k := range keys {
		p2pSignerEncs, ok := nop2Nodes[k]
		require.True(t, ok, "missing node operator %s", k.Name)
		for _, p2pSignerEnc := range p2pSignerEncs {
			_, exists := p2pToCapabilities[p2pSignerEnc.P2PKey]
			require.True(t, exists, "missing capabilities for p2pID %s", p2pSignerEnc.P2PKey)

			nodeParams = append(nodeParams, capabilities_registry.CapabilitiesRegistryNodeParams{
				Signer:              p2pSignerEnc.Signer,
				P2pId:               p2pSignerEnc.P2PKey,
				EncryptionPublicKey: p2pSignerEnc.EncryptionPublicKey,
				HashedCapabilityIds: p2pToCapabilities[p2pSignerEnc.P2PKey],
				NodeOperatorId:      i + 1, // one-indexed
			})
		}
		i++
	}

	return nodeParams
}

func ToNodeOps(
	t *testing.T,
	nop2Nodes map[capabilities_registry.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc,
) []capabilities_registry.CapabilitiesRegistryNodeOperator {
	t.Helper()

	nops := make([]capabilities_registry.CapabilitiesRegistryNodeOperator, 0)
	for nop := range nop2Nodes {
		nops = append(nops, nop)
	}

	sort.Slice(nops, func(i, j int) bool {
		return nops[i].Name < nops[j].Name
	})
	return nops
}

// MustAddCapabilities adds the capabilities to the registry and returns the registered capabilities
// if the capability is already registered, this call will fail.
func MustAddCapabilities(
	t *testing.T,
	lggr logger.Logger,
	in map[p2pkey.PeerID][]capabilities_registry.CapabilitiesRegistryCapability,
	chain deployment.Chain,
	registry *capabilities_registry.CapabilitiesRegistry,
) ([]internal.RegisteredCapability, *CapabilityCache) {
	t.Helper()
	cache := NewCapabiltyCache(t, registry)
	var capabilities []capabilities_registry.CapabilitiesRegistryCapability
	for _, caps := range in {
		capabilities = append(capabilities, caps...)
	}

	registeredCapabilities := cache.AddCapabilities(lggr, chain, registry, capabilities)
	expectedDeduped := make(map[capabilities_registry.CapabilitiesRegistryCapability]struct{})
	for _, c := range capabilities {
		expectedDeduped[c] = struct{}{}
	}
	require.Len(t, registeredCapabilities, len(expectedDeduped))
	return registeredCapabilities, cache
}

// GetRegisteredCapabilities returns the registered capabilities for the given capabilities.  Each
// capability must exist on the cache already.
func GetRegisteredCapabilities(
	t *testing.T,
	lggr logger.Logger,
	in map[p2pkey.PeerID][]capabilities_registry.CapabilitiesRegistryCapability,
	cache *CapabilityCache,
) []internal.RegisteredCapability {
	t.Helper()

	var capabilities []capabilities_registry.CapabilitiesRegistryCapability
	for _, caps := range in {
		capabilities = append(capabilities, caps...)
	}

	registeredCapabilities := make([]internal.RegisteredCapability, 0)
	for _, c := range capabilities {
		id, exists := cache.Get(c)
		require.True(t, exists, "capability not found in cache %v", c)
		registeredCapabilities = append(registeredCapabilities, internal.RegisteredCapability{
			CapabilitiesRegistryCapability: c,
			ID:                             id,
		})
	}

	return registeredCapabilities
}

func ToP2PToCapabilities(
	t *testing.T,
	in map[p2pkey.PeerID][]capabilities_registry.CapabilitiesRegistryCapability,
	registry *capabilities_registry.CapabilitiesRegistry,
	caps []internal.RegisteredCapability,
) map[p2pkey.PeerID][][32]byte {
	t.Helper()
	out := make(map[p2pkey.PeerID][][32]byte)
	for p2pID := range in {
		out[p2pID] = mustCapabilityIds(t, registry, caps)
	}
	return out
}

func deployCapReg(t *testing.T, chain deployment.Chain) *capabilities_registry.CapabilitiesRegistry {
	capabilitiesRegistryDeployer, err := internal.NewCapabilitiesRegistryDeployer()
	require.NoError(t, err)
	_, err = capabilitiesRegistryDeployer.Deploy(internal.DeployRequest{Chain: chain})
	require.NoError(t, err)
	return capabilitiesRegistryDeployer.Contract()
}

func addNops(t *testing.T, lggr logger.Logger, chain deployment.Chain, registry *capabilities_registry.CapabilitiesRegistry, nops []capabilities_registry.CapabilitiesRegistryNodeOperator) *internal.RegisterNOPSResponse {
	env := &deployment.Environment{
		Logger: lggr,
		Chains: map[uint64]deployment.Chain{
			chain.Selector: chain,
		},
		ExistingAddresses: deployment.NewMemoryAddressBookFromMap(map[uint64]map[string]deployment.TypeAndVersion{
			chain.Selector: {
				registry.Address().String(): deployment.TypeAndVersion{
					Type:    internal.CapabilitiesRegistry,
					Version: deployment.Version1_0_0,
				},
			},
		}),
	}
	resp, err := internal.RegisterNOPS(context.TODO(), lggr, internal.RegisterNOPSRequest{
		Env:                   env,
		RegistryChainSelector: chain.Selector,
		Nops:                  nops,
	})
	require.NoError(t, err)
	return resp
}

func AddNodes(
	t *testing.T,
	lggr logger.Logger,
	chain deployment.Chain,
	registry *capabilities_registry.CapabilitiesRegistry,
	nodes []capabilities_registry.CapabilitiesRegistryNodeParams,
) {
	tx, err := registry.AddNodes(chain.DeployerKey, nodes)
	if err != nil {
		err2 := deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)
		require.Fail(t, fmt.Sprintf("failed to call AddNodes: %s:  %s", err, err2))
	}
	_, err = chain.Confirm(tx)
	require.NoError(t, err)
}

func addDons(
	t *testing.T,
	_ logger.Logger,
	chain deployment.Chain,
	registry *capabilities_registry.CapabilitiesRegistry,
	capCache *CapabilityCache,
	dons []Don,
) {
	for _, don := range dons {
		acceptsWorkflows := false
		// lookup the capabilities
		var capConfigs []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration
		for _, ccfg := range don.CapabilityConfigs {
			var cc = capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
				CapabilityId: [32]byte{},
				Config:       ccfg.Config,
			}
			if cc.Config == nil {
				cc.Config = defaultCapConfig(t, ccfg.Capability)
			}
			var exists bool
			cc.CapabilityId, exists = capCache.Get(ccfg.Capability)
			require.True(t, exists, "capability not found in cache %v", ccfg.Capability)
			capConfigs = append(capConfigs, cc)
			if ccfg.Capability.CapabilityType == 2 { // ocr3 capabilities
				acceptsWorkflows = true
			}
		}
		// add the don
		isPublic := true
		f := len(don.P2PIDs)/3 + 1
		tx, err := registry.AddDON(chain.DeployerKey, internal.PeerIDsToBytes(don.P2PIDs), capConfigs, isPublic, acceptsWorkflows, uint8(f))
		if err != nil {
			err2 := deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)
			require.Fail(t, fmt.Sprintf("failed to call AddDON: %s:  %s", err, err2))
		}
		_, err = chain.Confirm(tx)
		require.NoError(t, err)
	}
}

func defaultCapConfig(t *testing.T, cap capabilities_registry.CapabilitiesRegistryCapability) []byte {
	empty := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
	}
	emptyb, err := proto.Marshal(empty)
	require.NoError(t, err)
	return emptyb
}

// CapabilityCache tracks registered capabilities by name
type CapabilityCache struct {
	t        *testing.T
	nameToId map[string][32]byte
}

func NewCapabiltyCache(t *testing.T, registry *capabilities_registry.CapabilitiesRegistry) *CapabilityCache {
	cache := &CapabilityCache{
		t:        t,
		nameToId: make(map[string][32]byte),
	}
	caps, err := registry.GetCapabilities(nil)
	require.NoError(t, err)
	for _, capb := range caps {
		c := capabilities_registry.CapabilitiesRegistryCapability{
			LabelledName:   capb.LabelledName,
			Version:        capb.Version,
			CapabilityType: capb.CapabilityType,
		}
		id, err := registry.GetHashedCapabilityId(&bind.CallOpts{}, c.LabelledName, c.Version)
		require.NoError(t, err)
		cache.nameToId[internal.CapabilityID(c)] = id
	}
	return cache
}
func (cc *CapabilityCache) Get(c capabilities_registry.CapabilitiesRegistryCapability) ([32]byte, bool) {
	id, exists := cc.nameToId[internal.CapabilityID(c)]
	return id, exists
}

// AddCapabilities adds the capabilities to the registry and returns the registered capabilities
// if the capability is already registered, it will not be re-registered
// if duplicate capabilities are passed, they will be deduped
func (cc *CapabilityCache) AddCapabilities(_ logger.Logger, chain deployment.Chain, registry *capabilities_registry.CapabilitiesRegistry, capabilities []capabilities_registry.CapabilitiesRegistryCapability) []internal.RegisteredCapability {
	t := cc.t
	var out []internal.RegisteredCapability
	// get the registered capabilities & dedup
	seen := make(map[capabilities_registry.CapabilitiesRegistryCapability]struct{})
	var toRegister []capabilities_registry.CapabilitiesRegistryCapability
	for _, c := range capabilities {
		id, cached := cc.nameToId[internal.CapabilityID(c)]
		if cached {
			out = append(out, internal.RegisteredCapability{
				CapabilitiesRegistryCapability: c,
				ID:                             id,
				Config:                         GetDefaultCapConfig(t, c),
			})
			continue
		}
		// dedup
		if _, exists := seen[c]; !exists {
			seen[c] = struct{}{}
			toRegister = append(toRegister, c)
		}
	}
	if len(toRegister) == 0 {
		return out
	}
	tx, err := registry.AddCapabilities(chain.DeployerKey, toRegister)
	if err != nil {
		err2 := deployment.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)
		require.Fail(t, fmt.Sprintf("failed to call AddCapabilities: %s:  %s", err, err2))
	}
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	// get the registered capabilities
	for _, capb := range toRegister {
		capb := capb
		id, err := registry.GetHashedCapabilityId(&bind.CallOpts{}, capb.LabelledName, capb.Version)
		require.NoError(t, err)
		out = append(out, internal.RegisteredCapability{
			CapabilitiesRegistryCapability: capb,
			ID:                             id,
			Config:                         GetDefaultCapConfig(t, capb),
		})
		// cache the id
		cc.nameToId[internal.CapabilityID(capb)] = id
	}
	return out
}

func testChain(t *testing.T) deployment.Chain {
	chains, _ := memory.NewMemoryChains(t, 1, 5)
	var chain deployment.Chain
	for _, c := range chains {
		chain = c
		break
	}
	require.NotEmpty(t, chain)
	return chain
}

func capabilityIds(registry *capabilities_registry.CapabilitiesRegistry, rcs []internal.RegisteredCapability) ([][32]byte, error) {
	out := make([][32]byte, len(rcs))
	for i := range rcs {
		id, err := registry.GetHashedCapabilityId(&bind.CallOpts{}, rcs[i].LabelledName, rcs[i].Version)
		if err != nil {
			return nil, fmt.Errorf("failed to get capability id: %w", err)
		}
		out[i] = id
	}
	return out, nil
}

func mustCapabilityIds(t *testing.T, registry *capabilities_registry.CapabilitiesRegistry, rcs []internal.RegisteredCapability) [][32]byte {
	t.Helper()
	out, err := capabilityIds(registry, rcs)
	require.NoError(t, err)
	return out
}

func MustCapabilityID(t *testing.T, registry *capabilities_registry.CapabilitiesRegistry, c capabilities_registry.CapabilitiesRegistryCapability) [32]byte {
	t.Helper()
	id, err := registry.GetHashedCapabilityId(&bind.CallOpts{}, c.LabelledName, c.Version)
	require.NoError(t, err)
	return id
}

func GetDefaultCapConfig(t *testing.T, capability capabilities_registry.CapabilitiesRegistryCapability) *capabilitiespb.CapabilityConfig {
	t.Helper()
	defaultCfg := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
	}
	switch capability.CapabilityType {
	case uint8(0): // trigger
		defaultCfg.RemoteConfig = &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh:     durationpb.New(20 * time.Second),
				RegistrationExpiry:      durationpb.New(60 * time.Second),
				MinResponsesToAggregate: uint32(10),
			},
		}
	case uint8(3): // target
		defaultCfg.RemoteConfig = &capabilitiespb.CapabilityConfig_RemoteTargetConfig{
			RemoteTargetConfig: &capabilitiespb.RemoteTargetConfig{
				RequestHashExcludedAttributes: []string{"signed_report.Signatures"},
			},
		}
	case uint8(2): // consensus
	default:
	}
	return defaultCfg
}
