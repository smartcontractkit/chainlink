package devmachine

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// tuiMessage carries reporter output into the Bubble Tea model.
type tuiMessage struct {
	kind string // "phase", "info", "warn"
	text string
}

// tuiDoneMsg signals that the orchestrator finished (or failed).
type tuiDoneMsg struct {
	err error
}

// tuiReporter is a Reporter that forwards events to the TUI channel.
type tuiReporter struct {
	ch chan<- tea.Msg
}

func (t *tuiReporter) Phase(name string) {
	t.ch <- tuiMessage{kind: "phase", text: name}
}

func (t *tuiReporter) Info(format string, args ...any) {
	t.ch <- tuiMessage{kind: "info", text: fmt.Sprintf(format, args...)}
}

func (t *tuiReporter) Warn(format string, args ...any) {
	t.ch <- tuiMessage{kind: "warn", text: fmt.Sprintf(format, args...)}
}

type tuiModel struct {
	ch    <-chan tea.Msg
	phase string
	logs  []string
	done  bool
	err   error
}

func newTUIModel(ch <-chan tea.Msg) tuiModel {
	return tuiModel{ch: ch}
}

func (m tuiModel) Init() tea.Cmd {
	return waitForTUI(m.ch)
}

func waitForTUI(ch <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return tuiDoneMsg{}
		}
		return msg
	}
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil
	case tuiMessage:
		switch msg.kind {
		case "phase":
			m.phase = msg.text
		default:
			m.logs = append(m.logs, msg.text)
		}
		return m, waitForTUI(m.ch)
	case tuiDoneMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	}
	return m, nil
}

func (m tuiModel) View() string {
	var b strings.Builder
	b.WriteString("CRE DevMachine provisioning\n\n")
	if m.phase != "" {
		fmt.Fprintf(&b, "  -> %s\n\n", m.phase)
	}
	start := 0
	const maxLines = 15
	if len(m.logs) > maxLines {
		start = len(m.logs) - maxLines
	}
	for _, l := range m.logs[start:] {
		fmt.Fprintf(&b, "  %s\n", l)
	}
	if m.done {
		if m.err != nil {
			b.WriteString("\n  failed\n")
		} else {
			b.WriteString("\n  done\n")
		}
	}
	return b.String()
}

// runTUI drives phases from..to inside a Bubble Tea program and returns the
// orchestrator error (if any).
func runTUI(ctx context.Context, o *Orchestrator, from, to phaseID) error {
	ch := make(chan tea.Msg, 256)
	o.rep = &tuiReporter{ch: ch}

	p := tea.NewProgram(newTUIModel(ch), tea.WithContext(ctx))
	go func() {
		err := o.runPhaseRange(ctx, from, to)
		ch <- tuiDoneMsg{err: err}
	}()

	final, err := p.Run()
	if err != nil {
		return err
	}
	if m, ok := final.(tuiModel); ok {
		return m.err
	}
	return nil
}
