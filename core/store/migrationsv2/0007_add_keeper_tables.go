package migrationsv2

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const up7 = `
		CREATE TABLE public.keeper_registries (
			id SERIAL PRIMARY KEY,
			keeper_index int NOT NULL,
			reference_id uuid UNIQUE NOT NULL,
			address bytea UNIQUE NOT NULL,
			"from" bytea NOT NULL,
			check_gas int NOT NULL,
			block_count_per_turn int NOT NULL,
			job_spec_id uuid UNIQUE NOT NULL REFERENCES job_specs (id),
			num_keepers int NOT NULL
		);

		CREATE TABLE public.keeper_registrations (
			id SERIAL PRIMARY KEY,
			registry_id INT NOT NULL REFERENCES keeper_registries (id) ON DELETE CASCADE,
			execute_gas int NOT NULL,
			check_data bytea NOT NULL,
			upkeep_id bigint NOT NULL,
			positioning_constant int NOT NULL
		);

		CREATE UNIQUE INDEX idx_keeper_registrations_unique_upkeep_ids_per_keeper ON keeper_registrations(upkeep_id, registry_id);
	`

const down7 = "DROP TABLE IF EXISTS keeper_registries, keeper_registrations;"

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "0007_add_keeper_tables",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up7).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down7).Error
		},
	})
}
