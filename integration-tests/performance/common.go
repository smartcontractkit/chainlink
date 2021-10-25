package performance

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"golang.org/x/sync/errgroup"

	"github.com/montanaflynn/stats"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
)

// Test is the interface to be implemented for performance tests
type Test interface {
	Setup() error
	Run() error
	RecordValues(b ginkgo.Benchmarker) error
}

// TestOptions common perf/soak test options
// either TestDuration can be set or NumberOfRounds, or both
type TestOptions struct {
	NumberOfContracts    int
	NumberOfRounds       int
	RoundTimeout         time.Duration
	TestDuration         time.Duration
	GracefulStopDuration time.Duration
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

// Contract is just a basic contract interface
type Contract interface {
	Address() string
}

// PerfJobRunResult contains the start & end time of the round submission to calculate latency
type PerfJobRunResult struct {
	StartTime time.Time
	EndTime   time.Time
}

// PerfRoundTestResults is a complex map that holds all test data in a map by the round ID, then contract instance and
// then the Chainlink client
type PerfRoundTestResults struct {
	mutex   *sync.Mutex
	results map[int64]map[Contract]map[client.Chainlink]*PerfJobRunResult
}

// NewPerfTestResults returns an instance PerfRoundTestResults
func NewPerfTestResults() PerfRoundTestResults {
	return PerfRoundTestResults{
		mutex:   &sync.Mutex{},
		results: map[int64]map[Contract]map[client.Chainlink]*PerfJobRunResult{},
	}
}

// Get a value from the test results map with nil checking to avoid panics
func (f PerfRoundTestResults) Get(
	roundID int64,
	contract Contract,
	chainlink client.Chainlink,
) *PerfJobRunResult {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if _, ok := f.results[roundID]; !ok {
		f.results[roundID] = map[Contract]map[client.Chainlink]*PerfJobRunResult{}
	}
	if _, ok := f.results[roundID][contract]; !ok {
		f.results[roundID][contract] = map[client.Chainlink]*PerfJobRunResult{}
	}
	if f.results[roundID][contract][chainlink] == nil {
		f.results[roundID][contract][chainlink] = &PerfJobRunResult{}
	}
	return f.results[roundID][contract][chainlink]
}

// GetAll returns the full map, not safe for concurrent actions
func (f PerfRoundTestResults) GetAll() map[int64]map[Contract]map[client.Chainlink]*PerfJobRunResult {
	return f.results
}

// PerfRequestIDTestResults is results traced and aggregated by request id, see models.DecodeLogTaskRun
type PerfRequestIDTestResults struct {
	mutex   *sync.Mutex
	results map[string]*PerfJobRunResult
}

// NewPerfRequestIDTestResults returns an instance NewPerfRequestIDTestResults
func NewPerfRequestIDTestResults() *PerfRequestIDTestResults {
	return &PerfRequestIDTestResults{
		mutex:   &sync.Mutex{},
		results: map[string]*PerfJobRunResult{},
	}
}

// Get a value from the test results map with nil checking to avoid panics
func (r *PerfRequestIDTestResults) Get(requestID string) *PerfJobRunResult {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.results[requestID]; !ok {
		r.results[requestID] = &PerfJobRunResult{}
	}
	return r.results[requestID]
}

// GetAll returns all test results
func (r *PerfRequestIDTestResults) GetAll() map[string]*PerfJobRunResult {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.results
}

func (r *PerfRequestIDTestResults) setResultStartTimes(clients []client.Chainlink, jobMap ContractsNodesJobsMap) error {
	g := errgroup.Group{}
	for contract := range jobMap {
		contract := contract
		g.Go(func() error {
			return r.setResultStartTimeByContract(clients, jobMap, contract)
		})
	}
	return g.Wait()
}

func (r *PerfRequestIDTestResults) setResultStartTimeByContract(clients []client.Chainlink, jobMap ContractsNodesJobsMap, contract interface{}) error {
	for _, chainlink := range clients {
		jobRuns, err := chainlink.ReadRunsByJob(jobMap[contract][chainlink].GetJobID())
		if err != nil {
			return err
		}
		log.Debug().
			Str("Node", chainlink.URL()).
			Int("Runs", len(jobRuns.Data)).
			Msg("Total runs")
		for _, jobDecodeData := range jobRuns.Data {
			rqInts, err := actions.ExtractRequestIDFromJobRun(jobDecodeData)
			if err != nil {
				return err
			}
			rqID := common.Bytes2Hex(rqInts)
			loc, _ := time.LoadLocation("UTC")
			startTime := jobDecodeData.Attributes.CreatedAt.In(loc)
			log.Debug().
				Time("StartTime", startTime).
				Str("RequestID", rqID).
				Msg("Request found")
			d := r.Get(rqID)
			d.StartTime = startTime
		}
	}
	return nil
}

func (r *PerfRequestIDTestResults) calculateLatencies(b ginkgo.Benchmarker) error {
	var latencies []time.Duration
	for rqID, testResult := range r.GetAll() {
		latency := testResult.EndTime.Sub(testResult.StartTime)
		log.Debug().
			Str("RequestID", rqID).
			Time("StartTime", testResult.StartTime).
			Time("EndTime", testResult.EndTime).
			Dur("Duration", latency).
			Msg("Calculating latencies for request id")
		if testResult.StartTime.IsZero() {
			log.Warn().
				Str("RequestID", rqID).
				Msg("Start time zero")
		}
		if testResult.EndTime.IsZero() {
			log.Warn().
				Str("RequestID", rqID).
				Msg("End time zero")
		}
		if latency.Seconds() < 0 {
			log.Warn().
				Str("RequestID", rqID).
				Msg("Latency below zero")
		} else {
			latencies = append(latencies, latency)
		}
	}
	return recordResults(b, "Request latency", latencies)
}

// NodeData common node data
type NodeData interface {
	GetJobID() string
	GetProvingKeyHash() [32]byte
}

// RunlogNodeData node data required for runlog test
type RunlogNodeData struct {
	JobID string
}

// GetJobID gets internal job id
func (n RunlogNodeData) GetJobID() string {
	return n.JobID
}

// GetProvingKeyHash gets proving key hash for VRF
func (n RunlogNodeData) GetProvingKeyHash() [32]byte {
	return [32]byte{}
}

// VRFNodeData VRF node data
type VRFNodeData struct {
	ProvingKeyHash [32]byte
	JobID          string
}

// GetJobID gets internal job id
func (n VRFNodeData) GetJobID() string {
	return n.JobID
}

// GetProvingKeyHash gets proving key hash for VRF
func (n VRFNodeData) GetProvingKeyHash() [32]byte {
	return n.ProvingKeyHash
}

// ContractsNodesJobsMap common contract to node to job id mapping for perf/soak tests
type ContractsNodesJobsMap map[interface{}]map[client.Chainlink]NodeData

// FromJobsChan fills ContractsNodesJobsMap from a chan used in parallel deployment
func (c ContractsNodesJobsMap) FromJobsChan(jobsChan chan ContractsNodesJobsMap) {
	for jobMap := range jobsChan {
		for contractAddr, m := range jobMap {
			if _, ok := c[contractAddr]; !ok {
				c[contractAddr] = map[client.Chainlink]NodeData{}
			}
			for k, v := range m {
				c[contractAddr][k] = v
			}
		}
	}
}

// LimitErrGroup implements the errgroup.Group interface, but limits goroutines to a
// execute at a max throughput.
type LimitErrGroup struct {
	ticker *time.Ticker
	eg     *errgroup.Group
}

// NewLimitErrGroup initializes and returns a new rate limited errgroup
func NewLimitErrGroup(rps int) *LimitErrGroup {
	eg := &errgroup.Group{}
	r := &LimitErrGroup{
		ticker: time.NewTicker(time.Second / time.Duration(rps)),
		eg:     eg,
	}
	return r
}

// Go runs a new job as a goroutine. It will wait until the next available
// time so that the ratelimit is not exceeded.
func (e *LimitErrGroup) Go(fn func() error) {
	<-e.ticker.C
	go e.eg.Go(fn)
}

// Wait will wait until all jobs are processed. Once Wait() is called, no more jobs
// can be added.
func (e *LimitErrGroup) Wait() error {
	defer e.ticker.Stop()
	return e.eg.Wait()
}
