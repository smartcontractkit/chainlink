package blockhashstore

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type batchBHSConfig interface {
	EvmGasLimitDefault() uint32
}

type BatchBlockhashStore struct {
	config   batchBHSConfig
	txm      txmgr.EvmTxManager
	abi      *abi.ABI
	batchbhs batch_blockhash_store.BatchBlockhashStoreInterface
	lggr     logger.Logger
}

func NewBatchBHS(
	config batchBHSConfig,
	fromAddresses []ethkey.EIP55Address,
	txm txmgr.EvmTxManager,
	batchbhs batch_blockhash_store.BatchBlockhashStoreInterface,
	chainID *big.Int,
	gethks keystore.Eth,
	lggr logger.Logger,
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
		lggr:     lggr,
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

	_, err = b.txm.CreateEthTransaction(txmgr.EvmNewTx{
		FromAddress:    fromAddress,
		ToAddress:      b.batchbhs.Address(),
		EncodedPayload: payload,
		FeeLimit:       b.config.EvmGasLimitDefault(),
		Strategy:       txmgr.NewSendEveryStrategy(),
	}, pg.WithParentCtx(ctx))

	if err != nil {
		return errors.Wrap(err, "creating transaction")
	}

	return nil
}
