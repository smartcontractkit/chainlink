package runner

import (
	"bufio"
	"encoding/csv"
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

	"charm.land/lipgloss/v2"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/termstyle"
)

// timeoutPanic appears in go test -json output when the test binary's
// -timeout fires. It may be attached to a running test or to the package.
const timeoutPanic = "panic: test timed out"

// TestEvent mirrors cmd/internal/test2json's TestEvent; only fields we need.
type TestEvent struct {
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
	skips        int
	maxElapsed   time.Duration
	timedOut     bool
	iterations   map[int]struct{}
	failedIters  map[int]bool
	timeoutIters map[int]bool
	skipIters    map[int]bool
	outputs      map[int]*strings.Builder
	elapseds     []time.Duration
}

// TestEntry is a single row in the analysis report.
type TestEntry struct {
	Package    string        `json:"package"`
	Test       string        `json:"test,omitempty"`
	Runs       int           `json:"runs"`
	Successes  int           `json:"successes"`
	Fails      int           `json:"fails"`
	Skips      int           `json:"skips"`
	Timeouts   int           `json:"timeouts"`
	MinElapsed time.Duration `json:"min_elapsed"`
	MaxElapsed time.Duration `json:"max_elapsed"`
	P50Elapsed time.Duration `json:"p50_elapsed"`
	LogFiles   []string      `json:"log_files,omitempty"`
}

// IterationSummary captures high-level stats for a single diagnose iteration.
// Duration and ShuffleSeed are populated by the runner after analysis.
type IterationSummary struct {
	Index        int           `json:"index"`
	Duration     time.Duration `json:"duration,omitempty"`
	Result       string        `json:"result"` // "pass", "fail", "timeout"
	FailingTests []string      `json:"failing_tests,omitempty"`
	ShuffleSeed  int64         `json:"shuffle_seed,omitempty"`
}

// Report classifies tests across iterations of a diagnose run.
type Report struct {
	Iterations         int                `json:"iterations"`
	SlowThreshold      time.Duration      `json:"slow_threshold"`
	IterationSummaries []IterationSummary `json:"iteration_summaries,omitempty"`
	Flakes             []TestEntry        `json:"flakes,omitempty"`
	Failures           []TestEntry        `json:"failures,omitempty"`
	Timeouts           []TestEntry        `json:"timeouts,omitempty"`
	Slow               []TestEntry        `json:"slow,omitempty"`
}

// LogMap maps (package,test) → iteration → raw interleaved output.
// Returned alongside Report so callers can write per-test log files without
// coupling the parser to the filesystem.
type LogMap map[testKey]map[int]string

// Analyze reads per-iteration test2json streams and classifies tests.
// Malformed lines are silently skipped (go test can interleave non-JSON).
func Analyze(iterations []io.Reader, slowThreshold time.Duration) (*Report, LogMap, error) {
	aggs := map[testKey]*aggregate{}
	newAgg := func() *aggregate {
		return &aggregate{
			iterations:   map[int]struct{}{},
			failedIters:  map[int]bool{},
			timeoutIters: map[int]bool{},
			skipIters:    map[int]bool{},
			outputs:      map[int]*strings.Builder{},
		}
	}

	for i, r := range iterations {
		// Line-based scan + per-line Unmarshal: go test -json can interleave
		// non-JSON output (stderr warnings, build errors); streaming decoder
		// can't recover from those. Skip unparsable lines silently.
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 || line[0] != '{' {
				continue
			}
			var ev TestEvent
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
				d := seconds(ev.Elapsed)
				a.elapseds = append(a.elapseds, d)
				if d > a.maxElapsed {
					a.maxElapsed = d
				}
			case "fail":
				a.fails++
				a.iterations[i] = struct{}{}
				a.failedIters[i] = true
				d := seconds(ev.Elapsed)
				a.elapseds = append(a.elapseds, d)
				if d > a.maxElapsed {
					a.maxElapsed = d
				}
			case "skip":
				a.skips++
				a.iterations[i] = struct{}{}
				a.skipIters[i] = true
				d := seconds(ev.Elapsed)
				a.elapseds = append(a.elapseds, d)
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
			return nil, nil, fmt.Errorf("reading iteration %d: %w", i, err)
		}
	}

	reattributeTimeouts(aggs, newAgg)

	rep := &Report{
		Iterations:    len(iterations),
		SlowThreshold: slowThreshold,
	}

	for key, a := range aggs {
		minE, p50 := stats(a.elapseds)
		base := TestEntry{
			Package:    key.Package,
			Test:       key.Test,
			Runs:       len(a.iterations),
			Successes:  a.passes,
			Fails:      a.fails,
			Skips:      a.skips,
			Timeouts:   len(a.timeoutIters),
			MinElapsed: minE,
			MaxElapsed: a.maxElapsed,
			P50Elapsed: p50,
		}
		switch {
		case a.timedOut:
			rep.Timeouts = append(rep.Timeouts, base)
		case key.Test == "" && a.fails == 0:
			// Package-level pass summary or benign events (no failing tests).
		case key.Test == "" && a.fails > 0:
			// Build failures, TestMain/init failures: Test is empty in go test -json.
			if a.passes > 0 {
				rep.Flakes = append(rep.Flakes, base)
			} else {
				rep.Failures = append(rep.Failures, base)
			}
		case a.passes > 0 && a.fails > 0:
			rep.Flakes = append(rep.Flakes, base)
		case a.fails > 0 && a.passes == 0:
			rep.Failures = append(rep.Failures, base)
		}
		if !a.timedOut && key.Test != "" && slowThreshold > 0 && a.maxElapsed > slowThreshold {
			rep.Slow = append(rep.Slow, base)
		}
	}

	sortEntries(rep.Flakes)
	sortEntries(rep.Failures)
	sortEntries(rep.Timeouts)
	sortEntries(rep.Slow)

	// Build per-iteration summaries from aggregated failure/timeout data.
	iterFails := make(map[int][]string, len(iterations))
	iterTimedOut := make(map[int]bool, len(iterations))
	for key, a := range aggs {
		for i := range a.timeoutIters {
			iterTimedOut[i] = true
		}
		failName := key.Test
		if failName == "" {
			failName = key.Package
		}
		for i := range a.failedIters {
			iterFails[i] = append(iterFails[i], failName)
		}
	}
	summaries := make([]IterationSummary, len(iterations))
	for i := range iterations {
		s := IterationSummary{Index: i}
		switch {
		case iterTimedOut[i]:
			s.Result = "timeout"
		case len(iterFails[i]) > 0:
			s.Result = "fail"
			sort.Strings(iterFails[i])
			s.FailingTests = iterFails[i]
		default:
			s.Result = "pass"
		}
		summaries[i] = s
	}
	rep.IterationSummaries = summaries

	logs := buildLogMap(aggs)
	return rep, logs, nil
}

// AnalyzeResults opens every `iteration-*.log.jsonl` file in resultsDir, in
// numeric-iteration order, and delegates to Analyze.
func AnalyzeResults(resultsDir string, slowThreshold time.Duration) (*Report, LogMap, error) {
	matches, err := filepath.Glob(filepath.Join(resultsDir, "iteration-*.log.jsonl"))
	if err != nil {
		return nil, nil, err
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
			return nil, nil, err
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

// WriteLogFiles writes per-test per-iteration log files under <resultsDir>/logs/
// for flagged tests and populates each flagged TestEntry's LogFiles slice with
// paths relative to resultsDir. One file is written for each iteration that has
// any captured output in logs (including iterations that passed but produced stderr
// or other output captured into the aggregate).
func WriteLogFiles(resultsDir string, rep *Report, logs LogMap) error {
	if rep == nil {
		return nil
	}
	logsDir := filepath.Join(resultsDir, "logs")
	if err := os.MkdirAll(logsDir, 0700); err != nil {
		return err
	}
	groups := [][]TestEntry{rep.Flakes, rep.Failures, rep.Timeouts, rep.Slow}
	for gi, group := range groups {
		for ei, entry := range group {
			key := testKey{Package: entry.Package, Test: entry.Test}
			m, ok := logs[key]
			if !ok || len(m) == 0 {
				continue
			}
			iterations := make([]int, 0, len(m))
			for it, out := range m {
				if out != "" {
					iterations = append(iterations, it)
				}
			}
			sort.Ints(iterations)
			paths := make([]string, 0, len(iterations))
			for _, it := range iterations {
				out := m[it]
				name := fmt.Sprintf("%s__%s__iter-%d.log",
					sanitize(shortPackage(entry.Package)), sanitize(entry.Test), it)
				abs := filepath.Join(logsDir, name)
				if err := os.WriteFile(abs, []byte(out), 0600); err != nil {
					return err
				}
				paths = append(paths, filepath.Join("logs", name))
			}
			if len(paths) > 0 {
				groups[gi][ei].LogFiles = paths
			}
		}
	}
	rep.Flakes = groups[0]
	rep.Failures = groups[1]
	rep.Timeouts = groups[2]
	rep.Slow = groups[3]
	return nil
}

// WriteCSV writes a human-readable CSV of every flagged test
// (Flakes ∪ Failures ∪ Timeouts ∪ Slow) to <resultsDir>/report.csv.
// Rows sort worst-first: (timeouts+fails) desc, then package, then test.
func WriteCSV(resultsDir string, rep *Report) error {
	if rep == nil {
		return nil
	}
	rows := flaggedRows(rep)
	f, err := os.Create(filepath.Join(resultsDir, "report.csv"))
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	if err := w.Write([]string{
		"package", "test", "category",
		"runs", "successes", "fails", "skips", "timeouts",
		"min", "max", "p50",
	}); err != nil {
		return err
	}
	for _, r := range rows {
		if err := w.Write(r.record()); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

type csvRow struct {
	Package   string
	Test      string
	Category  string
	Runs      int
	Successes int
	Fails     int
	Skips     int
	Timeouts  int
	Min       time.Duration
	Max       time.Duration
	P50       time.Duration
}

func (r csvRow) record() []string {
	return []string{
		r.Package, r.Test, r.Category,
		strconv.Itoa(r.Runs),
		strconv.Itoa(r.Successes),
		strconv.Itoa(r.Fails),
		strconv.Itoa(r.Skips),
		strconv.Itoa(r.Timeouts),
		r.Min.Round(time.Millisecond).String(),
		r.Max.Round(time.Millisecond).String(),
		r.P50.Round(time.Millisecond).String(),
	}
}

// flaggedRows builds the deduped CSV row set. A test in both Flakes and Slow
// is categorized as "flake" (primary signal wins over "slow").
func flaggedRows(rep *Report) []csvRow {
	seen := map[testKey]struct{}{}
	var rows []csvRow
	add := func(entries []TestEntry, cat string) {
		for _, e := range entries {
			k := testKey{Package: e.Package, Test: e.Test}
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			rows = append(rows, csvRow{
				Package:   e.Package,
				Test:      e.Test,
				Category:  cat,
				Runs:      e.Runs,
				Successes: e.Successes,
				Fails:     e.Fails,
				Skips:     e.Skips,
				Timeouts:  e.Timeouts,
				Min:       e.MinElapsed,
				Max:       e.MaxElapsed,
				P50:       e.P50Elapsed,
			})
		}
	}
	// Order matters: first category wins on dup.
	add(rep.Timeouts, "timeout")
	add(rep.Failures, "failure")
	add(rep.Flakes, "flake")
	add(rep.Slow, "slow")

	sort.SliceStable(rows, func(i, j int) bool {
		li := rows[i].Timeouts + rows[i].Fails
		lj := rows[j].Timeouts + rows[j].Fails
		if li != lj {
			return li > lj
		}
		if rows[i].Package != rows[j].Package {
			return rows[i].Package < rows[j].Package
		}
		return rows[i].Test < rows[j].Test
	})
	return rows
}

// PrintSummary writes a human-readable summary: headings and tests grouped by
// package under a common path prefix (tree). Broken/Flaky/Slow test lines use
// red / yellow / grey; package path rows are muted.
// Broken and Timeout entries are sorted alphabetically by package then test.
// Flaky entries are sorted by fails/runs (desc), then fails (desc), then name.
// Slow entries are sorted by max runtime (desc), then name.
func PrintSummary(w io.Writer, rep *Report) {
	if rep == nil {
		return
	}

	if n := len(rep.Failures); n > 0 {
		fails := append([]TestEntry(nil), rep.Failures...)
		sort.Slice(fails, func(i, j int) bool {
			if fails[i].Package != fails[j].Package {
				return fails[i].Package < fails[j].Package
			}
			return fails[i].Test < fails[j].Test
		})
		printSummarySectionTree(w, "Broken", n, fails, termstyle.Bad, termstyle.Bad, formatBrokenTestLine)
	}

	if n := len(rep.Flakes); n > 0 {
		flakes := append([]TestEntry(nil), rep.Flakes...)
		sort.Slice(flakes, func(i, j int) bool {
			ri := flakeFailRatio(flakes[i])
			rj := flakeFailRatio(flakes[j])
			if ri != rj {
				return ri > rj
			}
			if flakes[i].Fails != flakes[j].Fails {
				return flakes[i].Fails > flakes[j].Fails
			}
			return entryFQName(flakes[i]) < entryFQName(flakes[j])
		})
		printSummarySectionTree(w, "Flaky", n, flakes, termstyle.Flaky, termstyle.Flaky, formatFlakyTestLine)
	}

	if n := len(rep.Timeouts); n > 0 {
		touts := append([]TestEntry(nil), rep.Timeouts...)
		sort.Slice(touts, func(i, j int) bool {
			if touts[i].Package != touts[j].Package {
				return touts[i].Package < touts[j].Package
			}
			return touts[i].Test < touts[j].Test
		})
		printSummarySectionTree(w, "Timeout", n, touts, termstyle.Accent, termstyle.Accent, formatTimeoutTestLine)
	}

	if n := len(rep.Slow); n > 0 {
		slow := append([]TestEntry(nil), rep.Slow...)
		sort.Slice(slow, func(i, j int) bool {
			if slow[i].MaxElapsed != slow[j].MaxElapsed {
				return slow[i].MaxElapsed > slow[j].MaxElapsed
			}
			if slow[i].Package != slow[j].Package {
				return slow[i].Package < slow[j].Package
			}
			return slow[i].Test < slow[j].Test
		})
		printSummarySectionTree(w, "Slow", n, slow, termstyle.Muted, termstyle.Muted, formatSlowTestLine)
	}
}

func formatBrokenTestLine(e TestEntry) string {
	if e.Test == "" {
		return e.Package
	}
	return e.Test
}

func formatFlakyTestLine(e TestEntry) string {
	runs := e.Runs
	if runs < 1 {
		runs = e.Successes + e.Fails
	}
	if runs < 1 {
		runs = 1
	}
	if e.Test == "" {
		return fmt.Sprintf("%s (%d/%d)", e.Package, e.Fails, runs)
	}
	return fmt.Sprintf("%s (%d/%d)", e.Test, e.Fails, runs)
}

func formatTimeoutTestLine(e TestEntry) string {
	if e.Test == "" {
		return e.Package
	}
	return e.Test
}

func formatSlowTestLine(e TestEntry) string {
	if e.Test == "" {
		return fmt.Sprintf("%s %s", e.Package, e.MaxElapsed.Round(time.Millisecond))
	}
	return fmt.Sprintf("%s %s", e.Test, e.MaxElapsed.Round(time.Millisecond))
}

// pipeBranch returns a tree prefix: depth 1 -> "|-- ", depth 2 -> "|---- ", etc.
func pipeBranch(depth int) string {
	if depth < 1 {
		return ""
	}
	return "|" + strings.Repeat("-", 2*depth) + " "
}

// longestCommonPathPrefix returns the longest shared prefix ending at a '/'
// so grouped packages can share one root line. Empty if no '/' in common.
func longestCommonPathPrefix(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	p := append([]string(nil), paths...)
	sort.Strings(p)
	first, last := p[0], p[len(p)-1]
	cmpLen := min(len(last), len(first))
	i := 0
	for i < cmpLen && first[i] == last[i] {
		i++
	}
	common := first[:i]
	if j := strings.LastIndex(common, "/"); j >= 0 {
		return common[:j+1]
	}
	return ""
}

func printSummarySectionTree(w io.Writer, title string, n int, entries []TestEntry, headingStyle, testStyle lipgloss.Style, formatTest func(TestEntry) string) {
	fmt.Fprintln(w, headingStyle.Render(fmt.Sprintf("%s (%d)", title, n)))

	byPkg := make(map[string][]TestEntry)
	var pkgs []string
	seen := map[string]struct{}{}
	for _, e := range entries {
		if _, ok := seen[e.Package]; !ok {
			seen[e.Package] = struct{}{}
			pkgs = append(pkgs, e.Package)
		}
		byPkg[e.Package] = append(byPkg[e.Package], e)
	}
	sort.Strings(pkgs)
	for _, pkg := range pkgs {
		sort.Slice(byPkg[pkg], func(i, j int) bool { return byPkg[pkg][i].Test < byPkg[pkg][j].Test })
	}

	lcp := longestCommonPathPrefix(pkgs)
	if lcp == "" && len(pkgs) > 0 {
		lcp = pkgs[0]
		if j := strings.LastIndex(lcp, "/"); j >= 0 {
			lcp = lcp[:j+1]
		} else {
			lcp = ""
		}
	}

	if lcp != "" {
		fmt.Fprintln(w, termstyle.Muted.Render("- "+lcp))
	}

	for _, pkg := range pkgs {
		suffix := strings.TrimPrefix(pkg, lcp)
		suffix = strings.TrimPrefix(suffix, "/")
		segments := strings.Split(suffix, "/")
		var nonEmpty []string
		for _, s := range segments {
			if s != "" {
				nonEmpty = append(nonEmpty, s)
			}
		}
		depth := 0
		for _, seg := range nonEmpty {
			depth++
			fmt.Fprintln(w, termstyle.Muted.Render(pipeBranch(depth)+seg+"/"))
		}
		testDepth := len(nonEmpty) + 1
		if len(nonEmpty) == 0 {
			testDepth = 1
		}
		for _, e := range byPkg[pkg] {
			line := pipeBranch(testDepth) + formatTest(e)
			fmt.Fprintln(w, testStyle.Render(line))
		}
	}
	fmt.Fprintln(w)
}

func entryFQName(e TestEntry) string {
	if e.Test == "" {
		return e.Package
	}
	return e.Package + "." + e.Test
}

func flakeFailRatio(e TestEntry) float64 {
	runs := e.Runs
	if runs < 1 {
		runs = e.Successes + e.Fails
	}
	if runs < 1 {
		return 0
	}
	return float64(e.Fails) / float64(runs)
}

// reattributeTimeouts fixes the go-test-json quirk where a `panic: test timed out`
// is attached to whichever test most recently emitted events rather than the
// actually-stuck one. The real culprits are listed in the panic's
// "running tests:" block — move the timeout mark (and the captured stack
// trace) onto those tests.
func reattributeTimeouts(aggs map[testKey]*aggregate, newAgg func() *aggregate) {
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
			if slices.Contains(names, key.Test) {
				continue
			}
			delete(a.timeoutIters, i)
			if len(a.timeoutIters) == 0 {
				a.timedOut = false
			}
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
				continue
			}
			break
		}
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

// buildLogMap returns the raw per-iteration output for every (pkg, test) that
// has any output recorded. Callers use this to write per-test log files.
func buildLogMap(aggs map[testKey]*aggregate) LogMap {
	out := LogMap{}
	for k, a := range aggs {
		if len(a.outputs) == 0 {
			continue
		}
		m := map[int]string{}
		for i, buf := range a.outputs {
			if buf != nil && buf.Len() > 0 {
				m[i] = buf.String()
			}
		}
		if len(m) > 0 {
			out[k] = m
		}
	}
	return out
}

// stats computes min and p50 from a sample of durations.
// Returns (0, 0) for an empty sample.
func stats(samples []time.Duration) (minDur, p50 time.Duration) {
	if len(samples) == 0 {
		return 0, 0
	}
	sorted := append([]time.Duration(nil), samples...)
	slices.Sort(sorted)
	minDur = sorted[0]
	n := len(sorted)
	if n%2 == 1 {
		p50 = sorted[n/2]
	} else {
		p50 = (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return minDur, p50
}

// shortPackage keeps the last two path segments of a Go import path so log
// filenames stay under the OS NAME_MAX (255 on most filesystems). Deeply
// nested packages like github.com/.../core/services/ocr2/plugins/llo collapse
// to plugins/llo.
func shortPackage(pkg string) string {
	if pkg == "" {
		return ""
	}
	parts := strings.Split(pkg, "/")
	if len(parts) <= 2 {
		return pkg
	}
	return strings.Join(parts[len(parts)-2:], "/")
}

// sanitize turns a package path or test name into a filename-safe token.
// Replaces path separators and other hostile characters with '_'.
func sanitize(s string) string {
	if s == "" {
		return "_"
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '-', r == '.':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}

func seconds(f float64) time.Duration {
	return time.Duration(f * float64(time.Second))
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
