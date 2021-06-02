package directrequest

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/smartcontractkit/chainlink/core/services/postgres"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/assets"
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

type (
	Delegate struct {
		logBroadcaster log.Broadcaster
		pipelineRunner pipeline.Runner
		pipelineORM    pipeline.ORM
		db             *gorm.DB
		ethClient      eth.Client
		chHeads        chan models.Head
		config         Config
	}

	Config interface {
		MinIncomingConfirmations() uint32
		MinimumContractPayment() *assets.Link
	}
)

func NewDelegate(
	logBroadcaster log.Broadcaster,
	pipelineRunner pipeline.Runner,
	pipelineORM pipeline.ORM,
	ethClient eth.Client,
	db *gorm.DB,
	config Config,
) *Delegate {
	return &Delegate{
		logBroadcaster,
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
		return nil, errors.Errorf("directrequest.Delegate expects a *job.DirectRequestSpec to be present, got %v", job)
	}
	concreteSpec := job.DirectRequestSpec

	oracle, err := oracle_wrapper.NewOracle(concreteSpec.ContractAddress.Address(), d.ethClient)
	if err != nil {
		return
	}

	minIncomingConfirmations := d.config.MinIncomingConfirmations()

	if concreteSpec.MinIncomingConfirmations.Uint32 > minIncomingConfirmations {
		minIncomingConfirmations = concreteSpec.MinIncomingConfirmations.Uint32
	}

	logListener := &listener{
		config:                   d.config,
		logBroadcaster:           d.logBroadcaster,
		oracle:                   oracle,
		pipelineRunner:           d.pipelineRunner,
		db:                       d.db,
		pipelineORM:              d.pipelineORM,
		job:                      job,
		onChainJobSpecID:         job.DirectRequestSpec.OnChainJobSpecID.Hash(),
		mbLogs:                   utils.NewMailbox(50),
		minIncomingConfirmations: uint64(minIncomingConfirmations),
		chStop:                   make(chan struct{}),
	}
	services = append(services, logListener)

	return
}

var (
	_ log.Listener = &listener{}
	_ job.Service  = &listener{}
)

type listener struct {
	config                   Config
	logBroadcaster           log.Broadcaster
	oracle                   oracle_wrapper.OracleInterface
	pipelineRunner           pipeline.Runner
	db                       *gorm.DB
	pipelineORM              pipeline.ORM
	job                      job.Job
	onChainJobSpecID         common.Hash
	runs                     sync.Map
	shutdownWaitGroup        sync.WaitGroup
	mbLogs                   *utils.Mailbox
	minIncomingConfirmations uint64
	chStop                   chan struct{}
	utils.StartStopOnce
}

// Start complies with job.Service
func (l *listener) Start() error {
	return l.StartOnce("DirectRequestListener", func() error {
		unsubscribeLogs := l.logBroadcaster.Register(l, log.ListenerOpts{
			Contract: l.oracle,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				oracle_wrapper.OracleOracleRequest{}.Topic():       {{log.Topic(l.onChainJobSpecID)}},
				oracle_wrapper.OracleCancelOracleRequest{}.Topic(): {{log.Topic(l.onChainJobSpecID)}},
			},
			NumConfirmations: l.minIncomingConfirmations,
		})
		l.shutdownWaitGroup.Add(2)
		go l.run()

		go func() {
			<-l.chStop
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
		case <-l.mbLogs.Notify():
			l.handleReceivedLogs()
		}
	}
}

func (l *listener) handleReceivedLogs() {
	for {
		i, exists := l.mbLogs.Retrieve()
		if !exists {
			return
		}
		lb, ok := i.(log.Broadcast)
		if !ok {
			panic(errors.Errorf("DirectRequestListener: invariant violation, expected log.Broadcast but got %T", lb))
		}
		ctx, cancel := postgres.DefaultQueryCtx()
		defer cancel()
		was, err := l.logBroadcaster.WasAlreadyConsumed(l.db.WithContext(ctx), lb)
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
			l.handleOracleRequest(log, lb)
		case *oracle_wrapper.OracleCancelOracleRequest:
			l.handleCancelOracleRequest(log, lb)
		default:
			logger.Warnf("unexpected log type %T", log)
		}
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

func (l *listener) handleOracleRequest(request *oracle_wrapper.OracleOracleRequest, lb log.Broadcast) {
	minimumContractPayment := l.config.MinimumContractPayment()
	if minimumContractPayment != nil {
		requestPayment := assets.Link(*request.Payment)
		if minimumContractPayment.Cmp(&requestPayment) > 0 {
			logger.Infow("Rejected run for insufficient payment",
				"minimumContractPayment", minimumContractPayment.String(),
				"requestPayment", requestPayment.String(),
			)
			ctx, cancel := postgres.DefaultQueryCtx()
			defer cancel()
			if err := l.logBroadcaster.MarkConsumed(l.db.WithContext(ctx), lb); err != nil {
				logger.Errorw("DirectRequest: unable to mark log consumed", "err", err, "log", lb.String())
			}
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
		run, trrs, err := l.pipelineRunner.ExecuteRun(ctx, *l.job.PipelineSpec, nil, pipeline.JSONSerializable{Val: meta, Null: false}, *logger)
		if ctx.Err() != nil {
			return
		} else if err != nil {
			logger.Errorw("DirectRequest failed to create run", "err", err)
		}
		ctx, cancel = context.WithTimeout(ctx, postgres.DefaultQueryTimeout)
		defer cancel()
		err = postgres.GormTransaction(ctx, l.db, func(tx *gorm.DB) error {
			_, err = l.pipelineRunner.InsertFinishedRun(tx, run, trrs, true)
			if err != nil {
				return err
			}
			return l.logBroadcaster.MarkConsumed(tx, lb)
		})
		if ctx.Err() != nil {
			return
		} else if err != nil {
			logger.Errorw("DirectRequest failed to create run", "err", err)
		}
	}()
}

// Cancels runs that haven't been started yet, with the given request ID
func (l *listener) handleCancelOracleRequest(request *oracle_wrapper.OracleCancelOracleRequest, lb log.Broadcast) {
	runCloserChannelIf, loaded := l.runs.LoadAndDelete(formatRequestId(request.RequestId))
	if loaded {
		close(runCloserChannelIf.(chan struct{}))
	}
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	if err := l.logBroadcaster.MarkConsumed(l.db.WithContext(ctx), lb); err != nil {
		logger.Errorw("DirectRequest: failed to mark log consumed", "log", lb.String())
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
