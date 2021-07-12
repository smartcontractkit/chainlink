package migrations

import (
	"gorm.io/gorm"
)

const up25 = `
	CREATE TYPE log_level AS ENUM (
		'debug',
		'info',
		'warn',
		'error',
		'panic'
	);
	
	CREATE TABLE log_configs (
		"id" BIGSERIAL PRIMARY KEY,
		"service_name" text NOT NULL UNIQUE,
		"log_level" log_level NOT NULL,
		"created_at" timestamp with time zone,
		"updated_at" timestamp with time zone
	);
`

const down25 = `
	DROP TABLE IF EXISTS log_configs;

	DROP TYPE IF EXISTS log_level;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0025_create_log_config_table",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up25).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down25).Error
		},
	})
}
