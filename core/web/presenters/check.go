package presenters

type Check struct {
	JAID
	Name   string `json:"name"`
	Status string `json:"status"`
	Output string `json:"output"`
}

func (c Check) GetName() string {
	return "checks"
}
