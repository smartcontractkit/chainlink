package nix

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Shell is a wrapper around a nix shell process. It allows to run commands
// in the same context, preserving the environment variables set in the shell
// and the state set by initial execution of "nix develop".
type Shell struct {
	cmd    *exec.Cmd
	stdin  *bufio.Writer
	stdout *bufio.Reader
	mu     sync.Mutex
}

const ErrCommandFailed = "command failed with exit code"

func NewNixShell(folder string, globalEnvVars map[string]string) (*Shell, error) {
	cmd := exec.Command("nix", "develop", "--command", "sh")
	cmd.Dir = folder

	// Set global environment variables available to all subsequent commands
	cmd.Env = os.Environ()
	fmt.Println("Nix shell will use the following global environment variables:")
	for key, value := range globalEnvVars {
		fmt.Printf("%s=%s\n", key, value)
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Shell{
		cmd:    cmd,
		stdin:  bufio.NewWriter(stdin),
		stdout: bufio.NewReader(stdout),
	}, nil
}

func (ns *Shell) RunCommand(command string) (string, error) {
	return ns.RunCommandWithEnvVars(command, map[string]string{})
}

func (ns *Shell) RunCommandWithEnvVars(command string, envVars map[string]string) (string, error) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	// send stderr to stdout, append exit code to the end of the output and
	// add end marker to signal the end of the command output
	endMarker := "END_OF_COMMAND_OUTPUT"
	fullCommand := fmt.Sprintf("%s 2>&1; echo %s $?\n", command, endMarker)

	// Set command-specific environment variables
	if len(envVars) > 0 {
		fmt.Println("Setting the following command-specific environment variables:")
	}
	for key, value := range envVars {
		fmt.Printf("%s=%s\n", key, value)
		_, err := ns.stdin.WriteString(fmt.Sprintf("export %s=%s\n", key, value))
		if err != nil {
			return "", err
		}
	}

	_, err := ns.stdin.WriteString(fullCommand)
	if err != nil {
		return "", err
	}
	if err := ns.stdin.Flush(); err != nil {
		return "", err
	}

	// read output until the end marker is found
	var output strings.Builder
	var exitCode int
	for {
		line, err := ns.stdout.ReadString('\n')
		fmt.Print(line)
		if err != nil {
			return "", err
		}
		if strings.HasPrefix(line, endMarker) {
			_, scanRrr := fmt.Sscanf(line, endMarker+" %d", &exitCode)
			if scanRrr != nil {
				exitCode = 1
			}
			break
		}
		output.WriteString(line)
	}

	if exitCode != 0 {
		return output.String(), fmt.Errorf("%s %d", ErrCommandFailed, exitCode)
	}

	return strings.TrimSpace(output.String()), nil
}

func (ns *Shell) Close() error {
	return ns.cmd.Process.Kill()
}
