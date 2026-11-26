package jobs

import (
	"context"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	job_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_pkg "github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

// CapabilitiesRegistryPayload builds the offchain capabilities registry payload (proto-JSON) for the
// local environment: DON name (the onchain name) -> capability ID -> config. Versions must strictly
// increase across payloads delivered to the same environment.
func CapabilitiesRegistryPayload(env *cldf.Environment, version uint64, dons map[string]map[string]*capabilitiespb.CapabilityConfig) (string, error) {
	return job_pkg.BuildCapabilitiesRegistryPayload(offchain.ProductLabel, env.Name, version, dons)
}

// SpecConfigCapability returns a capability config whose spec_config is specConfig, for capabilities
// whose config is a plain map (for example cron's).
func SpecConfigCapability(specConfig map[string]any) (*capabilitiespb.CapabilityConfig, error) {
	m, err := values.WrapMap(specConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap spec_config: %w", err)
	}
	return &capabilitiespb.CapabilityConfig{SpecConfig: values.ProtoMap(m)}, nil
}

// ProposeAndApproveCapabilitiesRegistry delivers payload to the nodes of target through JD as a
// capabilities_registry cresettings job and approves it on each node. cresettings specs are
// auto-approved on proposal, which Approve tolerates.
func ProposeAndApproveCapabilitiesRegistry(ctx context.Context, env *cldf.Environment, dons *cre.Dons, target *cre.Don, payload string) error {
	report, err := operations.ExecuteOperation(
		env.OperationsBundle,
		job_ops.ProposeCapabilitiesRegistryJob,
		job_ops.ProposeCapabilitiesRegistryJobDeps{Env: *env},
		job_ops.ProposeCapabilitiesRegistryJobInput{
			Domain:         offchain.ProductLabel,
			Environment:    env.Name,
			DONName:        target.Name,
			DONFilters:     []offchain.TargetDONFilter{{Key: offchain.FilterKeyDONName, Value: target.Name}},
			ExtraLabels:    map[string]string{cre.CapabilityLabelKey: "capabilities-registry"},
			OffchainConfig: payload,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to propose capabilities registry job to DON %s: %w", target.Name, err)
	}

	if err = Approve(ctx, env.Offchain, dons, report.Output.Specs); err != nil {
		return fmt.Errorf("failed to approve capabilities registry job on DON %s: %w", target.Name, err)
	}
	return nil
}
