package evm

import (
	"database/sql"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// RequestRoundDB stores requested rounds for querying by the median plugin.
type RequestRoundDB interface {
	SaveLatestRoundRequested(tx pg.Queryer, rr ocr2aggregator.OCR2AggregatorRoundRequested) error
	LoadLatestRoundRequested() (rr ocr2aggregator.OCR2AggregatorRoundRequested, err error)
}

var _ RequestRoundDB = &requestRoundDB{}

//go:generate mockery --quiet --name RequestRoundDB --output ./mocks/ --case=underscore
type requestRoundDB struct {
	*sql.DB
	oracleSpecID int32
	lggr         logger.Logger
}

// NewDB returns a new DB scoped to this oracleSpecID
func NewRoundRequestedDB(sqldb *sql.DB, oracleSpecID int32, lggr logger.Logger) *requestRoundDB {
	return &requestRoundDB{sqldb, oracleSpecID, lggr}
}

func (d *requestRoundDB) SaveLatestRoundRequested(tx pg.Queryer, rr ocr2aggregator.OCR2AggregatorRoundRequested) error {
	rawLog, err := json.Marshal(rr.Raw)
	if err != nil {
		return errors.Wrap(err, "could not marshal log as JSON")
	}
	_, err = tx.Exec(`
INSERT INTO ocr2_latest_round_requested (ocr2_oracle_spec_id, requester, config_digest, epoch, round, raw)
VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (ocr2_oracle_spec_id) DO UPDATE SET
	requester = EXCLUDED.requester,
	config_digest = EXCLUDED.config_digest,
	epoch = EXCLUDED.epoch,
	round = EXCLUDED.round,
	raw = EXCLUDED.raw
`, d.oracleSpecID, rr.Requester, rr.ConfigDigest[:], rr.Epoch, rr.Round, rawLog)

	return errors.Wrap(err, "could not save latest round requested")
}

func (d *requestRoundDB) LoadLatestRoundRequested() (ocr2aggregator.OCR2AggregatorRoundRequested, error) {
	rr := ocr2aggregator.OCR2AggregatorRoundRequested{}
	rows, err := d.Query(`
SELECT requester, config_digest, epoch, round, raw
FROM ocr2_latest_round_requested
WHERE ocr2_oracle_spec_id = $1
LIMIT 1
`, d.oracleSpecID)
	if err != nil {
		return rr, errors.Wrap(err, "LoadLatestRoundRequested failed to query rows")
	}

	for rows.Next() {
		var configDigest []byte
		var rawLog []byte

		err = rows.Scan(&rr.Requester, &configDigest, &rr.Epoch, &rr.Round, &rawLog)
		if err != nil {
			return rr, errors.Wrap(err, "LoadLatestRoundRequested failed to scan row")
		}

		rr.ConfigDigest, err = ocrtypes.BytesToConfigDigest(configDigest)
		if err != nil {
			return rr, errors.Wrap(err, "LoadLatestRoundRequested failed to decode config digest")
		}

		err = json.Unmarshal(rawLog, &rr.Raw)
		if err != nil {
			return rr, errors.Wrap(err, "LoadLatestRoundRequested failed to unmarshal raw log")
		}
	}

	if err = rows.Err(); err != nil {
		return rr, err
	}

	return rr, nil
}
