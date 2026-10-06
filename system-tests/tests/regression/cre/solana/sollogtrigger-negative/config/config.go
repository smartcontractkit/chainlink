package config

import solanago "github.com/gagliardetto/solana-go"

// Config is the configuration of the Solana log trigger negative workflow.
// Every case builds a log trigger filter that the Solana capability is expected to reject when
// the workflow engine registers the trigger.
type Config struct {
	ChainSelector uint64
	// FunctionToTest must match one of the switch-case labels in sollogtrigger-negative/main.go
	FunctionToTest string
	// InvalidInput is interpreted based on FunctionToTest:
	//   - addresses: hex-encoded raw bytes of any length (an empty string means "no bytes")
	// Cases whose invalid input is a missing field ignore it.
	InvalidInput string
	// ProgramID is a well-formed program address, used by the cases whose invalid input is
	// something other than the address itself.
	ProgramID solanago.PublicKey
}
