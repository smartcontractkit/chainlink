package cache

import (
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestNewTokenPools(t *testing.T) {
	ctx := testutils.Context(t)

	tk1src := utils.RandomAddress()
	tk1dst := utils.RandomAddress()
	tk1pool := utils.RandomAddress()

	tk2src := utils.RandomAddress()
	tk2dst := utils.RandomAddress()
	tk2pool := utils.RandomAddress()

	testCases := []struct {
		name               string
		sourceToDestTokens map[common.Address]common.Address // offramp
		feeTokens          []common.Address                  // price registry
		tokenToPool        map[common.Address]common.Address // offramp
		expRes             map[common.Address]common.Address
		expErr             bool
	}{
		{
			name:               "no tokens",
			sourceToDestTokens: map[common.Address]common.Address{},
			feeTokens:          []common.Address{},
			tokenToPool:        map[common.Address]common.Address{},
			expRes:             map[common.Address]common.Address{},
			expErr:             false,
		},
		{
			name: "happy flow",
			sourceToDestTokens: map[common.Address]common.Address{
				tk1src: tk1dst,
				tk2src: tk2dst,
			},
			feeTokens: []common.Address{tk1dst, tk2dst},
			tokenToPool: map[common.Address]common.Address{
				tk1dst: tk1pool,
				tk2dst: tk2pool,
			},
			expRes: map[common.Address]common.Address{
				tk1dst: tk1pool,
				tk2dst: tk2pool,
			},
			expErr: false,
		},
		{
			name: "token pool not found",
			sourceToDestTokens: map[common.Address]common.Address{
				tk1src: tk1dst,
			},
			feeTokens:   []common.Address{tk1dst},
			tokenToPool: map[common.Address]common.Address{},
			expErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockLp := mocks.NewLogPoller(t)
			mockLp.On("LatestBlock", mock.Anything).Return(int64(100), nil)

			offRamp, _ := testhelpers.NewFakeOffRamp(t)
			offRamp.SetSourceToDestTokens(tc.sourceToDestTokens)
			offRamp.SetTokenPools(tc.tokenToPool)

			priceReg, _ := testhelpers.NewFakePriceRegistry(t)
			priceReg.SetFeeTokens(tc.feeTokens)

			c := NewTokenPools(logger.TestLogger(t), mockLp, offRamp, 0, 5)

			res, err := c.Get(ctx)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, len(tc.expRes), len(res))
			for k, v := range tc.expRes {
				assert.Equal(t, v, res[k])
			}
		})
	}
}

func Test_tokenPools_CallOrigin_concurrency(t *testing.T) {
	numDestTokens := rand.Intn(500)
	numWorkers := rand.Intn(500)

	sourceToDestTokens := make(map[common.Address]common.Address, numDestTokens)
	tokenToPool := make(map[common.Address]common.Address)
	for i := 0; i < numDestTokens; i++ {
		sourceToken := utils.RandomAddress()
		destToken := utils.RandomAddress()
		destPool := utils.RandomAddress()
		sourceToDestTokens[sourceToken] = destToken
		tokenToPool[destToken] = destPool
	}

	offRamp, _ := testhelpers.NewFakeOffRamp(t)
	offRamp.SetSourceToDestTokens(sourceToDestTokens)
	offRamp.SetTokenPools(tokenToPool)

	origin := newTokenPoolsOrigin(logger.TestLogger(t), offRamp, numWorkers)
	res, err := origin.CallOrigin(testutils.Context(t))
	assert.NoError(t, err)

	assert.Equal(t, len(tokenToPool), len(res))
	for k, v := range tokenToPool {
		assert.Equal(t, v, res[k])
	}
}
