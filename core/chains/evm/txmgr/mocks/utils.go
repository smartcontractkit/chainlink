package mocks

import evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"

type MockEvmTxManager = TxManager[evmtypes.Address, *evmtypes.TxHash, evmtypes.BlockHash]
