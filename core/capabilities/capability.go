package capabilities

import (
	"fmt"
	"regexp"
	"golang.org/x/mod/semver"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type CapabilityType string

const (
	CapabilityTypeTrigger CapabilityType = "trigger"
	CapabilityTypeAction CapabilityType = "action"
	CapabilityTypeReport CapabilityType = "report"
	CapabilityTypeTarget CapabilityType = "target"
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

type CapabilityInfo struct {
	// We use `fmt.Stringer` for the ID, since an ID can take
	// one of two forms (namely a fully-qualified ID expressed as a
	// string), or a tags object.
	id fmt.Stringer
	capabilityType CapabilityType
	description string
	version string
}

type Capability interface {
	Validatable
	Info() (CapabilityInfo)
}

type CapabilityRegistry interface {
	ListCapabilities() ([]CapabilityInfo)
	Get(id string) (Capability, error)
	Add(capability Capability) error
}

type CapabilityInfoProvider struct {
	info CapabilityInfo
}

func (c *CapabilityInfoProvider) Info() CapabilityInfo {
	return c.info
}

var idRegex = regexp.MustCompile("[a-z0-9_-:]")

func NewCapabilityInfoProvider(
	id fmt.Stringer,
	capabilityType CapabilityType,
	description string,
	version string,
) (*CapabilityInfoProvider, error) {
	if !idRegex.MatchString(id.String()) {
		return nil, fmt.Errorf("invalid id: %s. Allowed: %s", id, idRegex)
	}

	if ok := semver.IsValid(version); !ok {
		return nil, fmt.Errorf("invalid version: %+v", version)
	}
	
	if err := capabilityType.IsValid(); err != nil {
		return nil, err
	}

	return &CapabilityInfoProvider{
		info: CapabilityInfo{
			id: id,
			capabilityType: capabilityType,
			description: description,
			version: version,
		},
	}, nil
}
