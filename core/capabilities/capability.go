package capabilities

import (
	"context"
	"fmt"
	"regexp"

	"golang.org/x/mod/semver"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type CapabilityType string

const (
	CapabilityTypeTrigger CapabilityType = "trigger"
	CapabilityTypeAction  CapabilityType = "action"
	CapabilityTypeReport  CapabilityType = "report"
	CapabilityTypeTarget  CapabilityType = "target"
)

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

type Validatable interface {
	ValidateConfig(config values.Map) error
	ExampleOutput() values.Value
	ValidateInput(inputs values.Map) error
}

type Capability interface {
	Validatable
	Info() CapabilityInfo
}

type SynchronousCapability interface {
	Capability

	Start(ctx context.Context, config values.Map) (values.Value, error)
	Execute(ctx context.Context, inputs values.Map) (values.Value, error)
	Stop(ctx context.Context) error
}

type AsynchronousCapability interface {
	Capability

	Start(ctx context.Context, config values.Map) (values.Value, error)
	Execute(ctx context.Context, callback chan values.Map, inputs values.Map) (values.Value, error)
	Stop(ctx context.Context) error
}

type stringer struct {
	s string
}

func (s stringer) String() string {
	return s.s
}

func Stringer(s string) fmt.Stringer {
	return stringer{s: s}
}

type CapabilityInfo struct {
	// We use `fmt.Stringer` for the ID, since an ID can take
	// one of two forms (namely a fully-qualified ID expressed as a
	// string), or a tags object.
	Id             fmt.Stringer
	CapabilityType CapabilityType
	Description    string
	Version        string
}

func (c CapabilityInfo) Info() CapabilityInfo {
	return c
}

var idRegex = regexp.MustCompile("[a-z0-9_\\-:]")

func NewCapabilityInfo(
	id fmt.Stringer,
	capabilityType CapabilityType,
	description string,
	version string,
) (CapabilityInfo, error) {
	if !idRegex.MatchString(id.String()) {
		return CapabilityInfo{}, fmt.Errorf("invalid id: %s. Allowed: %s", id, idRegex)
	}

	if ok := semver.IsValid(version); !ok {
		return CapabilityInfo{}, fmt.Errorf("invalid version: %+v", version)
	}

	if err := capabilityType.IsValid(); err != nil {
		return CapabilityInfo{}, err
	}

	return CapabilityInfo{
		Id:             id,
		CapabilityType: capabilityType,
		Description:    description,
		Version:        version,
	}, nil
}
