package ui_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/ui"
)

func TestUI_PlainFallback(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	u := ui.NewWithMode(&buf, ui.ColorNever)

	assert.False(t, u.ColorEnabled())

	t.Run("badges in plain mode", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "[ OK ]", u.BadgeSuccess("OK"))
		assert.Equal(t, "[ WARN ]", u.BadgeWarning("WARN"))
		assert.Equal(t, "[ LARGE ]", u.BadgeDanger("LARGE"))
		assert.Equal(t, "[ SMALL ]", u.BadgeInfo("SMALL"))
		assert.Equal(t, "[LINT]", u.BadgeTag("LINT"))
	})

	t.Run("runner header in plain mode", func(t *testing.T) {
		t.Parallel()
		got := u.RunnerHeader("LINT", "core", "4 packages")
		assert.Equal(t, "[LINT] module 'core' (4 packages)", got)
	})

	t.Run("status item in plain mode", func(t *testing.T) {
		t.Parallel()
		got := u.StatusItem("FIXED", "core/services/app.go", "")
		assert.Equal(t, "[FIXED] core/services/app.go", got)

		gotWithReason := u.StatusItem("CHECK FAIL", "core/config.go", "trailing whitespace")
		assert.Equal(t, "[CHECK FAIL] core/config.go: trailing whitespace", gotWithReason)
	})

	t.Run("format diff in plain mode", func(t *testing.T) {
		t.Parallel()
		got := u.FormatDiff(120, 100, 20, 3, "per-file-max")
		assert.Equal(t, "120 effective lines (+100, -20 in 3 files) [strategy: per-file-max]", got)
	})

	t.Run("warning box in plain mode", func(t *testing.T) {
		t.Parallel()
		box := u.WarningBox("LARGE PR WARNING", []string{
			"PR diff size: 850 effective lines (+800, -50 across 12 files)",
			"Limit exceeded: 500 lines",
			"Split into smaller, focused PRs",
		})
		assert.Contains(t, box, "LARGE PR WARNING")
		assert.Contains(t, box, "850 effective lines")
		assert.Contains(t, box, "Split into smaller, focused PRs")
	})
}

func TestUI_ColorAlways(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	u := ui.NewWithMode(&buf, ui.ColorAlways)

	assert.True(t, u.ColorEnabled())

	t.Run("renders content with styling enabled", func(t *testing.T) {
		assert.Contains(t, u.BadgeSuccess("OK"), "OK")
		assert.Contains(t, u.FormatDiff(120, 100, 20, 3, "per-file-max"), "120 effective lines")
		assert.Contains(t, u.WarningBox("TITLE", []string{"tip"}), "TITLE")
	})
}
