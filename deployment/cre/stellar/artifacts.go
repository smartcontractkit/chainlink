package stellar

// Artifact filenames produced by `stellar contract build` in the chainlink-stellar
// cargo workspace (target/wasm32v1-none/release/). These are duplicated from
// chainlink-stellar/deployment/cre/artifacts.go so callers in the deployment
// module and downstream system-tests/lib do not need to import chainlink-stellar.
const (
	// ReadFixtureWasm is the CRE ReadContract test fixture (contracts/cre/test/read_fixture).
	ReadFixtureWasm = "read_fixture.wasm"

	// ForwarderWasm is the CRE forwarder (contracts/cre/forwarder).
	ForwarderWasm = "forwarder.wasm"

	// ReceiverWasm is the CRE test receiver (contracts/cre/test/receiver).
	ReceiverWasm = "receiver.wasm"

	// RejectingReceiverWasm is the CRE test receiver that always rejects on_report.
	RejectingReceiverWasm = "rejecting_receiver.wasm"
)
