package v2

import "strings"

// Engine initialization failure reasons labeling the
// platform_engine_workflow_initialization_failures_total counter. Values must
// stay low-cardinality: these are metric labels, not logs.
const (
	// initFailureWorkflowLimit labels the node's workflow-count limit being
	// reached before initialization began.
	initFailureWorkflowLimit = "workflow_limit"
	// initFailureDONSubscribe labels a failure subscribing to DON
	// configuration notifications.
	initFailureDONSubscribe = "don_subscribe"
	// initFailureSubscribe labels the module's Subscribe (trigger
	// subscription) request failing for any reason not classified below.
	initFailureSubscribe = "subscribe"
	// initFailureTriggerRegistration labels a failure registering the
	// module's trigger subscriptions with the capabilities registry.
	initFailureTriggerRegistration = "trigger_registration"
	// initFailureDisallowedCapabilityCall labels a module that attempted a
	// capability call during the trigger subscription phase.
	initFailureDisallowedCapabilityCall = "disallowed_capability_call_during_subscription"
	// initFailureDisallowedSecretsCall labels a module that attempted a
	// secrets call during the trigger subscription phase.
	initFailureDisallowedSecretsCall = "disallowed_secrets_call_during_subscription"
)

// initFailureReasonForSubscribe classifies a failed Subscribe phase for the
// initialization failure counter. The disallowed-host-call rejections cross
// the WASM boundary as plain strings (the guest re-reports the host's error
// text), so they are matched by text rather than errors.Is.
func initFailureReasonForSubscribe(err error) string {
	if err == nil {
		return initFailureSubscribe
	}
	switch {
	case strings.Contains(err.Error(), ErrSecretsCallDuringSubscription.Error()):
		return initFailureDisallowedSecretsCall
	case strings.Contains(err.Error(), ErrCapabilityCallDuringSubscription.Error()):
		return initFailureDisallowedCapabilityCall
	default:
		return initFailureSubscribe
	}
}
