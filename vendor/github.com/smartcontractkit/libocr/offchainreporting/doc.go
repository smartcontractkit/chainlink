// Package offchainreporting implements the Chainlink Offchain Reporting Protocol
//
// A note about concurrency
//
// We spawn lots of goroutines in this package. As a general rule, we keep track
// of all of them using the subprocesses package. We typically signal shutdowns
// using contexts and then do a subprocesses.Wait() to ensure that we're not
// leaking goroutines.

package offchainreporting
