package runner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

// timeoutPanic appears in go test -json output when the test binary's
// -timeout fires. It may be attached to a running test or to the package.
const timeoutPanic = "panic: test timed out"

// testEvent mirrors cmd/internal/test2json's TestEvent; only fields we need.
type testEvent struct {
	Action  string  `json:"Action"`
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Elapsed float64 `json:"Elapsed"`
	Output  string  `json:"Output"`
}

type testKey struct {
	Package string
	Test    string
}

type aggregate struct {
	passes       int
	fails        int
	maxElapsed   time.Duration
	timedOut     bool
	iterations   map[int]struct{}
	failedIters  map[int]bool
	timeoutIters map[int]bool
	outputs      map[int]*strings.Builder
}

// IterationLog captures a test's interleaved `output` events for one iteration.
type IterationLog struct {
	Iteration int    `json:"iteration"`
	Output    string `json:"output"`
}

// TestEntry is a single row in the analysis report.
type TestEntry struct {
	Package    string         `json:"package"`
	Test       string         `json:"test,omitempty"`
	Passes     int            `json:"passes,omitempty"`
	Fails      int            `json:"fails,omitempty"`
	MaxElapsed time.Duration  `json:"max_elapsed,omitempty"`
	Iterations []int          `json:"iterations,omitempty"`
	Logs       []IterationLog `json:"logs,omitempty"`
}

// Report classifies tests across iterations of a survey run.
type Report struct {
	Iterations    int           `json:"iterations"`
	SlowThreshold time.Duration `json:"slow_threshold"`
	Flakes        []TestEntry   `json:"flakes,omitempty"`
	Failures      []TestEntry   `json:"failures,omitempty"`
	Timeouts      []TestEntry   `json:"timeouts,omitempty"`
	Slow          []TestEntry   `json:"slow,omitempty"`
}

// Analyze reads per-iteration test2json streams and classifies tests.
// Malformed lines are silently skipped (go test can interleave non-JSON).
func Analyze(iterations []io.Reader, slowThreshold time.Duration) (*Report, error) {
	aggs := map[testKey]*aggregate{}
	newAgg := func() *aggregate {
		return &aggregate{
			iterations:   map[int]struct{}{},
			failedIters:  map[int]bool{},
			timeoutIters: map[int]bool{},
			outputs:      map[int]*strings.Builder{},
		}
	}

	for i, r := range iterations {
		// Line-based scan + per-line Unmarshal: go test -json can interleave
		// non-JSON output (stderr warnings, build errors); streaming decoder
		// can't recover from those. Skip unparsable lines silently.
		scanner := bufio.NewScanner(r)
		scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 || line[0] != '{' {
				continue
			}
			var ev testEvent
			if err := json.Unmarshal(line, &ev); err != nil {
				continue
			}
			key := testKey{Package: ev.Package, Test: ev.Test}
			a := aggs[key]
			if a == nil {
				a = newAgg()
				aggs[key] = a
			}
			switch ev.Action {
			case "pass":
				a.passes++
				a.iterations[i] = struct{}{}
				if d := seconds(ev.Elapsed); d > a.maxElapsed {
					a.maxElapsed = d
				}
			case "fail":
				a.fails++
				a.iterations[i] = struct{}{}
				a.failedIters[i] = true
				if d := seconds(ev.Elapsed); d > a.maxElapsed {
					a.maxElapsed = d
				}
			case "output":
				if strings.Contains(ev.Output, timeoutPanic) {
					a.timedOut = true
					a.iterations[i] = struct{}{}
					a.timeoutIters[i] = true
				}
				buf := a.outputs[i]
				if buf == nil {
					buf = &strings.Builder{}
					a.outputs[i] = buf
				}
				buf.WriteString(ev.Output)
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading iteration %d: %w", i, err)
		}
	}

	reattributeTimeouts(aggs, newAgg)

	rep := &Report{
		Iterations:    len(iterations),
		SlowThreshold: slowThreshold,
	}

	for key, a := range aggs {
		base := TestEntry{
			Package:    key.Package,
			Test:       key.Test,
			Passes:     a.passes,
			Fails:      a.fails,
			MaxElapsed: a.maxElapsed,
			Iterations: sortedKeys(a.iterations),
		}
		switch {
		case a.timedOut:
			entry := base
			entry.Logs = collectLogs(a, a.timeoutIters)
			rep.Timeouts = append(rep.Timeouts, entry)
		case key.Test == "":
			// package-level event without timeout: skip (suite-level pass/fail is noise)
		case a.passes > 0 && a.fails > 0:
			entry := base
			entry.Logs = collectLogs(a, a.failedIters)
			rep.Flakes = append(rep.Flakes, entry)
		case a.fails > 0 && a.passes == 0:
			entry := base
			entry.Logs = collectLogs(a, a.failedIters)
			rep.Failures = append(rep.Failures, entry)
		}
		if !a.timedOut && key.Test != "" && slowThreshold > 0 && a.maxElapsed > slowThreshold {
			rep.Slow = append(rep.Slow, base)
		}
	}

	sortEntries(rep.Flakes)
	sortEntries(rep.Failures)
	sortEntries(rep.Timeouts)
	sortEntries(rep.Slow)
	return rep, nil
}

// AnalyzeResults opens every `iteration-*.log.jsonl` file in resultsDir, in
// numeric-iteration order, and delegates to Analyze.
func AnalyzeResults(resultsDir string, slowThreshold time.Duration) (*Report, error) {
	matches, err := filepath.Glob(filepath.Join(resultsDir, "iteration-*.log.jsonl"))
	if err != nil {
		return nil, err
	}
	sort.Slice(matches, func(i, j int) bool {
		return iterNumber(matches[i]) < iterNumber(matches[j])
	})
	readers := make([]io.Reader, 0, len(matches))
	files := make([]*os.File, 0, len(matches))
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()
	for _, p := range matches {
		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
		readers = append(readers, f)
	}
	return Analyze(readers, slowThreshold)
}

// WriteReport writes the report as pretty JSON to <resultsDir>/report.json.
func WriteReport(resultsDir string, rep *Report) error {
	b, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(resultsDir, "report.json"), b, 0600)
}

// PrintSummary writes a short human summary.
func PrintSummary(w io.Writer, rep *Report) {
	fmt.Fprintf(w, "flakes (%d)%s\n", len(rep.Flakes), renderEntries(rep.Flakes, renderCounts))
	fmt.Fprintf(w, "failures (%d)%s\n", len(rep.Failures), renderEntries(rep.Failures, renderCounts))
	fmt.Fprintf(w, "timeouts (%d)%s\n", len(rep.Timeouts), renderEntries(rep.Timeouts, renderIterations))
	fmt.Fprintf(w, "slow >%s (%d)%s\n", rep.SlowThreshold, len(rep.Slow), renderEntries(rep.Slow, renderElapsed))
}

type renderMode int

const (
	renderCounts renderMode = iota
	renderElapsed
	renderIterations
)

func renderEntries(entries []TestEntry, mode renderMode) string {
	if len(entries) == 0 {
		return ""
	}
	parts := make([]string, 0, len(entries))
	for _, e := range entries {
		name := e.Package
		if e.Test != "" {
			name = e.Package + "." + e.Test
		}
		switch mode {
		case renderElapsed:
			parts = append(parts, fmt.Sprintf("%s %s", name, e.MaxElapsed.Round(time.Millisecond)))
		case renderIterations:
			if len(e.Iterations) == 0 {
				parts = append(parts, name)
			} else {
				parts = append(parts, fmt.Sprintf("%s (iter %s)", name, joinInts(e.Iterations)))
			}
		default: // renderCounts
			if e.Passes > 0 || e.Fails > 0 {
				parts = append(parts, fmt.Sprintf("%s (%dp/%df)", name, e.Passes, e.Fails))
			} else {
				parts = append(parts, name)
			}
		}
	}
	return ": " + strings.Join(parts, ", ")
}

func joinInts(ints []int) string {
	parts := make([]string, len(ints))
	for i, n := range ints {
		parts[i] = strconv.Itoa(n)
	}
	return strings.Join(parts, ",")
}

// reattributeTimeouts fixes the go-test-json quirk where a `panic: test timed out`
// is attached to whichever test most recently emitted events rather than the
// actually-stuck one. The real culprits are listed in the panic's
// "running tests:" block — move the timeout mark (and the captured stack
// trace) onto those tests.
func reattributeTimeouts(aggs map[testKey]*aggregate, newAgg func() *aggregate) {
	// Snapshot keys before we start mutating the map.
	keys := make([]testKey, 0, len(aggs))
	for k := range aggs {
		keys = append(keys, k)
	}
	for _, key := range keys {
		a := aggs[key]
		if !a.timedOut {
			continue
		}
		for i := range a.timeoutIters {
			buf := a.outputs[i]
			if buf == nil {
				continue
			}
			output := buf.String()
			names := parseRunningTests(output)
			if len(names) == 0 {
				continue
			}
			// If this entry is itself one of the running tests, it's correctly
			// attributed; leave it alone.
			if slices.Contains(names, key.Test) {
				continue
			}
			// Strip the bogus timeout attribution.
			delete(a.timeoutIters, i)
			if len(a.timeoutIters) == 0 {
				a.timedOut = false
			}
			// Attribute to the real culprit(s). Copy the output so each
			// culprit's Logs contain the full panic+stack for context.
			for _, name := range names {
				nk := testKey{Package: key.Package, Test: name}
				na := aggs[nk]
				if na == nil {
					na = newAgg()
					aggs[nk] = na
				}
				na.timedOut = true
				na.timeoutIters[i] = true
				na.iterations[i] = struct{}{}
				if na.outputs[i] == nil {
					na.outputs[i] = &strings.Builder{}
				}
				na.outputs[i].WriteString(output)
			}
		}
	}
}

// parseRunningTests extracts test names from a `panic: test timed out` block:
//
//	running tests:
//	        TestName (5s)
//	        TestOther/sub (4s)
func parseRunningTests(output string) []string {
	const marker = "running tests:"
	_, tail, found := strings.Cut(output, marker)
	if !found {
		return nil
	}
	var names []string
	for line := range strings.SplitSeq(tail, "\n") {
		trim := strings.TrimLeft(line, "\t ")
		if trim == "" {
			if len(names) == 0 {
				continue // skip leading blank line
			}
			break
		}
		// Expect "TestName (duration)" — split on the last " (".
		open := strings.LastIndex(trim, " (")
		if open < 0 || !strings.HasSuffix(trim, ")") {
			break
		}
		name := trim[:open]
		if name == "" {
			break
		}
		names = append(names, name)
	}
	return names
}

// collectLogs returns buffered output for each iteration where `want[i]` is true,
// sorted by iteration. Iterations with no buffered output are skipped.
func collectLogs(a *aggregate, want map[int]bool) []IterationLog {
	iters := make([]int, 0, len(want))
	for i := range want {
		iters = append(iters, i)
	}
	sort.Ints(iters)
	out := make([]IterationLog, 0, len(iters))
	for _, i := range iters {
		buf := a.outputs[i]
		if buf == nil || buf.Len() == 0 {
			continue
		}
		out = append(out, IterationLog{Iteration: i, Output: buf.String()})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func seconds(f float64) time.Duration {
	return time.Duration(f * float64(time.Second))
}

func sortedKeys(m map[int]struct{}) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Ints(out)
	return out
}

func sortEntries(entries []TestEntry) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Package != entries[j].Package {
			return entries[i].Package < entries[j].Package
		}
		return entries[i].Test < entries[j].Test
	})
}

func iterNumber(path string) int {
	base := filepath.Base(path)
	base = strings.TrimPrefix(base, "iteration-")
	base = strings.TrimSuffix(base, ".log.jsonl")
	n := 0
	for _, c := range base {
		if c < '0' || c > '9' {
			return -1
		}
		n = n*10 + int(c-'0')
	}
	return n
}
