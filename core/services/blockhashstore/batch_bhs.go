package blockhashstore

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
)

type batchBHSConfig interface {
	LimitDefault() uint64
}

type BatchBlockhashStore struct {
	config   batchBHSConfig
	txm      txmgr.TxManager
	abi      *abi.ABI
	batchbhs batch_blockhash_store.BatchBlockhashStoreInterface
}

func NewBatchBHS(
	config batchBHSConfig,
	txm txmgr.TxManager,
	batchbhs batch_blockhash_store.BatchBlockhashStoreInterface,
) (*BatchBlockhashStore, error) {
	abi, err := batch_blockhash_store.BatchBlockhashStoreMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "building ABI")
	}
	return &BatchBlockhashStore{
		config:   config,
		txm:      txm,
		abi:      abi,
		batchbhs: batchbhs,
	}, nil
}

func (b *BatchBlockhashStore) GetBlockhashes(ctx context.Context, blockNumbers []*big.Int) ([][32]byte, error) {
	blockhashes, err := b.batchbhs.GetBlockhashes(&bind.CallOpts{Context: ctx}, blockNumbers)
	if err != nil {
		return nil, errors.Wrap(err, "getting blockhashes")
	}
	return blockhashes, nil
}

func (b *BatchBlockhashStore) StoreVerifyHeader(ctx context.Context, blockNumbers []*big.Int, blockHeaders [][]byte, fromAddress common.Address) error {
	payload, err := b.abi.Pack("storeVerifyHeader", blockNumbers, blockHeaders)
	if err != nil {
		return errors.Wrap(err, "packing args")
	}

	_, err = b.txm.CreateTransaction(ctx, txmgr.TxRequest{
		FromAddress:    fromAddress,
		ToAddress:      b.batchbhs.Address(),
		EncodedPayload: payload,
		FeeLimit:       b.config.LimitDefault(),
		Strategy:       txmgrcommon.NewSendEveryStrategy(),
	})

	if err != nil {
		return errors.Wrap(err, "creating transaction")
	}

	return nil
}
