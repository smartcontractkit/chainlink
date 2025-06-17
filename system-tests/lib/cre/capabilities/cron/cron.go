package cron

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

var CronCapabilityFactoryFn = func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
	var capabilities []keystone_changeset.DONCapabilityWithConfig

	if flags.HasFlag(donFlags, types.CronCapability) {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "cron-trigger",
				Version:        "1.0.0",
				CapabilityType: 0, // TRIGGER
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities
}
