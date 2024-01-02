package logpoller

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/stretchr/testify/require"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	lp_helpers "github.com/smartcontractkit/chainlink/integration-tests/universal/log_poller"
)

func TestLoadTestLogPoller(t *testing.T) {
	config, err := tc.GetConfig("Load", tc.LogPoller)
	require.NoError(t, err)

	eventsToEmit := []abi.Event{}
	for _, event := range lp_helpers.EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	config.LogPoller.General.EventsToEmit = eventsToEmit

	lp_helpers.ExecuteBasicLogPollerTest(t, &config)
}
