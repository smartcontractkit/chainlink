package presenters

type Metric struct {
	JAID
	Name   string   `json:"name"`
	Help   string   `json:"help"`
	Type   string   `json:"type"`
	Labels []string `json:"labels"`
}

// GetName implements the api2go EntityNamer interface
func (m Metric) GetName() string {
	return "metrics"
}
