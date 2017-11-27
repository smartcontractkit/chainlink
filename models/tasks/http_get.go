package tasks

type HttpGet struct {
	Endpoint string `json:"endpoint"`
}

func (self *HttpGet) Perform() {
}
