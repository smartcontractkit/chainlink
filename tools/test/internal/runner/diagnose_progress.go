package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	for i, j := 0, len(pkgs)-1; i < j; i, j = i+1, j-1 {
		pkgs[i], pkgs[j] = pkgs[j], pkgs[i]
	}
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

func newDiagnoseProgress(totalPackages int) *diagnoseProgress {
	return &diagnoseProgress{
		done:       make(map[string]struct{}),
		pkgOutcome: make(map[string]string),
		total:      totalPackages,
	}
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

// renderDiagnoseProgressLine writes one status line to w.
func renderDiagnoseProgressLine(w io.Writer, iteration, iterations int, elapsed time.Duration, prog *diagnoseProgress, isTTY bool) {
	if !isTTY {
		return // no-op when not a TTY (ai output doesn't need this)
	}
	completed, total, lastPkg, outcome := prog.snapshot()

	meta := fmt.Sprintf("iter %d/%d", iteration, iterations)
	var countStr string
	if total < 0 {
		countStr = fmt.Sprintf("%d/?", completed)
	} else {
		pct := 0
		if total > 0 {
			pct = completed * 100 / total
		}
		countStr = fmt.Sprintf("%d/%d %d%%", completed, total, pct)
	}

	const pkgMaxChars = 42
	displayPkg := shortenChainlinkImportPath(lastPkg)
	mark := packageOutcomeMark(outcome, displayPkg)
	markReserve := 0
	if displayPkg != "" {
		markReserve = 8 // room for terminal marks or hourglass (display width approx)
	}
	shortPkg := ellipsizeRight(displayPkg, pkgMaxChars-markReserve) + mark

	line := termstyle.Label.Render(meta) + "  " + termstyle.Accent.Render(countStr)
	if shortPkg != "" {
		line += "  " + termstyle.Muted.Render(shortPkg) // path + ⌛ while running, or ✅/❌/⏭ when done
	}
	line += "  " + termstyle.Muted.Render(elapsed.Round(time.Second).String())
	fmt.Fprint(w, "\r\033[K")
	fmt.Fprint(w, line)
}

func ellipsizeRight(s string, maxLen int) string {
	if maxLen <= 3 || len(s) <= maxLen {
		return s
	}
	return "…" + s[len(s)-(maxLen-3):]
}
