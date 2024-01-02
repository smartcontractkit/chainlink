package smoke

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	logpoller "github.com/smartcontractkit/chainlink/integration-tests/universal/log_poller"
)

// consistency test with no network disruptions with approximate emission of 1500-1600 logs per second for ~110-120 seconds
// 6 filters are registered
func TestLogPollerFewFiltersFixedDepth(t *testing.T) {
	executeLogPollerBasicTest(t)
}

func TestLogPollerFewFiltersFinalityTag(t *testing.T) {
	executeLogPollerBasicTest(t)
}

// consistency test with no network disruptions with approximate emission of 1000-1100 logs per second for ~110-120 seconds
// 900 filters are registered
func TestLogManyFiltersPollerFixedDepth(t *testing.T) {
	executeLogPollerBasicTest(t)
}

func TestLogManyFiltersPollerFinalityTag(t *testing.T) {
	executeLogPollerBasicTest(t)
}

// consistency test that introduces random distruptions by pausing either Chainlink or Postgres containers for random interval of 5-20 seconds
// with approximate emission of 520-550 logs per second for ~110 seconds
// 6 filters are registered
func TestLogPollerWithChaosFixedDepth(t *testing.T) {
	executeLogPollerBasicTest(t)
}

func TestLogPollerWithChaosFinalityTag(t *testing.T) {
	executeLogPollerBasicTest(t)
}

// consistency test that registers filters after events were emitted and then triggers replay via API
// unfortunately there is no way to make sure that logs that are indexed are only picked up by replay
// and not by backup poller
// with approximate emission of 24 logs per second for ~110 seconds
// 6 filters are registered
func TestLogPollerReplayFixedDepth(t *testing.T) {
	executeLogPollerReplayTest(t, "5m")
}

func TestLogPollerReplayFinalityTag(t *testing.T) {
	executeLogPollerReplayTest(t, "5m")
}

func executeLogPollerBasicTest(t *testing.T) {
	eventsToEmit := []abi.Event{}
	for _, event := range logpoller.EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg, err := tc.GetConfig(t.Name(), tc.LogPoller)
	if err != nil {
		t.Fatal(err)
	}
	cfg.LogPoller.General.EventsToEmit = eventsToEmit

	logpoller.ExecuteBasicLogPollerTest(t, &cfg)
}

func executeLogPollerReplayTest(t *testing.T, duration string) {
	eventsToEmit := []abi.Event{}
	for _, event := range logpoller.EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg, err := tc.GetConfig(t.Name(), tc.LogPoller)
	if err != nil {
		t.Fatal(err)
	}
	cfg.LogPoller.General.EventsToEmit = eventsToEmit

	logpoller.ExecuteLogPollerReplay(t, &cfg, duration)
}
