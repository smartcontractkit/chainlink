package web2

import (
	"database/sql"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
)

type orm struct {
	chainID *big.Int
	q       pg.Q
}

func newORM(chainID *big.Int, db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) *orm {
	q := pg.NewQ(db, lggr.Named("VRFWeb2ORM"), cfg)
	return &orm{
		chainID: chainID,
		q:       q,
	}
}

func (o *orm) InsertRequest(
	clientRequestID []byte,
	lotteryType uint8,
	vrfExternalRequestID []byte,
	lotteryAddress common.Address,
	requestTxHash common.Hash,
	qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	return q.ExecQ(
		`
INSERT INTO vrf_web2_requests (client_request_id, lottery_type, vrf_external_request_id, lottery_contract_address, evm_chain_id, request_tx_hash)
VALUES ($1, $2, $3, $4, $5, $6)
`, clientRequestID, lotteryType, vrfExternalRequestID, lotteryAddress, utils.NewBig(o.chainID), requestTxHash,
	)
}

func (o *orm) InsertFulfillment(
	clientRequestID []byte,
	lotteryType uint8,
	vrfExternalRequestID []byte,
	lotteryAddress common.Address,
	winningNumbers []uint8,
	fulfillmentTxHash common.Hash,
	qopts ...pg.QOpt,
) error {
	q := o.q.WithOpts(qopts...)
	return q.ExecQ(
		`
INSERT INTO vrf_web2_fulfillments (client_request_id, lottery_type, vrf_external_request_id, lottery_contract_address, evm_chain_id, winning_numbers, fulfillment_tx_hash)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`, clientRequestID, lotteryType, vrfExternalRequestID, lotteryAddress, utils.NewBig(o.chainID), winningNumbers, fulfillmentTxHash,
	)
}

func (o *orm) GetFulfillment(
	clientRequestID []byte,
	lotteryType uint8,
	qopts ...pg.QOpt,
) (winningNumbers []uint8, err error) {
	q := o.q.WithOpts(qopts...)
	rows, err := q.Query(
		`
SELECT winning_numbers
FROM vrf_web2_fulfillments
WHERE client_request_id = $1
	AND lottery_type = $2
	AND vrf_external_request_id = $3
`,
	)
	if err == sql.ErrNoRows {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&winningNumbers)
		if err != nil {
			return nil, err
		} else {
			break // should only be one TODO: confirm
		}
	}

	return
}
