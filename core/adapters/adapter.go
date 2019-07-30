package adapters

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

var (
	// TaskTypeCopy is the identifier for the Copy adapter.
	TaskTypeCopy = models.MustNewTaskType("copy")
	// TaskTypeEthBool is the identifier for the EthBool adapter.
	TaskTypeEthBool = models.MustNewTaskType("ethbool")
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
	// TaskTypeRandom is the identifier for the Random adapter.
	TaskTypeRandom = models.MustNewTaskType("random")
	// TaskTypeCondition is the identifier for the Condition adapter.
	TaskTypeCondition = models.MustNewTaskType("condition")
)

// BaseAdapter is the minimum interface required to create an adapter. Only core
// adapters have this minimum requirement.
type BaseAdapter interface {
	Perform(models.RunResult, *store.Store) models.RunResult
}

// PipelineAdapter wraps a BaseAdapter with requirements for execution in the pipeline.
type PipelineAdapter struct {
	BaseAdapter
	minConfs           uint32
	minContractPayment *assets.Link
}

// MinConfs returns the private attribute
func (p PipelineAdapter) MinConfs() uint32 {
	return p.minConfs
}

// MinContractPayment returns the private attribute
func (p PipelineAdapter) MinContractPayment() *assets.Link {
	return p.minContractPayment
}

// For determines the adapter type to use for a given task.
func For(task models.TaskSpec, store *store.Store) (*PipelineAdapter, error) {
	var ba BaseAdapter
	var err error
	mic := store.Config.MinIncomingConfirmations()
	mcp := assets.NewLink(0)

	switch task.Type {
	case TaskTypeCopy:
		ba = &Copy{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeEthBool:
		ba = &EthBool{}
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
		mcp = store.Config.MinimumContractPayment()
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
	case TaskTypeRandom:
		ba = &Random{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeCondition:
		ba = &Condition{}
		err = unmarshalParams(task.Params, ba)
	default:
		bt, err := store.FindBridge(task.Type)
		if err != nil {
			return nil, fmt.Errorf("%s is not a supported adapter type", task.Type)
		}
		b := Bridge{BridgeType: &bt, Params: &task.Params}
		ba = &b
		mic = b.Confirmations
		mcp = bt.MinimumContractPayment
	}

	pa := &PipelineAdapter{
		BaseAdapter:        ba,
		minConfs:           mic,
		minContractPayment: mcp,
	}

	return pa, err
}

func unmarshalParams(params models.JSON, dst interface{}) error {
	bytes, err := params.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}
