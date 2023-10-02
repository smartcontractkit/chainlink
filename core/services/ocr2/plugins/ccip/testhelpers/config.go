// Package with set of configs that should be used only within tests suites

package testhelpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var PermissionLessExecutionThresholdSeconds = uint32(FirstBlockAge.Seconds())

func (c *CCIPContracts) CreateDefaultCommitOnchainConfig(t *testing.T) []byte {
	config, err := abihelpers.EncodeAbiStruct(ccipconfig.CommitOnchainConfig{
		PriceRegistry: c.Dest.PriceRegistry.Address(),
	})
	require.NoError(t, err)
	return config
}

func (c *CCIPContracts) CreateDefaultCommitOffchainConfig(t *testing.T) []byte {
	return c.createCommitOffchainConfig(t, 10*time.Second, 5*time.Second)
}

func (c *CCIPContracts) createCommitOffchainConfig(t *testing.T, feeUpdateHearBeat time.Duration, inflightCacheExpiry time.Duration) []byte {
	config, err := ccipconfig.EncodeOffchainConfig(ccipconfig.CommitOffchainConfig{
		SourceFinalityDepth:      1,
		DestFinalityDepth:        1,
		GasPriceHeartBeat:        models.MustMakeDuration(feeUpdateHearBeat),
		DAGasPriceDeviationPPB:   1,
		ExecGasPriceDeviationPPB: 1,
		TokenPriceHeartBeat:      models.MustMakeDuration(feeUpdateHearBeat),
		TokenPriceDeviationPPB:   1,
		MaxGasPrice:              200e9,
		InflightCacheExpiry:      models.MustMakeDuration(inflightCacheExpiry),
	})
	require.NoError(t, err)
	return config
}

func (c *CCIPContracts) CreateDefaultExecOnchainConfig(t *testing.T) []byte {
	config, err := abihelpers.EncodeAbiStruct(ccipconfig.ExecOnchainConfig{
		PermissionLessExecutionThresholdSeconds: PermissionLessExecutionThresholdSeconds,
		Router:                                  c.Dest.Router.Address(),
		PriceRegistry:                           c.Dest.PriceRegistry.Address(),
		MaxDataSize:                             1e5,
		MaxTokensLength:                         5,
	})
	require.NoError(t, err)
	return config
}

func (c *CCIPContracts) CreateDefaultExecOffchainConfig(t *testing.T) []byte {
	return c.createExecOffchainConfig(t, 1*time.Minute, 1*time.Minute)
}

func (c *CCIPContracts) createExecOffchainConfig(t *testing.T, inflightCacheExpiry time.Duration, rootSnoozeTime time.Duration) []byte {
	config, err := ccipconfig.EncodeOffchainConfig(ccipconfig.ExecOffchainConfig{
		SourceFinalityDepth:         1,
		DestOptimisticConfirmations: 1,
		DestFinalityDepth:           1,
		BatchGasLimit:               5_000_000,
		RelativeBoostPerWaitHour:    0.07,
		MaxGasPrice:                 200e9,
		InflightCacheExpiry:         models.MustMakeDuration(inflightCacheExpiry),
		RootSnoozeTime:              models.MustMakeDuration(rootSnoozeTime),
	})
	require.NoError(t, err)
	return config
}
