package tasks

import (
	"encoding/json"
)

type HttpGet struct {
	Endpoint string `json:"endpoint"`
}

func (t *Task) AsHttpGet() (HttpGet, error) {
	rval := HttpGet{}
	err := json.Unmarshal(t.Params, &rval)
	return rval, err
}
