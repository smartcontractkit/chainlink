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
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"

	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	internal "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type Don struct {
	Name              string
	P2PIDs            []p2pkey.PeerID
	CapabilityConfigs []internal.CapabilityConfig
}
type SetupTestRegistryRequest struct {
	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability
	NopToNodes        map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc
	Dons              []Don
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

	// make the nodes and register node
	var nodeParams []kcr.CapabilitiesRegistryNodeParams
	initialp2pToCapabilities := make(map[p2pkey.PeerID][][32]byte)
	for p2pID := range req.P2pToCapabilities {
		initialp2pToCapabilities[p2pID] = mustCapabilityIds(t, registry, registeredCapabilities)
	}
	// create node with initial capabilities assigned to nop
	for i, nop := range nops {
		if _, exists := req.NopToNodes[nop]; !exists {
			require.Fail(t, "missing nopToNodes for %s", nop.Name)
		}
		for _, p2pSignerEnc := range req.NopToNodes[nop] {
			nodeParams = append(nodeParams, kcr.CapabilitiesRegistryNodeParams{
				Signer:              p2pSignerEnc.Signer,
				P2pId:               p2pSignerEnc.P2PKey,
				EncryptionPublicKey: p2pSignerEnc.EncryptionPublicKey,
				HashedCapabilityIds: initialp2pToCapabilities[p2pSignerEnc.P2PKey],
				NodeOperatorId:      uint32(i + 1), // nopid in contract is 1-indexed
			})
		}
	}
	addNodes(t, lggr, chain, registry, nodeParams)

	// add the Dons
	addDons(t, lggr, chain, registry, capCache, req.Dons)

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
	resp, err := kslib.RegisterNOPS(context.TODO(), lggr, kslib.RegisterNOPSRequest{
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

func addDons(t *testing.T, lggr logger.Logger, chain deployment.Chain, registry *kcr.CapabilitiesRegistry, capCache *CapabilityCache, dons []Don) {
	for _, don := range dons {
		acceptsWorkflows := false
		// lookup the capabilities
		var capConfigs []kcr.CapabilitiesRegistryCapabilityConfiguration
		for _, ccfg := range don.CapabilityConfigs {
			var cc = kcr.CapabilitiesRegistryCapabilityConfiguration{
				CapabilityId: [32]byte{},
				Config:       ccfg.Config,
			}
			if cc.Config == nil {
				cc.Config = defaultCapConfig(t, ccfg.Capability)
			}
			var exists bool
			//var cc kcr.CapabilitiesRegistryCapabilityConfiguration{}
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
func (cc *CapabilityCache) Get(cap kcr.CapabilitiesRegistryCapability) ([32]byte, bool) {
	id, exists := cc.nameToId[kslib.CapabilityID(cap)]
	return id, exists
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

func capabilityIds(registry *capabilities_registry.CapabilitiesRegistry, rcs []kslib.RegisteredCapability) ([][32]byte, error) {
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

func mustCapabilityIds(t *testing.T, registry *capabilities_registry.CapabilitiesRegistry, rcs []kslib.RegisteredCapability) [][32]byte {
	t.Helper()
	out, err := capabilityIds(registry, rcs)
	require.NoError(t, err)
	return out
}

func MustCapabilityId(t *testing.T, registry *capabilities_registry.CapabilitiesRegistry, cap capabilities_registry.CapabilitiesRegistryCapability) [32]byte {
	t.Helper()
	id, err := registry.GetHashedCapabilityId(&bind.CallOpts{}, cap.LabelledName, cap.Version)
	require.NoError(t, err)
	return id
}
