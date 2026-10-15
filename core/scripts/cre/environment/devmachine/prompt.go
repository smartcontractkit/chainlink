package devmachine

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Prompter provides simple line based interactive input. It is only used for
// values that cannot be auto-detected (team, cost center, repo locations).
type Prompter struct {
	in  *bufio.Reader
	out io.Writer
}

// NewPrompter returns a Prompter reading from stdin.
func NewPrompter(out io.Writer) *Prompter {
	if out == nil {
		out = os.Stdout
	}
	return &Prompter{in: bufio.NewReader(os.Stdin), out: out}
}

// Ask prompts for a value, returning def when the user enters nothing.
func (p *Prompter) Ask(label, def string) (string, error) {
	if def != "" {
		fmt.Fprintf(p.out, "%s [%s]: ", label, def)
	} else {
		fmt.Fprintf(p.out, "%s: ", label)
	}
	line, err := p.in.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	val := strings.TrimSpace(line)
	if val == "" {
		return def, nil
	}
	return val, nil
}

// AskYesNo prompts for a yes/no answer.
func (p *Prompter) AskYesNo(label string, def bool) (bool, error) {
	hint := "y/N"
	if def {
		hint = "Y/n"
	}
	for {
		fmt.Fprintf(p.out, "%s [%s]: ", label, hint)
		line, err := p.in.ReadString('\n')
		if err != nil && err != io.EOF {
			return false, err
		}
		switch strings.ToLower(strings.TrimSpace(line)) {
		case "":
			return def, nil
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		}
		fmt.Fprintln(p.out, "please answer 'y' or 'n'")
	}
}

// WaitForEnter blocks until the user presses enter.
func (p *Prompter) WaitForEnter(label string) error {
	fmt.Fprintf(p.out, "%s (press enter to continue) ", label)
	_, err := p.in.ReadString('\n')
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
