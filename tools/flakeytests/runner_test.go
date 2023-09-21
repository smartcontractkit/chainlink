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

func TestRunner_WithFlake(t *testing.T) {
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}`
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
	assert.Equal(t, m.entries["github.com/smartcontractkit/chainlink/v2/core/assets"], []string{"TestLink"})
}

func TestRunner_WithFailedPackage(t *testing.T) {
	output := `
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Elapsed":0}
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
	assert.Equal(t, m.entries["github.com/smartcontractkit/chainlink/v2/core/assets"], []string{"TestLink"})
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
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}`

	rerunOutput := `
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}
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
	assert.Equal(t, m.entries["github.com/smartcontractkit/chainlink/v2/core/assets"], []string{"TestLink"})
}

func TestRunner_RootLevelTest(t *testing.T) {
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/","Test":"TestConfigDocs","Elapsed":0}`

	rerunOutput := ``
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
	assert.Equal(t, m.entries["github.com/smartcontractkit/chainlink/v2/"], []string{"TestConfigDocs"})
}

type exitError struct{}

func (e *exitError) ExitCode() int { return 1 }

func (e *exitError) Error() string { return "exit code: 1" }

func TestRunner_RerunFailsWithNonzeroExitCode(t *testing.T) {
	output := `{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}`

	rerunOutput := `
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}
`
	m := newMockReporter()
	r := &Runner{
		numReruns: 2,
		readers:   []io.Reader{strings.NewReader(output)},
		runTestFn: func(pkg string, testNames []string, numReruns int, w io.Writer) error {
			_, _ = w.Write([]byte(rerunOutput))
			return &exitError{}
		},
		parse:    parseOutput,
		reporter: m,
	}

	err := r.Run()
	require.NoError(t, err)
	assert.Equal(t, m.entries["github.com/smartcontractkit/chainlink/v2/core/assets"], []string{"TestLink"})
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
		`
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink/v2/core/assets","Test":"TestLink","Elapsed":0}
`,
		`
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2","Test":"TestMaybeReservedLinkV2","Elapsed":0}
{"Time":"2023-09-07T15:39:46.378315+01:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2","Test":"TestMaybeReservedLinkV2","Elapsed":0}
`,
	}

	index := 0
	m := newMockReporter()
	r := &Runner{
		numReruns: 2,
		readers:   outputs,
		runTestFn: func(pkg string, testNames []string, numReruns int, w io.Writer) error {

			_, _ = w.Write([]byte(rerunOutputs[index]))
			index++
			return &exitError{}
		},
		parse:    parseOutput,
		reporter: m,
	}

	err := r.Run()
	require.NoError(t, err)
	calls := index
	assert.Equal(t, 2, calls)
	assert.Equal(t, m.entries["github.com/smartcontractkit/chainlink/v2/core/assets"], []string{"TestLink"})
}
