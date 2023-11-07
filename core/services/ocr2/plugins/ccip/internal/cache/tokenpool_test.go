package cache

import (
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
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

			offRamp := ccipdatamocks.NewOffRampReader(t)
			offRamp.On("TokenEvents").Return([]common.Hash{})
			offRamp.On("Address").Return(utils.RandomAddress())
			destTokens := make([]common.Address, 0, len(tc.sourceToDestTokens))
			for _, tk := range tc.sourceToDestTokens {
				destTokens = append(destTokens, tk)
			}
			for destToken, pool := range tc.tokenToPool {
				offRamp.On("GetPoolByDestToken", mock.Anything, destToken).Return(pool, nil)
			}
			for _, destTk := range tc.sourceToDestTokens {
				if _, exists := tc.tokenToPool[destTk]; !exists {
					offRamp.On("GetPoolByDestToken", mock.Anything, destTk).Return(nil, errors.New("not found"))
				}
			}
			offRamp.On("GetDestinationTokens", mock.Anything).Return(destTokens, nil)

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
	destTokens := make([]common.Address, 0, numDestTokens)
	tokenToPool := make(map[common.Address]common.Address)
	for i := 0; i < numDestTokens; i++ {
		sourceToken := utils.RandomAddress()
		destToken := utils.RandomAddress()
		destPool := utils.RandomAddress()
		sourceToDestTokens[sourceToken] = destToken
		tokenToPool[destToken] = destPool
		destTokens = append(destTokens, destToken)
	}

	offRamp := ccipdatamocks.NewOffRampReader(t)
	offRamp.On("GetDestinationTokens", mock.Anything).Return(destTokens, nil)
	for destToken, pool := range tokenToPool {
		offRamp.On("GetPoolByDestToken", mock.Anything, destToken).Return(pool, nil)
	}

	origin := newTokenPoolsOrigin(logger.TestLogger(t), offRamp, numWorkers)
	res, err := origin.CallOrigin(testutils.Context(t))
	assert.NoError(t, err)

	assert.Equal(t, len(tokenToPool), len(res))
	for k, v := range tokenToPool {
		assert.Equal(t, v, res[k])
	}
}
