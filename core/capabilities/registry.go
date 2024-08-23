package capabilities

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

var (
	ErrCapabilityAlreadyExists = errors.New("capability already exists")
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
	m                map[string]capabilities.BaseCapability
	mu               sync.RWMutex
}

func (r *Registry) LocalNode(ctx context.Context) (capabilities.Node, error) {
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

// Get gets a capability from the registry.
func (r *Registry) Get(_ context.Context, id string) (capabilities.BaseCapability, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	r.lggr.Debugw("get capability", "id", id)
	c, ok := r.m[id]
	if !ok {
		return nil, fmt.Errorf("capability not found with id %s", id)
	}

	return c, nil
}

// GetTrigger gets a capability from the registry and tries to coerce it to the TriggerCapability interface.
func (r *Registry) GetTrigger(ctx context.Context, id string) (capabilities.TriggerCapability, error) {
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	tc, ok := c.(capabilities.TriggerCapability)
	if !ok {
		return nil, fmt.Errorf("capability with id: %s does not satisfy the capability interface", id)
	}

	return tc, nil
}

// GetAction gets a capability from the registry and tries to coerce it to the ActionCapability interface.
func (r *Registry) GetAction(ctx context.Context, id string) (capabilities.ActionCapability, error) {
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	ac, ok := c.(capabilities.ActionCapability)
	if !ok {
		return nil, fmt.Errorf("capability with id: %s does not satisfy the capability interface", id)
	}

	return ac, nil
}

// GetConsensus gets a capability from the registry and tries to coerce it to the ConsensusCapability interface.
func (r *Registry) GetConsensus(ctx context.Context, id string) (capabilities.ConsensusCapability, error) {
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	cc, ok := c.(capabilities.ConsensusCapability)
	if !ok {
		return nil, fmt.Errorf("capability with id: %s does not satisfy the capability interface", id)
	}

	return cc, nil
}

// GetTarget gets a capability from the registry and tries to coerce it to the TargetCapability interface.
func (r *Registry) GetTarget(ctx context.Context, id string) (capabilities.TargetCapability, error) {
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	tc, ok := c.(capabilities.TargetCapability)
	if !ok {
		return nil, fmt.Errorf("capability with id: %s does not satisfy the capability interface", id)
	}

	return tc, nil
}

// List lists all the capabilities in the registry.
func (r *Registry) List(_ context.Context) ([]capabilities.BaseCapability, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cl := []capabilities.BaseCapability{}
	for _, v := range r.m {
		cl = append(cl, v)
	}

	return cl, nil
}

// Add adds a capability to the registry.
func (r *Registry) Add(ctx context.Context, c capabilities.BaseCapability) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	info, err := c.Info(ctx)
	if err != nil {
		return err
	}

	switch info.CapabilityType {
	case capabilities.CapabilityTypeTrigger:
		_, ok := c.(capabilities.TriggerCapability)
		if !ok {
			return fmt.Errorf("trigger capability does not satisfy TriggerCapability interface")
		}
	case capabilities.CapabilityTypeAction:
		_, ok := c.(capabilities.ActionCapability)
		if !ok {
			return fmt.Errorf("action does not satisfy ActionCapability interface")
		}
	case capabilities.CapabilityTypeConsensus:
		_, ok := c.(capabilities.ConsensusCapability)
		if !ok {
			return fmt.Errorf("consensus capability does not satisfy ConsensusCapability interface")
		}
	case capabilities.CapabilityTypeTarget:
		_, ok := c.(capabilities.TargetCapability)
		if !ok {
			return fmt.Errorf("target capability does not satisfy TargetCapability interface")
		}
	default:
		return fmt.Errorf("unknown capability type: %s", info.CapabilityType)
	}

	id := info.ID
	_, ok := r.m[id]
	if ok {
		return fmt.Errorf("%w: id %s found in registry", ErrCapabilityAlreadyExists, id)
	}

	r.m[id] = c
	r.lggr.Infow("capability added", "id", id, "type", info.CapabilityType, "description", info.Description, "version", info.Version())
	return nil
}

// NewRegistry returns a new Registry.
func NewRegistry(lggr logger.Logger) *Registry {
	return &Registry{
		m:    map[string]capabilities.BaseCapability{},
		lggr: lggr.Named("CapabilitiesRegistry"),
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
