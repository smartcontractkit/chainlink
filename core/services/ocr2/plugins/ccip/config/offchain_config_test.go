package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func TestCommitOffchainConfig_Encoding(t *testing.T) {
	tests := map[string]struct {
		want      CommitOffchainConfig
		expectErr bool
	}{
		"encodes and decodes config with all fields set": {
			want: CommitOffchainConfig{
				SourceFinalityDepth:      3,
				DestFinalityDepth:        3,
				GasPriceHeartBeat:        models.MustMakeDuration(1 * time.Hour),
				DAGasPriceDeviationPPB:   5e7,
				ExecGasPriceDeviationPPB: 5e7,
				TokenPriceHeartBeat:      models.MustMakeDuration(1 * time.Hour),
				TokenPriceDeviationPPB:   5e7,
				MaxGasPrice:              200e9,
				InflightCacheExpiry:      models.MustMakeDuration(23456 * time.Second),
			},
		},
		"fails decoding when all fields present but with 0 values": {
			want: CommitOffchainConfig{
				SourceFinalityDepth:      0,
				DestFinalityDepth:        0,
				GasPriceHeartBeat:        models.MustMakeDuration(0),
				DAGasPriceDeviationPPB:   0,
				ExecGasPriceDeviationPPB: 0,
				TokenPriceHeartBeat:      models.MustMakeDuration(0),
				TokenPriceDeviationPPB:   0,
				MaxGasPrice:              0,
				InflightCacheExpiry:      models.MustMakeDuration(0),
			},
			expectErr: true,
		},
		"fails decoding when all fields are missing": {
			want:      CommitOffchainConfig{},
			expectErr: true,
		},
		"fails decoding when some fields are missing": {
			want: CommitOffchainConfig{
				SourceFinalityDepth:      3,
				GasPriceHeartBeat:        models.MustMakeDuration(1 * time.Hour),
				DAGasPriceDeviationPPB:   5e7,
				ExecGasPriceDeviationPPB: 5e7,
				TokenPriceHeartBeat:      models.MustMakeDuration(1 * time.Hour),
				TokenPriceDeviationPPB:   5e7,
				MaxGasPrice:              200e9,
			},
			expectErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			encode, err := EncodeOffchainConfig(tc.want)
			require.NoError(t, err)
			got, err := DecodeOffchainConfig[CommitOffchainConfig](encode)

			if tc.expectErr {
				require.ErrorContains(t, err, "must set")
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}

func TestExecOffchainConfig_Encoding(t *testing.T) {
	tests := map[string]struct {
		want      ExecOffchainConfig
		expectErr bool
	}{
		"encodes and decodes config with all fields set": {
			want: ExecOffchainConfig{
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
			want: ExecOffchainConfig{
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
			want:      ExecOffchainConfig{},
			expectErr: true,
		},
		"fails decoding when some fields are missing": {
			want: ExecOffchainConfig{
				SourceFinalityDepth: 99999999,
				InflightCacheExpiry: models.MustMakeDuration(64 * time.Second),
			},
			expectErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			exp := tc.want
			encode, err := EncodeOffchainConfig(&exp)
			require.NoError(t, err)
			got, err := DecodeOffchainConfig[ExecOffchainConfig](encode)

			if tc.expectErr {
				require.ErrorContains(t, err, "must set")
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}
