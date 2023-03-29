package mercury

// ResourcesConfig used to configure different CPU/MEM settings for load/chaos/smoke
type ResourcesConfig struct {
	/* Chainlink nodes resources */
	DONResources   map[string]interface{}
	DONDBResources map[string]interface{}
	/* Mercury server resources */
	MercuryResources   map[string]interface{}
	MercuryDBResources map[string]interface{}
}

var (
	DefaultResources = &ResourcesConfig{
		DONResources: map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "500m",
				"memory": "1024Mi",
			},
			"limits": map[string]interface{}{
				"cpu":    "500m",
				"memory": "1024Mi",
			},
		},
		DONDBResources: map[string]interface{}{
			"stateful": "true",
			"capacity": "2Gi",
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu":    "500m",
					"memory": "1024Mi",
				},
				"limits": map[string]interface{}{
					"cpu":    "500m",
					"memory": "1024Mi",
				},
			},
		},
		MercuryResources: map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "500m",
				"memory": "1024Mi",
			},
			"limits": map[string]interface{}{
				"cpu":    "500m",
				"memory": "1024Mi",
			},
		},
		MercuryDBResources: map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "500m",
				"memory": "1024Mi",
			},
			"limits": map[string]interface{}{
				"cpu":    "500m",
				"memory": "1024Mi",
			},
		},
	}
)

// Get mockserver resources depending on number of feeds in the DON
func GetMockserverResources(feedCount int) map[string]interface{} {
	var cpu, memory string

	if feedCount > 4 {
		cpu = "8000m"
		memory = "8048Mi"
	} else if feedCount > 1 {
		cpu = "4000m"
		memory = "4048Mi"
	} else {
		cpu = "2000m"
		memory = "2560Mi"
	}

	return map[string]interface{}{
		"app": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    cpu,
				"memory": memory,
			},
			"limits": map[string]interface{}{
				"cpu":    cpu,
				"memory": memory,
			},
		},
	}
}
