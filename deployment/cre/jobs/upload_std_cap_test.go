package jobs_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestUploadStandardCapabilityJob_VerifyPreconditions(t *testing.T) {
	j := jobs.UploadStandardCapabilityJob{}
	var env cldf.Environment

	// missing job name
	err := j.VerifyPreconditions(env, jobs.UploadStandardCapabilityJobInput{Command: "run", TargetDON: &offchain.DONFilter{DONName: "d", EnvLabel: "e", ProductLabel: offchain.ProductLabel, Size: 1}})
	require.Error(t, err)
	// missing command
	err = j.VerifyPreconditions(env, jobs.UploadStandardCapabilityJobInput{JobName: "name", TargetDON: &offchain.DONFilter{DONName: "d", EnvLabel: "e", ProductLabel: offchain.ProductLabel, Size: 1}})
	require.Error(t, err)
	// missing target DON
	err = j.VerifyPreconditions(env, jobs.UploadStandardCapabilityJobInput{JobName: "name", Command: "run"})
	require.Error(t, err)
	// valid
	err = j.VerifyPreconditions(env, jobs.UploadStandardCapabilityJobInput{JobName: "name", Command: "run", TargetDON: &offchain.DONFilter{DONName: "d", EnvLabel: "e", ProductLabel: offchain.ProductLabel, Size: 1}})
	require.NoError(t, err)
}

func TestUploadStandardCapabilityJob_Apply(t *testing.T) {
	testEnv := test.SetupEnvV2(t, false)

	// Build minimal environment
	env := testEnv.Env

	input := jobs.UploadStandardCapabilityJobInput{
		JobName: "cron-cap-job",
		Command: "cron",
		TargetDON: &offchain.DONFilter{
			DONName:      test.DONName,
			EnvLabel:     "test",
			ProductLabel: offchain.ProductLabel,
			Size:         4,
		},
	}

	out, err := jobs.UploadStandardCapabilityJob{}.Apply(*env, input)
	require.NoError(t, err)
	require.Len(t, out.Reports, 1)

	//env.Offchain.ListJobs(t.Context(), )
}
