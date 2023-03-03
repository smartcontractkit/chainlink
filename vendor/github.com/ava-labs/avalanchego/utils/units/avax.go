// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package units

// Denominations of value
const (
	NanoAvax  uint64 = 1
	MicroAvax uint64 = 1000 * NanoAvax
	Schmeckle uint64 = 49*MicroAvax + 463*NanoAvax
	MilliAvax uint64 = 1000 * MicroAvax
	Avax      uint64 = 1000 * MilliAvax
	KiloAvax  uint64 = 1000 * Avax
	MegaAvax  uint64 = 1000 * KiloAvax
)
