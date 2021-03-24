package presenters

type LogResource struct {
	JAID
	Level      string `json:"level"`
	SqlEnabled bool   `json:"sqlEnabled"`
}

// GetName implements the api2go EntityNamer interface
func (r LogResource) GetName() string {
	return "logs"
}
