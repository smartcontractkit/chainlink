package bm

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func Test_DummyTransmitter(t *testing.T) {
	lggr, observedLogs := logger.TestObservedSugared(t, zapcore.DebugLevel)
	tr := NewTransmitter(lggr, "dummy")

	servicetest.Run(t, tr)

	err := tr.Transmit(
		t.Context(),
		types.ConfigDigest{},
		42,
		ocr3types.ReportWithInfo[llotypes.ReportInfo]{},
		[]types.AttributedOnchainSignature{},
	)
	require.NoError(t, err)

	tests.RequireLogMessage(t, observedLogs, "Transmit")
}
