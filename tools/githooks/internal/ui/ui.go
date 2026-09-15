package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// ColorMode specifies when color output should be enabled.
type ColorMode int

const (
	ColorAuto ColorMode = iota
	ColorAlways
	ColorNever
)

// UI encapsulates terminal styling and formatting.
type UI struct {
	w        io.Writer
	renderer *lipgloss.Renderer
	mode     ColorMode

	// Cached styles
	styleSuccess    lipgloss.Style
	styleWarning    lipgloss.Style
	styleDanger     lipgloss.Style
	styleInfo       lipgloss.Style
	styleDim        lipgloss.Style
	styleBold       lipgloss.Style
	styleAdd        lipgloss.Style
	styleDel        lipgloss.Style
	styleTag        lipgloss.Style
	styleDangerPill lipgloss.Style
	styleCardBorder lipgloss.Style
}

// New creates a UI instance bound to w with ColorAuto mode.
func New(w io.Writer) *UI {
	return NewWithMode(w, ColorAuto)
}

// NewWithMode creates a UI instance with explicit ColorMode.
func NewWithMode(w io.Writer, mode ColorMode) *UI {
	if w == nil {
		w = io.Discard
	}

	r := lipgloss.NewRenderer(w)
	switch mode {
	case ColorNever:
		r.SetColorProfile(termenv.Ascii)
	case ColorAlways:
		if r.ColorProfile() == termenv.Ascii {
			r.SetColorProfile(termenv.ANSI256)
		}
	case ColorAuto:
		// Check NO_COLOR env var explicitly if termenv didn't already
		if os.Getenv("NO_COLOR") != "" {
			r.SetColorProfile(termenv.Ascii)
		}
	}

	u := &UI{
		w:        w,
		renderer: r,
		mode:     mode,
	}
	u.initStyles()
	return u
}

func (u *UI) initStyles() {
	r := u.renderer
	u.styleSuccess = r.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))  // bright green
	u.styleWarning = r.NewStyle().Bold(true).Foreground(lipgloss.Color("214")) // vibrant amber/orange
	u.styleDanger = r.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))  // bright red
	u.styleInfo = r.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))     // cyan/blue
	u.styleDim = r.NewStyle().Faint(true)
	u.styleBold = r.NewStyle().Bold(true)
	u.styleAdd = r.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	u.styleDel = r.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
	u.styleTag = r.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	u.styleDangerPill = r.NewStyle().Bold(true).Background(lipgloss.Color("196")).Foreground(lipgloss.Color("231")).Padding(0, 1)
	u.styleCardBorder = r.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("208")).
		Padding(0, 1)
}

// ColorEnabled returns whether color styling is active for this UI instance.
func (u *UI) ColorEnabled() bool {
	return u.renderer.ColorProfile() != termenv.Ascii
}

// BadgeSuccess renders a success badge, e.g. "[ OK ]" or "[ PASS ]".
func (u *UI) BadgeSuccess(text string) string {
	if !u.ColorEnabled() {
		return fmt.Sprintf("[ %s ]", text)
	}
	return u.styleSuccess.Render(fmt.Sprintf("[ %s ]", text))
}

// BadgeWarning renders a warning badge, e.g. "[ MEDIUM ]" or "[ WARN ]".
func (u *UI) BadgeWarning(text string) string {
	if !u.ColorEnabled() {
		return fmt.Sprintf("[ %s ]", text)
	}
	return u.styleWarning.Render(fmt.Sprintf("[ %s ]", text))
}

// BadgeDanger renders a danger badge, e.g. "[ LARGE ]" or "[ FAIL ]".
func (u *UI) BadgeDanger(text string) string {
	if !u.ColorEnabled() {
		return fmt.Sprintf("[ %s ]", text)
	}
	return u.styleDanger.Render(fmt.Sprintf("[ %s ]", text))
}

// BadgeInfo renders an info badge, e.g. "[ SMALL ]" or "[ INFO ]".
func (u *UI) BadgeInfo(text string) string {
	if !u.ColorEnabled() {
		return fmt.Sprintf("[ %s ]", text)
	}
	return u.styleInfo.Render(fmt.Sprintf("[ %s ]", text))
}

// BadgeTag renders a command or status tag, e.g. "[LINT]" or "[FIXED]".
func (u *UI) BadgeTag(tag string) string {
	if !u.ColorEnabled() {
		return fmt.Sprintf("[%s]", tag)
	}
	return u.styleTag.Render(fmt.Sprintf("[%s]", tag))
}

// FormatDiff formats diff metrics with colored +/- counters.
func (u *UI) FormatDiff(effectiveLines, additions, deletions, filesChanged int, strategy string) string {
	var addStr, delStr, fileStr, effStr, strategyStr string

	if u.ColorEnabled() {
		effStr = u.styleBold.Render(fmt.Sprintf("%d effective lines", effectiveLines))
		addStr = u.styleAdd.Render(fmt.Sprintf("+%d", additions))
		delStr = u.styleDel.Render(fmt.Sprintf("-%d", deletions))
		fileStr = fmt.Sprintf("in %d files", filesChanged)
		strategyStr = u.styleDim.Render(fmt.Sprintf("[strategy: %s]", strategy))
	} else {
		effStr = fmt.Sprintf("%d effective lines", effectiveLines)
		addStr = fmt.Sprintf("+%d", additions)
		delStr = fmt.Sprintf("-%d", deletions)
		fileStr = fmt.Sprintf("in %d files", filesChanged)
		strategyStr = fmt.Sprintf("[strategy: %s]", strategy)
	}

	return fmt.Sprintf("%s (%s, %s %s) %s", effStr, addStr, delStr, fileStr, strategyStr)
}

// Additions renders additions count, e.g. "+100" in green when color enabled.
func (u *UI) Additions(n int) string {
	if !u.ColorEnabled() {
		return fmt.Sprintf("+%d", n)
	}
	return u.styleAdd.Render(fmt.Sprintf("+%d", n))
}

// Deletions renders deletions count, e.g. "-20" in red when color enabled.
func (u *UI) Deletions(n int) string {
	if !u.ColorEnabled() {
		return fmt.Sprintf("-%d", n)
	}
	return u.styleDel.Render(fmt.Sprintf("-%d", n))
}

// BadgeDangerPill renders a solid background pill badge, e.g. " LARGE PR ".
func (u *UI) BadgeDangerPill(text string) string {
	if !u.ColorEnabled() {
		return fmt.Sprintf("[ %s ]", text)
	}
	return u.styleDangerPill.Render(text)
}

// CompactWarningCard renders a compact framed warning card (Variation 1A).
func (u *UI) CompactWarningCard(title string, lines []string) string {
	if u.ColorEnabled() {
		var content strings.Builder
		content.WriteString(title)
		content.WriteString("\n")
		for i, line := range lines {
			content.WriteString(line)
			if i < len(lines)-1 {
				content.WriteString("\n")
			}
		}
		return u.styleCardBorder.Render(content.String())
	}

	// Plain text box fallback
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s\n", title)
	for _, line := range lines {
		fmt.Fprintf(&sb, "   %s\n", line)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// WarningBox renders a framed warning box with a title and bulleted guidance.
func (u *UI) WarningBox(title string, items []string) string {
	var body strings.Builder
	if u.ColorEnabled() {
		body.WriteString(u.styleDanger.Render("⚠️  " + title))
		body.WriteString("\n\n")
		for _, item := range items {
			body.WriteString("  • ")
			body.WriteString(item)
			body.WriteString("\n")
		}
		return u.styleCardBorder.Render(strings.TrimRight(body.String(), "\n"))
	}

	// Plain text box fallback
	fmt.Fprintf(&body, "⚠️  [%s]\n", title)
	for _, item := range items {
		fmt.Fprintf(&body, "   • %s\n", item)
	}
	return strings.TrimRight(body.String(), "\n")
}

// RunnerHeader formats a step header for tool runners like lint and test.
func (u *UI) RunnerHeader(tag, target, details string) string {
	tagStr := u.BadgeTag(tag)
	if details != "" {
		if u.ColorEnabled() {
			return fmt.Sprintf("%s module %s %s", tagStr, u.styleBold.Render("'"+target+"'"), u.styleDim.Render("("+details+")"))
		}
		return fmt.Sprintf("%s module '%s' (%s)", tagStr, target, details)
	}
	if u.ColorEnabled() {
		return fmt.Sprintf("%s module %s", tagStr, u.styleBold.Render("'"+target+"'"))
	}
	return fmt.Sprintf("%s module '%s'", tagStr, target)
}

// StatusItem formats an item status update (e.g. whitespace or EOF fixer).
func (u *UI) StatusItem(status, path, note string) string {
	tagStr := u.BadgeTag(status)
	if note != "" {
		if u.ColorEnabled() {
			return fmt.Sprintf("%s %s: %s", tagStr, path, u.styleDim.Render(note))
		}
		return fmt.Sprintf("%s %s: %s", tagStr, path, note)
	}
	return fmt.Sprintf("%s %s", tagStr, path)
}

// Dim renders faint text when color is enabled.
func (u *UI) Dim(text string) string {
	if !u.ColorEnabled() {
		return text
	}
	return u.styleDim.Render(text)
}

// Bold renders bold text when color is enabled.
func (u *UI) Bold(text string) string {
	if !u.ColorEnabled() {
		return text
	}
	return u.styleBold.Render(text)
}

// Success renders success colored text.
func (u *UI) Success(text string) string {
	if !u.ColorEnabled() {
		return text
	}
	return u.styleSuccess.Render(text)
}

// Danger renders danger colored text.
func (u *UI) Danger(text string) string {
	if !u.ColorEnabled() {
		return text
	}
	return u.styleDanger.Render(text)
}
