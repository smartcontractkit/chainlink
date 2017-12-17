package adapters

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink-go/config"
	"github.com/smartcontractkit/chainlink-go/models"
	"gopkg.in/guregu/null.v3"
)

type Adapter interface {
	Perform(models.RunResult) models.RunResult
}

type adapterBase struct {
	*config.Config
}

type Output map[string]null.String

func For(task models.Task, cf config.Config) (Adapter, error) {
	switch task.Type {
	case "HttpGet":
		temp := &HttpGet{}
		err := json.Unmarshal(task.Params, temp)
		return temp, err
	case "JsonParse":
		temp := &JsonParse{}
		err := json.Unmarshal(task.Params, temp)
		return temp, err
	case "EthBytes32":
		temp := &EthBytes32{}
		err := unmarshalOrEmpty(task.Params, temp)
		return temp, err
	case "EthSendTx":
		temp := &EthSendTx{}
		err := json.Unmarshal(task.Params, temp)
		temp.Config = cf
		return temp, err
	case "NoOp":
		return &NoOp{}, nil
	}

	return nil, fmt.Errorf("%s is not a supported adapter type", task.Type)
}

func unmarshalOrEmpty(params json.RawMessage, dst interface{}) error {
	if len(params) > 0 {
		return json.Unmarshal(params, dst)
	}
	return nil
}

func Validate(job models.Job) error {
	var err error
	for _, task := range job.Tasks {
		err = validateTask(task)
		if err != nil {
			break
		}
	}

	return err
}

func validateTask(task models.Task) error {
	_, err := For(task, config.Config{})
	return err
}
