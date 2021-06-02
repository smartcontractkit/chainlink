package adapters

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
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
	// TaskTypeRandom is the identifier for the Random adapter.
	TaskTypeRandom = models.MustNewTaskType("random")
	// TaskTypeCompare is the identifier for the Compare adapter.
	TaskTypeCompare = models.MustNewTaskType("compare")
	// TaskTypeQuotient is the identifier for the Quotient adapter.
	TaskTypeQuotient = models.MustNewTaskType("quotient")
	// TaskTypeResultCollect is the identifier for the ResultCollect adapter.
	TaskTypeResultCollect = models.MustNewTaskType("resultcollect")
)

// BaseAdapter is the minimum interface required to create an adapter. Only core
// adapters have this minimum requirement.
type BaseAdapter interface {
	TaskType() models.TaskType
	Perform(models.RunInput, *store.Store, *keystore.Master) models.RunOutput
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
	var err error
	mic := config.MinIncomingConfirmations()
	var mp *assets.Link

	ba := FindNativeAdapterFor(task)
	if ba != nil { // task is for native adapter
		err = unmarshalParams(task.Params, ba)
	} else { // task is for external adapter
		bt, bErr := orm.FindBridge(task.Type)
		if bErr != nil {
			return nil, fmt.Errorf("%s is not a supported adapter type", task.Type)
		}
		b := Bridge{BridgeType: bt, Params: task.Params}
		ba = &b
		mp = bt.MinimumContractPayment
		mic = b.Confirmations
	}

	if ba == nil {
		return nil, fmt.Errorf("%s is not a supported adapter type", task.Type)
	}

	pa := &PipelineAdapter{
		BaseAdapter: ba,
		minConfs:    mic,
		minPayment:  mp,
	}

	return pa, err
}

// FindNativeAdapterFor find the native adapter for a given task
func FindNativeAdapterFor(task models.TaskSpec) BaseAdapter {
	switch task.Type {
	case TaskTypeCopy:
		return &Copy{}
	case TaskTypeEthBool:
		return &EthBool{}
	case TaskTypeEthBytes32:
		return &EthBytes32{}
	case TaskTypeEthInt256:
		return &EthInt256{}
	case TaskTypeEthUint256:
		return &EthUint256{}
	case TaskTypeEthTx:
		return &EthTx{}
	case TaskTypeHTTPGetWithUnrestrictedNetworkAccess:
		return &HTTPGet{AllowUnrestrictedNetworkAccess: true}
	case TaskTypeHTTPPostWithUnrestrictedNetworkAccess:
		return &HTTPPost{AllowUnrestrictedNetworkAccess: true}
	case TaskTypeHTTPGet:
		return &HTTPGet{}
	case TaskTypeHTTPPost:
		return &HTTPPost{}
	case TaskTypeJSONParse:
		return &JSONParse{}
	case TaskTypeMultiply:
		return &Multiply{}
	case TaskTypeNoOp:
		return &NoOp{}
	case TaskTypeNoOpPendOutgoing:
		return &NoOpPendOutgoing{}
	case TaskTypeSleep:
		return &Sleep{}
	case TaskTypeRandom:
		return &Random{}
	case TaskTypeCompare:
		return &Compare{}
	case TaskTypeQuotient:
		return &Quotient{}
	case TaskTypeResultCollect:
		return &ResultCollect{}
	default:
		return nil
	}
}

func unmarshalParams(params models.JSON, dst interface{}) error {
	bytes, err := params.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}
