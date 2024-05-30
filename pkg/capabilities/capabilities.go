package capabilities

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"golang.org/x/mod/semver"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

// CapabilityType is an enum for the type of capability.
type CapabilityType int

// CapabilityType enum values.
const (
	CapabilityTypeTrigger CapabilityType = iota
	CapabilityTypeAction
	CapabilityTypeConsensus
	CapabilityTypeTarget
)

// String returns a string representation of CapabilityType
func (c CapabilityType) String() string {
	switch c {
	case CapabilityTypeTrigger:
		return "trigger"
	case CapabilityTypeAction:
		return "action"
	case CapabilityTypeConsensus:
		return "consensus"
	case CapabilityTypeTarget:
		return "target"
	}

	// Panic as this should be unreachable.
	panic("unknown capability type")
}

// IsValid checks if the capability type is valid.
func (c CapabilityType) IsValid() error {
	switch c {
	case CapabilityTypeTrigger,
		CapabilityTypeAction,
		CapabilityTypeConsensus,
		CapabilityTypeTarget:
		return nil
	}

	return fmt.Errorf("invalid capability type: %s", c)
}

// CapabilityResponse is a struct for the Execute response of a capability.
type CapabilityResponse struct {
	Value values.Value
	Err   error
}

type RequestMetadata struct {
	WorkflowID          string
	WorkflowOwner       string
	WorkflowExecutionID string
}

type RegistrationMetadata struct {
	WorkflowID    string
	WorkflowOwner string
}

// CapabilityRequest is a struct for the Execute request of a capability.
type CapabilityRequest struct {
	Metadata RequestMetadata
	Config   *values.Map
	Inputs   *values.Map
}

type TriggerEvent struct {
	TriggerType string
	ID          string
	Timestamp   string
	// Trigger-specific payload+metadata
	Metadata values.Value
	Payload  values.Value
}

type RegisterToWorkflowRequest struct {
	Metadata RegistrationMetadata
	Config   *values.Map
}

type UnregisterFromWorkflowRequest struct {
	Metadata RegistrationMetadata
	Config   *values.Map
}

// CallbackExecutable is an interface for executing a capability.
type CallbackExecutable interface {
	RegisterToWorkflow(ctx context.Context, request RegisterToWorkflowRequest) error
	UnregisterFromWorkflow(ctx context.Context, request UnregisterFromWorkflowRequest) error
	// Capability must respect context.Done and cleanup any request specific resources
	// when the context is cancelled. When a request has been completed the capability
	// is also expected to close the callback channel.
	// Request specific configuration is passed in via the request parameter.
	// A successful response must always return a value. An error is assumed otherwise.
	// The intent is to make the API explicit.
	Execute(ctx context.Context, request CapabilityRequest) (<-chan CapabilityResponse, error)
}

type Validatable interface {
	// ValidateSchema returns the JSON schema for the capability.
	//
	// This schema includes the configuration, input and output schemas.
	Schema() (string, error)
}

// BaseCapability interface needs to be implemented by all capability types.
// Capability interfaces are intentionally duplicated to allow for an easy change
// or extension in the future.
type BaseCapability interface {
	Info(ctx context.Context) (CapabilityInfo, error)
}

type TriggerExecutable interface {
	RegisterTrigger(ctx context.Context, request CapabilityRequest) (<-chan CapabilityResponse, error)
	UnregisterTrigger(ctx context.Context, request CapabilityRequest) error
}

// TriggerCapability interface needs to be implemented by all trigger capabilities.
type TriggerCapability interface {
	BaseCapability
	TriggerExecutable
}

// CallbackCapability is the interface implemented by action, consensus and target
// capabilities. This interface is useful when trying to capture capabilities of varying types.
type CallbackCapability interface {
	BaseCapability
	CallbackExecutable
}

// ActionCapability interface needs to be implemented by all action capabilities.
type ActionCapability interface {
	CallbackCapability
}

// ConsensusCapability interface needs to be implemented by all consensus capabilities.
type ConsensusCapability interface {
	CallbackCapability
}

// TargetsCapability interface needs to be implemented by all target capabilities.
type TargetCapability interface {
	CallbackCapability
}

type DONConfig struct {
	SharedSecret [16]byte
}

type DON struct {
	ID      string
	Members []p2ptypes.PeerID
	F       uint8

	Config DONConfig
}

// CapabilityInfo is a struct for the info of a capability.
type CapabilityInfo struct {
	// The capability ID is a fully qualified identifier for the capability.
	//
	// It takes the form of `{name}:{label1_key}_{labe1_value}:{label2_key}_{label2_value}@{version}`
	//
	// The labels within the ID are ordered alphanumerically.
	ID             string
	CapabilityType CapabilityType
	Description    string
	Version        string
	DON            *DON
}

// Info returns the info of the capability.
func (c CapabilityInfo) Info(ctx context.Context) (CapabilityInfo, error) {
	return c, nil
}

var idRegex = regexp.MustCompile(`[a-z0-9_\-:]`)

const (
	// TODO: this length was largely picked arbitrarily.
	// Consider what a realistic/desirable value should be.
	// See: https://smartcontract-it.atlassian.net/jira/software/c/projects/KS/boards/182
	idMaxLength = 128
)

// NewCapabilityInfo returns a new CapabilityInfo.
func NewCapabilityInfo(
	id string,
	capabilityType CapabilityType,
	description string,
	version string,
	don *DON,
) (CapabilityInfo, error) {
	if len(id) > idMaxLength {
		return CapabilityInfo{}, fmt.Errorf("invalid id: %s exceeds max length %d", id, idMaxLength)
	}
	if !idRegex.MatchString(id) {
		return CapabilityInfo{}, fmt.Errorf("invalid id: %s. Allowed: %s", id, idRegex)
	}

	if ok := semver.IsValid(version); !ok {
		return CapabilityInfo{}, fmt.Errorf("invalid version: %+v", version)
	}

	if err := capabilityType.IsValid(); err != nil {
		return CapabilityInfo{}, err
	}

	return CapabilityInfo{
		ID:             id,
		CapabilityType: capabilityType,
		Description:    description,
		Version:        version,
		DON:            don,
	}, nil
}

// MustNewCapabilityInfo returns a new CapabilityInfo,
// `panic`ing if we could not instantiate a CapabilityInfo.
func MustNewCapabilityInfo(
	id string,
	capabilityType CapabilityType,
	description string,
	version string,
	don *DON,
) CapabilityInfo {
	c, err := NewCapabilityInfo(id, capabilityType, description, version, don)
	if err != nil {
		panic(err)
	}

	return c
}

// TODO: this timeout was largely picked arbitrarily.
// Consider what a realistic/desirable value should be.
// See: https://smartcontract-it.atlassian.net/jira/software/c/projects/KS/boards/182
var maximumExecuteTimeout = 60 * time.Second

// ExecuteSync executes a capability synchronously.
// We are not handling a case where a capability panics and crashes.
// There is default timeout of 10 seconds. If a capability takes longer than
// that then it should be executed asynchronously.
func ExecuteSync(ctx context.Context, c CallbackExecutable, request CapabilityRequest) (*values.List, error) {
	ctxWithT, cancel := context.WithTimeout(ctx, maximumExecuteTimeout)
	defer cancel()

	responseCh, err := c.Execute(ctxWithT, request)
	if err != nil {
		return nil, fmt.Errorf("error executing capability: %w", err)
	}

	vs := make([]values.Value, 0)
outerLoop:
	for {
		select {
		case response, isOpen := <-responseCh:
			if !isOpen {
				break outerLoop
			}
			// An error means execution has been interrupted.
			// We'll return the value discarding values received
			// until now.
			if response.Err != nil {
				return nil, response.Err
			}

			vs = append(vs, response.Value)
		// Timeout when a capability exceeds maximum permitted execution time or the caller cancels the context and does not close the channel.
		case <-ctxWithT.Done():
			return nil, fmt.Errorf("context timed out after %f seconds", maximumExecuteTimeout.Seconds())
		}
	}

	// If the capability did not return any values, we deem it as an error.
	// The intent is for the API to be explicit.
	if len(vs) == 0 {
		return nil, errors.New("capability did not return any values")
	}

	return &values.List{Underlying: vs}, nil
}
