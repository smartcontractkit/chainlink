package jobs_test

import (
	"strings"
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
		Domain:  offchain.ProductLabel,
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: test.DONName},
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

func TestProposeStandardCapabilityJob_Apply_HTTPTrigger(t *testing.T) {
	testEnv := test.SetupEnvV2(t, false)
	env := testEnv.Env

	input := jobs.ProposeStandardCapabilityJobInput{
		JobName:       "http-trigger-job",
		Command:       "http_trigger",
		Config:        `{}`,
		ExternalJobID: "http-trigger-external-id",
		DONName:       "test-don",
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

	// Verify the job specs contain expected HTTP trigger configuration
	for _, req := range reqs {
		if !strings.Contains(req.Spec, `name = "http-trigger-job"`) {
			continue
		}
		t.Logf("HTTP Trigger Job Spec:\n%s", req.Spec)
		assert.Contains(t, req.Spec, `command = "http_trigger"`)
		assert.Contains(t, req.Spec, `config = """{}"""`)
		assert.Contains(t, req.Spec, `externalJobID = "http-trigger-external-id"`)
	}
}

func TestProposeStandardCapabilityJob_Apply_HTTPAction(t *testing.T) {
	testEnv := test.SetupEnvV2(t, false)
	env := testEnv.Env

	input := jobs.ProposeStandardCapabilityJobInput{
		JobName:       "http-action-job",
		Command:       "http_action",
		Config:        `{"proxyMode": "direct"}`,
		ExternalJobID: "http-action-external-id",
		DONName:       "test-don",
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

	// Verify the job specs contain expected HTTP action configuration
	for _, req := range reqs {
		if !strings.Contains(req.Spec, `name = "http-action-job"`) {
			continue
		}
		t.Logf("HTTP Action Job Spec:\n%s", req.Spec)
		assert.Contains(t, req.Spec, `command = "http_action"`)
		assert.Contains(t, req.Spec, `config = """{"proxyMode": "direct"}"""`)
		assert.Contains(t, req.Spec, `externalJobID = "http-action-external-id"`)
	}
}

func TestProposeStandardCapabilityJob_VerifyPreconditions_HTTPJobs(t *testing.T) {
	j := jobs.ProposeStandardCapabilityJob{}
	var env cldf.Environment

	testCases := []struct {
		name        string
		input       jobs.ProposeStandardCapabilityJobInput
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid http trigger job",
			input: jobs.ProposeStandardCapabilityJobInput{
				JobName: "http-trigger-test",
				Command: "http_trigger",
				Config:  `{}`,
				DONName: "test-don",
				DONFilters: []offchain.TargetDONFilter{
					{Key: offchain.FilterKeyDONName, Value: "d"},
					{Key: "environment", Value: "e"},
					{Key: "product", Value: offchain.ProductLabel},
				},
			},
			expectError: false,
		},
		{
			name: "valid http action job",
			input: jobs.ProposeStandardCapabilityJobInput{
				JobName: "http-action-test",
				Command: "http_action",
				Config:  `{"proxyMode": "direct"}`,
				DONName: "test-don",
				DONFilters: []offchain.TargetDONFilter{
					{Key: offchain.FilterKeyDONName, Value: "d"},
					{Key: "environment", Value: "e"},
					{Key: "product", Value: offchain.ProductLabel},
				},
			},
			expectError: false,
		},
		{
			name: "http trigger with empty config",
			input: jobs.ProposeStandardCapabilityJobInput{
				JobName: "http-trigger-test",
				Command: "http_trigger",
				Config:  "",
				DONName: "test-don",
				DONFilters: []offchain.TargetDONFilter{
					{Key: offchain.FilterKeyDONName, Value: "d"},
				},
			},
			expectError: false, // Empty config is allowed
		},
		{
			name: "http action missing command",
			input: jobs.ProposeStandardCapabilityJobInput{
				JobName: "http-action-test",
				Config:  `{"proxyMode": "direct"}`,
				DONName: "test-don",
				DONFilters: []offchain.TargetDONFilter{
					{Key: offchain.FilterKeyDONName, Value: "d"},
				},
			},
			expectError: true,
			errorMsg:    "command is required",
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
