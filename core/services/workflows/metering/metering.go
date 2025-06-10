package metering

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/workflows/go/events"
)

type ReportStepRef string

func (s ReportStepRef) String() string {
	return string(s)
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

type ProtoDetail struct {
	Schema string
	Domain string
	Entity string
}

type ReportStep struct {
	Reserve map[SpendUnit]SpendValue
	Spend   map[SpendUnit][]ReportStepDetail
}

type ReportStepDetail struct {
	Peer2PeerID string
	SpendValue  SpendValue
}

type Report struct {
	balance *balanceStore
	mu      sync.RWMutex
	steps   map[ReportStepRef]ReportStep
	lggr    logger.Logger
}

func NewReport(lggr logger.Logger) *Report {
	logger := logger.Named(lggr, "Metering")
	balanceStore := NewBalanceStore(0, map[string]decimal.Decimal{}, logger)
	return &Report{
		balance: balanceStore,
		steps:   make(map[ReportStepRef]ReportStep),
		lggr:    logger,
	}
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

// ReserveStep earmarks the maximum spend for a given capability invocation in the engine.
// We expect to only set this value once - an error is returned if a step would be overwritten
func (r *Report) ReserveStep(ref ReportStepRef, capInfo capabilities.CapabilityInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.steps[ref]; ok {
		return errors.New("step reserve already exists")
	}

	// TODO: handle extra reserves for write step

	r.steps[ref] = ReportStep{
		Reserve: make(map[SpendUnit]SpendValue),
		Spend:   nil,
	}

	return nil
}

// SetStep sets the recorded spends for a given capability invocation in the engine.
// ReserveStep must be called before SetStep
// We expect to only set this value once - an error is returned if a step would be overwritten
func (r *Report) SetStep(ref ReportStepRef, steps []capabilities.MeteringNodeDetail) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	step, ok := r.steps[ref]
	if !ok {
		return errors.New("must call Report.ReserveStep first")
	}

	if step.Spend != nil {
		return errors.New("step spend already exists")
	}

	spend := make(map[SpendUnit][]ReportStepDetail)

	for _, detail := range steps {
		unit := SpendUnit(detail.SpendUnit)
		value, err := unit.StringToSpendValue(detail.SpendValue)
		if err != nil {
			r.lggr.Error(fmt.Sprintf("failed to get spend value from %s: %s", detail.SpendValue, err))
		}
		spend[unit] = append(spend[unit], ReportStepDetail{
			Peer2PeerID: detail.Peer2PeerID,
			SpendValue:  value,
		})
	}

	step.Spend = spend
	r.steps[ref] = step

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

		protoReport.Steps[key.String()] = &events.MeteringReportStep{
			Nodes: nodeDetails,
		}
	}

	return protoReport
}

// Reports is a concurrency-safe wrapper around map[string]*Report.
type Reports struct {
	mu      sync.RWMutex
	reports map[string]*Report
}

// NewReports initializes and returns a new Reports.
func NewReports() *Reports {
	return &Reports{
		reports: make(map[string]*Report),
	}
}

// Get retrieves a Report for a given key (if it exists).
func (s *Reports) Get(key string) (*Report, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.reports[key]
	return val, ok
}

// Add inserts or updates a Report under the specified key.
func (s *Reports) Add(key string, report *Report) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.reports[key] = report
}

// Delete removes the Report with the specified key.
func (s *Reports) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.reports, key)
}

func (s *Reports) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.reports)
}
