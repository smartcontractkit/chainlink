package versioning

import (
	"database/sql"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
)

// Version ORM manages the node_versions table
// NOTE: If you just need the current application version, consider using static.Version instead
// The database version is ONLY useful for managing versioning specific to the database e.g. for backups or migrations

type ORM interface {
	FindLatestNodeVersion() (*NodeVersion, error)
	UpsertNodeVersion(version NodeVersion) error
}

type orm struct {
	db *sqlx.DB
}

func NewORM(db *sqlx.DB, lggr logger.Logger) *orm {
	return &orm{
		db: db,
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

	return postgres.SqlxTransactionWithDefaultCtx(o.db, func(tx postgres.Queryer) error {
		if err := CheckVersion(tx, logger.NullLogger, version.Version); err != nil {
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
func CheckVersion(q postgres.Queryer, lggr logger.Logger, appVersion string) error {
	lggr = lggr.Named("Version")
	var dbVersion string
	err := q.QueryRowx(`SELECT version FROM node_versions ORDER BY created_at DESC LIMIT 1 FOR UPDATE`).Scan(&dbVersion)
	if errors.Is(err, sql.ErrNoRows) {
		lggr.Debug("No previous version set")
		return nil
	} else if err != nil {
		pqErr, ok := err.(*pgconn.PgError)
		if ok && pqErr.Code == "42P01" && pqErr.Message == `relation "node_versions" does not exist` {
			lggr.Debug("Previous version not set; node_versions table does not exist")
			return nil
		}
		return err
	}

	dbv, dberr := semver.NewVersion(dbVersion)
	appv, apperr := semver.NewVersion(appVersion)
	if dberr != nil {
		lggr.Warnf("Database version %q is not valid semver; skipping version check", dbVersion)
		return nil
	}
	if apperr != nil {
		return errors.Errorf("Application version %q is not valid semver", appVersion)
	}
	if dbv.GreaterThan(appv) {
		return errors.Errorf("Application version (%s) is older than database version (%s). Only Chainlink %s or later can be run on this database", appv, dbv, dbv)
	}
	return nil
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
