package blockhashstore

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type batchBHSConfig interface {
	EvmGasLimitDefault() uint32
}

type BatchBlockhashStore struct {
	config        batchBHSConfig
	fromAddresses []ethkey.EIP55Address
	txm           txmgr.TxManager
	abi           *abi.ABI
	chainID       *big.Int
	gethks        keystore.Eth
	batchbhs      batch_blockhash_store.BatchBlockhashStoreInterface
}

func NewBatchBHS(
	config batchBHSConfig,
	fromAddresses []ethkey.EIP55Address,
	txm txmgr.TxManager,
	batchbhs batch_blockhash_store.BatchBlockhashStoreInterface,
	chainID *big.Int,
	gethks keystore.Eth,
) (*BatchBlockhashStore, error) {
	abi, err := batch_blockhash_store.BatchBlockhashStoreMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "building ABI")
	}
	return &BatchBlockhashStore{
		config:        config,
		fromAddresses: fromAddresses,
		txm:           txm,
		abi:           abi,
		batchbhs:      batchbhs,
		chainID:       chainID,
		gethks:        gethks,
	}, nil
}

func (b *BatchBlockhashStore) GetBlockhashes(ctx context.Context, blockNumbers []*big.Int) ([][32]byte, error) {
	blockhashes, err := b.batchbhs.GetBlockhashes(&bind.CallOpts{Context: ctx}, blockNumbers)
	if err != nil {
		return nil, errors.Wrap(err, "getting blockhashes")
	}
	return blockhashes, nil
}

func (b *BatchBlockhashStore) StoreVerifyHeader(ctx context.Context, blockNumbers []*big.Int, blockHeaders [][]byte) error {
	payload, err := b.abi.Pack("storeVerifyHeader", blockNumbers, blockHeaders)
	if err != nil {
		return errors.Wrap(err, "packing args")
	}

	// TODO: re-evaluate if it's better to have 1 fromAddress
	// because ordering matters for StoreVerifyHeader
	fromAddress, err := b.gethks.GetRoundRobinAddress(b.chainID, b.sendingKeys()...)
	if err != nil {
		return errors.Wrap(err, "getting next from address")
	}

	_, err = b.txm.CreateEthTransaction(txmgr.NewTx{
		FromAddress:    fromAddress,
		ToAddress:      b.batchbhs.Address(),
		EncodedPayload: payload,
		GasLimit:       b.config.EvmGasLimitDefault(),
		Strategy:       txmgr.NewSendEveryStrategy(),
	}, pg.WithParentCtx(ctx))

	if err != nil {
		return errors.Wrap(err, "creating transaction")
	}

	return nil
}

func (b *BatchBlockhashStore) sendingKeys() []common.Address {
	var keys []common.Address
	for _, a := range b.fromAddresses {
		keys = append(keys, a.Address())
	}
	return keys
}
