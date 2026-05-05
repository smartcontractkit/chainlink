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

// chainlinkModulePrefix is trimmed from import paths in the diagnose progress line
// so the status shows repo-relative paths (e.g. core/foo).
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
}

type parallelIterationProgress struct {
	completed int
	total     int
	lastPkg   string
	outcome   string
}

func newDiagnoseProgress(totalPackages int) *diagnoseProgress {
	return &diagnoseProgress{
		done:       make(map[string]struct{}),
		pkgOutcome: make(map[string]string),
		total:      totalPackages,
	}
}

func newParallelDiagnoseProgress(totalIterations int) *parallelDiagnoseProgress {
	return &parallelDiagnoseProgress{
		totalIterations: totalIterations,
		active:          make(map[int]parallelIterationProgress),
	}
}

func (p *parallelDiagnoseProgress) start(iteration, totalPackages int) {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.active[iteration] = parallelIterationProgress{total: totalPackages}
}

func (p *parallelDiagnoseProgress) update(iteration int, prog *diagnoseProgress) {
	if p == nil || prog == nil {
		return
	}
	completed, total, lastPkg, outcome := prog.snapshot()
	p.mu.Lock()
	defer p.mu.Unlock()
	p.active[iteration] = parallelIterationProgress{
		completed: completed,
		total:     total,
		lastPkg:   lastPkg,
		outcome:   outcome,
	}
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

func (p *parallelDiagnoseProgress) withRenderLock(fn func()) {
	if p == nil {
		fn()
		return
	}
	p.renderMu.Lock()
	defer p.renderMu.Unlock()
	fn()
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

// packageProgressCountSegment returns styled "a/b · p%" (or "a/?") for bracketed progress lines.
func packageProgressCountSegment(completed, total int) string {
	if total < 0 {
		return termstyle.Accent.Render(fmt.Sprintf("%d/?", completed))
	}
	pct := 0
	if total > 0 {
		pct = completed * 100 / total
	}
	return termstyle.Accent.Render(fmt.Sprintf("%d/%d", completed, total)) +
		termstyle.Muted.Render(" · ") +
		termstyle.Accent.Render(fmt.Sprintf("%d%%", pct))
}

// packageOutcomeMark returns a short suffix after the displayed package path:
// pass/fail/skip from package-level JSON events, or an hourglass while that path
// is active but no terminal result is recorded yet, or empty when there is no path.
func packageOutcomeMark(action, displayPkg string) string {
	if displayPkg != "" && action == "" {
		return " ⌛"
	}
	switch action {
	case "pass":
		return " ✅"
	case "fail":
		return " ❌"
	case "skip":
		return " ⏭"
	default:
		return ""
	}
}

// renderDiagnoseProgressLine writes one status line to w when liveInline is true
// (TTY stderr in human mode). Otherwise it is a no-op so logs are not spammed.
func renderDiagnoseProgressLine(w io.Writer, iteration, iterations int, elapsed time.Duration, prog *diagnoseProgress, liveInline bool) {
	if !liveInline {
		return
	}
	completed, total, lastPkg, outcome := prog.snapshot()

	meta := fmt.Sprintf("%d/%d", iteration, iterations)

	const pkgMaxChars = 42
	displayPkg := shortenChainlinkImportPath(lastPkg)
	mark := packageOutcomeMark(outcome, displayPkg)
	markReserve := 0
	if displayPkg != "" {
		markReserve = 8 // room for terminal marks or hourglass (display width approx)
	}
	shortPkg := ellipsizeRight(displayPkg, pkgMaxChars-markReserve) + mark

	line := progressBracket(termstyle.Label.Render(meta)) + "  " +
		progressBracket(packageProgressCountSegment(completed, total))
	if shortPkg != "" {
		line += "  " + progressBracket(termstyle.Muted.Render(shortPkg)) // path + ⌛ while running, or ✅/❌/⏭ when done
	}
	line += "  " + progressBracket(termstyle.Muted.Render(elapsed.Round(time.Second).String()))
	fmt.Fprint(w, "\r\033[K")
	fmt.Fprint(w, line)
}

func renderParallelDiagnoseProgressLine(w io.Writer, prog *parallelDiagnoseProgress, elapsed time.Duration, liveInline bool) {
	if !liveInline || prog == nil {
		return
	}
	completed, totalIterations, activeCount, iteration, current := prog.snapshot()
	line := progressBracket(termstyle.Label.Render(fmt.Sprintf("done %d/%d", completed, totalIterations))) + "  " +
		progressBracket(termstyle.Accent.Render(fmt.Sprintf("active %d", activeCount)))
	if activeCount > 0 {
		line += "  " + progressBracket(termstyle.Label.Render(fmt.Sprintf("iter %d", iteration+1))) + "  " +
			progressBracket(packageProgressCountSegment(current.completed, current.total))
		displayPkg := shortenChainlinkImportPath(current.lastPkg)
		if displayPkg != "" {
			mark := packageOutcomeMark(current.outcome, displayPkg)
			line += "  " + progressBracket(termstyle.Muted.Render(ellipsizeRight(displayPkg, 42)+mark))
		}
	}
	line += "  " + progressBracket(termstyle.Muted.Render(elapsed.Round(time.Second).String()))
	fmt.Fprint(w, "\r\033[K")
	fmt.Fprint(w, line)
}

func (p *parallelDiagnoseProgress) snapshot() (completed, totalIterations, activeCount, iteration int, current parallelIterationProgress) {
	p.mu.Lock()
	defer p.mu.Unlock()
	iteration = -1
	for i, prog := range p.active {
		if iteration == -1 || i < iteration {
			iteration = i
			current = prog
		}
	}
	return p.completed, p.totalIterations, len(p.active), iteration, current
}

func ellipsizeRight(s string, maxLen int) string {
	if maxLen <= 3 || len(s) <= maxLen {
		return s
	}
	return "…" + s[len(s)-(maxLen-3):]
}
