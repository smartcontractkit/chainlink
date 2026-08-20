package stellar

import "github.com/smartcontractkit/chainlink-stellar/deployment/cre"

// Artifact filenames of the CRE contract WASM committed in chainlink-stellar
// (deployment/cre/artifacts/). These are duplicated from
// chainlink-stellar/deployment/cre/artifacts.go so callers in the deployment
// module and downstream system-tests/lib do not need to import chainlink-stellar.
const (
	// ReadFixtureWasm is the CRE ReadContract test fixture (contracts/cre/test/read_fixture).
	ReadFixtureWasm = cre.ReadFixtureWasm

	// ForwarderWasm is the CRE forwarder (contracts/cre/forwarder).
	ForwarderWasm = cre.ForwarderWasm

	// ReceiverWasm is the CRE test receiver (contracts/cre/test/receiver).
	ReceiverWasm = cre.ReceiverWasm

	// RejectingReceiverWasm is the CRE test receiver that always rejects on_report.
	RejectingReceiverWasm = cre.RejectingReceiverWasm
)

// Artifact returns the compiled WASM for one of the filename constants above.
// The bytes are those embedded in the pinned chainlink-stellar module, kept in
// sync with the contract sources by that repo's check-generated CI job —
// nothing is compiled, downloaded, or resolved at deploy time.
func Artifact(name string) ([]byte, error) {
	return cre.Artifact(name)
}
