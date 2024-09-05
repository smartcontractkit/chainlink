package workflows_test

import (
	"testing"

	commonworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
)

const anyYamlSpec = `
name: "wf-name"
owner: "0x00000000000000000000000000000000000000aa"
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

func TestYamlSpecFactory_GetSpec(t *testing.T) {
	t.Parallel()

	actual, err := workflows.YAMLSpecFactory{}.GetSpec([]byte(anyYamlSpec), []byte{})
	require.NoError(t, err)

	expected, err := commonworkflows.ParseWorkflowSpecYaml(anyYamlSpec)
	require.NoError(t, err)

	require.Equal(t, expected, actual)
}

func TestYamlSpecFactory_GetRawSpec(t *testing.T) {
	t.Parallel()

	actual, err := workflows.YAMLSpecFactory{}.GetRawSpec(anyYamlSpec)
	require.NoError(t, err)
	require.Equal(t, []byte(anyYamlSpec), actual)
}
