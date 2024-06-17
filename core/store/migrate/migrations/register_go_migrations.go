package migrations

import "github.com/pressly/goose/v3"

func init() {
	RegisterCustomMigrations()
}

func RegisterCustomMigrations() {
	goose.ResetGlobalMigrations() // reset global state that is mutated by both core and plugin migrations
	// add by name so we can run this function multiple times
	// in order for this name lookup to work, the migrations must be read from a filesystem object rooted at this package, eg an embed.FS
	goose.AddNamedMigrationContext("0036_external_job_id.go", Up36, Down36)
	goose.AddNamedMigrationContext("0054_remove_legacy_pipeline.go", Up54, Down54)
	goose.AddNamedMigrationContext("0056_multichain.go", Up56, Down56)
	goose.AddNamedMigrationContext("0195_add_not_null_to_evm_chain_id_in_job_specs.go", Up195, Down195)
}
