package jobs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestProposeStandardCapabilityJob_VerifyPreconditions(t *testing.T) {
	j := jobs.ProposeStandardCapabilityJob{}
	var env cldf.Environment

	// missing job name
	err := j.VerifyPreconditions(env, jobs.ProposeStandardCapabilityJobInput{
		Command: "run",
	})
	require.Error(t, err)
	// missing command
	err = j.VerifyPreconditions(env, jobs.ProposeStandardCapabilityJobInput{
		JobName: "name",
	})
	require.Error(t, err)
	// missing DON name
	err = j.VerifyPreconditions(env, jobs.ProposeStandardCapabilityJobInput{
		JobName: "name",
		Command: "run",
	})
	require.Error(t, err)
	// missing DON Filters
	err = j.VerifyPreconditions(env, jobs.ProposeStandardCapabilityJobInput{JobName: "name", Command: "run", DONName: "test-don"})
	require.Error(t, err)
	// valid
	err = j.VerifyPreconditions(env, jobs.ProposeStandardCapabilityJobInput{
		JobName: "name",
		Command: "run",
		DONName: "test-don",
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: "d"},
			{Key: "environment", Value: "e"},
			{Key: "product", Value: offchain.ProductLabel},
		},
	})
	require.NoError(t, err)
}

func TestProposeStandardCapabilityJob_Apply(t *testing.T) {
	testEnv := test.SetupEnvV2(t, false)

	// Build minimal environment
	env := testEnv.Env

	input := jobs.ProposeStandardCapabilityJobInput{
		JobName: "cron-cap-job",
		Command: "cron",
		DONName: "test-don",
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: "don-" + test.DONName},
			{Key: "environment", Value: "test"},
			{Key: "product", Value: offchain.ProductLabel},
		},
	}

	out, err := jobs.ProposeStandardCapabilityJob{}.Apply(*env, input)
	require.NoError(t, err)
	assert.Len(t, out.Reports, 1)

	reqs, err := testEnv.TestJD.ListProposedJobRequests()
	require.NoError(t, err)
	assert.Len(t, reqs, 4)
}
