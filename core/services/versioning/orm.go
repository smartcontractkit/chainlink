package versioning

import (
	"time"

	"github.com/smartcontractkit/sqlx"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	FindLatestNodeVersion() (*NodeVersion, error)
	UpsertNodeVersion(version NodeVersion) error
}

type orm struct {
	db *sqlx.DB
}

func NewORM(db *sqlx.DB) *orm {
	return &orm{
		db: db,
	}
}

// UpsertNodeVersion inserts a new NodeVersion
func (o *orm) UpsertNodeVersion(version NodeVersion) error {
	now := time.Now()

	stmt := `
INSERT INTO node_versions (version, created_at)
VALUES ($1, $2)
ON CONFLICT
DO NOTHING
`

	_, err := o.db.Exec(stmt, version.Version, now)
	if err != nil {
		return err
	}

	return nil
}

// FindLatestNodeVersion looks up the latest node version
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
