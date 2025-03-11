package types

type CapabilityFlag = string

// DON types
const (
	WorkflowDON     CapabilityFlag = "workflow"
	CapabilitiesDON CapabilityFlag = "capabilities"
	GatewayDON      CapabilityFlag = "gateway"
)

// Capabilities
const (
	OCR3Capability          CapabilityFlag = "ocr3"
	CronCapability          CapabilityFlag = "cron"
	CustomComputeCapability CapabilityFlag = "custom-compute"
	WriteEVMCapability      CapabilityFlag = "write-evm"

	// Add more capabilities as needed
)

var (
	// Add new capabilities here as well, if single DON should have them by default
	SingleDonFlags = []string{"capabilities", "ocr3", "cron", "custom-compute", "write-evm"}
)
