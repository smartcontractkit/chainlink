package synchronization

// TelemetryType defines supported telemetry types
type TelemetryType string

const (
	EnhancedEA     TelemetryType = "enhanced-ea"
	OCR            TelemetryType = "ocr"
	OCR2Automation TelemetryType = "ocr2-automation"
	OCR2Functions  TelemetryType = "ocr2-functions"
	OCR2Median     TelemetryType = "ocr2-median"
	OCR2Mercury    TelemetryType = "ocr2-mercury"
	OCR2VRF        TelemetryType = "ocr2-vrf"
)
