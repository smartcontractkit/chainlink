package dontime

import (
	"context"

	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
)

const flag = cre.DONTimeCapability

const donTimeLabelledName = "dontime"

type DONTime struct{}

func (o *DONTime) Flag() cre.CapabilityFlag {
	return flag
}

func (o *DONTime) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName: donTimeLabelledName,
			Version:      "1.0.0",
		},
		Config: &capabilitiespb.CapabilityConfig{
			LocalOnly: don.HasOnlyLocalCapabilities(),
		},
		UseCapRegOCRConfig: true,
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
		CapabilityToOCR3Config: map[string]*ocr3.OracleConfig{
			donTimeLabelledName: contracts.DefaultOCR3Config(),
		},
	}, nil
}

func (o *DONTime) PostEnvStartup(
	_ context.Context,
	_ zerolog.Logger,
	_ *cre.Don,
	_ *cre.Dons,
	_ *cre.Environment,
) error {
	return nil
}
