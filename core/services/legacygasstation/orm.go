package legacygasstation

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation"
	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var _ legacygasstation.ORM = &orm{}

type orm struct {
	q    pg.Q
	lggr logger.Logger
}

const InsertLegacyGaslessTxQuery = `INSERT INTO legacy_gasless_txs (legacy_gasless_tx_id, forwarder_address, from_address, target_address, receiver_address, nonce, amount, source_chain_id, destination_chain_id, valid_until_time, tx_signature, tx_status, token_name, token_version, eth_tx_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW(), NOW())`

const UpdateLegacyGaslessTxQuery = `UPDATE legacy_gasless_txs SET 
tx_status = $2,
ccip_message_id = $3,
failure_reason = $4,
tx_hash = $5,
updated_at = NOW()
WHERE legacy_gasless_tx_id = $1`

// NewORM creates an ORM scoped to chainID.
// TODO: implement pruning logic if needed
func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) legacygasstation.ORM {
	namedLogger := lggr.Named("LegacyGasStation")
	q := pg.NewQ(db, namedLogger, cfg)
	return &orm{
		q:    q,
		lggr: namedLogger,
	}
}

type DbLegacyGaslessTx struct {
	ID                 string         `db:"legacy_gasless_tx_id"` // UUID
	Forwarder          common.Address `db:"forwarder_address"`    // forwarder contract
	From               common.Address `db:"from_address"`         // token sender
	Target             common.Address `db:"target_address"`       // token contract
	Receiver           common.Address `db:"receiver_address"`     // token receiver
	Nonce              *utils.Big     // forwarder nonce
	Amount             *utils.Big     // token amount to be transferred
	SourceChainID      uint64         // meta-transaction source chain ID. This is CCIP chain selector instead of EVM chain ID.
	DestinationChainID uint64         // meta-transaction destination chain ID. This is CCIP chain selector instead of EVM chain ID.
	ValidUntilTime     *utils.Big     // unix timestamp of meta-transaction expiry in seconds
	Signature          []byte         `db:"tx_signature"` // EIP712 signature
	Status             string         `db:"tx_status"`    // status of meta-transaction
	FailureReason      *string        // failure reason of meta-transaction TODO: change this to sql.NullString
	TokenName          string         // name of token used to generate EIP712 domain separator hash
	TokenVersion       string         // version of token used to generate EIP712 domain separator hash
	CCIPMessageID      *common.Hash   `db:"ccip_message_id"` // CCIP message ID
	EthTxID            string         `db:"eth_tx_id"`       // tx ID in transaction manager
	TxHash             *common.Hash   `db:"tx_hash"`         // transaction hash on source chain
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func toDbLegacyGaslessTx(tx types.LegacyGaslessTx) DbLegacyGaslessTx {
	return DbLegacyGaslessTx{
		ID:                 tx.ID,
		Forwarder:          tx.Forwarder,
		From:               tx.From,
		Target:             tx.Target,
		Receiver:           tx.Receiver,
		Nonce:              utils.NewBig(tx.Nonce),
		Amount:             utils.NewBig(tx.Amount),
		SourceChainID:      tx.SourceChainID,
		DestinationChainID: tx.DestinationChainID,
		ValidUntilTime:     utils.NewBig(tx.ValidUntilTime),
		Signature:          tx.Signature[:],
		Status:             tx.Status.String(),
		FailureReason:      tx.FailureReason,
		TokenName:          tx.TokenName,
		TokenVersion:       tx.TokenVersion,
		CCIPMessageID:      tx.CCIPMessageID,
		EthTxID:            tx.EthTxID,
		TxHash:             tx.TxHash,
		CreatedAt:          tx.CreatedAt,
		UpdatedAt:          tx.UpdatedAt,
	}
}

func toLegacyGaslessTx(dbTx DbLegacyGaslessTx) (*types.LegacyGaslessTx, error) {
	var status types.Status
	err := status.Scan(dbTx.Status)
	if err != nil {
		return nil, err
	}
	return &types.LegacyGaslessTx{
		ID:                 dbTx.ID,
		Forwarder:          dbTx.Forwarder,
		From:               dbTx.From,
		Target:             dbTx.Target,
		Receiver:           dbTx.Receiver,
		Nonce:              dbTx.Nonce.ToInt(),
		Amount:             dbTx.Amount.ToInt(),
		SourceChainID:      dbTx.SourceChainID,
		DestinationChainID: dbTx.DestinationChainID,
		ValidUntilTime:     dbTx.ValidUntilTime.ToInt(),
		Signature:          dbTx.Signature,
		Status:             status,
		FailureReason:      dbTx.FailureReason,
		TokenName:          dbTx.TokenName,
		TokenVersion:       dbTx.TokenVersion,
		CCIPMessageID:      dbTx.CCIPMessageID,
		EthTxID:            dbTx.EthTxID,
		TxHash:             dbTx.TxHash,
		CreatedAt:          dbTx.CreatedAt,
		UpdatedAt:          dbTx.UpdatedAt,
	}, nil
}

// DbLegacyGaslessTxPlus has additional fieds from evm.txes and evm.tx_attempts table
type DbLegacyGaslessTxPlus struct {
	DbLegacyGaslessTx
	EthTxStatus txmgrtypes.TxState `db:"etx_state"`
	EthTxHash   *common.Hash       `db:"etx_hash"`
	EthTxError  *string            `db:"etx_error"`
	Receipt     *evmtypes.Receipt  `db:"etx_receipt"`
}

func (o *orm) SelectBySourceChainIDAndStatus(ctx context.Context, sourceChainID uint64, status types.Status) (txs []types.LegacyGaslessTx, err error) {
	q := o.q.WithOpts(pg.WithParentCtx(ctx))
	var dbTxs []DbLegacyGaslessTx
	err = q.Select(&dbTxs, `
		SELECT * FROM legacy_gasless_txs 
			WHERE legacy_gasless_txs.source_chain_id = $1 
			AND legacy_gasless_txs.tx_status = $2
		`, sourceChainID, status.String())
	for _, dbTx := range dbTxs {
		tx, err := toLegacyGaslessTx(dbTx)
		if err != nil {
			return nil, err
		}
		txs = append(txs, *tx)
	}
	return
}

func (o *orm) SelectByDestChainIDAndStatus(ctx context.Context, destChainID uint64, status types.Status) (txs []types.LegacyGaslessTx, err error) {
	q := o.q.WithOpts(pg.WithParentCtx(ctx))
	var dbTxs []DbLegacyGaslessTx
	err = q.Select(&dbTxs, `
		SELECT * FROM legacy_gasless_txs
			WHERE legacy_gasless_txs.destination_chain_id = $1 
			AND legacy_gasless_txs.tx_status = $2
		`, destChainID, status.String())
	for _, dbTx := range dbTxs {
		tx, err := toLegacyGaslessTx(dbTx)
		if err != nil {
			return nil, err
		}
		txs = append(txs, *tx)
	}
	return
}

// InsertLegacyGaslessTx is idempotent
func (o *orm) InsertLegacyGaslessTx(ctx context.Context, lgsTx types.LegacyGaslessTx) error {
	q := o.q.WithOpts(pg.WithParentCtx(ctx))
	dbTx := toDbLegacyGaslessTx(lgsTx)
	return q.ExecQ(InsertLegacyGaslessTxQuery,
		dbTx.ID,
		dbTx.Forwarder,
		dbTx.From,
		dbTx.Target,
		dbTx.Receiver,
		dbTx.Nonce,
		dbTx.Amount,
		dbTx.SourceChainID,
		dbTx.DestinationChainID,
		dbTx.ValidUntilTime,
		dbTx.Signature[:],
		dbTx.Status,
		dbTx.TokenName,
		dbTx.TokenVersion,
		dbTx.EthTxID,
	)
}

// UpdateLegacyGaslessTx updates legacy gasless transaction with status, ccip message ID (optional), failure reason (optional)
func (o *orm) UpdateLegacyGaslessTx(ctx context.Context, lgsTx types.LegacyGaslessTx) error {
	q := o.q.WithOpts(pg.WithParentCtx(ctx))
	dbTx := toDbLegacyGaslessTx(lgsTx)
	return q.ExecQ(UpdateLegacyGaslessTxQuery,
		dbTx.ID,
		dbTx.Status,
		dbTx.CCIPMessageID,
		dbTx.FailureReason,
		dbTx.TxHash,
	)
}

func (o *orm) SelectBySourceChainIDAndEthTxStates(ctx context.Context, sourceChainID uint64, states []legacygasstation.EtxStatus) ([]types.LegacyGaslessTxPlus, error) {
	var dbLgps []DbLegacyGaslessTxPlus
	q := o.q.WithOpts(pg.WithParentCtx(ctx))
	err := q.Select(&dbLgps, `SELECT 
		lgt.*,
		etx.state as etx_state,
		eta.hash as etx_hash,
		etx.error as etx_error,
		etr.receipt as etx_receipt
	FROM legacy_gasless_txs lgt
	LEFT JOIN evm.txes etx ON etx.id = lgt.eth_tx_id
	LEFT JOIN evm.tx_attempts eta ON etx.id = eta.eth_tx_id
	LEFT JOIN evm.receipts etr ON eta.hash = etr.tx_hash
	WHERE lgt.source_chain_id = $1
	AND etx.state = any($2)
	ORDER BY eta.broadcast_before_block_num ASC, etr.block_number ASC
	`, sourceChainID, pq.Array(states))
	if err != nil {
		return nil, errors.Wrap(err, "select eth txs by source chain id and states")
	}

	// result of the query above is sorted by broadcast_before_block_num in ascending order
	// this map de-duplicates tx attempts so that only most recent attempt tx attempt for the given tx remains
	// TODO: current implementation does not guarantee that the tx hash is on the canonical chain
	recentLgps := make(map[string]DbLegacyGaslessTxPlus)
	for _, dbLgp := range dbLgps {
		if _, ok := recentLgps[dbLgp.ID]; ok {
			// found a transaction with multiple attempts
			o.lggr.Debugw("found a gasless transaction with multiple attempts", "RequestID", dbLgp.ID, "txHash", dbLgp.TxHash)
		}
		recentLgps[dbLgp.ID] = dbLgp
	}

	var dedupedLgps []types.LegacyGaslessTxPlus
	for _, dbLgp := range recentLgps {
		lgp, err := toLegacyGaslessTxPlus(dbLgp)
		if err != nil {
			return nil, err
		}
		dedupedLgps = append(dedupedLgps, *lgp)
	}

	return dedupedLgps, nil
}

func toLegacyGaslessTxPlus(dbLgp DbLegacyGaslessTxPlus) (*types.LegacyGaslessTxPlus, error) {
	lgTx, err := toLegacyGaslessTx(dbLgp.DbLegacyGaslessTx)
	if err != nil {
		return nil, err
	}
	var status *uint64
	if dbLgp.Receipt != nil {
		status = &dbLgp.Receipt.Status
	}
	return &types.LegacyGaslessTxPlus{
		LegacyGaslessTx: *lgTx,
		EthTxStatus:     string(dbLgp.EthTxStatus),
		EthTxHash:       dbLgp.EthTxHash,
		EthTxError:      dbLgp.EthTxError,
		ReceiptStatus:   status,
	}, nil
}
