package trigger_test_utils

import "github.com/smartcontractkit/chainlink-common/pkg/values"

func NewWorkflowTriggerConfig(addresses []string, topics []string) (*values.Map, error) {
	var rateLimitConfig, err = values.NewMap(map[string]any{
		"GlobalRPS":      100.0,
		"GlobalBurst":    101,
		"PerSenderRPS":   102.0,
		"PerSenderBurst": 103,
	})
	if err != nil {
		return nil, err
	}

	triggerRegistrationConfig, err := values.NewMap(map[string]interface{}{
		"RateLimiter":    rateLimitConfig,
		"AllowedSenders": addresses,
		"AllowedTopics":  topics,
		"RequiredParams": []string{"bid", "ask"},
	})
	return triggerRegistrationConfig, err
}
