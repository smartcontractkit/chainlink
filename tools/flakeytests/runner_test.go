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

func TestParser_PanicInTest(t *testing.T) {
	output := `
?   	github.com/smartcontractkit/chainlink/v2/tools/flakeytests/cmd/runner	[no test files]
--- FAIL: TestParser (0.00s)
panic: foo [recovered]
	panic: foo

goroutine 21 [running]:
testing.tRunner.func1.2({0x1009953c0, 0x1009d1e40})
	/opt/homebrew/Cellar/go/1.20.3/libexec/src/testing/testing.go:1526 +0x1c8
testing.tRunner.func1()
	/opt/homebrew/Cellar/go/1.20.3/libexec/src/testing/testing.go:1529 +0x384
panic({0x1009953c0, 0x1009d1e40})
	/opt/homebrew/Cellar/go/1.20.3/libexec/src/runtime/panic.go:884 +0x204
github.com/smartcontractkit/chainlink/v2/tools/flakeytests.TestParser(0x0?)
	/Users/ccordenier/Development/chainlink/tools/flakeytests/runner_test.go:50 +0xa4
testing.tRunner(0x14000083520, 0x1009d1588)
	/opt/homebrew/Cellar/go/1.20.3/libexec/src/testing/testing.go:1576 +0x10c
created by testing.(*T).Run
	/opt/homebrew/Cellar/go/1.20.3/libexec/src/testing/testing.go:1629 +0x368
FAIL	github.com/smartcontractkit/chainlink/v2/tools/flakeytests	0.197s
FAIL`

	r := strings.NewReader(output)
	ts, err := parseOutput(r)
	require.NoError(t, err)

	assert.Len(t, ts, 1)
	assert.Len(t, ts["tools/flakeytests"], 1)
	assert.Equal(t, ts["tools/flakeytests"]["TestParser"], 1)
}

func TestParser_PanicDueToLogging(t *testing.T) {
	output := `
panic: Log in goroutine after TestIntegration_LogEventProvider_Backfill has completed: 2023-07-19T10:10:45.925Z	WARN	KeepersRegistry.LogEventProvider	logprovider/provider.go:218	failed to read logs	{"version": "2.3.0@d898528", "where": "reader", "err": "fetched logs with errors: context canceled"}

goroutine 4999 [running]:
testing.(*common).logDepth(0xc0051f6000, {0xc003011960, 0xd3}, 0x3)
	/opt/hostedtoolcache/go/1.20.5/x64/src/testing/testing.go:1003 +0x4e7
testing.(*common).log(...)
	/opt/hostedtoolcache/go/1.20.5/x64/src/testing/testing.go:985
testing.(*common).Logf(0xc0051f6000, {0x21ba777?, 0x41ac8a?}, {0xc00217c330?, 0x1e530c0?, 0x1?})
	/opt/hostedtoolcache/go/1.20.5/x64/src/testing/testing.go:1036 +0x5a
go.uber.org/zap/zaptest.testingWriter.Write({{0x7f4c5c94f018?, 0xc0051f6000?}, 0xa8?}, {0xc003017000?, 0xd4, 0xc00217c320?})
	/home/runner/go/pkg/mod/go.uber.org/zap@v1.24.0/zaptest/logger.go:130 +0xe6
go.uber.org/zap/zapcore.(*ioCore).Write(0xc0022f60f0, {0x1, {0xc1260b897723e1ca, 0x4bdc75e54, 0x3d56e40}, {0x22e96a1, 0x20}, {0x22c3204, 0x13}, {0x1, ...}, ...}, ...)
	/home/runner/go/pkg/mod/go.uber.org/zap@v1.24.0/zapcore/core.go:99 +0xb5
go.uber.org/zap/zapcore.(*CheckedEntry).Write(0xc001265ba0, {0xc0023ca100, 0x1, 0x2})
	/home/runner/go/pkg/mod/go.uber.org/zap@v1.24.0/zapcore/entry.go:255 +0x1d9
go.uber.org/zap.(*SugaredLogger).log(0xc00363a008, 0x1, {0x22c3204?, 0x13?}, {0x0?, 0x0?, 0x0?}, {0xc004ea3f80, 0x2, 0x2})
	/home/runner/go/pkg/mod/go.uber.org/zap@v1.24.0/sugar.go:295 +0xee
go.uber.org/zap.(*SugaredLogger).Warnw(...)
	/home/runner/go/pkg/mod/go.uber.org/zap@v1.24.0/sugar.go:216
github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider.(*logEventProvider).startReader(0xc00076d730, {0x2917018?, 0xc0043c4000?}, 0xc003b2a000)
	/home/runner/work/chainlink/chainlink/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider/provider.go:218 +0x29f
created by github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider.(*logEventProvider).Start
	/home/runner/work/chainlink/chainlink/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider/provider.go:108 +0x133
FAIL	github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider	20.380s
FAIL
`

	r := strings.NewReader(output)
	ts, err := parseOutput(r)
	require.NoError(t, err)

	assert.Len(t, ts, 1)
	assert.Len(t, ts["core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"], 1)
	assert.Equal(t, ts["core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"]["TestIntegration_LogEventProvider_Backfill"], 1)
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
