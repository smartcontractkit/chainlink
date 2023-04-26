// The blockhash store package provides a service that stores blockhashes such that they are available
// for on-chain proofs beyond the EVM 256 block limit.
package blockhashstore

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var _ BHS = &BulletproofBHS{}

type bpBHSConfig interface {
	EvmGasLimitDefault() uint32
	DatabaseDefaultQueryTimeout() time.Duration
}

// BulletproofBHS is an implementation of BHS that writes "store" transactions to a bulletproof
// transaction manager, and reads BlockhashStore state from the contract.
type BulletproofBHS struct {
	config        bpBHSConfig
	jobID         uuid.UUID
	fromAddresses []ethkey.EIP55Address
	txm           txmgr.EvmTxManager
	abi           *abi.ABI
	bhs           blockhash_store.BlockhashStoreInterface
	chainID       *big.Int
	gethks        keystore.Eth
}

// NewBulletproofBHS creates a new instance with the given transaction manager and blockhash store.
func NewBulletproofBHS(
	config bpBHSConfig,
	fromAddresses []ethkey.EIP55Address,
	txm txmgr.EvmTxManager,
	bhs blockhash_store.BlockhashStoreInterface,
	chainID *big.Int,
	gethks keystore.Eth,
) (*BulletproofBHS, error) {
	bhsABI, err := blockhash_store.BlockhashStoreMetaData.GetAbi()
	if err != nil {
		// blockhash_store.BlockhashStoreABI is generated code, this should never happen
		return nil, errors.Wrap(err, "building ABI")
	}

	return &BulletproofBHS{
		config:        config,
		fromAddresses: fromAddresses,
		txm:           txm,
		abi:           bhsABI,
		bhs:           bhs,
		chainID:       chainID,
		gethks:        gethks,
	}, nil
}

// Store satisfies the BHS interface.
func (c *BulletproofBHS) Store(ctx context.Context, blockNum uint64) error {
	payload, err := c.abi.Pack("store", new(big.Int).SetUint64(blockNum))
	if err != nil {
		return errors.Wrap(err, "packing args")
	}

	fromAddress, err := c.gethks.GetRoundRobinAddress(c.chainID, SendingKeys(c.fromAddresses)...)
	if err != nil {
		return errors.Wrap(err, "getting next from address")
	}

	_, err = c.txm.CreateEthTransaction(txmgr.EvmNewTx{
		FromAddress:    fromAddress,
		ToAddress:      c.bhs.Address(),
		EncodedPayload: payload,
		FeeLimit:       c.config.EvmGasLimitDefault(),

		// Set a queue size of 256. At most we store the blockhash of every block, and only the
		// latest 256 can possibly be stored.
		Strategy: txmgr.NewQueueingTxStrategy(c.jobID, 256, c.config.DatabaseDefaultQueryTimeout()),
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

func (c *BulletproofBHS) sendingKeys() []common.Address {
	var keys []common.Address
	for _, a := range c.fromAddresses {
		keys = append(keys, a.Address())
	}
	return keys
}

func (c *BulletproofBHS) StoreEarliest(ctx context.Context) error {
	payload, err := c.abi.Pack("storeEarliest")
	if err != nil {
		return errors.Wrap(err, "packing args")
	}

	fromAddress, err := c.gethks.GetRoundRobinAddress(c.chainID, c.sendingKeys()...)
	if err != nil {
		return errors.Wrap(err, "getting next from address")
	}

	_, err = c.txm.CreateEthTransaction(txmgr.EvmNewTx{
		FromAddress:    fromAddress,
		ToAddress:      c.bhs.Address(),
		EncodedPayload: payload,
		FeeLimit:       c.config.EvmGasLimitDefault(),
		Strategy:       txmgr.NewSendEveryStrategy(),
	}, pg.WithParentCtx(ctx))
	if err != nil {
		return errors.Wrap(err, "creating transaction")
	}

	return nil
}
