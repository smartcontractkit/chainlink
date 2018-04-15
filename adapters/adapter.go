package adapters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Adapter interface applies to all core adapters.
// Each implementation must return a RunResult.
type Adapter interface {
	Perform(models.RunResult, *store.Store) models.RunResult
}

// AdapterWithMinConfs is the interface required for an adapter to be run in
// the job pipeline. In addition to the Adapter interface, implementers must
// specify the number of confirmations required before the Adapter can be run.
type AdapterWithMinConfs interface {
	Adapter
	MinConfs() uint64
}

// MinConfsWrappedAdapter allows for an adapter to be wrapped so that it meets
// the AdapterWithMinConfsInterface.
type MinConfsWrappedAdapter struct {
	Adapter
	ConfiguredConfirmations uint64
}

// MinConfs specifies the number of block confirmations
// needed to run the adapter.
func (wa MinConfsWrappedAdapter) MinConfs() uint64 {
	return wa.ConfiguredConfirmations
}

// For determines the adapter type to use for a given task.
func For(task models.TaskSpec, store *store.Store) (AdapterWithMinConfs, error) {
	var ac Adapter
	var err error
	switch strings.ToLower(task.Type) {
	case "httpget":
		ac = &HTTPGet{}
		err = unmarshalParams(task.Params, ac)
	case "httppost":
		ac = &HTTPPost{}
		err = unmarshalParams(task.Params, ac)
	case "jsonparse":
		ac = &JSONParse{}
		err = unmarshalParams(task.Params, ac)
	case "ethbytes32":
		ac = &EthBytes32{}
		err = unmarshalParams(task.Params, ac)
	case "ethuint256":
		ac = &EthUint256{}
		err = unmarshalParams(task.Params, ac)
	case "ethtx":
		ac = &EthTx{}
		err = unmarshalParams(task.Params, ac)
	case "multiply":
		ac = &Multiply{}
		err = unmarshalParams(task.Params, ac)
	case "noop":
		ac = &NoOp{}
		err = unmarshalParams(task.Params, ac)
	case "nooppend":
		ac = &NoOpPend{}
		err = unmarshalParams(task.Params, ac)
	default:
		if bt, err := store.BridgeTypeFor(task.Type); err != nil {
			return nil, fmt.Errorf("%s is not a supported adapter type", task.Type)
		} else {
			return &Bridge{bt}, nil
		}
	}
	wa := MinConfsWrappedAdapter{
		Adapter:                 ac,
		ConfiguredConfirmations: store.Config.TaskMinConfirmations,
	}
	return wa, err
}

func unmarshalParams(params models.JSON, dst interface{}) error {
	bytes, err := params.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}
