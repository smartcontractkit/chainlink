package registrysyncer

import (
	"testing"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

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
		transmitter := []byte("0xabc")

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
		assert.Equal(t, []ocrtypes.Account{ocrtypes.Account(transmitter)}, cfg.Transmitters)
		assert.Equal(t, uint8(2), cfg.F)
		assert.Equal(t, []byte("onchain"), cfg.OnchainConfig)
		assert.Equal(t, uint64(5), cfg.OffchainConfigVersion)
		assert.Equal(t, []byte("offchain"), cfg.OffchainConfig)
		assert.Equal(t, uint64(3), cfg.ConfigCount)
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
