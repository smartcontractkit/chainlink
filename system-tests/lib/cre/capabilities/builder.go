package capabilities

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

type registryConfigFns struct {
	V1 cre.CapabilityRegistryConfigFn
	V2 cre.CapabilityRegistryConfigFn
}

type Capability struct {
	flag                      cre.CapabilityFlag
	jobSpecFn                 cre.JobSpecFn
	nodeConfigFn              cre.NodeConfigTransformerFn
	gatewayJobHandlerConfigFn cre.GatewayHandlerConfigFn
	registryConfigFns         *registryConfigFns
	validateFn                func(*Capability) error
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return c.flag
}

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return c.jobSpecFn
}

func (c *Capability) NodeConfigTransformerFn() cre.NodeConfigTransformerFn {
	return c.nodeConfigFn
}

func (c *Capability) GatewayJobHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return c.gatewayJobHandlerConfigFn
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	if c.registryConfigFns != nil {
		return c.registryConfigFns.V1
	}
	return nil
}

func (c *Capability) CapabilityRegistryV2ConfigFn() cre.CapabilityRegistryConfigFn {
	if c.registryConfigFns != nil {
		return c.registryConfigFns.V2
	}
	return nil
}

type Option func(*Capability)

func WithJobSpecFn(jobSpecFn cre.JobSpecFn) Option {
	return func(c *Capability) {
		c.jobSpecFn = jobSpecFn
	}
}

func WithNodeConfigTransformerFn(nodeConfigFn cre.NodeConfigTransformerFn) Option {
	return func(c *Capability) {
		c.nodeConfigFn = nodeConfigFn
	}
}

func WithGatewayJobHandlerConfigFn(gatewayJobHandlerConfigFn cre.GatewayHandlerConfigFn) Option {
	return func(c *Capability) {
		c.gatewayJobHandlerConfigFn = gatewayJobHandlerConfigFn
	}
}

func WithCapabilityRegistryV1ConfigFn(fn cre.CapabilityRegistryConfigFn) Option {
	return func(c *Capability) {
		if c.registryConfigFns == nil {
			c.registryConfigFns = &registryConfigFns{
				V1: fn,
			}
			return
		}
		c.registryConfigFns.V1 = fn
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
