package fluxmonitor

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services/postgres"

	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flags_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitor/promfm"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const hibernationPollPeriod = 24 * time.Hour

var fluxAggregatorABI = eth.MustGetABI(flux_aggregator_wrapper.FluxAggregatorABI)

//go:generate mockery --name Service --output ../../internal/mocks/ --case=underscore
//go:generate mockery --name DeviationCheckerFactory --output ../../internal/mocks/ --case=underscore
//go:generate mockery --name DeviationChecker --output ../../internal/mocks/ --case=underscore

type RunManager interface {
	Create(
		jobSpecID models.JobID,
		initiator *models.Initiator,
		creationHeight *big.Int,
		runRequest *models.RunRequest,
	) (*models.JobRun, error)
}

// Service is the interface encapsulating all functionality
// needed to listen to price deviations and new round requests.
type Service interface {
	AddJob(models.JobSpec) error
	RemoveJob(models.JobID)
	service.Service
	SetLogger(logger *logger.Logger)
}

type concreteFluxMonitor struct {
	muLogger       sync.RWMutex
	store          *store.Store
	runManager     RunManager
	logBroadcaster log.Broadcaster
	log            *logger.Logger
	checkerFactory DeviationCheckerFactory
	chAdd          chan addEntry
	chRemove       chan models.JobID
	chConnect      chan *models.Head
	chDisconnect   chan struct{}
	chStop         chan struct{}
	chDone         chan struct{}
	started        bool
}

type addEntry struct {
	jobID    models.JobID
	checkers []DeviationChecker
}

// New creates a service that manages a collection of DeviationCheckers,
// one per initiator of type InitiatorFluxMonitor for added jobs.
func New(
	store *store.Store,
	runManager RunManager,
	logBroadcaster log.Broadcaster,
) Service {
	return &concreteFluxMonitor{
		store:          store,
		runManager:     runManager,
		logBroadcaster: logBroadcaster,
		checkerFactory: pollingDeviationCheckerFactory{
			store:          store,
			logBroadcaster: logBroadcaster,
		},
		chAdd:        make(chan addEntry),
		chRemove:     make(chan models.JobID),
		chConnect:    make(chan *models.Head),
		chDisconnect: make(chan struct{}),
		chStop:       make(chan struct{}),
		chDone:       make(chan struct{}),
	}
}

// SetLogger sets and reconfigures the log for the flux monitor service
func (fm *concreteFluxMonitor) SetLogger(logger *logger.Logger) {
	fm.muLogger.Lock()
	defer fm.muLogger.Unlock()
	fm.log = logger
}

func (fm *concreteFluxMonitor) logger() *logger.Logger {
	fm.muLogger.RLock()
	defer fm.muLogger.RUnlock()
	return fm.log
}

func (fm *concreteFluxMonitor) Start() error {
	go fm.serveInternalRequests()
	fm.started = true

	var wg sync.WaitGroup
	err := fm.store.Jobs(func(j *models.JobSpec) bool {
		if j == nil {
			err := errors.New("received nil job")
			fm.logger().Error(err)
			return true
		}
		job := *j

		wg.Add(1)
		go func() {
			defer wg.Done()

			err := fm.AddJob(job)
			if err != nil {
				logger.Errorf("error adding FluxMonitor job: %v", err)
			}
		}()
		return true
	}, models.InitiatorFluxMonitor)

	wg.Wait()

	return err
}

func (fm *concreteFluxMonitor) Ready() error {
	if fm.started {
		return nil
	}
	return utils.ErrNotStarted
}

func (fm *concreteFluxMonitor) Healthy() error {
	if fm.started {
		return nil
	}
	return utils.ErrNotStarted
}

// Disconnect cleans up running deviation checkers.
func (fm *concreteFluxMonitor) Close() error {
	close(fm.chStop)
	if fm.started {
		fm.started = false
		<-fm.chDone
	}
	return nil
}

// serveInternalRequests handles internal requests for state change via
// channels.  Inspired by the ideas of Communicating Sequential Processes, or
// CSP.
func (fm *concreteFluxMonitor) serveInternalRequests() {
	defer close(fm.chDone)

	jobMap := map[models.JobID][]DeviationChecker{}

	for {
		select {
		case entry := <-fm.chAdd:
			if _, ok := jobMap[entry.jobID]; ok {
				logger.Errorf("job '%s' has already been added to flux monitor", entry.jobID.String())
				continue
			}
			for _, checker := range entry.checkers {
				checker.Start()
			}
			jobMap[entry.jobID] = entry.checkers

		case jobID := <-fm.chRemove:
			checkers, ok := jobMap[jobID]
			if !ok {
				logger.Debugf("job '%s' is missing from the flux monitor", jobID.String())
				continue
			}
			waiter := sync.WaitGroup{}
			waiter.Add(len(checkers))
			for _, checker := range checkers {
				go func(checker DeviationChecker) {
					checker.Stop()
					waiter.Done()
				}(checker)
			}
			waiter.Wait()
			delete(jobMap, jobID)

		case <-fm.chStop:
			waiter := sync.WaitGroup{}
			for _, checkers := range jobMap {
				waiter.Add(len(checkers))
				for _, checker := range checkers {
					go func(checker DeviationChecker) {
						checker.Stop()
						waiter.Done()
					}(checker)
				}
			}
			waiter.Wait()
			return
		}
	}
}

// AddJob created a DeviationChecker for any job initiators of type
// InitiatorFluxMonitor.
func (fm *concreteFluxMonitor) AddJob(job models.JobSpec) error {
	if job.ID.IsZero() {
		err := errors.New("received job with zero ID")
		logger.Error(err)
		return err
	}

	var validCheckers []DeviationChecker
	for _, initr := range job.InitiatorsFor(models.InitiatorFluxMonitor) {
		logger.Debugw("Adding job to flux monitor",
			"job", job.ID.String(),
			"initr", initr.ID,
		)

		timeout := fm.store.Config.DefaultHTTPTimeout()
		checker, err := fm.checkerFactory.New(
			initr,
			job.MinPayment,
			fm.runManager,
			fm.store.ORM,
			timeout,
		)
		if err != nil {
			fm.store.UpsertErrorFor(job.ID, "Unable to create deviation checker")
			return errors.Wrap(err, "factory unable to create checker")
		}
		validCheckers = append(validCheckers, checker)
	}
	if len(validCheckers) == 0 {
		return nil
	}

	fm.chAdd <- addEntry{job.ID, validCheckers}
	return nil
}

// RemoveJob stops and removes the checker for all Flux Monitor initiators belonging
// to the passed job ID.
func (fm *concreteFluxMonitor) RemoveJob(id models.JobID) {
	if id.IsZero() {
		logger.Warn("nil job ID passed to FluxMonitor#RemoveJob")
		return
	}
	fm.chRemove <- id
}

// DeviationCheckerFactory holds the New method needed to create a new instance
// of a DeviationChecker.
type DeviationCheckerFactory interface {
	New(models.Initiator, *assets.Link, RunManager, *orm.ORM, models.Duration) (DeviationChecker, error)
}

type pollingDeviationCheckerFactory struct {
	store          *store.Store
	logBroadcaster log.Broadcaster
}

func (f pollingDeviationCheckerFactory) New(
	initr models.Initiator,
	minJobPayment *assets.Link,
	runManager RunManager,
	orm *orm.ORM,
	timeout models.Duration,
) (DeviationChecker, error) {
	minimumPollingInterval := models.Duration(f.store.Config.DefaultHTTPTimeout())

	if !initr.PollTimer.Disabled &&
		initr.PollTimer.Period.Shorter(minimumPollingInterval) {
		return nil, fmt.Errorf("pollTimer.period must be equal or greater than %s", minimumPollingInterval)
	}

	urls, err := ExtractFeedURLs(initr.Feeds, orm)
	if err != nil {
		return nil, err
	}

	requestData, err := initr.RequestData.AsMap()
	if err != nil {
		return nil, err
	}

	fetcher, err := newMedianFetcherFromURLs(
		timeout,
		requestData,
		urls,
		f.store.Config.DefaultHTTPLimit())
	if err != nil {
		return nil, err
	}

	fluxAggregator, err := flux_aggregator_wrapper.NewFluxAggregator(initr.Address, f.store.EthClient)
	if err != nil {
		return nil, err
	}

	var flagsContract *flags_wrapper.Flags
	if f.store.Config.FlagsContractAddress() != "" {
		flagsContractAddress := common.HexToAddress(f.store.Config.FlagsContractAddress())
		flagsContract, err = flags_wrapper.NewFlags(flagsContractAddress, f.store.EthClient)
		errorMsg := fmt.Sprintf("unable to create Flags contract instance, check address: %s", f.store.Config.FlagsContractAddress())
		logger.ErrorIf(err, errorMsg)
	}

	min, err := fluxAggregator.MinSubmissionValue(nil)
	if err != nil {
		return nil, err
	}

	max, err := fluxAggregator.MaxSubmissionValue(nil)
	if err != nil {
		return nil, err
	}

	return NewPollingDeviationChecker(
		f.store,
		fluxAggregator,
		flagsContract,
		f.logBroadcaster,
		initr,
		minJobPayment,
		runManager,
		fetcher,
		min,
		max,
	)
}

// ExtractFeedURLs extracts a list of url.URLs from the feeds parameter of the initiator params
func ExtractFeedURLs(feeds models.Feeds, orm *orm.ORM) ([]*url.URL, error) {
	var feedsData []interface{}
	var urls []*url.URL

	err := json.Unmarshal(feeds.Bytes(), &feedsData)
	if err != nil {
		return nil, err
	}

	for _, entry := range feedsData {
		var bridgeURL *url.URL
		var err error

		switch feed := entry.(type) {
		case string: // feed url - ex: "http://example.com"
			bridgeURL, err = url.ParseRequestURI(feed)
		case map[string]interface{}: // named feed - ex: {"bridge": "bridgeName"}
			bridgeName, ok := feed["bridge"].(string)
			if !ok {
				return nil, errors.New("failed to convert bridge type into string")
			}
			bridgeURL, err = GetBridgeURLFromName(bridgeName, orm) // XXX: currently an n query
		default:
			err = errors.New("unable to extract feed URLs from json")
		}

		if err != nil {
			return nil, err
		}
		urls = append(urls, bridgeURL)
	}

	return urls, nil
}

// GetBridgeURLFromName looks up a bridge in the DB by name, then extracts the url
func GetBridgeURLFromName(name string, orm *orm.ORM) (*url.URL, error) {
	task := models.TaskType(name)
	bridge, err := orm.FindBridge(task)
	if err != nil {
		return nil, err
	}
	bridgeURL := url.URL(bridge.URL)
	return &bridgeURL, nil
}

// DeviationChecker encapsulate methods needed to initialize and check prices
// for price deviations.
type DeviationChecker interface {
	Start()
	Stop()
}

// PollingDeviationChecker polls external price adapters via HTTP to check for price swings.
type PollingDeviationChecker struct {
	store          *store.Store
	fluxAggregator flux_aggregator_wrapper.FluxAggregatorInterface
	flags          flags_wrapper.FlagsInterface
	runManager     RunManager
	logBroadcaster log.Broadcaster
	fetcher        Fetcher
	oracleAddress  common.Address

	initr         models.Initiator
	minJobPayment *assets.Link
	requestData   models.JSON
	precision     int32

	isHibernating    bool
	backlog          *utils.BoundedPriorityQueue
	chProcessLogs    chan struct{}
	pollTicker       utils.PausableTicker
	hibernationTimer utils.ResettableTimer
	idleTimer        utils.ResettableTimer
	roundTimer       utils.ResettableTimer

	minSubmission, maxSubmission *big.Int

	chStop     chan struct{}
	waitOnStop chan struct{}
}

// NewPollingDeviationChecker returns a new instance of PollingDeviationChecker.
func NewPollingDeviationChecker(
	store *store.Store,
	fluxAggregator flux_aggregator_wrapper.FluxAggregatorInterface,
	flags flags_wrapper.FlagsInterface,
	logBroadcaster log.Broadcaster,
	initr models.Initiator,
	minJobPayment *assets.Link,
	runManager RunManager,
	fetcher Fetcher,
	minSubmission, maxSubmission *big.Int,
) (*PollingDeviationChecker, error) {
	var idleTimer = utils.NewResettableTimer()
	if !initr.IdleTimer.Disabled {
		idleTimer.Reset(initr.IdleTimer.Duration.Duration())
	}
	pdc := &PollingDeviationChecker{
		store:            store,
		logBroadcaster:   logBroadcaster,
		fluxAggregator:   fluxAggregator,
		initr:            initr,
		minJobPayment:    minJobPayment,
		requestData:      initr.RequestData,
		precision:        initr.Precision,
		runManager:       runManager,
		fetcher:          fetcher,
		pollTicker:       utils.NewPausableTicker(initr.PollTimer.Period.Duration()),
		hibernationTimer: utils.NewResettableTimer(),
		idleTimer:        idleTimer,
		roundTimer:       utils.NewResettableTimer(),
		minSubmission:    minSubmission,
		maxSubmission:    maxSubmission,
		isHibernating:    false,
		backlog: utils.NewBoundedPriorityQueue(map[uint]uint{
			// We want reconnecting nodes to be able to submit to a round
			// that hasn't hit maxAnswers yet, as well as the newest round.
			PriorityNewRoundLog:      2,
			PriorityAnswerUpdatedLog: 1,
			PriorityFlagChangedLog:   2,
		}),
		chProcessLogs: make(chan struct{}, 1),
		chStop:        make(chan struct{}),
		waitOnStop:    make(chan struct{}),
	}
	// This is necessary due to the unfortunate fact that assigning `nil` to an
	// interface variable causes `x == nil` checks to always return false. If we
	// do this here, in the constructor, we can avoid using reflection when we
	// check `p.flags == nil` later in the code.
	if flags != nil && !reflect.ValueOf(flags).IsNil() {
		pdc.flags = flags
	}
	return pdc, nil
}

const (
	PriorityFlagChangedLog   uint = 0
	PriorityNewRoundLog      uint = 1
	PriorityAnswerUpdatedLog uint = 2
)

// Start begins the CSP consumer in a single goroutine to
// poll the price adapters and listen to NewRound events.
func (p *PollingDeviationChecker) Start() {
	logger.Debugw("Starting checker for job",
		"job", p.initr.JobSpecID.String(),
		"initr", p.initr.ID,
	)

	go gracefulpanic.WrapRecover(func() {
		p.consume()
	})

}

func (p *PollingDeviationChecker) setIsHibernatingStatus() {
	if p.flags == nil {
		p.isHibernating = false
		return
	}
	isFlagLowered, err := p.isFlagLowered()
	if err != nil {
		logger.Errorf("unable to set hibernation status: %v", err)
		p.isHibernating = false
	} else {
		p.isHibernating = !isFlagLowered
	}
}

func (p *PollingDeviationChecker) isFlagLowered() (bool, error) {
	if p.flags == nil {
		return true, nil
	}
	flags, err := p.flags.GetFlags(nil, []common.Address{utils.ZeroAddress, p.initr.Address})
	if err != nil {
		return true, err
	}
	return !flags[0] || !flags[1], nil
}

// Stop stops this instance from polling, cleaning up resources.
func (p *PollingDeviationChecker) Stop() {
	p.pollTicker.Destroy()
	p.hibernationTimer.Stop()
	p.idleTimer.Stop()
	p.roundTimer.Stop()
	close(p.chStop)
	<-p.waitOnStop
}

func (p *PollingDeviationChecker) JobID() models.JobID {
	return p.initr.JobSpecID
}
func (p *PollingDeviationChecker) JobIDV2() int32 { return 0 }
func (p *PollingDeviationChecker) IsV2Job() bool  { return false }

func (p *PollingDeviationChecker) HandleLog(broadcast log.Broadcast) {
	log := broadcast.DecodedLog()
	if log == nil || reflect.ValueOf(log).IsNil() {
		logger.Error("HandleLog: ignoring nil value")
		return
	}

	switch log := log.(type) {
	case *flux_aggregator_wrapper.FluxAggregatorNewRound:
		p.backlog.Add(PriorityNewRoundLog, broadcast)

	case *flux_aggregator_wrapper.FluxAggregatorAnswerUpdated:
		p.backlog.Add(PriorityAnswerUpdatedLog, broadcast)

	case *flags_wrapper.FlagsFlagRaised:
		if log.Subject == utils.ZeroAddress || log.Subject == p.initr.Address {
			p.backlog.Add(PriorityFlagChangedLog, broadcast)
		}

	case *flags_wrapper.FlagsFlagLowered:
		if log.Subject == utils.ZeroAddress || log.Subject == p.initr.Address {
			p.backlog.Add(PriorityFlagChangedLog, broadcast)
		}

	default:
		logger.Warnf("unexpected log type %T", log)
		return
	}

	select {
	case p.chProcessLogs <- struct{}{}:
	default:
	}
}

func (p *PollingDeviationChecker) consume() {
	defer close(p.waitOnStop)

	if err := p.SetOracleAddress(); err != nil {
		logger.Warnw("unable to set oracle address, this flux monitor job may not work correctly", "err", err)
	}

	// subscribe to contract logs
	unsubscribe := p.logBroadcaster.Register(p, log.ListenerOpts{
		Contract: p.fluxAggregator,
		Logs: []generated.AbigenLog{
			flux_aggregator_wrapper.FluxAggregatorNewRound{},
			flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{},
		},
		NumConfirmations: 1,
	})
	defer unsubscribe()

	if p.flags != nil {
		unsubscribe := p.logBroadcaster.Register(p, log.ListenerOpts{
			Contract: p.flags,
			Logs: []generated.AbigenLog{
				flags_wrapper.FlagsFlagLowered{},
				flags_wrapper.FlagsFlagRaised{},
			},
			NumConfirmations: 1,
		})
		defer unsubscribe()
	}

	p.setIsHibernatingStatus()
	p.setInitialTickers()
	p.performInitialPoll()

	for {
		select {
		case <-p.chStop:
			return

		case <-p.chProcessLogs:
			p.processLogs()

		case <-p.pollTicker.Ticks():
			logger.Debugw("Poll ticker fired",
				"pollPeriod", p.initr.PollTimer.Period,
				"idleDuration", p.initr.IdleTimer.Duration,
				"contract", p.initr.Address.Hex(),
			)
			p.pollIfEligible(DeviationThresholds{
				Rel: float64(p.initr.Threshold),
				Abs: float64(p.initr.AbsoluteThreshold),
			})

		case <-p.idleTimer.Ticks():
			logger.Debugw("Idle ticker fired",
				"pollPeriod", p.initr.PollTimer.Period,
				"idleDuration", p.initr.IdleTimer.Duration,
				"contract", p.initr.Address.Hex(),
			)
			p.pollIfEligible(DeviationThresholds{Rel: 0, Abs: 0})

		case <-p.roundTimer.Ticks():
			logger.Debugw("Round timeout ticker fired",
				"pollPeriod", p.initr.PollTimer.Period,
				"idleDuration", p.initr.IdleTimer.Duration,
				"contract", p.initr.Address.Hex(),
			)
			p.pollIfEligible(DeviationThresholds{
				Rel: float64(p.initr.Threshold),
				Abs: float64(p.initr.AbsoluteThreshold),
			})

		case <-p.hibernationTimer.Ticks():
			p.pollIfEligible(DeviationThresholds{Rel: 0, Abs: 0})
		}
	}
}

func (p *PollingDeviationChecker) SetOracleAddress() error {
	log := logger.Default.With(
		"jobID", p.initr.JobSpecID.String(),
		"contract", p.initr.Address.Hex(),
	)

	oracleAddrs, err := p.fluxAggregator.GetOracles(nil)
	if err != nil {
		log.Error("failed to get list of oracles from FluxAggregator contract")

		return errors.Wrap(err, "failed to get list of oracles from FluxAggregator contract")
	}
	keys, err := p.store.KeyStore.SendingKeys()
	if err != nil {
		return errors.Wrap(err, "failed to load send keys")
	}
	for _, k := range keys {
		for _, oracleAddr := range oracleAddrs {
			if k.Address.Hex() == oracleAddr.Hex() {
				p.oracleAddress = oracleAddr
				return nil
			}
		}
	}

	log = log.With(
		"keys", keys,
		"oracleAddresses", oracleAddrs,
	)

	if len(keys) > 0 {
		addr := keys[0].Address.Address()

		log.Warnw(
			"None of the node's keys matched any oracle addresses, using first available key. This flux monitor job may not work correctly",
			"address", addr.Hex(),
		)
		p.oracleAddress = addr

		return nil
	}

	log.Error("No keys found. This flux monitor job may not work correctly")
	return errors.New("No keys found")
}

func (p *PollingDeviationChecker) performInitialPoll() {
	if p.shouldPerformInitialPoll() {
		p.pollIfEligible(DeviationThresholds{
			Rel: float64(p.initr.Threshold),
			Abs: float64(p.initr.AbsoluteThreshold),
		})
	}
}

func (p *PollingDeviationChecker) shouldPerformInitialPoll() bool {
	return !(p.initr.PollTimer.Disabled && p.initr.IdleTimer.Disabled || p.isHibernating)
}

// hibernate restarts the PollingDeviationChecker in hibernation mode
func (p *PollingDeviationChecker) hibernate() {
	logger.Infof("entering hibernation mode for contract: %s", p.initr.Address.Hex())
	p.isHibernating = true
	p.resetTickers(flux_aggregator_wrapper.OracleRoundState{})
}

// reactivate restarts the PollingDeviationChecker without hibernation mode
func (p *PollingDeviationChecker) reactivate() {
	logger.Infof("exiting hibernation mode, reactivating contract: %s", p.initr.Address.Hex())
	p.isHibernating = false
	p.setInitialTickers()
	p.pollIfEligible(DeviationThresholds{Rel: 0, Abs: 0})
}

func (p *PollingDeviationChecker) processLogs() {
	for !p.backlog.Empty() {
		maybeBroadcast := p.backlog.Take()
		broadcast, ok := maybeBroadcast.(log.Broadcast)
		if !ok {
			logger.Errorf("Failed to convert backlog into LogBroadcast.  Type is %T", maybeBroadcast)
		}

		// If the log is a duplicate of one we've seen before, ignore it (this
		// happens because of the LogBroadcaster's backfilling behavior).
		ctx, cancel := postgres.DefaultQueryCtx()
		defer cancel()
		consumed, err := p.logBroadcaster.WasAlreadyConsumed(p.store.DB.WithContext(ctx), broadcast)
		if err != nil {
			logger.Errorf("Error determining if log was already consumed: %v", err)
			continue
		} else if consumed {
			logger.Debug("Log was already consumed by Flux Monitor, skipping")
			continue
		}

		ctx, cancel = postgres.DefaultQueryCtx()
		defer cancel()
		db := p.store.DB.WithContext(ctx)
		switch log := broadcast.DecodedLog().(type) {
		case *flux_aggregator_wrapper.FluxAggregatorNewRound:
			p.respondToNewRoundLog(*log)
			err = p.logBroadcaster.MarkConsumed(db, broadcast)
			if err != nil {
				logger.Errorf("Error marking log as consumed: %v", err)
			}

		case *flux_aggregator_wrapper.FluxAggregatorAnswerUpdated:
			p.respondToAnswerUpdatedLog(*log)
			err = p.logBroadcaster.MarkConsumed(db, broadcast)
			if err != nil {
				logger.Errorf("Error marking log as consumed: %v", err)
			}

		case *flags_wrapper.FlagsFlagRaised:
			// check the contract before hibernating, because one flag could be lowered
			// while the other flag remains raised
			var isFlagLowered bool
			isFlagLowered, err = p.isFlagLowered()
			logger.ErrorIf(err, "Error determining if flag is still raised")
			if !isFlagLowered {
				p.hibernate()
			}
			err = p.logBroadcaster.MarkConsumed(db, broadcast)
			logger.ErrorIf(err, "Error marking log as consumed")

		case *flags_wrapper.FlagsFlagLowered:
			if p.isHibernating {
				p.reactivate()
			}

			err = p.logBroadcaster.MarkConsumed(db, broadcast)
			logger.ErrorIf(err, "Error marking log as consumed")

		default:
			logger.Errorf("unknown log %v of type %T", log, log)
		}
	}
}

// The AnswerUpdated log tells us that round has successfully closed with a new
// answer.  We update our view of the oracleRoundState in case this log was
// generated by a chain reorg.
func (p *PollingDeviationChecker) respondToAnswerUpdatedLog(log flux_aggregator_wrapper.FluxAggregatorAnswerUpdated) {
	logger.Debugw("AnswerUpdated log", p.loggerFieldsForAnswerUpdated(log)...)

	roundState, err := p.roundState(0)
	if err != nil {
		logger.Errorw(fmt.Sprintf("could not fetch oracleRoundState: %v", err), p.loggerFieldsForAnswerUpdated(log)...)
		return
	}
	p.resetTickers(roundState)
}

// The NewRound log tells us that an oracle has initiated a new round.  This tells us that we
// need to poll and submit an answer to the contract regardless of the deviation.
func (p *PollingDeviationChecker) respondToNewRoundLog(log flux_aggregator_wrapper.FluxAggregatorNewRound) {
	l := logger.Default.With(
		"round", log.RoundId,
		"startedBy", log.StartedBy.Hex(),
		"startedAt", log.StartedAt.String(),
		"contract", p.fluxAggregator.Address().Hex(),
		"jobID", p.initr.JobSpecID,
	)
	l.Debugw("NewRound log")

	promfm.SetBigInt(promfm.SeenRound.WithLabelValues(p.initr.JobSpecID.String()), log.RoundId)

	//
	// NewRound answer submission logic:
	//   - Any log that reaches this point, regardless of chain reorgs or log backfilling, is one that we have
	//         not seen before.  Therefore, we should consider acting upon it.
	//   - We always take the round ID from the log, rather than the round ID suggested by `.RoundState`.  The
	//         reason is that if two NewRound logs come in in rapid succession, and we submit a tx for the first,
	//         the `.ReportableRoundID` field in the roundState() response for the 2nd log will not reflect the
	//         fact that we've submitted for the first round (assuming it hasn't been mined yet).
	//   - In the event of a reorg that pushes our previous submissions back into the mempool, we can rely on the
	//         TxManager to ensure they end up being mined into blocks, but this may cause them to revert if they
	//         are mined in an order that violates certain conditions in the FluxAggregator (restartDelay, etc.).
	//         Therefore, the cleanest solution at present is to resubmit for the reorged rounds.  The drawback
	//         of this approach is that one or the other submission tx for a given round will revert, costing the
	//         node operator some gas.  The benefit is that those submissions are guaranteed to be made, ensuring
	//         that we have high data availability (and also ensuring that node operators get paid).
	//   - There are a few straightforward cases where we don't want to submit:
	//         - When we're not eligible
	//         - When the aggregator is underfunded
	//         - When we were the initiator of the round (i.e. we've received our own NewRound log)
	//   - There are a few more nuanced cases as well:
	//         - When our node polls at the same time as another node, and both attempt to start a round.  In that
	//               case, it's possible that the other node will start the round, and our node will see the NewRound
	//               log and try to submit again.
	//         - When the poll ticker fires very soon after we've responded to a NewRound log.
	//
	//         To handle these more nuanced cases, we record round IDs and whether we've submitted for those rounds
	//         in the DB.  If we see we've already submitted for a given round, we simply bail out.
	//
	//         However, in the case of a chain reorganization, we might see logs with round IDs that we've already
	//         seen.  As mentioned above, we want to re-respond to these rounds to ensure high data availability.
	//         Therefore, if a log arrives with a round ID that is < the most recent that we submitted to, we delete
	//         all of the round IDs in the DB back to (and including) the incoming round ID.  This essentially
	//         rewinds the system back to a state wherein those reorg'ed rounds never occurred, allowing it to move
	//         forward normally.
	//
	//         There is one small exception: if the reorg is fairly shallow, and only un-starts a single round, we
	//         do not need to resubmit, because the TxManager will ensure that our existing submission gets back
	//         into the chain.  There is a very small risk that one of the nodes in the quorum (namely, whichever
	//         one started the previous round) will have its existing submission mined first, thereby violating
	//         the restartDelay, but as this risk is isolated to a single node, the round will not time out and
	//         go stale.  We consider this acceptable.
	//

	logRoundID := uint32(log.RoundId.Uint64())

	// We always want to reset the idle timer upon receiving a NewRound log, so we do it before any `return` statements.
	p.resetIdleTimer(log.StartedAt.Uint64())

	mostRecentRoundID, err := p.store.MostRecentFluxMonitorRoundID(p.initr.Address)
	if err != nil && err != gorm.ErrRecordNotFound {
		l.Errorw("error fetching Flux Monitor most recent round ID from DB", "err", err)
		return
	}

	if logRoundID < mostRecentRoundID {
		err = p.store.DeleteFluxMonitorRoundsBackThrough(p.initr.Address, logRoundID)
		if err != nil {
			l.Errorw("error deleting reorged Flux Monitor rounds from DB", "err", err)
			return
		}
	}

	roundStats, jobRunStatus, err := p.statsAndStatusForRound(logRoundID)
	if err != nil {
		l.Errorw("error determining round stats / run status for round", "err", err)
		return
	}

	if roundStats.NumSubmissions > 0 {
		// This indicates either that:
		//     - We tried to start a round at the same time as another node, and their transaction was mined first, or
		//     - The chain experienced a shallow reorg that unstarted the current round.
		// If our previous attempt is still pending, return early and don't re-submit
		// If our previous attempt is already over (completed or errored), we should retry
		if !jobRunStatus.Finished() {
			l.Debugw("Ignoring new round request: started round simultaneously with another node")
			return
		}
	}

	// Ignore rounds we started
	if p.oracleAddress == log.StartedBy {
		l.Infow("Ignoring new round request: we started this round")
		return
	}

	// Ignore rounds we're not eligible for, or for which we won't be paid
	roundState, err := p.roundState(logRoundID)
	if err != nil {
		l.Errorw("Ignoring new round request: error fetching eligibility from contract", "err", err)
		return
	}
	p.resetTickers(roundState)
	err = p.checkEligibilityAndAggregatorFunding(roundState)
	if err != nil {
		l.Infow("Ignoring new round request", "err", err)
		return
	}

	l.Infow("Responding to new round request")

	// Best effort to attach metadata.
	var metaDataForBridge map[string]interface{}
	lrd, err := p.fluxAggregator.LatestRoundData(nil)
	if err != nil {
		l.Warnw("Couldn't read latest round data for request meta", "err", err)
	} else {
		metaDataForBridge, err = models.MarshalBridgeMetaData(lrd.Answer, lrd.UpdatedAt)
		if err != nil {
			l.Warnw("Error marshalling roundState for request meta", "err", err)
		}
	}

	ctx, cancel := utils.CombinedContext(p.chStop)
	defer cancel()
	polledAnswer, err := p.fetcher.Fetch(ctx, metaDataForBridge, *logger.CreateLogger(l))
	if err != nil {
		l.Errorw("unable to fetch median price", "err", err)
		return
	}

	if !p.isValidSubmission(logger.Default.SugaredLogger, polledAnswer) {
		return
	}

	var payment assets.Link
	if roundState.PaymentAmount == nil {
		l.Error("roundState.PaymentAmount shouldn't be nil")
	} else {
		payment = assets.Link(*roundState.PaymentAmount)
	}

	err = p.createJobRun(polledAnswer, logRoundID, &payment)
	if err != nil {
		l.Errorw("unable to create job run", "err", err)
		return
	}
}

var (
	ErrNotEligible   = errors.New("not eligible to submit")
	ErrUnderfunded   = errors.New("aggregator is underfunded")
	ErrPaymentTooLow = errors.New("round payment amount < minimum contract payment")
)

func (p *PollingDeviationChecker) checkEligibilityAndAggregatorFunding(roundState flux_aggregator_wrapper.OracleRoundState) error {
	if !roundState.EligibleToSubmit {
		return ErrNotEligible
	} else if !p.sufficientFunds(roundState) {
		return ErrUnderfunded
	} else if !p.sufficientPayment(roundState.PaymentAmount) {
		return ErrPaymentTooLow
	}
	return nil
}

const MinFundedRounds int64 = 3

// sufficientFunds checks if the contract has sufficient funding to pay all the oracles on a
// conract for a minimum number of rounds, based on the payment amount in the contract
func (p *PollingDeviationChecker) sufficientFunds(state flux_aggregator_wrapper.OracleRoundState) bool {
	min := big.NewInt(int64(state.OracleCount))
	min = min.Mul(min, big.NewInt(MinFundedRounds))
	min = min.Mul(min, state.PaymentAmount)
	return state.AvailableFunds.Cmp(min) >= 0
}

// sufficientPayment checks if the available payment is enough to submit an answer. It compares
// the payment amount on chain with the min payment amount listed in the job spec / ENV var.
func (p *PollingDeviationChecker) sufficientPayment(payment *big.Int) bool {
	aboveOrEqMinGlobalPayment := payment.Cmp(p.store.Config.MinimumContractPayment().ToInt()) >= 0
	aboveOrEqMinJobPayment := true
	if p.minJobPayment != nil {
		aboveOrEqMinJobPayment = payment.Cmp(p.minJobPayment.ToInt()) >= 0
	}
	return aboveOrEqMinGlobalPayment && aboveOrEqMinJobPayment
}

// DeviationThresholds carries parameters used by the threshold-trigger logic
type DeviationThresholds struct {
	Rel float64 // Relative change required, i.e. |new-old|/|old| >= Rel
	Abs float64 // Absolute change required, i.e. |new-old| >= Abs
}

func (p *PollingDeviationChecker) pollIfEligible(thresholds DeviationThresholds) {
	l := logger.Default.With(
		"jobID", p.initr.JobSpecID,
		"address", p.initr.InitiatorParams.Address,
		"threshold", thresholds.Rel,
		"absoluteThreshold", thresholds.Abs,
	)

	if !p.logBroadcaster.IsConnected() {
		l.Warnw("FluxMonitor: LogBroadcaster is not connected to Ethereum node, skipping poll")
		return
	}

	//
	// Poll ticker submission logic:
	//   - We avoid saving on-chain state wherever possible.  Therefore, we do not know which round we should be
	//         submitting for when the pollTicker fires.
	//   - We pass 0 into `roundState()`, and the FluxAggregator returns a suggested roundID for us to
	//         submit to, as well as our eligibility to submit to that round.
	//   - If the poll ticker fires very soon after we've responded to a NewRound log, and our tx has not been
	//         mined, we risk double-submitting for a round.  To detect this, we check the DB to see whether
	//         we've responded to this round already, and bail out if so.
	//

	// Ask the FluxAggregator which round we should be submitting to, and what the state of that round is.
	roundState, err := p.roundState(0)
	if err != nil {
		l.Errorw("unable to determine eligibility to submit from FluxAggregator contract", "err", err)
		p.store.UpsertErrorFor(p.JobID(), "Unable to call roundState method on provided contract. Check contract address.")
		return
	}
	p.resetTickers(roundState)
	l = l.With("reportableRound", roundState.RoundId)

	roundStats, jobRunStatus, err := p.statsAndStatusForRound(roundState.RoundId)
	if err != nil {
		l.Errorw("error determining round stats / run status for round", "err", err)
		return
	}

	// If we've already successfully submitted to this round (ie through a NewRound log)
	// and the associated JobRun hasn't errored, skip polling
	if roundStats.NumSubmissions > 0 && !jobRunStatus.Errored() {
		l.Infow("skipping poll: round already answered, tx unconfirmed")
		return
	}

	// Don't submit if we're not eligible, or won't get paid
	err = p.checkEligibilityAndAggregatorFunding(roundState)
	if err != nil {
		l.Infow(fmt.Sprintf("skipping poll: %v", err))
		return
	}

	// Best effort to attach metadata.
	var metaDataForBridge map[string]interface{}
	lrd, err := p.fluxAggregator.LatestRoundData(nil)
	if err != nil {
		l.Warnw("Couldn't read latest round data for request meta", "err", err)
	} else {
		metaDataForBridge, err = models.MarshalBridgeMetaData(lrd.Answer, lrd.UpdatedAt)
		if err != nil {
			l.Warnw("Error marshalling roundState for request meta", "err", err)
		}
	}

	ctx, cancel := utils.CombinedContext(p.chStop)
	defer cancel()
	polledAnswer, err := p.fetcher.Fetch(ctx, metaDataForBridge, *logger.CreateLogger(l))
	if err != nil {
		l.Errorw("can't fetch answer", "err", err)
		p.store.UpsertErrorFor(p.JobID(), "Error polling")
		return
	}

	if !p.isValidSubmission(l, polledAnswer) {
		return
	}

	jobSpecID := p.initr.JobSpecID.String()
	latestAnswer := decimal.NewFromBigInt(roundState.LatestSubmission, -p.precision)

	promfm.SetDecimal(promfm.SeenValue.WithLabelValues(jobSpecID), polledAnswer)
	l = l.With(
		"latestAnswer", latestAnswer,
		"polledAnswer", polledAnswer,
	)
	if roundState.RoundId > 1 && !OutsideDeviation(latestAnswer, polledAnswer, thresholds) {
		l.Debugw("deviation < threshold, not submitting")
		return
	}

	if roundState.RoundId > 1 {
		l.Infow("deviation > threshold, starting new round")
	} else {
		l.Infow("starting first round")
	}

	var payment assets.Link
	if roundState.PaymentAmount == nil {
		l.Error("roundState.PaymentAmount shouldn't be nil")
	} else {
		payment = assets.Link(*roundState.PaymentAmount)
	}

	err = p.createJobRun(polledAnswer, roundState.RoundId, &payment)
	if err != nil {
		l.Errorw("can't create job run", "err", err)
		return
	}

	promfm.SetDecimal(promfm.ReportedValue.WithLabelValues(jobSpecID), polledAnswer)
	promfm.SetUint32(promfm.ReportedRound.WithLabelValues(jobSpecID), roundState.RoundId)
}

// If the polledAnswer is outside the allowable range, log an error and don't submit.
// to avoid an onchain reversion.
func (p *PollingDeviationChecker) isValidSubmission(l *zap.SugaredLogger, polledAnswer decimal.Decimal) bool {
	max := decimal.NewFromBigInt(p.maxSubmission, -p.precision)
	min := decimal.NewFromBigInt(p.minSubmission, -p.precision)

	if polledAnswer.GreaterThan(max) || polledAnswer.LessThan(min) {
		l.Errorw("polled value is outside acceptable range", "min", min, "max", max, "polled value", polledAnswer)
		p.store.UpsertErrorFor(p.JobID(), "Polled value is outside acceptable range")
		return false
	}
	return true
}

func (p *PollingDeviationChecker) roundState(roundID uint32) (flux_aggregator_wrapper.OracleRoundState, error) {
	return p.fluxAggregator.OracleRoundState(nil, p.oracleAddress, roundID)
}

// initialRoundState fetches the round information that the fluxmonitor should use when starting
// new jobs. Choosing the correct round on startup is key to setting timers correctly.
func (p *PollingDeviationChecker) initialRoundState() flux_aggregator_wrapper.OracleRoundState {
	defaultRoundState := flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
	}
	latestRoundData, err := p.fluxAggregator.LatestRoundData(nil)
	if err != nil {
		logger.Warnf(
			"unable to retrieve latestRoundData for FluxAggregator contract %s - defaulting "+
				"to current time for tickers: %v",
			p.initr.Address.Hex(),
			err,
		)
		return defaultRoundState
	}
	roundID := uint32(latestRoundData.RoundId.Uint64())
	latestRoundState, err := p.fluxAggregator.OracleRoundState(nil, p.oracleAddress, roundID)
	if err != nil {
		logger.Warnf(
			"unable to call roundState for latest round, contract: %s, round: %d, err: %v",
			p.initr.Address.Hex(),
			latestRoundData.RoundId,
			err,
		)
		return defaultRoundState
	}
	return latestRoundState
}

func (p *PollingDeviationChecker) resetTickers(roundState flux_aggregator_wrapper.OracleRoundState) {
	p.resetPollTicker()
	p.resetHibernationTimer()
	p.resetIdleTimer(roundState.StartedAt)
	p.resetRoundTimer(roundStateTimesOutAt(roundState))
}

func (p *PollingDeviationChecker) setInitialTickers() {
	p.resetTickers(p.initialRoundState())
}

func (p *PollingDeviationChecker) resetPollTicker() {
	if !p.initr.PollTimer.Disabled && !p.isHibernating {
		p.pollTicker.Resume()
	} else {
		p.pollTicker.Pause()
	}
}

func (p *PollingDeviationChecker) resetHibernationTimer() {
	if !p.isHibernating {
		p.hibernationTimer.Stop()
	} else {
		p.hibernationTimer.Reset(hibernationPollPeriod)
	}
}

func (p *PollingDeviationChecker) resetRoundTimer(roundTimesOutAt uint64) {
	if p.isHibernating {
		p.roundTimer.Stop()
		return
	}

	loggerFields := p.loggerFields("timesOutAt", roundTimesOutAt)

	if roundTimesOutAt == 0 {
		p.roundTimer.Stop()
		logger.Debugw("disabling roundTimer, no active round", loggerFields...)

	} else {
		timesOutAt := time.Unix(int64(roundTimesOutAt), 0)
		timeUntilTimeout := time.Until(timesOutAt)

		if timeUntilTimeout <= 0 {
			p.roundTimer.Stop()
			logger.Debugw("roundTimer has run down; disabling", loggerFields...)
		} else {
			p.roundTimer.Reset(timeUntilTimeout)
			loggerFields = append(loggerFields, "value", roundTimesOutAt)
			logger.Debugw("updating roundState.TimesOutAt", loggerFields...)
		}
	}
}

func (p *PollingDeviationChecker) resetIdleTimer(roundStartedAtUTC uint64) {
	if p.isHibernating || p.initr.IdleTimer.Disabled {
		p.idleTimer.Stop()
		return
	} else if roundStartedAtUTC == 0 {
		// There is no active round, so keep using the idleTimer we already have
		return
	}

	startedAt := time.Unix(int64(roundStartedAtUTC), 0)
	idleDeadline := startedAt.Add(p.initr.IdleTimer.Duration.Duration())
	timeUntilIdleDeadline := time.Until(idleDeadline)
	loggerFields := p.loggerFields(
		"startedAt", roundStartedAtUTC,
		"timeUntilIdleDeadline", timeUntilIdleDeadline,
	)

	if timeUntilIdleDeadline <= 0 {
		logger.Debugw("not resetting idleTimer, negative duration", loggerFields...)
		return
	}
	p.idleTimer.Reset(timeUntilIdleDeadline)
	logger.Debugw("resetting idleTimer", loggerFields...)
}

// jobRunRequest is the request used to trigger a Job Run by the Flux Monitor.
type jobRunRequest struct {
	Result           decimal.Decimal `json:"result"`
	Address          string          `json:"address"`
	FunctionSelector string          `json:"functionSelector"`
	DataPrefix       string          `json:"dataPrefix"`
}

func (p *PollingDeviationChecker) createJobRun(
	polledAnswer decimal.Decimal,
	roundID uint32,
	paymentAmount *assets.Link,
) error {

	methodID := fluxAggregatorABI.Methods["submit"].ID
	roundIDData := utils.EVMWordUint64(uint64(roundID))

	payload, err := json.Marshal(jobRunRequest{
		Result:           polledAnswer,
		Address:          p.initr.Address.Hex(),
		FunctionSelector: hexutil.Encode(methodID),
		DataPrefix:       hexutil.Encode(roundIDData),
	})
	if err != nil {
		return errors.Wrapf(err, "unable to encode Job Run request in JSON")
	}
	runData, err := models.ParseJSON(payload)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to start chainlink run with payload %s", payload))
	}
	runRequest := models.NewRunRequest(runData)
	runRequest.Payment = paymentAmount

	jobRun, err := p.runManager.Create(p.initr.JobSpecID, &p.initr, nil, runRequest)
	if err != nil {
		return err
	}

	err = p.store.UpdateFluxMonitorRoundStats(p.initr.Address, roundID, jobRun.ID)
	if err != nil {
		logger.Errorw(fmt.Sprintf("error updating FM round submission count: %v", err),
			"address", p.initr.Address.Hex(),
			"roundID", roundID,
			"jobID", p.initr.JobSpecID.String(),
		)
		return err
	}

	return nil
}

func (p *PollingDeviationChecker) loggerFields(added ...interface{}) []interface{} {
	return append(added, []interface{}{
		"pollFrequency", p.initr.PollTimer.Period,
		"idleDuration", p.initr.IdleTimer.Duration,
		"contract", p.initr.Address.Hex(),
		"jobID", p.initr.JobSpecID.String(),
	}...)
}

func (p *PollingDeviationChecker) loggerFieldsForAnswerUpdated(log flux_aggregator_wrapper.FluxAggregatorAnswerUpdated) []interface{} {
	return []interface{}{
		"round", log.RoundId,
		"answer", log.Current.String(),
		"timestamp", log.UpdatedAt.String(),
		"contract", p.fluxAggregator.Address().Hex(),
		"job", p.initr.JobSpecID,
	}
}

// OutsideDeviation checks whether the next price is outside the threshold.
// If both thresholds are zero (default value), always returns true.
func OutsideDeviation(curAnswer, nextAnswer decimal.Decimal, thresholds DeviationThresholds) bool {
	loggerFields := []interface{}{
		"threshold", thresholds.Rel,
		"absoluteThreshold", thresholds.Abs,
		"currentAnswer", curAnswer,
		"nextAnswer", nextAnswer,
	}

	if thresholds.Rel == 0 && thresholds.Abs == 0 {
		logger.Debugw(
			"Deviation thresholds both zero; short-circuiting deviation checker to "+
				"true, regardless of feed values", loggerFields...)
		return true
	}
	diff := curAnswer.Sub(nextAnswer).Abs()
	loggerFields = append(loggerFields, "absoluteDeviation", diff)

	if !diff.GreaterThan(decimal.NewFromFloat(thresholds.Abs)) {
		logger.Debugw("Absolute deviation threshold not met", loggerFields...)
		return false
	}

	if curAnswer.IsZero() {
		if nextAnswer.IsZero() {
			logger.Debugw("Relative deviation is undefined; can't satisfy threshold", loggerFields...)
			return false
		}
		logger.Infow("Threshold met: relative deviation is ∞", loggerFields...)
		return true
	}

	// 100*|new-old|/|old|: Deviation (relative to curAnswer) as a percentage
	percentage := diff.Div(curAnswer.Abs()).Mul(decimal.NewFromInt(100))

	loggerFields = append(loggerFields, "percentage", percentage)

	if percentage.LessThan(decimal.NewFromFloat(thresholds.Rel)) {
		logger.Debugw("Relative deviation threshold not met", loggerFields...)
		return false
	}
	logger.Infow("Relative and absolute deviation thresholds both met", loggerFields...)
	return true
}

// MakeIdleTimer checks the log timestamp and calculates the idle time
// from that.
//
// This function makes the assumption that the local system time is
// relatively accurate (to within a second or so) and all participating nodes
// agree on that.
//
// If system time is not accurate (compared to the cluster) then you should
// expect poor behaviour here.
func MakeIdleTimer(log flux_aggregator_wrapper.FluxAggregatorNewRound, idleThreshold models.Duration, clock utils.AfterNower) <-chan time.Time {
	timeNow := clock.Now()
	if log.StartedAt == nil {
		return defaultIdleTimer(idleThreshold, clock)
	}
	if !log.StartedAt.IsInt64() {
		logger.Errorf("Value for log.StartedAt %s would overflow int64, using default idle timer instead.", log.StartedAt.String())
		return defaultIdleTimer(idleThreshold, clock)
	}
	roundStarted := time.Unix(log.StartedAt.Int64(), 0)
	if roundStarted.After(timeNow) {
		logger.Warnf("Round started time of %s is later than current system time of %s, setting idle timer to %s from now. Most likely scenario is that this machine's clock is running slow. This is suboptimal! Please ensure your system clock is accurate.", roundStarted.String(), timeNow.String(), idleThreshold.Duration().String())
		return defaultIdleTimer(idleThreshold, clock)
	}
	// duration from now until idle threshold = log timestamp + idle threshold - current time
	durationUntilIdleThreshold := roundStarted.Add(idleThreshold.Duration()).Sub(timeNow)
	if durationUntilIdleThreshold < 0 {
		logger.Warnf("Idle threshold already passed, current time is %s and idle timer expired at %s (round started at %s with idle threshold of %s). It's possible you are processing an old round, or this machine has a fast clock. If this keeps happening, check your system clock and make sure it is accurate.", timeNow, roundStarted.Add(idleThreshold.Duration()).String(), roundStarted.String(), idleThreshold.Duration().String())
	}
	return clock.After(durationUntilIdleThreshold)
}

func defaultIdleTimer(idleThreshold models.Duration, clock utils.AfterNower) <-chan time.Time {
	return clock.After(idleThreshold.Duration())
}

func (p *PollingDeviationChecker) statsAndStatusForRound(roundID uint32) (
	models.FluxMonitorRoundStats,
	models.RunStatus,
	error,
) {
	roundStats, err := p.store.FindOrCreateFluxMonitorRoundStats(p.initr.Address, roundID)
	if err != nil {
		return models.FluxMonitorRoundStats{}, "", err
	}
	// JobRun will not exist if this is the first time responding to this round
	var jobRun models.JobRun
	if roundStats.JobRunID.Valid {
		jobRun, err = p.store.FindJobRun(roundStats.JobRunID.UUID)
		if err != nil {
			return models.FluxMonitorRoundStats{}, "", err
		}
	}
	return roundStats, jobRun.Status, nil
}

func roundStateTimesOutAt(rs flux_aggregator_wrapper.OracleRoundState) uint64 {
	return rs.StartedAt + rs.Timeout
}
