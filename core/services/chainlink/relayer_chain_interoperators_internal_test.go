package chainlink

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// Readers must not iterate loopRelayers while Get lazily inserts dummy relayers.
// Run with -race: an unlocked read shows up as a data race or a concurrent map fault.
func TestCoreRelayerChainInteroperators_ReadersDoNotRaceLazyGet(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	cr, err := NewCoreRelayerChainInteroperators(InitDummy(RelayerFactory{Logger: logger.Test(t)}))
	require.NoError(t, err)

	const n = 50
	var wg sync.WaitGroup
	for i := range n {
		wg.Add(3)
		id := types.RelayID{Network: relay.NetworkDummy, ChainID: strconv.Itoa(i)}
		go func() {
			defer wg.Done()
			_, getErr := cr.Get(id)
			assert.NoError(t, getErr)
		}()
		go func() {
			defer wg.Done()
			_, _, statusErr := cr.NodeStatuses(ctx, 0, 0)
			assert.NoError(t, statusErr)
		}()
		go func() {
			defer wg.Done()
			_ = cr.Slice()
		}()
	}
	wg.Wait()

	require.Len(t, cr.Slice(), n)
	_, _, err = cr.NodeStatuses(ctx, 0, 0, types.RelayID{Network: relay.NetworkDummy, ChainID: "0"})
	require.NoError(t, err)
}
