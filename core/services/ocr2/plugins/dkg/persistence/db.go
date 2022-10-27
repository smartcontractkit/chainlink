package persistence

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"github.com/smartcontractkit/ocr2vrf/types/hash"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

var (
	_        ocr2vrftypes.DKGSharePersistence = &shareDB{}
	zeroHash hash.Hash
)

type shareDB struct {
	q    pg.Q
	lggr logger.Logger
}

// NewShareDB creates a new DKG share database.
func NewShareDB(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) ocr2vrftypes.DKGSharePersistence {
	return &shareDB{
		q:    pg.NewQ(db, lggr, cfg),
		lggr: lggr,
	}
}

// WriteShareRecords writes the provided (already encrypted)
// share records to the Chainlink database.
func (s *shareDB) WriteShareRecords(
	ctx context.Context,
	cfgDgst ocrtypes.ConfigDigest,
	keyID [32]byte,
	shareRecords []ocr2vrftypes.PersistentShareSetRecord,
) error {
	lggr := s.lggr.With(
		"configDigest", hexutil.Encode(cfgDgst[:]),
		"keyID", hexutil.Encode(keyID[:]))

	start := time.Now()
	defer func() {
		lggr.Infow("Inserted DKG shares into DB", "duration", time.Since(start))
	}()

	var named []dkgShare
	for _, record := range shareRecords {
		if bytes.Equal(record.Hash[:], zeroHash[:]) {
			lggr.Warnw("skipping record with zero hash",
				"player", record.Dealer.String(),
				"hash", hexutil.Encode(record.Hash[:]),
			)
			continue
		}

		// XXX: this might be expensive, but is a good sanity check.
		localHash := hash.GetHash(record.MarshaledShareRecord)
		if !bytes.Equal(record.Hash[:], localHash[:]) {
			return fmt.Errorf("local hash doesn't match given hash in record, expected: %x, got: %x",
				localHash[:], record.Hash[:])
		}

		var h hash.Hash
		if copied := copy(h[:], record.Hash[:]); copied != 32 {
			return fmt.Errorf("wrong number of bytes copied in hash (dealer:%s) %x: %d",
				record.Dealer.String(), record.Hash[:], copied)
		}

		named = append(named, dkgShare{
			ConfigDigest:         cfgDgst[:],
			KeyID:                keyID[:],
			Dealer:               record.Dealer.Marshal(),
			MarshaledShareRecord: record.MarshaledShareRecord,
			/* TODO/WTF: can't do "record.Hash[:]": this leads to store the last record's hash for all the records! */
			RecordHash: h[:],
		})
	}

	if len(named) == 0 {
		lggr.Infow("No valid share records to insert")
		return nil
	}

	lggr.Infow("Inserting DKG shares into DB",
		"shareHashes", shareHashes(shareRecords),
		"numRecords", len(shareRecords),
		"numNamed", len(named))

	// Always upsert because we want the number of rows in the table to match
	// the number of members of the committee.
	query := `
INSERT INTO dkg_shares (config_digest, key_id, dealer, marshaled_share_record, record_hash)
VALUES (:config_digest, :key_id, :dealer, :marshaled_share_record, :record_hash)
ON CONFLICT ON CONSTRAINT dkg_shares_pkey
DO UPDATE SET marshaled_share_record = EXCLUDED.marshaled_share_record, record_hash = EXCLUDED.record_hash
`
	return s.q.ExecQNamed(query, named[:])
}

// ReadShareRecords retrieves any share records in the database that correspond
// to the provided config digest and DKG key ID.
func (s *shareDB) ReadShareRecords(
	cfgDgst ocrtypes.ConfigDigest,
	keyID [32]byte,
) (
	retrievedShares []ocr2vrftypes.PersistentShareSetRecord,
	err error,
) {
	lggr := s.lggr.With(
		"configDigest", hexutil.Encode(cfgDgst[:]),
		"keyID", hexutil.Encode(keyID[:]))

	start := time.Now()
	defer func() {
		lggr.Debugw("Finished reading DKG shares from DB", "duration", time.Since(start))
	}()

	a := map[string]any{
		"config_digest": cfgDgst[:],
		"key_id":        keyID[:],
	}
	query, args, err := sqlx.Named(
		`
SELECT *
FROM dkg_shares
WHERE config_digest = :config_digest
	AND key_id = :key_id
`, a)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx Named")
	}
	query = s.q.Rebind(query)
	var dkgShares []dkgShare
	err = s.q.Select(&dkgShares, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	for _, share := range dkgShares {
		playerIdx, _, err := ocr2vrftypes.UnmarshalPlayerIdx(share.Dealer)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshalling %x", share.Dealer)
		}
		var h hash.Hash
		if copied := copy(h[:], share.RecordHash); copied != 32 {
			return nil, fmt.Errorf("wrong number of bytes copied in hash %x: %d", share.RecordHash, copied)
		}

		// NOTE: no integrity check on share.MarshaledShareRecord
		// because caller will do it anyways, so it'd be wasteful.
		retrievedShares = append(retrievedShares, ocr2vrftypes.PersistentShareSetRecord{
			Dealer:               *playerIdx,
			MarshaledShareRecord: share.MarshaledShareRecord,
			Hash:                 h,
		})
	}

	lggr.Debugw("Read DKG shares from DB",
		"shareRecords", shareHashes(retrievedShares),
		"numRecords", len(dkgShares),
	)

	return retrievedShares, nil
}

func shareHashes(shareRecords []ocr2vrftypes.PersistentShareSetRecord) []string {
	r := make([]string, len(shareRecords))
	for i, record := range shareRecords {
		r[i] = hexutil.Encode(record.Hash[:])
	}
	return r
}
