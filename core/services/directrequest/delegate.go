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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/operator_wrapper"
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

var _ job.Delegate = (*Delegate)(nil)

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

func (Delegate) AfterJobCreated(spec job.Job)  {}
func (Delegate) BeforeJobDeleted(spec job.Job) {}

// ServicesForSpec returns the log listener service for a direct request job
func (d *Delegate) ServicesForSpec(job job.Job) (services []job.Service, err error) {
	if job.DirectRequestSpec == nil {
		return nil, errors.Errorf("directrequest.Delegate expects a *job.DirectRequestSpec to be present, got %v", job)
	}
	concreteSpec := job.DirectRequestSpec

	oracle, err := operator_wrapper.NewOperator(concreteSpec.ContractAddress.Address(), d.ethClient)
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
		mbOracleRequests:         utils.NewHighCapacityMailbox(),
		mbOracleCancelRequests:   utils.NewHighCapacityMailbox(),
		minIncomingConfirmations: uint64(minIncomingConfirmations),
		requesters:               concreteSpec.Requesters,
		minContractPayment:       concreteSpec.MinContractPayment,
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
	oracle                   operator_wrapper.OperatorInterface
	pipelineRunner           pipeline.Runner
	db                       *gorm.DB
	pipelineORM              pipeline.ORM
	job                      job.Job
	runs                     sync.Map
	shutdownWaitGroup        sync.WaitGroup
	mbOracleRequests         *utils.Mailbox
	mbOracleCancelRequests   *utils.Mailbox
	minIncomingConfirmations uint64
	requesters               models.AddressCollection
	minContractPayment       *assets.Link
	chStop                   chan struct{}
	utils.StartStopOnce
}

// Start complies with job.Service
func (l *listener) Start() error {
	return l.StartOnce("DirectRequestListener", func() error {
		unsubscribeLogs := l.logBroadcaster.Register(l, log.ListenerOpts{
			Contract: l.oracle.Address(),
			ParseLog: l.oracle.ParseLog,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				operator_wrapper.OperatorOracleRequest{}.Topic():       {{log.Topic(l.job.ExternalIDEncodeBytesToTopic()), log.Topic(l.job.ExternalIDEncodeStringToTopic())}},
				operator_wrapper.OperatorCancelOracleRequest{}.Topic(): {{log.Topic(l.job.ExternalIDEncodeBytesToTopic()), log.Topic(l.job.ExternalIDEncodeStringToTopic())}},
			},
			NumConfirmations: l.minIncomingConfirmations,
		})
		l.shutdownWaitGroup.Add(3)
		go l.processOracleRequests()
		go l.processCancelOracleRequests()

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
	log := lb.DecodedLog()
	if log == nil || reflect.ValueOf(log).IsNil() {
		logger.Error("DirectRequest: HandleLog: ignoring nil value")
		return
	}

	switch log := log.(type) {
	case *operator_wrapper.OperatorOracleRequest:
		wasOverCapacity := l.mbOracleRequests.Deliver(lb)
		if wasOverCapacity {
			logger.Error("DirectRequestListener: OracleRequest log mailbox is over capacity - dropped the oldest log")
		}
	case *operator_wrapper.OperatorCancelOracleRequest:
		wasOverCapacity := l.mbOracleCancelRequests.Deliver(lb)
		if wasOverCapacity {
			logger.Error("DirectRequest: CancelOracleRequest log mailbox is over capacity - dropped the oldest log")
		}
	default:
		logger.Warnf("DirectRequest: unexpected log type %T", log)
	}
}

func (l *listener) processOracleRequests() {
	for {
		select {
		case <-l.chStop:
			l.shutdownWaitGroup.Done()
			return
		case <-l.mbOracleRequests.Notify():
			l.handleReceivedLogs(l.mbOracleRequests)
		}
	}
}

func (l *listener) processCancelOracleRequests() {
	for {
		select {
		case <-l.chStop:
			l.shutdownWaitGroup.Done()
			return
		case <-l.mbOracleCancelRequests.Notify():
			l.handleReceivedLogs(l.mbOracleCancelRequests)
		}
	}
}

func (l *listener) handleReceivedLogs(mailbox *utils.Mailbox) {
	for {
		i, exists := mailbox.Retrieve()
		if !exists {
			return
		}
		lb, ok := i.(log.Broadcast)
		if !ok {
			panic(errors.Errorf("DirectRequestListener: invariant violation, expected log.Broadcast but got %T", lb))
		}
		ctx, cancel := postgres.DefaultQueryCtx()
		was, err := l.logBroadcaster.WasAlreadyConsumed(l.db.WithContext(ctx), lb)
		cancel()
		if err != nil {
			logger.Errorw("DirectRequestListener: could not determine if log was already consumed", "error", err)
			return
		} else if was {
			return
		}

		logJobSpecID := lb.RawLog().Topics[1]
		if logJobSpecID == (common.Hash{}) || (logJobSpecID != l.job.ExternalIDEncodeStringToTopic() && logJobSpecID != l.job.ExternalIDEncodeBytesToTopic()) {
			logger.Debugw("DirectRequestListener: Skipping Run for Log with wrong Job ID", "logJobSpecID", logJobSpecID, "actualJobID", l.job.ExternalJobID)
			l.markLogConsumed(nil, lb)
			return
		}

		log := lb.DecodedLog()
		if log == nil || reflect.ValueOf(log).IsNil() {
			logger.Error("DirectRequestListener: HandleLog: ignoring nil value")
			return
		}

		switch log := log.(type) {
		case *operator_wrapper.OperatorOracleRequest:
			l.handleOracleRequest(log, lb)
		case *operator_wrapper.OperatorCancelOracleRequest:
			l.handleCancelOracleRequest(log, lb)
		default:
			logger.Warnf("unexpected log type %T", log)
		}
	}
}

func oracleRequestToMap(request *operator_wrapper.OperatorOracleRequest) map[string]interface{} {
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

func (l *listener) handleOracleRequest(request *operator_wrapper.OperatorOracleRequest, lb log.Broadcast) {
	logger.Infow("Direct request: oracle request received",
		"specId", fmt.Sprintf("%0x", request.SpecId),
		"requester", request.Requester,
		"requestId", fmt.Sprintf("%0x", request.RequestId),
		"payment", request.Payment,
		"callbackAddr", request.CallbackAddr,
		"callbackFunctionId", fmt.Sprintf("%0x", request.CallbackFunctionId),
		"cancelExpiration", request.CancelExpiration,
		"dataVersion", request.DataVersion,
		"data", fmt.Sprintf("%0x", request.Data),
	)

	if !l.allowRequester(request.Requester) {
		logger.Infow("Rejected run for invalid requester",
			"requester", request.Requester,
			"allowedRequesters", l.requesters.ToStrings(),
		)
		l.markLogConsumed(nil, lb)
		return
	}

	var minContractPayment *assets.Link
	if l.minContractPayment != nil {
		minContractPayment = l.minContractPayment
	} else {
		minContractPayment = l.config.MinimumContractPayment()
	}
	if minContractPayment != nil {
		requestPayment := assets.Link(*request.Payment)
		if minContractPayment.Cmp(&requestPayment) > 0 {
			logger.Infow("Rejected run for insufficient payment",
				"minContractPayment", minContractPayment.String(),
				"requestPayment", requestPayment.String(),
			)
			l.markLogConsumed(nil, lb)
			return
		}
	}

	meta := make(map[string]interface{})
	meta["oracleRequest"] = oracleRequestToMap(request)

	runCloserChannel := make(chan struct{})
	runCloserChannelIf, loaded := l.runs.LoadOrStore(formatRequestId(request.RequestId), runCloserChannel)
	if loaded {
		runCloserChannel, _ = runCloserChannelIf.(chan struct{})
	}
	ctx, cancel := utils.CombinedContext(runCloserChannel, context.Background())
	defer cancel()

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"databaseID":    l.job.ID,
			"externalJobID": l.job.ExternalJobID,
			"name":          l.job.Name.ValueOrZero(),
		},
		"jobRun": map[string]interface{}{
			"meta":           meta,
			"logBlockHash":   request.Raw.BlockHash,
			"logBlockNumber": request.Raw.BlockNumber,
			"logTxHash":      request.Raw.TxHash,
			"logAddress":     request.Raw.Address,
			"logTopics":      request.Raw.Topics,
			"logData":        request.Raw.Data,
		},
	})
	run := pipeline.NewRun(*l.job.PipelineSpec, vars)
	_, err := l.pipelineRunner.Run(ctx, &run, *logger.Default, true)
	l.markLogConsumed(nil, lb)
	if ctx.Err() != nil {
		return
	} else if err != nil {
		logger.Errorw("DirectRequest: failed executing run", "err", err)
	}
}

func (l *listener) allowRequester(requester common.Address) bool {
	if len(l.requesters) == 0 {
		return true
	}
	for _, addr := range l.requesters {
		if addr == requester {
			return true
		}
	}
	return false
}

// Cancels runs that haven't been started yet, with the given request ID
func (l *listener) handleCancelOracleRequest(request *operator_wrapper.OperatorCancelOracleRequest, lb log.Broadcast) {
	runCloserChannelIf, loaded := l.runs.LoadAndDelete(formatRequestId(request.RequestId))
	if loaded {
		close(runCloserChannelIf.(chan struct{}))
	}
	l.markLogConsumed(nil, lb)
}

func (l *listener) markLogConsumed(db *gorm.DB, lb log.Broadcast) {
	if db == nil {
		ctx, cancel := postgres.DefaultQueryCtx()
		defer cancel()
		db = l.db.WithContext(ctx)
	}
	if err := l.logBroadcaster.MarkConsumed(db, lb); err != nil {
		logger.Errorw("DirectRequest: unable to mark log consumed", "err", err, "log", lb.String())
	}
}

// JobID complies with log.Listener
func (*listener) JobID() models.JobID {
	return models.NilJobID
}

// JobID - Job complies with log.Listener
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
