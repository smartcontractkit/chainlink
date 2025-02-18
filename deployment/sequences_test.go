package deployment

import (
	"strconv"
	"testing"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
)

type SeqDeps struct {
	Operation func(a, b int) (int, error)
}

type SeqInput struct {
	A int
	B int
	C int
	D int
}

// Create an instance using NewOperation
var myop = NewOperation("V1", "Operation that computes two numbers", func(ctx OpContext, deps OpDeps, input OpInput) (int, error) {
	return deps.Operation(input.A, input.B)
})

// The sequence uses the underlying "myop" operation to calculate more complex operation
var seq = NewSequence("V1", "Adds four numbers in sequence using the myop operation", func(env OpEnv, deps SeqDeps, input SeqInput) (output int, err error) {
	// Add first two numbers
	opDeps := OpDeps{
		Operation: deps.Operation,
	}
	rep, err := ExecuteOp(env, myop, opDeps, OpInput{A: input.A, B: input.B})
	if err != nil {
		return 0, err
	}

	// Add the result to the third number
	rep, err = ExecuteOp(env, myop, opDeps, OpInput{A: rep.Output, B: input.C})
	if err != nil {
		return 0, err
	}

	// Add the result to the fourth number
	rep, err = ExecuteOp(env, myop, opDeps, OpInput{A: rep.Output, B: input.D})
	if err != nil {
		return 0, err
	}

	return rep.Output, nil
})

func TestCreateNewSequence(t *testing.T) {
	env := OpEnv{
		Log:      logger.TestLogger(t),
		Reporter: NewMemoryReporter([]ReportAny{}),
	}

	deps := SeqDeps{
		Operation: func(a, b int) (int, error) { return a + b, nil },
	}
	// Execute the sequence
	report := ExecuteSeq(env, *seq, deps, SeqInput{A: 1, B: 2, C: 3, D: 4})

	// Check the report
	require.NoError(t, report.Err)
	require.Equal(t, 10, report.Output)
	require.Equal(t, 3, len(report.SubOps))
}

func TestSequenceRetry(t *testing.T) {
	env := OpEnv{
		Log:      logger.TestLogger(t),
		Reporter: NewMemoryReporter([]ReportAny{}),
	}

	// We force an error in the second operation, but only once
	counter := 0
	cache := make(map[string]bool)
	deps := SeqDeps{
		Operation: func(a, b int) (int, error) {
			counter++
			// We force second operation to fail
			if counter == 2 {
				return 0, errors.New("forced error")
			}

			// Inputs are cached. If a successful operation is retried with the same inputs, this errors out.
			// It should never enter here as the Execute will not compute and return the cached report
			id := strconv.Itoa(a) + strconv.Itoa(b)
			if cache[id] {
				return 0, errors.New("op already executed error")
			}
			cache[id] = true

			return a + b, nil
		},
	}

	// Execute the sequence
	report := ExecuteSeq(env, *seq, deps, SeqInput{A: 1, B: 2, C: 3, D: 4})

	// Check the report
	require.Error(t, report.Err)
	require.Equal(t, 2, len(report.SubOps))

	// Retry the sequence
	report = ExecuteSeq(env, *seq, deps, SeqInput{A: 1, B: 2, C: 3, D: 4}, report.ID)

	// Check the report
	require.NoError(t, report.Err)
	require.Equal(t, 10, report.Output)
	require.Equal(t, 3, len(report.SubOps))

	// The total number of reports (2 sequences and 3+1 ops)
	require.Equal(t, 6, len(env.Reporter.GetReports()))
}

type SeqInput2 struct {
	A int
	B int
	C int
	D int
	E int
}

// We define a new sequence that uses the previous sequence
var seq2 = NewSequence("V1", "Adds five numbers in sequence using the another sequence", func(env OpEnv, deps SeqDeps, input SeqInput2) (output int, err error) {
	// Add first two numbers
	opDeps := OpDeps{
		Operation: deps.Operation,
	}
	seqRep := ExecuteSeq(env, *seq, deps, SeqInput{A: input.A, B: input.B, C: input.C, D: input.D})
	if seqRep.Err != nil {
		return 0, seqRep.Err
	}

	// Add the result to the third number
	opRep, err := ExecuteOp(env, myop, opDeps, OpInput{A: seqRep.Output, B: input.B})
	if err != nil {
		return 0, err
	}

	return opRep.Output, nil
})

func TestSequenceComposition(t *testing.T) {
	deps := SeqDeps{
		Operation: func(a, b int) (int, error) { return a + b, nil },
	}

	env := OpEnv{
		Log:      logger.TestLogger(t),
		Reporter: NewMemoryReporter([]ReportAny{}),
	}
	// Execute the sequence
	report := ExecuteSeq(env, *seq2, deps, SeqInput2{A: 1, B: 2, C: 3, D: 4, E: 5})

	// Check the report
	require.NoError(t, report.Err)
	require.Equal(t, 12, report.Output)
	// The total number of subops (1 seq (3 ops inside) and 1 op)
	require.Equal(t, 5, len(report.SubOps))

	// The total number of reports (2 sequences + (3 + 1) ops)
	reports := env.Reporter.GetReports()
	require.Equal(t, 6, len(reports))
}

func TestSequenceCompositionRetry(t *testing.T) {

	// We force an error in the second operation, but only once
	counter := 0
	cache := make(map[string]bool)
	deps := SeqDeps{
		Operation: func(a, b int) (int, error) {
			counter++
			// We force second operation to fail
			if counter == 2 {
				return 0, errors.New("forced error")
			}

			// Inputs are cached. If a successful operation is retried with the same inputs, this errors out.
			// It should never enter here as the Execute will not compute and return the cached report
			id := strconv.Itoa(a) + strconv.Itoa(b)
			if cache[id] {
				return 0, errors.New("op already executed error")
			}
			cache[id] = true

			return a + b, nil
		},
	}

	env := OpEnv{
		Log:      logger.TestLogger(t),
		Reporter: NewMemoryReporter([]ReportAny{}),
	}

	// Execute the sequence
	report := ExecuteSeq(env, *seq2, deps, SeqInput2{A: 1, B: 2, C: 3, D: 4, E: 5})
	// Check the report
	require.Error(t, report.Err)
	require.Equal(t, 3, len(report.SubOps))
	require.Equal(t, 4, len(env.Reporter.GetReports()))

	// Retry the sequence
	report = ExecuteSeq(env, *seq2, deps, SeqInput2{A: 1, B: 2, C: 3, D: 4, E: 5}, report.ID)
	// Check the report
	require.NoError(t, report.Err)
	require.Equal(t, 12, report.Output)
	// The total number of subops (1 seq (3 ops inside) and 1 op)
	require.Equal(t, 5, len(report.SubOps))

	// The total number of reports (sq1(4) + seq2(5) = 9; seq2(2 seqs + (3 + 1 ops))
	reports := env.Reporter.GetReports()
	require.Equal(t, 9, len(reports))
}
