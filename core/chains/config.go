package chains

import (
	"errors"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrMultipleChainFound = errors.New("multiple chains found with the same chain ID")
)
