package ocrcommon

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	ocrnetworking "github.com/smartcontractkit/libocr/networking/types"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

var _ ocrnetworking.DiscovererDatabase = &DiscovererDatabase{}

const (
	// ocrDiscovererTable is the name of the table used to store OCR announcements
	ocrDiscovererTable = "ocr_discoverer_announcements"
	// don2donDiscovererTable is the name of the table used to store DON2DON announcements
	don2donDiscovererTable = "don2don_discoverer_announcements"
)

// DiscovererDatabase is a key-value store for p2p announcements
// that are based on the RageP2P library and bootstrap nodes
type DiscovererDatabase struct {
	ds        sqlutil.DataSource
	peerID    string
	tableName string
}

// NewOCRDiscovererDatabase creates a new DiscovererDatabase for OCR announcements
func NewOCRDiscovererDatabase(ds sqlutil.DataSource, peerID string) *DiscovererDatabase {
	return &DiscovererDatabase{
		ds:        ds,
		peerID:    peerID,
		tableName: ocrDiscovererTable,
	}
}

// NewDON2DONDiscovererDatabase creates a new DiscovererDatabase for DON2DON announcements
func NewDON2DONDiscovererDatabase(ds sqlutil.DataSource, peerID string) *DiscovererDatabase {
	return &DiscovererDatabase{
		ds:        ds,
		peerID:    peerID,
		tableName: don2donDiscovererTable,
	}
}

// StoreAnnouncement has key-value-store semantics and stores a peerID (key) and an associated serialized
// announcement (value).
func (d *DiscovererDatabase) StoreAnnouncement(ctx context.Context, peerID string, ann []byte) error {
	q := fmt.Sprintf(`
INSERT INTO %s (local_peer_id, remote_peer_id, ann, created_at, updated_at)
VALUES ($1,$2,$3,NOW(),NOW()) ON CONFLICT (local_peer_id, remote_peer_id) DO UPDATE SET
ann = EXCLUDED.ann,
updated_at = EXCLUDED.updated_at
;`, d.tableName)

	_, err := d.ds.ExecContext(ctx,
		q, d.peerID, peerID, ann)
	return errors.Wrap(err, "DiscovererDatabase failed to StoreAnnouncement")
}

// ReadAnnouncements returns one serialized announcement (if available) for each of the peerIDs in the form of a map
// keyed by each announcement's corresponding peer ID.
func (d *DiscovererDatabase) ReadAnnouncements(ctx context.Context, peerIDs []string) (results map[string][]byte, err error) {
	q := fmt.Sprintf(`SELECT remote_peer_id, ann FROM %s WHERE remote_peer_id = ANY($1) AND local_peer_id = $2`, d.tableName)

	rows, err := d.ds.QueryContext(ctx, q, pq.Array(peerIDs), d.peerID)
	if err != nil {
		return nil, errors.Wrap(err, "DiscovererDatabase failed to ReadAnnouncements")
	}
	defer func() { err = multierr.Combine(err, rows.Close()) }()
	results = make(map[string][]byte)
	for rows.Next() {
		var peerID string
		var ann []byte
		err = rows.Scan(&peerID, &ann)
		if err != nil {
			return
		}
		results[peerID] = ann
	}
	if err = rows.Err(); err != nil {
		return
	}
	return results, nil
}
