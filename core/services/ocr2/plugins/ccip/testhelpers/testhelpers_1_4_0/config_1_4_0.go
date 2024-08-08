// Package with set of configs that should be used only within tests suites

package testhelpers_1_4_0

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
)

var PermissionLessExecutionThresholdSeconds = uint32(testhelpers.FirstBlockAge.Seconds())

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
	config, err := NewCommitOffchainConfig(
		*config.MustNewDuration(feeUpdateHearBeat),
		1,
		1,
		*config.MustNewDuration(feeUpdateHearBeat),
		1,
		*config.MustNewDuration(inflightCacheExpiry),
	).Encode()
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
	config, err := NewExecOffchainConfig(
		1,
		5_000_000,
		0.07,
		*config.MustNewDuration(inflightCacheExpiry),
		*config.MustNewDuration(rootSnoozeTime),
	).Encode()
	require.NoError(t, err)
	return config
}
