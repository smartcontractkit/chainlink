package tasks

type JsonParse struct {
	Path []string `json:"path"`
}

func (self *JsonParse) Perform() {
}
