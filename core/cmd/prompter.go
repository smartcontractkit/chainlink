package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// Prompter implements the Prompt function to be used to display at
// the console.
type Prompter interface {
	Prompt(string) string
	PasswordPrompt(string) string
	IsTerminal() bool
}

// terminalPrompter is used to display and read input from the user.
type terminalPrompter struct{}

// NewTerminalPrompter prompts the user via terminal.
func NewTerminalPrompter() Prompter {
	return terminalPrompter{}
}

// Prompt displays the prompt for the user to enter the password and
// reads their input.
func (tp terminalPrompter) Prompt(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	line, err := reader.ReadString('\n')
	if err != nil {
		logger.Fatal(err)
	}
	clearLine()
	return strings.TrimSpace(line)
}

// PasswordPrompt displays the prompt for the user to enter the password and
// reads their input.
func (tp terminalPrompter) PasswordPrompt(prompt string) string {
	var rval string
	withTerminalResetter(func() {
		fmt.Print(prompt)
		bytePwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			logger.Fatal(err)
		}
		clearLine()
		rval = string(bytePwd)
	})
	return rval
}

// IsTerminal checks if the current process is executing in a terminal, this
// should be used to decide when to use PasswordPrompt.
func (tp terminalPrompter) IsTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

// Explicitly reset terminal state in the event of a signal (CTRL+C)
// to ensure typed characters are echoed in terminal:
// https://groups.google.com/forum/#!topic/Golang-nuts/kTVAbtee9UA
func withTerminalResetter(f func()) {
	osSafeStdin := int(os.Stdin.Fd())

	initialTermState, err := terminal.GetState(osSafeStdin)
	if err != nil {
		logger.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		err := terminal.Restore(osSafeStdin, initialTermState)
		logger.ErrorIf(err, "failed when restore terminal")
		os.Exit(1)
	}()

	f()
	signal.Stop(c)
}

func clearLine() {
	fmt.Printf("\r" + strings.Repeat(" ", 60) + "\r")
}
