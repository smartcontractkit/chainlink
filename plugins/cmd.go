package plugins

import (
	"fmt"
	"os/exec"
)

// CmdConfig is configuration used to register the LOOP and generate an exec
type CmdConfig struct {
	ID  string   // unique string used by the node to track the LOOP. typically supplied by the loop logger name
	Cmd string   // string value of executable to exec
	Env []string // environment variables as described in [exec.Cmd.Env]
}

// NewCmdFactory is helper to ensure synchronization between the loop registry and os cmd to exec the LOOP
func NewCmdFactory(register func(id string) (*RegisteredLoop, error), lcfg CmdConfig) (func() *exec.Cmd, error) {
	registeredLoop, err := register(lcfg.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to register %s LOOP plugin: %w", lcfg.ID, err)
	}
	return func() *exec.Cmd {
		cmd := exec.Command(lcfg.Cmd) //#nosec G204 -- we control the value of the cmd so the lint/sec error is a false positive
		cmd.Env = append(cmd.Env, lcfg.Env...)
		cmd.Env = append(cmd.Env, registeredLoop.EnvCfg.AsCmdEnv()...)
		return cmd
	}, nil
}
