package bridgestatus

// EAResponse represents the response schema from External Adapter status endpoint
type EAResponse struct {
	Adapter struct {
		Name          string  `json:"name"`
		Version       string  `json:"version"`
		UptimeSeconds float64 `json:"uptimeSeconds"`
	} `json:"adapter"`
	Endpoints []struct {
		Name       string   `json:"name"`
		Aliases    []string `json:"aliases"`
		Transports []string `json:"transports"`
	} `json:"endpoints"`
	DefaultEndpoint string `json:"defaultEndpoint"`
	Configuration   []struct {
		Name               string `json:"name"`
		Value              any    `json:"value"`
		Type               string `json:"type"`
		Description        string `json:"description"`
		Required           bool   `json:"required"`
		Default            any    `json:"default"`
		CustomSetting      bool   `json:"customSetting"`
		EnvDefaultOverride any    `json:"envDefaultOverride"`
	} `json:"configuration"`
	Runtime struct {
		NodeVersion  string `json:"nodeVersion"`
		Platform     string `json:"platform"`
		Architecture string `json:"architecture"`
		Hostname     string `json:"hostname"`
	} `json:"runtime"`
	Metrics struct {
		Enabled  bool    `json:"enabled"`
		Port     *int    `json:"port,omitempty"`
		Endpoint *string `json:"endpoint,omitempty"`
	} `json:"metrics"`
}

// JobInfo represents job information for a bridge
type JobInfo struct {
	ExternalJobID string
	Name          string
}
