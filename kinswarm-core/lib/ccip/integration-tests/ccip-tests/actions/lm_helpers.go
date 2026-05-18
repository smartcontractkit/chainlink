package actions

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/contracts"
)

type LMCommon struct {
	ChainClient       blockchain.EVMClient
	MockArm           *common.Address
	ArmProxy          *contracts.ArmProxy
	CcipRouter        *contracts.Router
	LM                *contracts.LiquidityManager
	TokenPool         *contracts.TokenPool
	BridgeAdapterAddr *common.Address
	WrapperNative     *common.Address
	MinimumLiquidity  *big.Int
	ChainSelectror    uint64
}

func DefaultLMModule(chainClient blockchain.EVMClient, minimumLiquidity *big.Int, chainSelector uint64) (*LMCommon, error) {
	return &LMCommon{
		ChainClient:      chainClient,
		MinimumLiquidity: minimumLiquidity,
		ChainSelectror:   chainSelector,
	}, nil
}
