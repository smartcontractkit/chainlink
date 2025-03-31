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

	ReadContractCapability  CapabilityFlag = "read-contract"
	LogTriggerCapability    CapabilityFlag = "log_trigger"
	WebAPITargetCapability  CapabilityFlag = "web_api_target"
	WebAPITriggerCapability CapabilityFlag = "web_api_trigger"
	// Add more capabilities as needed
)
