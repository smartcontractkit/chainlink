package flags

import "github.com/smartcontractkit/chainlink/system-tests/lib/cre"

type DefaultCapbilityFlagsProvider struct {
	globalCapabilities        []cre.CapabilityFlag
	chainSpecificCapabilities []cre.CapabilityFlag
}

func NewDefaultCapabilityFlagsProvider() *DefaultCapbilityFlagsProvider {
	return &DefaultCapbilityFlagsProvider{
		globalCapabilities: []cre.CapabilityFlag{
			cre.ConsensusCapability,
			cre.ConsensusCapabilityV2,
			cre.CronCapability,
			cre.CustomComputeCapability,
			cre.WebAPITargetCapability,
			cre.WebAPITriggerCapability,
			cre.MockCapability,
			cre.VaultCapability,
			cre.HTTPTriggerCapability,
			cre.HTTPActionCapability,
			cre.WriteSolanaCapability,
		},
		chainSpecificCapabilities: []cre.CapabilityFlag{
			cre.EVMCapability,
			cre.WriteEVMCapability,
			cre.ReadContractCapability,
			cre.LogTriggerCapability,
		},
	}
}

func (p *DefaultCapbilityFlagsProvider) SupportedCapabilityFlags() []cre.CapabilityFlag {
	return append(p.globalCapabilities, p.chainSpecificCapabilities...)
}

func (p *DefaultCapbilityFlagsProvider) GlobalCapabilityFlags() []cre.CapabilityFlag {
	return p.globalCapabilities
}

func (p *DefaultCapbilityFlagsProvider) ChainSpecificCapabilityFlags() []cre.CapabilityFlag {
	return p.chainSpecificCapabilities
}

type ExtensibleCapbilityFlagsProvider struct {
	globalCapabilities        []cre.CapabilityFlag
	chainSpecificCapabilities []cre.CapabilityFlag
}

func NewExtensibleCapabilityFlagsProvider(extraGlobalFlags []string) *ExtensibleCapbilityFlagsProvider {
	return &ExtensibleCapbilityFlagsProvider{
		globalCapabilities: append([]cre.CapabilityFlag{
			cre.ConsensusCapability,
			cre.ConsensusCapabilityV2,
			cre.CronCapability,
			cre.CustomComputeCapability,
			cre.WebAPITargetCapability,
			cre.WebAPITriggerCapability,
			cre.MockCapability,
			cre.VaultCapability,
			cre.HTTPTriggerCapability,
			cre.HTTPActionCapability,
			cre.WriteSolanaCapability,
		}, extraGlobalFlags...),
		chainSpecificCapabilities: []cre.CapabilityFlag{
			cre.EVMCapability,
			cre.WriteEVMCapability,
			cre.ReadContractCapability,
			cre.LogTriggerCapability,
		},
	}
}

func (p *ExtensibleCapbilityFlagsProvider) SupportedCapabilityFlags() []cre.CapabilityFlag {
	return append(p.globalCapabilities, p.chainSpecificCapabilities...)
}

func (p *ExtensibleCapbilityFlagsProvider) GlobalCapabilityFlags() []cre.CapabilityFlag {
	return p.globalCapabilities
}

func (p *ExtensibleCapbilityFlagsProvider) ChainSpecificCapabilityFlags() []cre.CapabilityFlag {
	return p.chainSpecificCapabilities
}

// NewSwappableCapabilityFlagsProvider returns a capability flags provider that supports all capabilities that can be swapped (hot-reloaded)
// All of these capabilities are provided as external binaries
func NewSwappableCapabilityFlagsProvider() *DefaultCapbilityFlagsProvider {
	return &DefaultCapbilityFlagsProvider{
		globalCapabilities: []cre.CapabilityFlag{
			cre.ConsensusCapability,
			cre.ConsensusCapabilityV2,
			cre.CronCapability,
			cre.MockCapability,
			cre.HTTPTriggerCapability,
			cre.HTTPActionCapability,
		},
		chainSpecificCapabilities: []cre.CapabilityFlag{
			cre.EVMCapability,
			cre.ReadContractCapability,
			cre.LogTriggerCapability,
		},
	}
}
