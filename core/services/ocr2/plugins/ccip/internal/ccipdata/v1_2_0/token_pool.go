package v1_2_0

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"

	burn_mint_token_pool_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

var (
	poolABI = abihelpers.MustParseABI(burn_mint_token_pool_1_2_0.BurnMintTokenPoolABI)
)

var _ ccipdata.TokenPoolReader = &TokenPool{}

type TokenPool struct {
	addr           common.Address
	OffRampAddress common.Address
	poolType       string
}

func NewTokenPool(poolType string, addr common.Address, offRampAddress common.Address) *TokenPool {
	return &TokenPool{
		addr:           addr,
		OffRampAddress: offRampAddress,
		poolType:       poolType,
	}
}

func (p *TokenPool) Address() common.Address {
	return p.addr
}

func (p *TokenPool) Type() string {
	return p.poolType
}

func GetInboundTokenPoolRateLimitCall(tokenPoolAddress common.Address, offRampAddress common.Address) rpclib.EvmCall {
	return rpclib.NewEvmCall(
		poolABI,
		"currentOffRampRateLimiterState",
		tokenPoolAddress,
		offRampAddress,
	)
}
