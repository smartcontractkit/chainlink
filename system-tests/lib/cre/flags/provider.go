package flags

import "github.com/smartcontractkit/chainlink/system-tests/lib/cre"

type DefaultCapbilityFlagsProvider struct {
	capabilities []cre.CapabilityFlag
}

func NewDefaultCapabilityFlagsProvider() *DefaultCapbilityFlagsProvider {
	return &DefaultCapbilityFlagsProvider{
		capabilities: []cre.CapabilityFlag{
			cre.ConsensusCapability,
			cre.ConsensusCapabilityV2,
			cre.CronCapability,
			cre.CustomComputeCapability,
			cre.DONTimeCapability,
			cre.WebAPITargetCapability,
			cre.WebAPITriggerCapability,
			cre.MockCapability,
			cre.VaultCapability,
			cre.HTTPTriggerCapability,
			cre.HTTPActionCapability,
			cre.SolanaCapability,
			cre.EVMCapability,
			cre.WriteEVMCapability,
			cre.ReadContractCapability,
			cre.LogEventTriggerCapability,
			cre.AptosCapability,
		},
	}
}

func (p *DefaultCapbilityFlagsProvider) SupportedCapabilityFlags() []cre.CapabilityFlag {
	return p.capabilities
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
		}, extraGlobalFlags...),
		chainSpecificCapabilities: []cre.CapabilityFlag{
			cre.EVMCapability,
			cre.WriteEVMCapability,
			cre.SolanaCapability,
			cre.ReadContractCapability,
			cre.LogEventTriggerCapability,
			cre.AptosCapability,
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
		capabilities: []cre.CapabilityFlag{
			cre.ConsensusCapability,
			cre.ConsensusCapabilityV2,
			cre.CronCapability,
			cre.MockCapability,
			cre.HTTPTriggerCapability,
			cre.HTTPActionCapability,
			cre.EVMCapability,
			cre.ReadContractCapability,
			cre.LogEventTriggerCapability,
			cre.SolanaCapability,
			cre.AptosCapability,
		},
	}
}
