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
