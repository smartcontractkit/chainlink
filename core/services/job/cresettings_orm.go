package job

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
)

// CRESettingsConfigTypeCapRegistry is the cresettings config_type carrying the offchain
// capabilities registry payload in CRESettingsSpec.OffchainConfig.
const CRESettingsConfigTypeCapRegistry = "capabilities_registry"

// capRegMetrics is the metrics sink for stale capabilities_registry submissions (a variable so
// tests can observe it).
var capRegMetrics = globalconfig.DefaultMetrics

// creSettingsSingleCapRegistryIndex enforces at most one config_type=capabilities_registry job.
const creSettingsSingleCapRegistryIndex = "idx_cre_settings_specs_single_capabilities_registry"

var (
	// ErrCRESettingsCapRegistryExists is returned when creating a second capabilities_registry
	// cresettings job. Replace the existing job (delete, then create) instead.
	ErrCRESettingsCapRegistryExists = errors.New("a cresettings job with config_type \"capabilities_registry\" already exists")
	// ErrCRESettingsCapRegistryStale is returned when creating a capabilities_registry job whose
	// payload version is not newer than the newest version ever committed on this node.
	ErrCRESettingsCapRegistryStale = errors.New("stale capabilities_registry payload")
)

func (o *orm) insertCRESettingsSpec(ctx context.Context, spec *CRESettingsSpec) (specID int32, err error) {
	if spec.ConfigType == CRESettingsConfigTypeCapRegistry {
		if err = o.advanceCapRegistryHighWater(ctx, spec); err != nil {
			if errors.Is(err, ErrCRESettingsCapRegistryStale) {
				domain, env := globalconfig.PayloadLabels(spec.OffchainConfig)
				capRegMetrics().RecordValidationError(ctx, domain, env)
			}
			return 0, err
		}
	}
	specID, err = o.prepareQuerySpecID(ctx, `INSERT INTO cre_settings_specs (settings, hash, config_type, offchain_config, created_at, updated_at) VALUES (:settings, :hash, :config_type, :offchain_config, NOW(), NOW()) RETURNING id;`, spec)
	var pqErr *pgconn.PgError
	if errors.As(err, &pqErr) && pqErr.Code == "23505" && pqErr.ConstraintName == creSettingsSingleCapRegistryIndex {
		return 0, ErrCRESettingsCapRegistryExists
	}
	return specID, err
}

// advanceCapRegistryHighWater enforces version monotonicity durably. It runs in the same
// transaction as the spec insert (CreateJob's), locking the single high-water row, so:
//   - a rolled-back create leaves the high-water mark untouched,
//   - concurrent creates are serialized,
//   - the mark survives job deletion and node restarts.
//
// A payload is accepted when its version is greater than the mark, or when it is exactly the
// payload that set the mark (re-creating the same job after deleting it).
func (o *orm) advanceCapRegistryHighWater(ctx context.Context, spec *CRESettingsSpec) error {
	reg, err := globalconfig.Parse(spec.OffchainConfig)
	if err != nil {
		return fmt.Errorf("invalid capabilities_registry config: %w", err)
	}
	sum := sha256.Sum256([]byte(spec.OffchainConfig))
	hash := hex.EncodeToString(sum[:])
	if spec.Hash != "" && spec.Hash != hash {
		return fmt.Errorf("invalid capabilities_registry hash %s: calculated %s", spec.Hash, hash)
	}
	version := reg.GetVersion()

	var row struct {
		Version string `db:"version"`
		Hash    string `db:"hash"`
	}
	if err = o.ds.GetContext(ctx, &row, `SELECT version::text AS version, hash FROM cre_offchain_registry_high_water WHERE id = 1 FOR UPDATE`); err != nil {
		return fmt.Errorf("failed to read capabilities_registry high-water mark: %w", err)
	}
	applied, err := strconv.ParseUint(row.Version, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid capabilities_registry high-water mark %q: %w", row.Version, err)
	}
	switch {
	case version > applied:
	case version == applied && hash == row.Hash:
		return nil // re-creating the payload that set the mark
	default:
		return fmt.Errorf("%w: version %d is not newer than committed version %d", ErrCRESettingsCapRegistryStale, version, applied)
	}
	if _, err = o.ds.ExecContext(ctx, `UPDATE cre_offchain_registry_high_water SET version = $1::numeric, hash = $2 WHERE id = 1`,
		strconv.FormatUint(version, 10), hash); err != nil {
		return fmt.Errorf("failed to advance capabilities_registry high-water mark: %w", err)
	}
	return nil
}

// LoadCapabilitiesRegistrySpec returns the committed capabilities_registry spec visible to ds,
// or nil when there is none. Callers that need committed state must pass a DataSource that is
// not inside the transaction creating or deleting the job.
func LoadCapabilitiesRegistrySpec(ctx context.Context, ds sqlutil.DataSource) (*CRESettingsSpec, error) {
	var spec CRESettingsSpec
	err := ds.GetContext(ctx, &spec, `SELECT s.* FROM cre_settings_specs s
		JOIN jobs j ON j.cre_settings_spec_id = s.id
		WHERE s.config_type = $1`, CRESettingsConfigTypeCapRegistry)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load capabilities_registry spec: %w", err)
	}
	return &spec, nil
}
