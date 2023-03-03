// Package services contain the key components of the Chainlink
// node. This includes the Application, JobRunner, LogListener,
// and Scheduler.
//
// Application
//
// The Application is the main component used for starting and
// stopping the Chainlink node.
//
// JobRunner
//
// The JobRunner keeps track of Runs within a Job and ensures
// that they're executed in order. Within each Run, the tasks
// are also executed from the JobRunner.
//
// JobSubscriber
//
// The JobSubscriber coordinates running job events with
// the EventLog in the Store, and also subscribes to the given
// address on the Ethereum blockchain.
//
// Scheduler
//
// The Scheduler ensures that recurring events are executed
// according to their schedule, and one-time events occur only
// when the specified time has passed.
package services
