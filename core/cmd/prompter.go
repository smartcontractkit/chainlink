package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/term"
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
		fmt.Print(err)
		os.Exit(1)
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
		bytePwd, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		clearLine()
		rval = string(bytePwd)
	})
	return rval
}

// IsTerminal checks if the current process is executing in a terminal, this
// should be used to decide when to use PasswordPrompt.
func (tp terminalPrompter) IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// Explicitly reset terminal state in the event of a signal (CTRL+C)
// to ensure typed characters are echoed in terminal:
// https://groups.google.com/forum/#!topic/Golang-nuts/kTVAbtee9UA
func withTerminalResetter(f func()) {
	osSafeStdin := int(os.Stdin.Fd())

	initialTermState, err := term.GetState(osSafeStdin)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		err := term.Restore(osSafeStdin, initialTermState)
		if err != nil {
			fmt.Printf("Error restoring terminal: %v", err)
		}
		os.Exit(1)
	}()

	f()
	signal.Stop(c)
}

func clearLine() {
	fmt.Print("\r" + strings.Repeat(" ", 60) + "\r")
}
