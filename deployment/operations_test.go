package deployment

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
)

type OpDeps struct {
	Operation func(a, b int) (int, error)
}

type OpInput struct {
	A int
	B int
}

// Create an instance using NewOperation
var op = NewOperation("V1", "Description", func(ctx OpContext, deps OpDeps, input OpInput) (int, error) {
	return deps.Operation(input.A, input.B)
})

func TestOperation(t *testing.T) {
	// Build context
	ctx := OpContext{
		Log: logger.TestLogger(t),
	}

	// Build Deps
	deps := OpDeps{
		// The operation this time adds two numbers
		Operation: func(a, b int) (int, error) { return a + b, nil },
	}

	res, err := op.Execute(ctx, deps, OpInput{A: 1, B: 2})
	require.NoError(t, err)
	require.Equal(t, 3, res)
}

func TestExecuteOperation(t *testing.T) {
	// Env is defined by the user
	env := OpEnv{
		Log:      logger.TestLogger(t),
		Reporter: NewMemoryReporter([]ReportAny{}),
	}

	// Build Deps
	deps := OpDeps{
		// The operation this time multiply two numbers
		Operation: func(a, b int) (int, error) { return a * b, nil },
	}
	// Execute the operation
	report, err := ExecuteOp(env, op, deps, OpInput{A: 2, B: 2})
	require.NoError(t, err)
	require.Equal(t, 4, report.Output)

	// Check the report
	reports := env.Reporter.GetReports()
	require.Len(t, reports, 1)
	require.Equal(t, 4, reports[0].Output)
}
