package remote

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

// CombinedClient represents a remote capability V2 accessed from a local node (by the Engine).
// The capability can have multiple methods, each one being a trigger or an executable.
// The CombinedClient holds method-specific shims for each method and forwards capability API calls
// to them. Responses are passed directly to method-specific shims from the Dispatcher.
type combinedClient struct {
	info               capabilities.CapabilityInfo
	triggerSubscribers map[string]capabilities.TriggerCapability
	executableClients  map[string]capabilities.ExecutableCapability
}

var _ capabilities.ExecutableAndTriggerCapability = &combinedClient{}

func (c *combinedClient) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return c.info, nil
}

func (c *combinedClient) RegisterTrigger(ctx context.Context, request capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	if _, ok := c.triggerSubscribers[request.Method]; !ok {
		return nil, fmt.Errorf("method %s not defined", request.Method)
	}
	return c.triggerSubscribers[request.Method].RegisterTrigger(ctx, request)
}

func (c *combinedClient) UnregisterTrigger(ctx context.Context, request capabilities.TriggerRegistrationRequest) error {
	if _, ok := c.triggerSubscribers[request.Method]; !ok {
		return fmt.Errorf("method %s not defined", request.Method)
	}
	return c.triggerSubscribers[request.Method].UnregisterTrigger(ctx, request)
}

func (c *combinedClient) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return errors.New("RegisterToWorkflow is not supported by remote capabilities")
}

func (c *combinedClient) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return errors.New("UnregisterFromWorkflow is not supported by remote capabilities")
}

func (c *combinedClient) Execute(ctx context.Context, request capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	if _, ok := c.executableClients[request.Method]; !ok {
		return capabilities.CapabilityResponse{}, fmt.Errorf("method %s not defined", request.Method)
	}
	return c.executableClients[request.Method].Execute(ctx, request)
}

func NewCombinedClient(info capabilities.CapabilityInfo) *combinedClient {
	return &combinedClient{
		info:               info,
		triggerSubscribers: make(map[string]capabilities.TriggerCapability),
		executableClients:  make(map[string]capabilities.ExecutableCapability),
	}
}

func (c *combinedClient) AddTriggerSubscriber(method string, subscriber capabilities.TriggerCapability) {
	c.triggerSubscribers[method] = subscriber
}

func (c *combinedClient) AddExecutableClient(method string, client capabilities.ExecutableCapability) {
	c.executableClients[method] = client
}
