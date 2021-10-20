// Package models contain the key job components used by the
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
// JobSpec
//
// A JobSpec is the largest unit of work that a Chainlink node can take
// on. It will have Initiators, which is how a JobRun is started from
// the job definition, and Tasks, which are the specific instructions
// for what work needs to be performed.
// The BridgeType is also located here, and is used to define the location
// (URL) of external adapters.
//
// ORM
//
// The ORM is the wrapper around the database. It gives a limited
// set of functions to allow for safe storing and withdrawing of
// information.
//
// Run
//
// A Run is the actual invocation of work being done on the Job and Task.
// This comprises of JobRuns and TaskRuns. A JobRun is like a workflow where
// the steps are the TaskRuns.
//
// i.e. We have a Scheduler Initiator that creates a JobRun every monday
// based on a JobDefinition. And in turn, those JobRuns have TaskRuns based
// on the JobDefinition's TaskDefinitions.
package models
