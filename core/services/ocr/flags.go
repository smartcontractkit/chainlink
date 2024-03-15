package ocr

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/flags_wrapper"
)

// ContractFlags wraps the a contract
type ContractFlags struct {
	flags_wrapper.FlagsInterface
}

// NewFlags constructs a new Flags from a flags contract address
func NewFlags(addrHex string, ethClient evmclient.Client) (*ContractFlags, error) {
	flags := &ContractFlags{}

	if addrHex == "" {
		return flags, nil
	}

	contractAddr := common.HexToAddress(addrHex)
	contract, err := flags_wrapper.NewFlags(contractAddr, ethClient)
	if err != nil {
		return flags, errors.Wrap(err, "Failed to create flags wrapper")
	}
	flags.FlagsInterface = contract
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
		return true, errors.Wrap(err, "Failed to call GetFlags in the contract")
	}

	return !flags[0] || !flags[1], nil
}
