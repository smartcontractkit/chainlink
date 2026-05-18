package directrequest

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type (
	Delegate struct {
		logger         logger.Logger
		pipelineRunner pipeline.Runner
		pipelineORM    pipeline.ORM
		chHeads        chan *evmtypes.Head
		legacyChains   legacyevm.LegacyChainContainer
		mailMon        *mailbox.Monitor
	}

	Config interface {
		MinIncomingConfirmations() uint32
		MinContractPayment() *assets.Link
	}
)

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(
	logger logger.Logger,
	pipelineRunner pipeline.Runner,
	pipelineORM pipeline.ORM,
	legacyChains legacyevm.LegacyChainContainer,
	mailMon *mailbox.Monitor,
) *Delegate {
	return &Delegate{
		logger:         logger.Named("DirectRequest"),
		pipelineRunner: pipelineRunner,
		pipelineORM:    pipelineORM,
		chHeads:        make(chan *evmtypes.Head, 1),
		legacyChains:   legacyChains,
		mailMon:        mailMon,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.DirectRequest
}

func (d *Delegate) BeforeJobCreated(spec job.Job)              {}
func (d *Delegate) AfterJobCreated(spec job.Job)               {}
func (d *Delegate) BeforeJobDeleted(spec job.Job)              {}
func (d *Delegate) OnDeleteJob(context.Context, job.Job) error { return nil }

// ServicesForSpec returns the log listener service for a direct request job
func (d *Delegate) ServicesForSpec(ctx context.Context, jb job.Job) ([]job.ServiceCtx, error) {
	if jb.DirectRequestSpec == nil {
		return nil, errors.Errorf("DirectRequest: directrequest.Delegate expects a *job.DirectRequestSpec to be present, got %v", jb)
	}
	chain, err := d.legacyChains.Get(jb.DirectRequestSpec.EVMChainID.String())
	if err != nil {
		return nil, err
	}
	concreteSpec := job.SetDRMinIncomingConfirmations(chain.Config().EVM().MinIncomingConfirmations(), *jb.DirectRequestSpec)

	oracle, err := operator_wrapper.NewOperator(concreteSpec.ContractAddress.Address(), chain.Client())
	if err != nil {
		return nil, errors.Wrapf(err, "DirectRequest: failed to create an operator wrapper for address: %v", concreteSpec.ContractAddress.Address().String())
	}

	svcLogger := d.logger.Named(jb.ExternalJobID.String()).
		With(
			"contract", concreteSpec.ContractAddress.Address().String(),
			"jobName", jb.PipelineSpec.JobName,
			"jobID", jb.PipelineSpec.JobID,
			"externalJobID", jb.ExternalJobID,
		)

	logListener := &listener{
		logger:                   svcLogger.Named("Listener"),
		config:                   chain.Config().EVM(),
		logBroadcaster:           chain.LogBroadcaster(),
		oracle:                   oracle,
		pipelineRunner:           d.pipelineRunner,
		pipelineORM:              d.pipelineORM,
		mailMon:                  d.mailMon,
		job:                      jb,
		mbOracleRequests:         mailbox.NewHighCapacity[log.Broadcast](),
		mbOracleCancelRequests:   mailbox.NewHighCapacity[log.Broadcast](),
		minIncomingConfirmations: concreteSpec.MinIncomingConfirmations.Uint32,
		requesters:               concreteSpec.Requesters,
		minContractPayment:       concreteSpec.MinContractPayment,
		chStop:                   make(chan struct{}),
	}
	var services []job.ServiceCtx
	services = append(services, logListener)

	return services, nil
}

var (
	_ log.Listener   = &listener{}
	_ job.ServiceCtx = &listener{}
)

type listener struct {
	services.StateMachine
	logger                   logger.Logger
	config                   Config
	logBroadcaster           log.Broadcaster
	oracle                   operator_wrapper.OperatorInterface
	pipelineRunner           pipeline.Runner
	pipelineORM              pipeline.ORM
	mailMon                  *mailbox.Monitor
	job                      job.Job
	runs                     sync.Map // map[string]services.StopChan
	shutdownWaitGroup        sync.WaitGroup
	mbOracleRequests         *mailbox.Mailbox[log.Broadcast]
	mbOracleCancelRequests   *mailbox.Mailbox[log.Broadcast]
	minIncomingConfirmations uint32
	requesters               models.AddressCollection
	minContractPayment       *assets.Link
	chStop                   services.StopChan
}

func (l *listener) HealthReport() map[string]error {
	return map[string]error{l.Name(): l.Healthy()}
}

func (l *listener) Name() string { return l.logger.Name() }

// Start complies with job.Service
func (l *listener) Start(context.Context) error {
	return l.StartOnce("DirectRequestListener", func() error {
		unsubscribeLogs := l.logBroadcaster.Register(l, log.ListenerOpts{
			Contract: l.oracle.Address(),
			ParseLog: l.oracle.ParseLog,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				operator_wrapper.OperatorOracleRequest{}.Topic():       {{log.Topic(l.job.ExternalIDEncodeBytesToTopic()), log.Topic(l.job.ExternalIDEncodeStringToTopic())}},
				operator_wrapper.OperatorCancelOracleRequest{}.Topic(): {{log.Topic(l.job.ExternalIDEncodeBytesToTopic()), log.Topic(l.job.ExternalIDEncodeStringToTopic())}},
			},
			MinIncomingConfirmations: l.minIncomingConfirmations,
		})
		l.shutdownWaitGroup.Add(3)
		go l.processOracleRequests()
		go l.processCancelOracleRequests()

		go func() {
			<-l.chStop
			unsubscribeLogs()
			l.shutdownWaitGroup.Done()
		}()

		l.mailMon.Monitor(l.mbOracleRequests, "DirectRequest", "Requests", fmt.Sprint(l.job.PipelineSpec.JobID))
		l.mailMon.Monitor(l.mbOracleCancelRequests, "DirectRequest", "Cancel", fmt.Sprint(l.job.PipelineSpec.JobID))

		return nil
	})
}

// Close complies with job.Service
func (l *listener) Close() error {
	return l.StopOnce("DirectRequestListener", func() error {
		l.runs.Range(func(key, runCloserChannelIf interface{}) bool {
			runCloserChannel := runCloserChannelIf.(services.StopChan)
			close(runCloserChannel)
			return true
		})
		l.runs = sync.Map{}

		close(l.chStop)
		l.shutdownWaitGroup.Wait()

		return services.CloseAll(l.mbOracleRequests, l.mbOracleCancelRequests)
	})
}

func (l *listener) HandleLog(ctx context.Context, lb log.Broadcast) {
	log := lb.DecodedLog()
	if log == nil || reflect.ValueOf(log).IsNil() {
		l.logger.Error("HandleLog: ignoring nil value")
		return
	}

	switch log := log.(type) {
	case *operator_wrapper.OperatorOracleRequest:
		wasOverCapacity := l.mbOracleRequests.Deliver(lb)
		if wasOverCapacity {
			l.logger.Error("OracleRequest log mailbox is over capacity - dropped the oldest log")
		}
	case *operator_wrapper.OperatorCancelOracleRequest:
		wasOverCapacity := l.mbOracleCancelRequests.Deliver(lb)
		if wasOverCapacity {
			l.logger.Error("CancelOracleRequest log mailbox is over capacity - dropped the oldest log")
		}
	default:
		l.logger.Warnf("Unexpected log type %T", log)
	}
}

func (l *listener) processOracleRequests() {
	ctx, cancel := l.chStop.NewCtx()
	defer cancel()
	for {
		select {
		case <-l.chStop:
			l.shutdownWaitGroup.Done()
			return
		case <-l.mbOracleRequests.Notify():
			l.handleReceivedLogs(ctx, l.mbOracleRequests)
		}
	}
}

func (l *listener) processCancelOracleRequests() {
	ctx, cancel := l.chStop.NewCtx()
	defer cancel()
	for {
		select {
		case <-l.chStop:
			l.shutdownWaitGroup.Done()
			return
		case <-l.mbOracleCancelRequests.Notify():
			l.handleReceivedLogs(ctx, l.mbOracleCancelRequests)
		}
	}
}

func (l *listener) handleReceivedLogs(ctx context.Context, mailbox *mailbox.Mailbox[log.Broadcast]) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		lb, exists := mailbox.Retrieve()
		if !exists {
			return
		}
		was, err := l.logBroadcaster.WasAlreadyConsumed(ctx, lb)
		if err != nil {
			l.logger.Errorw("Could not determine if log was already consumed", "err", err)
			continue
		} else if was {
			continue
		}

		logJobSpecID := lb.RawLog().Topics[1]
		if logJobSpecID == (common.Hash{}) || (logJobSpecID != l.job.ExternalIDEncodeStringToTopic() && logJobSpecID != l.job.ExternalIDEncodeBytesToTopic()) {
			l.logger.Debugw("Skipping Run for Log with wrong Job ID", "logJobSpecID", logJobSpecID)
			l.markLogConsumed(ctx, nil, lb)
			continue
		}

		log := lb.DecodedLog()
		if log == nil || reflect.ValueOf(log).IsNil() {
			l.logger.Error("HandleLog: ignoring nil value")
			continue
		}

		switch log := log.(type) {
		case *operator_wrapper.OperatorOracleRequest:
			l.handleOracleRequest(ctx, log, lb)
		case *operator_wrapper.OperatorCancelOracleRequest:
			l.handleCancelOracleRequest(ctx, nil, log, lb)
		default:
			l.logger.Warnf("Unexpected log type %T", log)
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

func (l *listener) handleOracleRequest(ctx context.Context, request *operator_wrapper.OperatorOracleRequest, lb log.Broadcast) {
	l.logger.Infow("Oracle request received",
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
		l.logger.Infow("Rejected run for invalid requester",
			"requester", request.Requester,
			"allowedRequesters", l.requesters.ToStrings(),
		)
		l.markLogConsumed(ctx, nil, lb)
		return
	}

	var minContractPayment *assets.Link
	if l.minContractPayment != nil {
		minContractPayment = l.minContractPayment
	} else {
		minContractPayment = l.config.MinContractPayment()
	}
	if minContractPayment != nil && request.Payment != nil {
		requestPayment := assets.Link(*request.Payment)
		if minContractPayment.Cmp(&requestPayment) > 0 {
			l.logger.Warnw("Rejected run for insufficient payment",
				"minContractPayment", minContractPayment.String(),
				"requestPayment", requestPayment.String(),
			)
			l.markLogConsumed(ctx, nil, lb)
			return
		}
	}

	meta := make(map[string]interface{})
	meta["oracleRequest"] = oracleRequestToMap(request)

	runCloserChannel := make(services.StopChan)
	runCloserChannelIf, loaded := l.runs.LoadOrStore(formatRequestId(request.RequestId), runCloserChannel)
	if loaded {
		runCloserChannel = runCloserChannelIf.(services.StopChan)
	}
	ctx, cancel := runCloserChannel.NewCtx()
	defer cancel()

	evmChainID := lb.EVMChainID()
	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"databaseID":    l.job.ID,
			"externalJobID": l.job.ExternalJobID,
			"name":          l.job.Name.ValueOrZero(),
			"pipelineSpec": &pipeline.Spec{
				ForwardingAllowed: l.job.ForwardingAllowed,
			},
			"evmChainID": evmChainID.String(),
		},
		"jobRun": map[string]interface{}{
			"meta":                  meta,
			"logBlockHash":          request.Raw.BlockHash,
			"logBlockNumber":        request.Raw.BlockNumber,
			"logTxHash":             request.Raw.TxHash,
			"logAddress":            request.Raw.Address,
			"logTopics":             request.Raw.Topics,
			"logData":               request.Raw.Data,
			"blockReceiptsRoot":     lb.ReceiptsRoot(),
			"blockTransactionsRoot": lb.TransactionsRoot(),
			"blockStateRoot":        lb.StateRoot(),
		},
	})
	run := pipeline.NewRun(*l.job.PipelineSpec, vars)
	_, err := l.pipelineRunner.Run(ctx, run, true, func(tx sqlutil.DataSource) error {
		l.markLogConsumed(ctx, tx, lb)
		return nil
	})
	if ctx.Err() != nil {
		return
	} else if err != nil {
		l.logger.Errorw("Failed executing run", "err", err)
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
func (l *listener) handleCancelOracleRequest(ctx context.Context, ds sqlutil.DataSource, request *operator_wrapper.OperatorCancelOracleRequest, lb log.Broadcast) {
	runCloserChannelIf, loaded := l.runs.LoadAndDelete(formatRequestId(request.RequestId))
	if loaded {
		close(runCloserChannelIf.(services.StopChan))
	}
	l.markLogConsumed(ctx, ds, lb)
}

func (l *listener) markLogConsumed(ctx context.Context, ds sqlutil.DataSource, lb log.Broadcast) {
	if err := l.logBroadcaster.MarkConsumed(ctx, ds, lb); err != nil {
		l.logger.Errorw("Unable to mark log consumed", "err", err, "log", lb.String())
	}
}

// JobID - Job complies with log.Listener
func (l *listener) JobID() int32 {
	return l.job.ID
}

func formatRequestId(requestId [32]byte) string {
	return fmt.Sprintf("0x%x", requestId)
}
