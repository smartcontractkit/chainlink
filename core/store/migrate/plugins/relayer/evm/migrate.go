package evm

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
	"gopkg.in/guregu/null.v4"
)

func setupPluginMigrations(cfg Cfg) {
	// reset the base fs and the global migrations
	goose.SetBaseFS(nil) // we don't want to use the base fs for plugin migrations because the embedded fs contains templates, not sql files
	goose.ResetGlobalMigrations()
	goose.SetTableName(fmt.Sprintf("goose_migration_relayer_%s_%s", cfg.Schema, cfg.ChainID.String()))
	Register0002(cfg)
}

// Migrate migrates a subsystem of the chainlink database.
// It generates migrations based on the template for the subsystem and applies them to the database.
func Migrate(ctx context.Context, db *sql.DB, cfg Cfg) error {
	tmpDir := os.TempDir()
	defer os.RemoveAll(tmpDir)

	//defer setupCoreMigrations()
	setupPluginMigrations(cfg)

	d := filepath.Join(tmpDir, cfg.Schema, fmt.Sprintf("%s", cfg.ChainID))
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", d, err)
	}
	migrations, err := generateMigrations(embeddedTmplFS, MigrationRootDir, d, cfg)
	if err != nil {
		return fmt.Errorf("failed to generate migrations for opt %v: %w", cfg, err)
	}
	fmt.Printf("Generated migrations: %v\n", migrations)

	err = goose.Up(db, d)
	if err != nil {
		return fmt.Errorf("failed to do database migration: %w", err)
	}

	return nil
}

func Rollback(ctx context.Context, db *sql.DB, version null.Int, cfg Cfg) error {
	tmpDir := os.TempDir()
	defer os.RemoveAll(tmpDir)

	setupPluginMigrations(cfg)

	// TODO: should these be saved somewhere? if so where, if not if the db itself?)
	d := filepath.Join(tmpDir, cfg.Schema, fmt.Sprintf("%s", cfg.ChainID))
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", d, err)
	}
	migrations, err := generateMigrations(embeddedTmplFS, MigrationRootDir, d, cfg)
	if err != nil {
		return fmt.Errorf("failed to generate migrations for opt %v: %w", cfg, err)
	}

	fmt.Printf("Generated migrations: %v\n", migrations)

	if version.Valid {
		return goose.DownTo(db, d, version.Int64)
	}

	return goose.Down(db, d)
}

func Current(ctx context.Context, db *sql.DB, cfg Cfg) (int64, error) {
	setupPluginMigrations(cfg)
	return goose.EnsureDBVersion(db)
}

func Status(ctx context.Context, db *sql.DB, cfg Cfg) error {
	setupPluginMigrations(cfg)
	// set the base fs only for status so that the templates are listed
	// an alternative would be to somehow keep track of the generated sql files, but that would be more complex
	// and error prone WRT to restarts
	goose.SetBaseFS(embeddedTmplFS)
	return goose.Status(db, MigrationRootDir)
}
