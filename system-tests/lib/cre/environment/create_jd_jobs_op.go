package environment

import (
	"github.com/Masterminds/semver/v3"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	common "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
)

type CreateJobsWithJdOpDeps struct {
	Logger                    zerolog.Logger
	SingleFileLogger          common.Logger
	HomeChainBlockchainOutput *blockchain.Output
	JobSpecFactoryFunctions   []cre.JobSpecFn
	CreEnvironment            *cre.Environment
	Dons                      *cre.Dons
	CapabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet
	Capabilities              []cre.InstallableCapability
}

type CreateJobsWithJdOpInput struct {
}

type CreateJobsWithJdOpOutput struct {
}

var CreateJobsWithJdOp = operations.NewOperation(
	"create-jobs-op",
	semver.MustParse("1.0.0"),
	"Create Jobs",
	func(b operations.Bundle, deps CreateJobsWithJdOpDeps, input CreateJobsWithJdOpInput) (CreateJobsWithJdOpOutput, error) {
		for _, jobSpecGeneratingFn := range deps.JobSpecFactoryFunctions {
			if jobSpecGeneratingFn == nil {
				continue
			}

			for idx, don := range deps.Dons.List() {
				jobSpecs, jobSpecsErr := jobSpecGeneratingFn(&cre.JobSpecInput{
					CreEnvironment: deps.CreEnvironment,
					Don:            don,
					NodeSet:        cre.ConvertToNodeSetWithChainCapabilities(deps.CapabilitiesAwareNodeSets)[idx],
					Capabilities:   deps.Capabilities, // TODO how was this used?
				})
				if jobSpecsErr != nil {
					return CreateJobsWithJdOpOutput{}, pkgerrors.Wrap(jobSpecsErr, "failed to generate job specs")
				}

				createErr := jobs.Create(b.GetContext(), deps.CreEnvironment.CldfEnvironment.Offchain, deps.Dons, jobSpecs)
				if createErr != nil {
					return CreateJobsWithJdOpOutput{}, pkgerrors.Wrapf(createErr, "failed to create jobs for DON %d", don.ID)
				}
			}
		}

		return CreateJobsWithJdOpOutput{}, nil
	},
)
