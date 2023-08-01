package types

import (
	"fmt"
)

// Opt is an option for a gas estimator
type Opt int

const (
	// OptForceRefetch forces the estimator to bust a cache if necessary
	OptForceRefetch Opt = iota
)

type Fee fmt.Stringer

// Constants for chains with unique fee denominations
const (
	EVM = "EVM"
	SOL = "Solana"
)
