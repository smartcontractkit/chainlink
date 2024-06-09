package workflows_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
)

func TestDelegate_JobSpecValidator(t *testing.T) {
	t.Parallel()
	validName := "ten bytes!"
	var tt = []struct {
		name           string
		workflowTomlFn func() string
		valid          bool
	}{
		{
			"not a hex owner",
			func() string {
				workflowId := "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef"
				workflowOwner := "00000000000000000000000000000000000000aZ"
				return testspecs.GenerateWorkflowSpec(workflowId, workflowOwner, "1234567890", "").Toml()
			},
			false,
		},
		{
			"missing workflow field",
			func() string {
				id := "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef"
				owner := "00000000000000000000000000000000000000aa"
				return testspecs.GenerateWorkflowSpec(id, owner, validName, "").Toml()
			},
			false,
		},

		{
			"null workflow",
			func() string {
				id := "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef"
				owner := "00000000000000000000000000000000000000aa"
				return testspecs.GenerateWorkflowSpec(id, owner, validName, "{}").Toml()
			},
			false,
		},

		{
			"missing name",
			func() string {
				id := "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef"
				owner := "00000000000000000000000000000000000000aa"
				wf := `
triggers: []
consensus: []
targets: []
`
				return testspecs.GenerateWorkflowSpec(id, owner, "", wf).Toml()
			},
			false,
		},

		{
			"minimal passing workflow",
			func() string {
				id := "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef"
				owner := "00000000000000000000000000000000000000aa"
				wf := `
triggers: []
consensus: []
targets: []
`
				return testspecs.GenerateWorkflowSpec(id, owner, validName, wf).Toml()
			},
			true,
		},

		{
			"name too long",
			func() string {
				id := "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef"
				owner := "00000000000000000000000000000000000000aa"
				wf := `
triggers: []
consensus: []
targets: []
`
				return testspecs.GenerateWorkflowSpec(id, owner, validName+"1", wf).Toml()
			},
			false,
		},

		// Taken from jobs controller test, as we want to fail early without a db / slow test dependency
		{
			"valid full spec",
			func() string {
				id := "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef"
				owner := "00000000000000000000000000000000000000aa"
				workflow := `
triggers:
  - id: "mercury-trigger@1.0.0"
    config:
      feedIds:
        - "0x1111111111111111111100000000000000000000000000000000000000000000"
        - "0x2222222222222222222200000000000000000000000000000000000000000000"
        - "0x3333333333333333333300000000000000000000000000000000000000000000"

consensus:
  - id: "offchain_reporting@2.0.0"
    ref: "evm_median"
    inputs:
      observations:
        - "$(trigger.outputs)"
    config:
      aggregation_method: "data_feeds_2_0"
      aggregation_config:
        "0x1111111111111111111100000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
        "0x2222222222222222222200000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
        "0x3333333333333333333300000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
      encoder: "EVM"
      encoder_config:
        abi: "mercury_reports bytes[]"

targets:
  - id: "write_polygon-testnet-mumbai@3.0.0"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x3F3554832c636721F1fD1822Ccca0354576741Ef"
      params: ["$(report)"]
      abi: "receive(report bytes)"
  - id: "write_ethereum-testnet-sepolia@4.0.0"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x54e220867af6683aE6DcBF535B4f952cB5116510"
      params: ["$(report)"]
      abi: "receive(report bytes)"
`
				return testspecs.GenerateWorkflowSpec(id, owner, validName, workflow).Toml()
			},
			true,
		},

		{
			"parse error",
			func() string {
				return `
invalid syntax{{{{
`
			},
			false,
		},

		{
			"invalid job type",
			func() string {
				return `
type = "work flows"
schemaVersion = 1
`
			},
			false,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := workflows.ValidatedWorkflowSpec(tc.workflowTomlFn())
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
