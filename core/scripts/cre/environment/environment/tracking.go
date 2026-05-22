package environment

const (
	MetricChipIngressStackStart = "cre.local.chip_ingress_stack.startup.result"
	// MetricBeholderStart is kept for backwards-compatible DX dashboards; emitted alongside MetricChipIngressStackStart.
	MetricBeholderStart = "cre.local.beholder.startup.result"
	MetricBillingStart  = "cre.local.billing.startup.result"

	MetricWorkflowDeploy = "cre.local.workflow.deploy"

	MetricStartupResult = "cre.local.startup.result"
	MetricStartupTime   = "cre.local.startup.time"

	MetricSetupResult = "cre.local.setup.result"

	MetricCapabilitySwap = "cre.local.env.swap.capability"
	MetricNodeSwap       = "cre.local.env.swap.nodes"

	// getDX configuration details
	GetDXGitHubVariableName = "API_TOKEN_LOCAL_CRE"
	GetDXProductName        = "local_cre"
)
