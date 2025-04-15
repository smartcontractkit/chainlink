package globals

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
)

func TestWithCommitOffchainOverrides(t *testing.T) {
	baseConfig := pluginconfig.CommitOffchainConfig{
		ChainFeeAsyncObserverSyncFreq:      1 * time.Second,
		ChainFeeAsyncObserverSyncTimeout:   2 * time.Second,
		TokenPriceAsyncObserverSyncFreq:    *config.MustNewDuration(3 * time.Second),
		TokenPriceAsyncObserverSyncTimeout: *config.MustNewDuration(4 * time.Second),
		TransmissionDelayMultiplier:        5 * time.Second,
		RMNEnabled:                         false,
		MaxReportTransmissionCheckAttempts: 5,
	}

	overrideConfig := pluginconfig.CommitOffchainConfig{
		// Only override specific fields
		ChainFeeAsyncObserverSyncFreq:      10 * time.Second,
		ChainFeeAsyncObserverSyncTimeout:   20 * time.Second,
		TokenPriceAsyncObserverSyncFreq:    *config.MustNewDuration(30 * time.Second),
		TokenPriceAsyncObserverSyncTimeout: *config.MustNewDuration(40 * time.Second),
		TransmissionDelayMultiplier:        50 * time.Second,
	}

	// Act
	result := withCommitOffchainOverrides(baseConfig, overrideConfig)

	// Assert
	// Overridden fields should have new values
	assert.Equal(t, 10*time.Second, result.ChainFeeAsyncObserverSyncFreq)
	assert.Equal(t, 20*time.Second, result.ChainFeeAsyncObserverSyncTimeout)
	assert.Equal(t, *config.MustNewDuration(30 * time.Second), result.TokenPriceAsyncObserverSyncFreq)
	assert.Equal(t, *config.MustNewDuration(40 * time.Second), result.TokenPriceAsyncObserverSyncTimeout)
	assert.Equal(t, 50*time.Second, result.TransmissionDelayMultiplier)

	// Non-overridden fields should retain base values
	assert.False(t, result.RMNEnabled)
	assert.Equal(t, uint(5), result.MaxReportTransmissionCheckAttempts)
}

func TestWithCommitOffchainOverrides_EmptyOverrides(t *testing.T) {
	baseConfig := pluginconfig.CommitOffchainConfig{
		ChainFeeAsyncObserverSyncFreq:      1 * time.Second,
		ChainFeeAsyncObserverSyncTimeout:   2 * time.Second,
		TokenPriceAsyncObserverSyncFreq:    *config.MustNewDuration(3 * time.Second),
		TokenPriceAsyncObserverSyncTimeout: *config.MustNewDuration(4 * time.Second),
		TransmissionDelayMultiplier:        5 * time.Second,
		RMNEnabled:                         true,
	}

	emptyOverrides := pluginconfig.CommitOffchainConfig{}

	// Act
	result := withCommitOffchainOverrides(baseConfig, emptyOverrides)

	// Assert - check each field individually
	assert.Equal(t, baseConfig.ChainFeeAsyncObserverSyncFreq, result.ChainFeeAsyncObserverSyncFreq)
	assert.Equal(t, baseConfig.ChainFeeAsyncObserverSyncTimeout, result.ChainFeeAsyncObserverSyncTimeout)
	assert.Equal(t, baseConfig.TokenPriceAsyncObserverSyncFreq.Duration(), result.TokenPriceAsyncObserverSyncFreq.Duration())
	assert.Equal(t, baseConfig.TokenPriceAsyncObserverSyncTimeout.Duration(), result.TokenPriceAsyncObserverSyncTimeout.Duration())
	assert.Equal(t, baseConfig.TransmissionDelayMultiplier, result.TransmissionDelayMultiplier)
	assert.Equal(t, baseConfig.RMNEnabled, result.RMNEnabled)
}

func TestCommitOffChainCfgForEthereum(t *testing.T) {
	// This tests that the Ethereum-specific config correctly overrides the default config

	// Assert that the Ethereum config has different values for the overridden fields
	assert.NotEqual(t, DefaultCommitOffChainCfg.ChainFeeAsyncObserverSyncFreq,
		CommitOffChainCfgForEthereum.ChainFeeAsyncObserverSyncFreq)
	assert.NotEqual(t, DefaultCommitOffChainCfg.ChainFeeAsyncObserverSyncTimeout,
		CommitOffChainCfgForEthereum.ChainFeeAsyncObserverSyncTimeout)
	assert.NotEqual(t, DefaultCommitOffChainCfg.TokenPriceAsyncObserverSyncFreq,
		CommitOffChainCfgForEthereum.TokenPriceAsyncObserverSyncFreq)
	assert.NotEqual(t, DefaultCommitOffChainCfg.TokenPriceAsyncObserverSyncTimeout,
		CommitOffChainCfgForEthereum.TokenPriceAsyncObserverSyncTimeout)

	// Check the expected values
	assert.Equal(t, 4*time.Second, CommitOffChainCfgForEthereum.ChainFeeAsyncObserverSyncFreq)
	assert.Equal(t, 3*time.Second, CommitOffChainCfgForEthereum.ChainFeeAsyncObserverSyncTimeout)
	assert.Equal(t, *config.MustNewDuration(4 * time.Second), CommitOffChainCfgForEthereum.TokenPriceAsyncObserverSyncFreq)
	assert.Equal(t, *config.MustNewDuration(3 * time.Second), CommitOffChainCfgForEthereum.TokenPriceAsyncObserverSyncTimeout)

	// Non-overridden fields should remain the same
	assert.Equal(t, DefaultCommitOffChainCfg.RMNEnabled, CommitOffChainCfgForEthereum.RMNEnabled)
	assert.Equal(t, DefaultCommitOffChainCfg.MaxReportTransmissionCheckAttempts,
		CommitOffChainCfgForEthereum.MaxReportTransmissionCheckAttempts)
}
