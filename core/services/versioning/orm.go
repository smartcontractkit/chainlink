package versioning

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Version ORM manages the node_versions table
// NOTE: If you just need the current application version, consider using static.Version instead
// The database version is ONLY useful for managing versioning specific to the database e.g. for backups or migrations

type ORM interface {
	FindLatestNodeVersion(ctx context.Context) (*NodeVersion, error)
	UpsertNodeVersion(ctx context.Context, version NodeVersion) error
}

type orm struct {
	ds   sqlutil.DataSource
	lggr logger.Logger
}

func NewORM(ds sqlutil.DataSource, lggr logger.Logger) *orm {
	return &orm{
		ds:   ds,
		lggr: lggr.Named("VersioningORM"),
	}
}

// UpsertNodeVersion inserts a new NodeVersion, returning error if the DB
// version is newer than the current one
// NOTE: If you just need the current application version, consider using static.Version instead
// The database version is ONLY useful for managing versioning specific to the database e.g. for backups or migrations
func (o *orm) UpsertNodeVersion(ctx context.Context, version NodeVersion) error {
	now := time.Now()

	if _, err := semver.NewVersion(version.Version); err != nil {
		return errors.Wrapf(err, "%q is not valid semver", version.Version)
	}

	return sqlutil.TransactDataSource(ctx, o.ds, nil, func(tx sqlutil.DataSource) error {
		if env.SkipAppVersionCheck.IsTrue() {
			o.lggr.Warnw("Skipping app version check", "appVersion", version.Version)
		} else {
			if _, _, err := CheckVersion(ctx, tx, logger.NullLogger, version.Version, env.IgnorePrereleaseVersionCheck.IsTrue()); err != nil {
				return err
			}
		}

		stmt := `
INSERT INTO node_versions (version, created_at)
VALUES ($1, $2)
ON CONFLICT ((version IS NOT NULL)) DO UPDATE SET
version = EXCLUDED.version,
created_at = EXCLUDED.created_at
`

		_, err := tx.ExecContext(ctx, stmt, version.Version, now)
		return err
	})
}

// CheckVersion returns an error if there is a valid semver version in the
// node_versions table that is higher than the current app version.
// If ignorePrerelease is true, pre-release information is ignored when comparing versions.
func CheckVersion(ctx context.Context, ds sqlutil.DataSource, lggr logger.Logger, appVersion string, ignorePrerelease bool) (appv, dbv *semver.Version, err error) {
	lggr = lggr.Named("Version")
	var dbVersion string
	err = ds.GetContext(ctx, &dbVersion, `SELECT version FROM node_versions ORDER BY created_at DESC LIMIT 1 FOR UPDATE`)
	if errors.Is(err, sql.ErrNoRows) {
		lggr.Debugw("No previous version set", "appVersion", appVersion)
		return nil, nil, nil
	} else if err != nil {
		var pqErr *pgconn.PgError
		ok := errors.As(err, &pqErr)
		if ok && pqErr.Code == "42P01" && pqErr.Message == `relation "node_versions" does not exist` {
			lggr.Debugw("Previous version not set; node_versions table does not exist", "appVersion", appVersion)
			return nil, nil, nil
		}
		return nil, nil, err
	}

	dbv, dberr := semver.NewVersion(dbVersion)
	appv, apperr := semver.NewVersion(appVersion)
	if dberr != nil {
		lggr.Warnf("Database version %q is not valid semver; skipping version check", dbVersion)
		return nil, nil, nil
	}
	if apperr != nil {
		return nil, nil, errors.Errorf("Application version %q is not valid semver", appVersion)
	}

	var dbvToCompare, appvToCompare *semver.Version
	if ignorePrerelease {
		dbvWithoutPre, err := dbv.SetPrerelease("")
		if err != nil {
			return nil, nil, fmt.Errorf("could not generate db version without pre-release: %w", err)
		}
		appvWithoutPre, err := appv.SetPrerelease("")
		if err != nil {
			return nil, nil, fmt.Errorf("could not generate app version without pre-release: %w", err)
		}
		dbvToCompare = &dbvWithoutPre
		appvToCompare = &appvWithoutPre
	} else {
		dbvToCompare = dbv
		appvToCompare = appv
	}

	if dbvToCompare.GreaterThan(appvToCompare) {
		return nil, nil, errors.Errorf("Application version (%s) is lower than database version (%s). Only Chainlink %s or higher can be run on this database", appv, dbv, dbv)
	}
	return appv, dbv, nil
}

// FindLatestNodeVersion looks up the latest node version
// NOTE: If you just need the current application version, consider using static.Version instead
// The database version is ONLY useful for managing versioning specific to the database e.g. for backups or migrations
func (o *orm) FindLatestNodeVersion(ctx context.Context) (*NodeVersion, error) {
	stmt := `
SELECT version, created_at
FROM node_versions
ORDER BY created_at DESC
`

	var nodeVersion NodeVersion
	err := o.ds.GetContext(ctx, &nodeVersion, stmt)
	if err != nil {
		return nil, err
	}

	return &nodeVersion, err
}
