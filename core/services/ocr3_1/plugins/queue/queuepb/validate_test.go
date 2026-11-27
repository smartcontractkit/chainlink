package queuepb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3_1/plugins/queue/queuelimits"
)

func validOCRLimits() ocr3_1types.ReportingPluginLimits {
	return ocr3_1types.ReportingPluginLimits{
		MaxObservationBytes:                             512 * 1024,
		MaxKeyValueModifiedKeys:                         queuelimits.NumKVKeys,
		MaxKeyValueModifiedKeysPlusValuesBytes:          10 * 1024 * 1024,
		MaxBlobPayloadBytes:                             5 * 1024 * 1024,
		MaxPerOracleUnexpiredBlobCount:                  1000,
		MaxPerOracleUnexpiredBlobCumulativePayloadBytes: 100 * 1024 * 1024,
	}
}

func TestValidateConfig(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, ValidateConfig(validOCRLimits()))
	})

	minBlob := queuelimits.MaxEventFramingOverhead(queuelimits.CoordinatedEventsFieldNumber) + MinEventPayloadBytes
	tests := []struct {
		name    string
		mutate  func(*ocr3_1types.ReportingPluginLimits)
		wantErr string
	}{
		{
			name: "observation below one event",
			mutate: func(l *ocr3_1types.ReportingPluginLimits) {
				l.MaxObservationBytes = MinObservationBytes - 1
			},
			wantErr: "MaxObservationBytes",
		},
		{
			name: "fewer keys than written",
			mutate: func(l *ocr3_1types.ReportingPluginLimits) {
				l.MaxKeyValueModifiedKeys = queuelimits.NumKVKeys - 1
			},
			wantErr: "MaxKeyValueModifiedKeys",
		},
		{
			name: "kv bytes below the keys",
			mutate: func(l *ocr3_1types.ReportingPluginLimits) {
				l.MaxKeyValueModifiedKeysPlusValuesBytes = MinKeyValueModifiedKeysPlusValuesBytes - 1
			},
			wantErr: "MaxKeyValueModifiedKeysPlusValuesBytes",
		},
		{
			name:    "blob payload below one event",
			mutate:  func(l *ocr3_1types.ReportingPluginLimits) { l.MaxBlobPayloadBytes = minBlob - 1 },
			wantErr: "MaxBlobPayloadBytes",
		},
		{
			name:    "no unexpired blobs allowed",
			mutate:  func(l *ocr3_1types.ReportingPluginLimits) { l.MaxPerOracleUnexpiredBlobCount = 0 },
			wantErr: "MaxPerOracleUnexpiredBlobCount",
		},
		{
			name:    "no unexpired blob payload allowed",
			mutate:  func(l *ocr3_1types.ReportingPluginLimits) { l.MaxPerOracleUnexpiredBlobCumulativePayloadBytes = 0 },
			wantErr: "MaxPerOracleUnexpiredBlobCumulativePayloadBytes",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := validOCRLimits()
			tt.mutate(&l)
			require.ErrorContains(t, ValidateConfig(l), tt.wantErr)
		})
	}

	t.Run("boundaries are inclusive", func(t *testing.T) {
		t.Parallel()
		l := validOCRLimits()
		l.MaxObservationBytes = MinObservationBytes
		l.MaxKeyValueModifiedKeys = queuelimits.NumKVKeys
		l.MaxKeyValueModifiedKeysPlusValuesBytes = MinKeyValueModifiedKeysPlusValuesBytes
		l.MaxBlobPayloadBytes = minBlob
		require.NoError(t, ValidateConfig(l))
	})
}

func TestValidateForDon(t *testing.T) {
	t.Parallel()

	for _, don := range []struct{ n, f int }{{4, 1}, {7, 2}, {31, 10}} {
		t.Run("observation must hold the blob handle reservation", func(t *testing.T) {
			t.Parallel()
			handleBytes := ocr3_1types.BlobHandleMarshalledBytesUpperBound(don.n, don.f)
			minBytes := queuelimits.MaxBlobHandlesPerObservation*handleBytes + MinObservationBytes

			l := validOCRLimits()
			l.MaxObservationBytes = minBytes
			require.NoError(t, ValidateForDon(l, don.n, don.f), "exactly enough")

			l.MaxObservationBytes = minBytes - 1
			require.ErrorContains(t, ValidateForDon(l, don.n, don.f), "MaxObservationBytes")
		})
	}
}
