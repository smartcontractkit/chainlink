package aptos

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
)

func TestSetRuntimeSpecConfig_ReplacesLegacyKey(t *testing.T) {
	specConfig := values.EmptyMap()
	legacy, err := values.Wrap([]string{"0x1"})
	require.NoError(t, err)
	specConfig.Underlying[legacyTransmittersKey] = legacy

	capConfig := &capabilitiespb.CapabilityConfig{
		SpecConfig: values.ProtoMap(specConfig),
	}

	expectedMap := map[string]string{
		"peer-a": "0x000000000000000000000000000000000000000000000000000000000000000a",
	}
	require.NoError(t, setRuntimeSpecConfig(capConfig, methodConfigSettings{
		TransmissionSchedule: capabilitiespb.TransmissionSchedule_AllAtOnce,
		DeltaStage:           1500 * time.Millisecond,
	}, expectedMap))

	decoded, err := values.FromMapValueProto(capConfig.SpecConfig)
	require.NoError(t, err)
	require.NotNil(t, decoded)
	require.NotContains(t, decoded.Underlying, legacyTransmittersKey)

	rawSchedule, ok := decoded.Underlying[specConfigScheduleKey]
	require.True(t, ok)
	schedule, err := rawSchedule.Unwrap()
	require.NoError(t, err)
	require.Equal(t, "allAtOnce", schedule)

	rawDeltaStage, ok := decoded.Underlying[specConfigDeltaStageKey]
	require.True(t, ok)
	deltaStage, err := rawDeltaStage.Unwrap()
	require.NoError(t, err)
	require.EqualValues(t, 1500*time.Millisecond, deltaStage)

	rawMap, ok := decoded.Underlying[specConfigP2PMapKey]
	require.True(t, ok)
	unwrapped, err := rawMap.Unwrap()
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"peer-a": "0x000000000000000000000000000000000000000000000000000000000000000a",
	}, unwrapped)
}

func TestBuildCapabilityConfig_UsesMethodConfigsAndSpecConfig(t *testing.T) {
	capConfig, err := BuildCapabilityConfig(
		map[string]any{
			requestTimeoutKey: "45s",
			deltaStageKey:     "2500ms",
		},
		map[string]string{
			"peer-a": "0x000000000000000000000000000000000000000000000000000000000000000a",
		},
		false,
	)
	require.NoError(t, err)
	require.False(t, capConfig.LocalOnly)
	require.Nil(t, capConfig.Ocr3Configs)
	require.Contains(t, capConfig.MethodConfigs, "View")
	require.Contains(t, capConfig.MethodConfigs, "WriteReport")

	writeCfg := capConfig.MethodConfigs["WriteReport"].GetRemoteExecutableConfig()
	require.NotNil(t, writeCfg)
	require.Equal(t, capabilitiespb.TransmissionSchedule_AllAtOnce, writeCfg.TransmissionSchedule)
	require.Equal(t, 2500*time.Millisecond, writeCfg.DeltaStage.AsDuration())
	require.Equal(t, 45*time.Second, writeCfg.RequestTimeout.AsDuration())

	specConfig, err := values.FromMapValueProto(capConfig.SpecConfig)
	require.NoError(t, err)
	require.NotNil(t, specConfig)

	rawSchedule, ok := specConfig.Underlying[specConfigScheduleKey]
	require.True(t, ok)
	schedule, err := rawSchedule.Unwrap()
	require.NoError(t, err)
	require.Equal(t, "allAtOnce", schedule)

	rawDeltaStage, ok := specConfig.Underlying[specConfigDeltaStageKey]
	require.True(t, ok)
	deltaStage, err := rawDeltaStage.Unwrap()
	require.NoError(t, err)
	require.EqualValues(t, 2500*time.Millisecond, deltaStage)

	rawMap, ok := specConfig.Underlying[specConfigP2PMapKey]
	require.True(t, ok)
	unwrapped, err := rawMap.Unwrap()
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"peer-a": "0x000000000000000000000000000000000000000000000000000000000000000a",
	}, unwrapped)
}

func TestBuildCapabilityConfig_WithoutP2PMap_StillSetsRuntimeSpecConfig(t *testing.T) {
	capConfig, err := BuildCapabilityConfig(nil, nil, true)
	require.NoError(t, err)
	require.True(t, capConfig.LocalOnly)
	require.Nil(t, capConfig.Ocr3Configs)
	require.Contains(t, capConfig.MethodConfigs, "View")
	require.Contains(t, capConfig.MethodConfigs, "WriteReport")

	specConfig, err := values.FromMapValueProto(capConfig.SpecConfig)
	require.NoError(t, err)
	require.NotNil(t, specConfig)
	require.NotContains(t, specConfig.Underlying, specConfigP2PMapKey)
	require.Contains(t, specConfig.Underlying, specConfigScheduleKey)
	require.Contains(t, specConfig.Underlying, specConfigDeltaStageKey)
}

func TestBuildWorkerConfigJSON_IncludesLocalRuntimeValues(t *testing.T) {
	configStr, err := buildWorkerConfigJSON(
		4,
		"0x000000000000000000000000000000000000000000000000000000000000000a",
		methodConfigSettings{DeltaStage: 2500 * time.Millisecond},
		map[string]string{"peer-a": "0x1"},
		true,
	)
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal([]byte(configStr), &got))
	require.Equal(t, "4", got["chainId"])
	require.Equal(t, "aptos", got["network"])
	require.Equal(t, true, got["isLocal"])
	require.EqualValues(t, (2500 * time.Millisecond).Nanoseconds(), got["deltaStage"])
	require.Equal(t, map[string]any{"peer-a": "0x1"}, got[specConfigP2PMapKey])
}

func TestNormalizeTransmitter(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "short address is normalized",
			input: "0xa",
			want:  "0x000000000000000000000000000000000000000000000000000000000000000a",
		},
		{
			name:  "whitespace is trimmed",
			input: " 0xB ",
			want:  "0x000000000000000000000000000000000000000000000000000000000000000b",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := normalizeTransmitter(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	_, err := normalizeTransmitter("not-an-address")
	require.Error(t, err)
}

func TestP2PToTransmitterMapForWorkers(t *testing.T) {
	key := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1))
	workers := []*cre.NodeMetadata{
		{
			Keys: &secrets.NodeKeys{
				P2PKey: &crypto.P2PKey{
					PeerID: key.PeerID(),
				},
				Aptos: &crypto.AptosKey{
					Account: "0xa",
				},
			},
		},
	}

	got, err := p2pToTransmitterMapForWorkers(workers)
	require.NoError(t, err)

	peerID := key.PeerID()
	expectedPeerKey := hex.EncodeToString(peerID[:])
	require.Equal(t, map[string]string{
		expectedPeerKey: "0x000000000000000000000000000000000000000000000000000000000000000a",
	}, got)
}

func TestResolveMethodConfigSettings_Defaults(t *testing.T) {
	settings, err := resolveMethodConfigSettings(nil)
	require.NoError(t, err)
	require.Equal(t, defaultRequestTimeout, settings.RequestTimeout)
	require.Equal(t, defaultWriteDeltaStage, settings.DeltaStage)
	require.Equal(t, capabilitiespb.TransmissionSchedule_AllAtOnce, settings.TransmissionSchedule)
}

func TestResolveMethodConfigSettings_Overrides(t *testing.T) {
	settings, err := resolveMethodConfigSettings(map[string]any{
		requestTimeoutKey:       "45s",
		deltaStageKey:           "2500ms",
		transmissionScheduleKey: "oneAtATime",
	})
	require.NoError(t, err)
	require.Equal(t, 45*time.Second, settings.RequestTimeout)
	require.Equal(t, 2500*time.Millisecond, settings.DeltaStage)
	require.Equal(t, capabilitiespb.TransmissionSchedule_OneAtATime, settings.TransmissionSchedule)
}

func TestResolveMethodConfigSettings_InvalidDuration(t *testing.T) {
	_, err := resolveMethodConfigSettings(map[string]any{
		requestTimeoutKey: "not-a-duration",
	})
	require.Error(t, err)
}

func TestResolveMethodConfigSettings_InvalidTransmissionSchedule(t *testing.T) {
	_, err := resolveMethodConfigSettings(map[string]any{
		transmissionScheduleKey: "staggered",
	})
	require.Error(t, err)
}
