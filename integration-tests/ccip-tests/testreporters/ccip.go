package testreporters

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/slack-go/slack"

	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
)

type Phase string
type Status string

const (
	E2E                Phase  = "CommitAndExecute"
	TX                 Phase  = "CCIP-Send Transaction"
	CCIPSendRe         Phase  = "CCIPSendRequested"
	SourceLogFinalized Phase  = "SourceLogFinalizedTentatively"
	Commit             Phase  = "Commit-ReportAccepted"
	ExecStateChanged   Phase  = "ExecutionStateChanged"
	ReportBlessed      Phase  = "ReportBlessedByARM"
	Success            Status = "✅"
	Failure            Status = "❌"
	slackFile          string = "payload_ccip.json"
)

type AggregatorMetrics struct {
	Min   float64 `json:"min_duration_for_successful_requests(s),omitempty"`
	Max   float64 `json:"max_duration_for_successful_requests(s),omitempty"`
	Avg   float64 `json:"avg_duration_for_successful_requests(s),omitempty"`
	sum   float64
	count int
}
type TransactionStats struct {
	Fee                string `json:"fee,omitempty"`
	GasUsed            uint64 `json:"gas_used,omitempty"`
	TxHash             string `json:"tx_hash,omitempty"`
	NoOfTokensSent     int    `json:"no_of_tokens_sent,omitempty"`
	MessageBytesLength int    `json:"message_bytes_length,omitempty"`
	FinalizedByBlock   string `json:"finalized_block_num,omitempty"`
	FinalizedAt        string `json:"finalized_at,omitempty"`
	CommitRoot         string `json:"commit_root,omitempty"`
}

type PhaseStat struct {
	SeqNum               uint64           `json:"seq_num,omitempty"`
	Duration             float64          `json:"duration,omitempty"`
	Status               Status           `json:"success"`
	SendTransactionStats TransactionStats `json:"ccip_send_data,omitempty"`
}

type RequestStat struct {
	ReqNo         int64
	SeqNum        uint64
	SourceNetwork string
	DestNetwork   string
	StatusByPhase map[Phase]PhaseStat `json:"status_by_phase,omitempty"`
}

func (stat *RequestStat) UpdateState(lggr zerolog.Logger, seqNum uint64, step Phase, duration time.Duration, state Status, sendTransactionStats ...TransactionStats) {
	durationInSec := duration.Seconds()
	stat.SeqNum = seqNum
	phaseDetails := PhaseStat{
		SeqNum:   seqNum,
		Duration: durationInSec,
		Status:   state,
	}
	if len(sendTransactionStats) > 0 {
		phaseDetails.SendTransactionStats = sendTransactionStats[0]
	}
	stat.StatusByPhase[step] = phaseDetails
	event := lggr.Info()
	if seqNum != 0 {
		event.Uint64("seq num", seqNum)
	}
	// if any of the phase fails mark the E2E as failed
	if state == Failure {
		stat.StatusByPhase[E2E] = PhaseStat{
			SeqNum: seqNum,
			Status: state,
		}
		lggr.Info().
			Str(fmt.Sprint(E2E), string(Failure)).
			Msgf("reqNo %d", stat.ReqNo)
		event.Str(string(step), string(Failure)).Msgf("reqNo %d", stat.ReqNo)
	} else {
		event.Str(string(step), string(Success)).Msgf("reqNo %d", stat.ReqNo)
		if step == Commit || step == ReportBlessed || step == ExecStateChanged {
			stat.StatusByPhase[E2E] = PhaseStat{
				SeqNum:   seqNum,
				Status:   state,
				Duration: stat.StatusByPhase[step].Duration + stat.StatusByPhase[E2E].Duration,
			}
			if step == ExecStateChanged {
				lggr.Info().
					Str(fmt.Sprint(E2E), string(Success)).
					Msgf("reqNo %d", stat.ReqNo)
			}
		}
	}
}

func NewCCIPRequestStats(reqNo int64, source, dest string) *RequestStat {
	return &RequestStat{
		ReqNo:         reqNo,
		StatusByPhase: make(map[Phase]PhaseStat),
		SourceNetwork: source,
		DestNetwork:   dest,
	}
}

type CCIPLaneStats struct {
	lane                    string
	lggr                    zerolog.Logger
	TotalRequests           int64                       `json:"total_requests,omitempty"`          // TotalRequests is the total number of requests made
	SuccessCountsByPhase    map[Phase]int64             `json:"success_counts_by_phase,omitempty"` // SuccessCountsByPhase is the number of requests that succeeded in each phase
	FailedCountsByPhase     map[Phase]int64             `json:"failed_counts_by_phase,omitempty"`  // FailedCountsByPhase is the number of requests that failed in each phase
	DurationStatByPhase     map[Phase]AggregatorMetrics `json:"duration_stat_by_phase,omitempty"`  // DurationStatByPhase is the duration statistics for each phase
	statusByPhaseByRequests sync.Map
}

func (testStats *CCIPLaneStats) UpdatePhaseStatsForReq(stat *RequestStat) {
	testStats.statusByPhaseByRequests.Store(stat.ReqNo, stat.StatusByPhase)
}

func (testStats *CCIPLaneStats) Aggregate(phase Phase, durationInSec float64) {
	if prevDur, ok := testStats.DurationStatByPhase[phase]; !ok {
		testStats.DurationStatByPhase[phase] = AggregatorMetrics{
			Min:   durationInSec,
			Max:   durationInSec,
			sum:   durationInSec,
			count: 1,
		}
	} else {
		if prevDur.Min > durationInSec {
			prevDur.Min = durationInSec
		}
		if prevDur.Max < durationInSec {
			prevDur.Max = durationInSec
		}
		prevDur.sum = prevDur.sum + durationInSec
		prevDur.count++
		testStats.DurationStatByPhase[phase] = prevDur
	}
}

func (testStats *CCIPLaneStats) Finalize(lane string) {
	phases := []Phase{E2E, TX, CCIPSendRe, SourceLogFinalized, Commit, ReportBlessed, ExecStateChanged}
	events := make(map[Phase]*zerolog.Event)
	testStats.statusByPhaseByRequests.Range(func(key, value interface{}) bool {
		if reqNo, ok := key.(int64); ok {
			if stat, ok := value.(map[Phase]PhaseStat); ok {
				for phase, phaseStat := range stat {
					if phaseStat.Status == Success {
						testStats.SuccessCountsByPhase[phase]++
						testStats.Aggregate(phase, phaseStat.Duration)
					} else {
						testStats.FailedCountsByPhase[phase]++
					}
				}
			}
			if reqNo > testStats.TotalRequests {
				testStats.TotalRequests = reqNo
			}
		}
		return true
	})

	testStats.lggr.Info().Int64("Total Requests Triggerred", testStats.TotalRequests).Msg("Test Run Completed")
	for _, phase := range phases {
		events[phase] = testStats.lggr.Info().Str("Phase", string(phase))
		if phaseStat, ok := testStats.DurationStatByPhase[phase]; ok {
			testStats.DurationStatByPhase[phase] = AggregatorMetrics{
				Min: phaseStat.Min,
				Max: phaseStat.Max,
				Avg: phaseStat.sum / float64(phaseStat.count),
			}
			events[phase].
				Str("Min Duration for Successful Requests", fmt.Sprintf("%.02f", testStats.DurationStatByPhase[phase].Min)).
				Str("Max Duration for Successful Requests", fmt.Sprintf("%.02f", testStats.DurationStatByPhase[phase].Max)).
				Str("Average Duration for Successful Requests", fmt.Sprintf("%.02f", testStats.DurationStatByPhase[phase].Avg))
		}
		if failed, ok := testStats.FailedCountsByPhase[phase]; ok {
			events[phase].Int64("Failed Count", failed)
		}
		if s, ok := testStats.SuccessCountsByPhase[phase]; ok {
			events[phase].Int64("Successful Count", s)
		}
		events[phase].Msgf("Phase Stats for Lane %s", lane)
	}
}

type CCIPTestReporter struct {
	t               *testing.T
	logger          zerolog.Logger
	namespace       string
	reportFilePath  string
	duration        time.Duration             // duration is the duration of the test
	FailedLanes     map[string]Phase          `json:"failed_lanes_and_phases,omitempty"` // FailedLanes is the list of lanes that failed and the phase at which it failed
	LaneStats       map[string]*CCIPLaneStats `json:"lane_stats"`                        // LaneStats is the statistics for each lane
	mu              *sync.Mutex
	sendSlackReport bool
}

func (r *CCIPTestReporter) SetSendSlackReport(sendSlackReport bool) {
	r.sendSlackReport = sendSlackReport
}

func (r *CCIPTestReporter) SendSlackNotification(t *testing.T, slackClient *slack.Client, _ testreporters.GrafanaURLProvider) error {
	if r.sendSlackReport {
		r.logger.Info().Msg("Sending Slack notification")
	} else {
		r.logger.Info().Msg("Slack notification not enabled")
		return nil
	}
	if testreporters.SlackAPIKey == "" || testreporters.SlackChannel == "" || testreporters.SlackUserID == "" {
		r.logger.Warn().Msg("Slack API Key, Channel or User ID not set. Skipping Slack notification")
		return nil
	}
	if slackClient == nil {
		slackClient = slack.New(testreporters.SlackAPIKey)
	}

	var msgTexts []string
	headerText := ":white_check_mark: CCIP Test PASSED :white_check_mark:"
	if t.Failed() {
		headerText = ":x: CCIP Test FAILED :x:"
	}
	for name, lane := range r.LaneStats {
		if lane.FailedCountsByPhase[E2E] > 0 {
			msgTexts = append(msgTexts,
				fmt.Sprintf("lane %s :x:", name),
				fmt.Sprintf(
					"Run Duration = %.0fm "+
						"\nNumber of ccip-send= %d"+
						"\nNo of failed requests = %d", r.duration.Minutes(), lane.TotalRequests, lane.FailedCountsByPhase[E2E]))
		}
	}

	msgTexts = append(msgTexts, fmt.Sprintf(
		"\nTest Run Summary created on _remote-test-runner_ at _%s_\nNotifying <@%s>",
		r.reportFilePath, testreporters.SlackUserID))
	if r.namespace == "" {
		r.SetNamespace("ccip")
	}
	messageBlocks := testreporters.SlackNotifyBlocks(headerText, r.namespace, msgTexts)
	ts, err := testreporters.SendSlackMessage(slackClient, slack.MsgOptionBlocks(messageBlocks...))
	if err != nil {
		return fmt.Errorf("failed to send slack message: %w messageBlocks = %v", err, messageBlocks)
	}

	return testreporters.UploadSlackFile(slackClient, slack.FileUploadParameters{
		Title:           fmt.Sprintf("CCIP Test Report %s", r.namespace),
		Filetype:        "json",
		Filename:        fmt.Sprintf("ccip_report_%s.csv", r.namespace),
		File:            r.reportFilePath,
		InitialComment:  fmt.Sprintf("CCIP Test Report %s.", r.namespace),
		Channels:        []string{testreporters.SlackChannel},
		ThreadTimestamp: ts,
	})
}

func (r *CCIPTestReporter) WriteReport(folderPath string) error {
	l := r.logger
	l.Debug().Str("Folder Path", folderPath).Msg("Writing CCIP Test Report")
	if err := testreporters.MkdirIfNotExists(folderPath); err != nil {
		return err
	}
	reportLocation := filepath.Join(folderPath, slackFile)
	r.reportFilePath = reportLocation
	slackFile, err := os.Create(reportLocation)
	defer func() {
		err = slackFile.Close()
		if err != nil {
			l.Error().Err(err).Msg("Error closing slack file")
		}
	}()
	if err != nil {
		return err
	}
	for k := range r.LaneStats {
		r.LaneStats[k].Finalize(k)
		// if E2E for the lane has failed
		if _, ok := r.LaneStats[k].FailedCountsByPhase[E2E]; ok {
			// find the phase at which it failed
			for phase, count := range r.LaneStats[k].FailedCountsByPhase {
				if count > 0 && phase != E2E {
					r.FailedLanes[k] = phase
					break
				}
			}
		}
	}
	if len(r.FailedLanes) > 0 {
		r.logger.Info().Interface("List of Failed Lanes", r.FailedLanes).Msg("Failed Lanes")
	} else {
		r.logger.Info().Msg("All Lanes Passed")
	}
	stats, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	_, err = slackFile.Write(stats)
	if err != nil {
		return err
	}
	return nil
}

// SetNamespace sets the namespace of the report for clean reports
func (r *CCIPTestReporter) SetNamespace(namespace string) {
	r.namespace = namespace
}

// SetDuration sets the duration of the test
func (r *CCIPTestReporter) SetDuration(d time.Duration) {
	r.duration = d
}

func (r *CCIPTestReporter) AddNewLane(name string, lggr zerolog.Logger) *CCIPLaneStats {
	r.mu.Lock()
	defer r.mu.Unlock()
	i := &CCIPLaneStats{
		lane:                 name,
		lggr:                 lggr,
		FailedCountsByPhase:  make(map[Phase]int64),
		SuccessCountsByPhase: make(map[Phase]int64),
		DurationStatByPhase:  make(map[Phase]AggregatorMetrics),
	}
	r.LaneStats[name] = i
	return i
}

func (r *CCIPTestReporter) SendReport(t *testing.T, namespace string, slackSend bool) error {
	logsPath := filepath.Join("logs", fmt.Sprintf("%s-%s-%d", t.Name(), namespace, time.Now().Unix()))
	r.SetNamespace(namespace)
	r.SetSendSlackReport(namespace != "" && slackSend)
	return testreporters.SendReport(t, namespace, logsPath, r, nil)
}

func NewCCIPTestReporter(t *testing.T, lggr zerolog.Logger) *CCIPTestReporter {
	return &CCIPTestReporter{
		LaneStats:   make(map[string]*CCIPLaneStats),
		logger:      lggr,
		t:           t,
		mu:          &sync.Mutex{},
		FailedLanes: make(map[string]Phase),
	}
}
