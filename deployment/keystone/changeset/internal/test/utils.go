package test

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"

	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	internal "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type SetupTestRegistryRequest struct {
	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability
	NopToNodes        map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc
	DonToNodes        map[string][]*internal.P2PSignerEnc
}

type SetupTestRegistryResponse struct {
	Registry         *kcr.CapabilitiesRegistry
	Chain            deployment.Chain
	RegistrySelector uint64
}

func SetupTestRegistry(t *testing.T, lggr logger.Logger, req *SetupTestRegistryRequest) *SetupTestRegistryResponse {
	chain := testChain(t)
	// deploy the registry
	registry := deployCapReg(t, lggr, chain)
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

	// add capabilities to registry
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

	// add the nodes with the phony capabilities. cannot register a node without a capability and capability must exist
	// to do this make an initial phony request and extract the node params
	initialp2pToCapabilities := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
	for p2pID := range req.P2pToCapabilities {
		initialp2pToCapabilities[p2pID] = vanillaCapabilities(registeredCapabilities)
	}
	phonyRequest := &internal.UpdateNodesRequest{
		Chain:             chain,
		Registry:          registry,
		P2pToCapabilities: req.P2pToCapabilities,
		NopToNodes:        req.NopToNodes,
	}
	nodeParams, err := phonyRequest.NodeParams()
	require.NoError(t, err)
	addNodes(t, lggr, chain, registry, nodeParams)
	return &SetupTestRegistryResponse{
		Registry:         registry,
		Chain:            chain,
		RegistrySelector: chain.Selector,
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
