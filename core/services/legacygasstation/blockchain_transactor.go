package legacygasstation

import (
	"context"
	"encoding/hex"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation"
	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation/types"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/sqlx"
)

var _ legacygasstation.BlockchainTransactor = &blockchainTransactor{}

type blockchainTransactor struct {
	lggr          logger.Logger
	db            *sqlx.DB
	txm           txmgr.TxManager
	gethks        keystore.Eth
	fromAddresses []ethkey.EIP55Address
	chainID       uint64
	orm           legacygasstation.ORM
}

func NewBlockchainTransactor(
	lggr logger.Logger,
	db *sqlx.DB,
	txm txmgr.TxManager,
	gethks keystore.Eth,
	fromAddresses []ethkey.EIP55Address,
	chainID uint64,
	orm legacygasstation.ORM,
) (*blockchainTransactor, error) {
	return &blockchainTransactor{
		lggr:          lggr,
		db:            db,
		txm:           txm,
		gethks:        gethks,
		fromAddresses: fromAddresses,
		chainID:       chainID,
		orm:           orm,
	}, nil
}

// CreateAndStoreTransaction creates eth transaction and persists data in a transaction
// to avoid partial failures, which would leave the persistence layer in inconsistent state
func (t *blockchainTransactor) CreateAndStoreTransaction(
	ctx context.Context,
	address gethcommon.Address,
	payload []byte,
	gasLimit uint32,
	req types.SendTransactionRequest,
	requestID string,
) error {
	fromAddresses := t.sendingKeys()
	fromAddress, err := t.gethks.GetRoundRobinAddress(big.NewInt(0).SetUint64(t.chainID), fromAddresses...)
	if err != nil {
		return err
	}

	// idempotencyKey is unique because payload contains nonce that can only be used once
	idempotencyKey := hex.EncodeToString(crypto.Keccak256(payload)[:])

	txmTx, err := t.txm.CreateTransaction(txmgr.TxRequest{
		IdempotencyKey: &idempotencyKey,
		FromAddress:    fromAddress,
		ToAddress:      address,
		EncodedPayload: payload,
		FeeLimit:       gasLimit,
		Strategy:       txmgrcommon.NewSendEveryStrategy(),
	})
	if err != nil {
		return err
	}

	t.lggr.Debugw("created Eth tx", "ethTxID", txmTx.GetID())

	gaslessTx := types.LegacyGaslessTx{
		ID:                 requestID,
		From:               req.From,
		Target:             req.Target,
		Forwarder:          address,
		Nonce:              req.Nonce,
		Receiver:           req.Receiver,
		Amount:             req.Amount,
		SourceChainID:      req.SourceChainID,
		DestinationChainID: req.DestinationChainID,
		ValidUntilTime:     req.ValidUntilTime,
		Signature:          req.Signature,
		Status:             types.Submitted,
		TokenName:          req.TargetName,
		TokenVersion:       req.Version,
		EthTxID:            txmTx.GetID(),
	}

	return t.orm.InsertLegacyGaslessTx(ctx, gaslessTx)
}

func (t *blockchainTransactor) sendingKeys() []gethcommon.Address {
	var addresses []gethcommon.Address
	for _, a := range t.fromAddresses {
		addresses = append(addresses, a.Address())
	}
	return addresses
}
