package blockhashstore

import (
	"context"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/core/logger"
)

var _ BackwardsBHS = &backwardsBHS{}

type backwardsBHSConfig interface {
	EvmGasLimitDefault() uint64
}

type backwardsBHS struct {
	txm         txmgr.TxManager
	ethClient   evmclient.Client
	fromAddress common.Address
	lg          logger.Logger
}

// BackwardsBHS provides methods to interact with the BatchBlockhashStore contract.
type BackwardsBHS interface {
	// Backwards fills the blockhash store backwards starting from the given startBlock
	// and stopping at the provided endBlock.
	// Preconditions:
	// 1. The startBlock parameter provided must have an associated blockhash in the blockhash store.
	// 2. endBlock < startBlock
	// 3. The batchBHSAddress must point to a valid BatchBlockhashStore contract.
	Backwards(
		startBlock, endBlock int64,
		batchBHSAddress common.Address,
		batchSize int64,
	) error
}

// NewBackwardsBHS creates a BackwardsBHS with the provided parameters.
// The fromAddress parameter is the address to use to send transactions.
// This address must be available in the Chainlink node's Ethereum keystore.
func NewBackwardsBHS(
	txm txmgr.TxManager,
	evmClient evmclient.Client,
	fromAddress common.Address,
	lg logger.Logger,
) BackwardsBHS {
	return &backwardsBHS{
		txm:         txm,
		ethClient:   evmClient,
		fromAddress: fromAddress,
		lg:          lg,
	}
}

// Backwards fills the blockhash store backwards starting from the given startBlock
// and stopping at the provided endBlock.
// Preconditions:
// 1. The startBlock parameter provided must have an associated blockhash in the blockhash store.
// 2. endBlock < startBlock
// 3. The batchBHSAddress must point to a valid BatchBlockhashStore contract.
func (b *backwardsBHS) Backwards(
	startBlock, endBlock int64,
	batchBHSAddress common.Address,
	batchSize int64,
) error {
	batchBHS, err := batch_blockhash_store.NewBatchBlockhashStore(batchBHSAddress, b.ethClient)
	if err != nil {
		return errors.Wrap(err, "creating batch bhs client")
	}

	batchBHSABI, err := batch_blockhash_store.BatchBlockhashStoreMetaData.GetAbi()
	if err != nil {
		// should never happen
		return errors.Wrap(err, "get batch bhs abi")
	}

	// check if startBlock is in the BHS
	bh, err := batchBHS.GetBlockhashes(
		&bind.CallOpts{Context: context.Background()},
		[]*big.Int{big.NewInt(startBlock)})
	if err != nil {
		return errors.Wrapf(err, "calling getBlockhashes([%d])", startBlock)
	}

	if reflect.DeepEqual(bh[0], [32]byte{}) {
		return fmt.Errorf("blockhash of block %d not in the blockhash store, please provide one in the store", startBlock)
	}

	// startBlock has a blockhash in the BHS, we're good to start feeding
	fromBlock := startBlock - 1
	blockRange, err := decreasingBlockRange(big.NewInt(fromBlock), big.NewInt(endBlock))
	if err != nil {
		return errors.Wrap(err, "create block range")
	}

	// iterate through the blockrange in batchSize batches, enqueue storeVerifyHeader tx into bptxm
	for i := 0; i < len(blockRange); i += int(batchSize) {
		j := i + int(batchSize)
		if j > len(blockRange) {
			j = len(blockRange)
		}

		blockNumbers := blockRange[i:j]
		blockHeaders, err := b.storeVerifyHeaders(blockNumbers)
		if err != nil {
			return errors.Wrapf(err, "getting block headers for blocks %+v", blockNumbers)
		}

		payload, err := batchBHSABI.Pack("storeVerifyHeader", blockNumbers, blockHeaders)
		if err != nil {
			return errors.Wrapf(err, "packing storeVerifyHeader(%+v, ...)", blockNumbers)
		}

		gasLimit, err := b.ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
			From: b.fromAddress,
			To:   &batchBHSAddress,
			Gas:  10e6, // overestimate due to potentially large batch sizes, real value will likely be smaller
			Data: payload,
		})
		if err != nil {
			return errors.Wrap(err, "estimating gas")
		}

		b.lg.Debugw("Sending storeVerifyHeader tx", "blockNumbers", blockNumbers)
		_, err = b.txm.CreateEthTransaction(txmgr.NewTx{
			FromAddress:    b.fromAddress,
			ToAddress:      batchBHSAddress,
			EncodedPayload: payload,
			GasLimit:       gasLimit,
		})
		if err != nil {
			return errors.Wrap(err, "enqueue bptxm")
		}
	}

	b.lg.Debug("Done feeding blocks", startBlock-1, "to", endBlock)

	return nil
}

func (b *backwardsBHS) storeVerifyHeaders(blockRange []*big.Int) (headers [][]byte, err error) {
	headers = [][]byte{}
	for _, blockNum := range blockRange {
		// Get child block since it's the one that has the parent hash in it's header.
		h, err := b.ethClient.HeaderByNumber(
			context.Background(),
			new(big.Int).Add(blockNum, big.NewInt(1)),
		)
		if err != nil {
			return nil, errors.Wrap(err, "get header")
		}
		rlpHeader, err := rlp.EncodeToBytes(h)
		if err != nil {
			return nil, errors.Wrap(err, "encode rlp")
		}
		headers = append(headers, rlpHeader)
	}
	return
}

// decreasingBlockRange creates a continugous block range starting with
// block `start` and ending at block `end`.
func decreasingBlockRange(start, end *big.Int) (ret []*big.Int, err error) {
	if start.Cmp(end) == -1 {
		return nil, fmt.Errorf("start (%s) must be greater than end (%s)", start.String(), end.String())
	}
	ret = []*big.Int{}
	for i := new(big.Int).Set(start); i.Cmp(end) >= 0; i.Sub(i, big.NewInt(1)) {
		ret = append(ret, new(big.Int).Set(i))
	}
	return
}
