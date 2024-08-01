package ocr2keepers

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// Generate types from third-party repos:
//
//go:generate mockery --name Logger --structname MockLogger --srcpkg "github.com/smartcontractkit/libocr/commontypes" --case underscore --filename logger.generated.go

type logWriter struct {
	l commontypes.Logger
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	l.l.Debug(string(p), nil)
	n = len(p)
	return
}

type ocrLogContextKey struct{}

type ocrLogContext struct {
	Epoch     uint32
	Round     uint8
	StartTime time.Time
}

func newOcrLogContext(rt types.ReportTimestamp) ocrLogContext {
	return ocrLogContext{
		Epoch:     rt.Epoch,
		Round:     rt.Round,
		StartTime: time.Now(),
	}
}

func (c ocrLogContext) String() string {
	return fmt.Sprintf("[epoch=%d, round=%d, completion=%dms]", c.Epoch, c.Round, time.Since(c.StartTime)/time.Millisecond)
}

func (c ocrLogContext) Short() string {
	return fmt.Sprintf("[epoch=%d, round=%d]", c.Epoch, c.Round)
}
