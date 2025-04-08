package fluxmonitorv2

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/flags_wrapper"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

type Flags interface {
	ContractExists() bool
	IsLowered(contractAddr common.Address) (bool, error)
	Address() common.Address
	ParseLog(log types.Log) (generated.AbigenLog, error)
}

// ContractFlags wraps the a contract
type ContractFlags struct {
	flags_wrapper.FlagsInterface
}

// NewFlags constructs a new Flags from a flags contract address
func NewFlags(addrHex string, ethClient evmclient.Client) (Flags, error) {
	flags := &ContractFlags{}

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
func (f *ContractFlags) Contract() flags_wrapper.FlagsInterface {
	return f.FlagsInterface
}

// ContractExists returns whether a flag contract exists
func (f *ContractFlags) ContractExists() bool {
	return f.FlagsInterface != nil
}

// IsLowered determines whether the flag is lowered for a given contract.
// If a contract does not exist, it is considered to be lowered
func (f *ContractFlags) IsLowered(contractAddr common.Address) (bool, error) {
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
