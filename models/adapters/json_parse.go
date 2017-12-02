package adapters

type JsonParse struct {
	Path []string `json:"path"`
}

func (self *JsonParse) Perform(input RunResult) RunResult {
	return RunResult{}
}
