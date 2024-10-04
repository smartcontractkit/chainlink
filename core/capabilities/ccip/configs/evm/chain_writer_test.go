package evm_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
)

func TestChainWriterConfigRaw(t *testing.T) {
	tests := []struct {
		name              string
		fromAddress       common.Address
		maxGasPrice       *assets.Wei
		commitGasLimit    uint64
		execBatchGasLimit uint64
		expectedError     string
	}{
		{
			name:              "valid input",
			fromAddress:       common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
			maxGasPrice:       assets.NewWeiI(1000000000),
			commitGasLimit:    21000,
			execBatchGasLimit: 42000,
			expectedError:     "",
		},
		{
			name:              "zero fromAddress",
			fromAddress:       common.HexToAddress("0x0"),
			maxGasPrice:       assets.NewWeiI(1000000000),
			commitGasLimit:    21000,
			execBatchGasLimit: 42000,
			expectedError:     "fromAddress cannot be zero",
		},
		{
			name:              "nil maxGasPrice",
			fromAddress:       common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
			maxGasPrice:       nil,
			commitGasLimit:    21000,
			execBatchGasLimit: 42000,
			expectedError:     "maxGasPrice cannot be nil",
		},
		{
			name:              "zero maxGasPrice",
			fromAddress:       common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
			maxGasPrice:       assets.NewWeiI(0),
			commitGasLimit:    21000,
			execBatchGasLimit: 42000,
			expectedError:     "maxGasPrice must be greater than zero",
		},
		{
			name:              "negative maxGasPrice",
			fromAddress:       common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
			maxGasPrice:       assets.NewWeiI(-1),
			commitGasLimit:    21000,
			execBatchGasLimit: 42000,
			expectedError:     "maxGasPrice must be greater than zero",
		},
		{
			name:              "zero commitGasLimit",
			fromAddress:       common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
			maxGasPrice:       assets.NewWeiI(1000000000),
			commitGasLimit:    0,
			execBatchGasLimit: 42000,
			expectedError:     "commitGasLimit must be greater than zero",
		},
		{
			name:              "zero execBatchGasLimit",
			fromAddress:       common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
			maxGasPrice:       assets.NewWeiI(1000000000),
			commitGasLimit:    21000,
			execBatchGasLimit: 0,
			expectedError:     "execBatchGasLimit must be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := evm.ChainWriterConfigRaw(tt.fromAddress, tt.maxGasPrice, tt.commitGasLimit, tt.execBatchGasLimit)
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
				assert.Equal(t,
					tt.maxGasPrice,
					config.MaxGasPrice)
			}
		})
	}
}
