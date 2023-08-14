package legacygasstation

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var _ ORM = &orm{}

type orm struct {
	q    pg.Q
	lggr logger.Logger
}

// NewORM creates an ORM scoped to chainID.
// TODO: implement pruning logic if needed
func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) ORM {
	namedLogger := lggr.Named("LegacyGasStation")
	q := pg.NewQ(db, namedLogger, cfg)
	return &orm{
		q:    q,
		lggr: namedLogger,
	}
}

func (o *orm) SelectBySourceChainIDAndStatus(sourceChainID uint64, status types.Status, qopts ...pg.QOpt) (txs []types.LegacyGaslessTx, err error) {
	q := o.q.WithOpts(qopts...)
	err = q.Select(&txs, `
		SELECT * FROM legacy_gasless_txs 
			WHERE legacy_gasless_txs.source_chain_id = $1 
			AND legacy_gasless_txs.tx_status = $2
		`, sourceChainID, status.String())
	return
}

func (o *orm) SelectByDestChainIDAndStatus(destChainID uint64, status types.Status, qopts ...pg.QOpt) (txs []types.LegacyGaslessTx, err error) {
	q := o.q.WithOpts(qopts...)
	err = q.Select(&txs, `
		SELECT * FROM legacy_gasless_txs
			WHERE legacy_gasless_txs.destination_chain_id = $1 
			AND legacy_gasless_txs.tx_status = $2
		`, destChainID, status.String())
	return
}

// InsertLegacyGaslessTx is idempotent
func (o *orm) InsertLegacyGaslessTx(tx types.LegacyGaslessTx, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	err := q.ExecQ(`INSERT INTO legacy_gasless_txs (legacy_gasless_tx_id, forwarder_address, from_address, target_address, receiver_address, nonce, amount, source_chain_id, destination_chain_id, valid_until_time, tx_signature, tx_status, token_name, token_version, eth_tx_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW(), NOW())`,
		tx.ID,
		tx.Forwarder,
		tx.From,
		tx.Target,
		tx.Receiver,
		tx.Nonce,
		tx.Amount,
		tx.SourceChainID,
		tx.DestinationChainID,
		tx.ValidUntilTime,
		tx.Signature[:],
		tx.Status.String(),
		tx.TokenName,
		tx.TokenVersion,
		tx.EthTxID,
	)
	return err
}

// UpdateLegacyGaslessTx updates legacy gasless transaction with status, ccip message ID (optional), failure reason (optional)
func (o *orm) UpdateLegacyGaslessTx(tx types.LegacyGaslessTx, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	_, err := q.Exec(`UPDATE legacy_gasless_txs SET 
	tx_status = $2,
	ccip_message_id = $3,
	failure_reason = $4,
	tx_hash = $5,
	updated_at = NOW()
	WHERE legacy_gasless_tx_id = $1`,
		tx.ID,
		tx.Status.String(),
		tx.CCIPMessageID,
		tx.FailureReason,
		tx.TxHash,
	)
	return err
}

func (o *orm) SelectBySourceChainIDAndEthTxStates(sourceChainID uint64, states []txmgrtypes.TxState, qopts ...pg.QOpt) ([]types.LegacyGaslessTxPlus, error) {
	var lgps []types.LegacyGaslessTxPlus
	q := o.q.WithOpts(qopts...)
	err := q.Select(&lgps, `SELECT 
		lgt.*,
		etx.state as etx_state,
		eta.hash as etx_hash,
		etx.error as etx_error
	FROM legacy_gasless_txs lgt
	LEFT JOIN eth_txes etx ON etx.id = lgt.eth_tx_id
	LEFT JOIN eth_tx_attempts eta ON etx.id = eta.eth_tx_id
	WHERE lgt.source_chain_id = $1
	AND etx.state = any($2)
	ORDER BY eta.broadcast_before_block_num ASC
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
