package ocrcommon

import (
	"context"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	ocrnetworking "github.com/smartcontractkit/libocr/networking/types"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

var _ ocrnetworking.DiscovererDatabase = &DiscovererDatabase{}

type DiscovererDatabase struct {
	ds     sqlutil.DataSource
	peerID string
}

func NewDiscovererDatabase(ds sqlutil.DataSource, peerID string) *DiscovererDatabase {
	return &DiscovererDatabase{
		ds,
		peerID,
	}
}

// StoreAnnouncement has key-value-store semantics and stores a peerID (key) and an associated serialized
// announcement (value).
func (d *DiscovererDatabase) StoreAnnouncement(ctx context.Context, peerID string, ann []byte) error {
	_, err := d.ds.ExecContext(ctx, `
INSERT INTO ocr_discoverer_announcements (local_peer_id, remote_peer_id, ann, created_at, updated_at)
VALUES ($1,$2,$3,NOW(),NOW()) ON CONFLICT (local_peer_id, remote_peer_id) DO UPDATE SET 
ann = EXCLUDED.ann,
updated_at = EXCLUDED.updated_at
;`, d.peerID, peerID, ann)
	return errors.Wrap(err, "DiscovererDatabase failed to StoreAnnouncement")
}

// ReadAnnouncements returns one serialized announcement (if available) for each of the peerIDs in the form of a map
// keyed by each announcement's corresponding peer ID.
func (d *DiscovererDatabase) ReadAnnouncements(ctx context.Context, peerIDs []string) (results map[string][]byte, err error) {
	rows, err := d.ds.QueryContext(ctx, `
SELECT remote_peer_id, ann FROM ocr_discoverer_announcements WHERE remote_peer_id = ANY($1) AND local_peer_id = $2`, pq.Array(peerIDs), d.peerID)
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
