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
	return fmt.Sprintf("%5s  %-8s  %8s  %8s  %8s  %10s",
		"Iter", "Result", "Failures", "Timeouts", "Slow", "Runtime")
}

func formatDiagnoseIterationTableRow(iter int, d IterationDigest, dur time.Duration) string {
	iterCol := lipgloss.PlaceHorizontal(diagnoseColIter, lipgloss.Right, termstyle.Label.Render(strconv.Itoa(iter)))
	resCol := lipgloss.PlaceHorizontal(diagnoseColResult, lipgloss.Left, renderIterationResultHuman(d.Result))
	failCol := lipgloss.PlaceHorizontal(diagnoseColCount, lipgloss.Right, diagnoseTableCountStyled(d.FailTests, "fail"))
	toCol := lipgloss.PlaceHorizontal(diagnoseColCount, lipgloss.Right, diagnoseTableCountStyled(d.TimeoutTests, "timeout"))
	slowCol := lipgloss.PlaceHorizontal(diagnoseColCount, lipgloss.Right, diagnoseTableCountStyled(d.SlowTests, "slow"))
	rt := termstyle.Muted.Render(dur.Round(time.Second).String())
	rtCol := lipgloss.PlaceHorizontal(diagnoseColRuntime, lipgloss.Right, rt)
	gap := "  "
	return lipgloss.JoinHorizontal(lipgloss.Top,
		iterCol, gap, resCol, gap, failCol, gap, toCol, gap, slowCol, gap, rtCol)
}

func diagnoseTableCountStyled(n int, kind string) string {
	var num string
	switch {
	case n > 0 && kind == "slow":
		num = termstyle.Flaky.Render(strconv.Itoa(n))
	case n > 0 && (kind == "fail" || kind == "timeout"):
		num = termstyle.Bad.Render(strconv.Itoa(n))
	default:
		num = termstyle.Muted.Render(strconv.Itoa(n))
	}
	return num
}
