package cron

import (
	"context"

	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

const flag = cre.CronCapability

type Cron struct{}

func (c *Cron) Flag() cre.CapabilityFlag {
	return flag
}

func (c *Cron) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "cron-trigger",
			Version:        "1.0.0",
			CapabilityType: 0, // TRIGGER
		},
		Config: &capabilitiespb.CapabilityConfig{
			LocalOnly: don.HasOnlyLocalCapabilities(),
		},
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

func (c *Cron) PostEnvStartup(
	_ context.Context,
	_ zerolog.Logger,
	_ *cre.Don,
	_ *cre.Dons,
	_ *cre.Environment,
) error {
	return nil
}
