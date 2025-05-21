package capabilities

import (
	"context"
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

type metadataRegistry interface {
	LocalNode(ctx context.Context) (capabilities.Node, error)
	ConfigForCapability(ctx context.Context, capabilityID string, donID uint32) (registrysyncer.CapabilityConfiguration, error)
}

// Registry is a struct for the registry of capabilities.
// Registry is safe for concurrent use.
type Registry struct {
	metadataRegistry metadataRegistry
	lggr             logger.Logger
	core.CapabilitiesRegistryBase
	mu sync.RWMutex
}

func (r *Registry) LocalNode(ctx context.Context) (capabilities.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.metadataRegistry == nil {
		return capabilities.Node{}, errors.New("metadataRegistry information not available")
	}

	return r.metadataRegistry.LocalNode(ctx)
}

func (r *Registry) ConfigForCapability(ctx context.Context, capabilityID string, donID uint32) (capabilities.CapabilityConfiguration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.metadataRegistry == nil {
		return capabilities.CapabilityConfiguration{}, errors.New("metadataRegistry information not available")
	}

	cfc, err := r.metadataRegistry.ConfigForCapability(ctx, capabilityID, donID)
	if err != nil {
		return capabilities.CapabilityConfiguration{}, err
	}

	return unmarshalCapabilityConfig(cfc.Config)
}

// SetLocalRegistry sets a local copy of the offchain registry for the registry to use.
// This is only public for testing purposes; the only production use should be from the CapabilitiesLauncher.
func (r *Registry) SetLocalRegistry(lr metadataRegistry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metadataRegistry = lr
}

// NewRegistry returns a new Registry.
func NewRegistry(lggr logger.Logger) *Registry {
	return &Registry{
		CapabilitiesRegistryBase: registry.NewBaseRegistry(lggr),
		lggr:                     lggr.Named("CapabilitiesRegistry"),
	}
}

// TestMetadataRegistry is a test implementation of the metadataRegistry
// interface. It is used when ExternalCapabilitiesRegistry is not available.
type TestMetadataRegistry struct{}

func (t *TestMetadataRegistry) LocalNode(ctx context.Context) (capabilities.Node, error) {
	peerID := p2ptypes.PeerID{}
	workflowDON := capabilities.DON{
		ID:            1,
		ConfigVersion: 1,
		Members: []p2ptypes.PeerID{
			peerID,
		},
		F:                0,
		IsPublic:         false,
		AcceptsWorkflows: true,
	}
	return capabilities.Node{
		PeerID:         &peerID,
		WorkflowDON:    workflowDON,
		CapabilityDONs: []capabilities.DON{},
	}, nil
}

func (t *TestMetadataRegistry) ConfigForCapability(ctx context.Context, capabilityID string, donID uint32) (registrysyncer.CapabilityConfiguration, error) {
	return registrysyncer.CapabilityConfiguration{}, nil
}
