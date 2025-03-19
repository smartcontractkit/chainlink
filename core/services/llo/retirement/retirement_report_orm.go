package retirement

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

type RetirementReportCacheORM interface {
	StoreAttestedRetirementReport(ctx context.Context, cd ocr2types.ConfigDigest, attestedRetirementReport []byte) error
	LoadAttestedRetirementReports(ctx context.Context) (map[ocr2types.ConfigDigest][]byte, error)
	StoreConfig(ctx context.Context, cd ocr2types.ConfigDigest, signers [][]byte, f uint8) error
	LoadConfigs(ctx context.Context) ([]Config, error)
}

type retirementReportCacheORM struct {
	ds sqlutil.DataSource
}

func (o *retirementReportCacheORM) StoreAttestedRetirementReport(ctx context.Context, cd ocr2types.ConfigDigest, attestedRetirementReport []byte) error {
	_, err := o.ds.ExecContext(ctx, `
INSERT INTO llo_retirement_report_cache (config_digest, attested_retirement_report, updated_at)
VALUES ($1, $2, NOW())
ON CONFLICT (config_digest) DO NOTHING
`, cd, attestedRetirementReport)
	if err != nil {
		return fmt.Errorf("StoreAttestedRetirementReport failed: %w", err)
	}
	return nil
}

func (o *retirementReportCacheORM) LoadAttestedRetirementReports(ctx context.Context) (map[ocr2types.ConfigDigest][]byte, error) {
	rows, err := o.ds.QueryContext(ctx, "SELECT config_digest, attested_retirement_report FROM llo_retirement_report_cache")
	if err != nil {
		return nil, fmt.Errorf("LoadAttestedRetirementReports failed: %w", err)
	}
	defer rows.Close()

	reports := make(map[ocr2types.ConfigDigest][]byte)
	for rows.Next() {
		var rawCd []byte
		var arr []byte
		if err := rows.Scan(&rawCd, &arr); err != nil {
			return nil, fmt.Errorf("LoadAttestedRetirementReports failed: %w", err)
		}
		cd, err := ocr2types.BytesToConfigDigest(rawCd)
		if err != nil {
			return nil, fmt.Errorf("LoadAttestedRetirementReports failed to scan config digest: %w", err)
		}
		reports[cd] = arr
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("LoadAttestedRetirementReports failed: %w", err)
	}

	return reports, nil
}

func (o *retirementReportCacheORM) StoreConfig(ctx context.Context, cd ocr2types.ConfigDigest, signers [][]byte, f uint8) error {
	// do nothing on overwrite since configs are supposedly immutable
	_, err := o.ds.ExecContext(ctx, `INSERT INTO llo_retirement_report_cache_configs (config_digest, signers, f, updated_at) VALUES ($1, $2, $3, NOW()) ON CONFLICT (config_digest) DO NOTHING`, cd, signers, f)
	if err != nil {
		return fmt.Errorf("StoreConfig failed: %w", err)
	}
	return nil
}

type Config struct {
	Digest  [32]byte      `db:"config_digest"`
	Signers pq.ByteaArray `db:"signers"`
	F       uint8         `db:"f"`
}

type scannableConfigDigest [32]byte

func (s *scannableConfigDigest) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	cd, err := ocr2types.BytesToConfigDigest(b)
	if err != nil {
		return err
	}
	copy(s[:], cd[:])
	return nil
}

func (o *retirementReportCacheORM) LoadConfigs(ctx context.Context) (configs []Config, err error) {
	type config struct {
		Digest  scannableConfigDigest `db:"config_digest"`
		Signers pq.ByteaArray         `db:"signers"`
		F       uint8                 `db:"f"`
	}
	var rawCfgs []config
	err = o.ds.SelectContext(ctx, &rawCfgs, `SELECT config_digest, signers, f FROM llo_retirement_report_cache_configs ORDER BY config_digest`)
	if err != nil {
		return nil, fmt.Errorf("LoadConfigs failed: %w", err)
	}
	for _, rawCfg := range rawCfgs {
		var cfg Config
		copy(cfg.Digest[:], rawCfg.Digest[:])
		cfg.Signers = rawCfg.Signers
		cfg.F = rawCfg.F
		configs = append(configs, cfg)
	}
	return
}
