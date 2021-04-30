package directrequest

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/oracle_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

type (
	Delegate struct {
		logBroadcaster  log.Broadcaster
		headBroadcaster *services.HeadBroadcaster
		pipelineRunner  pipeline.Runner
		pipelineORM     pipeline.ORM
		db              *gorm.DB
		ethClient       eth.Client
		chHeads         chan models.Head
		config          Config
	}

	Config interface {
		MinRequiredOutgoingConfirmations() uint64
		MinimumContractPayment() *assets.Link
	}
)

func NewDelegate(logBroadcaster log.Broadcaster, headBroadcaster *services.HeadBroadcaster,
	pipelineRunner pipeline.Runner, pipelineORM pipeline.ORM,
	ethClient eth.Client, db *gorm.DB, config Config) *Delegate {
	return &Delegate{
		logBroadcaster,
		headBroadcaster,
		pipelineRunner,
		pipelineORM,
		db,
		ethClient,
		make(chan models.Head, 1),
		config,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.DirectRequest
}

// ServicesForSpec returns the log listener service for a direct request job
func (d *Delegate) ServicesForSpec(job job.Job) (services []job.Service, err error) {
	if job.DirectRequestSpec == nil {
		return nil, errors.Errorf("services.Delegate expects a *job.DirectRequestSpec to be present, got %v", job)
	}
	concreteSpec := job.DirectRequestSpec

	oracle, err := oracle_wrapper.NewOracle(concreteSpec.ContractAddress.Address(), d.ethClient)
	if err != nil {
		return
	}

	minConfirmations := d.config.MinRequiredOutgoingConfirmations()

	if concreteSpec.NumConfirmations.Uint32 > uint32(minConfirmations) {
		minConfirmations = uint64(concreteSpec.NumConfirmations.Uint32)
	}

	logListener := &listener{
		config:          d.config,
		logBroadcaster:  d.logBroadcaster,
		headBroadcaster: d.headBroadcaster,
		oracle:          oracle,
		pipelineRunner:  d.pipelineRunner,
		db:              d.db,
		pipelineORM:     d.pipelineORM,
		job:             job,

		// At the moment the mailbox would start skipping if there were
		// too many relevant logs for the same job (> 50) in each block.
		// This is going to get fixed after new LB changes are merged.
		mbLogs:           utils.NewMailbox(50),
		chHeads:          d.chHeads,
		minConfirmations: minConfirmations,
		chStop:           make(chan struct{}),
	}
	copy(logListener.onChainJobSpecID[:], job.DirectRequestSpec.OnChainJobSpecID.Bytes())
	services = append(services, logListener)

	return
}

var (
	_ log.Listener = &listener{}
	_ job.Service  = &listener{}
)

type listener struct {
	config            Config
	logBroadcaster    log.Broadcaster
	headBroadcaster   *services.HeadBroadcaster
	oracle            oracle_wrapper.OracleInterface
	pipelineRunner    pipeline.Runner
	db                *gorm.DB
	pipelineORM       pipeline.ORM
	job               job.Job
	onChainJobSpecID  common.Hash
	runs              sync.Map
	shutdownWaitGroup sync.WaitGroup
	mbLogs            *utils.Mailbox
	chHeads           chan models.Head
	minConfirmations  uint64
	chStop            chan struct{}
	utils.StartStopOnce
}

// Start complies with job.Service
func (l *listener) Start() error {
	return l.StartOnce("DirectRequestListener", func() error {
		unsubscribeLogs := l.logBroadcaster.Register(l, log.ListenerOpts{
			Contract: l.oracle,
			Logs: []generated.AbigenLog{
				oracle_wrapper.OracleOracleRequest{},
				oracle_wrapper.OracleCancelOracleRequest{},
			},
			NumConfirmations: 1,
		})
		l.shutdownWaitGroup.Add(2)
		go l.run()
		unsubscribeHeads := l.headBroadcaster.Subscribe(l)

		go func() {
			<-l.chStop
			unsubscribeHeads()
			unsubscribeLogs()
			l.shutdownWaitGroup.Done()
		}()

		return nil
	})
}

// Close complies with job.Service
func (l *listener) Close() error {
	return l.StopOnce("DirectRequestListener", func() error {
		l.runs.Range(func(key, runCloserChannelIf interface{}) bool {
			runCloserChannel, _ := runCloserChannelIf.(chan struct{})
			close(runCloserChannel)
			return true
		})
		l.runs = sync.Map{}

		close(l.chStop)
		l.shutdownWaitGroup.Wait()

		return nil
	})
}

// OnConnect complies with log.Listener
func (*listener) OnConnect() {}

// OnDisconnect complies with log.Listener
func (*listener) OnDisconnect() {}

func (l *listener) OnNewLongestChain(ctx context.Context, head models.Head) {
	select {
	case l.chHeads <- head:
	default:
	}
}

func (l *listener) HandleLog(lb log.Broadcast) {
	wasOverCapacity := l.mbLogs.Deliver(lb)
	if wasOverCapacity {
		logger.Error("DirectRequestListener: log mailbox is over capacity - dropped the oldest log")
	}
}

func (l *listener) run() {
	for {
		select {
		case <-l.chStop:
			l.shutdownWaitGroup.Done()
			return
		case head := <-l.chHeads:
			l.handleReceivedLogs(head)
		}
	}
}

func (l *listener) handleReceivedLogs(head models.Head) {
	oldEnough := isOldEnoughConstructor(head, l.minConfirmations)
	for {
		i := l.mbLogs.RetrieveIf(oldEnough)
		if i == nil {
			return
		}
		lb, ok := i.(log.Broadcast)
		if !ok {
			panic(errors.Errorf("DirectRequestListener: invariant violation, expected log.Broadcast but got %T", lb))
		}
		was, err := lb.WasAlreadyConsumed()
		if err != nil {
			logger.Errorw("DirectRequestListener: could not determine if log was already consumed", "error", err)
			return
		} else if was {
			return
		}

		logJobSpecID := lb.RawLog().Topics[1]
		if logJobSpecID == (common.Hash{}) || logJobSpecID != l.onChainJobSpecID {
			logger.Debugw("DirectRequestListener: Skipping Run for Log with wrong Job ID", "logJobSpecID", logJobSpecID, "actualJobID", l.onChainJobSpecID)
			return
		}

		log := lb.DecodedLog()
		if log == nil || reflect.ValueOf(log).IsNil() {
			logger.Error("DirectRequestListener: HandleLog: ignoring nil value")
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
}

func isOldEnoughConstructor(head models.Head, minConfirmations uint64) func(interface{}) bool {
	return func(i interface{}) bool {
		broadcast, ok := i.(log.Broadcast)
		if !ok {
			panic(errors.Errorf("DirectRequestListener: Invalid type received - expected Broadcast, got %T", i))
		}
		logHeight := broadcast.RawLog().BlockNumber
		return (logHeight + uint64(minConfirmations) - 1) <= uint64(head.Number)
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
	minimumContractPayment := l.config.MinimumContractPayment()
	if minimumContractPayment != nil {
		requestPayment := assets.Link(*request.Payment)
		if minimumContractPayment.Cmp(&requestPayment) > 0 {
			logger.Infow("Rejected run for insufficient payment",
				"minimumContractPayment", minimumContractPayment.String(),
				"requestPayment", requestPayment.String(),
			)
			return
		}
	}

	meta := make(map[string]interface{})
	meta["oracleRequest"] = oracleRequestToMap(request)

	logger := logger.CreateLogger(logger.Default.With(
		"jobName", l.job.PipelineSpec.JobName,
		"jobID", l.job.PipelineSpec.JobID,
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

		_, _, err := l.pipelineRunner.ExecuteAndInsertFinishedRun(ctx, *l.job.PipelineSpec, pipeline.JSONSerializable{Val: meta, Null: false}, *logger, true)
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
	return l.job.ID
}

// IsV2Job complies with log.Listener
func (*listener) IsV2Job() bool {
	return true
}

func formatRequestId(requestId [32]byte) string {
	return fmt.Sprintf("0x%x", requestId)
}
