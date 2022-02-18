package evm

import (
	"database/sql"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ OCRContractTrackerDB = &contractDB{}

type contractDB struct {
	*sql.DB
	oracleSpecID int32
	lggr         logger.Logger
}

// NewDB returns a new DB scoped to this oracleSpecID
func NewContractDB(sqldb *sql.DB, oracleSpecID int32, lggr logger.Logger) *contractDB {
	return &contractDB{sqldb, oracleSpecID, lggr}
}

func (d *contractDB) SaveLatestRoundRequested(tx pg.Queryer, rr ocr2aggregator.OCR2AggregatorRoundRequested) error {
	rawLog, err := json.Marshal(rr.Raw)
	if err != nil {
		return errors.Wrap(err, "could not marshal log as JSON")
	}
	_, err = tx.Exec(`
INSERT INTO offchainreporting2_latest_round_requested (offchainreporting2_oracle_spec_id, requester, config_digest, epoch, round, raw)
VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (offchainreporting2_oracle_spec_id) DO UPDATE SET
	requester = EXCLUDED.requester,
	config_digest = EXCLUDED.config_digest,
	epoch = EXCLUDED.epoch,
	round = EXCLUDED.round,
	raw = EXCLUDED.raw
`, d.oracleSpecID, rr.Requester, rr.ConfigDigest[:], rr.Epoch, rr.Round, rawLog)

	return errors.Wrap(err, "could not save latest round requested")
}

func (d *contractDB) LoadLatestRoundRequested() (ocr2aggregator.OCR2AggregatorRoundRequested, error) {
	rr := ocr2aggregator.OCR2AggregatorRoundRequested{}
	rows, err := d.Query(`
SELECT requester, config_digest, epoch, round, raw
FROM offchainreporting2_latest_round_requested
WHERE offchainreporting2_oracle_spec_id = $1
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
