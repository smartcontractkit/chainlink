package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"
)

func TestP2PConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	p2p := cfg.P2P()
	assert.Equal(t, "p2p_12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw", p2p.PeerID().String())
	assert.Equal(t, 13, p2p.IncomingMessageBufferSize())
	assert.Equal(t, 17, p2p.OutgoingMessageBufferSize())
	assert.True(t, p2p.TraceLogging())

	v1 := p2p.V1()
	assert.True(t, v1.Enabled())
	assert.Equal(t, "1.2.3.4", v1.AnnounceIP().String())
	assert.Equal(t, uint16(1234), v1.AnnouncePort())
	assert.Equal(t, time.Minute, v1.BootstrapCheckInterval())
	p, err := v1.DefaultBootstrapPeers()
	require.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar", "should", "these", "be", "typed"}, p)
	assert.Equal(t, uint32(4321), v1.DHTAnnouncementCounterUserPrefix())
	assert.Equal(t, 9, v1.DHTLookupInterval())
	assert.Equal(t, "4.3.2.1", v1.ListenIP().String())
	assert.Equal(t, uint16(9), v1.ListenPort())
	assert.Equal(t, time.Second, v1.NewStreamTimeout())
	assert.Equal(t, time.Minute, v1.PeerstoreWriteInterval())

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
	assert.Equal(t, time.Second, v2.DeltaReconcile().Duration())
	assert.Equal(t, []string{"foo", "bar"}, v2.ListenAddresses())
}
