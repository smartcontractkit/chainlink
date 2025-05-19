package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"
)

func TestCapabilitiesConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	p2p := cfg.Capabilities().Peering()
	assert.Equal(t, "p2p_12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw", p2p.PeerID().String())
	assert.Equal(t, 13, p2p.IncomingMessageBufferSize())
	assert.Equal(t, 17, p2p.OutgoingMessageBufferSize())
	assert.True(t, p2p.TraceLogging())

	v2 := p2p.V2()
	assert.False(t, v2.Enabled())
	assert.Equal(t, []string{"a", "b", "c"}, v2.AnnounceAddresses())
	assert.ElementsMatch(
		t,
		[]commontypes.BootstrapperLocator{
			{
				PeerID: "12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw",
				Addrs:  []string{"test:99"},
			},
			{
				PeerID: "12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw",
				Addrs:  []string{"foo:42", "bar:10"},
			},
		},
		v2.DefaultBootstrappers(),
	)
	assert.Equal(t, time.Minute, v2.DeltaDial().Duration())
	assert.Equal(t, 2*time.Second, v2.DeltaReconcile().Duration())
	assert.Equal(t, []string{"foo", "bar"}, v2.ListenAddresses())
}
