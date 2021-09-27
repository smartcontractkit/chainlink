package migrations

import (
	"gorm.io/gorm"
)

const up53 = `
	CREATE TABLE web_authns (
		"id" BIGSERIAL PRIMARY KEY,
		"email" text NOT NULL,
		"public_key_data" jsonb NOT NULL,
		CONSTRAINT fk_email
			FOREIGN KEY(email)
			REFERENCES users(email)
	);

	CREATE UNIQUE INDEX web_authns_email_idx ON web_authns (email);
`

const down53 = `
	DROP TABLE IF EXISTS web_authns;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0053_add_web_authns_table",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up53).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down53).Error
		},
	})
}
