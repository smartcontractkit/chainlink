package directrequest

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/oracle_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

type Delegate struct {
	logBroadcaster log.Broadcaster
	pipelineRunner pipeline.Runner
	pipelineORM    pipeline.ORM
	db             *gorm.DB
	ethClient      eth.Client
}

func NewDelegate(logBroadcaster log.Broadcaster, pipelineRunner pipeline.Runner, pipelineORM pipeline.ORM, ethClient eth.Client, db *gorm.DB) *Delegate {
	return &Delegate{
		logBroadcaster,
		pipelineRunner,
		pipelineORM,
		db,
		ethClient,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.DirectRequest
}

// ServicesForSpec returns the log listener service for a direct request job
func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.Service, err error) {
	if spec.DirectRequestSpec == nil {
		return nil, errors.Errorf("services.Delegate expects a *job.DirectRequestSpec to be present, got %v", spec)
	}
	concreteSpec := spec.DirectRequestSpec

	oracle, err := oracle_wrapper.NewOracle(concreteSpec.ContractAddress.Address(), d.ethClient)
	if err != nil {
		return
	}

	logListener := &listener{
		logBroadcaster: d.logBroadcaster,
		oracle:         oracle,
		pipelineRunner: d.pipelineRunner,
		db:             d.db,
		pipelineORM:    d.pipelineORM,
		spec:           *spec.PipelineSpec,
	}
	copy(logListener.onChainJobSpecID[:], spec.DirectRequestSpec.OnChainJobSpecID.Bytes())
	services = append(services, logListener)

	return
}

var (
	_ log.Listener = &listener{}
	_ job.Service  = &listener{}
)

type listener struct {
	logBroadcaster    log.Broadcaster
	unsubscribeLogs   func()
	oracle            oracle_wrapper.OracleInterface
	pipelineRunner    pipeline.Runner
	db                *gorm.DB
	pipelineORM       pipeline.ORM
	spec              pipeline.Spec
	onChainJobSpecID  common.Hash
	runs              sync.Map
	shutdownWaitGroup sync.WaitGroup
}

// Start complies with job.Service
func (l *listener) Start() error {
	connected, unsubscribe := l.logBroadcaster.Register(l, log.ListenerOpts{
		Contract: l.oracle,
		Logs: []generated.AbigenLog{
			oracle_wrapper.OracleOracleRequest{},
			oracle_wrapper.OracleCancelOracleRequest{},
		},
	})
	if !connected {
		return errors.New("Failed to register listener with logBroadcaster")
	}
	l.unsubscribeLogs = unsubscribe
	return nil
}

// Close complies with job.Service
func (l *listener) Close() error {
	if l.unsubscribeLogs != nil {
		l.unsubscribeLogs()
	}
	l.runs.Range(func(key, runCloserChannelIf interface{}) bool {
		runCloserChannel, _ := runCloserChannelIf.(chan struct{})
		close(runCloserChannel)
		return true
	})
	l.runs = sync.Map{}
	l.shutdownWaitGroup.Wait()
	return nil
}

// OnConnect complies with log.Listener
func (*listener) OnConnect() {}

// OnDisconnect complies with log.Listener
func (*listener) OnDisconnect() {}

// OnConnect complies with log.Listener
func (l *listener) HandleLog(lb log.Broadcast) {
	was, err := lb.WasAlreadyConsumed()
	if err != nil {
		logger.Errorw("DirectRequestListener: could not determine if log was already consumed", "error", err)
		return
	} else if was {
		return
	}

	logJobSpecID := lb.RawLog().Topics[1]
	if logJobSpecID == (common.Hash{}) || logJobSpecID != l.onChainJobSpecID {
		logger.Debugw("Skipping Run for Log with wrong Job ID", "logJobSpecID", logJobSpecID, "actualJobID", l.onChainJobSpecID)
		return
	}

	log := lb.DecodedLog()
	if log == nil || reflect.ValueOf(log).IsNil() {
		logger.Error("HandleLog: ignoring nil value")
		return
	}

	switch log := log.(type) {
	case *oracle_wrapper.OracleOracleRequest:
		l.handleOracleRequest(log)
		err = lb.MarkConsumed()
		if err != nil {
			logger.Errorf("Error marking log as consumed: %v", err)
		}
	case *oracle_wrapper.OracleCancelOracleRequest:
		l.handleCancelOracleRequest(log)
		err = lb.MarkConsumed()
		if err != nil {
			logger.Errorf("Error marking log as consumed: %v", err)
		}

	default:
		logger.Warnf("unexpected log type %T", log)
	}
}

func oracleRequestToMap(request *oracle_wrapper.OracleOracleRequest) map[string]interface{} {
	result := make(map[string]interface{})
	result["specId"] = fmt.Sprintf("0x%x", request.SpecId)
	result["requester"] = request.Requester.Hex()
	result["requestId"] = formatRequestId(request.RequestId)
	result["payment"] = fmt.Sprintf("%v", request.Payment)
	result["callbackAddr"] = request.CallbackAddr.Hex()
	result["callbackFunctionId"] = fmt.Sprintf("0x%x", request.CallbackFunctionId)
	result["cancelExpiration"] = fmt.Sprintf("%v", request.CancelExpiration)
	result["dataVersion"] = fmt.Sprintf("%v", request.DataVersion)
	result["data"] = fmt.Sprintf("0x%x", request.Data)
	return result
}

func (l *listener) handleOracleRequest(request *oracle_wrapper.OracleOracleRequest) {
	meta := make(map[string]interface{})
	meta["oracleRequest"] = oracleRequestToMap(request)

	logger := logger.CreateLogger(logger.Default.With(
		"jobName", l.spec.JobName,
		"jobID", l.spec.JobID,
	))

	l.shutdownWaitGroup.Add(1)
	go func() {
		defer l.shutdownWaitGroup.Done()

		runCloserChannel := make(chan struct{})
		runCloserChannelIf, loaded := l.runs.LoadOrStore(formatRequestId(request.RequestId), runCloserChannel)
		if loaded {
			runCloserChannel, _ = runCloserChannelIf.(chan struct{})
		}
		ctx, cancel := utils.CombinedContext(runCloserChannel, context.Background())
		defer cancel()

		_, _, err := l.pipelineRunner.ExecuteAndInsertNewRun(ctx, l.spec, pipeline.JSONSerializable{Val: meta, Null: false}, *logger)
		if ctx.Err() != nil {
			return
		} else if err != nil {
			logger.Errorw("DirectRequest failed to create run", "err", err)
		}
	}()
}

// Cancels runs that haven't been started yet, with the given request ID
func (l *listener) handleCancelOracleRequest(request *oracle_wrapper.OracleCancelOracleRequest) {
	runCloserChannelIf, loaded := l.runs.LoadAndDelete(formatRequestId(request.RequestId))
	if loaded {
		close(runCloserChannelIf.(chan struct{}))
	}
}

// JobID complies with log.Listener
func (*listener) JobID() models.JobID {
	return models.NilJobID
}

// Job complies with log.Listener
func (l *listener) JobIDV2() int32 {
	return l.spec.ID
}

// IsV2Job complies with log.Listener
func (*listener) IsV2Job() bool {
	return true
}

func formatRequestId(requestId [32]byte) string {
	return fmt.Sprintf("0x%x", requestId)
}
