package trigger_test_utils

import "github.com/smartcontractkit/chainlink-common/pkg/values"

func NewWorkflowTriggerConfig(addresses []string, topics []string) (*values.Map, error) {
	var rateLimitConfig, err = values.NewMap(map[string]any{
		"globalRPS":      100.0,
		"globalBurst":    101,
		"perSenderRPS":   102.0,
		"perSenderBurst": 103,
	})
	if err != nil {
		return nil, err
	}

	triggerRegistrationConfig, err := values.NewMap(map[string]interface{}{
		"rateLimiter":    rateLimitConfig,
		"allowedSenders": addresses,
		"allowedTopics":  topics,
		"requiredParams": []string{"bid", "ask"},
	})
	return triggerRegistrationConfig, err
}
