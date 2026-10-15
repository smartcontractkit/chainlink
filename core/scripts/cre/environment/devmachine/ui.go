package devmachine

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Reporter receives progress notifications from the phase runner. Two
// implementations exist: a plain text reporter and a Bubble Tea driven one.
type Reporter interface {
	// Phase announces the start of a named phase.
	Phase(name string)
	// Info prints an informational line.
	Info(format string, args ...any)
	// Warn prints a warning line.
	Warn(format string, args ...any)
}

// plainReporter writes human readable progress to an io.Writer.
type plainReporter struct {
	out io.Writer
}

// NewPlainReporter returns a Reporter that writes to out (defaults to stdout).
func NewPlainReporter(out io.Writer) Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &plainReporter{out: out}
}

func (p *plainReporter) Phase(name string) {
	fmt.Fprintf(p.out, "\n==> %s\n", name)
}

func (p *plainReporter) Info(format string, args ...any) {
	fmt.Fprintf(p.out, "    "+format+"\n", args...)
}

func (p *plainReporter) Warn(format string, args ...any) {
	fmt.Fprintf(p.out, "    ! "+format+"\n", args...)
}

// phaseID is a stable identifier for each provisioning phase.
type phaseID string

const (
	phasePreflight phaseID = "P0"
	phaseMachine   phaseID = "P1"
	phaseBootstrap phaseID = "P2"
	phaseArtifacts phaseID = "P3"
	phaseTransfer  phaseID = "P4"
	phaseLoad      phaseID = "P5"
	phaseHandoff   phaseID = "P6"
)

// orderedPhases is the execution order used by --from/--only.
var orderedPhases = []phaseID{
	phasePreflight,
	phaseMachine,
	phaseBootstrap,
	phaseArtifacts,
	phaseTransfer,
	phaseLoad,
	phaseHandoff,
}

func phaseTitle(id phaseID) string {
	switch id {
	case phasePreflight:
		return "Preflight checks"
	case phaseMachine:
		return "DevMachine"
	case phaseBootstrap:
		return "Remote bootstrap"
	case phaseArtifacts:
		return "Build artifacts"
	case phaseTransfer:
		return "Transfer to machine"
	case phaseLoad:
		return "Load images"
	case phaseHandoff:
		return "Handoff"
	}
	return string(id)
}

// parsePhase resolves a user supplied phase string (e.g. "P3", "3",
// "artifacts") to a phaseID.
func parsePhase(s string) (phaseID, error) {
	t := strings.TrimSpace(strings.ToUpper(s))
	t = strings.TrimPrefix(t, "PHASE")
	switch t {
	case "0", "P0", "PREFLIGHT":
		return phasePreflight, nil
	case "1", "P1", "MACHINE", "DEVMACHINE":
		return phaseMachine, nil
	case "2", "P2", "BOOTSTRAP":
		return phaseBootstrap, nil
	case "3", "P3", "ARTIFACTS", "BUILD":
		return phaseArtifacts, nil
	case "4", "P4", "TRANSFER":
		return phaseTransfer, nil
	case "5", "P5", "LOAD":
		return phaseLoad, nil
	case "6", "P6", "HANDOFF":
		return phaseHandoff, nil
	}
	return "", fmt.Errorf("unknown phase %q", s)
}

// shouldRun reports whether phase id should execute given the --from/--only
// selectors.
func shouldRun(id phaseID, only, from phaseID) bool {
	if only != "" {
		return id == only
	}
	if from != "" {
		return phaseIndex(id) >= phaseIndex(from)
	}
	return true
}

func phaseIndex(id phaseID) int {
	for i, p := range orderedPhases {
		if p == id {
			return i
		}
	}
	return -1
}
