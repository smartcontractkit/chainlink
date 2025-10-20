package metering

import (
	"context"
	"errors"
	"fmt"
	"log"
	"maps"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	protoEvents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	wfEvents "github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

const (
	RatiosKey = "spendRatios"
	// the default decimal precision is a fixed number defined in the billing service. if this gets changed
	// in the billing service project, the value here needs to change.
	defaultDecimalPrecision = 10 // one thousandth of a dollar

	EngineVersionV1 = "v1"
	EngineVersionV2 = "v2"
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
	ErrDeductOptionRequired  = errors.New("deduct option required")
	ErrEmptyRateCard         = errors.New("empty rate card")
)

type BillingClient interface {
	GetOrganizationCreditsByWorkflow(context.Context, *billing.GetOrganizationCreditsByWorkflowRequest) (*billing.GetOrganizationCreditsByWorkflowResponse, error)
	GetWorkflowExecutionRates(context.Context, *billing.GetWorkflowExecutionRatesRequest) (*billing.GetWorkflowExecutionRatesResponse, error)
	ReserveCredits(context.Context, *billing.ReserveCreditsRequest) (*billing.ReserveCreditsResponse, error)
	SubmitWorkflowReceipt(context.Context, *billing.SubmitWorkflowReceiptRequest) (*emptypb.Empty, error)
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
	// The ID of the capability being used in this step
	CapabilityID string
	// CapDONN is the total number of nodes in a capability DON.
	CapdonN uint32
	// The maximum amount of universal credits that should be used in this step
	Deduction decimal.Decimal
	// The actual resource spend that each node used for this step
	Spends           map[string][]ReportStepDetail
	AggregatedSpends map[string]AggregatedStepDetail
}

type ReportStepDetail struct {
	Peer2PeerID   string
	SpendValue    string
	CRESpendValue decimal.Decimal
}

type AggregatedStepDetail struct {
	SpendUnit     string
	SpendValue    decimal.Decimal
	CRESpendValue decimal.Decimal
}

type Report struct {
	// descriptive properties
	labels        map[string]string
	engineVersion string

	// dependencies
	balance *balanceStore
	client  BillingClient
	lggr    logger.Logger
	metrics *monitoring.WorkflowsMetricLabeler

	// internal state
	mu       sync.RWMutex
	reserved bool

	// meteringMode turns off double spend checks.
	// In meteringMode, no accounting wrt universal credits is required;
	// only gathering resource types and spends from capabilities.
	// note: meteringMode == true allows negative balances.
	meteringMode    bool
	meteringModeErr error
	steps           map[string]ReportStep
	rateCard        map[string]decimal.Decimal
	stepRefLookup   []string

	// WorkflowRegistryAddress is the address of the workflow registry contract
	workflowRegistryAddress string
	// WorkflowRegistryChainSelector is the chain selector for the workflow registry
	workflowRegistryChainSelector uint64
}

func NewReport(
	ctx context.Context,
	labels map[string]string,
	lggr logger.Logger,
	client BillingClient,
	metrics *monitoring.WorkflowsMetricLabeler,
	workflowRegistryAddress, workflowRegistryChainSelector, engineVersion string,
) (*Report, error) {
	requiredLabels := []string{platform.KeyWorkflowOwner, platform.KeyWorkflowID, platform.KeyWorkflowExecutionID}
	for _, label := range requiredLabels {
		_, ok := labels[label]
		if !ok {
			return nil, ErrMissingLabels
		}
	}

	report := &Report{
		labels:                  labels,
		lggr:                    logger.Sugared(lggr).Named("Metering").With(platform.KeyWorkflowExecutionID, labels[platform.KeyWorkflowExecutionID]),
		metrics:                 metrics,
		workflowRegistryAddress: workflowRegistryAddress,
		rateCard:                make(map[string]decimal.Decimal),
		engineVersion:           engineVersion,

		reserved: false,
		steps:    make(map[string]ReportStep),
	}

	// for safety in evaluating the client interface.
	// the client could be a nil interface or a nil value that satisfies the interface.
	valOf := reflect.ValueOf(client)
	if valOf.IsValid() && valOf.IsNil() {
		client = nil
	}

	if client == nil {
		report.switchToMeteringMode(ErrNoBillingClient)
	}

	chainSelector, err := strconv.ParseUint(workflowRegistryChainSelector, 10, 64)
	if err != nil {
		report.switchToMeteringMode(fmt.Errorf("failed to parse registry chain selector: %w", err))
	}

	report.workflowRegistryChainSelector = chainSelector

	if client != nil {
		report.client = client

		var resp *billing.GetWorkflowExecutionRatesResponse

		resp, err = report.client.GetWorkflowExecutionRates(ctx, &billing.GetWorkflowExecutionRatesRequest{
			WorkflowOwner:           labels[platform.KeyWorkflowOwner],
			WorkflowRegistryAddress: report.workflowRegistryAddress,
			ChainSelector:           report.workflowRegistryChainSelector,
		})
		if err != nil {
			report.switchToMeteringMode(err)
		}

		report.rateCard, err = toRateCard(resp)
		if err != nil {
			report.switchToMeteringMode(err)
		}
	}

	if len(report.rateCard) == 0 {
		report.switchToMeteringMode(ErrEmptyRateCard)
	}

	report.balance, err = NewBalanceStore(decimal.Zero, report.rateCard)
	if err != nil {
		report.switchToMeteringMode(fmt.Errorf("failed to create balance store: %w", err))

		// we can recover with an empty rate card and in metering mode
		report.balance, err = NewBalanceStore(decimal.Zero, map[string]decimal.Decimal{})
		if err != nil {
			// this should never happen, but if it does, we cannot proceed
			return nil, err
		}
	}

	return report, nil
}

// Reserve calls the billing service for the initial credit balance that can be used in an execution.
// This method must be called before Deduct or Settle.
func (r *Report) Reserve(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// always indicate that reserve was called.
	r.reserved = true

	if r.client == nil {
		r.switchToMeteringMode(ErrNoBillingClient)

		return nil
	}

	// TODO: https://smartcontract-it.atlassian.net/browse/CRE-427 more robust check of billing service health

	// If there is no credit limit defined in the workflow, then open an empty reservation
	// TODO: https://smartcontract-it.atlassian.net/browse/CRE-284 consume user defined workflow execution limit

	req := billing.ReserveCreditsRequest{
		WorkflowOwner:                 r.labels[platform.KeyWorkflowOwner],
		WorkflowId:                    r.labels[platform.KeyWorkflowID],
		WorkflowExecutionId:           r.labels[platform.KeyWorkflowExecutionID],
		WorkflowRegistryAddress:       r.workflowRegistryAddress,
		WorkflowRegistryChainSelector: r.workflowRegistryChainSelector,
		Credits:                       nil,
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

	creditsStr := resp.GetCredits()
	if creditsStr == "" {
		r.lggr.Debug("empty credits; setting default of 0")
		creditsStr = "0"
	}

	credits, err := decimal.NewFromString(creditsStr)
	if err != nil {
		r.switchToMeteringMode(fmt.Errorf("%w: failed to parse credits %s", err, resp.GetCredits()))

		return nil
	}

	r.balance.Set(credits)

	return nil
}

// DeductOpt changes both the functional behavior of the Deduct method. We chose to do DeductOpt because the standard deduction
// in the v2 engine mucked up the metering interface and the Deduct input params. This approach allows specific behavior
// based on the desired deduct operation.
type DeductOpt func(string, *Report) ([]capabilities.SpendLimit, error)

// ByResource returns a DeductOpt that earmarks a specified amount of local universal credit balance for a given spend
// type.
func ByResource(
	spendType, capabilityID string,
	amount decimal.Decimal,
) func(string, *Report) ([]capabilities.SpendLimit, error) {
	return func(ref string, r *Report) ([]capabilities.SpendLimit, error) {
		step := ReportStep{
			CapabilityID:     capabilityID,
			Deduction:        decimal.Zero,
			AggregatedSpends: make(map[string]AggregatedStepDetail),
		}

		defer func() {
			r.steps[ref] = step
		}()

		bal, err := r.balance.ConvertToBalance(spendType, amount)
		if err != nil {
			// Fail open, continue optimistically
			r.switchToMeteringMode(fmt.Errorf("failed to convert to balance [%s]: %w", spendType, err))
		}

		step.Deduction = bal

		// if in metering mode, exit early without modifying local balance
		if r.meteringMode {
			return []capabilities.SpendLimit{}, nil
		}

		return []capabilities.SpendLimit{}, r.balance.Minus(bal)
	}
}

// ByDerivedAvailability returns a DeductOpt that derives the maximum spend limit based on the user spend limit and
// the number of open concurrent call slots.
func ByDerivedAvailability(
	userSpendLimit decimal.NullDecimal,
	openConcurrentCallSlots int,
	info capabilities.CapabilityInfo,
	config *values.Map,
) func(string, *Report) ([]capabilities.SpendLimit, error) {
	return func(ref string, r *Report) ([]capabilities.SpendLimit, error) {
		step := ReportStep{
			CapabilityID:     info.ID,
			Deduction:        decimal.Zero,
			AggregatedSpends: make(map[string]AggregatedStepDetail),
		}

		defer func() {
			r.steps[ref] = step
		}()

		limit, err := r.getMaxSpendForInvocation(userSpendLimit, openConcurrentCallSlots)
		if err != nil {
			return nil, err
		}

		if !limit.Valid {
			return []capabilities.SpendLimit{}, nil
		}

		step.Deduction = limit.Decimal

		// if in metering mode, exit early without modifying local balance
		if r.meteringMode {
			return []capabilities.SpendLimit{}, nil
		}

		return r.creditToSpendingLimits(info, config, limit.Decimal), r.balance.Minus(limit.Decimal)
	}
}

// Deduct earmarks an amount of local universal credit balance. The amount provided is expected to be in native units.
// An option of 0 indicates a max spend should be derived from user limits and concurrent call slots. We expect to only
// set this value once - an error is returned if a step would be overwritten.
func (r *Report) Deduct(ref string, opt DeductOpt) ([]capabilities.SpendLimit, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.reserved {
		return nil, ErrNoReserve
	}

	if opt == nil {
		return nil, ErrDeductOptionRequired
	}

	if _, ok := r.steps[ref]; ok {
		return nil, ErrStepDeductExists
	}

	return opt(ref, r)
}

// Settle handles the actual spends that each node used for a given capability invocation in the engine,
// by returning earmarked local balance to the available to use pool and adding the spend to the metering report.
// The Deduct method must be called before Settle.
// We expect to only set this value once - an error is returned if a step would be overwritten.
func (r *Report) Settle(ref string, metadata capabilities.ResponseMetadata) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.reserved {
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
	for _, nodeDetail := range metadata.Metering {
		resourceSpends[nodeDetail.SpendUnit] = append(resourceSpends[nodeDetail.SpendUnit], ReportStepDetail{
			Peer2PeerID:   nodeDetail.Peer2PeerID,
			SpendValue:    nodeDetail.SpendValue,
			CRESpendValue: decimal.Zero,
		})
	}

	// Aggregate node responses to a single number
	for unit, spendDetails := range resourceSpends {
		aggregated := AggregatedStepDetail{
			SpendUnit:  unit,
			SpendValue: decimal.Zero,
		}

		deciVals := []decimal.Decimal{}
		for idx, detail := range spendDetails {
			value, err := decimal.NewFromString(detail.SpendValue)
			if err != nil {
				r.lggr.Info(fmt.Sprintf("failed to get spend value from %s: %s", detail.SpendValue, err))
				// throw out invalid values for local balance settlement. they will still be included in metering report.
				continue
			}

			if isGasSpendType(unit) {
				// TODO: this decimal shift should be temporary and converted when write capabilities
				// are converted to provide spend as big.Int fixed point values
				// WARNING: 18 is a magic number here and assumes all gas tokens will have the same level of precision
				value = value.Shift(18) // shift to fixed point value
			}

			if val, err := r.balance.ConvertToBalance(unit, value); err == nil {
				resourceSpends[unit][idx].CRESpendValue = val
			}

			deciVals = append(deciVals, value)

			if isGasSpendType(unit) && len(deciVals) > 1 {
				r.switchToMeteringMode(fmt.Errorf("multiple executions for single execution unit [%s]: %w", unit, err))
			}
		}

		// TODO: explicitly ignore RPC_EVM spend types for now -
		// this check causes TestEngine_Metering_ValidBillingClient/billing_type_and_capability_settle_spend_type_mismatch ./core/services/workflows/v2
		// to fail because the capability is returning a spend type that isn't gas or compute
		// This should be removed when we have proper support for non-gas/compute spend types
		if unit == "RPC_EVM" {
			continue
		}

		aggregated.SpendValue = medianSpend(deciVals)
		value := aggregated.SpendValue

		// if N is not set, assume 1
		if metadata.CapDON_N == 0 {
			metadata.CapDON_N = 1
		}

		// TODO: indicate in the registry config that a capability is single execution or not
		// https://smartcontract-it.atlassian.net/browse/CRE-1037
		if !isGasSpendType(unit) {
			value = value.Mul(decimal.NewFromUint64(uint64(metadata.CapDON_N)))
		}

		bal, err := r.balance.ConvertToBalance(unit, value)

		if err != nil {
			r.switchToMeteringMode(fmt.Errorf("attempted to Settle [%s]: %w", unit, err))
		} else {
			aggregated.CRESpendValue = bal
			spentCredits = spentCredits.Add(bal)
		}

		step.AggregatedSpends[unit] = aggregated
	}

	step.Spends = resourceSpends
	step.CapdonN = metadata.CapDON_N
	r.steps[ref] = step

	// if in metering mode, exit early without modifying local balance
	if r.meteringMode {
		return nil
	}

	// Refund the difference between what local balance had been earmarked and the actual spend
	if err := r.balance.Add(step.Deduction.Sub(spentCredits)); err != nil {
		// invariant: capability should not let spend exceed reserve
		r.lggr.Info("invariant: spend exceeded reserve")
	}

	r.balance.AddSpent(spentCredits)

	return nil
}

func labelToInt32(label string) int32 {
	if value, err := strconv.ParseInt(label, 10, 32); err == nil {
		return int32(value)
	}

	return -1
}

func (r *Report) FormatReport() *protoEvents.MeteringReport {
	protoReport := &protoEvents.MeteringReport{
		Steps: map[string]*protoEvents.MeteringReportStep{},
		Metadata: &protoEvents.WorkflowMetadata{
			WorkflowOwner:           r.labels[platform.KeyWorkflowOwner],
			WorkflowName:            r.labels[platform.KeyWorkflowID],
			Version:                 r.labels[platform.KeyWorkflowVersion],
			WorkflowID:              r.labels[platform.KeyWorkflowID],
			WorkflowExecutionID:     r.labels[platform.KeyWorkflowExecutionID],
			DonID:                   labelToInt32(r.labels[platform.KeyDonID]),
			DonF:                    labelToInt32(r.labels[platform.KeyDonF]),
			DonN:                    labelToInt32(r.labels[platform.KeyDonN]),
			P2PID:                   r.labels[platform.KeyP2PID],
			WorkflowRegistryAddress: r.workflowRegistryAddress,
			WorkflowRegistryVersion: "", // TODO: r.workflowRegistryVersion,
			WorkflowRegistryChain:   strconv.FormatUint(r.workflowRegistryChainSelector, 10),
			EngineVersion:           r.engineVersion,
			DonVersion:              "", // TODO: r.donVersion,
			Trigger: &protoEvents.TriggerDetail{
				TriggerID: r.labels[platform.KeyTriggerID],
			},
		},
		MeteringMode: r.meteringMode,
	}

	if r.meteringModeErr != nil {
		protoReport.Message = r.meteringModeErr.Error()
	}

	r.stepRefLookup = []string{}

	for ref, step := range r.steps {
		stepDetails := &protoEvents.MeteringReportStep{}
		nodeDetails := []*protoEvents.MeteringReportNodeDetail{}
		r.stepRefLookup = append(r.stepRefLookup, ref+":"+step.CapabilityID)

		// since map key order is non-deterministic, order the keys to help make tests deterministic
		orderedUnits := make([]string, 0, len(step.Spends))
		for unit := range step.Spends {
			orderedUnits = append(orderedUnits, unit)
		}

		sort.Slice(orderedUnits, func(i, j int) bool {
			return orderedUnits[i] > orderedUnits[j]
		})

		for _, unit := range orderedUnits {
			details := step.Spends[unit]

			for _, detail := range details {
				nodeDetails = append(nodeDetails, &protoEvents.MeteringReportNodeDetail{
					Peer_2PeerId:  detail.Peer2PeerID,
					SpendUnit:     unit,
					SpendValue:    detail.SpendValue,
					SpendValueCre: detail.CRESpendValue.StringFixed(defaultDecimalPrecision),
				})
			}

			if aggregated, ok := step.AggregatedSpends[unit]; ok {
				// TODO: remove the inaccurate aggregated fields in favor of the repeated field
				stepDetails.AggSpendUnit = aggregated.SpendUnit
				stepDetails.AggSpendValue = aggregated.SpendValue.StringFixed(defaultDecimalPrecision)
				stepDetails.AggSpendValueCre = aggregated.CRESpendValue.StringFixed(defaultDecimalPrecision)

				stepDetails.AggSpend = append(stepDetails.AggSpend, &protoEvents.AggregatedSpendDetail{
					SpendUnit:     aggregated.SpendUnit,
					SpendValue:    aggregated.SpendValue.StringFixed(defaultDecimalPrecision),
					SpendValueCre: aggregated.CRESpendValue.StringFixed(defaultDecimalPrecision),
				})
			}
		}

		stepDetails.Nodes = nodeDetails
		stepDetails.CapdonN = step.CapdonN
		protoReport.Steps[ref] = stepDetails
	}

	return protoReport
}

func (r *Report) SendReceipt(ctx context.Context) error {
	if !r.reserved {
		return ErrNoReserve
	}

	if r.client == nil {
		return ErrNoBillingClient
	}

	r.metrics.UpdateWorkflowMeteringModeGauge(ctx, r.meteringMode)

	// TODO: https://smartcontract-it.atlassian.net/browse/CRE-427 more robust check of billing service health

	req := billing.SubmitWorkflowReceiptRequest{
		WorkflowOwner:                 r.labels[platform.KeyWorkflowOwner],
		WorkflowId:                    r.labels[platform.KeyWorkflowID],
		WorkflowExecutionId:           r.labels[platform.KeyWorkflowExecutionID],
		WorkflowRegistryAddress:       r.workflowRegistryAddress,
		WorkflowRegistryChainSelector: r.workflowRegistryChainSelector,
		Metering:                      r.FormatReport(),
		CreditsConsumed:               r.balance.GetSpent().String(),
	}

	resp, err := r.client.SubmitWorkflowReceipt(ctx, &req)
	if err != nil {
		return err
	}

	if resp == nil {
		return ErrReceiptFailed
	}

	return nil
}

func (r *Report) EmitReceipt(ctx context.Context) error {
	if !r.reserved {
		return ErrNoReserve
	}

	rpt := r.FormatReport()

	r.lggr.Debug("Emitting metering report", "report", rpt, "stepRefs", strings.Join(r.stepRefLookup, ","))

	return wfEvents.EmitMeteringReport(ctx, r.labels, rpt)
}

// creditToSpendingLimits returns a slice of spend limits where the amount is applied to the spend types from the
// provided info. Amount should be specified in universal credits and will be converted to spend type credits within
// this function.
func (r *Report) creditToSpendingLimits(
	info capabilities.CapabilityInfo,
	capConfig *values.Map,
	amount decimal.Decimal,
) []capabilities.SpendLimit {
	if r.meteringMode {
		return []capabilities.SpendLimit{}
	}

	// no spend types results in no limits and is not a failure case
	if len(info.SpendTypes) == 0 {
		return []capabilities.SpendLimit{}
	}

	ratios, err := ratiosFromConfig(info, capConfig)
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
			r.switchToMeteringMode(fmt.Errorf("attempted to create spending limits [%s]: %w", spendType, err))

			return []capabilities.SpendLimit{}
		}

		formattedLimit := spendLimit.StringFixed(defaultDecimalPrecision)
		if isGasSpendType(string(spendType)) {
			formattedLimit = spendLimit.StringFixed(0)
		}

		limits = append(limits, capabilities.SpendLimit{SpendType: spendType, Limit: formattedLimit})
	}

	return limits
}

func isGasSpendType(spendType string) bool {
	return strings.HasPrefix(spendType, "GAS.")
}

// getMaxSpendForInvocation returns the amount of credits that can be used based on the minimum between an optionally
// provided max spend by the user or the available credit balance. The available credit balance is determined by
// dividing unearmarked local credit balance by the number of potential concurrent calls.
func (r *Report) getMaxSpendForInvocation(
	userSpendLimit decimal.NullDecimal,
	openConcurrentCallSlots int,
) (decimal.NullDecimal, error) {
	nullCapSpendLimit := decimal.NewNullDecimal(decimal.Zero)
	nullCapSpendLimit.Valid = false

	if openConcurrentCallSlots == 0 {
		// invariant: this should be managed by the consumer (engine)
		return nullCapSpendLimit, ErrNoOpenCalls
	}

	if !r.reserved {
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

func (r *Report) switchToMeteringMode(err error) {
	// add to the reported errors to indicate all errors encountered during metering
	r.meteringModeErr = errors.Join(r.meteringModeErr, err)

	if r.meteringMode {
		return
	}

	// only log a single metering mode switch error. this error should indicate the first metering related error
	// encountered
	r.lggr.Errorf("switching to metering mode: %s", err)

	r.meteringMode = true
}

func toRateCard(resp *billing.GetWorkflowExecutionRatesResponse) (map[string]decimal.Decimal, error) {
	rates := resp.GetRateCards()

	rateCard := map[string]decimal.Decimal{}
	for _, rate := range rates {
		conversionDeci, err := decimal.NewFromString(rate.UnitsPerCredit)
		if err != nil {
			return map[string]decimal.Decimal{}, fmt.Errorf("could not convert unit %s's value %s to decimal", rate.ResourceType, rate.UnitsPerCredit)
		}

		rateCard[rate.ResourceType.String()] = conversionDeci
	}

	// credits per gas are provided in the form of map[chainselector] -> <gasRate>string
	// each entry should be converted to a usable rate card with form of GAS.[chainselector] -> <unitsPerCredit>decimal
	gasCredits := resp.GetGasTokensPerCredit()

	for chainSelector, gasRate := range gasCredits {
		conversionDeci, err := decimal.NewFromString(gasRate)
		if err != nil {
			return map[string]decimal.Decimal{}, fmt.Errorf("could not convert gas rate %d's value %s to decimal", chainSelector, gasRate)
		}

		rateCard[fmt.Sprintf("GAS.%d", chainSelector)] = conversionDeci
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

	// WorkflowRegistryAddress is the address of the workflow registry contract
	workflowRegistryAddress string
	// WorkflowRegistryChainSelector is the chain selector for the workflow registry
	workflowRegistryChainSelector string
	engineVersion                 string
}

// NewReports initializes and returns a new Reports.
func NewReports(
	client BillingClient,
	owner, workflowID string,
	lggr logger.Logger,
	labels map[string]string,
	metrics *monitoring.WorkflowsMetricLabeler,
	workflowRegistryAddress,
	workflowRegistryChainSelector, engineVersion string,
) *Reports {
	valOf := reflect.ValueOf(client)
	if valOf.IsValid() && valOf.IsNil() {
		client = nil
	}

	return &Reports{
		reports: make(map[string]*Report),
		client:  client,
		lggr:    lggr,
		metrics: metrics,

		owner:      owner,
		workflowID: workflowID,
		labelMap:   labels,

		workflowRegistryAddress:       workflowRegistryAddress,
		workflowRegistryChainSelector: workflowRegistryChainSelector,
		engineVersion:                 engineVersion,
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

	report, err := NewReport(ctx, labels, s.lggr, s.client, s.metrics, s.workflowRegistryAddress, s.workflowRegistryChainSelector, s.engineVersion)
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
		multiErr = errors.Join(multiErr, emitErr)
	}

	sendErr := report.SendReceipt(ctx)
	if sendErr != nil {
		s.metrics.IncrementWorkflowMissingMeteringReport(ctx)
		multiErr = errors.Join(multiErr, sendErr)
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
	capConfig *values.Map,
) (map[capabilities.CapabilitySpendType]decimal.Decimal, error) {
	ratios := make(map[capabilities.CapabilitySpendType]decimal.Decimal)

	// if info.SpendTypes has only 1, return ratio 100%
	if len(info.SpendTypes) == 1 {
		ratios[info.SpendTypes[0]] = decimal.NewFromInt(1)

		return ratios, nil
	}

	if capConfig == nil {
		return ratios, fmt.Errorf("%w: spending ratios not set; config is nil", ErrInvalidRatios)
	}

	rawRatiosValue, hasRatios := capConfig.Underlying[RatiosKey]
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
