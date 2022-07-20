package versioning

import (
	"database/sql"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// Version ORM manages the node_versions table
// NOTE: If you just need the current application version, consider using static.Version instead
// The database version is ONLY useful for managing versioning specific to the database e.g. for backups or migrations

type ORM interface {
	FindLatestNodeVersion() (*NodeVersion, error)
	UpsertNodeVersion(version NodeVersion) error
}

type orm struct {
	db   *sqlx.DB
	lggr logger.Logger
}

func NewORM(db *sqlx.DB, lggr logger.Logger) *orm {
	return &orm{
		db:   db,
		lggr: lggr.Named("VersioningORM"),
	}
}

// UpsertNodeVersion inserts a new NodeVersion, returning error if the DB
// version is newer than the current one
// NOTE: If you just need the current application version, consider using static.Version instead
// The database version is ONLY useful for managing versioning specific to the database e.g. for backups or migrations
func (o *orm) UpsertNodeVersion(version NodeVersion) error {
	now := time.Now()

	if _, err := semver.NewVersion(version.Version); err != nil {
		return errors.Wrapf(err, "%q is not valid semver", version.Version)
	}

	return pg.SqlxTransactionWithDefaultCtx(o.db, o.lggr, func(tx pg.Queryer) error {
		if _, _, err := CheckVersion(tx, logger.NullLogger, version.Version); err != nil {
			return err
		}

		stmt := `
INSERT INTO node_versions (version, created_at)
VALUES ($1, $2)
ON CONFLICT ((version IS NOT NULL)) DO UPDATE SET
version = EXCLUDED.version,
created_at = EXCLUDED.created_at
`

		_, err := tx.Exec(stmt, version.Version, now)
		return err
	})
}

// CheckVersion returns an error if there is a valid semver version in the
// node_versions table that is lower than the current app version
func CheckVersion(q pg.Queryer, lggr logger.Logger, appVersion string) (appv, dbv *semver.Version, err error) {
	lggr = lggr.Named("Version")
	var dbVersion string
	err = q.Get(&dbVersion, `SELECT version FROM node_versions ORDER BY created_at DESC LIMIT 1 FOR UPDATE`)
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
	if dbv.GreaterThan(appv) {
		return nil, nil, errors.Errorf("Application version (%s) is older than database version (%s). Only Chainlink %s or later can be run on this database", appv, dbv, dbv)
	}
	return appv, dbv, nil
}

// FindLatestNodeVersion looks up the latest node version
// NOTE: If you just need the current application version, consider using static.Version instead
// The database version is ONLY useful for managing versioning specific to the database e.g. for backups or migrations
func (o *orm) FindLatestNodeVersion() (*NodeVersion, error) {
	stmt := `
SELECT version, created_at
FROM node_versions
ORDER BY created_at DESC
`

	var nodeVersion NodeVersion
	err := o.db.Get(&nodeVersion, stmt)
	if err != nil {
		return nil, err
	}

	return &nodeVersion, err
}
