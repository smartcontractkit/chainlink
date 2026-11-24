package queue

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/ocr3_1"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
)

const flag = cre.ConsensusQueueCapability

const consensusQueueLabelledName = cre.ConsensusQueueCapability

type ConsensusQueue struct{}

func (o *ConsensusQueue) Flag() cre.CapabilityFlag { return flag }
func (o *ConsensusQueue) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	if !don.HasFlag(flag) {
		return nil, nil
	}

	workers, err := don.Workers()
	if err != nil {
		return nil, fmt.Errorf("failed to get workers for DON %s: %w", don.Name, err)
	}

	// The plugin reads its OCR limits only from the plugin config inside the registry's OCR config.
	ocrConfig := contracts.DefaultOCR3_1Config(len(workers))
	ocrConfig.ConsensusQueueOffchainConfig = ocr3_1.DefaultConsensusQueueOffchainConfig()

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: []keystone_changeset.DONCapabilityWithConfig{{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   consensusQueueLabelledName,
				Version:        "1.0.0",
				CapabilityType: 2, // CONSENSUS
				ResponseType:   0, // REPORT (not used)
			},
			Config: &capabilitiespb.CapabilityConfig{
				LocalOnly: don.HasOnlyLocalCapabilities(),
			},
			UseCapRegOCRConfig: true,
		}},
		// The plugin runs on the OCR 3.1 oracle, which rejects OCR3 (version 30) offchain configs.
		CapabilityToOCR3_1Config: map[string]*ocr3_1.V3_1OracleConfig{
			consensusQueueLabelledName: ocrConfig,
		},
	}, nil
}

// PostEnvStartup does nothing. Nodes start consensus-queue from the Capabilities Registry once the DON has the
// capability, because Local CRE adds it to their RegistryBasedLaunchAllowlist. There is no job spec or contract.
func (o *ConsensusQueue) PostEnvStartup(
	_ context.Context,
	_ zerolog.Logger,
	_ *cre.Don,
	_ *cre.Dons,
	_ *cre.Environment,
) error {
	return nil
}
