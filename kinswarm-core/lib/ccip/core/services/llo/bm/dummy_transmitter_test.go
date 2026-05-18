package bm

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_DummyTransmitter(t *testing.T) {
	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	tr := NewTransmitter(lggr, "dummy")

	servicetest.Run(t, tr)

	err := tr.Transmit(
		testutils.Context(t),
		types.ConfigDigest{},
		42,
		ocr3types.ReportWithInfo[llotypes.ReportInfo]{},
		[]types.AttributedOnchainSignature{},
	)
	require.NoError(t, err)

	testutils.RequireLogMessage(t, observedLogs, "Transmit")
}
