package tasks

import (
	"encoding/json"
	"errors"
)

type Adapter interface {
	Perform()
}

type Task struct {
	Type   string          `json:"type" storm:"index"`
	Params json.RawMessage `json:"params"`
	Adapter
}

func (t *Task) Valid() bool {
	switch t.Type {
	case "HttpGet":
		_, err := t.AsHttpGet()
		return err == nil
	}
	return false
}

func (self *Task) UnmarshalJSON(b []byte) error {
	type tempTask Task
	err := json.Unmarshal(b, (*tempTask)(self))
	if err != nil {
		return err
	}
	self.Adapter, err = self.adapterFromRaw()
	return err
}

func (self *Task) adapterFromRaw() (Adapter, error) {
	switch self.Type {
	case "HttpGet":
		temp := &HttpGet{}
		err := json.Unmarshal(self.Params, temp)
		return temp, err
	}

	return nil, errors.New(self.Type + " is not a supported adapter type")
}
