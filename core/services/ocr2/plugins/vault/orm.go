package vault

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

type ORM interface {
	dkgocrtypes.ResultPackageDatabase
	GetResultPackageCount(ctx context.Context) (int64, error)
}

type orm struct {
	ds sqlutil.DataSource
}

var _ ORM = (*orm)(nil)

func NewVaultORM(ds sqlutil.DataSource) ORM {
	return &orm{ds: ds}
}

// Sanity check the result package. Otherwise, trust the writer.
func (o *orm) validateResultPackage(value dkgocrtypes.ResultPackageDatabaseValue) error {
	var zeroDigest types.ConfigDigest
	if value.ConfigDigest == zeroDigest {
		return errors.New("config digest cannot be zero")
	}

	if value.SeqNr == 0 {
		return errors.New("sequence number cannot be zero")
	}

	if len(value.ReportWithResultPackage) == 0 {
		return errors.New("report with result package cannot be empty")
	}

	if len(value.Signatures) == 0 {
		return errors.New("signatures cannot be empty")
	}

	return nil
}

func (o *orm) GetResultPackageCount(ctx context.Context) (int64, error) {
	var count int64

	query := `SELECT COUNT(*) FROM dkg_results;`
	row := o.ds.QueryRowxContext(ctx, query)
	err := row.Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get dkg result count")
	}

	return count, nil
}

func (o *orm) ReadResultPackage(ctx context.Context, iid dkgocrtypes.InstanceID) (*dkgocrtypes.ResultPackageDatabaseValue, error) {
	var configDigest []byte
	var seqNr uint64
	var reportWithResultPackage []byte
	var signatures pq.ByteaArray
	var signerOracleIDs []byte

	query := `SELECT config_digest, seq_nr, report_with_result_package, signatures, signer_oracle_ids FROM dkg_results WHERE instance_id = $1;`
	row := o.ds.QueryRowxContext(ctx, query, iid)
	err := row.Scan(&configDigest, &seqNr, &reportWithResultPackage, &signatures, &signerOracleIDs)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to read dkg result")
	}

	var cd types.ConfigDigest
	copy(cd[:], configDigest)

	attributedSigs := make([]types.AttributedOnchainSignature, len(signatures))
	for i := range signatures {
		attributedSigs[i] = types.AttributedOnchainSignature{
			Signature: signatures[i],
			Signer:    commontypes.OracleID(signerOracleIDs[i]),
		}
	}

	value := &dkgocrtypes.ResultPackageDatabaseValue{
		ConfigDigest:            cd,
		SeqNr:                   seqNr,
		ReportWithResultPackage: reportWithResultPackage,
		Signatures:              attributedSigs,
	}

	return value, nil
}

func (o *orm) WriteResultPackage(ctx context.Context,
	instanceID dkgocrtypes.InstanceID,
	value dkgocrtypes.ResultPackageDatabaseValue,
) error {
	if err := o.validateResultPackage(value); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	signatures := make(pq.ByteaArray, len(value.Signatures))
	signerOracleIDs := make([]byte, len(value.Signatures))
	for i, sig := range value.Signatures {
		signatures[i] = sig.Signature
		signerOracleIDs[i] = byte(sig.Signer)
	}

	query := `
        INSERT INTO dkg_results (instance_id, config_digest, seq_nr, report_with_result_package, signatures, signer_oracle_ids, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
        ON CONFLICT (instance_id) DO UPDATE SET
        config_digest = EXCLUDED.config_digest,
        seq_nr = EXCLUDED.seq_nr,
        report_with_result_package = EXCLUDED.report_with_result_package,
        signatures = EXCLUDED.signatures,
        signer_oracle_ids = EXCLUDED.signer_oracle_ids,
        updated_at = NOW();
    `
	_, err := o.ds.ExecContext(ctx, query, instanceID, value.ConfigDigest[:], value.SeqNr, value.ReportWithResultPackage, signatures, signerOracleIDs)
	return errors.Wrap(err, "failed to write dkg result")
}
