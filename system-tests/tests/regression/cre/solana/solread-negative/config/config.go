package config

// Outcome is what a case expects the Solana capability to answer with.
type Outcome int

const (
	// OutcomeError expects the call to fail with an error containing ExpectedError. It is used
	// for malformed requests, which the capability rejects before they reach the RPC.
	OutcomeError Outcome = iota
	// OutcomeErrorOrEmpty expects either an error or an empty response. It is used for requests
	// that are well-formed but point at something that does not exist, where answering with
	// "nothing" is as correct as failing - only returning real data is a failure.
	OutcomeErrorOrEmpty
)

// Case is a single negative read case. A batch of them runs in one workflow execution, which is
// what keeps the suite from paying the workflow compile/deploy/trigger cost per case.
type Case struct {
	Name           string
	FunctionToTest string
	// InvalidInput is interpreted based on FunctionToTest:
	//   - account/program keys and transaction signatures: hex-encoded raw bytes of any length
	//     (an empty string means "no bytes", which reaches the capability as a nil field)
	//   - slots, commitment/encoding enum values and batch sizes: decimal string
	//   - encoded messages: passed through verbatim
	InvalidInput string
	Outcome      Outcome
	// ExpectedError is the substring the capability error must contain. Only used by OutcomeError.
	ExpectedError string
}

// Config is the configuration of the Solana read negative workflow.
type Config struct {
	ChainSelector uint64
	// BatchName identifies the batch in the workflow logs the test asserts on.
	BatchName string
	// Cases all run in a single execution, so their number must stay within
	// PerWorkflow.ChainRead.CallLimit - the engine counts rejected calls too.
	Cases []Case
}
