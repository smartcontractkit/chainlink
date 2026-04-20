package runner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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
	passes     int
	fails      int
	maxElapsed time.Duration
	timedOut   bool
	iterations map[int]struct{}
}

// TestEntry is a single row in the analysis report.
type TestEntry struct {
	Package    string        `json:"package"`
	Test       string        `json:"test,omitempty"`
	Passes     int           `json:"passes,omitempty"`
	Fails      int           `json:"fails,omitempty"`
	MaxElapsed time.Duration `json:"max_elapsed,omitempty"`
	Iterations []int         `json:"iterations,omitempty"`
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
				a = &aggregate{iterations: map[int]struct{}{}}
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
				if d := seconds(ev.Elapsed); d > a.maxElapsed {
					a.maxElapsed = d
				}
			case "output":
				if strings.Contains(ev.Output, timeoutPanic) {
					a.timedOut = true
					a.iterations[i] = struct{}{}
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading iteration %d: %w", i, err)
		}
	}

	rep := &Report{
		Iterations:    len(iterations),
		SlowThreshold: slowThreshold,
	}

	for key, a := range aggs {
		entry := TestEntry{
			Package:    key.Package,
			Test:       key.Test,
			Passes:     a.passes,
			Fails:      a.fails,
			MaxElapsed: a.maxElapsed,
			Iterations: sortedKeys(a.iterations),
		}
		switch {
		case a.timedOut:
			rep.Timeouts = append(rep.Timeouts, entry)
		case key.Test == "":
			// package-level event without timeout: skip (suite-level pass/fail is noise)
		case a.passes > 0 && a.fails > 0:
			rep.Flakes = append(rep.Flakes, entry)
		case a.fails > 0 && a.passes == 0:
			rep.Failures = append(rep.Failures, entry)
		}
		if !a.timedOut && key.Test != "" && slowThreshold > 0 && a.maxElapsed > slowThreshold {
			rep.Slow = append(rep.Slow, entry)
		}
	}

	sortEntries(rep.Flakes)
	sortEntries(rep.Failures)
	sortEntries(rep.Timeouts)
	sortEntries(rep.Slow)
	return rep, nil
}

// AnalyzeResults opens every `test-*.log.jsonl` file in resultsDir, in
// numeric-iteration order, and delegates to Analyze.
func AnalyzeResults(resultsDir string, slowThreshold time.Duration) (*Report, error) {
	matches, err := filepath.Glob(filepath.Join(resultsDir, "test-*.log.jsonl"))
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
	fmt.Fprintf(w, "flakes (%d)%s\n", len(rep.Flakes), renderEntries(rep.Flakes, false))
	fmt.Fprintf(w, "failures (%d)%s\n", len(rep.Failures), renderEntries(rep.Failures, false))
	fmt.Fprintf(w, "timeouts (%d)%s\n", len(rep.Timeouts), renderEntries(rep.Timeouts, false))
	fmt.Fprintf(w, "slow >%s (%d)%s\n", rep.SlowThreshold, len(rep.Slow), renderEntries(rep.Slow, true))
}

func renderEntries(entries []TestEntry, withElapsed bool) string {
	if len(entries) == 0 {
		return ""
	}
	parts := make([]string, 0, len(entries))
	for _, e := range entries {
		name := e.Package
		if e.Test != "" {
			name = e.Package + "." + e.Test
		}
		switch {
		case withElapsed:
			parts = append(parts, fmt.Sprintf("%s %s", name, e.MaxElapsed.Round(time.Millisecond)))
		case e.Passes > 0 || e.Fails > 0:
			parts = append(parts, fmt.Sprintf("%s (%dp/%df)", name, e.Passes, e.Fails))
		default:
			parts = append(parts, name)
		}
	}
	return ": " + strings.Join(parts, ", ")
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
	base = strings.TrimPrefix(base, "test-")
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
