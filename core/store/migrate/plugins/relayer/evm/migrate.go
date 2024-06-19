package evm

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"gopkg.in/guregu/null.v4"
)

// hacking, trying to make a provider instead of using global state
var mu sync.Mutex
var providerCache = make(map[string]*goose.Provider)

func newProvider(db *sql.DB, cfg Cfg) (*goose.Provider, error) {
	mTable := fmt.Sprintf("goose_migration_evmrelayer_%s_%s", cfg.Schema, cfg.ChainID.String())

	mu.Lock()
	defer mu.Unlock()
	if p, ok := providerCache[mTable]; ok {
		return p, nil
	}

	store, err := database.NewStore(goose.DialectPostgres, mTable)
	if err != nil {
		return nil, fmt.Errorf("failed to create goose store for table %s: %w", mTable, err)
	}

	goMigrations := make([]*goose.Migration, 0)
	up0002, down0002, err := generate0002(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate migration 0002 for cfg %v: %w", cfg, err)
	}
	goMigrations = append(goMigrations, goose.NewGoMigration(2, up0002, down0002))

	// note we are leaking here, but can't delete the temp dir until the migrations are actually executed
	// maybe update the cache to store the temp dir and delete it when cache is deleted
	tmpDir, err := os.MkdirTemp("", cfg.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	d := filepath.Join(tmpDir, cfg.Schema, cfg.ChainID.String())
	err = os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration directory %s: %w", d, err)
	}
	migrations, err := generateMigrations(embeddedTmplFS, MigrationRootDir, d, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate migrations for opt %v: %w", cfg, err)
	}
	fmt.Printf("Generated migrations: %v\n", migrations)
	fsys := os.DirFS(d)

	// hack to avoid global state. the goose lib doesn't allow to disable the global registry
	// and also pass custom go migrations (wtf the point of disabling the global registry then?)
	// https://github.com/pressly/goose/blob/3ad137847a4c242f09e425a12c15c7c7961d4b0f/provider.go#L119
	goose.ResetGlobalMigrations()
	p, err := goose.NewProvider(
		"",
		db, fsys,
		goose.WithStore(store),
		goose.WithDisableGlobalRegistry(true), // until/if we refactor the core migrations to use goose provider
		goose.WithGoMigrations(goMigrations...))
	if err != nil {
		return nil, fmt.Errorf("failed to create goose provider: %w", err)
	}
	providerCache[mTable] = p
	return p, nil
}

/*
func setupPluginMigrations(cfg Cfg) error {
	// reset the base fs and the global migrations
	goose.SetBaseFS(nil) // we don't want to use the base fs for plugin migrations because the embedded fs contains templates, not sql files
	goose.ResetGlobalMigrations()
	goose.SetTableName(fmt.Sprintf("goose_migration_relayer_%s_%s", cfg.Schema, cfg.ChainID.String()))
	err := Register0002(cfg)
	if err != nil {
		return fmt.Errorf("failed to register migration 0002: %w", err)
	}
	return nil
}
*/
// Migrate migrates a subsystem of the chainlink database.
// It generates migrations based on the template for the subsystem and applies them to the database.
func Migrate(ctx context.Context, db *sql.DB, cfg Cfg) error {
	p, err := newProvider(db, cfg)
	if err != nil {
		return fmt.Errorf("failed to create goose provider: %w", err)
	}
	if todo, _ := p.HasPending(ctx); !todo {
		return nil
	}
	/*
		tmpDir, err := os.MkdirTemp("", cfg.Schema)
		if err != nil {
			return fmt.Errorf("failed to create temp dir: %w", err)
		}
		defer os.RemoveAll(tmpDir)

		d := filepath.Join(tmpDir, cfg.Schema, cfg.ChainID.String())
		err = os.MkdirAll(d, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
		migrations, err := generateMigrations(embeddedTmplFS, MigrationRootDir, d, cfg)
		if err != nil {
			return fmt.Errorf("failed to generate migrations for opt %v: %w", cfg, err)
		}
		fmt.Printf("Generated migrations: %v\n", migrations)
	*/
	// seems to be upside about global go migrations?
	//goose.ResetGlobalMigrations()
	r, err := p.Up(ctx)
	if err != nil {
		return fmt.Errorf("failed to do database migration: %w", err)
	}
	// todo: logger
	for _, m := range r {
		fmt.Println(m)
	}
	return nil
}

func Rollback(ctx context.Context, db *sql.DB, version null.Int, cfg Cfg) error {
	p, err := newProvider(db, cfg)
	if err != nil {
		return fmt.Errorf("failed to create goose provider: %w", err)
	}
	if version.Valid {
		_, err = p.DownTo(ctx, version.Int64)
	} else {
		_, err = p.Down(ctx)
	}

	return err
}

func Current(ctx context.Context, db *sql.DB, cfg Cfg) (int64, error) {
	p, err := newProvider(db, cfg)
	if err != nil {
		return -1, fmt.Errorf("failed to create goose provider: %w", err)
	}
	return p.GetDBVersion(ctx)

}

func Status(ctx context.Context, db *sql.DB, cfg Cfg) error {
	p, err := newProvider(db, cfg)
	if err != nil {
		return fmt.Errorf("failed to create goose provider: %w", err)
	}
	_, err = p.Status(ctx)
	return err
}
