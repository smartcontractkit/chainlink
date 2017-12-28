package adapters

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"gopkg.in/guregu/null.v3"
)

type Adapter interface {
	Perform(models.RunResult) models.RunResult
}

type AdapterBase struct {
	Store *store.Store
}

type Output map[string]null.String

type storeSetter interface {
	setStore(*store.Store)
}

type adapterStoreSetter interface {
	Adapter
	storeSetter
}

func For(task models.Task, s *store.Store) (Adapter, error) {
	var ac adapterStoreSetter
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
	case "EthConfirmTx":
		ac = &EthConfirmTx{}
		err = unmarshalOrEmpty(task.Params, ac)
	case "EthSendRawTx":
		ac = &EthSendRawTx{}
		err = unmarshalOrEmpty(task.Params, ac)
	case "EthSignTx":
		ac = &EthSignTx{}
		err = json.Unmarshal(task.Params, ac)
	case "EthSignAndSendTx":
		ac = &EthSignAndSendTx{}
		err = unmarshalOrEmpty(task.Params, ac)
	case "NoOp":
		ac = &NoOp{}
		err = unmarshalOrEmpty(task.Params, ac)
	default:
		return nil, fmt.Errorf("%s is not a supported adapter type", task.Type)
	}

	ac.setStore(s)
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
	_, err := For(task, nil)
	return err
}

func (self *AdapterBase) setStore(s *store.Store) {
	self.Store = s
}
