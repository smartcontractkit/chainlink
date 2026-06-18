// Package output centralizes harness-owned terminal writes: human-rich vs sparse
// (--ai-output), and whether stderr supports carriage-return progress (TTY).
package output

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/x/term"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
)

// Printer writes CLI messages for tools/test. Child processes (go test) still
// attach os.Stdout/os.Stderr directly where passthrough is intended.
type Printer struct {
	aiOutput   bool
	stdout     io.Writer
	stderr     io.Writer
	liveInline bool // human mode and stderr is a TTY (safe for \r progress)
}

// SkipFD disables live-inline TTY detection (tests, builders as stderr).
const SkipFD = ^uintptr(0)

// New builds a Printer. stderrFD is the stderr file descriptor for TTY detection
// (e.g. os.Stderr.Fd()). Pass SkipFD when stderr is not a real terminal fd.
func New(aiOutput bool, stdout, stderr io.Writer, stderrFD uintptr) *Printer {
	var live bool
	if stderrFD != SkipFD && !aiOutput {
		live = term.IsTerminal(stderrFD)
	}
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}
	return &Printer{
		aiOutput:   aiOutput,
		stdout:     stdout,
		stderr:     stderr,
		liveInline: live,
	}
}

// NewForTest returns a Printer with explicit live-inline behavior for unit tests.
// Production code should use New/NewFromApp (TTY detection on stderr).
func NewForTest(aiOutput bool, stdout, stderr io.Writer, liveInline bool) *Printer {
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}
	li := liveInline && !aiOutput
	return &Printer{
		aiOutput:   aiOutput,
		stdout:     stdout,
		stderr:     stderr,
		liveInline: li,
	}
}

// NewFromApp uses os.Stdout/os.Stderr and stderr's TTY bit from conf.AIOutput.
func NewFromApp(conf *config.App) *Printer {
	return New(conf.AIOutput, os.Stdout, os.Stderr, os.Stderr.Fd())
}

// AIOutput reports sparse / agent-oriented mode.
func (p *Printer) AIOutput() bool {
	return p.aiOutput
}

// LiveInlineProgress is true when human mode may use carriage-return redraws on stderr.
func (p *Printer) LiveInlineProgress() bool {
	return p.liveInline
}

// IfHuman runs fn only when not in AI output mode.
func (p *Printer) IfHuman(fn func()) {
	if !p.aiOutput {
		fn()
	}
}

// HumanStderr writes a line to stderr when in human mode.
func (p *Printer) HumanStderr(a ...any) {
	if p.aiOutput {
		return
	}
	_, _ = fmt.Fprintln(p.stderr, a...)
}

// HumanStderrf formats to stderr when in human mode.
func (p *Printer) HumanStderrf(format string, a ...any) {
	if p.aiOutput {
		return
	}
	_, _ = fmt.Fprintf(p.stderr, format, a...)
}

// HumanStdout writes a line to stdout when in human mode.
func (p *Printer) HumanStdout(a ...any) {
	if p.aiOutput {
		return
	}
	_, _ = fmt.Fprintln(p.stdout, a...)
}

// HumanFprint writes to stderr without a trailing newline when in human mode.
func (p *Printer) HumanFprint(a ...any) {
	if p.aiOutput {
		return
	}
	_, _ = fmt.Fprint(p.stderr, a...)
}

// HumanStderrWriter returns stderr for human-mode helpers that need io.Writer.
// When AI output is enabled, returns io.Discard so callers can always write.
func (p *Printer) HumanStderrWriter() io.Writer {
	if p.aiOutput {
		return io.Discard
	}
	return p.stderr
}

// SparseStdoutln prints one line to stdout in AI mode (machine-oriented).
func (p *Printer) SparseStdoutln(a ...any) {
	if !p.aiOutput {
		return
	}
	_, _ = fmt.Fprintln(p.stdout, a...)
}

// Stderrf always writes formatted text to stderr (errors, diagnostics).
func (p *Printer) Stderrf(format string, a ...any) {
	_, _ = fmt.Fprintf(p.stderr, format, a...)
}

// Stdoutln always writes a line to stdout.
func (p *Printer) Stdoutln(a ...any) {
	_, _ = fmt.Fprintln(p.stdout, a...)
}

// ClearInline clears the current stderr line when live inline progress is active.
func (p *Printer) ClearInline() {
	if !p.liveInline {
		return
	}
	_, _ = fmt.Fprint(p.stderr, "\r\033[K")
}

// WarnWriter returns stderr for notes that should appear in human mode and stay
// quiet in AI mode (e.g. diagnose -count hints).
func (p *Printer) WarnWriter() io.Writer {
	if p.aiOutput {
		return io.Discard
	}
	return p.stderr
}
