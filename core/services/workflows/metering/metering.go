package metering

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	"github.com/smartcontractkit/chainlink-protos/workflows/go/events"
)

const (
	ComputeResourceDimension = "COMPUTE"
)

var (
	ErrNoBillingClient     = errors.New("no billing client has been configured")
	ErrInsufficientFunding = errors.New("insufficient funding")
	ErrReceiptFailed       = errors.New("failed to submit workflow receipt")
	ErrNoReserve           = errors.New("must call Reserve first")
	ErrStepDeductExists    = errors.New("step deduct already exists")
	ErrNoOpenCalls         = errors.New("openConcurrentCallSlots must be greater than 0")
	ErrNoDeduct            = errors.New("must call Deduct first")
	ErrStepSpendExists     = errors.New("step spend already exists")
	ErrReportNotFound      = errors.New("report not found")
	ErrReportExists        = errors.New("report already exists")
)

type BillingClient interface {
	SubmitWorkflowReceipt(context.Context, *billing.SubmitWorkflowReceiptRequest) (*billing.SubmitWorkflowReceiptResponse, error)
	ReserveCredits(context.Context, *billing.ReserveCreditsRequest) (*billing.ReserveCreditsResponse, error)
}

type SpendTuple struct {
	Unit  string
	Value int64
}

type ProtoDetail struct {
	Schema string
	Domain string
	Entity string
}

type ReportStep struct {
	// The maximum amount of universal credits that should be used in this step
	Deduction int64
	// The actual resource spend that each node used for this step
	Spends map[string][]ReportStepDetail
}

type ReportStepDetail struct {
	Peer2PeerID string
	SpendValue  string
}

type Report struct {
	// descriptive properties
	owner               string
	workflowID          string
	workflowExecutionID string

	// dependencies
	balance *balanceStore
	client  BillingClient
	lggr    logger.Logger

	// internal state
	ready bool
	mu    sync.RWMutex
	steps map[string]ReportStep
}

func NewReport(owner, workflowID, workflowExecutionID string, lggr logger.Logger, client BillingClient) *Report {
	return &Report{
		owner:               owner,
		workflowID:          workflowID,
		workflowExecutionID: workflowExecutionID,

		client: client,
		lggr:   logger.Sugared(lggr).Named("Metering").With("workflowExecutionID", workflowExecutionID),

		ready: false,
		steps: make(map[string]ReportStep),
	}
}

// Reserve calls the billing service for the initial credit balance that can be used in an execution.
// This method must be called before Deduct or Settle.
func (r *Report) Reserve(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.client == nil {
		return ErrNoBillingClient
	}

	// TODO: https://smartcontract-it.atlassian.net/browse/CRE-427 more robust check of billing service health

	// If there is no credit limit defined in the workflow, then open an empty reservation
	// TODO: https://smartcontract-it.atlassian.net/browse/CRE-284 consume user defined workflow execution limit
	req := billing.ReserveCreditsRequest{
		AccountId:           r.owner,
		WorkflowId:          r.workflowID,
		WorkflowExecutionId: r.workflowExecutionID,
		Credits:             []*billing.AccountCreditsInput{}, // TODO: https://smartcontract-it.atlassian.net/browse/CRE-290 send the credit balance, not resource types
	}

	resp, err := r.client.ReserveCredits(ctx, &req)

	// If there is an error communicating with the billing service, fail open
	if err != nil {
		r.lggr.Warnf("failed to reserve credits: %s", err)
		r.enterMeteringMode()
		return nil
	}

	if success := resp.GetSuccess(); !success {
		return ErrInsufficientFunding
	}

	rateCard, err := toRateCard(resp.GetRates())
	if err != nil {
		r.lggr.Warnf("failed to parse rate card: %s", err)
		r.enterMeteringMode()
		return nil
	}

	// TODO: https://smartcontract-it.atlassian.net/browse/CRE-290 once billing client response contains balance set using balanceStore.Add
	dummyInitialBalance := int64(10000)
	r.ready = true
	r.balance = NewBalanceStore(dummyInitialBalance, rateCard, r.lggr)
	return nil
}

func (r *Report) enterMeteringMode() {
	// TODO: https://smartcontract-it.atlassian.net/browse/CRE-453 pass through errors and persist cause of metering mode on to meteringReport
	balanceStore := NewBalanceStore(0, map[string]decimal.Decimal{}, r.lggr)
	balanceStore.AllowNegative()
	r.ready = true
	r.balance = balanceStore
}

// ConvertFromBalance converts a credit amount to a resource dimensions amount.
func (r *Report) ConvertFromBalance(toUnit string, amount int64) (resources int64, err error) {
	if !r.ready {
		return 0, ErrNoReserve
	}
	return r.balance.ConvertFromBalance(toUnit, amount), nil
}

// ConvertToBalance converts a resource dimensions amount to a credit amount.
func (r *Report) ConvertToBalance(fromUnit string, amount int64) (credits int64, err error) {
	if !r.ready {
		return 0, ErrNoReserve
	}
	return r.balance.ConvertToBalance(fromUnit, amount), nil
}

// Deduct earmarks an amount of local universal credit balance.
// We expect to only set this value once - an error is returned if a step would be overwritten.
func (r *Report) Deduct(ref string, amount int64) error {
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
	if r.balance.allowNegative {
		return nil
	}

	err := r.balance.Minus(amount)
	if err != nil {
		return err
	}

	return nil
}

// GetAvailableForInvocation returns the amount of credits that can be used based on the available credit balance.
// This is determined by dividing unearmarked local credit balance by the number of potential concurrent calls.
func (r *Report) GetAvailableForInvocation(openConcurrentCallSlots int) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if openConcurrentCallSlots == 0 {
		// invariant: this should be managed by the consumer (engine)
		return 0, ErrNoOpenCalls
	}

	if !r.ready {
		return 0, ErrNoReserve
	}

	if r.balance.allowNegative {
		return math.MaxInt64, nil
	}

	// Split the available local balance between the potential number of concurrent calls that can be made
	available := r.balance.Get()
	share := decimal.NewFromInt(available).Div(decimal.NewFromInt(int64(openConcurrentCallSlots)))
	roundedShare := share.RoundDown(0).IntPart()

	return roundedShare, nil
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

	spentCredits := int64(0)
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

		spentCredits += r.balance.ConvertToBalance(unit, aggregateSpend.IntPart())
	}

	step.Spends = resourceSpends
	r.steps[ref] = step

	// if in metering mode, exit early without modifying local balance
	if r.balance.allowNegative {
		return nil
	}

	// Refund the difference between what local balance had been earmarked and the actual spend
	err := r.balance.Add(step.Deduction - spentCredits)
	if err != nil {
		// invariant: capability should not let spend exceed reserve
		r.lggr.Error("invariant: spend exceeded reserve")
	}

	return nil
}

func (r *Report) FormatReport() *events.MeteringReport {
	protoReport := &events.MeteringReport{
		Steps:    map[string]*events.MeteringReportStep{},
		Metadata: &events.WorkflowMetadata{},
	}

	for ref, step := range r.steps {
		nodeDetails := []*events.MeteringReportNodeDetail{}

		for unit, details := range step.Spends {
			for _, detail := range details {
				nodeDetails = append(nodeDetails, &events.MeteringReportNodeDetail{
					Peer_2PeerId: detail.Peer2PeerID,
					SpendUnit:    unit,
					SpendValue:   detail.SpendValue,
				})
			}
		}

		protoReport.Steps[ref] = &events.MeteringReportStep{
			Nodes: nodeDetails,
		}
	}

	return protoReport
}

func (r *Report) SendReceipt(ctx context.Context) error {
	if !r.ready {
		return ErrNoReserve
	}

	// TODO: https://smartcontract-it.atlassian.net/browse/CRE-427 more robust check of billing service health

	req := billing.SubmitWorkflowReceiptRequest{
		AccountId:           r.owner,
		WorkflowId:          r.workflowID,
		WorkflowExecutionId: r.workflowExecutionID,
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

	// descriptive properties
	owner      string
	workflowID string
}

// NewReports initializes and returns a new Reports.
func NewReports(client BillingClient, owner, workflowID string, lggr logger.Logger) *Reports {
	return &Reports{
		reports: make(map[string]*Report),
		client:  client,

		lggr: lggr,

		owner:      owner,
		workflowID: workflowID,
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

	report := NewReport(s.owner, s.workflowID, workflowExecutionID, s.lggr, s.client)

	if s.client == nil {
		return nil, ErrNoBillingClient
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

	err := report.SendReceipt(ctx)

	delete(s.reports, workflowExecutionID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Reports) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.reports)
}
