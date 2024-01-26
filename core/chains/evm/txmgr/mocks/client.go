package mocks

import (
	"math/big"
	"testing"

	gethCommon "github.com/ethereum/go-ethereum/common"

	txmgrtypesmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

type EvmTxmClient = txmgrtypesmocks.TxmClient[*big.Int, gethCommon.Address, gethCommon.Hash, gethCommon.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee]

func NewEvmTxmClient(t testing.TB) *EvmTxmClient {
	return txmgrtypesmocks.NewTxmClient[*big.Int, gethCommon.Address, gethCommon.Hash, gethCommon.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee](t)
}

func NewEvmTxmClientWithDefaultChain(t testing.TB) *EvmTxmClient {
	c := NewEvmTxmClient(t)
	c.On("ConfiguredChainID").Return(testutils.FixtureChainID).Maybe()
	c.On("IsL2").Return(false).Maybe()
	return c
}
