package presenters

import "cmp"

type Check struct {
	JAID
	Name   string `json:"name"`
	Status string `json:"status"`
	Output string `json:"output"`
}

func (c Check) GetName() string {
	return "checks"
}

func CmpCheckName(a, b Check) int {
	return cmp.Compare(a.Name, b.Name)
}
