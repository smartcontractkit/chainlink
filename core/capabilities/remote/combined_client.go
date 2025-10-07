package remote

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

// CombinedClient represents a remote capability V2 accessed from a local node (by the Engine).
// The capability can have multiple methods, each one being a trigger or an executable.
// The CombinedClient holds method-specific shims for each method and forwards capability API calls
// to them. Responses are passed directly to method-specific shims from the Dispatcher.
type CombinedClient interface {
	capabilities.ExecutableAndTriggerCapability
	SetTriggerSubscriber(method string, subscriber capabilities.TriggerCapability)
	SetExecutableClient(method string, client capabilities.ExecutableCapability)
	GetTriggerSubscriber(method string) capabilities.TriggerCapability
	GetExecutableClient(method string) capabilities.ExecutableCapability
}

type combinedClient struct {
	info               capabilities.CapabilityInfo
	triggerSubscribers map[string]capabilities.TriggerCapability
	executableClients  map[string]capabilities.ExecutableCapability
	mu                 sync.RWMutex
}

var _ CombinedClient = &combinedClient{}

func (c *combinedClient) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return c.info, nil
}

func (c *combinedClient) RegisterTrigger(ctx context.Context, request capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	c.mu.RLock()
	subscriber, ok := c.triggerSubscribers[request.Method]
	c.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("method %s not defined", request.Method)
	}
	return subscriber.RegisterTrigger(ctx, request)
}

func (c *combinedClient) UnregisterTrigger(ctx context.Context, request capabilities.TriggerRegistrationRequest) error {
	c.mu.RLock()
	subscriber, ok := c.triggerSubscribers[request.Method]
	c.mu.RUnlock()

	if !ok {
		return fmt.Errorf("method %s not defined", request.Method)
	}
	return subscriber.UnregisterTrigger(ctx, request)
}

func (c *combinedClient) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return errors.New("RegisterToWorkflow is not supported by remote capabilities")
}

func (c *combinedClient) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return errors.New("UnregisterFromWorkflow is not supported by remote capabilities")
}

func (c *combinedClient) Execute(ctx context.Context, request capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	c.mu.RLock()
	client, ok := c.executableClients[request.Method]
	c.mu.RUnlock()

	if !ok {
		return capabilities.CapabilityResponse{}, fmt.Errorf("method %s not defined", request.Method)
	}
	return client.Execute(ctx, request)
}

func NewCombinedClient(info capabilities.CapabilityInfo) *combinedClient {
	return &combinedClient{
		info:               info,
		triggerSubscribers: make(map[string]capabilities.TriggerCapability),
		executableClients:  make(map[string]capabilities.ExecutableCapability),
	}
}

func (c *combinedClient) SetTriggerSubscriber(method string, subscriber capabilities.TriggerCapability) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.triggerSubscribers[method] = subscriber
}

func (c *combinedClient) SetExecutableClient(method string, client capabilities.ExecutableCapability) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.executableClients[method] = client
}

func (c *combinedClient) GetTriggerSubscriber(method string) capabilities.TriggerCapability {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.triggerSubscribers[method]
}

func (c *combinedClient) GetExecutableClient(method string) capabilities.ExecutableCapability {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.executableClients[method]
}
