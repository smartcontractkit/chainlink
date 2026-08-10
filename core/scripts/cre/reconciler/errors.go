package reconciler

import "errors"

// ErrBreakpoint is returned by Run ONLY in the --wait-at-breakpoint=false (two-invocation) flow, after the
// node-config layer has been written and the operator must re-deploy via Griddle before re-running.
// Run() (in cmd.go) maps it to ExitCodeTOMLPatch (42). In the default wait-at-breakpoint=true flow the
// breakpoint pauses in-process and this sentinel is never returned.
var ErrBreakpoint = errors.New("reconciler: node config written — apply via Griddle and re-run")

// ErrStepDeclined is returned by skipUnlessConfirmed when the user declines a step.
// Run() treats this as a clean stop (exit 0), not an error.
var ErrStepDeclined = errors.New("reconciler: step declined by operator")
