package core_node_components

type PlatformOpts struct {
	// Platform is infrastructure deployment platform: docker or k8s
	Platform     string
	LabelFilters map[string]string
	LabelFilter  string
	LegendString string
	LabelQuery   string
}

// PlatformPanelOpts generate different queries for "docker" and "k8s" deployment platforms
func PlatformPanelOpts() PlatformOpts {
	po := PlatformOpts{
		LabelFilters: map[string]string{
			"env":          `=~"${env}"`,
			"cluster":      `=~"${cluster}"`,
			"blockchain":   `=~"${blockchain}"`,
			"product":      `=~"${product}"`,
			"network_type": `=~"${network_type}"`,
			"component":    `=~"${component}"`,
			"service":      `=~"${service}"`,
		},
	}
	for key, value := range po.LabelFilters {
		po.LabelQuery += key + value + ", "
	}
	return po
}
