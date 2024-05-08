package types

// OCR2PluginType defines supported OCR2 plugin types.
type OCR2PluginType string

const (
	Median  OCR2PluginType = "median"
	DKG     OCR2PluginType = "dkg"
	OCR2VRF OCR2PluginType = "ocr2vrf"

	// TODO: sc-55296 to rename ocr2keeper to ocr2automation in code
	OCR2Keeper     OCR2PluginType = "ocr2automation"
	Functions      OCR2PluginType = "functions"
	Mercury        OCR2PluginType = "mercury"
	LLO            OCR2PluginType = "llo"
	GenericPlugin  OCR2PluginType = "plugin"
	OCR3Capability OCR2PluginType = "ocr3-capability"

	CCIPCommit    OCR2PluginType = "ccip-commit"
	CCIPExecution OCR2PluginType = "ccip-execution"
)
