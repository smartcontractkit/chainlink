package vault

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type mockGauge struct {
	gotValue int64
}

func (m *mockGauge) Record(ctx context.Context, value int64, options ...metric.RecordOption) {
	m.gotValue = value
}

func TestDiskMonitor_emitDirSizeMetric(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)

	dm := &DiskMonitor{
		sizeOfDir: func() (int64, error) {
			return 42, nil
		},
		gauge: &mockGauge{},
		lggr:  lggr,
	}

	dm.emitDirSizeMetric(context.Background())
	assert.Equal(t, int64(42), dm.gauge.(*mockGauge).gotValue)

	assert.Len(t, observed.FilterMessage("Emitting vault directory size metric").All(), 1)
}

func TestDiskMonitor_emitDirSizeMetric_error(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)

	dm := &DiskMonitor{
		sizeOfDir: func() (int64, error) {
			return 0, errors.New("disk read error")
		},
		gauge: &mockGauge{},
		lggr:  lggr,
	}

	dm.emitDirSizeMetric(context.Background())
	assert.Equal(t, int64(0), dm.gauge.(*mockGauge).gotValue)

	assert.Len(t, observed.FilterMessage("Failed to measure vault directory size").All(), 1)
}
