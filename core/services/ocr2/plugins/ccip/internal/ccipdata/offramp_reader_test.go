package ccipdata_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestOffRampFilters(t *testing.T) {
	assertFilterRegistration(t, new(lpmocks.LogPoller), func(lp *lpmocks.LogPoller, addr common.Address) ccipdata.Closer {
		c, err := ccipdata.NewOffRampV1_0_0(logger.TestLogger(t), addr, new(mocks.Client), lp, nil)
		require.NoError(t, err)
		return c
	}, 3)
	assertFilterRegistration(t, new(lpmocks.LogPoller), func(lp *lpmocks.LogPoller, addr common.Address) ccipdata.Closer {
		c, err := ccipdata.NewOffRampV1_2_0(logger.TestLogger(t), addr, new(mocks.Client), lp, nil)
		require.NoError(t, err)
		return c
	}, 3)
}

func TestExecOffchainConfig_Encoding(t *testing.T) {
	tests := map[string]struct {
		want      ccipdata.ExecOffchainConfig
		expectErr bool
	}{
		"encodes and decodes config with all fields set": {
			want: ccipdata.ExecOffchainConfig{
				SourceFinalityDepth:         3,
				DestOptimisticConfirmations: 6,
				DestFinalityDepth:           3,
				BatchGasLimit:               5_000_000,
				RelativeBoostPerWaitHour:    0.07,
				MaxGasPrice:                 200e9,
				InflightCacheExpiry:         models.MustMakeDuration(64 * time.Second),
				RootSnoozeTime:              models.MustMakeDuration(128 * time.Minute),
			},
		},
		"fails decoding when all fields present but with 0 values": {
			want: ccipdata.ExecOffchainConfig{
				SourceFinalityDepth:         0,
				DestFinalityDepth:           0,
				DestOptimisticConfirmations: 0,
				BatchGasLimit:               0,
				RelativeBoostPerWaitHour:    0,
				MaxGasPrice:                 0,
				InflightCacheExpiry:         models.MustMakeDuration(0),
				RootSnoozeTime:              models.MustMakeDuration(0),
			},
			expectErr: true,
		},
		"fails decoding when all fields are missing": {
			want:      ccipdata.ExecOffchainConfig{},
			expectErr: true,
		},
		"fails decoding when some fields are missing": {
			want: ccipdata.ExecOffchainConfig{
				SourceFinalityDepth: 99999999,
				InflightCacheExpiry: models.MustMakeDuration(64 * time.Second),
			},
			expectErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			exp := tc.want
			encode, err := ccipconfig.EncodeOffchainConfig(&exp)
			require.NoError(t, err)
			got, err := ccipconfig.DecodeOffchainConfig[ccipdata.ExecOffchainConfig](encode)

			if tc.expectErr {
				require.ErrorContains(t, err, "must set")
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}

func TestExecOnchainConfig100(t *testing.T) {
	tests := []struct {
		name      string
		want      ccipdata.ExecOnchainConfigV1_0_0
		expectErr bool
	}{
		{
			name: "encodes and decodes config with all fields set",
			want: ccipdata.ExecOnchainConfigV1_0_0{
				PermissionLessExecutionThresholdSeconds: rand.Uint32(),
				Router:                                  utils.RandomAddress(),
				PriceRegistry:                           utils.RandomAddress(),
				MaxTokensLength:                         uint16(rand.Uint32()),
				MaxDataSize:                             rand.Uint32(),
			},
		},
		{
			name: "encodes and fails decoding config with missing fields",
			want: ccipdata.ExecOnchainConfigV1_0_0{
				PermissionLessExecutionThresholdSeconds: rand.Uint32(),
				MaxDataSize:                             rand.Uint32(),
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := abihelpers.EncodeAbiStruct(tt.want)
			require.NoError(t, err)

			decoded, err := abihelpers.DecodeAbiStruct[ccipdata.ExecOnchainConfigV1_0_0](encoded)
			if tt.expectErr {
				require.ErrorContains(t, err, "must set")
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, decoded)
			}
		})
	}
}

func TestExecOnchainConfig120(t *testing.T) {
	tests := []struct {
		name      string
		want      ccipdata.ExecOnchainConfigV1_2_0
		expectErr bool
	}{
		{
			name: "encodes and decodes config with all fields set",
			want: ccipdata.ExecOnchainConfigV1_2_0{
				PermissionLessExecutionThresholdSeconds: rand.Uint32(),
				Router:                                  utils.RandomAddress(),
				PriceRegistry:                           utils.RandomAddress(),
				MaxNumberOfTokensPerMsg:                 uint16(rand.Uint32()),
				MaxDataBytes:                            rand.Uint32(),
				MaxPoolReleaseOrMintGas:                 rand.Uint32(),
			},
		},
		{
			name: "encodes and fails decoding config with missing fields",
			want: ccipdata.ExecOnchainConfigV1_2_0{
				PermissionLessExecutionThresholdSeconds: rand.Uint32(),
				MaxDataBytes:                            rand.Uint32(),
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := abihelpers.EncodeAbiStruct(tt.want)
			require.NoError(t, err)

			decoded, err := abihelpers.DecodeAbiStruct[ccipdata.ExecOnchainConfigV1_2_0](encoded)
			if tt.expectErr {
				require.ErrorContains(t, err, "must set")
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, decoded)
			}
		})
	}
}
