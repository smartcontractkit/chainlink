package globals

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_MergeWithOverrides(t *testing.T) {
	assert.Equal(t, ExecOCRParams.DeltaProgress, ExecOCRParamsForEthereum.DeltaProgress)
	assert.Equal(t, ExecOCRParams.DeltaResend, ExecOCRParamsForEthereum.DeltaResend)
	assert.Equal(t, ExecOCRParams.DeltaInitial, ExecOCRParamsForEthereum.DeltaInitial)
	assert.Equal(t, ExecOCRParams.DeltaGrace, ExecOCRParamsForEthereum.DeltaGrace)
	assert.Equal(t, ExecOCRParams.DeltaCertifiedCommitRequest, ExecOCRParamsForEthereum.DeltaCertifiedCommitRequest)
	assert.Equal(t, ExecOCRParams.MaxDurationQuery, ExecOCRParamsForEthereum.MaxDurationQuery)
	assert.Equal(t, ExecOCRParams.MaxDurationObservation, ExecOCRParamsForEthereum.MaxDurationObservation)
	assert.Equal(t, ExecOCRParams.MaxDurationShouldAcceptAttestedReport, ExecOCRParamsForEthereum.MaxDurationShouldAcceptAttestedReport)
	assert.Equal(t, ExecOCRParams.MaxDurationShouldTransmitAcceptedReport, ExecOCRParamsForEthereum.MaxDurationShouldTransmitAcceptedReport)
	assert.Equal(t, ExecOCRParams.MaxDurationQuery, ExecOCRParamsForEthereum.MaxDurationQuery)

	assert.Equal(t, 5*time.Second, ExecOCRParamsForEthereum.DeltaRound)
	assert.Equal(t, 25*time.Second, ExecOCRParamsForEthereum.DeltaStage)
	assert.Equal(t, 100*time.Millisecond, ExecOCRParams.MaxDurationQuery)
	assert.Equal(t, 100*time.Millisecond, ExecOCRParamsForEthereum.MaxDurationQuery)
}
