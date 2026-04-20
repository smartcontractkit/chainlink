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

	"github.com/charmbracelet/x/term"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/termstyle"
)

// listTestPackageCount runs `go list -test -e` for the same pattern as `go test`
// and returns how many import paths would be listed. Used as the progress bar
// denominator; on error or zero lines callers should treat the total as unknown.
func listTestPackageCount(ctx context.Context, repoRoot, pattern string) (int, error) {
	cmd := exec.CommandContext(ctx, "go", "list", "-test", "-e", "-f", "{{.ImportPath}}", pattern)
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

// surveyProgress tracks completed packages from a go test -json stream.
type surveyProgress struct {
	mu      sync.Mutex
	done    map[string]struct{}
	lastPkg string
	total   int // -1 when denominator is unknown (go list failed or empty)
}

func newSurveyProgress(totalPackages int) *surveyProgress {
	return &surveyProgress{
		done:  make(map[string]struct{}),
		total: totalPackages,
	}
}

// onTestJSONLine updates state from one JSONL line. Returns true if the number
// of completed packages increased (for throttled redraws).
func (p *surveyProgress) onTestJSONLine(line []byte) (completedIncreased bool) {
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

func (p *surveyProgress) snapshot() (completed int, total int, lastPkg string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.done), p.total, p.lastPkg
}

func stderrTermWidth() int {
	fd := os.Stderr.Fd()
	if !term.IsTerminal(fd) {
		return 80
	}
	w, _, err := term.GetSize(fd)
	if err != nil || w < 40 {
		return 80
	}
	return w
}

func surveyProgressBarWidth(termW, reserved int) int {
	const minBar, maxBar = 8, 36
	w := termW - reserved
	if w < minBar {
		return minBar
	}
	if w > maxBar {
		return maxBar
	}
	return w
}

// renderSurveyProgressLine writes one status line to w. When isTTY is true,
// caller should prefix with "\r\033[K" to overwrite the previous line.
func renderSurveyProgressLine(w io.Writer, iteration, iterations int, elapsed time.Duration, prog *surveyProgress, isTTY bool) {
	completed, total, lastPkg := prog.snapshot()

	termW := stderrTermWidth()
	meta := fmt.Sprintf("iter %d/%d", iteration, iterations)
	countStr := fmt.Sprintf("%d/%d", completed, total)
	if total < 0 {
		countStr = fmt.Sprintf("%d/?", completed)
	}

	const pkgMaxChars = 42
	shortPkg := ellipsizeRight(lastPkg, pkgMaxChars)

	barReserved := len(meta) + len(countStr) + 12 + len(shortPkg) + 8 // rough fixed + elapsed
	barW := surveyProgressBarWidth(termW, barReserved)

	var barStyled string
	if total > 0 {
		filled := completed * barW / total
		if completed > 0 && filled == 0 {
			filled = 1
		}
		if filled > barW {
			filled = barW
		}
		filledStr := strings.Repeat("█", filled)
		emptyStr := strings.Repeat("░", barW-filled)
		barStyled = termstyle.Filled.Render(filledStr) + termstyle.Empty.Render(emptyStr)
	} else {
		// Indeterminate: one bright cell shuttles in a dim track.
		track := make([]rune, barW)
		for i := range track {
			track[i] = '░'
		}
		if barW > 0 {
			pos := int(time.Now().UnixMilli()/200) % barW
			track[pos] = '█'
		}
		s := string(track)
		hi := strings.IndexRune(s, '█')
		if hi >= 0 {
			barStyled = termstyle.Empty.Render(s[:hi]) + termstyle.Filled.Render("█") + termstyle.Empty.Render(s[hi+1:])
		} else {
			barStyled = termstyle.Empty.Render(s)
		}
	}

	line := termstyle.Label.Render(meta) + "  " + barStyled + "  " + termstyle.Accent.Render(countStr)
	if shortPkg != "" {
		line += "  " + termstyle.Muted.Render(shortPkg)
	}
	line += "  " + termstyle.Muted.Render(elapsed.Round(time.Second).String())
	if isTTY {
		fmt.Fprint(w, "\r\033[K")
	} else {
		fmt.Fprint(w, "\n")
	}
	fmt.Fprint(w, line)
	if !isTTY {
		fmt.Fprint(w, "\n")
	}
}

func ellipsizeRight(s string, maxLen int) string {
	if maxLen <= 3 || len(s) <= maxLen {
		return s
	}
	return "…" + s[len(s)-(maxLen-3):]
}
