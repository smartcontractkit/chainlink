package ocrcommon

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type txManager interface {
	CreateEthTransaction(newTx txmgr.NewTx, qopts ...pg.QOpt) (etx txmgr.EthTx, err error)
}

type Transmitter interface {
	CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte) error
	FromAddress() common.Address
}
