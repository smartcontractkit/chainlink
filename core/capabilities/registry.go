package capabilities

import (
	"context"
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// Registry is a struct for the registry of capabilities.
// Registry is safe for concurrent use.
type Registry struct {
	core.UnimplementedCapabilitiesRegistryMetadata
	core.CapabilitiesRegistryBase

	metadataRegistry core.CapabilitiesRegistryMetadata
	lggr             logger.Logger
	mu               sync.RWMutex
}

func (r *Registry) LocalNode(ctx context.Context) (capabilities.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.metadataRegistry == nil {
		return capabilities.Node{}, errors.New("metadataRegistry information not available")
	}

	return r.metadataRegistry.LocalNode(ctx)
}

func (r *Registry) NodeByPeerID(ctx context.Context, peerID p2ptypes.PeerID) (capabilities.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.metadataRegistry == nil {
		return capabilities.Node{}, errors.New("metadataRegistry information not available")
	}
	return r.metadataRegistry.NodeByPeerID(ctx, peerID)
}

func (r *Registry) ConfigForCapability(ctx context.Context, capabilityID string, donID uint32) (capabilities.CapabilityConfiguration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.metadataRegistry == nil {
		return capabilities.CapabilityConfiguration{}, errors.New("metadataRegistry information not available")
	}

	return r.metadataRegistry.ConfigForCapability(ctx, capabilityID, donID)
}

func (r *Registry) DONsForCapability(ctx context.Context, capabilityID string) ([]capabilities.DONWithNodes, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.metadataRegistry == nil {
		return nil, errors.New("metadataRegistry information not available")
	}

	return r.metadataRegistry.DONsForCapability(ctx, capabilityID)
}

// SetLocalRegistry sets a local copy of the offchain registry for the registry to use.
// This is only public for testing purposes; the only production use should be from the CapabilitiesLauncher.
func (r *Registry) SetLocalRegistry(lr core.CapabilitiesRegistryMetadata) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metadataRegistry = lr
}

// NewRegistry returns a new Registry.
func NewRegistry(lggr logger.Logger) *Registry {
	return &Registry{
		CapabilitiesRegistryBase: registry.NewBaseRegistry(lggr),
		lggr:                     logger.Named(lggr, "CapabilitiesRegistry"),
	}
}

// TestMetadataRegistry is a test implementation of the metadataRegistry
// interface. It is used when ExternalCapabilitiesRegistry is not available.
type TestMetadataRegistry struct {
	core.UnimplementedCapabilitiesRegistryMetadata
	// WorkflowDONF allows local CRE to override the synthetic workflow DON fault
	// tolerance for compatibility paths that still expect a multi-signer shape.
	WorkflowDONF uint8
}

const (
	testWorkflowDONID            = 1
	testWorkflowDONConfigVersion = 1
)

func (t *TestMetadataRegistry) LocalNode(ctx context.Context) (capabilities.Node, error) {
	peerID := p2ptypes.PeerID{}
	return capabilities.Node{
		PeerID:         &peerID,
		WorkflowDON:    newTestWorkflowDON(peerID, t.WorkflowDONF),
		CapabilityDONs: []capabilities.DON{},
	}, nil
}

func newTestWorkflowDON(peerID p2ptypes.PeerID, faultTolerance uint8) capabilities.DON {
	return capabilities.DON{
		ID:            testWorkflowDONID,
		ConfigVersion: testWorkflowDONConfigVersion,
		Members: []p2ptypes.PeerID{
			peerID,
		},
		F:                faultTolerance,
		IsPublic:         false,
		AcceptsWorkflows: true,
	}
}

func (t *TestMetadataRegistry) NodeByPeerID(ctx context.Context, _ p2ptypes.PeerID) (capabilities.Node, error) {
	return t.LocalNode(ctx)
}

func (t *TestMetadataRegistry) ConfigForCapability(ctx context.Context, capabilityID string, donID uint32) (capabilities.CapabilityConfiguration, error) {
	return capabilities.CapabilityConfiguration{}, nil
}

func (t *TestMetadataRegistry) DONsForCapability(ctx context.Context, capabilityID string) ([]capabilities.DONWithNodes, error) {
	return []capabilities.DONWithNodes{}, nil
}
