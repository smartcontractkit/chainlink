package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const up24 = `
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

const down24 = `
	DROP TABLE IF EXISTS log_conf;

	DROP TABLE IF EXISTS log_services;
`

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "0024_create_log_config_table",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up24).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down24).Error
		},
	})
}
