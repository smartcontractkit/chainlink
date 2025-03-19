package aggregation

import "testing"

func TestSignedReportAggregator_Aggregate(t *testing.T) {
	// TODO
	// Prepare signed report in expected format
	//   TriggerResponse { TriggerEvent { OCRTriggerEvent { marshalled OCRTriggerReport { Outputs }}}}
	// Expect a "flattened" response - OCR structs removed and Outputs lifted up to TriggerEvent
	//   TriggerResponse { TriggerEvent { Outputs }}
}
