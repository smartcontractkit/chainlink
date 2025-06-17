package v1_4_0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

func TestTokenPool(t *testing.T) {
	addr := utils.RandomAddress()
	chainSelector := uint64(2000)
	poolType := "BurnMint"

	tokenPool := NewTokenPool(poolType, addr, chainSelector)

	assert.Equal(t, addr, tokenPool.Address())
	assert.Equal(t, poolType, tokenPool.Type())

	inboundRateLimitCall := GetInboundTokenPoolRateLimitCall(addr, chainSelector)

	assert.Equal(t, "getCurrentInboundRateLimiterState", inboundRateLimitCall.MethodName())
}
