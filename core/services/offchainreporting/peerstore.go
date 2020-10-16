package offchainreporting

import (
	"context"
	"database/sql"

	sqlds "github.com/ipfs/go-ds-sql"
	pgqueries "github.com/ipfs/go-ds-sql/postgres"
	p2ppeerstore "github.com/libp2p/go-libp2p-core/peerstore"
	p2ppeerstoreds "github.com/libp2p/go-libp2p-peerstore/pstoreds"
)

const tableName = "p2p_peerstore"

// NewPeerstore creates a new database-backed peerstore
// NOTE: You can get sql.DB from store with store.DB.DB()
func NewPeerstore(ctx context.Context, db *sql.DB) (p2ppeerstore.Peerstore, error) {
	queries := pgqueries.NewQueries(tableName)
	datastore := sqlds.NewDatastore(db, queries)
	opts := p2ppeerstoreds.DefaultOpts()
	return p2ppeerstoreds.NewPeerstore(ctx, datastore, opts)
}
