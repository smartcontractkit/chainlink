package types

import "errors"

var (
	ErrGlobalWorkflowCountLimitReached   = errors.New("global workflow count limit reached: the platform has reached its maximum number of active workflows. Please deactivate unused workflows or contact support to request a limit increase")
	ErrPerOwnerWorkflowCountLimitReached = errors.New("per-owner workflow count limit reached: you have reached the maximum number of active workflows allowed per owner. Please deactivate unused workflows before creating new ones")
)
