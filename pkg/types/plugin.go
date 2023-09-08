package types

// OCR2PluginType defines supported OCR2 plugin types.
type OCR2PluginType string

const (
	Median  OCR2PluginType = "median"
	DKG     OCR2PluginType = "dkg"
	OCR2VRF OCR2PluginType = "ocr2vrf"

	// TODO: sc-55296 to rename ocr2keeper to ocr2automation in code
	OCR2Keeper    OCR2PluginType = "ocr2automation"
	Functions     OCR2PluginType = "functions"
	Mercury       OCR2PluginType = "mercury"
	GenericPlugin OCR2PluginType = "plugin"
)
