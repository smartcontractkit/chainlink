package evm

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"

	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
)

var (
	offrampABI = evmtypes.MustGetABI(offramp.OffRampABI)
)

// ChainWriterConfigRaw returns a ChainWriterConfig that can be used to transmit commit and execute reports.
func ChainWriterConfigRaw(
	fromAddress common.Address,
	commitGasLimit,
	execBatchGasLimit uint64,
) (config.ChainWriterConfig, error) {
	if fromAddress == common.HexToAddress("0x0") {
		return config.ChainWriterConfig{}, errors.New("fromAddress cannot be zero")
	}
	if commitGasLimit == 0 {
		return config.ChainWriterConfig{}, errors.New("commitGasLimit must be greater than zero")
	}
	if execBatchGasLimit == 0 {
		return config.ChainWriterConfig{}, errors.New("execBatchGasLimit must be greater than zero")
	}

	return config.ChainWriterConfig{
		Contracts: map[string]*config.ContractConfig{
			consts.ContractNameOffRamp: {
				ContractABI: offramp.OffRampABI,
				Configs: map[string]*config.ChainWriterDefinition{
					consts.MethodCommit: {
						ChainSpecificName: mustGetMethodName("commit", offrampABI),
						FromAddress:       fromAddress,
						GasLimit:          commitGasLimit,
					},
					consts.MethodExecute: {
						ChainSpecificName: mustGetMethodName("execute", offrampABI),
						FromAddress:       fromAddress,
						GasLimit:          execBatchGasLimit,
					},
				},
			},
		},
	}, nil
}

// mustGetMethodName panics if the method name is not found in the provided ABI.
func mustGetMethodName(name string, tabi abi.ABI) (methodName string) {
	m, ok := tabi.Methods[name]
	if !ok {
		panic(fmt.Sprintf("missing method %s in the abi", name))
	}
	return m.Name
}
