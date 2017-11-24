package tasks

import (
	"encoding/json"
)

type Task struct {
	Type   string          `json:"type" storm:"index"`
	Params json.RawMessage `json:"params"`
}

func (t *Task) Valid() bool {
	switch t.Type {
	case "HttpGet":
		_, err := t.AsHttpGet()
		return err == nil
	}
	return false
}
