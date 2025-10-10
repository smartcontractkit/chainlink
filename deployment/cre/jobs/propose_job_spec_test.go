package jobs_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/quarantine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/sequences"
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
			name: "valid http trigger job",
			input: jobs.ProposeJobSpecInput{
				Environment: "test",
				JobName:     "http-trigger-test",
				Domain:      "cre",
				DONName:     "test-don",
				DONFilters: []offchain.TargetDONFilter{
					{Key: offchain.FilterKeyDONName, Value: "d"},
					{Key: "environment", Value: "e"},
					{Key: "product", Value: offchain.ProductLabel},
				},
				Template: job_types.HTTPTrigger,
				Inputs: job_types.JobSpecInput{
					"command":       "http_trigger",
					"config":        `{}`,
					"externalJobID": "http-trigger-job-id",
				},
			},
			expectError: false,
		},
		{
			name: "valid http action job",
			input: jobs.ProposeJobSpecInput{
				Environment: "test",
				JobName:     "http-action-test",
				Domain:      "cre",
				DONName:     "test-don",
				DONFilters: []offchain.TargetDONFilter{
					{Key: offchain.FilterKeyDONName, Value: "d"},
					{Key: "environment", Value: "e"},
					{Key: "product", Value: offchain.ProductLabel},
				},
				Template: job_types.HTTPAction,
				Inputs: job_types.JobSpecInput{
					"command":       "http_action",
					"config":        `{"proxyMode": "direct"}`,
					"externalJobID": "http-action-job-id",
				},
			},
			expectError: false,
		},
		{
			name: "valid evm job",
			input: jobs.ProposeJobSpecInput{
				Environment: "test",
				JobName:     "evm-test",
				Domain:      "cre",
				DONName:     "test-don",
				DONFilters: []offchain.TargetDONFilter{
					{Key: offchain.FilterKeyDONName, Value: "d"},
					{Key: "environment", Value: "e"},
					{Key: "product", Value: offchain.ProductLabel},
				},
				Template: job_types.EVM,
				Inputs: job_types.JobSpecInput{
					"command": "/usr/local/bin/evm",
					"config":  `{"chainId":1337,"network":"evm"}`,
					"oracleFactory": pkg.OracleFactory{
						Enabled:            true,
						BootstrapPeers:     []string{"12D3KooWDnZtWxJCSZNUyPRmEUdmks9FigetxVuvaB3xuxn1hwmW@some-node0:9001"},
						OCRContractAddress: "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853",
						OCRKeyBundleID:     "dadasdsuidnasiudasnduas",
						ChainID:            "1337",
						TransmitterID:      "0x27118799c7368C2018052CD29072C0478C76d0e5",
						OnchainSigningStrategy: pkg.OnchainSigningStrategy{
							StrategyName: "single-chain",
							Config:       map[string]string{"evm": "dadasdsuidnasiudasnduas"},
						},
					},
				},
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

func TestProposeJobSpec_VerifyPreconditions_EVM(t *testing.T) {
	j := jobs.ProposeJobSpec{}
	var env cldf.Environment

	base := jobs.ProposeJobSpecInput{
		Environment: "test",
		Domain:      "cre",
		DONName:     "test-don",
		JobName:     "evm-test",
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: "d"},
			{Key: "environment", Value: "e"},
			{Key: "product", Value: offchain.ProductLabel},
		},
		Template: job_types.EVM,
	}

	validEVMInputs := func() job_types.JobSpecInput {
		return job_types.JobSpecInput{
			"command": "/usr/local/bin/evm",
			"config":  `{"chainId":1337,"network":"evm"}`,
			"oracleFactory": pkg.OracleFactory{
				Enabled:            true,
				BootstrapPeers:     []string{"12D3KooWDnZtWxJCSZNUyPRmEUdmks9FigetxVuvaB3xuxn1hwmW@workflow-node0:5001"},
				OCRContractAddress: "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853",
				OCRKeyBundleID:     "qbwjdbdywefeiwfiewb",
				ChainID:            "1337",
				TransmitterID:      "0x27118799c7368C2018052CD29072C0478C76d0e5",
				OnchainSigningStrategy: pkg.OnchainSigningStrategy{
					StrategyName: "single-chain",
					Config:       map[string]string{"evm": "deadbeefcafebabefeedface"},
				},
			},
		}
	}

	t.Run("valid evm spec passes", func(t *testing.T) {
		in := base
		in.Inputs = validEVMInputs()
		require.NoError(t, j.VerifyPreconditions(env, in))
	})

	t.Run("valid evm spec passes (oracleFactory as map[string]any)", func(t *testing.T) {
		in := base
		in.Inputs = validEVMInputs()
		in.Inputs["oracleFactory"] = map[string]any{
			"enabled":            true,
			"bootstrapPeers":     []string{"12D3KooWDnZtWxJCSZNUyPRmEUdmks9FigetxVuvaB3xuxn1hwmW@workflow-node0:5001"},
			"ocrContractAddress": "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853",
			"ocrKeyBundleID":     "qbwjdbdywefeiwfiewb",
			"chainID":            "1337",
			"transmitterID":      "0x27118799c7368C2018052CD29072C0478C76d0e5",
			"onchainSigningStrategy": map[string]any{
				"strategyName": "single-chain",
				"config":       map[string]string{"evm": "deadbeefcafebabefeedface"},
			},
		}
		require.NoError(t, j.VerifyPreconditions(env, in))
	})

	type negCase struct {
		name    string
		mutate  func(job_types.JobSpecInput)
		wantEnd string // appended to the common prefix
	}

	const prefix = "invalid inputs for EVM job spec: "

	cases := []negCase{
		// command
		{"missing command", func(m job_types.JobSpecInput) { delete(m, "command") }, "command is required and must be a string"},
		{"empty command", func(m job_types.JobSpecInput) { m["command"] = "   " }, "command is required and must be a string"},
		{"non-string command", func(m job_types.JobSpecInput) { m["command"] = nil }, "command is required and must be a string"},

		// config
		{"missing config", func(m job_types.JobSpecInput) { delete(m, "config") }, "config is required and must be a string"},
		{"empty config", func(m job_types.JobSpecInput) { m["config"] = "" }, "config is required and must be a string"},
		{"non-string config", func(m job_types.JobSpecInput) { m["config"] = nil }, "config is required and must be a string"},

		// oracleFactory presence/type/enabled
		{"missing oracleFactory", func(m job_types.JobSpecInput) { delete(m, "oracleFactory") }, "oracleFactory is required"},
		{"oracleFactory wrong type", func(m job_types.JobSpecInput) { m["oracleFactory"] = "not-a-factory" }, "cannot unmarshal !!str `not-a-f...` into pkg.OracleFactory"},
		{"oracleFactory present but disabled", func(m job_types.JobSpecInput) {
			of := m["oracleFactory"].(pkg.OracleFactory)
			of.Enabled = false
			m["oracleFactory"] = of
		}, "oracleFactory.enabled must be true for EVM jobs"},

		// bootstrapPeers
		{"enabled=true but missing bootstrapPeers", func(m job_types.JobSpecInput) {
			of := m["oracleFactory"].(pkg.OracleFactory)
			of.BootstrapPeers = nil
			m["oracleFactory"] = of
		}, "oracleFactory.bootstrapPeers is required"},
		{"enabled=true but invalid bootstrapPeers format", func(m job_types.JobSpecInput) {
			of := m["oracleFactory"].(pkg.OracleFactory)
			of.BootstrapPeers = []string{"not-a-peer"}
			m["oracleFactory"] = of
		}, "oracleFactory.bootstrapPeers is invalid"},

		// ocrContractAddress
		{"enabled=true but missing ocrContractAddress", func(m job_types.JobSpecInput) {
			of := m["oracleFactory"].(pkg.OracleFactory)
			of.OCRContractAddress = "   "
			m["oracleFactory"] = of
		}, "oracleFactory.ocrContractAddress is required"},
		{"enabled=true but invalid ocrContractAddress", func(m job_types.JobSpecInput) {
			of := m["oracleFactory"].(pkg.OracleFactory)
			of.OCRContractAddress = "0xnotanaddress"
			m["oracleFactory"] = of
		}, "oracleFactory.ocrContractAddress is invalid"},

		// ocrKeyBundleID
		{"enabled=true but missing ocrKeyBundleID", func(m job_types.JobSpecInput) {
			of := m["oracleFactory"].(pkg.OracleFactory)
			of.OCRKeyBundleID = ""
			m["oracleFactory"] = of
		}, "oracleFactory.ocrKeyBundleID is required"},

		// chainID
		{"enabled=true but missing chainID", func(m job_types.JobSpecInput) {
			of := m["oracleFactory"].(pkg.OracleFactory)
			of.ChainID = ""
			m["oracleFactory"] = of
		}, "oracleFactory.chainID is required"},
		{"enabled=true but invalid chainID", func(m job_types.JobSpecInput) {
			of := m["oracleFactory"].(pkg.OracleFactory)
			of.ChainID = "bogus"
			m["oracleFactory"] = of
		}, "oracleFactory.chainID is invalid"},

		// transmitterID
		{"enabled=true but missing transmitterID", func(m job_types.JobSpecInput) {
			of := m["oracleFactory"].(pkg.OracleFactory)
			of.TransmitterID = " "
			m["oracleFactory"] = of
		}, "oracleFactory.transmitterID is required"},

		// signing strategy
		{"enabled=true but missing strategyName", func(m job_types.JobSpecInput) {
			of := m["oracleFactory"].(pkg.OracleFactory)
			of.OnchainSigningStrategy.StrategyName = ""
			m["oracleFactory"] = of
		}, "oracleFactory.onchainSigningStrategy.strategyName is required"},
		{"enabled=true but missing signing config map", func(m job_types.JobSpecInput) {
			of := m["oracleFactory"].(pkg.OracleFactory)
			of.OnchainSigningStrategy.Config = nil
			m["oracleFactory"] = of
		}, "oracleFactory.onchainSigningStrategy.config is required"},
		{"enabled=true but missing config.evm entry", func(m job_types.JobSpecInput) {
			of := m["oracleFactory"].(pkg.OracleFactory)
			of.OnchainSigningStrategy.Config = map[string]string{
				"something-else": "value"}
			m["oracleFactory"] = of
		}, "oracleFactory.onchainSigningStrategy.config.evm is required"},
		{"enabled=true but empty config.evm entry", func(m job_types.JobSpecInput) {
			of := m["oracleFactory"].(pkg.OracleFactory)
			of.OnchainSigningStrategy.Config = map[string]string{"evm": "   "}
			m["oracleFactory"] = of
		}, "oracleFactory.onchainSigningStrategy.config.evm is required"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := base
			in.Inputs = validEVMInputs()
			tc.mutate(in.Inputs)

			err := j.VerifyPreconditions(env, in)
			require.Error(t, err)
			assert.Contains(t, err.Error(), prefix)
			assert.Contains(t, err.Error(), tc.wantEnd)
		})
	}
}

func TestProposeJobSpec_Apply(t *testing.T) {
	quarantine.Flaky(t, "DX-1893")
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
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
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
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
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
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"contractQualifier": "ocr3-contract-qualifier",
				"chainSelector":     strconv.FormatUint(chainSelector, 10),
			},
		}

		out, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.NoError(t, err)
		assert.Len(t, out.Reports, 1)

		reqs, err := testEnv.TestJD.ListProposedJobRequests()
		require.NoError(t, err)

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
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				// Missing "chainSelector"
				"contractQualifier": "ocr-contract-qualifier",
			},
		}

		_, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get OCR3 contract address for chain selector 0 and qualifier ocr-contract-qualifier")
		assert.Contains(t, err.Error(), "failed to get OCR3 contract address for chain selector 0 and qualifier ocr-contract-qualifier")
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
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"templateName":       "worker-ocr3",
				"contractQualifier":  "ocr3-contract-qualifier",
				"chainSelectorEVM":   strconv.FormatUint(chainSelector, 10),
				"chainSelectorAptos": strconv.FormatUint(testEnv.AptosSelector, 10),
				"bootstrapperOCR3Urls": []string{
					"12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001",
				},
			},
		}

		out, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.NoError(t, err)
		assert.Len(t, out.Reports, 1)

		reqs, err := testEnv.TestJD.ListProposedJobRequests()
		require.NoError(t, err)

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
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				// missing `templateName`
				"contractQualifier": "ocr3-contract-qualifier",
			},
		}

		_, err = jobs.ProposeJobSpec{}.Apply(*env, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to convert inputs to OCR3 job input")
		assert.Contains(t, err.Error(), "templateName is required and must be a non-empty string")
	})

	t.Run("successful evm job distribution", func(t *testing.T) {
		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "capability_evm_1337-1337",
			DONName:     test.DONName,
			Template:    job_types.EVM,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"command":       "/usr/local/bin/evm",
				"config":        `{"chainId":1337,"network":"evm","logTriggerPollInterval":1500000000,"creForwarderAddress":"0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0","receiverGasMinimum":500,"nodeAddress":"0x27118799c7368C2018052CD29072C0478C76d0e5"}`,
				"externalJobID": "2d462183-acf3-489e-926c-464342578a38",
				"oracleFactory": pkg.OracleFactory{
					Enabled:            true,
					BootstrapPeers:     []string{"12D3KooWDnZtWxJCSZNUyPRmEUdmks9FigetxVuvaB3xuxn1hwmW@workflow-node0:5001"},
					OCRContractAddress: "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853",
					OCRKeyBundleID:     "c6f25ead88206730f25b3b44cfcf909f0a69b2408f3ea7de8e408bafce7ebae5",
					ChainID:            "1337",
					TransmitterID:      "0x27118799c7368C2018052CD29072C0478C76d0e5",
					OnchainSigningStrategy: pkg.OnchainSigningStrategy{
						StrategyName: "single-chain",
						Config:       map[string]string{"evm": "c6f25ead88206730f25b3b44cfcf909f0a69b2408f3ea7de8e408bafce7ebae5"},
					},
				},
			},
		}

		out, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.NoError(t, err)
		assert.Len(t, out.Reports, 1)

		reqs, err := testEnv.TestJD.ListProposedJobRequests()
		require.NoError(t, err)

		for _, req := range reqs {
			if !strings.Contains(req.Spec, `name = "capability_evm_1337-1337"`) {
				continue
			}

			t.Logf("Job Spec:\n%s", req.Spec)

			assert.Contains(t, req.Spec, `type = "standardcapabilities"`)
			assert.Contains(t, req.Spec, `name = "capability_evm_1337-1337"`)
			assert.Contains(t, req.Spec, `externalJobID = "2d462183-acf3-489e-926c-464342578a38"`)
			assert.Contains(t, req.Spec, `command = "/usr/local/bin/evm"`)

			// config (embedded JSON string)
			assert.Contains(t, req.Spec, `config = """`)
			assert.Contains(t, req.Spec, `"network":"evm"`)
			assert.Contains(t, req.Spec, `"chainId":1337`)
			assert.Contains(t, req.Spec, `"creForwarderAddress":"0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"`)
			assert.Contains(t, req.Spec, `"receiverGasMinimum":500`)
			assert.Contains(t, req.Spec, `"nodeAddress":"0x27118799c7368C2018052CD29072C0478C76d0e5"`)

			// oracle factory block
			assert.Contains(t, req.Spec, `[oracle_factory]`)
			assert.Contains(t, req.Spec, `enabled = true`)
			assert.Regexp(t,
				`bootstrap_peers\s*=\s*\[\s*"12D3KooWDnZtWxJCSZNUyPRmEUdmks9FigetxVuvaB3xuxn1hwmW@workflow-node0:5001"\s*\]`,
				req.Spec,
			)
			assert.Contains(t, req.Spec, `ocr_contract_address = "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853"`)
			assert.Contains(t, req.Spec, `ocr_key_bundle_id = "c6f25ead88206730f25b3b44cfcf909f0a69b2408f3ea7de8e408bafce7ebae5"`)
			assert.Contains(t, req.Spec, `chain_id = "1337"`)
			assert.Contains(t, req.Spec, `transmitter_id = "0x27118799c7368C2018052CD29072C0478C76d0e5"`)

			assert.Contains(t, req.Spec, `[oracle_factory.onchainSigningStrategy]`)
			assert.Contains(t, req.Spec, `strategyName = "single-chain"`)
			assert.Contains(t, req.Spec, `[oracle_factory.onchainSigningStrategy.config]`)
			assert.Contains(t, req.Spec, `evm = "c6f25ead88206730f25b3b44cfcf909f0a69b2408f3ea7de8e408bafce7ebae5"`)
		}
	})

	t.Run("failed evm job distribution due to bad input", func(t *testing.T) {
		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "capability_evm_1337-1337",
			DONName:     test.DONName,
			Template:    job_types.EVM,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				// Intentionally omit "command"
				"config":        `{"chainId":1337,"network":"evm"}`,
				"externalJobID": "an-evm-job-id",
				"oracleFactory": pkg.OracleFactory{
					Enabled:            true,
					BootstrapPeers:     []string{"12D3KooWDnZtWxJCSZNUyPRmEUdmks9FigetxVuvaB3xuxn1hwmW@workflow-node0:5001"},
					OCRContractAddress: "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853",
					OCRKeyBundleID:     "deadbeef",
					ChainID:            "1337",
					TransmitterID:      "0x0000000000000000000000000000000000000001",
					OnchainSigningStrategy: pkg.OnchainSigningStrategy{
						StrategyName: "single-chain",
						Config:       map[string]string{"evm": "deadbeef"},
					},
				},
			},
		}

		_, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to convert inputs to standard capability job")
		assert.Contains(t, err.Error(), "command is required and must be a string")
	})
	t.Run("successful http trigger job distribution", func(t *testing.T) {
		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "http-trigger-job",
			DONName:     test.DONName,
			Template:    job_types.HTTPTrigger,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"command":       "http_trigger",
				"config":        `{}`,
				"externalJobID": "http-trigger-job-id",
			},
		}

		out, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.NoError(t, err)
		assert.Len(t, out.Reports, 1)

		reqs, err := testEnv.TestJD.ListProposedJobRequests()
		require.NoError(t, err)

		for _, req := range reqs {
			if !strings.Contains(req.Spec, `name = "http-trigger-job"`) {
				continue
			}
			// log each spec in readable format
			t.Logf("HTTP Trigger Job Spec:\n%s", req.Spec)
			assert.Contains(t, req.Spec, `name = "http-trigger-job"`)
			assert.Contains(t, req.Spec, `command = "http_trigger"`)
			assert.Contains(t, req.Spec, `config = """{}"""`)
			assert.Contains(t, req.Spec, `externalJobID = "http-trigger-job-id"`)
		}
	})

	t.Run("successful http action job distribution", func(t *testing.T) {
		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "http-action-job",
			DONName:     test.DONName,
			Template:    job_types.HTTPAction,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"command":       "http_action",
				"config":        `{"proxyMode": "direct"}`,
				"externalJobID": "http-action-job-id",
			},
		}

		out, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.NoError(t, err)
		assert.Len(t, out.Reports, 1)

		reqs, err := testEnv.TestJD.ListProposedJobRequests()
		require.NoError(t, err)
		for _, req := range reqs {
			if !strings.Contains(req.Spec, `name = "http-action-job"`) {
				continue
			}
			// log each spec in readable format
			t.Logf("HTTP Action Job Spec:\n%s", req.Spec)
			assert.Contains(t, req.Spec, `name = "http-action-job"`)
			assert.Contains(t, req.Spec, `command = "http_action"`)
			assert.Contains(t, req.Spec, `config = """{"proxyMode": "direct"}"""`)
			assert.Contains(t, req.Spec, `externalJobID = "http-action-job-id"`)
		}
	})

	t.Run("failed http trigger job distribution due to bad input", func(t *testing.T) {
		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "http-trigger-job",
			DONName:     test.DONName,
			Template:    job_types.HTTPTrigger,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				// Missing "command"
				"config":        `{}`,
				"externalJobID": "http-trigger-job-id",
			},
		}

		_, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to convert inputs to standard capability job")
		assert.Contains(t, err.Error(), "command is required and must be a string")
	})

	t.Run("failed http action job distribution due to bad input", func(t *testing.T) {
		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "http-action-job",
			DONName:     test.DONName,
			Template:    job_types.HTTPAction,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"config":        `{"proxyMode": "direct"}`,
				"externalJobID": "http-action-job-id",
			},
		}

		_, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to convert inputs to standard capability job")
		assert.Contains(t, err.Error(), "command is required and must be a string")
	})

	t.Run("failed evm job distribution due to bad input", func(t *testing.T) {
		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "capability_evm_1337-1337",
			Template:    job_types.EVM, // if unavailable, use the same template you use for cron but with evm inputs.
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"config":        `{"chainId":1337,"network":"evm"}`,
				"externalJobID": "an-evm-job-id",
				"oracleFactory": pkg.OracleFactory{
					Enabled: true,
					// Provide partial/OK oracle factory so the error specifically points to missing command.
					BootstrapPeers:     []string{"12D3KooWDnZtWxJCSZNUyPRmEUdmks9FigetxVuvaB3xuxn1hwmW@workflow-node0:5001"},
					OCRContractAddress: "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853",
					OCRKeyBundleID:     "deadbeef",
					ChainID:            "1337",
					TransmitterID:      "0x0000000000000000000000000000000000000001",
					OnchainSigningStrategy: pkg.OnchainSigningStrategy{
						StrategyName: "single-chain",
						Config:       map[string]string{"evm": "deadbeef"},
					},
				},
			},
		}

		_, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to convert inputs to standard capability job")
		assert.Contains(t, err.Error(), "command is required and must be a string")
	})

	t.Run("successful bootstrap distribution", func(t *testing.T) {
		chainSelector := testEnv.RegistrySelector
		ds := datastore.NewMemoryDataStore()

		err := ds.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSelector,
			Type:          datastore.ContractType(ocr3.OCR3Capability),
			Version:       semver.MustParse("1.0.0"),
			Address:       "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853",
			Qualifier:     "vault_1_plugin",
		})
		require.NoError(t, err)

		err = ds.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSelector,
			Type:          datastore.ContractType(ocr3.OCR3Capability),
			Version:       semver.MustParse("1.0.0"),
			Address:       "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B",
			Qualifier:     "vault_1_dkg",
		})
		require.NoError(t, err)

		env.DataStore = ds.Seal()

		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "vault-bootstrappers",
			DONName:     test.DONName,
			Template:    job_types.BootstrapVault,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"contractQualifierPrefix": "vault_1",
				"chainSelector":           strconv.FormatUint(chainSelector, 10),
			},
		}

		out, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.NoError(t, err)

		output := out.Reports[0].Output.(sequences.ProposeVaultBootstrapJobsOutput)
		assert.Len(t, output.Specs, 1)

		jobs := []struct {
			Address       string
			JobNameSuffix string
		}{
			{
				Address:       "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853",
				JobNameSuffix: " (Plugin)",
			},
			{
				Address:       "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B",
				JobNameSuffix: " (DKG)",
			},
		}
		for _, specs := range output.Specs {
			for i, s := range specs {
				assert.Contains(t, s, `type = "bootstrap"`)
				assert.Contains(t, s, `name = "vault-bootstrappers`+jobs[i].JobNameSuffix)
				assert.Contains(t, s, fmt.Sprintf(`contractID = "%s"`, jobs[i].Address))
			}
		}

		propJobs, err := testEnv.TestJD.ListProposedJobRequests()
		require.NoError(t, err)

		foundSet := map[string]bool{}
		for _, p := range propJobs {
			if strings.Contains(p.Spec, `name = "vault-bootstrappers (Plugin)`) {
				assert.Contains(t, p.Spec, `contractID = "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853"`)
				foundSet["plugin"] = true
			}

			if strings.Contains(p.Spec, `name = "vault-bootstrappers (DKG)`) {
				assert.Contains(t, p.Spec, `contractID = "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"`)
				foundSet["dkg"] = true
			}
		}

		assert.Len(t, foundSet, 2)
	})

	t.Run("unsuccessful bootstrap distribution because contracts don't exist", func(t *testing.T) {
		chainSelector := testEnv.RegistrySelector
		ds := datastore.NewMemoryDataStore()

		env.DataStore = ds.Seal()

		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "vault-bootstrappers",
			DONName:     test.DONName,
			Template:    job_types.BootstrapVault,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"contractQualifierPrefix": "another_vault_1",
				"chainSelector":           strconv.FormatUint(chainSelector, 10),
			},
		}

		_, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		assert.ErrorContains(t, err, "failed to get Vault Plugin contract address")
	})

	t.Run("successful vault ocr3 job distribution", func(t *testing.T) {
		testEnv := test.SetupEnvV2(t, false)
		env := testEnv.Env

		chainSelector := testEnv.RegistrySelector
		ds := datastore.NewMemoryDataStore()

		err := ds.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSelector,
			Type:          datastore.ContractType(ocr3.OCR3Capability),
			Version:       semver.MustParse("1.0.0"),
			Address:       "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B",
			Qualifier:     "vault_1_plugin",
		})
		require.NoError(t, err)

		err = ds.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSelector,
			Type:          datastore.ContractType(ocr3.OCR3Capability),
			Version:       semver.MustParse("1.0.0"),
			Address:       "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853",
			Qualifier:     "vault_1_dkg",
		})
		require.NoError(t, err)

		env.DataStore = ds.Seal()

		input := jobs.ProposeJobSpecInput{
			Environment: "test",
			Domain:      "cre",
			JobName:     "vault-job",
			DONName:     test.DONName,
			Template:    job_types.OCR3,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"templateName":         "worker-vault",
				"contractQualifier":    "vault_1_plugin",
				"dkgContractQualifier": "vault_1_dkg",
				"chainSelectorEVM":     strconv.FormatUint(chainSelector, 10),
				"bootstrapperOCR3Urls": []string{
					"12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001",
				},
			},
		}

		out, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.NoError(t, err)
		assert.Len(t, out.Reports, 1)

		reqs, err := testEnv.TestJD.ListProposedJobRequests()
		require.NoError(t, err)

		expectedChainID := chainsel.TEST_90000001.EvmChainID

		for _, req := range reqs {
			if !strings.Contains(req.Spec, `pluginType = "vault-plugin"`) {
				continue
			}
			// log each spec in readable yaml format
			t.Logf("Job Spec:\n%s", req.Spec)
			assert.Contains(t, req.Spec, `name = "vault-job`)
			assert.Contains(t, req.Spec, `contractID = "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"`)
			assert.Contains(t, req.Spec, `p2pv2Bootstrappers = [
  "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001",
]`)
			assert.Contains(t, req.Spec, fmt.Sprintf(`chainID = "%d"`, expectedChainID))
			assert.Contains(t, req.Spec, `dkgContractID = "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853"`)
		}
	})

	t.Run("successful consensus job distribution", func(t *testing.T) {
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
			JobName:     "ocr3-consensus-job",
			DONName:     test.DONName,
			Template:    job_types.Consensus,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"command":            "consensus",
				"contractQualifier":  "ocr3-contract-qualifier",
				"chainSelectorEVM":   strconv.FormatUint(chainSelector, 10),
				"chainSelectorAptos": strconv.FormatUint(testEnv.AptosSelector, 10),
				"bootstrapPeers": []string{
					"12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001",
				},
			},
		}

		out, err := jobs.ProposeJobSpec{}.Apply(*env, input)
		require.NoError(t, err)
		assert.Len(t, out.Reports, 1)

		reqs, err := testEnv.TestJD.ListProposedJobRequests()
		require.NoError(t, err)

		expectedChainID := chainsel.TEST_90000001.EvmChainID

		for _, req := range reqs {
			if !strings.Contains(req.Spec, `name = "ocr3-consensus-job"`) {
				continue
			}
			// log each spec in readable yaml format
			t.Logf("Job Spec:\n%s", req.Spec)
			assert.Contains(t, req.Spec, `name = "ocr3-consensus-job"`)
			assert.Contains(t, req.Spec, `bootstrap_peers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]`)
			assert.Contains(t, req.Spec, fmt.Sprintf(`chain_id = "%d"`, expectedChainID))
			assert.Contains(t, req.Spec, `command = "consensus"`)
			assert.Contains(t, req.Spec, `config = """"""`)
			assert.Contains(t, req.Spec, `[oracle_factory]`)
			assert.Contains(t, req.Spec, `enabled = true`)
			assert.Contains(t, req.Spec, `ocr_contract_address = "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"`)
			assert.Contains(t, req.Spec, `strategyName = "multi-chain"`)
			assert.Contains(t, req.Spec, `ocr_key_bundle_id = "fake_orc_bundle_evm"`)
		}
	})

	t.Run("failed consensus job distribution", func(t *testing.T) {
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
			JobName:     "ocr3-consensus-job",
			DONName:     test.DONName,
			Template:    job_types.Consensus,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: "test"},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				// missing `command`
			},
		}

		_, err = jobs.ProposeJobSpec{}.Apply(*env, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to convert inputs to standard capability job")
		assert.Contains(t, err.Error(), "command is required and must be a string")
	})
}
