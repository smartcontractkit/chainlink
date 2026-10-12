package registrysyncer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

type ORM interface {
	AddRegistryMetadata(ctx context.Context, registryMetadata *registry.RegistryMetadata) error
	LatestRegistryMetadata(ctx context.Context) (*registry.RegistryMetadata, error)
}

type orm struct {
	ds   sqlutil.DataSource
	lggr logger.Logger
}

var _ ORM = (*orm)(nil)

func NewORM(ds sqlutil.DataSource, lggr logger.Logger) orm {
	namedLogger := logger.Named(lggr, "RegistrySyncerORM")
	return orm{
		ds:   ds,
		lggr: namedLogger,
	}
}

func (orm orm) AddRegistryMetadata(ctx context.Context, registryMetadata *registry.RegistryMetadata) error {
	orm.lggr.Debugw("Adding local registry to DB...")
	return sqlutil.TransactDataSource(ctx, orm.ds, nil, func(tx sqlutil.DataSource) error {
		localRegistryJSON, err := registryMetadata.MarshalJSON()
		if err != nil {
			return err
		}
		hash := sha256.Sum256(localRegistryJSON)
		// update if and only if the hash does not match the latest value
		r, err := tx.ExecContext(
			ctx,
			`INSERT INTO registry_syncer_states (data, data_hash) 
            SELECT $1, $2 
            WHERE $2 NOT IN (
                SELECT data_hash FROM registry_syncer_states 
                ORDER BY id DESC LIMIT 1
            )`,
			localRegistryJSON, hex.EncodeToString(hash[:]),
		)
		if err != nil {
			return fmt.Errorf("failed to insert into registry_syncer: %w", err)
		}

		n, _ := r.RowsAffected()
		if n != 0 {
			id, _ := r.LastInsertId()
			orm.lggr.Debugw("Inserted new local registry", "id", id, "hash", hex.EncodeToString(hash[:]), "registry", registryMetadata)
		} else {
			orm.lggr.Debugw("No rows affected, local registry updated. ", "hash", hex.EncodeToString(hash[:]))
		}
		_, err = tx.ExecContext(ctx, `DELETE FROM registry_syncer_states
WHERE data_hash NOT IN (
    SELECT data_hash FROM registry_syncer_states
    ORDER BY id DESC
    LIMIT 10
);`)
		return err
	})
}

func (orm orm) LatestRegistryMetadata(ctx context.Context) (*registry.RegistryMetadata, error) {
	var localRegistry registry.RegistryMetadata
	var localRegistryJSON string
	err := orm.ds.GetContext(ctx, &localRegistryJSON, `SELECT data FROM registry_syncer_states ORDER BY id DESC LIMIT 1`)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256([]byte(localRegistryJSON))
	err = localRegistry.UnmarshalJSON([]byte(localRegistryJSON))
	if err != nil {
		return nil, err
	}
	//nolint:govet // copylocks: logging the value copies RegistryMetadata's embedded RWMutex, which is not held here
	orm.lggr.Debugw("Fetched latest local registry from DB", "hash", hex.EncodeToString(hash[:]), "registry", localRegistry)

	return &localRegistry, nil
}
