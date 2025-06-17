package v1_2_0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

func TestTokenPool(t *testing.T) {
	addr := utils.RandomAddress()
	offRamp := utils.RandomAddress()
	poolType := "BurnMint"

	tokenPool := NewTokenPool(poolType, addr, offRamp)

	assert.Equal(t, addr, tokenPool.Address())
	assert.Equal(t, poolType, tokenPool.Type())

	inboundRateLimitCall := GetInboundTokenPoolRateLimitCall(addr, offRamp)

	assert.Equal(t, "currentOffRampRateLimiterState", inboundRateLimitCall.MethodName())
}
