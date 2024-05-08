package test_workflow_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/test_workflow"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/wasm/wrokflowtesting"
)

func TestWorkflow(t *testing.T) {
	workflow, err := test_workflow.CreateWorkflow()
	require.NoError(t, err)
	trigger := &test_workflow.MercuryTriggerResponse{
		Values: map[string]decimal.Decimal{
			"123": decimal.NewFromFloat(1.00),
			"456": decimal.NewFromFloat(1.25),
			"789": decimal.NewFromFloat(1.50),
		},
		Decimals: map[string]int{
			"123": 19,
			"456": 19,
			"789": 8,
		},
		Metadata: test_workflow.TriggerMetadata{TriggerRef: "Mercury"},
	}

	registry := wrokflowtesting.NewRegistry()
	require.NoError(t, AddMercuryTriggerToRegistry("mercury-trigger", registry, trigger))

	writeTarget := &ChainWriterMock{}
	writeTarget.AddTarget("write", registry)

	runner := wrokflowtesting.NewRunner(registry)
	require.NoError(t, runner.Run(workflow))

	assert.Len(t, writeTarget.Seen, 1)
	assert.Equal(t, trigger, writeTarget.Seen[0])
}
