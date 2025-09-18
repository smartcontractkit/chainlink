package jobs_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
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
				DONFilters: []offchain.TargetDONFilter{
					{Key: offchain.FilterKeyDONName, Value: "d"},
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
				Template:    job_types.Cron,
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
				Template:    job_types.Cron,
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
				DONFilters: []offchain.TargetDONFilter{
					{Key: offchain.FilterKeyDONName, Value: "d"},
					{Key: "environment", Value: "e"},
					{Key: "product", Value: offchain.ProductLabel},
				},
				Template: job_types.Cron,
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
				DONFilters: []offchain.TargetDONFilter{
					{Key: offchain.FilterKeyDONName, Value: "d"},
					{Key: "environment", Value: "e"},
					{Key: "product", Value: offchain.ProductLabel},
				},
				Template: 100,
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
				DONFilters: []offchain.TargetDONFilter{
					{Key: offchain.FilterKeyDONName, Value: "d"},
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
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: "don-" + test.DONName},
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
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: "don" + test.DONName},
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

	t.Run("successful ocr3 bootstrap job distribution", func(t *testing.T) {
		chainSelector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
		ds := datastore.NewMemoryDataStore()

		err := ds.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSelector,
			Type:          datastore.ContractType(ocr3.OCR3Capability),
			Version:       semver.MustParse("1.0.0"),
			Address:       "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B",
			Qualifier:     "ocr3-contract-qualifier",
		})
		require.NoError(t, err)

		env.DataStore = ds.Seal()

		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "ocr3-bootstrap-job",
			DONName:     test.DONName,
			Template:    job_types.BootstrapOCR3,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: "don-" + test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"contract_qualifier": "ocr3-contract-qualifier",
				"chain_selector":     strconv.FormatUint(chainSelector, 10),
			},
		}

		out, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.NoError(t, err)
		assert.Len(t, out.Reports, 1)

		reqs, err := testEnv.TestJD.ListProposedJobRequests()
		require.NoError(t, err)
		assert.Len(t, reqs, 5)

		expectedChainID := chainsel.ETHEREUM_TESTNET_SEPOLIA.EvmChainID

		for _, req := range reqs {
			if !strings.Contains(req.Spec, `type = "bootstrap"`) {
				continue
			}
			// log each spec in readable yaml format
			t.Logf("Job Spec:\n%s", req.Spec)
			assert.Contains(t, req.Spec, `name = "ocr3-bootstrap-job`)
			assert.Contains(t, req.Spec, `contractID = "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"`)
			assert.Contains(t, req.Spec, fmt.Sprintf("chainID = %d", expectedChainID))
		}
	})

	t.Run("failed ocr3 bootstrap job distribution", func(t *testing.T) {
		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "ocr3-bootstrap-job",
			DONName:     test.DONName,
			Template:    job_types.BootstrapOCR3,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: "don-" + test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				// Missing "chain_selector"
				"contract_qualifier": "ocr-contract-qualifier",
			},
		}

		_, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to convert inputs to OCR3 bootstrap job input")
		assert.Contains(t, err.Error(), "chain_selector is required and must be a string")
	})

	t.Run("successful ocr3 job distribution", func(t *testing.T) {
		chainSelector := testEnv.RegistrySelector
		ds := datastore.NewMemoryDataStore()

		err := ds.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSelector,
			Type:          datastore.ContractType(ocr3.OCR3Capability),
			Version:       semver.MustParse("1.0.0"),
			Address:       "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B",
			Qualifier:     "ocr3-contract-qualifier",
		})
		require.NoError(t, err)

		env.DataStore = ds.Seal()

		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "ocr3-job",
			DONName:     test.DONName,
			Template:    job_types.OCR3,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: "don-" + test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"template_name":        "worker-ocr3",
				"contract_qualifier":   "ocr3-contract-qualifier",
				"chain_selector_evm":   strconv.FormatUint(chainSelector, 10),
				"chain_selector_aptos": strconv.FormatUint(testEnv.AptosSelector, 10),
				"bootstrapper_ocr3_urls": []string{
					"12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001",
				},
			},
		}

		out, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.NoError(t, err)
		assert.Len(t, out.Reports, 1)

		reqs, err := testEnv.TestJD.ListProposedJobRequests()
		require.NoError(t, err)
		assert.Len(t, reqs, 9)

		expectedChainID := chainsel.TEST_90000001.EvmChainID

		for _, req := range reqs {
			if !strings.Contains(req.Spec, `type = "offchainreporting2"`) {
				continue
			}
			// log each spec in readable yaml format
			t.Logf("Job Spec:\n%s", req.Spec)
			assert.Contains(t, req.Spec, `name = "ocr3-job`)
			assert.Contains(t, req.Spec, `contractID = "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"`)
			assert.Contains(t, req.Spec, `p2pv2Bootstrappers = [
  "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001",
]`)
			assert.Contains(t, req.Spec, fmt.Sprintf(`chainID = "%d"`, expectedChainID))
			assert.Contains(t, req.Spec, `command = "/usr/local/bin/chainlink-ocr3-capability"`)
			assert.Contains(t, req.Spec, `pluginName = "ocr-capability"`)
			assert.Contains(t, req.Spec, `providerType = "ocr3-capability"`)
			assert.Contains(t, req.Spec, `strategyName = 'multi-chain'`)
		}
	})

	t.Run("failed ocr3 job distribution", func(t *testing.T) {
		chainSelector := testEnv.RegistrySelector
		ds := datastore.NewMemoryDataStore()

		err := ds.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSelector,
			Type:          datastore.ContractType(ocr3.OCR3Capability),
			Version:       semver.MustParse("1.0.0"),
			Address:       "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B",
			Qualifier:     "ocr3-contract-qualifier",
		})
		require.NoError(t, err)

		env.DataStore = ds.Seal()

		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "ocr3-job",
			DONName:     test.DONName,
			Template:    job_types.OCR3,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: "don-" + test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				// missing `template_name`
				"contract_qualifier": "ocr3-contract-qualifier",
			},
		}

		_, err = jobs.ProposeJobSpec{}.Apply(*env, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to convert inputs to OCR3 job input")
		assert.Contains(t, err.Error(), "template_name is required and must be a non-empty string")
	})
}
