package capabilities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	regpkg "github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types/mocks"
)

// requestTimeoutFor returns the requestTimeout the launcher wired for (capID, method),
// read from the executable client's SetConfig log line.
func requestTimeoutFor(t *testing.T, logs *observer.ObservedLogs, capID, method string) time.Duration {
	t.Helper()
	for _, entry := range logs.All() {
		if entry.Message != "SetConfig" {
			continue
		}
		if entry.ContextMap()["capabilityID"] != capID || entry.ContextMap()["capMethodName"] != method {
			continue
		}
		rt, ok := entry.ContextMap()["requestTimeout"].(time.Duration)
		if !ok {
			// zap encodes durations as strings in output; the observer may keep either.
			if s, isStr := entry.ContextMap()["requestTimeout"].(string); isStr {
				d, err := time.ParseDuration(s)
				require.NoError(t, err)
				return d
			}
			require.True(t, ok, "SetConfig log for %s/%s has no requestTimeout", capID, method)
		}
		return rt
	}
	t.Fatalf("no SetConfig log for %s/%s", capID, method)
	return 0
}

// offchainMethodCfg builds an offchain registry payload with method_configs for one capability
// on one DON, keyed by DON name.
func offchainMethodCfg(t *testing.T, version uint64, donName string, capID string, methodConfigs map[string]*capabilitiespb.CapabilityMethodConfig) *globalconfig.GlobalConfig {
	t.Helper()
	reg := &capabilitiespb.OffchainCapabilitiesRegistry{
		Domain:  "cre",
		Env:     "test",
		Version: version,
		Dons: map[string]*capabilitiespb.DONConfig{
			donName: {
				Capabilities: map[string]*capabilitiespb.CapabilityConfig{
					capID: {MethodConfigs: methodConfigs},
				},
			},
		},
	}
	raw, err := protojson.Marshal(reg)
	require.NoError(t, err)
	gc := globalconfig.New()
	require.NoError(t, gc.Store(globalconfig.Update{Raw: string(raw)}))
	return gc
}

func onchainMethodCfg(t *testing.T, methodConfigs map[string]*capabilitiespb.CapabilityMethodConfig) []byte {
	t.Helper()
	b, err := proto.Marshal(&capabilitiespb.CapabilityConfig{MethodConfigs: methodConfigs})
	require.NoError(t, err)
	return b
}

func TestOffchainMethodConfigs_Overlay(t *testing.T) {
	t.Parallel()

	onchain := map[string]capabilities.CapabilityMethodConfig{
		"Write": {
			RemoteExecutableConfig: &capabilities.RemoteExecutableConfig{RequestTimeout: 30 * time.Second},
		},
		"View": {
			RemoteExecutableConfig: &capabilities.RemoteExecutableConfig{RequestTimeout: 10 * time.Second},
		},
	}
	don := regpkg.DON{ID: 1, Name: "don-1"}

	t.Run("nil snapshot is a no-op", func(t *testing.T) {
		t.Parallel()
		var snap *offchainMethodConfigs
		assert.Equal(t, onchain, snap.overlay("cap@1.0.0", don, onchain))
	})

	t.Run("empty registry snapshot is a no-op", func(t *testing.T) {
		t.Parallel()
		snap := &offchainMethodConfigs{reg: nil, donNames: map[uint32]string{1: "don-1"}}
		assert.Equal(t, onchain, snap.overlay("cap@1.0.0", don, onchain))
	})

	t.Run("offchain wins per method, omitted methods keep on-chain config", func(t *testing.T) {
		t.Parallel()
		snap := &offchainMethodConfigs{
			reg: &capabilitiespb.OffchainCapabilitiesRegistry{
				Dons: map[string]*capabilitiespb.DONConfig{
					"don-1": {
						Capabilities: map[string]*capabilitiespb.CapabilityConfig{
							"cap@1.0.0": {
								MethodConfigs: map[string]*capabilitiespb.CapabilityMethodConfig{
									// overrides Write's timeout
									"Write": {
										RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
											RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
												RequestTimeout: durationpb.New(time.Minute),
											},
										},
									},
									// View omitted offchain
								},
							},
						},
					},
				},
			},
			donNames: map[uint32]string{1: "don-1"},
		}
		got := snap.overlay("cap@1.0.0", don, onchain)
		assert.Equal(t, time.Minute, got["Write"].RemoteExecutableConfig.RequestTimeout, "offchain entry replaces on-chain entry for the same method")
		assert.Equal(t, 10*time.Second, got["View"].RemoteExecutableConfig.RequestTimeout, "method omitted offchain keeps its on-chain config")
		assert.Len(t, got, 2)
	})

	t.Run("capability absent offchain keeps full on-chain map", func(t *testing.T) {
		t.Parallel()
		snap := &offchainMethodConfigs{
			reg: &capabilitiespb.OffchainCapabilitiesRegistry{
				Dons: map[string]*capabilitiespb.DONConfig{
					"don-1": {
						Capabilities: map[string]*capabilitiespb.CapabilityConfig{
							"other-cap@1.0.0": {MethodConfigs: map[string]*capabilitiespb.CapabilityMethodConfig{
								"Write": {},
							}},
						},
					},
				},
			},
			donNames: map[uint32]string{1: "don-1"},
		}
		assert.Equal(t, onchain, snap.overlay("cap@1.0.0", don, onchain))
	})

	t.Run("DON without a usable name keeps on-chain config", func(t *testing.T) {
		t.Parallel()
		snap := &offchainMethodConfigs{
			reg: &capabilitiespb.OffchainCapabilitiesRegistry{
				Dons: map[string]*capabilitiespb.DONConfig{
					"don-1": {
						Capabilities: map[string]*capabilitiespb.CapabilityConfig{
							"cap@1.0.0": {MethodConfigs: map[string]*capabilitiespb.CapabilityMethodConfig{
								"Write": {},
							}},
						},
					},
				},
			},
			donNames: map[uint32]string{}, // don 1 not resolvable
		}
		assert.Equal(t, onchain, snap.overlay("cap@1.0.0", don, onchain))
	})

	t.Run("offchain-only method is added", func(t *testing.T) {
		t.Parallel()
		snap := &offchainMethodConfigs{
			reg: &capabilitiespb.OffchainCapabilitiesRegistry{
				Dons: map[string]*capabilitiespb.DONConfig{
					"don-1": {
						Capabilities: map[string]*capabilitiespb.CapabilityConfig{
							"cap@1.0.0": {MethodConfigs: map[string]*capabilitiespb.CapabilityMethodConfig{
								"NewMethod": {
									RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig{
										RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
											MinResponsesToAggregate: 2,
										},
									},
								},
							}},
						},
					},
				},
			},
			donNames: map[uint32]string{1: "don-1"},
		}
		got := snap.overlay("cap@1.0.0", don, onchain)
		require.Len(t, got, 3)
		require.NotNil(t, got["NewMethod"].RemoteTriggerConfig)
		assert.Equal(t, uint32(2), got["NewMethod"].RemoteTriggerConfig.MinResponsesToAggregate)
	})
}

// TestOffchainMethodConfigs_GateOff_Unchanged asserts that with the cutover gate off, the
// launcher wires the on-chain method configs byte-for-byte (backwards compatibility).
func TestOffchainMethodConfigs_GateOff_Unchanged(t *testing.T) {
	t.Parallel()
	lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
	reg := regpkg.NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	workflowDonNodes := newNodes(4)
	capabilityDonNodes := newNodes(4)
	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeer.On("ID").Return(workflowDonNodes[0])
	sharedPeer.On("IsBootstrap").Return(false)
	sharedPeer.On("UpdateConnectionsByDONs", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	fullCapID := "write-chain_evm_1@1.0.0"
	capIDHash := RandomUTF8BytesWord()

	onchainCfg := onchainMethodCfg(t, map[string]*capabilitiespb.CapabilityMethodConfig{
		"Write": {
			RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
				RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
					RequestTimeout: durationpb.New(30 * time.Second),
				},
			},
		},
	})

	wfDONID, capDONID := uint32(1), uint32(2)
	localRegistry := buildLocalRegistry()
	addDON(localRegistry, wfDONID, uint32(0), uint8(1), true, true, workflowDonNodes, []string{"zone-a"}, 1, nil)
	addDON(localRegistry, capDONID, uint32(0), uint8(1), true, false, capabilityDonNodes, []string{"zone-a"}, 1, [][32]byte{capIDHash})
	addCapabilityToDON(localRegistry, capDONID, fullCapID, capabilities.CapabilityTypeTarget, onchainCfg)

	// Offchain payload carries a DIFFERENT Write timeout; with the gate off it must be ignored.
	gc := offchainMethodCfg(t, 1, "don-2", fullCapID, map[string]*capabilitiespb.CapabilityMethodConfig{
		"Write": {
			RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
				RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
					RequestTimeout: durationpb.New(time.Minute),
				},
			},
		},
	})

	launcher, err := NewLauncher(lggr, sharedPeer, nil, dispatcher, reg, &mockDonNotifier{}, limits.Factory{}, false, 0)
	require.NoError(t, err)
	launcher.SetOffchainRegistry(gc, false) // gate OFF
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()

	dispatcher.On("SetReceiverForMethod", fullCapID, capDONID, "Write", mock.AnythingOfType("*executable.client")).Return(nil)
	require.NoError(t, launcher.OnNewRegistry(t.Context(), localRegistry))

	// The wired method config must be the on-chain one (30s), not the offchain one (1m).
	assert.Equal(t, 30*time.Second, requestTimeoutFor(t, logs, fullCapID, "Write"),
		"gate off: on-chain method config must be used unchanged")
}

// TestOffchainMethodConfigs_GateOn_Overrides asserts that with the cutover gate on, the
// launcher wires the offchain method config (offchain-wins-with-fallback).
func TestOffchainMethodConfigs_GateOn_Overrides(t *testing.T) {
	t.Parallel()
	lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
	reg := regpkg.NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	workflowDonNodes := newNodes(4)
	capabilityDonNodes := newNodes(4)
	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeer.On("ID").Return(workflowDonNodes[0])
	sharedPeer.On("IsBootstrap").Return(false)
	sharedPeer.On("UpdateConnectionsByDONs", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	fullCapID := "write-chain_evm_1@1.0.0"
	capIDHash := RandomUTF8BytesWord()

	onchainCfg := onchainMethodCfg(t, map[string]*capabilitiespb.CapabilityMethodConfig{
		"Write": {
			RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
				RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
					RequestTimeout: durationpb.New(30 * time.Second),
				},
			},
		},
		"View": {
			RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
				RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
					RequestTimeout: durationpb.New(10 * time.Second),
				},
			},
		},
	})

	wfDONID, capDONID := uint32(1), uint32(2)
	localRegistry := buildLocalRegistry()
	addDON(localRegistry, wfDONID, uint32(0), uint8(1), true, true, workflowDonNodes, []string{"zone-a"}, 1, nil)
	addDON(localRegistry, capDONID, uint32(0), uint8(1), true, false, capabilityDonNodes, []string{"zone-a"}, 1, [][32]byte{capIDHash})
	addCapabilityToDON(localRegistry, capDONID, fullCapID, capabilities.CapabilityTypeTarget, onchainCfg)
	// Name the capability DON so the offchain payload can key it.
	capDON := localRegistry.IDsToDONs[regpkg.DonID(capDONID)]
	capDON.Name = "don-2"
	localRegistry.IDsToDONs[regpkg.DonID(capDONID)] = capDON

	// Offchain overrides Write's timeout; View omitted offchain keeps its on-chain config.
	gc := offchainMethodCfg(t, 1, "don-2", fullCapID, map[string]*capabilitiespb.CapabilityMethodConfig{
		"Write": {
			RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
				RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
					RequestTimeout: durationpb.New(time.Minute),
				},
			},
		},
	})

	launcher, err := NewLauncher(lggr, sharedPeer, nil, dispatcher, reg, &mockDonNotifier{}, limits.Factory{}, false, 0)
	require.NoError(t, err)
	launcher.SetOffchainRegistry(gc, true) // gate ON
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()

	dispatcher.On("SetReceiverForMethod", fullCapID, capDONID, "Write", mock.AnythingOfType("*executable.client")).Return(nil)
	dispatcher.On("SetReceiverForMethod", fullCapID, capDONID, "View", mock.AnythingOfType("*executable.client")).Return(nil)
	require.NoError(t, launcher.OnNewRegistry(t.Context(), localRegistry))

	assert.Equal(t, time.Minute, requestTimeoutFor(t, logs, fullCapID, "Write"),
		"gate on: offchain Write config must win")
	assert.Equal(t, 10*time.Second, requestTimeoutFor(t, logs, fullCapID, "View"),
		"gate on: method omitted offchain keeps its on-chain config")
}

// TestOffchainMethodConfigs_NoThrashAcrossPasses locks the central invariant of the method_configs
// overlay: the prune path (computeWantedShimKeys) and the create path (addRemoteCapabilities) must
// apply the SAME offchain overlay. An offchain-only method's shim created on the first OnNewRegistry
// pass must survive the second pass untouched.
//
// A single pass cannot catch a prune/create overlay mismatch: on the first pass the shim cache is
// empty, so pruneStaleShims has nothing to remove and the method is created regardless. The thrash
// only manifests on the second pass — if computeWantedShimKeys did NOT apply the overlay, the
// offchain-only method would be absent from the wanted set and torn down (RemoveReceiverForMethod)
// before being re-created, churning every tick. RemoveReceiverForMethod is deliberately left
// unexpected on the dispatcher mock, so any such prune fails the test.
func TestOffchainMethodConfigs_NoThrashAcrossPasses(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	reg := regpkg.NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	workflowDonNodes := newNodes(4)
	capabilityDonNodes := newNodes(4)
	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeer.On("ID").Return(workflowDonNodes[0])
	sharedPeer.On("IsBootstrap").Return(false)
	sharedPeer.On("UpdateConnectionsByDONs", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	fullCapID := "write-chain_evm_1@1.0.0"
	capIDHash := RandomUTF8BytesWord()

	// On-chain carries a single executable method, "Write".
	onchainCfg := onchainMethodCfg(t, map[string]*capabilitiespb.CapabilityMethodConfig{
		"Write": {
			RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
				RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
					RequestTimeout: durationpb.New(30 * time.Second),
				},
			},
		},
	})

	wfDONID, capDONID := uint32(1), uint32(2)
	localRegistry := buildLocalRegistry()
	addDON(localRegistry, wfDONID, uint32(0), uint8(1), true, true, workflowDonNodes, []string{"zone-a"}, 1, nil)
	addDON(localRegistry, capDONID, uint32(0), uint8(1), true, false, capabilityDonNodes, []string{"zone-a"}, 1, [][32]byte{capIDHash})
	addCapabilityToDON(localRegistry, capDONID, fullCapID, capabilities.CapabilityTypeTarget, onchainCfg)
	// Name the capability DON so the offchain payload can key it.
	capDON := localRegistry.IDsToDONs[regpkg.DonID(capDONID)]
	capDON.Name = "don-2"
	localRegistry.IDsToDONs[regpkg.DonID(capDONID)] = capDON

	// Offchain adds an executable method "Extra" that exists ONLY offchain (not on-chain).
	gc := offchainMethodCfg(t, 1, "don-2", fullCapID, map[string]*capabilitiespb.CapabilityMethodConfig{
		"Extra": {
			RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
				RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
					RequestTimeout: durationpb.New(45 * time.Second),
				},
			},
		},
	})

	launcher, err := NewLauncher(lggr, sharedPeer, nil, dispatcher, reg, &mockDonNotifier{}, limits.Factory{}, false, 0)
	require.NoError(t, err)
	launcher.SetOffchainRegistry(gc, true) // gate ON
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()

	// Both the on-chain "Write" and the offchain-only "Extra" client shims are created on the first
	// pass. SetReceiverForMethod fires only on shim creation, so it is called twice total (pass 1),
	// not again on pass 2. RemoveReceiverForMethod is intentionally NOT registered: if the overlay
	// is missing from the wanted-set computation, pass 2 prunes "Extra" and the unexpected call fails
	// the mock.
	dispatcher.On("SetReceiverForMethod", fullCapID, capDONID, "Write", mock.AnythingOfType("*executable.client")).Return(nil)
	dispatcher.On("SetReceiverForMethod", fullCapID, capDONID, "Extra", mock.AnythingOfType("*executable.client")).Return(nil)

	require.NoError(t, launcher.OnNewRegistry(t.Context(), localRegistry))
	// Second pass with identical inputs must be a no-op for the shim set.
	require.NoError(t, launcher.OnNewRegistry(t.Context(), localRegistry))

	dispatcher.AssertNotCalled(t, "RemoveReceiverForMethod", fullCapID, capDONID, "Extra")
	dispatcher.AssertNotCalled(t, "RemoveReceiverForMethod", fullCapID, capDONID, "Write")
}
