package standardcapabilities

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const testCapabilityID = "evm:ChainSelector:42@1.0.0"

func node(peerID p2ptypes.PeerID) capabilities.Node {
	return capabilities.Node{PeerID: &peerID}
}

func donWithNodes(id uint32, nodes ...capabilities.Node) capabilities.DONWithNodes {
	return capabilities.DONWithNodes{DON: capabilities.DON{ID: id}, Nodes: nodes}
}

func newTestDynamicResolver(t *testing.T, reg *mocks.CapabilitiesRegistry, getPeerID func() (p2ptypes.PeerID, error)) *dynamicCapDonIDResolver {
	r := newDynamicCapDonIDResolver(testCapabilityID, reg, getPeerID, logger.Test(t))
	r.retryInterval = time.Millisecond
	return r
}

func peerIDFn(id p2ptypes.PeerID) func() (p2ptypes.PeerID, error) {
	return func() (p2ptypes.PeerID, error) { return id, nil }
}

func TestDynamicCapDonIDResolver_ResolveCapDonID(t *testing.T) {
	t.Parallel()
	localPeerID := testPeerID(1)

	t.Run("returns matching DON ID", func(t *testing.T) {
		t.Parallel()
		reg := mocks.NewCapabilitiesRegistry(t)
		reg.EXPECT().DONsForCapability(t.Context(), testCapabilityID).Return([]capabilities.DONWithNodes{
			donWithNodes(10, node(testPeerID(2))),
			donWithNodes(20, node(localPeerID)),
		}, nil)

		got, err := newTestDynamicResolver(t, reg, peerIDFn(localPeerID)).resolveCapDonID(t.Context())
		require.NoError(t, err)
		assert.Equal(t, uint32(20), got)
	})

	t.Run("falls back to 0 when no DON matches local peer", func(t *testing.T) {
		t.Parallel()
		reg := mocks.NewCapabilitiesRegistry(t)
		reg.EXPECT().DONsForCapability(t.Context(), testCapabilityID).Return([]capabilities.DONWithNodes{
			donWithNodes(10, node(testPeerID(2))),
		}, nil)

		got, err := newTestDynamicResolver(t, reg, peerIDFn(localPeerID)).resolveCapDonID(t.Context())
		require.NoError(t, err)
		assert.Zero(t, got)
	})

	t.Run("falls back to 0 when local peer matches multiple DONs", func(t *testing.T) {
		t.Parallel()
		reg := mocks.NewCapabilitiesRegistry(t)
		reg.EXPECT().DONsForCapability(t.Context(), testCapabilityID).Return([]capabilities.DONWithNodes{
			donWithNodes(10, node(localPeerID)),
			donWithNodes(20, node(localPeerID)),
		}, nil)

		got, err := newTestDynamicResolver(t, reg, peerIDFn(localPeerID)).resolveCapDonID(t.Context())
		require.NoError(t, err)
		assert.Zero(t, got)
	})

	t.Run("errors for retry when getPeerID fails", func(t *testing.T) {
		t.Parallel()
		reg := mocks.NewCapabilitiesRegistry(t)
		r := newTestDynamicResolver(t, reg, func() (p2ptypes.PeerID, error) {
			return p2ptypes.PeerID{}, errors.New("dispatcher not ready")
		})

		_, err := r.resolveCapDonID(t.Context())
		require.Error(t, err)
	})

	t.Run("errors for retry when registry call fails", func(t *testing.T) {
		t.Parallel()
		reg := mocks.NewCapabilitiesRegistry(t)
		reg.EXPECT().DONsForCapability(t.Context(), testCapabilityID).Return(nil, errors.New("registry unavailable"))

		_, err := newTestDynamicResolver(t, reg, peerIDFn(localPeerID)).resolveCapDonID(t.Context())
		require.Error(t, err)
	})
}

// Right after node start the capability registry may not have synced yet, so
// resolution must be retried rather than settling on the workflow DON for the
// life of the process.
func TestDynamicCapDonIDResolver_AwaitCapDonID(t *testing.T) {
	t.Parallel()
	localPeerID := testPeerID(1)
	notSynced := errors.New("metadataRegistry information not available")

	t.Run("retries until the registry is available", func(t *testing.T) {
		t.Parallel()
		reg := mocks.NewCapabilitiesRegistry(t)
		reg.EXPECT().DONsForCapability(mock.Anything, testCapabilityID).Return(nil, notSynced).Times(2)
		reg.EXPECT().DONsForCapability(mock.Anything, testCapabilityID).Return([]capabilities.DONWithNodes{
			donWithNodes(10, node(localPeerID)),
		}, nil).Once()

		got, err := newTestDynamicResolver(t, reg, peerIDFn(localPeerID)).awaitCapDonID(t.Context())
		require.NoError(t, err)
		assert.Equal(t, uint32(10), got)
	})

	t.Run("does not retry a definitive zero", func(t *testing.T) {
		t.Parallel()
		reg := mocks.NewCapabilitiesRegistry(t)
		reg.EXPECT().DONsForCapability(mock.Anything, testCapabilityID).Return(nil, nil).Once()

		got, err := newTestDynamicResolver(t, reg, peerIDFn(localPeerID)).awaitCapDonID(t.Context())
		require.NoError(t, err)
		assert.Zero(t, got)
	})

	t.Run("falls back to zero on timeout", func(t *testing.T) {
		t.Parallel()
		reg := mocks.NewCapabilitiesRegistry(t)
		reg.EXPECT().DONsForCapability(mock.Anything, testCapabilityID).Return(nil, notSynced)
		r := newTestDynamicResolver(t, reg, peerIDFn(localPeerID))
		r.timeout = 20 * time.Millisecond

		got, err := r.awaitCapDonID(t.Context())
		require.NoError(t, err)
		assert.Zero(t, got)
	})

	t.Run("stops when the parent context is cancelled", func(t *testing.T) {
		t.Parallel()
		reg := mocks.NewCapabilitiesRegistry(t)
		reg.EXPECT().DONsForCapability(mock.Anything, testCapabilityID).Return(nil, notSynced).Maybe()
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		got, err := newTestDynamicResolver(t, reg, peerIDFn(localPeerID)).awaitCapDonID(ctx)
		require.NoError(t, err)
		assert.Zero(t, got)
	})
}

func TestDelegate_NewCapDonIDResolver(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		donID         uint32
		capabilityID  string
		registryAddr  string
		wantDynamic   bool
		wantStaticDON uint32
	}{
		"known DON ID":           {donID: 7, capabilityID: testCapabilityID, registryAddr: "0xabc", wantStaticDON: 7},
		"no capability ID":       {registryAddr: "0xabc"},
		"no external registry":   {capabilityID: testCapabilityID},
		"resolved from registry": {capabilityID: testCapabilityID, registryAddr: "0xabc", wantDynamic: true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			d := &Delegate{capRegistryAddress: tc.registryAddr}
			r := d.newCapDonIDResolver(logger.Test(t), tc.capabilityID, tc.donID)
			if tc.wantDynamic {
				require.IsType(t, &dynamicCapDonIDResolver{}, r)
				return
			}
			require.IsType(t, &staticCapDonIDResolver{}, r)
			got, err := r.awaitCapDonID(t.Context())
			require.NoError(t, err)
			assert.Equal(t, tc.wantStaticDON, got)
		})
	}
}
