package flakeytests

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockReporter struct {
	entries map[string][]string
}

func (m *mockReporter) Report(entries map[string][]string) error {
	m.entries = entries
	return nil
}

func newMockReporter() *mockReporter {
	return &mockReporter{entries: map[string][]string{}}
}

func TestParser(t *testing.T) {
	output := `
--- FAIL: TestLink (0.00s)
    --- FAIL: TestLink/1.1_link#01 (0.00s)
        currencies_test.go:325:
                Error Trace:    /Users/ccordenier/Development/chainlink/core/assets/currencies_test.go:325
                Error:          Not equal:
                                expected: "1.2 link"
                                actual  : "1.1 link"

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1 +1 @@
                                -1.2 link
                                +1.1 link
                Test:           TestLink/1.1_link#01
FAIL
FAIL    github.com/smartcontractkit/chainlink/v2/core/assets    0.338s
FAIL
`

	r := strings.NewReader(output)
	ts, err := parseOutput(r)
	require.NoError(t, err)

	assert.Len(t, ts, 1)
	assert.Len(t, ts["core/assets"], 1)
	assert.Equal(t, ts["core/assets"]["TestLink"], 1)
}

func TestParser_SuccessfulOutput(t *testing.T) {
	output := `
?       github.com/smartcontractkit/chainlink/v2/tools/flakeytests/cmd/runner   [no test files]
ok      github.com/smartcontractkit/chainlink/v2/tools/flakeytests      0.320s
`

	r := strings.NewReader(output)
	ts, err := parseOutput(r)
	require.NoError(t, err)
	assert.Len(t, ts, 0)
}

func TestRunner_WithFlake(t *testing.T) {
	output := `
--- FAIL: TestLink (0.00s)
    --- FAIL: TestLink/1.1_link#01 (0.00s)
        currencies_test.go:325:
                Error Trace:    /Users/ccordenier/Development/chainlink/core/assets/currencies_test.go:325
                Error:          Not equal:
                                expected: "1.2 link"
                                actual  : "1.1 link"

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1 +1 @@
                                -1.2 link
                                +1.1 link
                Test:           TestLink/1.1_link#01
FAIL
FAIL    github.com/smartcontractkit/chainlink/v2/core/assets    0.338s
FAIL
`
	m := newMockReporter()
	r := &Runner{
		numReruns: 2,
		readers:   []io.Reader{strings.NewReader(output)},
		runTestFn: func(pkg string, testNames []string, numReruns int, w io.Writer) error {
			_, err := w.Write([]byte(output))
			return err
		},
		parse:    parseOutput,
		reporter: m,
	}

	// This will report a flake since we've mocked the rerun
	// to only report one failure (not two as expected).
	err := r.Run()
	require.NoError(t, err)
	assert.Len(t, m.entries, 1)
	assert.Equal(t, m.entries["core/assets"], []string{"TestLink"})
}

func TestRunner_AllFailures(t *testing.T) {
	output := `
--- FAIL: TestLink (0.00s)
    --- FAIL: TestLink/1.1_link#01 (0.00s)
        currencies_test.go:325:
                Error Trace:    /Users/ccordenier/Development/chainlink/core/assets/currencies_test.go:325
                Error:          Not equal:
                                expected: "1.2 link"
                                actual  : "1.1 link"

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1 +1 @@
                                -1.2 link
                                +1.1 link
                Test:           TestLink/1.1_link#01
FAIL
FAIL    github.com/smartcontractkit/chainlink/v2/core/assets    0.338s
FAIL
`

	rerunOutput := `
--- FAIL: TestLink (0.00s)
    --- FAIL: TestLink/1.1_link#01 (0.00s)
        currencies_test.go:325:
                Error Trace:    /Users/ccordenier/Development/chainlink/core/assets/currencies_test.go:325
                Error:          Not equal:
                                expected: "1.2 link"
                                actual  : "1.1 link"

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1 +1 @@
                                -1.2 link
                                +1.1 link
                Test:           TestLink/1.1_link#01
--- FAIL: TestLink (0.00s)
    --- FAIL: TestLink/1.1_link#01 (0.00s)
        currencies_test.go:325:
                Error Trace:    /Users/ccordenier/Development/chainlink/core/assets/currencies_test.go:325
                Error:          Not equal:
                                expected: "1.2 link"
                                actual  : "1.1 link"

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1 +1 @@
                                -1.2 link
                                +1.1 link
                Test:           TestLink/1.1_link#01
FAIL
FAIL    github.com/smartcontractkit/chainlink/v2/core/assets    0.315s
FAIL
`
	m := newMockReporter()
	r := &Runner{
		numReruns: 2,
		readers:   []io.Reader{strings.NewReader(output)},
		runTestFn: func(pkg string, testNames []string, numReruns int, w io.Writer) error {
			_, err := w.Write([]byte(rerunOutput))
			return err
		},
		parse:    parseOutput,
		reporter: m,
	}

	err := r.Run()
	require.NoError(t, err)
	assert.Len(t, m.entries, 0)
}

func TestRunner_RerunSuccessful(t *testing.T) {
	output := `
--- FAIL: TestLink (0.00s)
    --- FAIL: TestLink/1.1_link#01 (0.00s)
        currencies_test.go:325:
                Error Trace:    /Users/ccordenier/Development/chainlink/core/assets/currencies_test.go:325
                Error:          Not equal:
                                expected: "1.2 link"
                                actual  : "1.1 link"

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1 +1 @@
                                -1.2 link
                                +1.1 link
                Test:           TestLink/1.1_link#01
FAIL
FAIL    github.com/smartcontractkit/chainlink/v2/core/assets    0.338s
FAIL
`

	rerunOutput := `
ok      github.com/smartcontractkit/chainlink/v2/core/assets      0.320s
`
	m := newMockReporter()
	r := &Runner{
		numReruns: 2,
		readers:   []io.Reader{strings.NewReader(output)},
		runTestFn: func(pkg string, testNames []string, numReruns int, w io.Writer) error {
			_, err := w.Write([]byte(rerunOutput))
			return err
		},
		parse:    parseOutput,
		reporter: m,
	}

	err := r.Run()
	require.NoError(t, err)
	assert.Equal(t, m.entries["core/assets"], []string{"TestLink"})
}
