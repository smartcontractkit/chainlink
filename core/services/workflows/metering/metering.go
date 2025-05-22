package metering

import (
	"errors"
	"sort"
	"sync"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
	Peer2PeerID string
	SpendUnit   SpendUnit
	SpendValue  SpendValue
}

type Report struct {
	balance *balanceStore
	mu      sync.RWMutex
	steps   map[ReportStepRef][]ReportStep
	lggr    logger.Logger
}

func NewReport(lggr logger.Logger) *Report {
	logger := lggr.Named("Metering")
	balanceStore := NewBalanceStore(0, map[string]decimal.Decimal{}, logger)
	return &Report{
		balance: balanceStore,
		steps:   make(map[ReportStepRef][]ReportStep),
		lggr:    logger,
	}
}

func (r *Report) MedianSpend() map[SpendUnit]SpendValue {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := map[SpendUnit][]SpendValue{}
	medians := map[SpendUnit]SpendValue{}

	for _, nodeVals := range r.steps {
		for _, step := range nodeVals {
			vals, ok := values[step.SpendUnit]
			if !ok {
				vals = []SpendValue{}
			}

			values[step.SpendUnit] = append(vals, step.SpendValue)
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

// SetStep sets the recorded spends for a given capability invocation in the engine.
// We expect to only set this value once - an error is returned if a step would be overwritten
func (r *Report) SetStep(ref ReportStepRef, steps []ReportStep) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.steps[ref]; ok {
		return errors.New("step already exists")
	}

	r.steps[ref] = steps

	return nil
}

func (r *Report) Message() *events.MeteringReport {
	protoReport := &events.MeteringReport{
		Steps:    map[string]*events.MeteringReportStep{},
		Metadata: &events.WorkflowMetadata{},
	}

	for key, step := range r.steps {
		nodeDetail := make([]*events.MeteringReportNodeDetail, len(step))

		for idx, nodeVal := range step {
			nodeDetail[idx] = &events.MeteringReportNodeDetail{
				Peer_2PeerId: nodeVal.Peer2PeerID,
				SpendUnit:    nodeVal.SpendUnit.String(),
				SpendValue:   nodeVal.SpendValue.String(),
			}
		}
		protoReport.Steps[key.String()] = &events.MeteringReportStep{
			Nodes: nodeDetail,
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
