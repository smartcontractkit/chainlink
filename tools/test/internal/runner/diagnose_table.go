package runner

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/output"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/termstyle"
)

// Fixed column widths for streaming rows (each row is formatted independently).
const (
	diagnoseColIter    = 5
	diagnoseColResult  = 8
	diagnoseColTests   = 8
	diagnoseColCount   = 8
	diagnoseColRuntime = 10
)

func printDiagnoseIterationTableHeader(out *output.Printer) {
	if out.AIOutput() {
		return
	}
	out.HumanStderr(termstyle.Muted.Render(diagnoseTableHeaderPlain()))
	out.HumanStderr(termstyle.Muted.Render(strings.Repeat("─", len(diagnoseTableHeaderPlain()))))
}

func diagnoseTableHeaderPlain() string {
	return fmt.Sprintf("%5s  %-8s  %8s  %8s  %8s  %8s  %10s",
		"Iter", "Result", "Tests", "Failures", "Timeouts", "Slow", "Runtime")
}

func formatDiagnoseIterationTableRow(iter int, d IterationDigest, dur time.Duration) string {
	iterCol := lipgloss.PlaceHorizontal(diagnoseColIter, lipgloss.Right, termstyle.Label.Render(strconv.Itoa(iter)))
	resCol := lipgloss.PlaceHorizontal(diagnoseColResult, lipgloss.Left, renderIterationResultHuman(d.Result))
	testsCol := lipgloss.PlaceHorizontal(diagnoseColTests, lipgloss.Right, termstyle.Muted.Render(strconv.Itoa(d.RanTests)))
	failCol := lipgloss.PlaceHorizontal(diagnoseColCount, lipgloss.Right, diagnoseTableCountStyled(d.FailTests, "fail"))
	toCol := lipgloss.PlaceHorizontal(diagnoseColCount, lipgloss.Right, diagnoseTableCountStyled(d.TimeoutTests, "timeout"))
	slowCol := lipgloss.PlaceHorizontal(diagnoseColCount, lipgloss.Right, diagnoseTableCountStyled(d.SlowTests, "slow"))
	rt := termstyle.Muted.Render(dur.Round(time.Second).String())
	rtCol := lipgloss.PlaceHorizontal(diagnoseColRuntime, lipgloss.Right, rt)
	gap := "  "
	return lipgloss.JoinHorizontal(lipgloss.Top,
		iterCol, gap, resCol, gap, testsCol, gap, failCol, gap, toCol, gap, slowCol, gap, rtCol)
}

func diagnoseTableCountStyled(n int, kind string) string {
	s := strconv.Itoa(n)
	switch kind {
	case "fail", "timeout":
		if n == 0 {
			return termstyle.OK.Render(s)
		}
		return termstyle.Bad.Render(s)
	case "slow":
		if n == 0 {
			return termstyle.OK.Render(s)
		}
		return termstyle.Flaky.Render(s)
	default:
		return termstyle.Muted.Render(s)
	}
}
