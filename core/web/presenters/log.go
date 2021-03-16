package presenters

type LogResource struct {
	JAID
	DebugEnabled bool `json:"debugEnabled"`
}

// GetName implements the api2go EntityNamer interface
func (r LogResource) GetName() string {
	return "logs"
}
