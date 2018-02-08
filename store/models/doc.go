// Package models contains the key job components used by the
// Chainlink application.
//
// Common
//
// Common contains types and functions that are useful across the
// application. Particularly dealing with the URL field, dates, and
// time.
//
// Eth
//
// Eth creates transactions and tracks transaction attempts on the
// Ethereum blockchain.
//
// Job
//
// A Job is the largest unit of work that a Chainlink node can take
// on. It will have Initiators, which is how the node knows what kind
// of work to perform, and Tasks, which are the specific instructions.
// The BridgeType is also located here, and is used for external adapters.
//
// ORM
//
// The ORM is the wrapper around the database. It gives a limited
// set of functions to allow for safe storing and withdrawing of
// information.
//
// Run
//
// A Run is the node's attempt at completing work. This comprises of
// JobRuns and TaskRuns. The JobRun keeps track of the entire job,
// and will only produce a result when all underlying Tasks are finished.
// The TaskRun keeps track of the status of an individual Task.
package models
