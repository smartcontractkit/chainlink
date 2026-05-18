package v1_0_0

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib/rpclibmocks"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
)

func TestExecOffchainConfig100_Encoding(t *testing.T) {
	tests := []struct {
		name      string
		want      ExecOffchainConfig
		expectErr bool
	}{
		{
			name: "encodes and decodes config with all fields set",
			want: ExecOffchainConfig{
				SourceFinalityDepth:         3,
				DestOptimisticConfirmations: 6,
				DestFinalityDepth:           3,
				BatchGasLimit:               5_000_000,
				RelativeBoostPerWaitHour:    0.07,
				InflightCacheExpiry:         *config.MustNewDuration(64 * time.Second),
				RootSnoozeTime:              *config.MustNewDuration(128 * time.Minute),
				MessageVisibilityInterval:   *config.MustNewDuration(6 * time.Hour),
			},
		},
		{
			name: "fails decoding when all fields present but with 0 values",
			want: ExecOffchainConfig{
				SourceFinalityDepth:         0,
				DestFinalityDepth:           0,
				DestOptimisticConfirmations: 0,
				BatchGasLimit:               0,
				RelativeBoostPerWaitHour:    0,
				InflightCacheExpiry:         *config.MustNewDuration(0),
				RootSnoozeTime:              *config.MustNewDuration(0),
				MessageVisibilityInterval:   *config.MustNewDuration(0),
			},
			expectErr: true,
		},
		{
			name:      "fails decoding when all fields are missing",
			want:      ExecOffchainConfig{},
			expectErr: true,
		},
		{
			name: "fails decoding when some fields are missing",
			want: ExecOffchainConfig{
				SourceFinalityDepth: 99999999,
				InflightCacheExpiry: *config.MustNewDuration(64 * time.Second),
			},
			expectErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			exp := tc.want
			encode, err := ccipconfig.EncodeOffchainConfig(&exp)
			require.NoError(t, err)
			got, err := ccipconfig.DecodeOffchainConfig[ExecOffchainConfig](encode)

			if tc.expectErr {
				require.ErrorContains(t, err, "must set")
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}

func TestExecOffchainConfig100_AllFieldsRequired(t *testing.T) {
	cfg := ExecOffchainConfig{
		SourceFinalityDepth:         3,
		DestOptimisticConfirmations: 6,
		DestFinalityDepth:           3,
		BatchGasLimit:               5_000_000,
		RelativeBoostPerWaitHour:    0.07,
		InflightCacheExpiry:         *config.MustNewDuration(64 * time.Second),
		RootSnoozeTime:              *config.MustNewDuration(128 * time.Minute),
		BatchingStrategyID:          0,
	}
	encoded, err := ccipconfig.EncodeOffchainConfig(&cfg)
	require.NoError(t, err)

	var configAsMap map[string]any
	err = json.Unmarshal(encoded, &configAsMap)
	require.NoError(t, err)
	for keyToDelete := range configAsMap {
		if keyToDelete == "MessageVisibilityInterval" {
			continue // this field is optional
		}

		partialConfig := make(map[string]any)
		for k, v := range configAsMap {
			if k != keyToDelete {
				partialConfig[k] = v
			}
		}
		encodedPartialConfig, err := json.Marshal(partialConfig)
		require.NoError(t, err)
		_, err = ccipconfig.DecodeOffchainConfig[ExecOffchainConfig](encodedPartialConfig)
		if keyToDelete == "BatchingStrategyID" {
			require.NoError(t, err)
		} else {
			require.ErrorContains(t, err, keyToDelete)
		}
	}
}

func Test_GetSendersNonce(t *testing.T) {
	sender1 := cciptypes.Address(utils.RandomAddress().String())
	sender2 := cciptypes.Address(utils.RandomAddress().String())

	tests := []struct {
		name           string
		addresses      []cciptypes.Address
		batchCaller    *rpclibmocks.EvmBatchCaller
		expectedResult map[cciptypes.Address]uint64
		expectedError  bool
	}{
		{
			name:           "return empty map when input is empty",
			addresses:      []cciptypes.Address{},
			batchCaller:    rpclibmocks.NewEvmBatchCaller(t),
			expectedResult: map[cciptypes.Address]uint64{},
		},
		{
			name:      "return error when batch call fails",
			addresses: []cciptypes.Address{sender1},
			batchCaller: func() *rpclibmocks.EvmBatchCaller {
				mockBatchCaller := rpclibmocks.NewEvmBatchCaller(t)
				mockBatchCaller.On("BatchCall", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("batch call error"))
				return mockBatchCaller
			}(),
			expectedError: true,
		},
		{
			name:      "return error when nonces dont match senders",
			addresses: []cciptypes.Address{sender1, sender2},
			batchCaller: func() *rpclibmocks.EvmBatchCaller {
				mockBatchCaller := rpclibmocks.NewEvmBatchCaller(t)
				results := []rpclib.DataAndErr{
					{
						Outputs: []any{uint64(1)},
						Err:     nil,
					},
				}
				mockBatchCaller.On("BatchCall", mock.Anything, mock.Anything, mock.Anything).
					Return(results, nil)
				return mockBatchCaller
			}(),
			expectedError: true,
		},
		{
			name:      "return error when single request from batch fails",
			addresses: []cciptypes.Address{sender1, sender2},
			batchCaller: func() *rpclibmocks.EvmBatchCaller {
				mockBatchCaller := rpclibmocks.NewEvmBatchCaller(t)
				results := []rpclib.DataAndErr{
					{
						Outputs: []any{uint64(1)},
						Err:     nil,
					},
					{
						Outputs: []any{},
						Err:     errors.New("request failed"),
					},
				}
				mockBatchCaller.On("BatchCall", mock.Anything, mock.Anything, mock.Anything).
					Return(results, nil)
				return mockBatchCaller
			}(),
			expectedError: true,
		},
		{
			name:      "return map of nonce per sender",
			addresses: []cciptypes.Address{sender1, sender2},
			batchCaller: func() *rpclibmocks.EvmBatchCaller {
				mockBatchCaller := rpclibmocks.NewEvmBatchCaller(t)
				results := []rpclib.DataAndErr{
					{
						Outputs: []any{uint64(1)},
						Err:     nil,
					},
					{
						Outputs: []any{uint64(2)},
						Err:     nil,
					},
				}
				mockBatchCaller.On("BatchCall", mock.Anything, mock.Anything, mock.Anything).
					Return(results, nil)
				return mockBatchCaller
			}(),
			expectedResult: map[cciptypes.Address]uint64{
				sender1: uint64(1),
				sender2: uint64(2),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			offramp := OffRamp{evmBatchCaller: test.batchCaller, Logger: logger.TestLogger(t)}
			nonce, err := offramp.ListSenderNonces(testutils.Context(t), test.addresses)

			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectedResult, nonce)
			}
		})
	}
}
