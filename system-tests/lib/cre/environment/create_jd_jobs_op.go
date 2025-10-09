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
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type CreateJobsWithJdOpDeps struct {
	Logger                    zerolog.Logger
	SingleFileLogger          common.Logger
	HomeChainBlockchainOutput *blockchain.Output
	JobSpecFactoryFunctions   []cre.JobSpecFn
	CreEnvironment            *cre.Environment
	CapabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet
	CapabilitiesConfigs       cre.CapabilityConfigs
	Capabilities              []cre.InstallableCapability
	InfraInput                infra.Provider
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
		donToJobSpecs := make(cre.DonsToJobSpecs)

		for _, jobSpecGeneratingFn := range deps.JobSpecFactoryFunctions {
			if jobSpecGeneratingFn == nil {
				continue
			}
			singleDonToJobSpecs, jobSpecsErr := jobSpecGeneratingFn(&cre.JobSpecInput{
				CldEnvironment:            deps.CreEnvironment.CldfEnvironment,
				BlockchainOutput:          deps.HomeChainBlockchainOutput,
				DonTopology:               deps.CreEnvironment.DonTopology,
				InfraInput:                deps.InfraInput,
				CapabilityConfigs:         deps.CapabilitiesConfigs,
				CapabilitiesAwareNodeSets: deps.CapabilitiesAwareNodeSets,
				Capabilities:              deps.Capabilities,
			})
			if jobSpecsErr != nil {
				return CreateJobsWithJdOpOutput{}, pkgerrors.Wrap(jobSpecsErr, "failed to generate job specs")
			}
			mergeJobSpecSlices(singleDonToJobSpecs, donToJobSpecs)
		}

		for _, don := range deps.CreEnvironment.DonTopology.Dons.List() {
			if jobSpecs, ok := donToJobSpecs[don.ID]; ok {
				createErr := jobs.Create(b.GetContext(), deps.CreEnvironment.CldfEnvironment.Offchain, deps.CreEnvironment.DonTopology, jobSpecs)
				if createErr != nil {
					return CreateJobsWithJdOpOutput{}, pkgerrors.Wrapf(createErr, "failed to create jobs for DON %d", don.ID)
				}
			} else {
				deps.Logger.Warn().Msgf("No job specs found for DON %d", don.ID)
			}
		}

		return CreateJobsWithJdOpOutput{}, nil
	},
)

// CreateJobsWithJdOpFactory creates a new operation with user-specified ID and version
func CreateJobsWithJdOpFactory(id string, version string) *operations.Operation[CreateJobsWithJdOpInput, CreateJobsWithJdOpOutput, CreateJobsWithJdOpDeps] {
	return operations.NewOperation(
		id,
		semver.MustParse(version),
		"Create Jobs",
		func(b operations.Bundle, deps CreateJobsWithJdOpDeps, input CreateJobsWithJdOpInput) (CreateJobsWithJdOpOutput, error) {
			donToJobSpecs := make(cre.DonsToJobSpecs)

			for _, jobSpecGeneratingFn := range deps.JobSpecFactoryFunctions {
				singleDonToJobSpecs, jobSpecsErr := jobSpecGeneratingFn(&cre.JobSpecInput{
					CldEnvironment:            deps.CreEnvironment.CldfEnvironment,
					BlockchainOutput:          deps.HomeChainBlockchainOutput,
					DonTopology:               deps.CreEnvironment.DonTopology,
					CapabilitiesAwareNodeSets: deps.CapabilitiesAwareNodeSets,
					CapabilityConfigs:         deps.CapabilitiesConfigs,
					InfraInput:                deps.InfraInput,
				})
				if jobSpecsErr != nil {
					return CreateJobsWithJdOpOutput{}, pkgerrors.Wrap(jobSpecsErr, "failed to generate job specs")
				}
				mergeJobSpecSlices(singleDonToJobSpecs, donToJobSpecs)
			}

			for _, don := range deps.CreEnvironment.DonTopology.Dons.List() {
				if jobSpecs, ok := donToJobSpecs[don.ID]; ok {
					createErr := jobs.Create(b.GetContext(), deps.CreEnvironment.CldfEnvironment.Offchain, deps.CreEnvironment.DonTopology, jobSpecs)
					if createErr != nil {
						return CreateJobsWithJdOpOutput{}, pkgerrors.Wrapf(createErr, "failed to create jobs for DON %d", don.ID)
					}
				} else {
					deps.Logger.Warn().Msgf("No job specs found for DON %d", don.ID)
				}
			}

			return CreateJobsWithJdOpOutput{}, nil
		},
	)
}
