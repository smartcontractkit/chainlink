package workflows

import (
	"sort"
	"sync"

	"github.com/shopspring/decimal"
)

type MeteringReportStepRef string

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

type MeteringReportStep struct {
	Peer2PeerID string
	SpendUnit   MeteringSpendUnit
	SpendValue  MeteringSpendValue
}

type MeteringReport struct {
	mu    sync.RWMutex
	steps map[MeteringReportStepRef]MeteringReportStep
}

func NewMeteringReport() *MeteringReport {
	return &MeteringReport{
		steps: make(map[MeteringReportStepRef]MeteringReportStep),
	}
}

func (r *MeteringReport) MedianSpend() map[MeteringSpendUnit]MeteringSpendValue {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := map[MeteringSpendUnit][]MeteringSpendValue{}
	medians := map[MeteringSpendUnit]MeteringSpendValue{}

	for _, step := range r.steps {
		vals, ok := values[step.SpendUnit]
		if !ok {
			vals = []MeteringSpendValue{}
		}

		values[step.SpendUnit] = append(vals, step.SpendValue)
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

func (r *MeteringReport) AddStep(ref MeteringReportStepRef, step MeteringReportStep) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.steps[ref] = step

	return nil
}
