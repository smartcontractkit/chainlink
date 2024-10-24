package trigger_test_utils

import (
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/webapicap"
)

func NewWorkflowTriggerConfig(addresses []string, topics []string) (webapicap.TriggerConfig, *values.Map, error) {
	triggerConfig := webapicap.TriggerConfig{
		AllowedSenders: addresses,
		AllowedTopics:  topics,
		RateLimiter: webapicap.RateLimiterConfig{
			GlobalRPS:      100.0,
			GlobalBurst:    101,
			PerSenderRPS:   102.0,
			PerSenderBurst: 103,
		},
		RequiredParams: []string{"bid", "ask"},
	}

	var rateLimitConfig, err = values.NewMap(map[string]any{
		"GlobalRPS":      100.0,
		"GlobalBurst":    101,
		"PerSenderRPS":   102.0,
		"PerSenderBurst": 103,
	})
	if err != nil {
		return triggerConfig, nil, err
	}

	triggerRegistrationConfig, err := values.NewMap(map[string]interface{}{
		"RateLimiter":    rateLimitConfig,
		"AllowedSenders": addresses,
		"AllowedTopics":  topics,
		"RequiredParams": []string{"bid", "ask"},
	})
	return triggerConfig, triggerRegistrationConfig, err
}
