package v1_6_1

import (
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_5_1"

	"github.com/ethereum/go-ethereum/common"

	lock_release_token_pool_v1_6_1 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_1/lock_release_token_pool"
)

func GenerateLockReleaseTokenPoolView(pool *lock_release_token_pool_v1_6_1.LockReleaseTokenPool, priceFeed common.Address) (v1_5_1.PoolView, error) {
	basePoolView, err := v1_5_1.GenerateTokenPoolView(pool, priceFeed)
	if err != nil {
		return v1_5_1.PoolView{}, err
	}
	poolView := v1_5_1.PoolView{
		TokenPoolView: basePoolView,
	}
	rebalancer, err := pool.GetRebalancer(nil)
	if err != nil {
		return poolView, err
	}
	poolView.LockReleaseTokenPoolView = v1_5_1.LockReleaseTokenPoolView{
		TokenPoolView:   basePoolView,
		AcceptLiquidity: false,
		Rebalancer:      rebalancer,
	}
	return poolView, nil
}
