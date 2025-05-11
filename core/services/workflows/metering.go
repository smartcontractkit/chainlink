package workflows

import (
	"errors"
	"sort"
	"sync"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-protos/workflows/go/events"
)

type MeteringReportStepRef string

func (s MeteringReportStepRef) String() string {
	return string(s)
}

type MeteringSpendUnit string

func (s MeteringSpendUnit) String() string {
	return string(s)
}

func (s MeteringSpendUnit) DecimalToSpendValue(value decimal.Decimal) MeteringSpendValue {
	return MeteringSpendValue{value: value, roundingPlace: 18}
}

func (s MeteringSpendUnit) IntToSpendValue(value int64) MeteringSpendValue {
	return MeteringSpendValue{value: decimal.NewFromInt(value), roundingPlace: 18}
}

func (s MeteringSpendUnit) StringToSpendValue(value string) (MeteringSpendValue, error) {
	dec, err := decimal.NewFromString(value)
	if err != nil {
		return MeteringSpendValue{}, err
	}

	return MeteringSpendValue{value: dec, roundingPlace: 18}, nil
}

type MeteringSpendValue struct {
	value         decimal.Decimal
	roundingPlace uint8
}

func (v MeteringSpendValue) Add(value MeteringSpendValue) MeteringSpendValue {
	return MeteringSpendValue{
		value:         v.value.Add(value.value),
		roundingPlace: v.roundingPlace,
	}
}

func (v MeteringSpendValue) Div(value MeteringSpendValue) MeteringSpendValue {
	return MeteringSpendValue{
		value:         v.value.Div(value.value),
		roundingPlace: v.roundingPlace,
	}
}

func (v MeteringSpendValue) GreaterThan(value MeteringSpendValue) bool {
	return v.value.GreaterThan(value.value)
}

func (v MeteringSpendValue) String() string {
	return v.value.StringFixedBank(int32(v.roundingPlace))
}

type ProtoDetail struct {
	Schema string
	Domain string
	Entity string
}

type MeteringReportStep struct {
	Peer2PeerID string
	SpendUnit   MeteringSpendUnit
	SpendValue  MeteringSpendValue
}

type MeteringReport struct {
	mu    sync.RWMutex
	steps map[MeteringReportStepRef][]MeteringReportStep
}

func NewMeteringReport() *MeteringReport {
	return &MeteringReport{
		steps: make(map[MeteringReportStepRef][]MeteringReportStep),
	}
}

func (r *MeteringReport) MedianSpend() map[MeteringSpendUnit]MeteringSpendValue {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := map[MeteringSpendUnit][]MeteringSpendValue{}
	medians := map[MeteringSpendUnit]MeteringSpendValue{}

	for _, nodeVals := range r.steps {
		for _, step := range nodeVals {
			vals, ok := values[step.SpendUnit]
			if !ok {
				vals = []MeteringSpendValue{}
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
func (r *MeteringReport) SetStep(ref MeteringReportStepRef, steps []MeteringReportStep) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.steps[ref]; ok {
		return errors.New("step already exists")
	}

	r.steps[ref] = steps

	return nil
}

func (r *MeteringReport) Message() *events.MeteringReport {
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

// MeterReports is a concurrency-safe wrapper around map[string]*MeteringReport.
type MeterReports struct {
	mu           sync.RWMutex
	meterReports map[string]*MeteringReport
}

// NewMeterReports initializes and returns a new MeterReports.
func NewMeterReports() *MeterReports {
	return &MeterReports{
		meterReports: make(map[string]*MeteringReport),
	}
}

// Get retrieves a MeteringReport for a given key (if it exists).
func (s *MeterReports) Get(key string) (*MeteringReport, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.meterReports[key]
	return val, ok
}

// Add inserts or updates a MeteringReport under the specified key.
func (s *MeterReports) Add(key string, report *MeteringReport) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.meterReports[key] = report
}

// Delete removes the MeteringReport with the specified key.
func (s *MeterReports) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.meterReports, key)
}

func (s *MeterReports) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.meterReports)
}
