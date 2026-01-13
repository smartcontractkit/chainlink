package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	f "github.com/smartcontractkit/chainlink-testing-framework/framework"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "ocr2"}).Logger()

type ResourceLeakCheckerConfig struct {
	PrometheusURL string
}

type ResourceLeakChecker struct {
	PrometheusURL string
	c             *f.PrometheusQueryClient
}

func WithPromURL(url string) func(*ResourceLeakChecker) {
	return func(rlc *ResourceLeakChecker) {
		rlc.PrometheusURL = url
	}
}

func NewResourceLeakChecker(opts ...func(*ResourceLeakChecker)) *ResourceLeakChecker {
	lc := &ResourceLeakChecker{}
	for _, o := range opts {
		o(lc)
	}
	var pc *f.PrometheusQueryClient
	if lc.PrometheusURL != "" {
		pc = f.NewPrometheusQueryClient(lc.PrometheusURL)
	} else {
		pc = f.NewPrometheusQueryClient(f.LocalPrometheusBaseURL)
	}
	lc.c = pc
	return lc
}

type ResourceLeakCheck struct {
	Query string
	Start time.Time
	End   time.Time
}

func (rc *ResourceLeakChecker) measureIncrease(
	c *ResourceLeakCheck,
	currentNodeNum int,
) (float64, error) {
	memStart, err := rc.c.Query(c.Query, c.Start)
	if err != nil {
		return 0, fmt.Errorf("failed to get memory for the test start: %w", err)
	}

	memEnd, err := rc.c.Query(c.Query, c.End)
	if err != nil {
		return 0, fmt.Errorf("failed to get memory for the test end: %w", err)
	}

	memStartVal := memStart.Data.Result[0].Value[1].(string)
	memEndVal := memEnd.Data.Result[0].Value[1].(string)

	memStartValFloat, err := strconv.ParseFloat(memStartVal, 64)
	if err != nil {
		return 0, fmt.Errorf("start quantile can't be parsed from string: %w", err)
	}
	memEndValFloat, err := strconv.ParseFloat(memEndVal, 64)
	if err != nil {
		return 0, fmt.Errorf("start quantile can't be parsed from string: %w", err)
	}

	totalIncreasePercentage := (memEndValFloat / memStartValFloat * 100) - 100

	L.Debug().
		Int("NodeNum", currentNodeNum).
		Float64("Start", memStartValFloat).
		Float64("End", memEndValFloat).
		Float64("Increase", totalIncreasePercentage).
		Msg("Memory increase total (percentage)")
	return totalIncreasePercentage, nil
}

func main() {
	memoryLeaks := make([]float64, 0)
	lc := NewResourceLeakChecker(WithPromURL(f.LocalPrometheusBaseURL))
	for i := range 4 {
		c := &ResourceLeakCheck{
			Query: fmt.Sprintf(`quantile_over_time(0.5, container_memory_rss{name="don-node%d"}[1h]) / 1024 / 1024`, i),
			Start: time.Now().Add(-2 * time.Hour),
			End:   time.Now(),
		}
		totalDiff, err := lc.measureIncrease(c, i)
		if err != nil {
			panic(err)
		}
		memoryLeaks = append(memoryLeaks, totalDiff)
	}
	fmt.Println(memoryLeaks)
}
