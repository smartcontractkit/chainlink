package migrate

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate/migrations" // Invoke init() functions within migrations pkg.
)

//go:embed migrations/*.sql migrations/*.go
var embedMigrations embed.FS

// go:embed plugins/relayers/**/*.tmpl.sql
var embedRelayerMigrations embed.FS

const PLUGIN_MIGRATIONS_DIR string = "plugins"

const MIGRATIONS_DIR string = "migrations"

// go:embed manifest.txt
var migrationManifest string

func init() {
	setupCoreMigrations()
	logMigrations := os.Getenv("CL_LOG_SQL_MIGRATIONS")
	verbose, _ := strconv.ParseBool(logMigrations)
	goose.SetVerbose(verbose)
}

func setupCoreMigrations() {
	goose.SetBaseFS(embedMigrations)
	goose.SetSequential(true)
	goose.SetTableName("goose_migrations")
}

// Ensure we migrated from v1 migrations to goose_migrations
func ensureMigrated(ctx context.Context, db *sql.DB) error {
	sqlxDB := pg.WrapDbWithSqlx(db)
	var names []string
	err := sqlxDB.SelectContext(ctx, &names, `SELECT id FROM migrations`)
	if err != nil {
		// already migrated
		return nil
	}
	// ensure that no legacy job specs are present: we _must_ bail out early if
	// so because otherwise we run the risk of dropping working jobs if the
	// user has not read the release notes
	err = migrations.CheckNoLegacyJobs(ctx, db)
	if err != nil {
		return err
	}

	// Look for the squashed migration. If not present, the db needs to be migrated on an earlier release first
	found := false
	for _, name := range names {
		if name == "1611847145" {
			found = true
		}
	}
	if !found {
		return pkgerrors.New("database state is too old. Need to migrate to chainlink version 0.9.10 first before upgrading to this version. This upgrade is NOT REVERSIBLE, so it is STRONGLY RECOMMENDED that you take a database backup before continuing")
	}

	// ensure a goose migrations table exists with it's initial v0
	if _, err = goose.GetDBVersionContext(ctx, db); err != nil {
		return err
	}

	// insert records for existing migrations
	//nolint
	sql := fmt.Sprintf(`INSERT INTO %s (version_id, is_applied) VALUES ($1, true);`, goose.TableName())
	return sqlutil.TransactDataSource(ctx, sqlxDB, nil, func(tx sqlutil.DataSource) error {
		for _, name := range names {
			var id int64
			// the first migration doesn't follow the naming convention
			if name == "1611847145" {
				id = 1
			} else {
				idx := strings.Index(name, "_")
				if idx < 0 {
					// old migration we don't care about
					continue
				}

				id, err = strconv.ParseInt(name[:idx], 10, 64)
				if err == nil && id <= 0 {
					return pkgerrors.New("migration IDs must be greater than zero")
				}
			}

			if _, err = tx.ExecContext(ctx, sql, id); err != nil {
				return err
			}
		}

		_, err = tx.ExecContext(ctx, "DROP TABLE migrations;")
		return err
	})
}

type MigrationConfig struct {
	Type     string //relayer,app
	Template string //chain family
	Schema   string // evm, optimism, arbitrum
	Dir      string
}

func (o MigrationConfig) validate() error {
	if o.Type != "relayer" {
		return fmt.Errorf("unknown migration type: %s", o.Type)
	}
	if o.Template != "evm" {
		return fmt.Errorf("unknown migration template: %s", o.Template)
	}
	return nil
}

func Migrate(ctx context.Context, db *sql.DB) error {
	if err := ensureMigrated(ctx, db); err != nil {
		return err
	}
	// WithAllowMissing is necessary when upgrading from 0.10.14 since it
	// includes out-of-order migrations
	err := goose.Up(db, MIGRATIONS_DIR, goose.WithAllowMissing())
	if err != nil {
		return fmt.Errorf("failed to do core database migration: %w", err)
	}
	return nil
}

// MigratePlugin migrates a subsystem of the chainlink database.
// It generates migrations based on the template for the subsystem and applies them to the database.
func MigratePlugin(ctx context.Context, db *sql.DB, cfg MigrationConfig) error {

	if err := cfg.validate(); err != nil {
		return fmt.Errorf("invalid migration option: %w", err)
	}

	tmpDir := os.TempDir()
	defer os.RemoveAll(tmpDir)

	defer setupCoreMigrations()
	setupPluginMigrations(cfg)

	d := filepath.Join(tmpDir, cfg.Template, cfg.Schema)
	migrations, err := generateMigrations(cfg.Dir, d, SQLConfig{Schema: cfg.Schema})
	if err != nil {
		return fmt.Errorf("failed to generate migrations for opt %v: %w", cfg, err)
	}
	fmt.Printf("Generated migrations: %v\n", migrations)

	err = goose.Up(db, d)
	if err != nil {
		return fmt.Errorf("failed to do %s database migration: %w", cfg.Type, err)
	}

	return nil
}

func Rollback(ctx context.Context, db *sql.DB, version null.Int) error {
	if err := ensureMigrated(ctx, db); err != nil {
		return err
	}
	if version.Valid {
		return goose.DownTo(db, MIGRATIONS_DIR, version.Int64)
	}
	return goose.Down(db, MIGRATIONS_DIR)
}

func RollbackPlugin(ctx context.Context, db *sql.DB, version null.Int, cfg MigrationConfig) error {
	if err := cfg.validate(); err != nil {
		return fmt.Errorf("invalid migration option: %w", err)
	}

	tmpDir := os.TempDir()
	defer os.RemoveAll(tmpDir)

	defer setupCoreMigrations()
	setupPluginMigrations(cfg)

	// TODO: should these be saved somewhere? if so where, if not if the db itself?)
	d := filepath.Join(tmpDir, cfg.Template, cfg.Schema)
	migrations, err := generateMigrations(cfg.Dir, d, SQLConfig{Schema: cfg.Schema})
	if err != nil {
		return fmt.Errorf("failed to generate migrations for opt %v: %w", cfg, err)
	}
	fmt.Printf("Generated migrations: %v\n", migrations)

	if version.Valid {
		return goose.DownTo(db, d, version.Int64)
	}
	return goose.Down(db, d)
}

func Current(ctx context.Context, db *sql.DB) (int64, error) {
	if err := ensureMigrated(ctx, db); err != nil {
		return -1, err
	}
	return goose.EnsureDBVersion(db)
}

func CurrentPlugin(ctx context.Context, db *sql.DB, cfg MigrationConfig) (int64, error) {
	if err := ensureMigrated(ctx, db); err != nil {
		return -1, err
	}
	defer setupCoreMigrations()
	setupPluginMigrations(cfg)

	return goose.EnsureDBVersion(db)
}

func setupPluginMigrations(cfg MigrationConfig) {
	goose.SetBaseFS(nil)
	goose.ResetGlobalMigrations()
	goose.SetTableName(fmt.Sprintf("goose_migration_%s_%s", cfg.Template, cfg.Schema))
}

func Status(ctx context.Context, db *sql.DB) error {
	if err := ensureMigrated(ctx, db); err != nil {
		return err
	}
	return goose.Status(db, MIGRATIONS_DIR)
}

func Create(db *sql.DB, name, migrationType string) error {
	return goose.Create(db, "core/store/migrate/migrations", name, migrationType)
}

// SetMigrationENVVars is used to inject values from config to goose migrations via env.
func SetMigrationENVVars(generalConfig chainlink.GeneralConfig) error {
	if generalConfig.EVMEnabled() {
		err := os.Setenv(env.EVMChainIDNotNullMigration0195, generalConfig.EVMConfigs()[0].ChainID.String())
		if err != nil {
			panic(pkgerrors.Wrap(err, "failed to set migrations env variables"))
		}
	}
	return nil
}

type Manifest struct {
	Entries []ManifestEntry
	m       map[string]ManifestEntry
}

func (m Manifest) Latest() (ManifestEntry, error) {
	if len(m.Entries) == 0 {
		return ManifestEntry{}, errors.New("no entries in manifest")
	}
	return m.Entries[0], nil
}

func (m Manifest) After(e ManifestEntry) ([]ManifestEntry, error) {
	indexed, exists := m.m[e.id()]
	if !exists {
		return nil, fmt.Errorf("entry not found in manifest: %v key %s", e, e.id())
	}
	var entries []ManifestEntry
	// reverse order index
	for i := len(m.Entries) - 1; i > indexed.index; i-- {
		entries = append(entries, m.Entries[i])
	}
	return entries, nil

}

func (m Manifest) Before(e ManifestEntry) ([]ManifestEntry, error) {
	var entries []ManifestEntry
	for _, entry := range m.Entries {
		if entry.Version < e.Version {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

type ManifestEntry struct {
	Type          string // core, plugin
	PluginKind    string // relayer, app
	PluginVariant string // evm, optimism, arbitrum, functions, ccip
	Version       int

	index int    // 0 ==> most recent
	path  string // migration path

}

func (m ManifestEntry) root() string {
	if m.Type == "core" {
		return MIGRATIONS_DIR
	}
	return filepath.Join(PLUGIN_MIGRATIONS_DIR, m.PluginKind, m.PluginVariant)
}

func (m ManifestEntry) id() string {
	return fmt.Sprintf("%s_%d", m.root(), m.Version)
}

func validateMigrationEntry(m ManifestEntry) error {
	if m.Version == 0 {
		return fmt.Errorf("missing version")
	}
	if m.Type != "core" && m.Type != "plugin" {
		return fmt.Errorf("unknown migration type: %s", m.Type)
	}
	if m.Type == "core" {
		if m.PluginKind != "" || m.PluginVariant != "" {
			return fmt.Errorf("core migration: expected empty plugin configruation but got plugin kind '%s', variant '%s'", m.PluginKind, m.PluginVariant)
		}
	}
	if m.Type == "plugin" {
		if m.PluginKind != "relayer" && m.PluginKind != "app" {
			return fmt.Errorf("unknown plugin kind: %s", m.PluginKind)
		}
		if m.PluginVariant == "" {
			return fmt.Errorf("missing plugin variant")
		}
	}
	return nil
}

func LoadManifest(txt string) (Manifest, error) {
	lines := strings.Split(txt, "\n")
	var m Manifest
	m.m = make(map[string]ManifestEntry, len(lines))
	for i, l := range lines {
		if l == "" {
			continue
		}
		e, err := parseEntry(l)
		if err != nil {
			return Manifest{}, fmt.Errorf("failed to parse line %s: %w", l, err)
		}
		e.index = i
		m.Entries = append(m.Entries, e)
		m.m[e.id()] = e
	}
	return m, nil
}

var (
	coreMigrationsRoot  = MIGRATIONS_DIR
	relayMigrationsRoot = filepath.Join(PLUGIN_MIGRATIONS_DIR, "relayers")
	appMigrationsRoot   = filepath.Join(PLUGIN_MIGRATIONS_DIR, "apps")
	regexGenerator      = func(root string) string {
		return fmt.Sprintf(`^%s/[0-9]{4}_.*$`, root)
	}
	coreRe                        = regexp.MustCompile(regexGenerator(coreMigrationsRoot))
	relayRe                       = regexp.MustCompile(regexGenerator(relayMigrationsRoot))
	appRe                         = regexp.MustCompile(regexGenerator(appMigrationsRoot))
	errInvalidManifestEntryName   = fmt.Errorf("invalid migration name")
	errInvalidPluginManifestEntry = fmt.Errorf("invalid plugin migration path")
)

func parseEntry(path string) (e ManifestEntry, err error) {
	p := strings.TrimPrefix(path, "core/store/migrate/")
	e, err = parseCoreEntry(p)
	if err != nil {
		var err2 error
		e, err2 = parsePluginEntry(p)
		if err2 != nil {
			return e, errors.Join(fmt.Errorf("failed to parse path '%s' into entry", path), err, err2)
		}
	}
	return e, validateMigrationEntry(e)
}

func parseCoreEntry(path string) (ManifestEntry, error) {
	version, err := extractVersion(filepath.Base(path))
	if err != nil {
		return ManifestEntry{}, fmt.Errorf("failed to extract version for %s: %w", path, err)
	}
	parts := strings.Split(path, "/")
	path = strings.TrimPrefix(path, "core/store/migrate/")
	if len(parts) != 2 {
		return ManifestEntry{}, fmt.Errorf("invalid core migration path: %s", path)
	}
	return ManifestEntry{
		path:    path,
		Type:    "core",
		Version: version,
	}, nil
}

func parsePluginEntry(path string) (ManifestEntry, error) {
	version, err := extractVersion(filepath.Base(path))
	if err != nil {
		return ManifestEntry{}, fmt.Errorf("failed to extract version for %s: %w", path, err)
	}
	path = strings.TrimPrefix(path, "core/store/migrate/")

	// plugins/<kind>/<variant>/<version>_<name>.sql
	parts := strings.Split(path, "/")
	if len(parts) != 4 {
		return ManifestEntry{}, fmt.Errorf("invalid plugin migration path: %s", path)
	}
	return ManifestEntry{
		path:          path,
		Type:          "plugin",
		PluginKind:    parts[1],
		PluginVariant: parts[2],
		Version:       version,
	}, nil
}

func extractVersion(migrationName string) (int, error) {
	if migrationName == "" {
		return 0, fmt.Errorf("%w: empty migration name", errInvalidManifestEntryName)
	}
	parts := strings.Split(migrationName, "_")
	if len(parts) < 2 {
		return 0, fmt.Errorf("%w: %s", errInvalidManifestEntryName, migrationName)
	}
	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("%w: could not parse version: %s", errInvalidManifestEntryName, migrationName)
	}
	return version, nil
}
