package capabilities

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

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
		return "report"
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
	WorkflowExecutionID string
}

type RegistrationMetadata struct {
	WorkflowID string
}

// CapabilityRequest is a struct for the Execute request of a capability.
type CapabilityRequest struct {
	Metadata RequestMetadata
	Config   *values.Map
	Inputs   *values.Map
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
	Execute(ctx context.Context, callback chan CapabilityResponse, request CapabilityRequest) error
}

// BaseCapability interface needs to be implemented by all capability types.
// Capability interfaces are intentionally duplicated to allow for an easy change
// or extension in the future.
type BaseCapability interface {
	Info() CapabilityInfo
}

// TriggerCapability interface needs to be implemented by all trigger capabilities.
type TriggerCapability interface {
	BaseCapability
	RegisterTrigger(ctx context.Context, callback chan CapabilityResponse, request CapabilityRequest) error
	UnregisterTrigger(ctx context.Context, request CapabilityRequest) error
}

// ActionCapability interface needs to be implemented by all action capabilities.
type ActionCapability interface {
	BaseCapability
	CallbackExecutable
}

// ConsensusCapability interface needs to be implemented by all consensus capabilities.
type ConsensusCapability interface {
	BaseCapability
	CallbackExecutable
}

// TargetsCapability interface needs to be implemented by all target capabilities.
type TargetCapability interface {
	BaseCapability
	CallbackExecutable
}

// CapabilityInfo is a struct for the info of a capability.
type CapabilityInfo struct {
	ID             string
	CapabilityType CapabilityType
	Description    string
	Version        string
}

// Info returns the info of the capability.
func (c CapabilityInfo) Info() CapabilityInfo {
	return c
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
	}, nil
}

// MustCapabilityInfo returns a new CapabilityInfo,
// panicking if we could not instantiate a CapabilityInfo.
func MustCapabilityInfo(
	id string,
	capabilityType CapabilityType,
	description string,
	version string,
) CapabilityInfo {
	c, err := NewCapabilityInfo(id, capabilityType, description, version)
	if err != nil {
		panic(err)
	}

	return c
}

// TODO: this timeout was largely picked arbitrarily.
// Consider what a realistic/desirable value should be.
// See: https://smartcontract-it.atlassian.net/jira/software/c/projects/KS/boards/182
var defaultExecuteTimeout = 10 * time.Second

// ExecuteSync executes a capability synchronously.
// We are not handling a case where a capability panics and crashes.
// There is default timeout of 10 seconds. If a capability takes longer than
// that then it should be executed asynchronously.
func ExecuteSync(ctx context.Context, c CallbackExecutable, request CapabilityRequest) (values.Value, error) {
	ctxWithT, cancel := context.WithTimeout(ctx, defaultExecuteTimeout)
	defer cancel()

	callback := make(chan CapabilityResponse)
	sec := make(chan error)

	go func(innerCtx context.Context, innerC CallbackExecutable, innerReq CapabilityRequest, innerCallback chan CapabilityResponse, errCh chan error) {
		setupErr := innerC.Execute(innerCtx, innerCallback, innerReq)
		sec <- setupErr
	}(ctxWithT, c, request, callback, sec)

	vs := make([]values.Value, 0)
outerLoop:
	for {
		select {
		case response, isOpen := <-callback:
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

		// Timeout when a capability panics, crashes, and does not close the channel.
		case <-ctxWithT.Done():
			return nil, fmt.Errorf("context timed out. If you did not set a timeout, be aware that the default ExecuteSync timeout is %f seconds", defaultExecuteTimeout.Seconds())
		}
	}

	setupErr := <-sec
	// Something went wrong when setting up a capability.
	if setupErr != nil {
		return nil, setupErr
	}

	// If the capability did not return any values, we deem it as an error.
	// The intent is for the API to be explicit.
	if len(vs) == 0 {
		return nil, errors.New("capability did not return any values")
	}

	// If the capability returned only one value,
	// let's unwrap it to improve usability.
	if len(vs) == 1 {
		return vs[0], nil
	}

	return &values.List{Underlying: vs}, nil
}
