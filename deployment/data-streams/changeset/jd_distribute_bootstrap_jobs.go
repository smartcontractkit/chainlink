package changeset

import (
	"context"
	"fmt"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jobs"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
)

var _ deployment.ChangeSetV2[CsDistributeBootstrapJobSpecsConfig] = CsDistributeBootstrapJobSpecs{}

type CsDistributeBootstrapJobSpecsConfig struct {
	ChainSelectorEVM    uint64
	Filter              *jd.ListFilter
	ConfiguratorAddress string
}

type CsDistributeBootstrapJobSpecs struct{}

func (CsDistributeBootstrapJobSpecs) Apply(e deployment.Environment, cfg CsDistributeBootstrapJobSpecsConfig) (deployment.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(e.GetContext(), defaultJobSpecsTimeout)
	defer cancel()

	chainID, _, err := chainAndAddresses(e, cfg.ChainSelectorEVM)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Add a label to the job spec to identify the related DON
	labels := append([]*ptypes.Label(nil),
		&ptypes.Label{
			Key: utils.DonIdentifier(cfg.Filter.DONID, cfg.Filter.DONName),
		})

	bootstrapSpec := jobs.NewBootstrapSpec(
		cfg.ConfiguratorAddress,
		cfg.Filter.DONID,
		jobs.RelayTypeEVM,
		jobs.RelayConfig{
			ChainID: chainID,
		},
	)

	renderedSpec, err := bootstrapSpec.MarshalTOML()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to marshal bootstrap spec: %w", err)
	}

	bootstrapNodes, err := jd.FetchDONBootstrappersFromJD(ctx, e.Offchain, cfg.Filter)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get workflow don nodes: %w", err)
	}

	var proposals []*jobv1.ProposeJobRequest
	for _, node := range bootstrapNodes {
		proposals = append(proposals, &jobv1.ProposeJobRequest{
			NodeId: node.Id,
			Spec:   string(renderedSpec),
			Labels: labels,
		})
	}
	proposedJobs, err := proposeAllOrNothing(ctx, e.Offchain, proposals)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to propose all bootstrap jobs: %w", err)
	}

	return deployment.ChangesetOutput{
		Jobs: proposedJobs,
	}, nil
}

func (f CsDistributeBootstrapJobSpecs) VerifyPreconditions(e deployment.Environment, config CsDistributeBootstrapJobSpecsConfig) error {
	return nil
}
