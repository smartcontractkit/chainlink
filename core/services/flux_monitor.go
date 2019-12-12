package services

import (
	"chainlink/core/eth"
	"chainlink/core/logger"
	"chainlink/core/store"
	"chainlink/core/store/models"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
)

// FluxMonitor is the interface encapsulating all functionality
// needed to listen to price deviations and new round requests.
type FluxMonitor interface {
	store.HeadTrackable // (Dis)Connect methods handle initial boot and intermittent connectivity.
	AddJob(models.JobSpec) error
	RemoveJob(*models.ID)
}

type concreteFluxMonitor struct {
	store          *store.Store
	runManager     RunManager
	jobs           map[string]models.JobSpec
	checkers       map[uint]DeviationChecker
	checkersMutex  *sync.RWMutex
	checkerFactory DeviationCheckerFactory
	ctx            context.Context
	cancel         context.CancelFunc
}

// NewFluxMonitor creates a service that manages a collection of DeviationCheckers,
// one per initiator of type InitiatorFluxMonitor for added jobs.
func NewFluxMonitor(store *store.Store, runManager RunManager) FluxMonitor {
	return &concreteFluxMonitor{
		store:          store,
		runManager:     runManager,
		jobs:           map[string]models.JobSpec{},
		checkers:       map[uint]DeviationChecker{},
		checkersMutex:  &sync.RWMutex{},
		checkerFactory: pollingDeviationCheckerFactory{},
	}
}

// Connect adds all persisted jobs and starts deviation checkers for each
// flux monitor initiator.
func (fm *concreteFluxMonitor) Connect(*models.Head) error {
	fm.checkersMutex.Lock()
	fm.ctx, fm.cancel = context.WithCancel(context.Background())
	fm.checkersMutex.Unlock()

	var merr error
	err := fm.store.Jobs(func(j *models.JobSpec) bool { // improve scoping of sql query
		merr = multierr.Combine(merr, fm.AddJob(*j))
		return true
	})
	return multierr.Append(err, merr)
}

// Disconnect cleans up running deviation checkers.
func (fm *concreteFluxMonitor) Disconnect() {
	fm.checkersMutex.Lock()
	defer fm.checkersMutex.Unlock()

	fm.cancel() // parent context stops all deviation checkers.
	fm.jobs = map[string]models.JobSpec{}
	fm.checkers = map[uint]DeviationChecker{}
}

// OnNewHead is a noop.
func (fm *concreteFluxMonitor) OnNewHead(*models.Head) {}

// AddJob created a DeviationChecker for any job initiators of type
// InitiatorFluxMonitor.
func (fm *concreteFluxMonitor) AddJob(job models.JobSpec) error {
	fm.checkersMutex.Lock()
	defer fm.checkersMutex.Unlock()

	if _, ok := fm.jobs[job.ID.String()]; ok {
		return fmt.Errorf("job %s has already been added to flux monitor", job.ID.String())
	}

	validCheckers := map[uint]DeviationChecker{}
	for _, initr := range job.InitiatorsFor(models.InitiatorFluxMonitor) {
		checker, err := fm.checkerFactory.New(fm.ctx, initr, fm.runManager)
		if err != nil {
			return err
		}
		err = checker.Initialize(fm.store.TxManager)
		if err != nil {
			return err
		}
		validCheckers[initr.ID] = checker
	}

	for id, checker := range validCheckers {
		fm.checkers[id] = checker
		go checker.Start()
	}

	fm.jobs[job.ID.String()] = job
	return nil
}

// RemoveJob stops and removes the checker for all Flux Monitor initiators belonging
// to the passed job ID.
func (fm *concreteFluxMonitor) RemoveJob(ID *models.ID) {
	fm.checkersMutex.Lock()
	defer fm.checkersMutex.Unlock()

	job, ok := fm.jobs[ID.String()]
	if !ok {
		return
	}
	delete(fm.jobs, ID.String())
	for _, initr := range job.InitiatorsFor(models.InitiatorFluxMonitor) {
		checker := fm.checkers[initr.ID]
		delete(fm.checkers, initr.ID)
		checker.Stop()
	}
}

// DeviationCheckerFactory holds the New method needed to create a new instance
// of a DeviationChecker.
type DeviationCheckerFactory interface {
	New(context.Context, models.Initiator, RunManager) (DeviationChecker, error)
}

type pollingDeviationCheckerFactory struct{}

func (f pollingDeviationCheckerFactory) New(parentCtx context.Context, initr models.Initiator, runManager RunManager) (DeviationChecker, error) {
	return NewPollingDeviationChecker(parentCtx, initr, runManager)
}

// DeviationChecker encapsulate methods needed to initialize and check prices
// for deviations, or swings.
type DeviationChecker interface {
	Initialize(eth.Client) error
	Start()
	Stop()
}

// PollingDeviationChecker polls external price adapters via HTTP to check for price swings.
type PollingDeviationChecker struct {
	initr         models.Initiator
	address       common.Address
	requestData   models.JSON
	threshold     float64
	precision     int32
	runManager    RunManager
	previousPrice decimal.Decimal
	fetcher       fetcher
	ctx           context.Context
	cancel        context.CancelFunc
}

// defaultHTTPTimeout is the timeout used by the price adapter fetcher for outgoing HTTP requests.
const defaultHTTPTimeout = 5 * time.Second

// NewPollingDeviationChecker returns a new instance of PollingDeviationChecker.
func NewPollingDeviationChecker(parentCtx context.Context, initr models.Initiator, runManager RunManager) (*PollingDeviationChecker, error) {
	fetcher, err := newMedianFetcherFromURLs(
		defaultHTTPTimeout,
		initr.InitiatorParams.RequestData.String(),
		initr.InitiatorParams.Feeds...)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(parentCtx)
	return &PollingDeviationChecker{
		initr:         initr,
		address:       initr.InitiatorParams.Address,
		requestData:   initr.InitiatorParams.RequestData,
		threshold:     float64(initr.InitiatorParams.Threshold),
		precision:     initr.InitiatorParams.Precision,
		runManager:    runManager,
		previousPrice: decimal.NewFromInt(0),
		ctx:           ctx,
		cancel:        cancel,
		fetcher:       fetcher,
	}, nil
}

// Initialize retrieves the price that's on-chain, with which we must check
// the deviation from.
func (p *PollingDeviationChecker) Initialize(client eth.Client) error {
	price, err := client.GetAggregatorPrice(p.address, p.precision)
	p.previousPrice = price
	return err
}

// Start begins a loop polling the price adapters set in the InitiatorFluxMonitor.
func (p *PollingDeviationChecker) Start() {
	for {
		logger.ErrorIf(p.Run())

		select {
		case <-p.ctx.Done():
			return
		case <-time.After(1 * time.Second):
		}
	}
}

// Stop stops this instance from polling, cleaning up resources.
func (p *PollingDeviationChecker) Stop() {
	p.cancel()
}

// PreviousPrice returns the price used to check deviations against.
func (p *PollingDeviationChecker) PreviousPrice() decimal.Decimal {
	return p.previousPrice
}

// Run walks through the steps to check for a deviation, early exiting if any particular
// step fails, or triggering a new job run.
// Uses a railway paradigm: https://fsharpforfunandprofit.com/rop/
func (p *PollingDeviationChecker) Run() error {
	_, err := railway(
		newData(p.previousPrice),
		railwayStep{"fetch current prices", p.fetchPrices},
		railwayStep{"check if outside deviation", p.checkIfOutsideDeviation},
		railwayStep{"create job run", p.createJobRun},
		railwayStep{"update previous price", p.updatePreviousPrice}, // only reached if outside deviation
	)
	return err
}

// railway walks through a set of steps linearly, with each step having one of three
// states: stop, continue, error.
func railway(d *data, steps ...railwayStep) (*data, error) {
	for _, step := range steps {
		exit, err := step.operation(d)
		if err != nil {
			return d, errors.Wrapf(err, "on step %s", step.label)
		}
		if exit != nil {
			logger.Infow(fmt.Sprintf("%s early exited: %s", step.label, exit.reason), "step", step.label)
			return d, nil
		}
	}
	return d, nil
}

type railwayStep struct {
	label     string
	operation func(*data) (*railwayExit, error)
}

type railwayExit struct {
	reason string
}

func newRailwayExit(reason string) *railwayExit {
	return &railwayExit{reason: reason}
}

func (p *PollingDeviationChecker) fetchPrices(d *data) (*railwayExit, error) {
	median, err := p.fetcher.Fetch()
	d.MedianPrice = decimal.NewFromFloat(median)
	return nil, errors.Wrap(err, "unable to fetch median price")
}

func (p *PollingDeviationChecker) checkIfOutsideDeviation(d *data) (*railwayExit, error) {
	prevPrice := d.PreviousPrice
	diff := prevPrice.Sub(d.MedianPrice).Abs()
	perc := diff.Div(prevPrice).Mul(decimal.NewFromInt(100))
	logger.Infow(
		fmt.Sprintf("deviation of %v%% for threshold %f%% with %s", perc, p.threshold, d),
		"threshold", p.threshold,
		"deviation", perc,
	)
	if perc.LessThan(decimal.NewFromFloat(p.threshold)) {
		reason := fmt.Sprintf("difference is %v%%, deviation threshold of %f%% not met", perc, p.threshold)
		return newRailwayExit(reason), nil
	}
	return nil, nil
}

func (p *PollingDeviationChecker) createJobRun(d *data) (*railwayExit, error) {
	runData, err := models.JSON{}.Add("result", fmt.Sprintf("%v", d.MedianPrice))
	if err != nil {
		return nil, errors.Wrap(err, "unable to start chainlink run")
	}
	_, err = p.runManager.Create(p.initr.JobSpecID, &p.initr, &runData, nil, models.NewRunRequest())
	return nil, err
}

func (p *PollingDeviationChecker) updatePreviousPrice(d *data) (*railwayExit, error) {
	p.previousPrice = d.MedianPrice
	return nil, nil
}

type data struct {
	MedianPrice   decimal.Decimal
	PreviousPrice decimal.Decimal
}

func newData(previousPrice decimal.Decimal) *data {
	return &data{
		PreviousPrice: previousPrice,
	}
}

func (d *data) String() string {
	return fmt.Sprintf("previous: %v, current median: %v", d.PreviousPrice, d.MedianPrice)
}
