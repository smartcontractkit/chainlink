package performance

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"time"
)

// Test is the interface to be implemented for performance tests
type Test interface {
	Setup() error
	Run() error
	RecordValues(b ginkgo.Benchmarker) error
}

type TestOptions struct {
	NumberOfContracts int
	NumberOfRounds    int64
}

// PercentileReport common percentile report
type PercentileReport struct {
	StdDev float64
	Max    float64
	Min    float64
	P99    float64
	P95    float64
	P90    float64
	P50    float64
}

// NewPercentileReport calculates percentiles for arbitrary float64 data
func NewPercentileReport(data []time.Duration) (*PercentileReport, error) {
	dataFloat64 := make([]float64, 0)
	for _, d := range data {
		dataFloat64 = append(dataFloat64, d.Seconds())
	}
	perc99, err := stats.Percentile(dataFloat64, 99)
	if err != nil {
		return nil, err
	}
	perc95, err := stats.Percentile(dataFloat64, 95)
	if err != nil {
		return nil, err
	}
	perc90, err := stats.Percentile(dataFloat64, 90)
	if err != nil {
		return nil, err
	}
	perc50, err := stats.Percentile(dataFloat64, 50)
	if err != nil {
		return nil, err
	}
	max, err := stats.Max(dataFloat64)
	if err != nil {
		return nil, err
	}
	min, err := stats.Min(dataFloat64)
	if err != nil {
		return nil, err
	}
	stdDev, err := stats.StandardDeviation(dataFloat64)
	if err != nil {
		return nil, err
	}
	return &PercentileReport{P99: perc99, P95: perc95, P90: perc90, P50: perc50, Max: max, Min: min, StdDev: stdDev}, nil
}

// PrintPercentileMetrics prints percentile metrics
func (m *PercentileReport) PrintPercentileMetrics() {
	log.Info().Float64("Latency", m.Max).Msg("Maximum")
	log.Info().Float64("Latency", m.P99).Msg("99th Percentile")
	log.Info().Float64("Latency", m.P95).Msg("95th Percentile")
	log.Info().Float64("Latency", m.P90).Msg("90th Percentile")
	log.Info().Float64("Latency", m.P50).Msg("50th Percentile")
	log.Info().Float64("Latency", m.Min).Msg("Minimum")
	log.Info().Float64("Latency", m.StdDev).Msg("Standard Deviation")
}

func recordResults(b ginkgo.Benchmarker, ID string, results []time.Duration) error {
	percentileReport, err := NewPercentileReport(results)
	if err != nil {
		return err
	}
	percentileReport.PrintPercentileMetrics()

	for _, result := range results {
		b.RecordValue(ID, result.Seconds())
	}
	b.RecordValue(fmt.Sprintf("%s_P50", ID), percentileReport.P50)
	b.RecordValue(fmt.Sprintf("%s_P90", ID), percentileReport.P90)
	b.RecordValue(fmt.Sprintf("%s_P95", ID), percentileReport.P95)
	b.RecordValue(fmt.Sprintf("%s_P99", ID), percentileReport.P99)
	b.RecordValue(fmt.Sprintf("%s_Max", ID), percentileReport.Max)

	return nil
}
