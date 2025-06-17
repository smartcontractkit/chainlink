package types

import "errors"

var (
	ErrGlobalWorkflowCountLimitReached   = errors.New("global workflow count limit reached")
	ErrPerOwnerWorkflowCountLimitReached = errors.New("per owner workflow count limit reached")
)
