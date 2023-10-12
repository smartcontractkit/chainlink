package flakeytests

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockReporter struct {
	entries map[string]map[string]struct{}
}

func (m *mockReporter) Report(entries map[string]map[string]struct{}) error {
	m.entries = entries
	return nil
}

func newMockReporter() *mockReporter {
	return &mockReporter{entries: map[string]map[string]struct{}{}}
}

func TestParser(t *testing.T) {
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}
`

	r := strings.NewReader(output)
	ts, err := parseOutput(r)
	require.NoError(t, err)

	assert.Len(t, ts, 1)
	assert.Len(t, ts["github.com/smartcontractkit/chainlink/v2/core/assets"], 1)
	assert.Equal(t, ts["github.com/smartcontractkit/chainlink/v2/core/assets"]["TestLink"], 1)
}

func TestParser_SkipsNonJSON(t *testing.T) {
	output := `Failed tests and panics:
-------
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}
`

	r := strings.NewReader(output)
	ts, err := parseOutput(r)
	require.NoError(t, err)

	assert.Len(t, ts, 1)
	assert.Len(t, ts["github.com/smartcontractkit/chainlink/v2/core/assets"], 1)
	assert.Equal(t, ts["github.com/smartcontractkit/chainlink/v2/core/assets"]["TestLink"], 1)
}

func TestParser_PanicDueToLogging(t *testing.T) {
	output := `
{"Time":"2023-09-07T16:01:40.649849+01:00","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestAssets_LinkScanValue","Output":"panic: foo\n"}
`

	r := strings.NewReader(output)
	ts, err := parseOutput(r)
	require.NoError(t, err)

	assert.Len(t, ts, 1)
	assert.Len(t, ts["github.com/smartcontractkit/chainlink/v2/core/assets"], 1)
	assert.Equal(t, ts["github.com/smartcontractkit/chainlink/v2/core/assets"]["TestAssets_LinkScanValue"], 1)
}

func TestParser_SuccessfulOutput(t *testing.T) {
	output := `
{"Time":"2023-09-07T16:22:52.556853+01:00","Action":"start","Package":"github.com/smartcontractkit/chainlink/v2/core/assets"}
{"Time":"2023-09-07T16:22:52.762353+01:00","Action":"run","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestAssets_NewLinkAndString"}
{"Time":"2023-09-07T16:22:52.762456+01:00","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestAssets_NewLinkAndString","Output":"=== RUN   TestAssets_NewLinkAndString\n"}
{"Time":"2023-09-07T16:22:52.76249+01:00","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestAssets_NewLinkAndString","Output":"=== PAUSE TestAssets_NewLinkAndString\n"}
{"Time":"2023-09-07T16:22:52.7625+01:00","Action":"pause","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestAssets_NewLinkAndString"}
{"Time":"2023-09-07T16:22:52.762511+01:00","Action":"cont","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestAssets_NewLinkAndString"}
{"Time":"2023-09-07T16:22:52.762528+01:00","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestAssets_NewLinkAndString","Output":"=== CONT  TestAssets_NewLinkAndString\n"}
{"Time":"2023-09-07T16:22:52.762546+01:00","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestAssets_NewLinkAndString","Output":"--- PASS: TestAssets_NewLinkAndString (0.00s)\n"}
{"Time":"2023-09-07T16:22:52.762557+01:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestAssets_NewLinkAndString","Elapsed":0}
{"Time":"2023-09-07T16:22:52.762566+01:00","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Output":"PASS\n"}
{"Time":"2023-09-07T16:22:52.762955+01:00","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Output":"ok  \tgithub.com/smartcontractkit/chainlink/v2/core/assets\t0.206s\n"}
{"Time":"2023-09-07T16:22:52.765598+01:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Elapsed":0.209}
`

	r := strings.NewReader(output)
	ts, err := parseOutput(r)
	require.NoError(t, err)
	assert.Len(t, ts, 0)
}

type testAdapter func(string, []string, io.Writer) error

func (t testAdapter) test(pkg string, tests []string, out io.Writer) error {
	return t(pkg, tests, out)
}

func TestRunner_WithFlake(t *testing.T) {
	initialOutput := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}`
	outputs := []string{
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}`,
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
	err := r.Run()
	require.NoError(t, err)
	assert.Len(t, m.entries, 1)
	_, ok := m.entries["github.com/smartcontractkit/chainlink/v2/core/assets"]["TestLink"]
	assert.True(t, ok)
}

func TestRunner_WithFailedPackage(t *testing.T) {
	initialOutput := `
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Elapsed":0}
`
	outputs := []string{`
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Elapsed":0}
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
	err := r.Run()
	require.NoError(t, err)
	assert.Len(t, m.entries, 1)
	_, ok := m.entries["github.com/smartcontractkit/chainlink/v2/core/assets"]["TestLink"]
	assert.True(t, ok)
}

func TestRunner_AllFailures(t *testing.T) {
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}`

	rerunOutput := `
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}
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

	err := r.Run()
	require.NoError(t, err)
	assert.Len(t, m.entries, 0)
}

func TestRunner_RerunSuccessful(t *testing.T) {
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}`

	rerunOutputs := []string{
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}`,
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}`,
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

	err := r.Run()
	require.NoError(t, err)
	_, ok := m.entries["github.com/smartcontractkit/chainlink/v2/core/assets"]["TestLink"]
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

	err := r.Run()
	require.NoError(t, err)
	_, ok := m.entries["github.com/smartcontractkit/chainlink/v2/"]["TestConfigDocs"]
	assert.True(t, ok)
}

func TestRunner_RerunFailsWithNonzeroExitCode(t *testing.T) {
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}`

	rerunOutputs := []string{
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}`,
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}`,
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

	err := r.Run()
	require.NoError(t, err)
	_, ok := m.entries["github.com/smartcontractkit/chainlink/v2/core/assets"]["TestLink"]
	assert.True(t, ok)
}

func TestRunner_RerunWithNonZeroExitCodeDoesntStopCommand(t *testing.T) {
	outputs := []io.Reader{
		strings.NewReader(`
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}
`),
		strings.NewReader(`
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2","Test":"TestMaybeReservedLinkV2","Elapsed":0}
`),
	}

	rerunOutputs := []string{
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}`,
		`{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}`,
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

	err := r.Run()
	require.NoError(t, err)
	calls := index
	assert.Equal(t, 4, calls)

	_, ok := m.entries["github.com/smartcontractkit/chainlink/v2/core/assets"]["TestLink"]
	assert.True(t, ok)
}

// Used for integration tests
func TestSkippedForTests(t *testing.T) {
	if os.Getenv("FLAKEY_TEST_RUNNER_RUN_FIXTURE_TEST") != "1" {
		t.Skip()
	}

	go func() {
		panic("skipped test")
	}()
}

// Used for integration tests
func TestSkippedForTests_Success(t *testing.T) {
	if os.Getenv("FLAKEY_TEST_RUNNER_RUN_FIXTURE_TEST") != "1" {
		t.Skip()
	}

	assert.True(t, true)
}

func TestParsesPanicCorrectly(t *testing.T) {
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/tools/flakeytests/","Test":"TestSkippedForTests","Elapsed":0}`

	m := newMockReporter()
	tc := &testCommand{
		repo:    "github.com/smartcontractkit/chainlink/v2/tools/flakeytests",
		command: "../bin/go_core_tests",
		overrides: func(cmd *exec.Cmd) {
			cmd.Env = append(cmd.Env, "FLAKEY_TESTRUNNER_RUN_FIXTURE_TEST=1")
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

	err := r.Run()
	require.NoError(t, err)
	_, ok := m.entries["github.com/smartcontractkit/chainlink/v2/tools/flakeytests"]["TestSkippedForTests"]
	assert.False(t, ok)
}

func TestIntegration(t *testing.T) {
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/tools/flakeytests/","Test":"TestSkippedForTests_Success","Elapsed":0}`

	m := newMockReporter()
	tc := &testCommand{
		repo:    "github.com/smartcontractkit/chainlink/v2/tools/flakeytests",
		command: "../bin/go_core_tests",
		overrides: func(cmd *exec.Cmd) {
			cmd.Env = append(cmd.Env, "FLAKEY_TESTRUNNER_RUN_FIXTURE_TEST=1")
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

	err := r.Run()
	require.NoError(t, err)
	_, ok := m.entries["github.com/smartcontractkit/chainlink/v2/tools/flakeytests"]["TestSkippedForTests_Success"]
	assert.False(t, ok)
}
