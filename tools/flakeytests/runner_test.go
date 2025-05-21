package flakeytests

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test suite has two sets of tests that are run in two different modes:
//  1. Normal mode: The tests are run as usual. All tests that start with TestSkippedForTests will be skipped.
//     There are a set of integration tests that will run this same test suite in integration mode.
//  2. Integration mode: The tests are run with an environment variable set to run the tests that were skipped in normal mode.
//     Their output is captured and parsed to check if the flake runner is behaving as expected.
//
// In summary, the first set of tests are executed by the user, and these tests will trigger the second set of tests to run.
var runInIntegrationTestModeEnvKey = "FLAKEY_TEST_RUNNER_RUN_FIXTURE_TEST"
var skipIntegrationTestMode = os.Getenv(runInIntegrationTestModeEnvKey) != "1"
var runInIntegrationTestModeEnv = runInIntegrationTestModeEnvKey + "=1"

type mockReporter struct {
	report *Report
}

func (m *mockReporter) Report(_ context.Context, report *Report) error {
	m.report = report
	return nil
}

func newMockReporter() *mockReporter {
	return &mockReporter{report: NewReport()}
}

func TestParser(t *testing.T) {
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}
`

	r := strings.NewReader(output)
	pr, err := parseOutput(r)
	require.NoError(t, err)

	ts := pr.tests
	assert.Len(t, ts, 1)
	assert.Len(t, ts["github.com/smartcontractkit/chainlink-evm/pkg/assets"], 1)
	assert.Equal(t, 1, ts["github.com/smartcontractkit/chainlink-evm/pkg/assets"]["TestLink"])
}

func TestParser_SkipsNonJSON(t *testing.T) {
	output := `Failed tests and panics:
-------
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}
`

	r := strings.NewReader(output)
	pr, err := parseOutput(r)
	require.NoError(t, err)

	ts := pr.tests
	assert.Len(t, ts, 1)
	assert.Len(t, ts["github.com/smartcontractkit/chainlink-evm/pkg/assets"], 1)
	assert.Equal(t, 1, ts["github.com/smartcontractkit/chainlink-evm/pkg/assets"]["TestLink"])
}

func TestParser_PanicDueToLogging(t *testing.T) {
	output := `
{"Time":"2023-09-07T16:01:40.649849+01:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestAssets_LinkScanValue","Output":"panic: foo\n"}
`

	r := strings.NewReader(output)
	pr, err := parseOutput(r)
	require.NoError(t, err)

	ts := pr.tests
	assert.Len(t, ts, 1)
	assert.Len(t, ts["github.com/smartcontractkit/chainlink-evm/pkg/assets"], 1)
	assert.Equal(t, 1, ts["github.com/smartcontractkit/chainlink-evm/pkg/assets"]["TestAssets_LinkScanValue"])
}

func TestParser_SuccessfulOutput(t *testing.T) {
	output := `
{"Time":"2023-09-07T16:22:52.556853+01:00","Action":"start","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets"}
{"Time":"2023-09-07T16:22:52.762353+01:00","Action":"run","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestAssets_NewLinkAndString"}
{"Time":"2023-09-07T16:22:52.762456+01:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestAssets_NewLinkAndString","Output":"=== RUN   TestAssets_NewLinkAndString\n"}
{"Time":"2023-09-07T16:22:52.76249+01:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestAssets_NewLinkAndString","Output":"=== PAUSE TestAssets_NewLinkAndString\n"}
{"Time":"2023-09-07T16:22:52.7625+01:00","Action":"pause","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestAssets_NewLinkAndString"}
{"Time":"2023-09-07T16:22:52.762511+01:00","Action":"cont","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestAssets_NewLinkAndString"}
{"Time":"2023-09-07T16:22:52.762528+01:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestAssets_NewLinkAndString","Output":"=== CONT  TestAssets_NewLinkAndString\n"}
{"Time":"2023-09-07T16:22:52.762546+01:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestAssets_NewLinkAndString","Output":"--- PASS: TestAssets_NewLinkAndString (0.00s)\n"}
{"Time":"2023-09-07T16:22:52.762557+01:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestAssets_NewLinkAndString","Elapsed":0}
{"Time":"2023-09-07T16:22:52.762566+01:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Output":"PASS\n"}
{"Time":"2023-09-07T16:22:52.762955+01:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Output":"ok  \tgithub.com/smartcontractkit/chainlink/v2/core/assets\t0.206s\n"}
{"Time":"2023-09-07T16:22:52.765598+01:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Elapsed":0.209}
`

	r := strings.NewReader(output)
	ts, err := parseOutput(r)
	require.NoError(t, err)
	assert.Empty(t, ts.tests)
}

type testAdapter func(string, []string, io.Writer) error

func (t testAdapter) test(pkg string, tests []string, out io.Writer) error {
	return t(pkg, tests, out)
}

func TestRunner_WithFlake(t *testing.T) {
	initialOutput := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}`
	outputs := []string{
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}`,
		``,
	}
	m := newMockReporter()
	i := 0
	r := &Runner{
		numReruns: 2,
		readers:   []io.Reader{strings.NewReader(initialOutput)},

		testCommand: testAdapter(func(pkg string, testNames []string, w io.Writer) error {
			_, err := w.Write([]byte(outputs[i]))
			i++
			return err
		}),
		parse:    parseOutput,
		reporter: m,
	}

	// This will report a flake since we've mocked the rerun
	// to only report one failure (not two as expected).
	err := r.Run(t.Context())
	require.NoError(t, err)
	assert.Len(t, m.report.tests, 1)
	_, ok := m.report.tests["github.com/smartcontractkit/chainlink-evm/pkg/assets"]["TestLink"]
	assert.True(t, ok)
}

func TestRunner_WithFailedPackage(t *testing.T) {
	initialOutput := `
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Elapsed":0}
`
	outputs := []string{`
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Elapsed":0}
`,
		``,
	}

	m := newMockReporter()
	i := 0
	r := &Runner{
		numReruns: 2,
		readers:   []io.Reader{strings.NewReader(initialOutput)},
		testCommand: testAdapter(func(pkg string, testNames []string, w io.Writer) error {
			_, err := w.Write([]byte(outputs[i]))
			i++
			return err
		}),
		parse:    parseOutput,
		reporter: m,
	}

	// This will report a flake since we've mocked the rerun
	// to only report one failure (not two as expected).
	err := r.Run(t.Context())
	require.NoError(t, err)
	assert.Len(t, m.report.tests, 1)
	_, ok := m.report.tests["github.com/smartcontractkit/chainlink-evm/pkg/assets"]["TestLink"]
	assert.True(t, ok)
}

func TestRunner_AllFailures(t *testing.T) {
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}`

	rerunOutput := `
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}
`
	m := newMockReporter()
	r := &Runner{
		numReruns: 2,
		readers:   []io.Reader{strings.NewReader(output)},
		testCommand: testAdapter(func(pkg string, testNames []string, w io.Writer) error {
			_, err := w.Write([]byte(rerunOutput))
			return err
		}),
		parse:    parseOutput,
		reporter: m,
	}

	err := r.Run(t.Context())
	require.NoError(t, err)
	assert.Empty(t, m.report.tests)
}

func TestRunner_RerunSuccessful(t *testing.T) {
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}`

	rerunOutputs := []string{
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}`,
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}`,
	}
	m := newMockReporter()
	i := 0
	r := &Runner{
		numReruns: 2,
		readers:   []io.Reader{strings.NewReader(output)},
		testCommand: testAdapter(func(pkg string, testNames []string, w io.Writer) error {
			_, err := w.Write([]byte(rerunOutputs[i]))
			i++
			return err
		}),
		parse:    parseOutput,
		reporter: m,
	}

	err := r.Run(t.Context())
	require.NoError(t, err)
	_, ok := m.report.tests["github.com/smartcontractkit/chainlink-evm/pkg/assets"]["TestLink"]
	assert.True(t, ok)
}

func TestRunner_RootLevelTest(t *testing.T) {
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/","Test":"TestConfigDocs","Elapsed":0}`

	rerunOutput := ``
	m := newMockReporter()
	r := &Runner{
		numReruns: 2,
		readers:   []io.Reader{strings.NewReader(output)},
		testCommand: testAdapter(func(pkg string, testNames []string, w io.Writer) error {
			_, err := w.Write([]byte(rerunOutput))
			return err
		}),
		parse:    parseOutput,
		reporter: m,
	}

	err := r.Run(t.Context())
	require.NoError(t, err)
	_, ok := m.report.tests["github.com/smartcontractkit/chainlink/v2/"]["TestConfigDocs"]
	assert.True(t, ok)
}

func TestRunner_RerunFailsWithNonzeroExitCode(t *testing.T) {
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}`

	rerunOutputs := []string{
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}`,
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}`,
	}
	m := newMockReporter()
	i := 0
	r := &Runner{
		numReruns: 2,
		readers:   []io.Reader{strings.NewReader(output)},
		testCommand: testAdapter(func(pkg string, testNames []string, w io.Writer) error {
			_, err := w.Write([]byte(rerunOutputs[i]))
			i++
			return err
		}),
		parse:    parseOutput,
		reporter: m,
	}

	err := r.Run(t.Context())
	require.NoError(t, err)
	_, ok := m.report.tests["github.com/smartcontractkit/chainlink-evm/pkg/assets"]["TestLink"]
	assert.True(t, ok)
}

func TestRunner_RerunWithNonZeroExitCodeDoesntStopCommand(t *testing.T) {
	outputs := []io.Reader{
		strings.NewReader(`
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}
`),
		strings.NewReader(`
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2","Test":"TestMaybeReservedLinkV2","Elapsed":0}
`),
	}

	rerunOutputs := []string{
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}`,
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-evm/pkg/assets","Test":"TestLink","Elapsed":0}`,
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2","Test":"TestMaybeReservedLinkV2","Elapsed":0}`,
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2","Test":"TestMaybeReservedLinkV2","Elapsed":0}`,
	}

	index := 0
	m := newMockReporter()
	r := &Runner{
		numReruns: 2,
		readers:   outputs,
		testCommand: testAdapter(func(pkg string, testNames []string, w io.Writer) error {
			_, err := w.Write([]byte(rerunOutputs[index]))
			index++
			return err
		}),
		parse:    parseOutput,
		reporter: m,
	}

	err := r.Run(t.Context())
	require.NoError(t, err)
	calls := index
	assert.Equal(t, 4, calls)

	_, ok := m.report.tests["github.com/smartcontractkit/chainlink-evm/pkg/assets"]["TestLink"]
	assert.True(t, ok)
}

// Used for integration tests
func TestSkippedForTests_Subtests(t *testing.T) {
	if skipIntegrationTestMode {
		t.Skip()
	}

	t.Run("1: should fail", func(t *testing.T) {
		assert.False(t, true)
	})

	t.Run("2: should fail", func(t *testing.T) {
		assert.False(t, true)
	})
}

// Used for integration tests
func TestSkippedForTests(t *testing.T) {
	if skipIntegrationTestMode {
		t.Skip()
	}

	go func() {
		panic("skipped test")
	}()
}

// Used for integration tests
func TestSkippedForTests_Success(t *testing.T) {
	if skipIntegrationTestMode {
		t.Skip()
	}

	assert.True(t, true)
}

func TestIntegration_DealsWithSubtests(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	output := `
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/tools/flakeytests/","Test":"TestSkippedForTests_Subtests/1:_should_fail","Elapsed":0}
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/tools/flakeytests/","Test":"TestSkippedForTests_Subtests","Elapsed":0}
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/tools/flakeytests/","Test":"TestSkippedForTests_Subtests/2:_should_fail","Elapsed":0}
`

	m := newMockReporter()
	tc := &testCommand{
		repo:    "github.com/smartcontractkit/chainlink/v2/tools/flakeytests",
		command: "../bin/go_core_tests",
		overrides: func(cmd *exec.Cmd) {
			cmd.Env = append(cmd.Env, runInIntegrationTestModeEnv)
			cmd.Stdout = io.Discard
			cmd.Stderr = io.Discard
		},
	}
	r := &Runner{
		numReruns:   2,
		readers:     []io.Reader{strings.NewReader(output)},
		testCommand: tc,
		parse:       parseOutput,
		reporter:    m,
	}

	err := r.Run(t.Context())
	require.NoError(t, err)
	expectedTests := map[string]map[string]int{
		"github.com/smartcontractkit/chainlink/v2/tools/flakeytests/": {
			"TestSkippedForTests_Subtests/1:_should_fail": 1,
			"TestSkippedForTests_Subtests/2:_should_fail": 1,
		},
	}
	assert.Equal(t, expectedTests, m.report.tests)
}

func TestIntegration_ParsesPanics(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/tools/flakeytests/","Test":"TestSkippedForTests","Elapsed":0}`

	m := newMockReporter()
	tc := &testCommand{
		repo:    "github.com/smartcontractkit/chainlink/v2/tools/flakeytests",
		command: "../bin/go_core_tests",
		overrides: func(cmd *exec.Cmd) {
			cmd.Env = append(cmd.Env, runInIntegrationTestModeEnv)
			cmd.Stdout = io.Discard
			cmd.Stderr = io.Discard
		},
	}
	r := &Runner{
		numReruns:   2,
		readers:     []io.Reader{strings.NewReader(output)},
		testCommand: tc,
		parse:       parseOutput,
		reporter:    m,
	}

	err := r.Run(t.Context())
	require.NoError(t, err)
	_, ok := m.report.tests["github.com/smartcontractkit/chainlink/v2/tools/flakeytests"]["TestSkippedForTests"]
	assert.False(t, ok)
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/tools/flakeytests/","Test":"TestSkippedForTests_Success","Elapsed":0}`

	m := newMockReporter()
	tc := &testCommand{
		repo:    "github.com/smartcontractkit/chainlink/v2/tools/flakeytests",
		command: "../bin/go_core_tests",
		overrides: func(cmd *exec.Cmd) {
			cmd.Env = append(cmd.Env, runInIntegrationTestModeEnv)
			cmd.Stdout = io.Discard
			cmd.Stderr = io.Discard
		},
	}
	r := &Runner{
		numReruns:   2,
		readers:     []io.Reader{strings.NewReader(output)},
		testCommand: tc,
		parse:       parseOutput,
		reporter:    m,
	}

	err := r.Run(t.Context())
	require.NoError(t, err)
	_, ok := m.report.tests["github.com/smartcontractkit/chainlink/v2/tools/flakeytests"]["TestSkippedForTests_Success"]
	assert.False(t, ok)
}
