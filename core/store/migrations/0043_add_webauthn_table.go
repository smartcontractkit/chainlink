package migrations

import (
	"gorm.io/gorm"
)

const up43 = `
	CREATE TABLE web_authns (
		"id" BIGSERIAL PRIMARY KEY, 
		"email" text NOT NULL,
		"public_key_data" text NOT NULL,
		"settings" text NOT NULL,
		CONSTRAINT fk_email
			FOREIGN KEY(email)
			REFERENCES users(email)
	);
`

const down43 = `
	DROP TABLE IF EXISTS web_authns;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0043_add_web_authns_table",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up43).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down43).Error
		},
	})
}
