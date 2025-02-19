package ccip

import (
	"context"
	"strconv"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"testing"
	"time"
)

const (
	LokiLoadLabel = "ccipv2_load_test"
	ErrLokiPush   = "failed to push metrics to Loki"
)

// MetricsManager is used for maintaining state of different sequence numbers
// Once we've received all expected timestamps, it pushes the metrics to Loki
type MetricManager struct {
	lggr      logger.Logger
	loki      *wasp.LokiClient
	InputChan chan messageData
	state     map[srcDstSeqNum]metricState
}

type metricState struct {
	timestamps [3]uint64
	round      int
}

type srcDstSeqNum struct {
	src    uint64
	dst    uint64
	seqNum uint64
}

type messageData struct {
	eventType int
	srcDstSeqNum
	timestamp uint64
	round     int
}

func NewMetricsManager(t *testing.T, l logger.Logger) *MetricManager {
	// initialize loki using endpoint from user defined env vars
	loki, err := wasp.NewLokiClient(wasp.NewEnvLokiConfig())
	require.NoError(t, err)

	return &MetricManager{
		lggr:      l,
		loki:      loki,
		InputChan: make(chan messageData),
		state:     make(map[srcDstSeqNum]metricState),
	}
}

func (mm *MetricManager) Stop() {
	close(mm.InputChan)
}

func (mm *MetricManager) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			mm.lggr.Infow("received timeout, pushing remaining state to loki")
			// any remaining data in state should be pushed to loki as incomplete
			for srcDstSeqNum, metricState := range mm.state {
				commitDuration, execDuration := uint64(0), uint64(0)
				timestamps := metricState.timestamps
				if timestamps[committed] != 0 && timestamps[transmitted] != 0 {
					commitDuration = timestamps[committed] - timestamps[transmitted]
				}
				if timestamps[executed] != 0 && timestamps[committed] != 0 {
					execDuration = timestamps[executed] - timestamps[committed]
				}

				lokiLabels, err := setLokiLabels(srcDstSeqNum.src, srcDstSeqNum.dst, metricState.round)
				if err != nil {
					mm.lggr.Error("error setting loki labels", "error", err)
					// don't return here, we still want to push metrics to loki
				}
				SendMetricsToLoki(mm.lggr, mm.loki, lokiLabels, &LokiMetric{
					ExecDuration:   execDuration,
					CommitDuration: commitDuration,
					SequenceNumber: srcDstSeqNum.seqNum,
				})
			}
			close(mm.InputChan)
			mm.loki.Stop()
			return
		case data := <-mm.InputChan:
			if _, ok := mm.state[data.srcDstSeqNum]; !ok {
				mm.state[data.srcDstSeqNum] = metricState{
					timestamps: [3]uint64{0, 0, 0},
				}
			}

			state := mm.state[data.srcDstSeqNum]
			state.timestamps[data.eventType] = data.timestamp
			if data.eventType == transmitted && data.round != -1 {
				state.round = data.round
			}
			mm.state[data.srcDstSeqNum] = state
			if data.eventType == executed {
				mm.lggr.Infow("new state for received seqNum is ", "dst", data.dst, "seqNum", data.seqNum, "round", state.round, "timestamps", state.timestamps)
			}
			// we have all data needed to push to Loki
			if state.timestamps[transmitted] != 0 && state.timestamps[committed] != 0 && state.timestamps[executed] != 0 {
				lokiLabels, err := setLokiLabels(data.src, data.dst, mm.state[data.srcDstSeqNum].round)
				if err != nil {
					mm.lggr.Error("error setting loki labels", "error", err)
				}
				SendMetricsToLoki(mm.lggr, mm.loki, lokiLabels, &LokiMetric{
					ExecDuration:   state.timestamps[executed] - state.timestamps[committed],
					CommitDuration: state.timestamps[committed] - state.timestamps[transmitted],
					SequenceNumber: data.seqNum,
				})

				delete(mm.state, data.srcDstSeqNum)
			}
		}
	}
}

func SendMetricsToLoki(l logger.Logger, lc *wasp.LokiClient, updatedLabels map[string]string, metrics *LokiMetric) {
	if err := lc.HandleStruct(wasp.LabelsMapToModel(updatedLabels), time.Now(), metrics); err != nil {
		l.Error(ErrLokiPush)
	}
}

func setLokiLabels(src, dst uint64, round int) (map[string]string, error) {
	srcChainID, err := chainselectors.GetChainIDFromSelector(src)
	if err != nil {
		return nil, err
	}
	dstChainID, err := chainselectors.GetChainIDFromSelector(dst)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"sourceEvmChainId": srcChainID,
		"destEvmChainId":   dstChainID,
		"roundNum":         strconv.Itoa(round),
		"testType":         LokiLoadLabel,
	}, nil
}
