package pkg

import a "github.com/smartcontractkit/chainlink-env/pkg/alias"

// Common labels for k8s envs
const (
	TTLLabelKey = "janitor/ttl"
)

// Environment types, envs got selected by having a label of that type
const (
	EnvTypeEVM5             = "evm-5-minimal"
	EnvTypeEVM5RemoteRunner = "evm-5-remote-runner"
)

func PGIsReadyCheck() *[]*string {
	return &[]*string{
		a.Str("pg_isready"),
		a.Str("-U"),
		a.Str("postgres"),
	}
}
