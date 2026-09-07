package stellar

import "github.com/smartcontractkit/chainlink-stellar/deployment/cre"

// Artifact filenames for Stellar contracts embedded in chainlink-stellar.
//
// Deployment changesets should use this package instead of importing the
// chainlink-stellar artifact package directly.
const (
	// MCMSWasm is the Stellar Many Chain MultiSig contract.
	MCMSWasm = cre.MCMSWasm

	// TimelockWasm is the Stellar Timelock contract.
	TimelockWasm = cre.TimelockWasm

	// ReadFixtureWasm is the CRE ReadContract test fixture.
	ReadFixtureWasm = cre.ReadFixtureWasm

	// ForwarderWasm is the CRE Forwarder contract.
	ForwarderWasm = cre.ForwarderWasm

	// ReceiverWasm is the CRE test receiver.
	ReceiverWasm = cre.ReceiverWasm

	// RejectingReceiverWasm is the CRE test receiver that always rejects
	// on_report calls.
	RejectingReceiverWasm = cre.RejectingReceiverWasm
)

// Artifact returns the compiled WASM for the requested Stellar contract.
//
// The artifacts are embedded in the pinned chainlink-stellar module. Nothing
// is compiled, downloaded, or resolved at deployment time.
func Artifact(name string) ([]byte, error) {
	return cre.Artifact(name)
}
