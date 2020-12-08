package services

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"gopkg.in/guregu/null.v4"
)

var (
	DirectRequestLogTopic = getDirectRequestLogTopic()
)

func getDirectRequestLogTopic() gethCommon.Hash {
	abi, err := abi.JSON(strings.NewReader(offchainaggregator.OffchainAggregatorABI))
	if err != nil {
		panic("could not parse OffchainAggregator ABI: " + err.Error())
	}
	return abi.Events["ConfigSet"].ID
}

// EthRequestEvent is a wrapper for `models.EthRequestEvent`, the DB
// representation of the job spec. It fulfills the job.Spec interface
// and has facilities for unmarshaling the pipeline DAG from the job spec text.
type EthRequestEventSpec struct {
	Type            string          `toml:"type"`
	SchemaVersion   uint32          `toml:"schemaVersion"`
	Name            null.String     `toml:"name"`
	MaxTaskDuration models.Interval `toml:"maxTaskDuration"`

	models.EthRequestEventSpec

	// The `jobID` field exists to cache the ID from the jobs table that joins
	// to the eth_request_events table.
	jobID int32

	// The `Pipeline` field is only used during unmarshaling.  A pipeline.TaskDAG
	// is a type that implements gonum.org/v1/gonum/graph#Graph, which means that
	// you can dot.Unmarshal(...) raw DOT source directly into it, and it will
	// be a fully-instantiated DAG containing information about all of the nodes
	// and edges described by the DOT.  Our pipeline.TaskDAG type has a method
	// called `.TasksInDependencyOrder()` which converts this node/edge data
	// structure into task specs which can then be saved to the database.
	Pipeline pipeline.TaskDAG `toml:"observationSource"`
}

// EthRequestEventSpec conforms to the job.Spec interface
var _ job.Spec = EthRequestEventSpec{}

func (spec EthRequestEventSpec) JobID() int32 {
	return spec.jobID
}

func (spec EthRequestEventSpec) JobType() job.Type {
	return models.EthRequestEventJobType
}

func (spec EthRequestEventSpec) TaskDAG() pipeline.TaskDAG {
	return spec.Pipeline
}

type ethRequestEventSpecDelegate struct {
	logBroadcaster eth.LogBroadcaster
	pipelineRunner pipeline.Runner
}

func (d *ethRequestEventSpecDelegate) JobType() job.Type {
	return models.EthRequestEventJobType
}

func (d *ethRequestEventSpecDelegate) ToDBRow(spec job.Spec) models.JobSpecV2 {
	concreteSpec, ok := spec.(EthRequestEventSpec)
	if !ok {
		panic(fmt.Sprintf("expected a services.EthRequestEventSpec, got %T", spec))
	}
	return models.JobSpecV2{
		EthRequestEventSpec: &concreteSpec.EthRequestEventSpec,
		Type:                string(models.EthRequestEventJobType),
		SchemaVersion:       concreteSpec.SchemaVersion,
		MaxTaskDuration:     concreteSpec.MaxTaskDuration,
	}
}

func (d *ethRequestEventSpecDelegate) FromDBRow(spec models.JobSpecV2) job.Spec {
	if spec.EthRequestEventSpec == nil {
		return nil
	}
	return &EthRequestEventSpec{
		EthRequestEventSpec: *spec.EthRequestEventSpec,
		jobID:               spec.ID,
	}
}

// ServicesForSpec TODO
func (d *ethRequestEventSpecDelegate) ServicesForSpec(spec job.Spec) (services []job.Service, err error) {
	concreteSpec, is := spec.(*EthRequestEventSpec)
	if !is {
		return nil, errors.Errorf("services.ethRequestEventSpecDelegate expects a *services.EthRequestEventSpec, got %T", spec)
	}

	logListener := directRequestListener{
		d.logBroadcaster,
		concreteSpec.ContractAddress.Address(),
		d.pipelineRunner,
		spec.JobID(),
	}
	services = append(services, logListener)

	return
}

func RegisterEthRequestEventDelegate(jobSpawner job.Spawner, logBroadcaster eth.LogBroadcaster, pipelineRunner pipeline.Runner) {
	jobSpawner.RegisterDelegate(
		NewEthRequestEventDelegate(jobSpawner, logBroadcaster, pipelineRunner),
	)
}

func NewEthRequestEventDelegate(jobSpawner job.Spawner, logBroadcaster eth.LogBroadcaster, pipelineRunner pipeline.Runner) *ethRequestEventSpecDelegate {
	return &ethRequestEventSpecDelegate{
		logBroadcaster,
		pipelineRunner,
	}
}

var (
	_ eth.LogListener = &directRequestListener{}
	_ job.Service     = &directRequestListener{}
)

type directRequestListener struct {
	logBroadcaster  eth.LogBroadcaster
	contractAddress gethCommon.Address
	pipelineRunner  pipeline.Runner
	jobID           int32
}

// Start complies with job.Service
func (d directRequestListener) Start() error {
	connected := d.logBroadcaster.Register(d.contractAddress, d)
	if !connected {
		return errors.New("Failed to register directRequestListener with logBroadcaster")
	}
	return nil
}

// Close complies with job.Service
func (d directRequestListener) Close() error {
	d.logBroadcaster.Unregister(d.contractAddress, d)
	return nil
}

// OnConnect complies with eth.LogListener
func (directRequestListener) OnConnect() {}

// OnDisconnect complies with eth.LogListener
func (directRequestListener) OnDisconnect() {}

// OnConnect complies with eth.LogListener
func (d directRequestListener) HandleLog(lb eth.LogBroadcast, err error) {
	if err != nil {
		logger.Errorw("DirectRequestListener: error in previous LogListener", "err", err)
		return
	}

	was, err := lb.WasAlreadyConsumed()
	if err != nil {
		logger.Errorw("DirectRequestListener: could not determine if log was already consumed", "error", err)
		return
	} else if was {
		return
	}

	log := lb.DecodedLog()
	if log == nil || reflect.ValueOf(log).IsNil() {
		logger.Error("HandleLog: ignoring nil value")
		return
	}

	switch log := log.(type) {
	case *contracts.LogOracleRequest:
		d.handleOracleRequest(log.ToOracleRequest())
		return
	case *contracts.LogCancelOracleRequest:
		d.handleCancelOracleRequest(log)
		return

	default:
		logger.Warnf("unexpected log type %T", log)
		return
	}

	err = lb.MarkConsumed()
	if err != nil {
		logger.Errorw("OCRContract: could not mark log consumed", "error", err)
	}
	return
}

// TODO: This _really_ need to be async otherwise it will block log broadcaster
func (d *directRequestListener) handleOracleRequest(req contracts.OracleRequest) {
	ctx := context.TODO()
	runID, err := d.pipelineRunner.CreateRun(ctx, d.jobID, nil)
	if err != nil {
		logger.Errorw("EthRequestEvent failed to create run", "err", err)
		return
	}

	err = d.pipelineRunner.AwaitRun(ctx, runID)
	if err != nil {
		logger.Errorw("EthRequestEvent failed awaiting run", "err", err)
		return
	}

	results, err := d.pipelineRunner.ResultsForRun(ctx, runID)
	if err != nil {
		logger.Errorw("EthRequestEvent failed to get results for run", "err", err)
		return
	} else if len(results) != 1 {
		logger.Errorw("EthRequestEvent should have a single output", "err", err, "jobID", d.JobID, "runID", runID, "results", results)
		return
	}
}

// JobID complies with eth.LogListener
func (directRequestListener) JobID() *models.ID {
	return nil
}

// JobSpecV2 complies with eth.LogListener
func (d directRequestListener) JobIDV2() int32 {
	return d.jobID
}

// IsV2Job complies with eth.LogListener
func (directRequestListener) IsV2Job() bool {
	return true
}
