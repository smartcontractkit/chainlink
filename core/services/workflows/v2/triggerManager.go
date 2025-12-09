package v2

// This will be used to register all triggers and is responsible for retrying subscription and keeping track of trigger status
// across all workflows.  Moving the retry behaviour here will ensure restarted standard capability triggers are re-registered.
// It is assumed that the register trigger function on all capabilities is idempotent.

// The initial call should block, but perhaps not later calls?
