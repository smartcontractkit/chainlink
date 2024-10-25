package chains

import "errors"

var (
	ErrLOOPPUnsupported = errors.New("LOOPP not yet supportedd")
	ErrChainDisabled    = errors.New("chain is disabled")
)
