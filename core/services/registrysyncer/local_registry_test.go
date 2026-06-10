package registrysyncer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

func TestLocalRegistry_LocalNode(t *testing.T) {
	lggr := logger.Test(t)
	localPeer := types.PeerID{0: 7}
	getPeerID := func() (types.PeerID, error) {
		return localPeer, nil
	}
	idsToDons := map[DonID]DON{
		1: {
			DON: capabilities.DON{
				ID:               1,
				F:                1,
				Members:          []types.PeerID{localPeer},
				AcceptsWorkflows: true,
			},
			CapabilityConfigurations: map[string]CapabilityConfiguration{
				"capabilityID@1.0.0": {},
			},
		},
	}
	idsToNodes := map[types.PeerID]NodeInfo{
		localPeer: {NodeOperatorID: 42},
	}
	idsToCapabilities := map[string]Capability{
		"capabilityID@1.0.0": {
			ID:             "capabilityID@1.0.0",
			CapabilityType: capabilities.CapabilityTypeAction,
		},
	}
	lr := NewLocalRegistry(lggr, getPeerID, idsToDons, idsToNodes, idsToCapabilities)

	ctx := t.Context()
	want, err := lr.NodeByPeerID(ctx, localPeer)
	require.NoError(t, err)

	got, err := lr.LocalNode(ctx)
	require.NoError(t, err)
	assert.Equal(t, want, got)

	gotAgain, err := lr.LocalNode(ctx)
	require.NoError(t, err)
	assert.Equal(t, want, gotAgain)

	t.Run("GetPeerID error", func(t *testing.T) {
		broken := NewLocalRegistry(lggr, func() (types.PeerID, error) {
			return types.PeerID{}, assert.AnError
		}, idsToDons, idsToNodes, idsToCapabilities)
		_, err := broken.LocalNode(context.Background())
		require.ErrorContains(t, err, "unable to get local node: peerWrapper hasn't started yet")
	})
}

func TestLocalRegistry_DONsForCapability(t *testing.T) {
	lggr := logger.Test(t)
	getPeerID := func() (types.PeerID, error) {
		return [32]byte{0: 1}, nil
	}
	idsToDons := map[DonID]DON{
		1: {
			DON: capabilities.DON{
				Name: "don1",
				ID:   1,
				F:    1,
				Members: []types.PeerID{
					{0: 1},
					{0: 2},
				},
			},
			CapabilityConfigurations: map[string]CapabilityConfiguration{
				"capabilityID@1.0.0": {},
			},
		},
		2: {
			DON: capabilities.DON{
				Name: "don2",
				ID:   2,
				F:    2,
				Members: []types.PeerID{
					{0: 3},
					{0: 4},
				},
			},
			CapabilityConfigurations: map[string]CapabilityConfiguration{
				"secondCapabilityID@1.0.0": {},
			},
		},
		3: {
			DON: capabilities.DON{
				Name: "don2",
				ID:   2,
				F:    2,
				Members: []types.PeerID{
					{0: 5},
					{0: 6},
				},
			},
			CapabilityConfigurations: map[string]CapabilityConfiguration{
				"thirdCapabilityID@1.0.0": {},
			},
		},
	}
	idsToNodes := map[types.PeerID]NodeInfo{
		{0: 1}: {
			NodeOperatorID: 0,
		},
		{0: 2}: {
			NodeOperatorID: 1,
		},
		{0: 3}: {
			NodeOperatorID: 2,
		},
		{0: 4}: {
			NodeOperatorID: 3,
		},
	}
	idsToCapabilities := map[string]Capability{
		"capabilityID@1.0.0": {
			ID:             "capabilityID@1.0.0",
			CapabilityType: capabilities.CapabilityTypeAction,
		},
		"secondCapabilityID@1.0.0": {
			ID:             "secondCapabilityID@1.0.0",
			CapabilityType: capabilities.CapabilityTypeAction,
		},
	}
	lr := NewLocalRegistry(lggr, getPeerID, idsToDons, idsToNodes, idsToCapabilities)

	gotDons, err := lr.DONsForCapability(t.Context(), "capabilityID@1.0.0")
	require.NoError(t, err)

	assert.Len(t, gotDons, 1)
	assert.Equal(t, idsToDons[1].DON, gotDons[0].DON)

	nodes := gotDons[0].Nodes
	assert.Len(t, nodes, 2)
	assert.Equal(t, types.PeerID{0: 1}, *nodes[0].PeerID)
	assert.Equal(t, types.PeerID{0: 2}, *nodes[1].PeerID)

	// Non-existent DON
	_, err = lr.DONsForCapability(t.Context(), "nonExistentCapabilityID@1.0.0")
	require.ErrorContains(t, err, "could not find DON for capability nonExistentCapabilityID@1.0.0")

	// thirdCapability is on a DON with invalid peers
	_, err = lr.DONsForCapability(t.Context(), "thirdCapabilityID@1.0.0")
	require.ErrorContains(t, err, "could not find node for peerID")
}

func mustMarshalProto(t *testing.T, msg proto.Message) []byte {
	t.Helper()
	b, err := proto.Marshal(msg)
	require.NoError(t, err)
	return b
}

func makeStringMap(kv map[string]string) *valuespb.Map {
	fields := make(map[string]*valuespb.Value, len(kv))
	for k, v := range kv {
		fields[k] = valuespb.NewStringValue(v)
	}
	return &valuespb.Map{Fields: fields}
}

func TestCapabilityConfiguration_Unmarshal(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		raw := mustMarshalProto(t, &capabilitiespb.CapabilityConfig{})
		cc := CapabilityConfiguration{Config: raw}

		got, err := cc.Unmarshal()
		require.NoError(t, err)
		assert.Nil(t, got.DefaultConfig)
		assert.Nil(t, got.SpecConfig)
		assert.Nil(t, got.RestrictedConfig)
		assert.Nil(t, got.Ocr3Configs)
		assert.Nil(t, got.OracleFactoryConfigs)
	})

	t.Run("invalid proto returns error", func(t *testing.T) {
		cc := CapabilityConfiguration{Config: []byte("not-valid-proto")}
		_, err := cc.Unmarshal()
		require.ErrorContains(t, err, "failed to unmarshal capability configuration")
	})

	t.Run("DefaultConfig and RestrictedConfig", func(t *testing.T) {
		raw := mustMarshalProto(t, &capabilitiespb.CapabilityConfig{
			DefaultConfig:    makeStringMap(map[string]string{"key": "val"}),
			RestrictedConfig: makeStringMap(map[string]string{"secret": "hidden"}),
			RestrictedKeys:   []string{"secret"},
		})
		cc := CapabilityConfiguration{Config: raw}

		got, err := cc.Unmarshal()
		require.NoError(t, err)

		require.NotNil(t, got.DefaultConfig)
		unwrapped, err := got.DefaultConfig.Unwrap()
		require.NoError(t, err)
		assert.Equal(t, "val", unwrapped.(map[string]any)["key"])

		require.NotNil(t, got.RestrictedConfig)
		unwrappedRC, err := got.RestrictedConfig.Unwrap()
		require.NoError(t, err)
		assert.Equal(t, "hidden", unwrappedRC.(map[string]any)["secret"])

		assert.Equal(t, []string{"secret"}, got.RestrictedKeys)
	})

	t.Run("SpecConfig", func(t *testing.T) {
		raw := mustMarshalProto(t, &capabilitiespb.CapabilityConfig{
			SpecConfig: makeStringMap(map[string]string{"interval": "60"}),
		})
		cc := CapabilityConfiguration{Config: raw}

		got, err := cc.Unmarshal()
		require.NoError(t, err)

		require.NotNil(t, got.SpecConfig)
		unwrapped, err := got.SpecConfig.Unwrap()
		require.NoError(t, err)
		assert.Equal(t, "60", unwrapped.(map[string]any)["interval"])
	})

	t.Run("Ocr3Configs", func(t *testing.T) {
		signer := []byte{0x01, 0x02, 0x03}
		transmitter := []byte{0xde, 0xad, 0xbe, 0xef}

		raw := mustMarshalProto(t, &capabilitiespb.CapabilityConfig{
			Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
				"__default__": {
					Signers:               [][]byte{signer},
					Transmitters:          [][]byte{transmitter},
					F:                     2,
					OnchainConfig:         []byte("onchain"),
					OffchainConfigVersion: 5,
					OffchainConfig:        []byte("offchain"),
					ConfigCount:           3,
				},
			},
		})
		cc := CapabilityConfiguration{Config: raw}

		got, err := cc.Unmarshal()
		require.NoError(t, err)

		require.Contains(t, got.Ocr3Configs, "__default__")
		cfg := got.Ocr3Configs["__default__"]
		assert.Equal(t, []ocrtypes.OnchainPublicKey{signer}, cfg.Signers)
		assert.Equal(t, []ocrtypes.Account{ocrtypes.Account("deadbeef")}, cfg.Transmitters)
		assert.Equal(t, uint8(2), cfg.F)
		assert.Equal(t, []byte("onchain"), cfg.OnchainConfig)
		assert.Equal(t, uint64(5), cfg.OffchainConfigVersion)
		assert.Equal(t, []byte("offchain"), cfg.OffchainConfig)
		assert.Equal(t, uint64(3), cfg.ConfigCount)
	})

	t.Run("Ocr3Configs normalizes transmitters to hex text", func(t *testing.T) {
		raw := mustMarshalProto(t, &capabilitiespb.CapabilityConfig{
			Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
				"aptos": {
					Transmitters: [][]byte{
						{0x00, 0xff},
						[]byte("ascii-bytes"),
					},
				},
			},
		})
		cc := CapabilityConfiguration{Config: raw}

		got, err := cc.Unmarshal()
		require.NoError(t, err)

		require.Contains(t, got.Ocr3Configs, "aptos")
		cfg := got.Ocr3Configs["aptos"]
		assert.Equal(t,
			[]ocrtypes.Account{
				ocrtypes.Account("00ff"),
				ocrtypes.Account("61736369692d6279746573"),
			},
			cfg.Transmitters,
		)
	})

	t.Run("OracleFactoryConfigs", func(t *testing.T) {
		raw := mustMarshalProto(t, &capabilitiespb.CapabilityConfig{
			OracleFactoryConfigs: map[string]*valuespb.Map{
				"blue": makeStringMap(map[string]string{"mode": "blue"}),
			},
		})
		cc := CapabilityConfiguration{Config: raw}

		got, err := cc.Unmarshal()
		require.NoError(t, err)

		require.Contains(t, got.OracleFactoryConfigs, "blue")
		blueMap := got.OracleFactoryConfigs["blue"]
		unwrapped, err := blueMap.Unwrap()
		require.NoError(t, err)
		assert.Equal(t, "blue", unwrapped.(map[string]any)["mode"])
	})

	t.Run("LocalOnly flag", func(t *testing.T) {
		raw := mustMarshalProto(t, &capabilitiespb.CapabilityConfig{
			LocalOnly: true,
		})
		cc := CapabilityConfiguration{Config: raw}

		got, err := cc.Unmarshal()
		require.NoError(t, err)
		assert.True(t, got.LocalOnly)
	})

	t.Run("all fields together", func(t *testing.T) {
		raw := mustMarshalProto(t, &capabilitiespb.CapabilityConfig{
			DefaultConfig: makeStringMap(map[string]string{"dc": "1"}),
			SpecConfig:    makeStringMap(map[string]string{"sc": "2"}),
			Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
				"__default__": {
					Signers:      [][]byte{{0xAA}},
					Transmitters: [][]byte{[]byte("tx1")},
					F:            1,
					ConfigCount:  7,
				},
			},
			OracleFactoryConfigs: map[string]*valuespb.Map{
				"green": makeStringMap(map[string]string{"x": "y"}),
			},
			LocalOnly: true,
		})
		cc := CapabilityConfiguration{Config: raw}

		got, err := cc.Unmarshal()
		require.NoError(t, err)

		assert.NotNil(t, got.DefaultConfig)
		assert.NotNil(t, got.SpecConfig)
		assert.Len(t, got.Ocr3Configs, 1)
		assert.Len(t, got.OracleFactoryConfigs, 1)
		assert.True(t, got.LocalOnly)
	})
}

// rawCapabilityConfigHex is the serialized capabilitiespb.CapabilityConfig taken
// from the "test value" comment above CapabilityConfiguration.Unmarshal in
// local_registry.go. It carries a nested DefaultConfig values.Map.
const rawCapabilityConfigHex = "0x0acd080aca080a08456e636c6176657312bd082aba081299042296040ada020a0d5472757374656456616c75657312c8022ac50212c2021abf027b2270637230223a22336130646135323235666439323134323666633864613065376139316534323631623431333137336239626530303237616462326666626231616661653430616462613662663431626133663066653537386339303032646134653565613237222c2270637231223a22346234643562333636316233656663313239323039303063383065313236653463653738336335323264653663303261326135626637616633613262393332376238363737366631383865346265316331633430346131323964626461343933222c2270637232223a22643938306538363964333066303333353130356563373236623234393430633636633130383130376362386636653661363638363561323562646464303537303363393231373964303334386661323936636166336332663633313036393631227d0a170a11456e636c6176654175746848656164657212020a000a160a10456e636c61766545787472614461746112021a000a2f0a09456e636c617665494412221a2066338ab3ebe1d246feb5849330e92022b9d66b0c0753b1f2ac6e356acbff20fc0a160a0b456e636c6176655479706512070a054e4954524f0a3d0a0a456e636c61766555524c122f0a2d68747470733a2f2f636f6e666964656e7469616c2d687474702e656e636c617665732e636861696e2e6c696e6b129b042298040a170a11456e636c6176654175746848656164657212020a000a160a10456e636c61766545787472614461746112021a000a2f0a09456e636c617665494412221a206401f4177859d87613cdc14bd3841049840abaf6d0163cf77080227b8866ebe20a160a0b456e636c6176655479706512070a054e4954524f0a3f0a0a456e636c61766555524c12310a2f68747470733a2f2f636f6e666964656e7469616c2d687474702d322e656e636c617665732e636861696e2e6c696e6b0ada020a0d5472757374656456616c75657312c8022ac50212c2021abf027b2270637230223a22336130646135323235666439323134323666633864613065376139316534323631623431333137336239626530303237616462326666626231616661653430616462613662663431626133663066653537386339303032646134653565613237222c2270637231223a22346234643562333636316233656663313239323039303063383065313236653463653738336335323264653663303261326135626637616633613262393332376238363737366631383865346265316331633430346131323964626461343933222c2270637232223a22643938306538363964333066303333353130356563373236623234393430633636633130383130376362386636653661363638363561323562646464303537303363393231373964303334386661323936636166336332663633313036393631227d4001"

func TestCapabilityConfiguration_Unmarshal_RawBytes(t *testing.T) {
	raw, err := hex.DecodeString(strings.TrimPrefix(rawCapabilityConfigHex, "0x"))
	require.NoError(t, err)

	cc := CapabilityConfiguration{Config: raw}
	got, err := cc.Unmarshal()
	require.NoError(t, err)

	t.Logf("LocalOnly: %v", got.LocalOnly)
	t.Logf("RestrictedKeys: %v", got.RestrictedKeys)
	t.Logf("DefaultConfig:\n%s", prettyValueMap(t, got.DefaultConfig))
	t.Logf("RestrictedConfig:\n%s", prettyValueMap(t, got.RestrictedConfig))
	t.Logf("SpecConfig:\n%s", prettyValueMap(t, got.SpecConfig))

	if len(got.OracleFactoryConfigs) == 0 {
		t.Log("OracleFactoryConfigs: <none>")
	}
	for name, m := range got.OracleFactoryConfigs {
		m := m
		t.Logf("OracleFactoryConfig[%q]:\n%s", name, prettyValueMap(t, &m))
	}

	if len(got.Ocr3Configs) == 0 {
		t.Log("Ocr3Configs: <none>")
	}
	for name, cfg := range got.Ocr3Configs {
		t.Logf("Ocr3Config[%q]: %+v", name, cfg)
	}

	// The raw bytes encode a DefaultConfig values.Map; assert we decoded it and
	// can fully traverse its nested keys/values.
	require.NotNil(t, got.DefaultConfig)
	unwrapped, err := got.DefaultConfig.Unwrap()
	require.NoError(t, err)
	assert.NotEmpty(t, unwrapped)
}

// prettyValueMap unwraps a values.Map into a fully-resolved map[string]any and
// renders it as indented JSON so every nested key/value is visible. Byte-slice
// leaves are rendered as base64 by encoding/json.
func prettyValueMap(t *testing.T, m *values.Map) string {
	t.Helper()
	if m == nil {
		return "<nil>"
	}
	unwrapped, err := m.Unwrap()
	require.NoError(t, err)
	b, err := json.MarshalIndent(unwrapped, "", "  ")
	require.NoError(t, err)
	return string(b)
}
