package environment

import (
	"time"

	"github.com/Masterminds/semver/v3"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	common "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
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

		createJobsInput := cre.CreateJobsInput{
			CldEnv:        deps.CreEnvironment.CldfEnvironment,
			DonTopology:   deps.CreEnvironment.DonTopology,
			DonToJobSpecs: donToJobSpecs,
		}

		jobsErr := libdon.CreateJobs(b.GetContext(), deps.Logger, createJobsInput)
		if jobsErr != nil {
			return CreateJobsWithJdOpOutput{}, pkgerrors.Wrap(jobsErr, "failed to create jobs")
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
			createJobsStartTime := time.Now()
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

			createJobsInput := cre.CreateJobsInput{
				CldEnv:        deps.CreEnvironment.CldfEnvironment,
				DonTopology:   deps.CreEnvironment.DonTopology,
				DonToJobSpecs: donToJobSpecs,
			}

			jobsErr := libdon.CreateJobs(b.GetContext(), deps.Logger, createJobsInput)
			if jobsErr != nil {
				return CreateJobsWithJdOpOutput{}, pkgerrors.Wrap(jobsErr, "failed to create jobs")
			}

			deps.Logger.Info().Msgf("Jobs created in %.2f seconds", time.Since(createJobsStartTime).Seconds())

			return CreateJobsWithJdOpOutput{}, nil
		},
	)
}
