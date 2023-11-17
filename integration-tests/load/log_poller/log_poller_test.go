package logpoller

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/stretchr/testify/require"

	lp_helpers "github.com/smartcontractkit/chainlink/integration-tests/universal/log_poller"
)

func TestLoadTestLogPoller(t *testing.T) {
	cfg, err := lp_helpers.ReadConfig(lp_helpers.DefaultConfigFilename)
	require.NoError(t, err)

	eventsToEmit := []abi.Event{}
	for _, event := range lp_helpers.EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg.General.EventsToEmit = eventsToEmit

	lp_helpers.ExecuteBasicLogPollerTest(t, cfg)
}
