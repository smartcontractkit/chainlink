package config

import solanago "github.com/gagliardetto/solana-go"

// Case is a single negative write case. A batch of them runs in one workflow execution, which is
// what keeps the suite from paying the workflow compile/deploy/trigger cost per case.
type Case struct {
	Name string
	// FunctionToTest must match one of the switch-case labels in solwrite-negative/main.go
	FunctionToTest string
	// InvalidInput is interpreted based on FunctionToTest:
	//   - receiver/account keys: hex-encoded raw bytes of any length (empty string means "no bytes")
	//   - compute limit: decimal string
	// Cases that build their invalid input from the config addresses ignore it.
	InvalidInput string
	// ExpectedError is the substring the capability error must contain.
	ExpectedError string
}

// Config is the configuration of the Solana write negative workflow.
type Config struct {
	ChainSelector uint64
	// BatchName identifies the batch in the workflow logs the test asserts on.
	BatchName string
	// Cases all run in a single execution, so their number must stay within
	// PerWorkflow.ChainWrite.TargetsLimit - the engine counts rejected calls too.
	Cases []Case

	// ForwarderProgramID is the deployed CRE forwarder program.
	ForwarderProgramID solanago.PublicKey
	// ForwarderState is the forwarder state account the Solana capability is configured with.
	ForwarderState solanago.PublicKey
	// Receiver is a deployed, executable program used as a valid receiver, so that cases
	// which are not about the receiver reach the validation step they actually target.
	Receiver solanago.PublicKey
	// NonExecutableAccount is an existing, funded, system-owned account: it exists on-chain
	// but is not a program, which is what the "non-executable receiver" case needs.
	NonExecutableAccount solanago.PublicKey
}
