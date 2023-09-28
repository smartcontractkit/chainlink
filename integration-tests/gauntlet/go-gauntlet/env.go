package gogauntlet

import "fmt"

type EnvVar string

const (
	ReportNameEnvVar EnvVar = "REPORT_NAME"
)

// Env is the environment that the gauntlet CLI will run in.  Use to pass environment variables to gauntlet.
type Env map[EnvVar]string

func (e Env) ToArr() []string {
	var arr []string
	for k, v := range e {
		arr = append(arr, fmt.Sprintf("%s=%s", k, v))
	}
	return arr
}

// CommandOpts holds the environment in which the gauntlet CLI will be run and additional CLI flags to pass to a command.
type CommandOpts struct {
	Env   Env
	Flags []string
}

func WithEnv(env Env) func(*CommandOpts) {
	return func(co *CommandOpts) {
		co.Env = env
	}
}

func WithFlags(flags []string) func(*CommandOpts) {
	return func(co *CommandOpts) {
		co.Flags = append(co.Flags, flags...)
	}
}
