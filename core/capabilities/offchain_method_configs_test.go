package capabilities

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	regpkg "github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
)

func TestOffchainMethodConfigs_Selection(t *testing.T) {
	t.Parallel()
	const capID = "cap@1.0.0"
	method := &capabilitiespb.CapabilityMethodConfig{RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig{RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{}}}
	raw, err := proto.Marshal(&capabilitiespb.CapabilityConfig{MethodConfigs: map[string]*capabilitiespb.CapabilityMethodConfig{"Legacy": method}})
	require.NoError(t, err)
	onchain := regpkg.CapabilityConfiguration{Config: raw}
	don := regpkg.DON{ID: 1, Name: "don-1"}
	snapshot := func(caps map[string]*capabilitiespb.CapabilityConfig) *offchainMethodConfigs {
		return newOffchainMethodConfigs(&capabilitiespb.OffchainCapabilitiesRegistry{Dons: map[string]*capabilitiespb.DONConfig{"don-1": {Capabilities: caps}}}, []regpkg.DON{don})
	}
	for _, tc := range []struct {
		name    string
		snap    *offchainMethodConfigs
		methods []string
	}{
		{"nil snapshot", nil, []string{"Legacy"}},
		{"missing capability", snapshot(map[string]*capabilitiespb.CapabilityConfig{}), []string{"Legacy"}},
		{"complete replacement", snapshot(map[string]*capabilitiespb.CapabilityConfig{capID: {MethodConfigs: map[string]*capabilitiespb.CapabilityMethodConfig{"New": method}}}), []string{"New"}},
		{"empty method map", snapshot(map[string]*capabilitiespb.CapabilityConfig{capID: {}}), nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config, selectErr := tc.snap.configuration(capID, don, onchain)
			require.NoError(t, selectErr)
			assert.Len(t, config.CapabilityMethodConfig, len(tc.methods))
			for _, method := range tc.methods {
				assert.Contains(t, config.CapabilityMethodConfig, method)
			}
		})
	}
	snap := snapshot(map[string]*capabilitiespb.CapabilityConfig{capID: {MethodConfigs: map[string]*capabilitiespb.CapabilityMethodConfig{"New": method}}})
	// An unused malformed on-chain config cannot block an offchain capability.
	config, err := snap.configuration(capID, don, regpkg.CapabilityConfiguration{Config: []byte{0xff}})
	require.NoError(t, err)
	assert.Contains(t, config.CapabilityMethodConfig, "New")
	// Unknown, unnamed, or ambiguous DON names retain the legacy source.
	for _, dons := range [][]regpkg.DON{nil, {{ID: 1}}, {{ID: 1, Name: "don-1"}, {ID: 2, Name: "don-1"}}} {
		unresolved := newOffchainMethodConfigs(snap.reg, dons)
		legacy, selectErr := unresolved.configuration(capID, don, onchain)
		require.NoError(t, selectErr)
		assert.Contains(t, legacy.CapabilityMethodConfig, "Legacy")
	}
}

func TestOffchainMethodConfigs_LocalOnly(t *testing.T) {
	t.Parallel()
	const capID = "cap@1.0.0"
	don := regpkg.DON{ID: 1, Name: "don-1"}
	for _, onchainLocalOnly := range []bool{false, true} {
		raw, err := proto.Marshal(&capabilitiespb.CapabilityConfig{LocalOnly: onchainLocalOnly})
		require.NoError(t, err)
		onchain := regpkg.CapabilityConfiguration{Config: raw}
		var disabled *offchainMethodConfigs
		config, err := disabled.configuration(capID, don, onchain)
		require.NoError(t, err)
		assert.Equal(t, onchainLocalOnly, config.LocalOnly, "gate off retains the on-chain flag")
		for _, offchainLocalOnly := range []bool{false, true} {
			gc := globalconfig.New()
			_, err = gc.Store(fmt.Sprintf(`{"version":1,"dons":{"don-1":{"capabilities":{"cap@1.0.0":{"localOnly":%t}}}}}`, offchainLocalOnly), "payload")
			require.NoError(t, err)
			snapshot := newOffchainMethodConfigs(gc.LoadParsed(), []regpkg.DON{don})
			config, err = snapshot.configuration(capID, don, onchain)
			require.NoError(t, err)
			assert.Equal(t, offchainLocalOnly, config.LocalOnly, "the selected offchain capability supplies the whole flag")
			config, err = snapshot.configuration("missing@1.0.0", don, onchain)
			require.NoError(t, err)
			assert.Equal(t, onchainLocalOnly, config.LocalOnly, "a missing capability retains the on-chain flag")
		}
	}
}
