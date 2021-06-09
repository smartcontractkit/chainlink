package migrations

import (
	"gorm.io/gorm"
)

const up35 = `
CREATE TABLE feeds_managers (
    id BIGSERIAL PRIMARY KEY,
	name VARCHAR (255) NOT NULL,
	uri VARCHAR (255) NOT NULL,
	public_key bytea CHECK (octet_length(public_key) = 32) NOT NULL UNIQUE,
	job_types TEXT [] NOT NULL,
	network VARCHAR (100) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);
`

const down35 = `
	DROP TABLE feeds_managers
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0035_create_feeds_managers",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up35).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down35).Error
		},
	})
}
