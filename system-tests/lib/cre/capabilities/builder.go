package capabilities

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

type Capability struct {
	flag                         cre.CapabilityFlag
	jobSpecFn                    cre.JobSpecFn
	nodeConfigFn                 cre.NodeConfigFn
	gatewayJobHandlerConfigFn    cre.GatewayHandlerConfigFn
	capabilityRegistryV1ConfigFn cre.CapabilityRegistryConfigFn
	validateFn                   func(*Capability) error
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return c.flag
}

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return c.jobSpecFn
}

func (c *Capability) NodeConfigFn() cre.NodeConfigFn {
	return c.nodeConfigFn
}

func (c *Capability) GatewayJobHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return c.gatewayJobHandlerConfigFn
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return c.capabilityRegistryV1ConfigFn
}

type Option func(*Capability)

func WithJobSpecFn(jobSpecFn cre.JobSpecFn) Option {
	return func(c *Capability) {
		c.jobSpecFn = jobSpecFn
	}
}

func WithNodeConfigFn(nodeConfigFn cre.NodeConfigFn) Option {
	return func(c *Capability) {
		c.nodeConfigFn = nodeConfigFn
	}
}

func WithGatewayJobHandlerConfigFn(gatewayJobHandlerConfigFn cre.GatewayHandlerConfigFn) Option {
	return func(c *Capability) {
		c.gatewayJobHandlerConfigFn = gatewayJobHandlerConfigFn
	}
}

func WithCapabilityRegistryV1ConfigFn(capabilityRegistryV1ConfigFn cre.CapabilityRegistryConfigFn) Option {
	return func(c *Capability) {
		c.capabilityRegistryV1ConfigFn = capabilityRegistryV1ConfigFn
	}
}

func WithValidateFn(validateFn func(*Capability) error) Option {
	return func(c *Capability) {
		c.validateFn = validateFn
	}
}

func New(flag cre.CapabilityFlag, opts ...Option) (*Capability, error) {
	capability := &Capability{
		flag: flag,
	}
	for _, opt := range opts {
		opt(capability)
	}

	if capability.validateFn != nil {
		if err := capability.validateFn(capability); err != nil {
			return nil, errors.Wrapf(err, "failed to validate capability %s", capability.flag)
		}
	}

	return capability, nil
}
