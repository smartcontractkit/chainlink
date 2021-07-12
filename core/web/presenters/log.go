package presenters

type ServiceLogConfigResource struct {
	JAID
	ServiceName []string `json:"serviceName"`
	LogLevel    []string `json:"logLevel"`
}

// GetName implements the api2go EntityNamer interface
func (r ServiceLogConfigResource) GetName() string {
	return "serviceLevelLogs"
}
