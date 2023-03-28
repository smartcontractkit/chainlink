package gogauntlet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// TODO: The below TODO was copied over from the go-gauntlet repo.  Not
// clear if there is a value add, but this change will be considered when revisiting
// the ticket here: https://smartcontract-it.atlassian.net/browse/FMS-1204
//
// TODO: Change path from string to an os.File
type Gauntlet struct {
	path     string
	Defaults DefaultConfig
}

func (g Gauntlet) Path() string {
	return g.path
}

// NewGauntletFromLocal is a constructor for a gauntlet object using a local file.
// We simply check if the path exists. We do not check whether it's a valid Gauntlet binary
func NewGauntletFromLocal(input string, dc *DefaultConfig) (*Gauntlet, error) {
	if !CheckPathExists(input) {
		return &Gauntlet{}, fmt.Errorf("Path %q does not exist", input)
	}

	return &Gauntlet{path: input, Defaults: *dc}, nil
}

func (g Gauntlet) ExecWriteCommandWithEnv(ctx context.Context, args []string, env Env) (Report, error) {
	return g.ExecWriteCommand(ctx, args, WithEnv(env))
}

func (g Gauntlet) ExecWriteCommand(ctx context.Context, args []string, opts ...func(*CommandOpts)) (Report, error) {
	co := &CommandOpts{
		Flags: g.Defaults.CreateDefaultFlags(),
	}

	for _, fn := range opts {
		fn(co)
	}

	allArgs := append(args, co.Flags...)

	err := g.execBlockingCommand(ctx, allArgs, co.Env.ToArr())
	if err != nil {
		return Report{}, err
	}

	reportName, ok := co.Env[ReportNameEnvVar]
	if !ok {
		reportName = "report"
	}

	return g.parseJsonReport(reportName)
}

func MakeMarshaledInput(in interface{}) (string, error) {
	if in == nil {
		return "", nil
	}
	marshaled, err := json.Marshal(in)
	if err != nil {
		return "", errors.New("Could not marshal input")
	}
	return ("--input=" + string(marshaled)), nil
}

func (g Gauntlet) ExecuteContractWithFlags(ctx context.Context, cmdName, contractAddress string, input interface{}, additionalFlags []string, opts ...func(*CommandOpts)) (Report, error) {
	marshaledInput, err := MakeMarshaledInput(input)
	if err != nil {
		return Report{}, err
	}

	allArgs := []string{cmdName, contractAddress, marshaledInput}
	allArgs = append(allArgs, additionalFlags...)

	return g.ExecWriteCommand(ctx, allArgs, opts...)
}

func (g Gauntlet) execBlockingCommand(ctx context.Context, args []string, envVar []string) error {
	// TODO(mstreet3): This seems like a major issue in using go-gauntlet
	// there is no verification of the binary and the path to execute comes
	// in as an environment variable.  We should verify the executable somehow
	// before pushing this to production.
	//
	// https://smartcontract-it.atlassian.net/browse/FMS-1204
	cmd := exec.CommandContext(ctx, g.path, args...) //nolint:all
	if g.Defaults.Out() != nil {                     //nolint:all
		cmd.Stdout = g.Defaults.Out()
	} else {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// add environment variable to this command specifically
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "SKIP_PROMPTS=true")
	cmd.Env = append(cmd.Env, envVar...)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// CheckPathExists checks if the directory or file passed into it is valid
func CheckPathExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}
