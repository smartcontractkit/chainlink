package capabilities

import (
	"context"
	"fmt"
	"regexp"

	"golang.org/x/mod/semver"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

// CapabilityType is an enum for the type of capability.
type CapabilityType string

// CapabilityType enum values.
const (
	CapabilityTypeTrigger CapabilityType = "trigger"
	CapabilityTypeAction  CapabilityType = "action"
	CapabilityTypeReport  CapabilityType = "report"
	CapabilityTypeTarget  CapabilityType = "target"
)

// IsValid checks if the capability type is valid.
func (c CapabilityType) IsValid() error {
	switch c {
	case CapabilityTypeTrigger,
		CapabilityTypeAction,
		CapabilityTypeReport,
		CapabilityTypeTarget:
		return nil
	}

	return fmt.Errorf("invalid capability type: %s", c)
}

// Validatable is an interface for validating the config and inputs of a capability.
type Validatable interface {
	ValidateConfig(config values.Map) error
	ExampleOutput() values.Value
	ValidateInput(inputs values.Map) error
}

// Executable is an interface for executing a capability.
type Executable interface {
	// Start will be called when the capability is loaded by the application.
	// Start will be called before the capability is added to the registry.
	Start(ctx context.Context, config values.Map) (values.Value, error)
	// Capability must respect context.Done and cleanup any request specific resources
	// when the context is cancelled. When a request has been completed the capability
	// is also expected to close the callback channel.
	// Request specific configuration is passed in via the inputs parameter.
	Execute(ctx context.Context, callback chan values.Value, inputs values.Map) error
	// Stop will be called before the application exits.
	// Stop will be called after the capability is removed from the registry.
	Stop(ctx context.Context) error
}

// Capability is an interface for a capability.
type Capability interface {
	Executable
	Validatable
	Info() CapabilityInfo
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

var idRegex = regexp.MustCompile("[a-z0-9_\\-:]")

// NewCapabilityInfo returns a new CapabilityInfo.
func NewCapabilityInfo(
	id string,
	capabilityType CapabilityType,
	description string,
	version string,
) (CapabilityInfo, error) {
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

// ExecuteSync executes a capability synchronously.
// We are not handling a case where a capability panics and crashes.
func ExecuteSync(ctx context.Context, c Capability, inputs values.Map) (values.Value, error) {
	callback := make(chan values.Value)
	vs := make([]values.Value, 0)

	var executionErr error
	go func(innerCtx context.Context, innerC Capability, innerInputs values.Map) {
		executionErr = innerC.Execute(ctx, callback, inputs)
	}(ctx, c, inputs)

	for value := range callback {
		// TODO: Handle the case where a capability returns an error as part of the callback.
		// if valErr, ok := value.(values.Error); ok {
		// 	return nil, valError.Underlying
		// }
		vs = append(vs, value)
	}

	// Something went wrong when executing a capability. If this happens at any point,
	// we want to stop the capability and return the error. We are discarding all values
	// returned up to the error.
	if executionErr != nil {
		return nil, executionErr
	}

	// TODO: This is a special case for when a capability returns no values.
	// It can happen when the channel is closed before a value is returned.
	// if len(vs) == 0 {
	// 	return nil, nil
	// }

	if len(vs) == 1 {
		return vs[0], nil
	}

	return &values.List{Underlying: vs}, nil
}
