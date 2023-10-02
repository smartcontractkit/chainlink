package synchronization

// TelemetryType defines supported telemetry types
type TelemetryType string

const (
	EnhancedEA        TelemetryType = "enhanced-ea"
	FunctionsRequests TelemetryType = "functions-requests"
	EnhancedEAMercury TelemetryType = "enhanced-ea-mercury"
	OCR               TelemetryType = "ocr"
	OCR2Automation    TelemetryType = "ocr2-automation"
	OCR2Functions     TelemetryType = "ocr2-functions"
	OCR2Threshold     TelemetryType = "ocr2-threshold"
	OCR2S4            TelemetryType = "ocr2-s4"
	OCR2Median        TelemetryType = "ocr2-median"
	OCR3Mercury       TelemetryType = "ocr3-mercury"
	OCR2VRF           TelemetryType = "ocr2-vrf"
	AutomationCustom  TelemetryType = "automation-custom"
)
