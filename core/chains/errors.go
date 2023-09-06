package chains

import "errors"

var (
	ErrLOOPPUnsupported = errors.New("LOOPP not yet supported")
	ErrChainDisabled    = errors.New("chain is disabled")
)
