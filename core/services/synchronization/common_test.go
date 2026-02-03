package synchronization

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelemetryTypeToDomainAndEntity(t *testing.T) {
	tests := []struct {
		name           string
		telemType      TelemetryType
		expectedDomain string
		expectedEntity string
		expectError    bool
	}{
		{
			name:           "OCR",
			telemType:      OCR,
			expectedDomain: "data-feeds.telemetry.ocr",
			expectedEntity: "offchainreporting.TelemetryWrapper",
		},
		{
			name:           "OCR2Median",
			telemType:      OCR2Median,
			expectedDomain: "data-feeds.telemetry.ocr2-median",
			expectedEntity: "offchainreporting2.TelemetryWrapper",
		},
		{
			name:           "OCR2Automation",
			telemType:      OCR2Automation,
			expectedDomain: "automations.telemetry.ocr2-automation",
			expectedEntity: "offchainreporting2.TelemetryWrapper",
		},
		{
			name:           "OCR2CCIPCommit",
			telemType:      OCR2CCIPCommit,
			expectedDomain: "ccip.telemetry.ocr2",
			expectedEntity: "offchainreporting2.TelemetryWrapper",
		},
		{
			name:           "OCR2CCIPExec",
			telemType:      OCR2CCIPExec,
			expectedDomain: "ccip.telemetry.ocr2",
			expectedEntity: "offchainreporting2.TelemetryWrapper",
		},
		{
			name:           "OCR2Functions",
			telemType:      OCR2Functions,
			expectedDomain: "functions.telemetry.ocr2-functions",
			expectedEntity: "offchainreporting2.TelemetryWrapper",
		},
		{
			name:           "OCR3Automation",
			telemType:      OCR3Automation,
			expectedDomain: "automations.telemetry.ocr3-automation",
			expectedEntity: "offchainreporting3.TelemetryWrapper",
		},
		{
			name:           "OCR3Mercury",
			telemType:      OCR3Mercury,
			expectedDomain: "data-streams.telemetry.ocr3-mercury",
			expectedEntity: "offchainreporting3.TelemetryWrapper",
		},
		{
			name:           "OCR3CCIPCommit",
			telemType:      OCR3CCIPCommit,
			expectedDomain: "ccip.telemetry.ocr3",
			expectedEntity: "offchainreporting3.TelemetryWrapper",
		},
		{
			name:           "OCR3CCIPExec",
			telemType:      OCR3CCIPExec,
			expectedDomain: "ccip.telemetry.ocr3",
			expectedEntity: "offchainreporting3.TelemetryWrapper",
		},
		{
			name:           "OCR3DataFeeds",
			telemType:      OCR3DataFeeds,
			expectedDomain: "data-feeds.telemetry.ocr3-data-feeds",
			expectedEntity: "offchainreporting3.TelemetryWrapper",
		},
		{
			name:           "EnhancedEA",
			telemType:      EnhancedEA,
			expectedDomain: "data-feeds.telemetry.enhanced-ea",
			expectedEntity: "telem.EnhancedEA",
		},
		{
			name:           "EnhancedEAMercury",
			telemType:      EnhancedEAMercury,
			expectedDomain: "data-streams.telemetry.enhanced-ea-mercury",
			expectedEntity: "telem.EnhancedEAMercury",
		},
		{
			name:           "LLOObservation",
			telemType:      LLOObservation,
			expectedDomain: "data-streams.telemetry.llo-observation",
			expectedEntity: "telem.LLOObservationTelemetry",
		},
		{
			name:           "LLOOutcome",
			telemType:      LLOOutcome,
			expectedDomain: "data-streams.telemetry.llo-outcome",
			expectedEntity: "telem.LLOOutcomeTelemetry",
		},
		{
			name:           "AutomationCustom",
			telemType:      AutomationCustom,
			expectedDomain: "automations.telemetry.automation-custom",
			expectedEntity: "telem.AutomationTelemWrapper",
		},
		{
			name:           "FunctionsRequests",
			telemType:      FunctionsRequests,
			expectedDomain: "functions.telemetry.functions-requests",
			expectedEntity: "telem.FunctionsRequest",
		},
		{
			name:           "HeadReport",
			telemType:      HeadReport,
			expectedDomain: "ccip.telemetry.head-report",
			expectedEntity: "telem.HeadReportRequest",
		},
		{
			name:           "PipelineBridge",
			telemType:      PipelineBridge,
			expectedDomain: "data-streams.telemetry.pipeline-bridge",
			expectedEntity: "telem.LLOBridgeTelemetry",
		},
		{
			name:        "Unknown telemetry type",
			telemType:   TelemetryType("unknown"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain, entity, err := TelemetryTypeToDomainAndEntity(tt.telemType)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "could not resolve domain and entity from telem type")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedDomain, domain)
				assert.Equal(t, tt.expectedEntity, entity)
			}
		})
	}
}
