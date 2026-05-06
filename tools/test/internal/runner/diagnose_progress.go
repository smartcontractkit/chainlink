package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/termstyle"
)

const chainlinkModulePrefix = "github.com/smartcontractkit/chainlink/v2"

// packagePatternsFromEnd returns trailing non-flag arguments. This matches the usual
// `go test [flags] [packages]` layout (package patterns last).
func packagePatternsFromEnd(args []string) []string {
	var pkgs []string
	for i := len(args) - 1; i >= 0; i-- {
		if strings.HasPrefix(args[i], "-") {
			break
		}
		pkgs = append(pkgs, args[i])
	}
	slices.Reverse(pkgs)
	return pkgs
}

func shortenChainlinkImportPath(importPath string) string {
	if importPath == "" {
		return ""
	}
	if importPath == chainlinkModulePrefix {
		return "."
	}
	p := chainlinkModulePrefix + "/"
	return strings.TrimPrefix(importPath, p)
}

// listTestPackageCount runs `go list -test -e` for the trailing package patterns
// in go test arguments (see packagePatternsFromEnd). On error or no patterns,
// returns an error or zero packages.
func listTestPackageCount(ctx context.Context, repoRoot string, goTestArgs []string) (int, error) {
	pkgs := packagePatternsFromEnd(goTestArgs)
	if len(pkgs) == 0 {
		return 0, errors.New("no package patterns in go test arguments (put packages last, after flags)")
	}
	//nolint:gosec // it's fine
	cmd := exec.CommandContext(ctx, "go", append([]string{"list", "-test", "-e", "-f", "{{.ImportPath}}"}, pkgs...)...)
	cmd.Dir = repoRoot
	cmd.Env = os.Environ()
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	n := 0
	for line := range strings.SplitSeq(string(out), "\n") {
		if strings.TrimSpace(line) != "" {
			n++
		}
	}
	if n == 0 {
		return 0, errors.New("go list returned no packages")
	}
	return n, nil
}

// diagnoseProgress tracks completed packages from a go test -json stream.
type diagnoseProgress struct {
	mu         sync.Mutex
	done       map[string]struct{}
	lastPkg    string
	pkgOutcome map[string]string // package import path → pass|fail|skip (package-level events only)
	total      int               // -1 when denominator is unknown (go list failed or empty)
}

type parallelDiagnoseProgress struct {
	mu              sync.Mutex
	renderMu        sync.Mutex
	totalIterations int
	completed       int
	active          map[int]parallelIterationProgress
	poolStartedAt   time.Time
}

type parallelIterationProgress struct {
	startedAt time.Time
}

type activeIterElapsed struct {
	iteration int           // 0-based diagnose iteration index
	elapsed   time.Duration // wall since this iteration's go test started
}

func newParallelDiagnoseProgress(totalIterations int) *parallelDiagnoseProgress {
	return newParallelDiagnoseProgressAt(totalIterations, time.Now())
}

func newParallelDiagnoseProgressAt(totalIterations int, poolStartedAt time.Time) *parallelDiagnoseProgress {
	return &parallelDiagnoseProgress{
		totalIterations: totalIterations,
		active:          make(map[int]parallelIterationProgress),
		poolStartedAt:   poolStartedAt,
	}
}

func newDiagnoseProgress(totalPackages int) *diagnoseProgress {
	return &diagnoseProgress{
		done:       make(map[string]struct{}),
		pkgOutcome: make(map[string]string),
		total:      totalPackages,
	}
}

func (p *parallelDiagnoseProgress) start(iteration int) {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.active[iteration] = parallelIterationProgress{startedAt: time.Now()}
}

// startAtForTest records an active iteration with a fixed startedAt (package runner tests).
func (p *parallelDiagnoseProgress) startAtForTest(iteration int, startedAt time.Time) {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.active[iteration] = parallelIterationProgress{startedAt: startedAt}
}

func (p *parallelDiagnoseProgress) finish(iteration int) {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.active, iteration)
	p.completed++
}

// bumpCompletedForTest increments the done counter (package runner tests).
func (p *parallelDiagnoseProgress) bumpCompletedForTest(n int) {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.completed += n
}

func (p *parallelDiagnoseProgress) withRenderLock(fn func()) {
	if p == nil {
		fn()
		return
	}
	p.renderMu.Lock()
	defer p.renderMu.Unlock()
	fn()
}

// renderSnapshot returns completed iteration count, total planned iterations,
// per-active-iteration elapsed (sorted by iteration index), and wall time since
// the parallel pool began.
func (p *parallelDiagnoseProgress) renderSnapshot(now time.Time) (completed, total int, actives []activeIterElapsed, poolElapsed time.Duration) {
	if p == nil {
		return 0, 0, nil, 0
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	completed = p.completed
	total = p.totalIterations
	poolElapsed = now.Sub(p.poolStartedAt)
	if poolElapsed < 0 {
		poolElapsed = 0
	}
	poolElapsed = poolElapsed.Round(time.Second)
	for iter, pr := range p.active {
		elapsed := now.Sub(pr.startedAt)
		if elapsed < 0 {
			elapsed = 0
		}
		actives = append(actives, activeIterElapsed{iteration: iter, elapsed: elapsed.Round(time.Second)})
	}
	slices.SortFunc(actives, func(a, b activeIterElapsed) int {
		return a.iteration - b.iteration
	})
	return completed, total, actives, poolElapsed
}

// onTestJSONLine updates state from one JSONL line. Returns true if the number
// of completed packages increased (for throttled redraws).
func (p *diagnoseProgress) onTestJSONLine(line []byte) (completedIncreased bool) {
	if len(line) == 0 || line[0] != '{' {
		return false
	}
	var ev TestEvent
	if err := json.Unmarshal(line, &ev); err != nil {
		return false
	}
	if ev.Package != "" {
		p.mu.Lock()
		p.lastPkg = ev.Package
		p.mu.Unlock()
	}
	if !isPackageTerminalEvent(&ev) {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pkgOutcome[ev.Package] = ev.Action
	before := len(p.done)
	p.done[ev.Package] = struct{}{}
	return len(p.done) > before
}

func isPackageTerminalEvent(ev *TestEvent) bool {
	if ev.Package == "" || ev.Test != "" {
		return false
	}
	switch ev.Action {
	case "pass", "fail", "skip":
		return true
	default:
		return false
	}
}

func (p *diagnoseProgress) snapshot() (completed int, total int, lastPkg string, outcome string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.done), p.total, p.lastPkg, p.pkgOutcome[p.lastPkg]
}

// progressBracket wraps inner (already styled) in muted square brackets.
func progressBracket(inner string) string {
	return termstyle.Muted.Render("[") + inner + termstyle.Muted.Render("]")
}

// renderDiagnoseProgressLine writes one status line to w when liveInline is true
// (TTY stderr in human mode). Otherwise it is a no-op so logs are not spammed.
// diagnoseRunStart is when the overall diagnose run began; if zero, only the
// per-iteration bracket is shown.
func renderDiagnoseProgressLine(w io.Writer, iteration, iterations int, iterElapsed time.Duration, diagnoseRunStart time.Time, now time.Time, liveInline bool) {
	if !liveInline {
		return
	}
	iterBracket := fmt.Sprintf("iter %d/%d (%s)", iteration, iterations, iterElapsed.Round(time.Second).String())
	line := progressBracket(termstyle.Label.Render(iterBracket))
	if !diagnoseRunStart.IsZero() {
		runEl := now.Sub(diagnoseRunStart)
		if runEl < 0 {
			runEl = 0
		}
		line += "  " + progressBracket(termstyle.Muted.Render(runEl.Round(time.Second).String()))
	}
	fmt.Fprint(w, "\r\033[K")
	fmt.Fprint(w, line)
}

func renderParallelDiagnoseProgressLine(w io.Writer, prog *parallelDiagnoseProgress, now time.Time, liveInline bool) {
	if !liveInline || prog == nil {
		return
	}
	completed, totalIters, actives, poolElapsed := prog.renderSnapshot(now)
	line := progressBracket(termstyle.Label.Render(fmt.Sprintf("done %d/%d", completed, totalIters)))
	var lineSb275 strings.Builder
	for _, a := range actives {
		lineSb275.WriteString("  " + progressBracket(termstyle.Label.Render(fmt.Sprintf("iter %d (%s)", a.iteration+1, a.elapsed.String()))))
	}
	line += lineSb275.String()
	line += "  " + progressBracket(termstyle.Muted.Render(poolElapsed.String()))
	fmt.Fprint(w, "\r\033[K")
	fmt.Fprint(w, line)
}

func ellipsizeRight(s string, maxLen int) string {
	if maxLen <= 3 || len(s) <= maxLen {
		return s
	}
	return "…" + s[len(s)-(maxLen-3):]
}
