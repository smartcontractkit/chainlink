package test

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	kslib "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type Don struct {
	Name              string
	P2PIDs            []p2pkey.PeerID
	CapabilityConfigs []kslib.CapabilityConfig
}

type SetupTestRegistryRequest struct {
	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability
	NopToNodes        map[kcr.CapabilitiesRegistryNodeOperator][]*kslib.P2PSignerEnc
	Dons              []Don
}

type SetupTestRegistryResponse struct {
	Registry *kcr.CapabilitiesRegistry
	Chain    deployment.Chain
}

// SetupTestRegistry sets up a test registry with the given capabilities and nodes
// It always creates a test chain and deploys a new registry it
// The configuration determines what, if anything, is added to the registry
func SetupTestRegistry(t *testing.T, lggr logger.Logger, req *SetupTestRegistryRequest) *SetupTestRegistryResponse {
	chain := testChain(t)
	// deploy the registry
	registry := deployCapReg(t, lggr, chain)

	// Register capabilities
	capCache := NewCapabiltyCache(t)
	var capabilities []kcr.CapabilitiesRegistryCapability
	for _, caps := range req.P2pToCapabilities {
		capabilities = append(capabilities, caps...)
	}
	registeredCapabilities := capCache.AddCapabilities(lggr, chain, registry, capabilities)
	expectedDeduped := make(map[kcr.CapabilitiesRegistryCapability]struct{})
	for _, cap := range capabilities {
		expectedDeduped[cap] = struct{}{}
	}
	require.Len(t, registeredCapabilities, len(expectedDeduped))

	// Register NOPs and Nodes
	// convert req to nodeoperators
	nops := make([]kcr.CapabilitiesRegistryNodeOperator, 0)
	for nop := range req.NopToNodes {
		nops = append(nops, nop)
	}
	sort.Slice(nops, func(i, j int) bool {
		return nops[i].Name < nops[j].Name
	})
	addNopsResp := addNops(t, lggr, chain, registry, nops)
	require.Len(t, addNopsResp.Nops, len(nops))

	// add the nodes with the phony capabilities. cannot register a node without a capability and capability must exist
	phonyRequest := &kslib.UpdateNodesRequest{
		Chain:             chain,
		Registry:          registry,
		P2pToCapabilities: req.P2pToCapabilities,
		NopToNodes:        req.NopToNodes,
	}
	nodeParams, err := phonyRequest.NodeParams()
	require.NoError(t, err)
	addNodes(t, lggr, chain, registry, nodeParams)

	// add dons
	addDons(t, lggr, chain, registry, capCache, req.Dons)

	return &SetupTestRegistryResponse{
		Registry: registry,
		Chain:    chain,
	}
}

func deployCapReg(t *testing.T, lggr logger.Logger, chain deployment.Chain) *kcr.CapabilitiesRegistry {
	capabilitiesRegistryDeployer := kslib.NewCapabilitiesRegistryDeployer(lggr)
	_, err := capabilitiesRegistryDeployer.Deploy(kslib.DeployRequest{Chain: chain})
	require.NoError(t, err)
	return capabilitiesRegistryDeployer.Contract()
}

func addNops(t *testing.T, lggr logger.Logger, chain deployment.Chain, registry *kcr.CapabilitiesRegistry, nops []kcr.CapabilitiesRegistryNodeOperator) *kslib.RegisterNOPSResponse {
	resp, err := kslib.RegisterNOPS(context.TODO(), kslib.RegisterNOPSRequest{
		Chain:    chain,
		Registry: registry,
		Nops:     nops,
	})
	require.NoError(t, err)
	return resp
}

func addNodes(t *testing.T, lggr logger.Logger, chain deployment.Chain, registry *kcr.CapabilitiesRegistry, nodes []kcr.CapabilitiesRegistryNodeParams) {
	tx, err := registry.AddNodes(chain.DeployerKey, nodes)
	if err != nil {
		err2 := kslib.DecodeErr(kcr.CapabilitiesRegistryABI, err)
		require.Fail(t, fmt.Sprintf("failed to call AddNodes: %s:  %s", err, err2))
	}
	_, err = chain.Confirm(tx)
	require.NoError(t, err)
}

func addDons(t *testing.T, lggr logger.Logger, chain deployment.Chain, registry *kcr.CapabilitiesRegistry, cc *CapabilityCache, dons []Don) {
	for _, don := range dons {
		acceptsWorkflows := false
		// lookup the capabilities
		var capConfigs []kcr.CapabilitiesRegistryCapabilityConfiguration
		for _, ccfg := range don.CapabilityConfigs {
			if ccfg.Config == nil {
				ccfg.Config = defaultCapConfig(t, ccfg.Capability)
			}
			var exists bool
			ccfg.CapabilityId, exists = cc.Get(ccfg.Capability)
			require.True(t, exists, "capability not found in cache %v", ccfg.Capability)
			capConfigs = append(capConfigs, ccfg.CapabilitiesRegistryCapabilityConfiguration)
			if ccfg.Capability.CapabilityType == 2 { // ocr3 capabilities
				acceptsWorkflows = true
			}
		}
		// add the don
		isPublic := true
		f := len(don.P2PIDs)/3 + 1
		tx, err := registry.AddDON(chain.DeployerKey, peerIDsToBytes(don.P2PIDs), capConfigs, isPublic, acceptsWorkflows, uint8(f))
		if err != nil {
			err2 := kslib.DecodeErr(kcr.CapabilitiesRegistryABI, err)
			require.Fail(t, fmt.Sprintf("failed to call AddDON: %s:  %s", err, err2))
		}
		_, err = chain.Confirm(tx)
		require.NoError(t, err)
	}
}

func defaultCapConfig(t *testing.T, cap kcr.CapabilitiesRegistryCapability) []byte {
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

func NewCapabiltyCache(t *testing.T) *CapabilityCache {
	return &CapabilityCache{
		t:        t,
		nameToId: make(map[string][32]byte),
	}
}

// AddCapabilities adds the capabilities to the registry and returns the registered capabilities
// if the capability is already registered, it will not be re-registered
// if duplicate capabilities are passed, they will be deduped
func (cc *CapabilityCache) AddCapabilities(lggr logger.Logger, chain deployment.Chain, registry *kcr.CapabilitiesRegistry, capabilities []kcr.CapabilitiesRegistryCapability) []kslib.RegisteredCapability {
	t := cc.t
	var out []kslib.RegisteredCapability
	// get the registered capabilities & dedup
	seen := make(map[kcr.CapabilitiesRegistryCapability]struct{})
	var toRegister []kcr.CapabilitiesRegistryCapability
	for _, cap := range capabilities {
		id, cached := cc.nameToId[kslib.CapabilityID(cap)]
		if cached {
			out = append(out, kslib.RegisteredCapability{
				CapabilitiesRegistryCapability: cap,
				ID:                             id,
			})
			continue
		}
		// dedup
		if _, exists := seen[cap]; !exists {
			seen[cap] = struct{}{}
			toRegister = append(toRegister, cap)
		}
	}
	if len(toRegister) == 0 {
		return out
	}
	tx, err := registry.AddCapabilities(chain.DeployerKey, toRegister)
	if err != nil {
		err2 := kslib.DecodeErr(kcr.CapabilitiesRegistryABI, err)
		require.Fail(t, fmt.Sprintf("failed to call AddCapabilities: %s:  %s", err, err2))
	}
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	// get the registered capabilities
	for _, capb := range toRegister {
		capb := capb
		id, err := registry.GetHashedCapabilityId(&bind.CallOpts{}, capb.LabelledName, capb.Version)
		require.NoError(t, err)
		out = append(out, kslib.RegisteredCapability{
			CapabilitiesRegistryCapability: capb,
			ID:                             id,
		})
		// cache the id
		cc.nameToId[kslib.CapabilityID(capb)] = id
	}
	return out
}

func (cc *CapabilityCache) Get(cap kcr.CapabilitiesRegistryCapability) ([32]byte, bool) {
	id, exists := cc.nameToId[kslib.CapabilityID(cap)]
	return id, exists
}

func testChain(t *testing.T) deployment.Chain {
	chains := memory.NewMemoryChains(t, 1)
	var chain deployment.Chain
	for _, c := range chains {
		chain = c
		break
	}
	require.NotEmpty(t, chain)
	return chain
}

func vanillaCapabilities(rcs []kslib.RegisteredCapability) []kcr.CapabilitiesRegistryCapability {
	out := make([]kcr.CapabilitiesRegistryCapability, len(rcs))
	for i := range rcs {
		out[i] = rcs[i].CapabilitiesRegistryCapability
	}
	return out
}

func peerIDsToBytes(p2pIDs []p2pkey.PeerID) [][32]byte {
	out := make([][32]byte, len(p2pIDs))
	for i, p2pID := range p2pIDs {
		out[i] = p2pID
	}
	return out
}
