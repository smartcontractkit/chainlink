package metering

import (
	"context"
	"errors"
	"fmt"
	"log"
	"maps"
	"sort"
	"sync"

	"github.com/shopspring/decimal"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	protoEvents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	wfEvents "github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

const (
	ComputeResourceDimension = "COMPUTE"
	RatiosKey                = "spendRatios"
	defaultDecimalPrecision  = 3 // one thousandth of a dollar
)

var (
	ErrMissingLabels         = errors.New("missing required labels: owner, workflowID, workflowExecutionID")
	ErrNoBillingClient       = errors.New("no billing client has been configured")
	ErrInsufficientFunding   = errors.New("insufficient funding")
	ErrReceiptFailed         = errors.New("failed to submit workflow receipt")
	ErrNoReserve             = errors.New("must call Reserve first")
	ErrStepDeductExists      = errors.New("step deduct already exists")
	ErrNoOpenCalls           = errors.New("openConcurrentCallSlots must be greater than 0")
	ErrNoDeduct              = errors.New("must call Deduct first")
	ErrStepSpendExists       = errors.New("step spend already exists")
	ErrReportNotFound        = errors.New("report not found")
	ErrReportExists          = errors.New("report already exists")
	ErrRatiosAndTypesNoMatch = errors.New("spending types and ratios do not match")
	ErrInvalidRatios         = errors.New("invalid spending type ratios")
)

type BillingClient interface {
	SubmitWorkflowReceipt(context.Context, *billing.SubmitWorkflowReceiptRequest) (*billing.SubmitWorkflowReceiptResponse, error)
	ReserveCredits(context.Context, *billing.ReserveCreditsRequest) (*billing.ReserveCreditsResponse, error)
}

type SpendTuple struct {
	Unit  string
	Value decimal.Decimal
}

type ProtoDetail struct {
	Schema string
	Domain string
	Entity string
}

type ReportStep struct {
	// The maximum amount of universal credits that should be used in this step
	Deduction decimal.Decimal
	// The actual resource spend that each node used for this step
	Spends map[string][]ReportStepDetail
}

type ReportStepDetail struct {
	Peer2PeerID string
	SpendValue  string
}

type Report struct {
	// descriptive properties
	labels map[string]string

	// dependencies
	balance *balanceStore
	client  BillingClient
	lggr    logger.Logger

	// internal state
	mu    sync.RWMutex
	ready bool

	// meteringMode turns off double spend checks.
	// In meteringMode, no accounting wrt universal credits is required;
	// only gathering resource types and spends from capabilities.
	// note: meteringMode == true allows negative balances.
	meteringMode bool
	steps        map[string]ReportStep
}

func NewReport(labels map[string]string, lggr logger.Logger, client BillingClient) (*Report, error) {
	requiredLabels := []string{platform.KeyWorkflowOwner, platform.KeyWorkflowID, platform.KeyWorkflowExecutionID}
	for _, label := range requiredLabels {
		_, ok := labels[label]
		if !ok {
			return nil, ErrMissingLabels
		}
	}

	balanceStore, err := NewBalanceStore(decimal.Zero, map[string]decimal.Decimal{})
	if err != nil {
		return nil, err
	}

	return &Report{
		labels: labels,

		balance: balanceStore,
		client:  client,
		lggr:    logger.Sugared(lggr).Named("Metering").With(platform.KeyWorkflowExecutionID, labels[platform.KeyWorkflowExecutionID]),

		ready:        false,
		meteringMode: false,
		steps:        make(map[string]ReportStep),
	}, nil
}

// Reserve calls the billing service for the initial credit balance that can be used in an execution.
// This method must be called before Deduct or Settle.
func (r *Report) Reserve(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.client == nil {
		r.switchToMeteringMode(ErrNoBillingClient)

		return nil
	}

	// TODO: https://smartcontract-it.atlassian.net/browse/CRE-427 more robust check of billing service health

	// If there is no credit limit defined in the workflow, then open an empty reservation
	// TODO: https://smartcontract-it.atlassian.net/browse/CRE-284 consume user defined workflow execution limit
	req := billing.ReserveCreditsRequest{
		AccountId:           r.labels[platform.KeyWorkflowOwner],
		WorkflowId:          r.labels[platform.KeyWorkflowID],
		WorkflowExecutionId: r.labels[platform.KeyWorkflowExecutionID],
		Credits:             0,
	}

	resp, err := r.client.ReserveCredits(ctx, &req)

	// If there is an error communicating with the billing service, fail open
	if err != nil {
		r.switchToMeteringMode(err)

		return nil
	}

	if success := resp.GetSuccess(); !success {
		return ErrInsufficientFunding
	}

	rateCard, err := toRateCard(resp.GetRates())
	if err != nil {
		r.switchToMeteringMode(err)

		return nil
	}

	balanceStore, err := NewBalanceStore(decimal.NewFromFloat32(resp.Credits), rateCard)
	if err != nil {
		r.switchToMeteringMode(err)

		return nil
	}

	r.ready = true
	r.balance = balanceStore

	return nil
}

// ConvertToBalance converts a resource dimensions amount to a credit amount.
func (r *Report) ConvertToBalance(fromUnit string, amount decimal.Decimal) (decimal.Decimal, error) {
	if !r.ready {
		return decimal.Zero, ErrNoReserve
	}

	bal, err := r.balance.ConvertToBalance(fromUnit, amount)
	if err != nil {
		// Fail open, continue optimistically
		r.switchToMeteringMode(err)
	}

	return bal, nil
}

// Deduct earmarks an amount of local universal credit balance. We expect to only set this value once - an error is
// returned if a step would be overwritten.
func (r *Report) Deduct(ref string, amount decimal.Decimal) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.ready {
		return ErrNoReserve
	}

	if _, ok := r.steps[ref]; ok {
		return ErrStepDeductExists
	}

	r.steps[ref] = ReportStep{
		Deduction: amount,
		Spends:    nil,
	}

	// if in metering mode, exit early without modifying local balance
	if r.meteringMode {
		return nil
	}

	return r.balance.Minus(amount)
}

// CreditToSpendingLimits returns a slice of spend limits where the amount is applied to the spend types from the
// provided info. Amount should be specified in universal credits and will be converted to spend type credits within
// this function.
func (r *Report) CreditToSpendingLimits(
	info capabilities.CapabilityInfo,
	config *values.Map,
	amount decimal.Decimal,
) []capabilities.SpendLimit {
	if r.meteringMode {
		return []capabilities.SpendLimit{}
	}

	// no spend types results in no limits and is not a failure case
	if len(info.SpendTypes) == 0 {
		return []capabilities.SpendLimit{}
	}

	ratios, err := ratiosFromConfig(info, config)
	if err != nil {
		r.switchToMeteringMode(err)

		return []capabilities.SpendLimit{}
	}

	// spend types do not have matching ratios; this is a bad configuration
	if len(info.SpendTypes) != len(ratios) {
		r.switchToMeteringMode(fmt.Errorf("%w: %d spend types and %d ratios", ErrRatiosAndTypesNoMatch, len(info.SpendTypes), len(ratios)))

		return []capabilities.SpendLimit{}
	}

	limits := []capabilities.SpendLimit{}

	for _, spendType := range info.SpendTypes {
		ratio, hasRatio := ratios[spendType]
		if !hasRatio {
			// the spend type does not exist in the ratios mapping; this is a bad configuration
			r.switchToMeteringMode(fmt.Errorf("%w: ratios missing %s spend type", ErrRatiosAndTypesNoMatch, spendType))

			return []capabilities.SpendLimit{}
		}

		// use rate card to convert capSpendLimit to native units
		spendLimit, err := r.balance.ConvertFromBalance(string(spendType), amount.Mul(ratio))
		if err != nil {
			r.switchToMeteringMode(err)

			return []capabilities.SpendLimit{}
		}

		limits = append(limits, capabilities.SpendLimit{SpendType: spendType, Limit: spendLimit.StringFixed(defaultDecimalPrecision)})
	}

	return limits
}

// GetMaxSpendForInvocation returns the amount of credits that can be used based on the minimum between an optionally
// provided max spend by the user or the available credit balance. The available credit balance is determined by
// dividing unearmarked local credit balance by the number of potential concurrent calls.
func (r *Report) GetMaxSpendForInvocation(
	userSpendLimit decimal.NullDecimal,
	openConcurrentCallSlots int,
) (decimal.NullDecimal, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	nullCapSpendLimit := decimal.NewNullDecimal(decimal.Zero)
	nullCapSpendLimit.Valid = false

	if openConcurrentCallSlots == 0 {
		// invariant: this should be managed by the consumer (engine)
		return nullCapSpendLimit, ErrNoOpenCalls
	}

	if !r.ready {
		return nullCapSpendLimit, ErrNoReserve
	}

	if r.meteringMode {
		return nullCapSpendLimit, nil
	}

	// Split the available local balance between the number of concurrent calls that can still be made
	spendLimit := r.balance.Get().Div(decimal.NewFromInt(int64(openConcurrentCallSlots)))

	if userSpendLimit.Valid {
		spendLimit = decimal.Min(spendLimit, userSpendLimit.Decimal)
	}

	return decimal.NewNullDecimal(spendLimit), nil
}

// Settle handles the actual spends that each node used for a given capability invocation in the engine,
// by returning earmarked local balance to the available to use pool and adding the spend to the metering report.
// The Deduct method must be called before Settle.
// We expect to only set this value once - an error is returned if a step would be overwritten.
func (r *Report) Settle(ref string, spendsByNode []capabilities.MeteringNodeDetail) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.ready {
		return ErrNoReserve
	}

	step, ok := r.steps[ref]
	if !ok {
		return ErrNoDeduct
	}

	if step.Spends != nil {
		return ErrStepSpendExists
	}

	spentCredits := decimal.NewFromInt(0)
	resourceSpends := make(map[string][]ReportStepDetail)

	// Group by resource dimension
	for _, nodeDetail := range spendsByNode {
		resourceSpends[nodeDetail.SpendUnit] = append(resourceSpends[nodeDetail.SpendUnit], ReportStepDetail{
			Peer2PeerID: nodeDetail.Peer2PeerID,
			SpendValue:  nodeDetail.SpendValue,
		})
	}

	// Aggregate node responses to a single number
	for unit, spendDetails := range resourceSpends {
		deciVals := []decimal.Decimal{}
		for _, detail := range spendDetails {
			value, err := decimal.NewFromString(detail.SpendValue)
			if err != nil {
				r.lggr.Error(fmt.Sprintf("failed to get spend value from %s: %s", detail.SpendValue, err))
				// throw out invalid values for local balance settlement. they will still be included in metering report.
				continue
			}

			deciVals = append(deciVals, value)
		}

		aggregateSpend := medianSpend(deciVals)
		bal, err := r.balance.ConvertToBalance(unit, aggregateSpend)
		if err != nil {
			r.switchToMeteringMode(err)
		}

		spentCredits = spentCredits.Add(bal)
	}

	step.Spends = resourceSpends
	r.steps[ref] = step

	// if in metering mode, exit early without modifying local balance
	if r.meteringMode {
		return nil
	}

	// Refund the difference between what local balance had been earmarked and the actual spend
	if err := r.balance.Add(step.Deduction.Sub(spentCredits)); err != nil {
		// invariant: capability should not let spend exceed reserve
		r.lggr.Error("invariant: spend exceeded reserve")
	}

	return nil
}

func (r *Report) FormatReport() *protoEvents.MeteringReport {
	protoReport := &protoEvents.MeteringReport{
		Steps:    map[string]*protoEvents.MeteringReportStep{},
		Metadata: &protoEvents.WorkflowMetadata{},
	}

	for ref, step := range r.steps {
		nodeDetails := []*protoEvents.MeteringReportNodeDetail{}

		for unit, details := range step.Spends {
			for _, detail := range details {
				nodeDetails = append(nodeDetails, &protoEvents.MeteringReportNodeDetail{
					Peer_2PeerId: detail.Peer2PeerID,
					SpendUnit:    unit,
					SpendValue:   detail.SpendValue,
				})
			}
		}

		protoReport.Steps[ref] = &protoEvents.MeteringReportStep{
			Nodes: nodeDetails,
		}
	}

	return protoReport
}

func (r *Report) SendReceipt(ctx context.Context) error {
	if !r.ready {
		return ErrNoReserve
	}

	if r.client == nil {
		return ErrNoBillingClient
	}

	// TODO: https://smartcontract-it.atlassian.net/browse/CRE-427 more robust check of billing service health

	req := billing.SubmitWorkflowReceiptRequest{
		AccountId:           r.labels[platform.KeyWorkflowOwner],
		WorkflowId:          r.labels[platform.KeyWorkflowID],
		WorkflowExecutionId: r.labels[platform.KeyWorkflowExecutionID],
		Metering:            r.FormatReport(),
	}

	resp, err := r.client.SubmitWorkflowReceipt(ctx, &req)
	if err != nil {
		return err
	}

	if resp == nil || !resp.Success {
		return ErrReceiptFailed
	}

	return nil
}

func (r *Report) EmitReceipt(ctx context.Context) error {
	if !r.ready {
		return ErrNoReserve
	}

	return wfEvents.EmitMeteringReport(ctx, r.labels, r.FormatReport())
}

func (r *Report) switchToMeteringMode(err error) {
	r.lggr.Errorf("switching to metering mode: %s", err)
	r.meteringMode = true
	r.ready = true
}

func toRateCard(rates []*billing.ResourceUnitRate) (map[string]decimal.Decimal, error) {
	rateCard := map[string]decimal.Decimal{}
	for _, rate := range rates {
		conversionDeci, err := decimal.NewFromString(rate.ConversionRate)
		if err != nil {
			return map[string]decimal.Decimal{}, fmt.Errorf("could not convert unit %s's value %s to decimal", rate.ResourceUnit, rate.ConversionRate)
		}
		rateCard[rate.ResourceUnit] = conversionDeci
	}
	return rateCard, nil
}

func medianSpend(spends []decimal.Decimal) decimal.Decimal {
	sort.Slice(spends, func(i, j int) bool {
		return spends[j].GreaterThan(spends[i])
	})

	if len(spends)%2 > 0 {
		return spends[len(spends)/2]
	}

	return spends[len(spends)/2-1].Add(spends[len(spends)/2]).Div(decimal.NewFromInt(2))
}

// Reports is a concurrency-safe wrapper around map[string]*Report.
type Reports struct {
	mu      sync.RWMutex
	reports map[string]*Report
	client  BillingClient
	lggr    logger.Logger
	metrics *monitoring.WorkflowsMetricLabeler

	// descriptive properties
	owner      string
	workflowID string
	labelMap   map[string]string
}

// NewReports initializes and returns a new Reports.
func NewReports(client BillingClient, owner, workflowID string, lggr logger.Logger, labels map[string]string, metrics *monitoring.WorkflowsMetricLabeler) *Reports {
	return &Reports{
		reports: make(map[string]*Report),
		client:  client,
		lggr:    lggr,
		metrics: metrics,

		owner:      owner,
		workflowID: workflowID,
		labelMap:   labels,
	}
}

// Get retrieves a Report for a given workflowExecutionID (if it exists).
func (s *Reports) Get(workflowExecutionID string) (*Report, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.reports[workflowExecutionID]
	return val, ok
}

// Start creates a new report and inserts it under the specified workflowExecutionID.
func (s *Reports) Start(ctx context.Context, workflowExecutionID string) (*Report, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.reports[workflowExecutionID]
	if ok {
		return nil, ErrReportExists
	}

	labels := map[string]string{}
	maps.Copy(labels, s.labelMap)
	labels[platform.KeyWorkflowExecutionID] = workflowExecutionID

	report, err := NewReport(labels, s.lggr, s.client)
	if err != nil {
		return nil, err
	}

	s.reports[workflowExecutionID] = report

	return report, nil
}

// End removes the Report with the specified workflowExecutionID.
func (s *Reports) End(ctx context.Context, workflowExecutionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	report, ok := s.reports[workflowExecutionID]
	if !ok {
		return ErrReportNotFound
	}

	var multiErr error

	emitErr := report.EmitReceipt(ctx)
	if emitErr != nil {
		s.metrics.IncrementWorkflowMissingMeteringReport(ctx)
		multiErr = multierr.Combine(multiErr, emitErr)
	}

	sendErr := report.SendReceipt(ctx)
	if sendErr != nil {
		s.metrics.IncrementWorkflowMissingMeteringReport(ctx)
		multiErr = multierr.Combine(multiErr, sendErr)
	}

	delete(s.reports, workflowExecutionID)

	if multiErr != nil {
		return multiErr
	}

	return nil
}

func (s *Reports) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.reports)
}

// ratiosFromConfig collects all ratios from a value map that match specified spend types. Any error will return an
// empty set of ratios with the error.
//
// CapabilityInfo contains information about the spend types while the registry config contains ratios for splitting
// spend types. This allows capability authors to not have to redeploy a capability to change spending ratios. The
// spending ratios was not put in the billing service because the ratios are not expected to change often. The registry
// is mutable enough for this purpose while the capability info.
func ratiosFromConfig(
	info capabilities.CapabilityInfo,
	config *values.Map,
) (map[capabilities.CapabilitySpendType]decimal.Decimal, error) {
	ratios := make(map[capabilities.CapabilitySpendType]decimal.Decimal)

	// if info.SpendTypes has only 1, return ratio 100%
	if len(info.SpendTypes) == 1 {
		ratios[info.SpendTypes[0]] = decimal.NewFromInt(1)

		return ratios, nil
	}

	if config == nil {
		return ratios, fmt.Errorf("%w: spending ratios not set; config is nil", ErrInvalidRatios)
	}

	rawRatiosValue, hasRatios := config.Underlying[RatiosKey]
	if !hasRatios {
		return ratios, fmt.Errorf("%w: spending ratios not set", ErrInvalidRatios)
	}

	rawRatiosAny, err := rawRatiosValue.Unwrap()
	if err != nil {
		return ratios, fmt.Errorf("%w: %w", ErrInvalidRatios, err)
	}

	rawRatios, ok := rawRatiosAny.(map[string]any)
	if !ok {
		return ratios, fmt.Errorf("%w: not a value map", ErrInvalidRatios)
	}

	for _, spendType := range info.SpendTypes {
		// using a namespace on the config key to distinguish billing specific keys
		value, hasRatio := rawRatios[string(spendType)]
		if !hasRatio {
			return make(map[capabilities.CapabilitySpendType]decimal.Decimal), fmt.Errorf("%w: ratio does not exist for: %s", ErrInvalidRatios, spendType)
		}

		strValue, ok := value.(string)
		if !ok {
			log.Println(strValue)
			return make(map[capabilities.CapabilitySpendType]decimal.Decimal), fmt.Errorf("%w: ratio for key '%s' should be type string", ErrInvalidRatios, spendType)
		}

		ratio, err := decimal.NewFromString(strValue)
		if err != nil {
			return make(map[capabilities.CapabilitySpendType]decimal.Decimal), fmt.Errorf("%w: could not unwrap decimal ratio value: %s", ErrInvalidRatios, value)
		}

		ratios[spendType] = ratio
	}

	return ratios, nil
}
