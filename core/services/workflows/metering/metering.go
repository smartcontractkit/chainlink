package metering

import (
	"context"
	"errors"
	"fmt"
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
)

type BillingClient interface {
	SubmitWorkflowReceipt(context.Context, *billing.SubmitWorkflowReceiptRequest) (*billing.SubmitWorkflowReceiptResponse, error)
	ReserveCredits(context.Context, *billing.ReserveCreditsRequest) (*billing.ReserveCreditsResponse, error)
}

type SpendUnit string

func (s SpendUnit) String() string {
	return string(s)
}

func (s SpendUnit) DecimalToSpendValue(value decimal.Decimal) SpendValue {
	return SpendValue{value: value, roundingPlace: 18}
}

func (s SpendUnit) IntToSpendValue(value int64) SpendValue {
	return SpendValue{value: decimal.NewFromInt(value), roundingPlace: 18}
}

func (s SpendUnit) StringToSpendValue(value string) (SpendValue, error) {
	dec, err := decimal.NewFromString(value)
	if err != nil {
		return SpendValue{}, err
	}

	return SpendValue{value: dec, roundingPlace: 18}, nil
}

type SpendValue struct {
	value         decimal.Decimal
	roundingPlace uint8
}

func (v SpendValue) Add(value SpendValue) SpendValue {
	return SpendValue{
		value:         v.value.Add(value.value),
		roundingPlace: v.roundingPlace,
	}
}

func (v SpendValue) Div(value SpendValue) SpendValue {
	return SpendValue{
		value:         v.value.Div(value.value),
		roundingPlace: v.roundingPlace,
	}
}

func (v SpendValue) GreaterThan(value SpendValue) bool {
	return v.value.GreaterThan(value.value)
}

func (v SpendValue) String() string {
	return v.value.StringFixedBank(int32(v.roundingPlace))
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
	// The actual spend of this step
	Spend map[SpendUnit][]ReportStepDetail
}

type ReportStepDetail struct {
	Peer2PeerID string
	SpendValue  SpendValue
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

// Reserve calls the billing service for the initial credit balance that can be used in an execution
// This method must be called before Deduct or Settle
func (r *Report) Reserve(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.client == nil {
		// TODO: https://smartcontract-it.atlassian.net/browse/CRE-427 more robust check of billing service health
		return ErrNoBillingClient
	}

	// TODO: https://smartcontract-it.atlassian.net/browse/CRE-460 get rate card from billing service
	rateCard := map[string]decimal.Decimal{}

	balanceStore := NewBalanceStore(0, rateCard, r.lggr)

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
		// TODO: https://smartcontract-it.atlassian.net/browse/CRE-453 track causes of metering mode
		balanceStore.AllowNegative()
		r.lggr.Warnf("failed to reserve credits: %s", err)
	} else {
		if success := resp.GetSuccess(); !success {
			return ErrInsufficientFunding
		}
		// TODO: https://smartcontract-it.atlassian.net/browse/CRE-290 once billing client response contains balance set using balanceStore.Add
		dummyInitialBalance := int64(10000)
		if addErr := balanceStore.Add(dummyInitialBalance); addErr != nil {
			return addErr
		}
	}

	r.ready = true
	r.balance = balanceStore
	return nil
}

func (r *Report) MedianSpend() map[SpendUnit]SpendValue {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := map[SpendUnit][]SpendValue{}
	medians := map[SpendUnit]SpendValue{}

	for _, step := range r.steps {
		for unit, details := range step.Spend {
			_, ok := values[unit]
			if !ok {
				values[unit] = []SpendValue{}
			}

			for _, detail := range details {
				values[unit] = append(values[unit], detail.SpendValue)
			}
		}
	}

	for unit, set := range values {
		sort.Slice(set, func(i, j int) bool {
			return set[j].GreaterThan(set[i])
		})

		if len(set)%2 > 0 {
			medians[unit] = set[len(set)/2]

			continue
		}

		medians[unit] = set[len(set)/2-1].Add(set[len(set)/2]).Div(unit.IntToSpendValue(2))
	}

	return medians
}

// ConvertFromBalance converts a credit amount to a resource dimensions amount
func (r *Report) ConvertFromBalance(toUnit string, amount int64) (resources int64) {
	return r.balance.ConvertFromBalance(toUnit, amount)
}

// ConvertToBalance converts a resource dimensions amount to a credit amount
func (r *Report) ConvertToBalance(fromUnit string, amount int64) (credits int64) {
	return r.balance.ConvertToBalance(fromUnit, amount)
}

// Deduct earmarks an amount of local universal credit balance
// We expect to only set this value once - an error is returned if a step would be overwritten
func (r *Report) Deduct(ref string, amount int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.ready {
		return ErrNoReserve
	}

	if _, ok := r.steps[ref]; ok {
		return ErrStepDeductExists
	}

	err := r.balance.Minus(amount)
	if err != nil {
		return err
	}

	r.steps[ref] = ReportStep{
		Deduction: amount,
		Spend:     nil,
	}

	return nil
}

// GetAvailablity returns the amount of balance that is available to be used when split among a number of potential concurrent calls
func (r *Report) GetAvailablity(openConcurrentCallSlots int) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if openConcurrentCallSlots == 0 {
		return 0, ErrNoOpenCalls
	}

	// Split the available local balance between the number of concurrent calls that can still be made
	available := r.balance.Get()
	share := decimal.NewFromInt(available).Div(decimal.NewFromInt(int64(openConcurrentCallSlots)))
	roundedShare := share.RoundDown(0).IntPart()

	return roundedShare, nil
}

// Settle records the actual spend for a given capability invocation in the engine.
// The Deduct method must be called before Settle
// We expect to only set this value once - an error is returned if a step would be overwritten
func (r *Report) Settle(ref string, steps []capabilities.MeteringNodeDetail) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.ready {
		return ErrNoReserve
	}

	step, ok := r.steps[ref]
	if !ok {
		return ErrNoDeduct
	}

	if step.Spend != nil {
		return ErrStepSpendExists
	}

	spent := int64(0)
	spends := make(map[SpendUnit][]ReportStepDetail)

	for _, detail := range steps {
		unit := SpendUnit(detail.SpendUnit)
		value, err := unit.StringToSpendValue(detail.SpendValue)
		if err != nil {
			r.lggr.Error(fmt.Sprintf("failed to get spend value from %s: %s", detail.SpendValue, err))
		}
		spends[unit] = append(spends[unit], ReportStepDetail{
			Peer2PeerID: detail.Peer2PeerID,
			SpendValue:  value,
		})
		spent += r.balance.ConvertToBalance(detail.SpendUnit, value.value.IntPart())
	}

	step.Spend = spends
	r.steps[ref] = step

	// Refund the difference between what local balance what earmarked and the actual spend
	err := r.balance.Add(step.Deduction - spent)
	if err != nil {
		// invariant: capability should not let spend exceed reserve
		r.lggr.Error("invariant: spend exceeded reserve")
	}

	return nil
}

func (r *Report) Message() *events.MeteringReport {
	protoReport := &events.MeteringReport{
		Steps:    map[string]*events.MeteringReportStep{},
		Metadata: &events.WorkflowMetadata{},
	}

	for key, step := range r.steps {
		nodeDetails := []*events.MeteringReportNodeDetail{}

		for unit, details := range step.Spend {
			for _, detail := range details {
				nodeDetails = append(nodeDetails, &events.MeteringReportNodeDetail{
					Peer_2PeerId: detail.Peer2PeerID,
					SpendUnit:    unit.String(),
					SpendValue:   detail.SpendValue.String(),
				})
			}
		}

		protoReport.Steps[key] = &events.MeteringReportStep{
			Nodes: nodeDetails,
		}
	}

	return protoReport
}

func (r *Report) SendReceipt(ctx context.Context) error {
	if !r.ready {
		return ErrNoReserve
	}

	req := billing.SubmitWorkflowReceiptRequest{
		AccountId:           r.owner,
		WorkflowId:          r.workflowID,
		WorkflowExecutionId: r.workflowExecutionID,
		Metering:            r.Message(),
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
	// sugaredLggr := lggr.Named("Metering").With("workflowExecutionID", workflowExecutionID)

	return &Reports{
		reports: make(map[string]*Report),
		client:  client,

		lggr: lggr,

		owner:      owner,
		workflowID: workflowID,
	}
}

// Get retrieves a Report for a given key (if it exists).
func (s *Reports) Get(key string) (*Report, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.reports[key]
	return val, ok
}

// Start creates a new report and inserts it under the specified key.
func (s *Reports) Start(ctx context.Context, key string) (*Report, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.reports[key]
	if ok {
		return nil, errors.New("report already exists")
	}

	report := NewReport(s.owner, s.workflowID, key, s.lggr, s.client)

	err := report.Reserve(ctx)
	if err != nil {
		return nil, err
	}

	s.reports[key] = report

	return report, nil
}

// End removes the Report with the specified key.
func (s *Reports) End(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	report, ok := s.reports[key]
	if !ok {
		return errors.New("report not found")
	}

	if err := report.SendReceipt(ctx); err != nil {
		return err
	}

	delete(s.reports, key)

	return nil
}

func (s *Reports) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.reports)
}
