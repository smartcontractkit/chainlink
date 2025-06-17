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
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/trusted_blockhash_store"
	evmkeystore "github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
)

var _ BHS = &BulletproofBHS{}

type bpBHSConfig interface {
	LimitDefault() uint64
}

type bpBHSDatabaseConfig interface {
	DefaultQueryTimeout() time.Duration
}

// BulletproofBHS is an implementation of BHS that writes "store" transactions to a bulletproof
// transaction manager, and reads BlockhashStore state from the contract.
type BulletproofBHS struct {
	config        bpBHSConfig
	dbConfig      bpBHSDatabaseConfig
	jobID         uuid.UUID
	fromAddresses []types.EIP55Address
	txm           txmgr.TxManager
	abi           *abi.ABI
	trustedAbi    *abi.ABI
	bhs           blockhash_store.BlockhashStoreInterface
	trustedBHS    *trusted_blockhash_store.TrustedBlockhashStore
	gethks        evmkeystore.RoundRobin
}

// NewBulletproofBHS creates a new instance with the given transaction manager and blockhash store.
func NewBulletproofBHS(
	config bpBHSConfig,
	dbConfig bpBHSDatabaseConfig,
	fromAddresses []types.EIP55Address,
	txm txmgr.TxManager,
	bhs blockhash_store.BlockhashStoreInterface,
	trustedBHS *trusted_blockhash_store.TrustedBlockhashStore,
	gethks evmkeystore.RoundRobin,
) (*BulletproofBHS, error) {
	bhsABI, err := blockhash_store.BlockhashStoreMetaData.GetAbi()
	if err != nil {
		// blockhash_store.BlockhashStoreABI is generated code, this should never happen
		return nil, errors.Wrap(err, "building ABI")
	}

	trustedBHSAbi, err := trusted_blockhash_store.TrustedBlockhashStoreMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "building trusted BHS ABI")
	}

	return &BulletproofBHS{
		config:        config,
		dbConfig:      dbConfig,
		fromAddresses: fromAddresses,
		txm:           txm,
		abi:           bhsABI,
		trustedAbi:    trustedBHSAbi,
		bhs:           bhs,
		trustedBHS:    trustedBHS,
		gethks:        gethks,
	}, nil
}

// Store satisfies the BHS interface.
func (c *BulletproofBHS) Store(ctx context.Context, blockNum uint64) error {
	payload, err := c.abi.Pack("store", new(big.Int).SetUint64(blockNum))
	if err != nil {
		return errors.Wrap(err, "packing args")
	}

	fromAddress, err := c.gethks.GetNextAddress(ctx, SendingKeys(c.fromAddresses)...)
	if err != nil {
		return errors.Wrap(err, "getting next from address")
	}

	_, err = c.txm.CreateTransaction(ctx, txmgr.TxRequest{
		FromAddress:    fromAddress,
		ToAddress:      c.bhs.Address(),
		EncodedPayload: payload,
		FeeLimit:       c.config.LimitDefault(),

		// Set a queue size of 256. At most we store the blockhash of every block, and only the
		// latest 256 can possibly be stored.
		Strategy: txmgrcommon.NewQueueingTxStrategy(c.jobID, 256),
	})
	if err != nil {
		return errors.Wrap(err, "creating transaction")
	}

	return nil
}

func (c *BulletproofBHS) StoreTrusted(
	ctx context.Context,
	blockNums []uint64,
	blockhashes []common.Hash,
	recentBlock uint64,
	recentBlockhash common.Hash,
) error {
	// Convert and pack arguments for a "storeTrusted" function call to the trusted BHS.
	var blockNumsBig []*big.Int
	for _, b := range blockNums {
		blockNumsBig = append(blockNumsBig, new(big.Int).SetUint64(b))
	}
	recentBlockBig := new(big.Int).SetUint64(recentBlock)
	payload, err := c.trustedAbi.Pack("storeTrusted", blockNumsBig, blockhashes, recentBlockBig, recentBlockhash)
	if err != nil {
		return errors.Wrap(err, "packing args")
	}

	// Create a transaction from the given batch and send it to the TXM.
	fromAddress, err := c.gethks.GetNextAddress(ctx, SendingKeys(c.fromAddresses)...)
	if err != nil {
		return errors.Wrap(err, "getting next from address")
	}
	_, err = c.txm.CreateTransaction(ctx, txmgr.TxRequest{
		FromAddress:    fromAddress,
		ToAddress:      c.trustedBHS.Address(),
		EncodedPayload: payload,
		FeeLimit:       c.config.LimitDefault(),

		Strategy: txmgrcommon.NewSendEveryStrategy(),
	})
	if err != nil {
		return errors.Wrap(err, "creating transaction")
	}

	return nil
}

func (c *BulletproofBHS) IsTrusted() bool {
	return c.trustedBHS != nil
}

// IsStored satisfies the BHS interface.
func (c *BulletproofBHS) IsStored(ctx context.Context, blockNum uint64) (bool, error) {
	var err error
	if c.IsTrusted() {
		_, err = c.trustedBHS.GetBlockhash(&bind.CallOpts{Context: ctx}, big.NewInt(int64(blockNum)))
	} else {
		_, err = c.bhs.GetBlockhash(&bind.CallOpts{Context: ctx}, big.NewInt(int64(blockNum)))
	}
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

	fromAddress, err := c.gethks.GetNextAddress(ctx, c.sendingKeys()...)
	if err != nil {
		return errors.Wrap(err, "getting next from address")
	}

	_, err = c.txm.CreateTransaction(ctx, txmgr.TxRequest{
		FromAddress:    fromAddress,
		ToAddress:      c.bhs.Address(),
		EncodedPayload: payload,
		FeeLimit:       c.config.LimitDefault(),
		Strategy:       txmgrcommon.NewSendEveryStrategy(),
	})
	if err != nil {
		return errors.Wrap(err, "creating transaction")
	}

	return nil
}
