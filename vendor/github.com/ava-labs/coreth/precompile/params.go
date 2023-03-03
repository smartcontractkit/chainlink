// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
)

// Gas costs for stateful precompiles
// can be added here eg.
// const MintGasCost = 30_000

// AddressRange represents a continuous range of addresses
type AddressRange struct {
	Start common.Address
	End   common.Address
}

// Contains returns true iff [addr] is contained within the (inclusive)
func (a *AddressRange) Contains(addr common.Address) bool {
	addrBytes := addr.Bytes()
	return bytes.Compare(addrBytes, a.Start[:]) >= 0 && bytes.Compare(addrBytes, a.End[:]) <= 0
}

// Designated addresses of stateful precompiles
// Note: it is important that none of these addresses conflict with each other or any other precompiles
// in core/vm/contracts.go.
// We start at 0x0100000000000000000000000000000000000000 and will increment by 1 from here to reduce
// the risk of conflicts.
var (
	UsedAddresses = []common.Address{
		// precompile contract addresses can be added here
	}

	// ReservedRanges contains addresses ranges that are reserved
	// for precompiles and cannot be used as EOA or deployed contracts.
	ReservedRanges = []AddressRange{
		{
			// reserved for coreth precompiles
			common.HexToAddress("0x0100000000000000000000000000000000000000"),
			common.HexToAddress("0x01000000000000000000000000000000000000ff"),
		},
	}
)
