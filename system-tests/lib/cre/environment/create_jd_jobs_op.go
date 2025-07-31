package environment

import (
	"time"

	"github.com/Masterminds/semver/v3"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	common "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
)

type CreateJobsWithJdOpDeps struct {
	Logger                    zerolog.Logger
	SingleFileLogger          common.Logger
	HomeChainBlockchainOutput *blockchain.Output
	AddressBook               deployment.AddressBook
	JobSpecFactoryFunctions   []cre.JobSpecFactoryFn
	FullCLDEnvOutput          *cre.FullCLDEnvironmentOutput
}

type CreateJobsWithJdOpInput struct {
}

type CreateJobsWithJdOpOutput struct {
}

var CreateJobsWithJdOp = operations.NewOperation[CreateJobsWithJdOpInput, CreateJobsWithJdOpOutput, CreateJobsWithJdOpDeps](
	"create-jobs-op",
	semver.MustParse("1.0.0"),
	"Create Jobs",
	func(b operations.Bundle, deps CreateJobsWithJdOpDeps, input CreateJobsWithJdOpInput) (CreateJobsWithJdOpOutput, error) {
		createJobsStartTime := time.Now()
		deps.Logger.Info().Msg("Creating jobs with Job Distributor")

		donToJobSpecs := make(cre.DonsToJobSpecs)

		for _, jobSpecGeneratingFn := range deps.JobSpecFactoryFunctions {
			singleDonToJobSpecs, jobSpecsErr := jobSpecGeneratingFn(&cre.JobSpecFactoryInput{
				CldEnvironment:   deps.FullCLDEnvOutput.Environment,
				BlockchainOutput: deps.HomeChainBlockchainOutput,
				DonTopology:      deps.FullCLDEnvOutput.DonTopology,
				AddressBook:      deps.AddressBook,
			})
			if jobSpecsErr != nil {
				return CreateJobsWithJdOpOutput{}, pkgerrors.Wrap(jobSpecsErr, "failed to generate job specs")
			}
			mergeJobSpecSlices(singleDonToJobSpecs, donToJobSpecs)
		}

		createJobsInput := cre.CreateJobsInput{
			CldEnv:        deps.FullCLDEnvOutput.Environment,
			DonTopology:   deps.FullCLDEnvOutput.DonTopology,
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

// CreateJobsWithJdOpFactory creates a new operation with user-specified ID and version
func CreateJobsWithJdOpFactory(id string, version string) *operations.Operation[CreateJobsWithJdOpInput, CreateJobsWithJdOpOutput, CreateJobsWithJdOpDeps] {
	return operations.NewOperation[CreateJobsWithJdOpInput, CreateJobsWithJdOpOutput, CreateJobsWithJdOpDeps](
		id,
		semver.MustParse(version),
		"Create Jobs",
		func(b operations.Bundle, deps CreateJobsWithJdOpDeps, input CreateJobsWithJdOpInput) (CreateJobsWithJdOpOutput, error) {
			createJobsStartTime := time.Now()
			deps.Logger.Info().Msg("Creating jobs with Job Distributor")

			donToJobSpecs := make(cre.DonsToJobSpecs)

			for _, jobSpecGeneratingFn := range deps.JobSpecFactoryFunctions {
				singleDonToJobSpecs, jobSpecsErr := jobSpecGeneratingFn(&cre.JobSpecFactoryInput{
					CldEnvironment:   deps.FullCLDEnvOutput.Environment,
					BlockchainOutput: deps.HomeChainBlockchainOutput,
					DonTopology:      deps.FullCLDEnvOutput.DonTopology,
					AddressBook:      deps.AddressBook,
				})
				if jobSpecsErr != nil {
					return CreateJobsWithJdOpOutput{}, pkgerrors.Wrap(jobSpecsErr, "failed to generate job specs")
				}
				mergeJobSpecSlices(singleDonToJobSpecs, donToJobSpecs)
			}

			createJobsInput := cre.CreateJobsInput{
				CldEnv:        deps.FullCLDEnvOutput.Environment,
				DonTopology:   deps.FullCLDEnvOutput.DonTopology,
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
