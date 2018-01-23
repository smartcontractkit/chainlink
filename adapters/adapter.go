package adapters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// The Adapter interface applies to all core adapters.
// Each implementation must return a RunResult.
type Adapter interface {
	Perform(models.RunResult, *store.Store) models.RunResult
}

// For determines the adapter type to use for a given task
func For(task models.Task, store *store.Store) (ac Adapter, err error) {
	switch strings.ToLower(task.Type) {
	case "httpget":
		ac = &HttpGet{}
		err = json.Unmarshal(task.Params, ac)
	case "jsonparse":
		ac = &JsonParse{}
		err = json.Unmarshal(task.Params, ac)
	case "ethbytes32":
		ac = &EthBytes32{}
		err = unmarshalOrEmpty(task.Params, ac)
	case "ethtx":
		ac = &EthTx{}
		err = unmarshalOrEmpty(task.Params, ac)
	case "noop":
		ac = &NoOp{}
		err = unmarshalOrEmpty(task.Params, ac)
	case "nooppend":
		ac = &NoOpPend{}
		err = unmarshalOrEmpty(task.Params, ac)
	default:
		if tt, err := store.TaskTypeFor(task.Type); err != nil {
			return nil, fmt.Errorf("%s is not a supported adapter type", task.Type)
		} else {
			ac = &ExternalBridge{tt}
			err = unmarshalOrEmpty(task.Params, ac)
		}
	}
	return ac, err
}

func unmarshalOrEmpty(params json.RawMessage, dst interface{}) error {
	if len(params) > 0 {
		return json.Unmarshal(params, dst)
	}
	return nil
}

// Validate that there were no errors in any of the tasks of a job
func Validate(job *models.Job, store *store.Store) error {
	var err error
	for _, task := range job.Tasks {
		err = validateTask(task, store)
		if err != nil {
			break
		}
	}

	return err
}

func validateTask(task models.Task, store *store.Store) error {
	_, err := For(task, store)
	return err
}
