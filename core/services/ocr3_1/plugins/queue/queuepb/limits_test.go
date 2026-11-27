package queuepb

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

func validLimits() *ReportingPluginLimits {
	return &ReportingPluginLimits{
		MaxQueryBytes:                                   1,
		MaxObservationBytes:                             512 * 1024,
		MaxReportsPlusPrecursorBytes:                    3,
		MaxReportBytes:                                  4,
		MaxReportCount:                                  5,
		MaxKeyValueModifiedKeys:                         4,
		MaxKeyValueModifiedKeysPlusValuesBytes:          10 * 1024 * 1024,
		MaxBlobPayloadBytes:                             5 * 1024 * 1024,
		MaxPerOracleUnexpiredBlobCumulativePayloadBytes: 100 * 1024 * 1024,
		MaxPerOracleUnexpiredBlobCount:                  1000,
	}
}

func mustMarshal(t *testing.T, cfg *ReportingPluginConfig) []byte {
	t.Helper()
	b, err := proto.Marshal(cfg)
	require.NoError(t, err)
	return b
}

// TestDecode asserts that Decode returns every limit from a valid config and rejects empty, malformed or invalid ones.
func TestDecode(t *testing.T) {
	t.Parallel()

	withLimits := func(mutate func(l *ReportingPluginLimits)) []byte {
		l := validLimits()
		mutate(l)
		return mustMarshal(t, &ReportingPluginConfig{Limits: l})
	}

	tests := []struct {
		name    string
		got     []byte
		wantErr string
	}{
		{name: "nil", got: nil, wantErr: "empty"},
		{name: "empty", got: []byte{}, wantErr: "empty"},
		{name: "malformed", got: []byte{0xff}, wantErr: "unmarshal"},
		// An empty message marshals to zero bytes, so an unknown field (number 15) keeps the input non-empty.
		{name: "missing limits", got: []byte{0x78, 0x01}, wantErr: "no limits"},
		{
			name:    "zero modified keys",
			got:     withLimits(func(l *ReportingPluginLimits) { l.MaxKeyValueModifiedKeys = 0 }),
			wantErr: "MaxKeyValueModifiedKeys",
		},
		{
			name:    "zero observation bytes",
			got:     withLimits(func(l *ReportingPluginLimits) { l.MaxObservationBytes = 0 }),
			wantErr: "MaxObservationBytes",
		},
		{
			name:    "zero modified keys plus values bytes",
			got:     withLimits(func(l *ReportingPluginLimits) { l.MaxKeyValueModifiedKeysPlusValuesBytes = 0 }),
			wantErr: "MaxKeyValueModifiedKeysPlusValuesBytes",
		},
		{
			name:    "observation bytes above libocr maximum",
			got:     withLimits(func(l *ReportingPluginLimits) { l.MaxObservationBytes = ocr3_1types.MaxMaxObservationBytes + 1 }),
			wantErr: "MaxObservationBytes",
		},
		{
			name:    "too few modified keys for the four writes",
			got:     withLimits(func(l *ReportingPluginLimits) { l.MaxKeyValueModifiedKeys = 3 }),
			wantErr: "MaxKeyValueModifiedKeys",
		},
		{
			name:    "too few key bytes for the four writes",
			got:     withLimits(func(l *ReportingPluginLimits) { l.MaxKeyValueModifiedKeysPlusValuesBytes = 1 }),
			wantErr: "MaxKeyValueModifiedKeysPlusValuesBytes",
		},
		{
			name:    "blob payload too small for one event",
			got:     withLimits(func(l *ReportingPluginLimits) { l.MaxBlobPayloadBytes = 1 }),
			wantErr: "MaxBlobPayloadBytes",
		},
		{
			name:    "no unexpired blobs allowed",
			got:     withLimits(func(l *ReportingPluginLimits) { l.MaxPerOracleUnexpiredBlobCount = 0 }),
			wantErr: "MaxPerOracleUnexpiredBlobCount",
		},
		{
			name:    "no unexpired blob payload allowed",
			got:     withLimits(func(l *ReportingPluginLimits) { l.MaxPerOracleUnexpiredBlobCumulativePayloadBytes = 0 }),
			wantErr: "MaxPerOracleUnexpiredBlobCumulativePayloadBytes",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg, err := Decode(tt.got)
			require.ErrorContains(t, err, tt.wantErr)
			require.Nil(t, cfg)
		})
	}

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		cfg, err := Decode(mustMarshal(t, &ReportingPluginConfig{Limits: validLimits()}))
		require.NoError(t, err)

		got, err := cfg.OCRLimits()
		require.NoError(t, err)
		want := ocr3_1types.ReportingPluginLimits{
			MaxQueryBytes:                                   1,
			MaxObservationBytes:                             512 * 1024,
			MaxReportsPlusPrecursorBytes:                    3,
			MaxReportBytes:                                  4,
			MaxReportCount:                                  5,
			MaxKeyValueModifiedKeys:                         4,
			MaxKeyValueModifiedKeysPlusValuesBytes:          10 * 1024 * 1024,
			MaxBlobPayloadBytes:                             5 * 1024 * 1024,
			MaxPerOracleUnexpiredBlobCumulativePayloadBytes: 100 * 1024 * 1024,
			MaxPerOracleUnexpiredBlobCount:                  1000,
		}
		require.Equal(t, want, got)
	})
}
