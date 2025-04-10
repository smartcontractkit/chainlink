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
	LogTriggerCapability    CapabilityFlag = "log-trigger"
	WebAPITargetCapability  CapabilityFlag = "web-api-target"
	WebAPITriggerCapability CapabilityFlag = "web-api-trigger"
	MockCapability          CapabilityFlag = "mock"
	// Add more capabilities as needed
)

// Job names for which there are no specific capabilities
const (
	GatewayJobName = "gateway"
)
