package adapters

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

var (
	// TaskTypeCopy is the identifier for the Copy adapter.
	TaskTypeCopy = models.MustNewTaskType("copy")
	// TaskTypeEthBytes32 is the identifier for the EthBytes32 adapter.
	TaskTypeEthBytes32 = models.MustNewTaskType("ethbytes32")
	// TaskTypeEthInt256 is the identifier for the EthInt256 adapter.
	TaskTypeEthInt256 = models.MustNewTaskType("ethint256")
	// TaskTypeEthUint256 is the identifier for the EthUint256 adapter.
	TaskTypeEthUint256 = models.MustNewTaskType("ethuint256")
	// TaskTypeEthTx is the identifier for the EthTx adapter.
	TaskTypeEthTx = models.MustNewTaskType("ethtx")
	// TaskTypeHTTPGet is the identifier for the HTTPGet adapter.
	TaskTypeHTTPGet = models.MustNewTaskType("httpget")
	// TaskTypeHTTPPost is the identifier for the HTTPPost adapter.
	TaskTypeHTTPPost = models.MustNewTaskType("httppost")
	// TaskTypeJSONParse is the identifier for the JSONParse adapter.
	TaskTypeJSONParse = models.MustNewTaskType("jsonparse")
	// TaskTypeMultiply is the identifier for the Multiply adapter.
	TaskTypeMultiply = models.MustNewTaskType("multiply")
	// TaskTypeNoOp is the identifier for the NoOp adapter.
	TaskTypeNoOp = models.MustNewTaskType("noop")
	// TaskTypeNoOpPend is the identifier for the NoOpPend adapter.
	TaskTypeNoOpPend = models.MustNewTaskType("nooppend")
	// TaskTypeSleep is the identifier for the Sleep adapter.
	TaskTypeSleep = models.MustNewTaskType("sleep")
	// TaskTypeWasm is the wasm interpereter adapter
	TaskTypeWasm = models.MustNewTaskType("wasm")
)

// BaseAdapter is the minimum interface required to create an adapter. Only core
// adapters have this minimum requirement.
type BaseAdapter interface {
	Perform(models.RunResult, *store.Store) models.RunResult
}

// PipelineAdapter is the interface required for an adapter to be run in
// the job pipeline with validation checks. In addition to the BaseAdapter
// interface, implementers must specify the number of confirmations required
// before the adapter can be run.
type PipelineAdapter interface {
	BaseAdapter
	MinConfs() uint64
}

// MinConfsWrappedAdapter allows for an adapter to be wrapped so that it meets
// the AdapterWithMinConfsInterface.
type MinConfsWrappedAdapter struct {
	BaseAdapter
	ConfiguredConfirmations uint64
}

// MinConfs specifies the number of block confirmations
// needed to run the adapter.
func (wa MinConfsWrappedAdapter) MinConfs() uint64 {
	return wa.ConfiguredConfirmations
}

// For determines the adapter type to use for a given task.
func For(task models.TaskSpec, store *store.Store) (PipelineAdapter, error) {
	var ba BaseAdapter
	var err error
	switch task.Type {
	case TaskTypeCopy:
		ba = &Copy{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeEthBytes32:
		ba = &EthBytes32{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeEthInt256:
		ba = &EthInt256{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeEthUint256:
		ba = &EthUint256{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeEthTx:
		ba = &EthTx{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeHTTPGet:
		ba = &HTTPGet{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeHTTPPost:
		ba = &HTTPPost{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeJSONParse:
		ba = &JSONParse{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeMultiply:
		ba = &Multiply{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeNoOp:
		ba = &NoOp{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeNoOpPend:
		ba = &NoOpPend{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeSleep:
		ba = &Sleep{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeWasm:
		ba = &Wasm{}
		err = unmarshalParams(task.Params, ba)
	default:
		bt, err := store.FindBridge(task.Type.String())
		if err != nil {
			return nil, fmt.Errorf("%s is not a supported adapter type", task.Type)
		}
		return &Bridge{BridgeType: bt, Params: &task.Params}, nil
	}
	a := MinConfsWrappedAdapter{
		BaseAdapter:             ba,
		ConfiguredConfirmations: store.Config.MinIncomingConfirmations,
	}
	return a, err
}

func unmarshalParams(params models.JSON, dst interface{}) error {
	bytes, err := params.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}
