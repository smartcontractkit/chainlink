package evm_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
)

func TestChainWriterConfigRaw(t *testing.T) {
	tests := []struct {
		name              string
		fromAddress       common.Address
		commitGasLimit    uint64
		execBatchGasLimit uint64
		expectedError     string
	}{
		{
			name:              "valid input",
			fromAddress:       common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
			commitGasLimit:    21000,
			execBatchGasLimit: 42000,
			expectedError:     "",
		},
		{
			name:              "zero fromAddress",
			fromAddress:       common.HexToAddress("0x0"),
			commitGasLimit:    21000,
			execBatchGasLimit: 42000,
			expectedError:     "fromAddress cannot be zero",
		},
		{
			name:              "zero commitGasLimit",
			fromAddress:       common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
			commitGasLimit:    0,
			execBatchGasLimit: 42000,
			expectedError:     "commitGasLimit must be greater than zero",
		},
		{
			name:              "zero execBatchGasLimit",
			fromAddress:       common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
			commitGasLimit:    21000,
			execBatchGasLimit: 0,
			expectedError:     "execBatchGasLimit must be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := evm.ChainWriterConfigRaw(tt.fromAddress, tt.commitGasLimit, tt.execBatchGasLimit)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t,
					tt.fromAddress,
					config.Contracts[consts.ContractNameOffRamp].Configs[consts.MethodCommit].FromAddress)
				assert.Equal(t,
					tt.commitGasLimit,
					config.Contracts[consts.ContractNameOffRamp].Configs[consts.MethodCommit].GasLimit)
				assert.Equal(t,
					tt.execBatchGasLimit,
					config.Contracts[consts.ContractNameOffRamp].Configs[consts.MethodExecute].GasLimit)
			}
		})
	}
}
