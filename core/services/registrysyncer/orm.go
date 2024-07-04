package registrysyncer

import (
	"context"
	"crypto/sha256"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type syncerORM struct {
	ds   sqlutil.DataSource
	lggr logger.Logger
}

func newORM(ds sqlutil.DataSource, lggr logger.Logger) syncerORM {
	namedLogger := lggr.Named("RegistrySyncerORM")
	return syncerORM{
		ds:   ds,
		lggr: namedLogger,
	}
}

func (orm syncerORM) addState(ctx context.Context, stateJSON string) error {
	hash := sha256.Sum256([]byte(stateJSON))
	_, err := orm.ds.ExecContext(
		ctx,
		`INSERT INTO registry_syncer_states (data, data_hash) VALUES ($1, $2) ON CONFLICT (data_hash) DO NOTHING`,
		stateJSON, hash[:],
	)
	return err
}

func (orm syncerORM) latestState(ctx context.Context) (*State, error) {
	var state State
	var stateJSON string
	err := orm.ds.GetContext(ctx, &stateJSON, `SELECT data FROM registry_syncer_states ORDER BY created_at DESC LIMIT 1`)
	if err != nil {
		return nil, err
	}
	err = state.UnmarshalJSON([]byte(stateJSON))
	if err != nil {
		return nil, err
	}
	return &state, nil
}
