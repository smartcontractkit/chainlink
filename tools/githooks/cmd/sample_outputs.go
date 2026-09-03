package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/ui"
)

func newSampleOutputsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "sample-outputs",
		Short:  "Showcase all terminal UI output styles (hidden developer tool)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			u := uiForCmd(cmd)
			out := cmd.OutOrStdout()

			renderSampleOutputs(out, u)
			return nil
		},
	}

	return cmd
}

func renderSampleOutputs(out io.Writer, u *ui.UI) {
	// 1. Headers & Section
	fmt.Fprintln(out, u.Bold("=== [githooks UI Output Showcase] ==="))
	fmt.Fprintln(out)

	// 2. Badges & Tags
	fmt.Fprintln(out, u.Bold("--- Status Badges & Tags ---"))
	fmt.Fprintf(out, "Success:     %s\n", u.BadgeSuccess("OK"))
	fmt.Fprintf(out, "Warning:     %s\n", u.BadgeWarning("MEDIUM"))
	fmt.Fprintf(out, "Danger:      %s\n", u.BadgeDanger("LARGE"))
	fmt.Fprintf(out, "Info:        %s\n", u.BadgeInfo("SMALL"))
	fmt.Fprintf(out, "Tags:        %s %s %s %s\n", u.BadgeTag("LINT"), u.BadgeTag("TEST"), u.BadgeTag("FIXED"), u.BadgeTag("CHECK FAIL"))
	fmt.Fprintln(out)

	// 3. PR Size Guard Samples
	fmt.Fprintln(out, u.Bold("--- PR Size Guard (pr-size) ---"))
	// Small
	fmt.Fprintf(out, "PR Diff Size: %s -> Classification: %s %s\n",
		u.FormatDiff(45, 45, 0, 2, "per-file-max"),
		u.BadgeInfo("SMALL"),
		u.BadgeSuccess("OK"),
	)
	// Medium
	fmt.Fprintf(out, "PR Diff Size: %s -> Classification: %s %s\n",
		u.FormatDiff(320, 280, 40, 6, "per-file-max"),
		u.BadgeWarning("MEDIUM"),
		u.BadgeSuccess("OK"),
	)
	fmt.Fprintf(out, "  %s\n", u.Dim("(Excluded 3 lock/ignored files)"))
	fmt.Fprintln(out)

	// Large PR Warning Card (Variation 1A)
	title := fmt.Sprintf("%s %s", u.BadgeDangerPill("LARGE PR"), u.Danger("(1250 lines > 500 limit)"))
	diffLine := fmt.Sprintf("%s      %s / %s in 14 files (per-file-max) • %s",
		u.Bold("Diff:"),
		u.Additions(1100),
		u.Deletions(150),
		u.Dim("Excluded: 2 lockfiles"),
	)
	bypassLine := fmt.Sprintf("%s    %s", u.Bold("Bypass:"), u.Dim("ALLOW_LARGE_PR=true git push"))
	promptLine := u.Bold("AI Prompt:") + " Execute the skill @tools/githooks/skills/split-pr/SKILL.md to break feature-branch branch up."

	fmt.Fprintln(out, u.CompactWarningCard(title, []string{diffLine, bypassLine, promptLine}))
	fmt.Fprintln(out)

	// 4. Runner Step Headers
	fmt.Fprintln(out, u.Bold("--- Runners (lint & test) ---"))
	fmt.Fprintln(out, u.RunnerHeader("LINT", "core", "4 packages"))
	fmt.Fprintln(out, u.RunnerHeader("LINT", "deployment", "./environment"))
	fmt.Fprintln(out, u.RunnerHeader("TEST", "plugins", "2 packages"))
	fmt.Fprintln(out, u.RunnerHeader("TEST", ".", "./core/logger ./core/services"))
	fmt.Fprintln(out)

	// 5. Fixer Status Items
	fmt.Fprintln(out, u.Bold("--- Fixers (whitespace & eof) ---"))
	fmt.Fprintln(out, u.StatusItem("FIXED", "core/services/app.go", ""))
	fmt.Fprintln(out, u.StatusItem("FIXED", "deployment/environment.go", ""))
	fmt.Fprintln(out, u.StatusItem("CHECK FAIL", "core/config/config.go", "erroneous trailing whitespace"))
	fmt.Fprintln(out, u.StatusItem("CHECK FAIL", "README.md", "missing trailing newline"))
	fmt.Fprintln(out)
}
