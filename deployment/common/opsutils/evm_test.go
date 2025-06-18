package opsutils

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/stretchr/testify/assert"
)

func TestCloneTransactOptsWithGas(t *testing.T) {
	orig := &bind.TransactOpts{
		GasLimit: 100,
		GasPrice: big.NewInt(123),
	}
	// Should clone and override both
	cloned := cloneTransactOptsWithGas(orig, 200, 456)
	assert.NotSame(t, orig, cloned)
	assert.Equal(t, uint64(200), cloned.GasLimit)
	assert.Equal(t, big.NewInt(456), cloned.GasPrice)
	// Should not override if zero
	cloned2 := cloneTransactOptsWithGas(orig, 0, 0)
	assert.Equal(t, orig.GasLimit, cloned2.GasLimit)
	assert.Equal(t, orig.GasPrice, cloned2.GasPrice)
	// Nil input
	assert.Nil(t, cloneTransactOptsWithGas(nil, 1, 1))
}

func TestGasBoostConfigsForChainMap(t *testing.T) {
	chainMap := map[uint64]string{1: "a", 2: "b"}
	gasBoostConfigs := map[uint64]commontypes.GasBoostConfig{
		1: {InitialGasLimit: 10},
	}
	cfgs := GasBoostConfigsForChainMap(chainMap, gasBoostConfigs)
	assert.Len(t, cfgs, 2)
	assert.NotNil(t, cfgs[1])
	assert.Nil(t, cfgs[2])
	// Nil configs
	assert.Empty(t, GasBoostConfigsForChainMap[string](chainMap, nil))
	assert.Empty(t, GasBoostConfigsForChainMap[string](nil, gasBoostConfigs))
}

func TestGetBoostedGasForAttempt_DefaultsAndOverrides(t *testing.T) {
	cfg := commontypes.GasBoostConfig{}
	limit, price := getBoostedGasForAttempt(cfg, 0)
	assert.Equal(t, uint64(200_000), limit)
	assert.Equal(t, uint64(20_000_000_000), price)
	limit, price = getBoostedGasForAttempt(cfg, 2)
	assert.Equal(t, uint64(200_000+2*50_000), limit)
	assert.Equal(t, uint64(20_000_000_000+2*10_000_000_000), price)

	cfg = commontypes.GasBoostConfig{
		InitialGasLimit:   1000,
		GasLimitIncrement: 100,
		InitialGasPrice:   2000,
		GasPriceIncrement: 100,
	}
	limit, price = getBoostedGasForAttempt(cfg, 3)
	assert.Equal(t, uint64(1000+3*100), limit)
	assert.Equal(t, uint64(2000+3*100), price)
}

func TestRetryDeploymentWithGasBoost(t *testing.T) {
	cfg := &commontypes.GasBoostConfig{
		InitialGasLimit:   1000,
		GasLimitIncrement: 100,
		InitialGasPrice:   2000,
		GasPriceIncrement: 100,
	}
	opt := RetryDeploymentWithGasBoost[any](cfg)
	// Should not panic and should be non-nil
	assert.NotNil(t, opt)
	// Should fallback to default if nil
	assert.NotNil(t, RetryDeploymentWithGasBoost[string](nil))
}
