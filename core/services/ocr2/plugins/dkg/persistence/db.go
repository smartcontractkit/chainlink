package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"github.com/smartcontractkit/ocr2vrf/types/hash"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

var _ ocr2vrftypes.DKGSharePersistence = &shareDB{}

type shareDB struct {
	q pg.Q
}

// NewShareDB creates a new DKG share database.
func NewShareDB(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) ocr2vrftypes.DKGSharePersistence {
	return &shareDB{
		q: pg.NewQ(db, lggr, cfg),
	}
}

func (s *shareDB) WriteShareRecords(
	ctx context.Context,
	cfgDgst ocrtypes.ConfigDigest,
	keyID [32]byte,
	shareRecords []ocr2vrftypes.PersistentShareSetRecord,
) error {
	named := make([]dkgShare, len(shareRecords))
	for i, record := range shareRecords {
		var h hash.Hash
		copy(h[:], record.Hash[:])
		named[i] = dkgShare{
			ConfigDigest:         cfgDgst[:],
			KeyID:                keyID[:],
			Dealer:               record.Dealer.Marshal(),
			MarshaledShareRecord: record.MarshaledShareRecord,
			RecordHash:           h[:], /* TODO/WTF: can't do "record.Hash[:]": this leads to store the last record's hash for all the records!*/
		}
	}
	q := s.q.WithOpts(pg.WithParentCtx(ctx))
	return q.ExecQNamed(`
INSERT INTO dkg_shares (config_digest, key_id, dealer, marshaled_share_record, record_hash)
VALUES (:config_digest, :key_id, :dealer, :marshaled_share_record, :record_hash)
`, named[:])
}

func (s *shareDB) ReadShareRecords(
	cfgDgst ocrtypes.ConfigDigest,
	keyID [32]byte,
) (
	retrievedShares []ocr2vrftypes.PersistentShareSetRecord,
	err error,
) {
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
		retrievedShares = append(retrievedShares, ocr2vrftypes.PersistentShareSetRecord{
			Dealer:               *playerIdx,
			MarshaledShareRecord: share.MarshaledShareRecord,
			Hash:                 h,
		})
	}

	return retrievedShares, nil
}
