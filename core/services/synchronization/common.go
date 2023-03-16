package synchronization

// TelemetryType defines supported telemetry types
type TelemetryType string

const (
	OCR            TelemetryType = "ocr"
	OCR2Median     TelemetryType = "ocr2-median"
	OCR2VRF        TelemetryType = "ocr2-vrf"
	OCR2Automation TelemetryType = "ocr2-automation"
	OCR2Functions  TelemetryType = "ocr2-functions"
)
