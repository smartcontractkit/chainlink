package plugins

import (
	"fmt"
	"os/exec"
)

// LoopExecConfig is configuration used to register the LOOP and generate an exec
type LoopExecConfig struct {
	Id  string // unique string used by the node to track the LOOP. typically supplied by the loop logger name
	Cmd string // string value of executable to exec
	LoggingConfig
}

// makeLoopCmd is helper to ensure synchronization between the loop registry and os cmd to exec the LOOP
func MakeLoopCmd(loopRegistry *LoopRegistry, lcfg LoopExecConfig) (func() *exec.Cmd, error) {
	registeredLoop, err := loopRegistry.Register(lcfg.Id, lcfg.LoggingConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to register %s LOOP plugin: %w", lcfg.Id, err)
	}
	return func() *exec.Cmd {
		cmd := exec.Command(lcfg.Cmd) //#nosec G204 -- we control the value of the cmd so the lint/sec error is a false positive
		SetCmdEnvFromConfig(cmd, registeredLoop.EnvCfg)
		return cmd
	}, nil
}
