package adapters

type HttpGet struct {
	Endpoint string `json:"endpoint"`
}

func (self *HttpGet) Perform(input RunResult) RunResult {
	return RunResult{}
}
