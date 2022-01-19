package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra"

	"github.com/smartcontractkit/chainlink/core/logger"
)

func TestConfig(t *testing.T) {
	def := terra.DefaultConfigSet

	cfg := NewConfig(ChainCfg{}, logger.TestLogger(t))
	assert.Equal(t, def.BlocksUntilTxTimeout, cfg.BlocksUntilTxTimeout())
	assert.Equal(t, def.ConfirmMaxPolls, cfg.ConfirmMaxPolls())
	assert.Equal(t, def.ConfirmPollPeriod, cfg.ConfirmPollPeriod())
	assert.Equal(t, def.FallbackGasPriceULuna, cfg.FallbackGasPriceULuna())
	assert.Equal(t, def.GasLimitMultiplier, cfg.GasLimitMultiplier())
	assert.Equal(t, def.MaxMsgsPerBatch, cfg.MaxMsgsPerBatch())

	updated := ChainCfg{
		BlocksUntilTxTimeout:  null.IntFrom(1000),
		FallbackGasPriceULuna: null.StringFrom("5.6"),
	}
	cfg.Update(updated)
	assert.Equal(t, updated.BlocksUntilTxTimeout.Int64, cfg.BlocksUntilTxTimeout())
	assert.Equal(t, def.ConfirmMaxPolls, cfg.ConfirmMaxPolls())
	assert.Equal(t, def.ConfirmPollPeriod, cfg.ConfirmPollPeriod())
	assert.Equal(t, sdk.MustNewDecFromStr(updated.FallbackGasPriceULuna.String), cfg.FallbackGasPriceULuna())
	assert.Equal(t, def.GasLimitMultiplier, cfg.GasLimitMultiplier())
	assert.Equal(t, def.MaxMsgsPerBatch, cfg.MaxMsgsPerBatch())

	updated = ChainCfg{
		FallbackGasPriceULuna: null.StringFrom("not-a-number"),
	}
	cfg.Update(updated)
	assert.Equal(t, def.FallbackGasPriceULuna, cfg.FallbackGasPriceULuna())
}
