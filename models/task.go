package models

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink-go/adapters"
)

type Task struct {
	Type   string          `json:"type" storm:"index"`
	Params json.RawMessage `json:"params,omitempty"`
}

type TaskRun struct {
	Task
	ID     string `storm:"id"`
	Status string
	Result adapters.RunResult
}

func (self Task) Validate() error {
	_, err := self.Adapter()
	return err
}

func (self Task) Adapter() (adapters.Adapter, error) {
	switch self.Type {
	case "HttpGet":
		temp := &adapters.HttpGet{}
		err := json.Unmarshal(self.Params, temp)
		return temp, err
	case "JsonParse":
		temp := &adapters.JsonParse{}
		err := json.Unmarshal(self.Params, temp)
		return temp, err
	case "EthBytes32":
		temp := &adapters.EthBytes32{}
		err := unmarshalOrEmpty(self.Params, temp)
		return temp, err
	case "EthSendTx":
		temp := &adapters.EthSendTx{}
		err := json.Unmarshal(self.Params, temp)
		return temp, err
	case "NoOp":
		return &adapters.NoOp{}, nil
	}

	return nil, fmt.Errorf("%s is not a supported adapter type", self.Type)
}

func unmarshalOrEmpty(params json.RawMessage, dst interface{}) error {
	if len(params) > 0 {
		return json.Unmarshal(params, dst)
	}
	return nil
}
