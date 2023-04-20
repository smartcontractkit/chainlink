package mocks

import common "github.com/ethereum/go-ethereum/common"

type MockEvmTxManager = TxManager[common.Address, common.Hash, common.Hash]
