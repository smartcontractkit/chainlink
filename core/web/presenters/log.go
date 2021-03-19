package presenters

type LogResource struct {
	JAID
	LogLevel string `json:"logLevel"`
	LogSql   bool   `json:"logSql"`
}

// GetName implements the api2go EntityNamer interface
func (r LogResource) GetName() string {
	return "logs"
}
