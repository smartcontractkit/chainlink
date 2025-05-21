package retirement

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type mockORM struct {
	storedAttestedRetirementReports map[ocr2types.ConfigDigest][]byte
	storedConfigs                   map[ocr2types.ConfigDigest]Config

	err error
}

func (m *mockORM) StoreAttestedRetirementReport(ctx context.Context, cd ocr2types.ConfigDigest, attestedRetirementReport []byte) error {
	m.storedAttestedRetirementReports[cd] = attestedRetirementReport
	return m.err
}
func (m *mockORM) LoadAttestedRetirementReports(ctx context.Context) (map[ocr2types.ConfigDigest][]byte, error) {
	return m.storedAttestedRetirementReports, m.err
}
func (m *mockORM) StoreConfig(ctx context.Context, cd ocr2types.ConfigDigest, signers [][]byte, f uint8) error {
	m.storedConfigs[cd] = Config{Signers: signers, F: f, Digest: cd}
	return m.err
}
func (m *mockORM) LoadConfigs(ctx context.Context) ([]Config, error) {
	configs := make([]Config, 0, len(m.storedConfigs))
	for _, config := range m.storedConfigs {
		configs = append(configs, config)
	}
	return configs, m.err
}

func Test_RetirementReportCache(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	lggr := logger.TestLogger(t)
	orm := &mockORM{
		make(map[ocrtypes.ConfigDigest][]byte),
		make(map[ocrtypes.ConfigDigest]Config),
		nil,
	}
	exampleRetirementReport := []byte{1, 2, 3}
	exampleRetirementReport2 := []byte{4, 5, 6}
	exampleSignatures := []ocrtypes.AttributedOnchainSignature{
		{Signature: []byte("signature0"), Signer: 0},
		{Signature: []byte("signature1"), Signer: 1},
		{Signature: []byte("signature2"), Signer: 2},
		{Signature: []byte("signature3"), Signer: 3},
	}
	// this is a serialized protobuf of report with 4 signers
	exampleAttestedRetirementReport := []byte{0xa, 0x3, 0x1, 0x2, 0x3, 0x10, 0x64, 0x1a, 0xc, 0xa, 0xa, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x30, 0x1a, 0xe, 0xa, 0xa, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x31, 0x10, 0x1, 0x1a, 0xe, 0xa, 0xa, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x32, 0x10, 0x2, 0x1a, 0xe, 0xa, 0xa, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x33, 0x10, 0x3}
	exampleDigest := ocrtypes.ConfigDigest{1}
	exampleDigest2 := ocrtypes.ConfigDigest{2}

	seqNr := uint64(100)

	t.Run("start loads from ORM", func(t *testing.T) {
		rrc := newRetirementReportCache(lggr, orm)

		t.Run("orm failure, errors", func(t *testing.T) {
			orm.err = errors.New("orm failed")
			err := rrc.start(ctx)
			assert.EqualError(t, err, "failed to load attested retirement reports: orm failed")
		})
		t.Run("orm success, loads both configs and attestedRetirementReports from orm", func(t *testing.T) {
			orm.err = nil
			orm.storedAttestedRetirementReports = map[ocr2types.ConfigDigest][]byte{
				exampleDigest:  exampleAttestedRetirementReport,
				exampleDigest2: exampleAttestedRetirementReport,
			}
			config1 := Config{Digest: exampleDigest, Signers: [][]byte{{1}, {2}, {3}, {4}}, F: 1}
			config2 := Config{Digest: exampleDigest2, Signers: [][]byte{{5}, {6}, {7}, {8}}, F: 2}
			orm.storedConfigs[exampleDigest] = config1
			orm.storedConfigs[exampleDigest2] = config2

			err := rrc.start(ctx)
			assert.NoError(t, err)

			assert.Len(t, rrc.arrs, 2)
			assert.Equal(t, exampleAttestedRetirementReport, rrc.arrs[exampleDigest])
			assert.Equal(t, exampleAttestedRetirementReport, rrc.arrs[exampleDigest2])

			assert.Len(t, rrc.configs, 2)
			assert.Equal(t, config1, rrc.configs[exampleDigest])
			assert.Equal(t, config2, rrc.configs[exampleDigest2])
		})
	})

	t.Run("StoreAttestedRetirementReport", func(t *testing.T) {
		rrc := newRetirementReportCache(lggr, orm)

		err := rrc.StoreAttestedRetirementReport(ctx, exampleDigest, seqNr, exampleRetirementReport, exampleSignatures)
		assert.NoError(t, err)

		assert.Len(t, rrc.arrs, 1)
		assert.Equal(t, exampleAttestedRetirementReport, rrc.arrs[exampleDigest])
		assert.Equal(t, exampleAttestedRetirementReport, orm.storedAttestedRetirementReports[exampleDigest])

		t.Run("does nothing if retirement report already exists for the given config digest", func(t *testing.T) {
			err = rrc.StoreAttestedRetirementReport(ctx, exampleDigest, seqNr, exampleRetirementReport2, exampleSignatures)
			assert.NoError(t, err)
			assert.Len(t, rrc.arrs, 1)
			assert.Equal(t, exampleAttestedRetirementReport, rrc.arrs[exampleDigest])
		})

		t.Run("returns error if ORM store fails", func(t *testing.T) {
			orm.err = errors.New("failed to store")
			err = rrc.StoreAttestedRetirementReport(ctx, exampleDigest2, seqNr, exampleRetirementReport, exampleSignatures)
			assert.Error(t, err)

			// it wasn't cached
			assert.Len(t, rrc.arrs, 1)
		})

		t.Run("second retirement report succeeds when orm starts working again", func(t *testing.T) {
			orm.err = nil
			err := rrc.StoreAttestedRetirementReport(ctx, exampleDigest2, seqNr, exampleRetirementReport, exampleSignatures)
			assert.NoError(t, err)

			assert.Len(t, rrc.arrs, 2)
			assert.Equal(t, exampleAttestedRetirementReport, rrc.arrs[exampleDigest2])
			assert.Equal(t, exampleAttestedRetirementReport, orm.storedAttestedRetirementReports[exampleDigest2])

			assert.Len(t, orm.storedAttestedRetirementReports, 2)
		})
	})
	t.Run("AttestedRetirementReport", func(t *testing.T) {
		rrc := newRetirementReportCache(lggr, orm)

		attestedRetirementReport, exists := rrc.AttestedRetirementReport(exampleDigest)
		assert.False(t, exists)
		assert.Nil(t, attestedRetirementReport)

		rrc.arrs[exampleDigest] = exampleAttestedRetirementReport

		attestedRetirementReport, exists = rrc.AttestedRetirementReport(exampleDigest)
		assert.True(t, exists)
		assert.Equal(t, exampleAttestedRetirementReport, attestedRetirementReport)
	})
	t.Run("StoreConfig", func(t *testing.T) {
		rrc := newRetirementReportCache(lggr, orm)

		signers := [][]byte{{1}, {2}, {3}, {4}}

		err := rrc.StoreConfig(ctx, exampleDigest, signers, 1)
		assert.NoError(t, err)

		assert.Len(t, rrc.configs, 1)
		assert.Equal(t, Config{Digest: exampleDigest, Signers: [][]byte{{1}, {2}, {3}, {4}}, F: 1}, rrc.configs[exampleDigest])
		assert.Equal(t, Config{Digest: exampleDigest, Signers: [][]byte{{1}, {2}, {3}, {4}}, F: 1}, orm.storedConfigs[exampleDigest])

		t.Run("Config", func(t *testing.T) {
			config, exists := rrc.Config(exampleDigest)
			assert.True(t, exists)
			assert.Equal(t, Config{Digest: exampleDigest, Signers: [][]byte{{1}, {2}, {3}, {4}}, F: 1}, config)
		})
	})
}
