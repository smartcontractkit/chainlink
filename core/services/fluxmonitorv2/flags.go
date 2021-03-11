package fluxmonitorv2

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flags_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Flags wraps the a contract
type Flags struct {
	flags_wrapper.FlagsInterface
}

// NewFlags constructs a new Flags from a flags contract address
func NewFlags(addrHex string, ethClient eth.Client) (*Flags, error) {
	flags := &Flags{}

	if addrHex == "" {
		return flags, nil
	}

	contractAddr := common.HexToAddress(addrHex)
	contract, err := flags_wrapper.NewFlags(contractAddr, ethClient)
	if err != nil {
		return flags, err
	}

	// This is necessary due to the unfortunate fact that assigning `nil` to an
	// interface variable causes `x == nil` checks to always return false. If we
	// do this here, in the constructor, we can avoid using reflection when we
	// check `p.flags == nil` later in the code.
	if contract != nil && !reflect.ValueOf(contract).IsNil() {
		flags.FlagsInterface = contract
	}

	return flags, nil
}

// Contract returns the flags contract
func (f *Flags) Contract() flags_wrapper.FlagsInterface {
	return f.FlagsInterface
}

// ContractExists returns whether a flag contract exists
func (f *Flags) ContractExists() bool {
	return f.FlagsInterface != nil
}

// IsLowered determines whether the flag is lowered for a given contract.
// If a contract does not exist, it is considered to be lowered
func (f *Flags) IsLowered(contractAddr common.Address) (bool, error) {
	if !f.ContractExists() {
		return true, nil
	}

	flags, err := f.GetFlags(nil,
		[]common.Address{utils.ZeroAddress, contractAddr},
	)
	if err != nil {
		return true, err
	}

	return !flags[0] || !flags[1], nil
}
