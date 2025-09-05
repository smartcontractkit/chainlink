package jobs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestProposeJobSpec_VerifyPreconditions(t *testing.T) {
	j := jobs.ProposeJobSpec{}
	var env cldf.Environment

	testCases := []struct {
		name        string
		input       jobs.ProposeJobSpecInput
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid cron job",
			input: jobs.ProposeJobSpecInput{
				Environment: "test",
				JobName:     "cron-test",
				Domain:      "cre",
				DONName:     "test-don",
				DONFilters: []operations.TargetDONFilter{
					{Key: operations.FilterKeyDONName, Value: "d"},
					{Key: "environment", Value: "e"},
					{Key: "product", Value: offchain.ProductLabel},
				},
				Template: job_types.Cron,
				Inputs:   job_types.JobSpecInput{},
			},
			expectError: false,
		},
		{
			name: "missing environment",
			input: jobs.ProposeJobSpecInput{
				Domain:   "cre",
				Template: job_types.Cron,
				Inputs:   job_types.JobSpecInput{},
			},
			expectError: true,
			errorMsg:    "environment is required",
		},
		{
			name: "missing domain",
			input: jobs.ProposeJobSpecInput{
				Environment: "test",
				Template:    job_types.Cron,
				Inputs:      job_types.JobSpecInput{},
			},
			expectError: true,
			errorMsg:    "domain is required",
		},
		{
			name: "missing don name",
			input: jobs.ProposeJobSpecInput{
				Environment: "test",
				Domain:      "cre",
				Template:    1,
				Inputs:      job_types.JobSpecInput{},
			},
			expectError: true,
			errorMsg:    "don_name is required",
		},
		{
			name: "missing don filters",
			input: jobs.ProposeJobSpecInput{
				Environment: "test",
				Domain:      "cre",
				DONName:     "test-don",
				Template:    1,
				Inputs:      job_types.JobSpecInput{},
			},
			expectError: true,
			errorMsg:    "don_filters is required",
		},
		{
			name: "missing job name",
			input: jobs.ProposeJobSpecInput{
				Environment: "test",
				Domain:      "cre",
				DONName:     "test-don",
				DONFilters: []operations.TargetDONFilter{
					{Key: operations.FilterKeyDONName, Value: "d"},
					{Key: "environment", Value: "e"},
					{Key: "product", Value: offchain.ProductLabel},
				},
				Template: 1,
				Inputs:   job_types.JobSpecInput{},
			},
			expectError: true,
			errorMsg:    "job_name is required",
		},
		{
			name: "unsupported template",
			input: jobs.ProposeJobSpecInput{
				Environment: "test",
				Domain:      "cre",
				DONName:     "test-don",
				JobName:     "cron-test",
				DONFilters: []operations.TargetDONFilter{
					{Key: operations.FilterKeyDONName, Value: "d"},
					{Key: "environment", Value: "e"},
					{Key: "product", Value: offchain.ProductLabel},
				},
				Template: 1,
				Inputs:   job_types.JobSpecInput{},
			},
			expectError: true,
			errorMsg:    "unsupported template",
		},
		{
			name: "missing inputs",
			input: jobs.ProposeJobSpecInput{
				Environment: "test",
				Domain:      "cre",
				DONName:     "test-don",
				JobName:     "cron-test",
				DONFilters: []operations.TargetDONFilter{
					{Key: operations.FilterKeyDONName, Value: "d"},
					{Key: "environment", Value: "e"},
					{Key: "product", Value: offchain.ProductLabel},
				},
				Template: job_types.Cron,
				Inputs:   nil,
			},
			expectError: true,
			errorMsg:    "inputs are required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := j.VerifyPreconditions(env, tc.input)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProposeJobSpec_Apply(t *testing.T) {
	testEnv := test.SetupEnvV2(t, false)
	env := testEnv.Env

	t.Run("successful cron job distribution", func(t *testing.T) {
		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "cron-cap-job",
			DONName:     test.DONName,
			Template:    job_types.Cron,
			DONFilters: []operations.TargetDONFilter{
				{Key: operations.FilterKeyDONName, Value: "don-" + test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"command":       "cron",
				"config":        "CRON_TZ=UTC * * * * *",
				"externalJobID": "a-cron-job-id",
				"oracleFactory": pkg.OracleFactory{
					Enabled: false,
				},
			},
		}

		out, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.NoError(t, err)
		assert.Len(t, out.Reports, 1)

		reqs, err := testEnv.TestJD.ListProposedJobRequests()
		require.NoError(t, err)
		assert.Len(t, reqs, 4)

		for _, req := range reqs {
			// log each spec in readable yaml format
			t.Logf("Job Spec:\n%s", req.Spec)
			assert.Contains(t, req.Spec, `name = "cron-cap-job"`)
			assert.Contains(t, req.Spec, `command = "cron"`)
			assert.Contains(t, req.Spec, `config = """CRON_TZ=UTC * * * * *"""`)
			assert.Contains(t, req.Spec, `externalJobID = "a-cron-job-id"`)
		}
	})

	t.Run("failed cron job distribution due to bad input", func(t *testing.T) {
		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "cron-cap-job",
			Template:    job_types.Cron,
			DONFilters: []operations.TargetDONFilter{
				{Key: operations.FilterKeyDONName, Value: "don" + test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				// Missing "command"
				"config":        "CRON_TZ=UTC * * * * *",
				"externalJobID": "a-cron-job-id",
				"oracleFactory": pkg.OracleFactory{
					Enabled: false,
				},
			},
		}

		_, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to convert inputs to standard capability job")
		assert.Contains(t, err.Error(), "command is required and must be a string")
	})
}
