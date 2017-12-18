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

type AdapterBase struct {
	Config config.Config
}

type Output map[string]null.String

type configSetter interface {
	setConfig(config.Config)
}

type adapterConfigSetter interface {
	Adapter
	configSetter
}

func For(task models.Task, cf config.Config) (Adapter, error) {
	var ac adapterConfigSetter
	var err error
	switch task.Type {
	case "HttpGet":
		ac = &HttpGet{}
		err = json.Unmarshal(task.Params, ac)
	case "JsonParse":
		ac = &JsonParse{}
		err = json.Unmarshal(task.Params, ac)
	case "EthBytes32":
		ac = &EthBytes32{}
		err = unmarshalOrEmpty(task.Params, ac)
	case "EthSendTx":
		ac = &EthSendTx{}
		err = json.Unmarshal(task.Params, ac)
	case "NoOp":
		ac, err = &NoOp{}, nil
	default:
		return nil, fmt.Errorf("%s is not a supported adapter type", task.Type)
	}

	ac.setConfig(cf)
	return ac, err
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

func (self AdapterBase) setConfig(cf config.Config) {
	self.Config = cf
}
