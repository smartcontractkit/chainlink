package legacygasstation

import (
	"context"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation"
	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var _ legacygasstation.ORM = &orm{}

type orm struct {
	q    pg.Q
	lggr logger.Logger
}

const InsertLegacyGaslessTxQuery = `INSERT INTO legacy_gasless_txs (legacy_gasless_tx_id, forwarder_address, from_address, target_address, receiver_address, nonce, amount, source_chain_id, destination_chain_id, valid_until_time, tx_signature, tx_status, token_name, token_version, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW(), NOW())`

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

func (o *orm) SelectBySourceChainIDAndStatus(ctx context.Context, sourceChainID uint64, status types.Status) (txs []types.LegacyGaslessTx, err error) {
	q := o.q.WithOpts(pg.WithParentCtx(ctx))
	err = q.Select(&txs, `
		SELECT * FROM legacy_gasless_txs 
			WHERE legacy_gasless_txs.source_chain_id = $1 
			AND legacy_gasless_txs.tx_status = $2
		`, sourceChainID, status.String())
	return
}

func (o *orm) SelectByDestChainIDAndStatus(ctx context.Context, destChainID uint64, status types.Status) (txs []types.LegacyGaslessTx, err error) {
	q := o.q.WithOpts(pg.WithParentCtx(ctx))
	err = q.Select(&txs, `
		SELECT * FROM legacy_gasless_txs
			WHERE legacy_gasless_txs.destination_chain_id = $1 
			AND legacy_gasless_txs.tx_status = $2
		`, destChainID, status.String())
	return
}

// InsertLegacyGaslessTx is idempotent
func (o *orm) InsertLegacyGaslessTx(ctx context.Context, lgsTx types.LegacyGaslessTx) error {
	q := o.q.WithOpts(pg.WithParentCtx(ctx))
	return q.ExecQ(InsertLegacyGaslessTxQuery,
		lgsTx.ID,
		lgsTx.Forwarder,
		lgsTx.From,
		lgsTx.Target,
		lgsTx.Receiver,
		lgsTx.Nonce,
		lgsTx.Amount,
		lgsTx.SourceChainID,
		lgsTx.DestinationChainID,
		lgsTx.ValidUntilTime,
		lgsTx.Signature[:],
		lgsTx.Status.String(),
		lgsTx.TokenName,
		lgsTx.TokenVersion,
	)
}

// UpdateLegacyGaslessTx updates legacy gasless transaction with status, ccip message ID (optional), failure reason (optional)
func (o *orm) UpdateLegacyGaslessTx(ctx context.Context, lgsTx types.LegacyGaslessTx) error {
	q := o.q.WithOpts(pg.WithParentCtx(ctx))
	return q.ExecQ(UpdateLegacyGaslessTxQuery,
		lgsTx.ID,
		lgsTx.Status.String(),
		lgsTx.CCIPMessageID,
		lgsTx.FailureReason,
		lgsTx.TxHash,
	)
}

func (o *orm) SelectBySourceChainIDAndEthTxStates(ctx context.Context, sourceChainID uint64, states []legacygasstation.EtxStatus) ([]types.LegacyGaslessTxPlus, error) {
	var lgps []types.LegacyGaslessTxPlus
	q := o.q.WithOpts(pg.WithParentCtx(ctx))
	err := q.Select(&lgps, `SELECT 
		lgt.*,
		etx.state as etx_state,
		eta.hash as etx_hash,
		etx.error as etx_error,
		etr.receipt.status as receipt_status
	FROM legacy_gasless_txs lgt
	LEFT JOIN eth_txes etx ON etx.id = lgt.eth_tx_id
	LEFT JOIN eth_tx_attempts eta ON etx.id = eta.eth_tx_id
	LEFT JOIN eth_receipts etr ON eta.hash = etr.tx_hash
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
	recentLgps := make(map[string]types.LegacyGaslessTxPlus)
	for _, lgp := range lgps {
		if _, ok := recentLgps[lgp.ID]; ok {
			// found a transaction with multiple attempts
			o.lggr.Debugw("found a gasless transaction with multiple attempts", "RequestID", lgp.ID, "txHash", lgp.TxHash)
		}
		recentLgps[lgp.ID] = lgp
	}

	var dedupedLgps []types.LegacyGaslessTxPlus
	for _, lgp := range recentLgps {
		dedupedLgps = append(dedupedLgps, lgp)
	}

	return dedupedLgps, nil
}
