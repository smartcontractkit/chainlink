package flags

import "github.com/smartcontractkit/chainlink/system-tests/lib/cre"

type DefaultCapbilityFlagsProvider struct {
	supportedCapabilities []cre.CapabilityFlag
}

func NewDefaultCapabilityFlagsProvider() *DefaultCapbilityFlagsProvider {
	return &DefaultCapbilityFlagsProvider{
		supportedCapabilities: []cre.CapabilityFlag{
			cre.ConsensusCapability,
			cre.ConsensusCapabilityV2,
			cre.CronCapability,
			cre.EVMCapability,
			cre.CustomComputeCapability,
			cre.WriteEVMCapability,
			cre.ReadContractCapability,
			cre.LogTriggerCapability,
			cre.WebAPITargetCapability,
			cre.WebAPITriggerCapability,
			cre.MockCapability,
			cre.VaultCapability,
			cre.HTTPTriggerCapability,
			cre.HTTPActionCapability,
			cre.WriteSolanaCapability,
		},
	}
}

func (p *DefaultCapbilityFlagsProvider) SupportedCapabilityFlags() []cre.CapabilityFlag {
	return p.supportedCapabilities
}
