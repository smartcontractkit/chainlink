package updater

import "errors"

var (
	// ErrModOperation indicates a failure in a module-related operation.
	ErrModOperation = errors.New("module operation failed")

	// ErrInvalidConfig indicates invalid configuration parameters.
	ErrInvalidConfig = errors.New("invalid configuration")
)
