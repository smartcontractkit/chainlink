// Package with set of configs that should be used only within tests suites

package testhelpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

var PermissionLessExecutionThresholdSeconds = uint32(FirstBlockAge.Seconds())

func (c *CCIPContracts) CreateDefaultCommitOnchainConfig(t *testing.T) []byte {
	config, err := abihelpers.EncodeAbiStruct(ccipdata.CommitOnchainConfig{
		PriceRegistry: c.Dest.PriceRegistry.Address(),
	})
	require.NoError(t, err)
	return config
}

func (c *CCIPContracts) CreateDefaultCommitOffchainConfig(t *testing.T) []byte {
	return c.createCommitOffchainConfig(t, 10*time.Second, 5*time.Second)
}

func (c *CCIPContracts) createCommitOffchainConfig(t *testing.T, feeUpdateHearBeat time.Duration, inflightCacheExpiry time.Duration) []byte {
	config, err := ccipconfig.EncodeOffchainConfig(v1_2_0.CommitOffchainConfig{
		SourceFinalityDepth:      1,
		DestFinalityDepth:        1,
		GasPriceHeartBeat:        *config.MustNewDuration(feeUpdateHearBeat),
		DAGasPriceDeviationPPB:   1,
		ExecGasPriceDeviationPPB: 1,
		TokenPriceHeartBeat:      *config.MustNewDuration(feeUpdateHearBeat),
		TokenPriceDeviationPPB:   1,
		MaxGasPrice:              200e9,
		InflightCacheExpiry:      *config.MustNewDuration(inflightCacheExpiry),
	})
	require.NoError(t, err)
	return config
}

func (c *CCIPContracts) CreateDefaultExecOnchainConfig(t *testing.T) []byte {
	config, err := abihelpers.EncodeAbiStruct(v1_2_0.ExecOnchainConfig{
		PermissionLessExecutionThresholdSeconds: PermissionLessExecutionThresholdSeconds,
		Router:                                  c.Dest.Router.Address(),
		PriceRegistry:                           c.Dest.PriceRegistry.Address(),
		MaxDataBytes:                            1e5,
		MaxNumberOfTokensPerMsg:                 5,
		MaxPoolReleaseOrMintGas:                 200_000,
	})
	require.NoError(t, err)
	return config
}

func (c *CCIPContracts) CreateDefaultExecOffchainConfig(t *testing.T) []byte {
	return c.createExecOffchainConfig(t, 1*time.Minute, 1*time.Minute)
}

func (c *CCIPContracts) createExecOffchainConfig(t *testing.T, inflightCacheExpiry time.Duration, rootSnoozeTime time.Duration) []byte {
	config, err := ccipconfig.EncodeOffchainConfig(v1_0_0.ExecOffchainConfig{
		SourceFinalityDepth:         1,
		DestOptimisticConfirmations: 1,
		DestFinalityDepth:           1,
		BatchGasLimit:               5_000_000,
		RelativeBoostPerWaitHour:    0.07,
		MaxGasPrice:                 200e9,
		InflightCacheExpiry:         *config.MustNewDuration(inflightCacheExpiry),
		RootSnoozeTime:              *config.MustNewDuration(rootSnoozeTime),
	})
	require.NoError(t, err)
	return config
}
