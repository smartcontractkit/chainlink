package v1_2_0

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
)

func modifyCopy[T any](c T, f func(c *T)) T {
	f(&c)
	return c
}

func TestExecOffchainConfig120_Encoding(t *testing.T) {
	t.Parallel()
	validConfig := JSONExecOffchainConfig{
		SourceFinalityDepth:         3,
		DestOptimisticConfirmations: 6,
		DestFinalityDepth:           3,
		BatchGasLimit:               5_000_000,
		RelativeBoostPerWaitHour:    0.07,
		MaxGasPrice:                 200e9,
		InflightCacheExpiry:         *config.MustNewDuration(64 * time.Second),
		RootSnoozeTime:              *config.MustNewDuration(128 * time.Minute),
	}

	tests := []struct {
		name       string
		want       JSONExecOffchainConfig
		errPattern string
	}{
		{
			name: "legacy offchain config format parses",
			want: validConfig,
		},
		{
			name: "can omit finality depth",
			want: modifyCopy(validConfig, func(c *JSONExecOffchainConfig) {
				c.SourceFinalityDepth = 0
				c.DestFinalityDepth = 0
			}),
		},
		{
			name: "can set the DestMaxGasPrice",
			want: modifyCopy(validConfig, func(c *JSONExecOffchainConfig) {
				c.MaxGasPrice = 0
				c.DestMaxGasPrice = 200e9
			}),
		},
		{
			name: "must set DestMaxGasPrice",
			want: modifyCopy(validConfig, func(c *JSONExecOffchainConfig) {
				c.MaxGasPrice = 0
				c.DestMaxGasPrice = 0
			}),
			errPattern: "DestMaxGasPrice",
		},
		{
			name: "cannot set both MaxGasPrice and DestMaxGasPrice",
			want: modifyCopy(validConfig, func(c *JSONExecOffchainConfig) {
				c.DestMaxGasPrice = c.MaxGasPrice
			}),
			errPattern: "MaxGasPrice and DestMaxGasPrice",
		},
		{
			name: "must set BatchGasLimit",
			want: modifyCopy(validConfig, func(c *JSONExecOffchainConfig) {
				c.BatchGasLimit = 0
			}),
			errPattern: "BatchGasLimit",
		},
		{
			name: "must set DestOptimisticConfirmations",
			want: modifyCopy(validConfig, func(c *JSONExecOffchainConfig) {
				c.DestOptimisticConfirmations = 0
			}),
			errPattern: "DestOptimisticConfirmations",
		},
		{
			name: "must set RelativeBoostPerWaitHour",
			want: modifyCopy(validConfig, func(c *JSONExecOffchainConfig) {
				c.RelativeBoostPerWaitHour = 0
			}),
			errPattern: "RelativeBoostPerWaitHour",
		},
		{
			name: "must set InflightCacheExpiry",
			want: modifyCopy(validConfig, func(c *JSONExecOffchainConfig) {
				c.InflightCacheExpiry = *config.MustNewDuration(0)
			}),
			errPattern: "InflightCacheExpiry",
		},
		{
			name: "must set RootSnoozeTime",
			want: modifyCopy(validConfig, func(c *JSONExecOffchainConfig) {
				c.RootSnoozeTime = *config.MustNewDuration(0)
			}),
			errPattern: "RootSnoozeTime",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			exp := tc.want
			encode, err := ccipconfig.EncodeOffchainConfig(&exp)
			require.NoError(t, err)
			got, err := ccipconfig.DecodeOffchainConfig[JSONExecOffchainConfig](encode)

			if tc.errPattern != "" {
				require.ErrorContains(t, err, tc.errPattern)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}

func TestExecOffchainConfig120_MaxGasPrice(t *testing.T) {
	config := JSONExecOffchainConfig{
		SourceFinalityDepth:         3,
		DestOptimisticConfirmations: 6,
		DestFinalityDepth:           3,
		BatchGasLimit:               5_000_000,
		RelativeBoostPerWaitHour:    0.07,
		MaxGasPrice:                 200e9,
		InflightCacheExpiry:         *config.MustNewDuration(64 * time.Second),
		RootSnoozeTime:              *config.MustNewDuration(128 * time.Minute),
	}
	require.NoError(t, config.Validate())
	require.Equal(t, uint64(200e9), config.ComputeDestMaxGasPrice())

	config.MaxGasPrice = 0
	config.DestMaxGasPrice = 250e9
	require.NoError(t, config.Validate())
	require.Equal(t, uint64(250e9), config.ComputeDestMaxGasPrice())
}

func TestExecOffchainConfig120_ParseRawJson(t *testing.T) {
	t.Parallel()
	decoded, err := ccipconfig.DecodeOffchainConfig[JSONExecOffchainConfig]([]byte(`{
		"DestOptimisticConfirmations": 6,
		"BatchGasLimit": 5000000,
		"RelativeBoostPerWaitHour": 0.07,
		"MaxGasPrice": 200000000000,
		"InflightCacheExpiry": "64s",
		"RootSnoozeTime": "128m"
	}`))
	require.NoError(t, err)
	require.Equal(t, JSONExecOffchainConfig{
		DestOptimisticConfirmations: 6,
		BatchGasLimit:               5_000_000,
		RelativeBoostPerWaitHour:    0.07,
		MaxGasPrice:                 200e9,
		InflightCacheExpiry:         *config.MustNewDuration(64 * time.Second),
		RootSnoozeTime:              *config.MustNewDuration(128 * time.Minute),
	}, decoded)
}
