package changeset

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployCCIPContracts(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := NewMemoryEnvironmentWithJobsAndContracts(t, lggr, 2, 4, nil)

	// Deploy all the CCIP contracts.
	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)
	snap, err := state.View(e.Env.AllChainSelectors())
	require.NoError(t, err)

	// Assert expect every deployed address to be in the address book.
	// TODO (CCIP-3047): Add the rest of CCIPv2 representation
	b, err := json.MarshalIndent(snap, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
	startBlocks := make(map[uint64]*uint64)
	for key, val := range e.ReplayBlocks {
		startBlocks[key] = &val
	}
	// wait for jobs to start
	time.Sleep(40 * time.Second)
	ConfirmTokenPriceUpdatedForAll(t, e.Env, state, startBlocks,
		DefaultInitialPrices.LinkPrice, DefaultInitialPrices.WethPrice)
}
