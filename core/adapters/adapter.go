package adapters

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
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
	// TaskTypeEthTxABIEncode is the identifier for the EthTxABIEncode adapter.
	TaskTypeEthTxABIEncode = models.MustNewTaskType("ethtxabiencode")
	// TaskTypeHTTPGetWithUnrestrictedNetworkAccess is the identifier for the HTTPGet adapter, with local/private IP access enabled.
	TaskTypeHTTPGetWithUnrestrictedNetworkAccess = models.MustNewTaskType("httpgetwithunrestrictednetworkaccess")
	// TaskTypeHTTPPostWithUnrestrictedNetworkAccess is the identifier for the HTTPPost adapter, with local/private IP access enabled.
	TaskTypeHTTPPostWithUnrestrictedNetworkAccess = models.MustNewTaskType("httppostwithunrestrictednetworkaccess")
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
	// TaskTypeNoOpPendOutgoing is the identifier for the NoOpPendOutgoing adapter.
	TaskTypeNoOpPendOutgoing = models.MustNewTaskType("nooppendoutgoing")
	// TaskTypeSleep is the identifier for the Sleep adapter.
	TaskTypeSleep = models.MustNewTaskType("sleep")
	// TaskTypeWasm is the wasm interpereter adapter
	TaskTypeWasm = models.MustNewTaskType("wasm")
	// TaskTypeRandom is the identifier for the Random adapter.
	TaskTypeRandom = models.MustNewTaskType("random")
	// TaskTypeCompare is the identifier for the Compare adapter.
	TaskTypeCompare = models.MustNewTaskType("compare")
	// TaskTypeQuotient is the identifier for the Quotient adapter.
	TaskTypeQuotient = models.MustNewTaskType("quotient")
)

// BaseAdapter is the minimum interface required to create an adapter. Only core
// adapters have this minimum requirement.
type BaseAdapter interface {
	TaskType() models.TaskType
	Perform(models.RunInput, *store.Store) models.RunOutput
}

// PipelineAdapter wraps a BaseAdapter with requirements for execution in the pipeline.
type PipelineAdapter struct {
	BaseAdapter
	minConfs   uint32
	minPayment *assets.Link
}

// MinConfs returns the private attribute
func (p PipelineAdapter) MinConfs() uint32 {
	return p.minConfs
}

// MinPayment returns the payment for this adapter (defaults to none)
func (p PipelineAdapter) MinPayment() *assets.Link {
	return p.minPayment
}

// For determines the adapter type to use for a given task.
func For(task models.TaskSpec, config orm.ConfigReader, orm *orm.ORM) (*PipelineAdapter, error) {
	var ba BaseAdapter
	var err error
	mic := config.MinIncomingConfirmations()
	var mp *assets.Link

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
		err = unmarshalParams(task.Params, ba)
	case TaskTypeEthTxABIEncode:
		ba = &EthTxABIEncode{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeHTTPGetWithUnrestrictedNetworkAccess:
		ba = &HTTPGet{AllowUnrestrictedNetworkAccess: true}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeHTTPPostWithUnrestrictedNetworkAccess:
		ba = &HTTPPost{AllowUnrestrictedNetworkAccess: true}
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
	case TaskTypeNoOpPendOutgoing:
		ba = &NoOpPendOutgoing{}
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
	case TaskTypeCompare:
		ba = &Compare{}
		err = unmarshalParams(task.Params, ba)
	case TaskTypeQuotient:
		ba = &Quotient{}
		err = unmarshalParams(task.Params, ba)
	default:
		bt, e := orm.FindBridge(task.Type)
		if e != nil {
			return nil, fmt.Errorf("%s is not a supported adapter type", task.Type)
		}
		b := Bridge{BridgeType: bt, Params: task.Params}
		ba = &b
		mp = bt.MinimumContractPayment
		mic = b.Confirmations
	}

	pa := &PipelineAdapter{
		BaseAdapter: ba,
		minConfs:    mic,
		minPayment:  mp,
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
