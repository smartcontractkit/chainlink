package workflows

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidMeteringSpendValue = errors.New("invalid metering spend value")
)

type MeteringReportStepRef string

type MeteringSpendUnit string

func (s MeteringSpendUnit) String() string {
	return string(s)
}

func (s MeteringSpendUnit) DecimalToSpendValue(value decimal.Decimal) MeteringSpendValue {
	return MeteringSpendValue(value.String())
}

type MeteringSpendValue string

func (s MeteringSpendValue) String() string {
	return string(s)
}

type MeteringReportStep struct {
	Peer2PeerID string
	SpendUnit   MeteringSpendUnit
	SpendValue  MeteringSpendValue
}

func (s MeteringReportStep) Value() (decimal.Decimal, error) {
	return decimal.NewFromString(s.SpendValue.String())
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

	values := map[MeteringSpendUnit][]decimal.Decimal{}
	medians := map[MeteringSpendUnit]MeteringSpendValue{}

	for _, step := range r.steps {
		vals, ok := values[step.SpendUnit]
		if !ok {
			vals = []decimal.Decimal{}
		}

		// ignoring the error here should be safe as long as AddStep verifies parsing
		value, _ := step.Value()

		values[step.SpendUnit] = append(vals, value)
	}

	for unit, set := range values {
		sort.Slice(set, func(i, j int) bool {
			return set[j].GreaterThan(set[i])
		})

		if len(set)%2 > 0 {
			medians[unit] = unit.DecimalToSpendValue(set[len(set)/2])

			continue
		}

		avg := set[len(set)/2-1].Add(set[len(set)/2]).Div(decimal.NewFromInt(2))
		medians[unit] = unit.DecimalToSpendValue(avg)
	}

	return medians
}

func (r *MeteringReport) AddStep(ref MeteringReportStepRef, step MeteringReportStep) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, err := step.Value(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidMeteringSpendValue, err)
	}

	r.steps[ref] = step

	return nil
}
