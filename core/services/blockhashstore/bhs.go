package blockhashstore

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

var _ BHS = &BulletproofBHS{}

type bpBHSConfig interface {
	EvmGasLimitDefault() uint64
}

// BulletproofBHS is an implementation of BHS that writes "store" transactions to a bulletproof
// transaction manager, and reads BlockhashStore state from the contract.
type BulletproofBHS struct {
	config      bpBHSConfig
	jobID       uuid.UUID
	fromAddress common.Address
	txm         txmgr.TxManager
	abi         *abi.ABI
	bhs         blockhash_store.BlockhashStoreInterface
}

// NewBulletproofBHS creates a new instance with the given transaction manager and blockhash store.
func NewBulletproofBHS(
	config bpBHSConfig,
	fromAddress common.Address,
	txm txmgr.TxManager,
	bhs blockhash_store.BlockhashStoreInterface,
) (*BulletproofBHS, error) {
	bhsABI, err := blockhash_store.BlockhashStoreMetaData.GetAbi()
	if err != nil {
		// blockhash_store.BlockhashStoreABI is generated code, this should never happen
		return nil, errors.Wrap(err, "building ABI")
	}

	return &BulletproofBHS{
		config:      config,
		fromAddress: fromAddress,
		txm:         txm,
		abi:         bhsABI,
		bhs:         bhs,
	}, nil
}

// Store satisfies the BHS interface.
func (c *BulletproofBHS) Store(ctx context.Context, blockNum uint64) error {
	payload, err := c.abi.Pack("store", new(big.Int).SetUint64(blockNum))
	if err != nil {
		return errors.Wrap(err, "packing args")
	}

	_, err = c.txm.CreateEthTransaction(txmgr.NewTx{
		FromAddress:    c.fromAddress,
		ToAddress:      c.bhs.Address(),
		EncodedPayload: payload,
		GasLimit:       c.config.EvmGasLimitDefault(),

		// Set a queue size of 256. At most we store the blockhash of every block, and only the
		// latest 256 can possibly be stored.
		Strategy: txmgr.NewQueueingTxStrategy(c.jobID, 256),
	}, pg.WithParentCtx(ctx))
	if err != nil {
		return errors.Wrap(err, "creating transaction")
	}

	return nil
}

// IsStored satisfies the BHS interface.
func (c *BulletproofBHS) IsStored(ctx context.Context, blockNum uint64) (bool, error) {
	_, err := c.bhs.GetBlockhash(&bind.CallOpts{Context: ctx}, big.NewInt(int64(blockNum)))
	if err != nil && strings.Contains(err.Error(), "reverted") {
		// Transaction reverted because the blockhash is not stored
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, "getting blockhash")
	}
	return true, nil
}
