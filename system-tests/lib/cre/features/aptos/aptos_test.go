package aptos

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	aptoschain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/aptos"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
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

func TestSetForwarderAddress_UpdatesMatchingAptosChain(t *testing.T) {
	cfg := corechainlink.Config{
		Aptos: corechainlink.RawConfigs{
			corechainlink.RawConfig{
				"ChainID": "4",
				"Workflow": corechainlink.RawConfig{
					"ForwarderAddress": "0xold",
					"Keep":             "yes",
				},
			},
			corechainlink.RawConfig{
				"ChainID": "8",
			},
		},
	}

	err := setForwarderAddress(&cfg, "4", "0xnew")
	require.NoError(t, err)

	workflow := workflowMap(t, cfg.Aptos[0]["Workflow"])
	require.Equal(t, "0xnew", workflow["ForwarderAddress"])
	require.Equal(t, "yes", workflow["Keep"])
	require.Nil(t, cfg.Aptos[1]["Workflow"])
}

func TestPatchNodeTOML_PatchesAllMatchingAptosChainsAndPreservesWorkflowFields(t *testing.T) {
	baseConfig := corechainlink.Config{
		Aptos: corechainlink.RawConfigs{
			corechainlink.RawConfig{
				"ChainID": "4",
				"Workflow": corechainlink.RawConfig{
					"ForwarderAddress": "0xold4",
					"Keep":             "value4",
				},
			},
			corechainlink.RawConfig{
				"ChainID": "8",
				"Workflow": corechainlink.RawConfig{
					"ForwarderAddress": "0xold8",
					"Keep":             "value8",
				},
			},
		},
	}
	rawConfig, err := toml.Marshal(baseConfig)
	require.NoError(t, err)

	don := testDonMetadata(t, string(rawConfig), string(rawConfig))
	err = patchNodeTOML(don, map[uint64]string{
		4: "0x0000000000000000000000000000000000000000000000000000000000000004",
		8: "0x0000000000000000000000000000000000000000000000000000000000000008",
	})
	require.NoError(t, err)

	for _, spec := range don.MustNodeSet().NodeSpecs {
		var patched corechainlink.Config
		require.NoError(t, toml.Unmarshal([]byte(spec.Node.TestConfigOverrides), &patched))

		workflow4 := workflowMap(t, patched.Aptos[0]["Workflow"])
		require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", workflow4["ForwarderAddress"])
		require.Equal(t, "value4", workflow4["Keep"])

		workflow8 := workflowMap(t, patched.Aptos[1]["Workflow"])
		require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000008", workflow8["ForwarderAddress"])
		require.Equal(t, "value8", workflow8["Keep"])
	}
}

func TestFindAptosChainByChainID_ReturnsTypedBlockchain(t *testing.T) {
	aptosBC := testAptosBlockchain(4, 4457093679053095497)

	got, err := findAptosChainByChainID([]creblockchains.Blockchain{aptosBC}, 4)
	require.NoError(t, err)
	require.Same(t, aptosBC, got)
}

func TestFindAptosChainByChainID_ErrorsOnTypeMismatch(t *testing.T) {
	chains := []creblockchains.Blockchain{fakeChain{family: "aptos", chainID: 4}}

	_, err := findAptosChainByChainID(chains, 4)
	require.Error(t, err)
	require.ErrorContains(t, err, "unexpected type")
}

func TestBuildCapabilityRegistrations_UsesCapRegOCRConfig(t *testing.T) {
	don := testDonMetadataWithCapabilities(t, []string{"[Aptos]\n"}, []string{cre.AptosCapability + "-4"}, cre.CapabilityConfigs{
		cre.AptosCapability: {
			BinaryName: "aptos",
			Values: map[string]any{
				requestTimeoutKey:       "45s",
				deltaStageKey:           "2500ms",
				transmissionScheduleKey: "allAtOnce",
			},
		},
	})

	caps, capabilityToOCR3Config, capabilityLabels, err := buildCapabilityRegistrations(
		don,
		[]creblockchains.Blockchain{testAptosBlockchain(4, 4457093679053095497)},
		[]uint64{4},
		map[string]string{"peer-a": "0x1"},
	)
	require.NoError(t, err)
	require.Len(t, caps, 1)
	require.Len(t, capabilityLabels, 1)
	require.True(t, caps[0].UseCapRegOCRConfig)
	require.Equal(t, CapabilityLabel(4457093679053095497), caps[0].Capability.LabelledName)
	require.Contains(t, capabilityToOCR3Config, capabilityLabels[0])
	require.NotNil(t, capabilityToOCR3Config[capabilityLabels[0]])
}

func TestNewAptosWorkerJobInput_UsesCapRegVersion(t *testing.T) {
	creEnv := &cre.Environment{
		RegistryChainSelector: 111,
		ContractVersions: map[cre.ContractType]*semver.Version{
			keystone_changeset.CapabilitiesRegistry.String(): semver.MustParse("2.0.0"),
		},
	}

	input, err := newAptosWorkerJobInput(
		creEnv,
		"workflow-don",
		"/usr/local/bin/aptos",
		`{"chainId":"4"}`,
		[]string{"peer@127.0.0.1:6690"},
		222,
		4,
	)
	require.NoError(t, err)
	require.Equal(t, "aptos-worker-4", input.JobName)
	require.Equal(t, true, input.Inputs["useCapRegOCRConfig"])
	require.Equal(t, "2.0.0", input.Inputs["capRegVersion"])
	require.Equal(t, uint64(111), input.Inputs["chainSelectorEVM"])
	require.Equal(t, uint64(222), input.Inputs["chainSelectorAptos"])
	_, hasQualifier := input.Inputs["contractQualifier"]
	require.False(t, hasQualifier)
}

func TestNewAptosWorkerJobInput_ErrorsWithoutCapRegVersion(t *testing.T) {
	_, err := newAptosWorkerJobInput(
		&cre.Environment{},
		"workflow-don",
		"/usr/local/bin/aptos",
		`{"chainId":"4"}`,
		nil,
		222,
		4,
	)
	require.Error(t, err)
	require.ErrorContains(t, err, "CapabilitiesRegistry version not found")
}

func testDonMetadata(t *testing.T, nodeConfigs ...string) *cre.DonMetadata {
	return testDonMetadataWithCapabilities(t, nodeConfigs, nil, nil)
}

func testDonMetadataWithCapabilities(t *testing.T, nodeConfigs []string, capabilities []string, capabilityConfigs cre.CapabilityConfigs) *cre.DonMetadata {
	t.Helper()

	nodeSpecs := make([]*cre.NodeSpecWithRole, len(nodeConfigs))
	for i, cfg := range nodeConfigs {
		nodeSpecs[i] = &cre.NodeSpecWithRole{
			Input: &clnode.Input{
				Node: &clnode.NodeInput{TestConfigOverrides: cfg},
			},
			Roles: []cre.NodeType{cre.WorkerNode},
		}
	}

	nodeSet := &cre.NodeSet{
		Input:        &ns.Input{Name: "aptos-don"},
		NodeSpecs:    nodeSpecs,
		Capabilities: capabilities,
	}

	don, err := cre.NewDonMetadata(nodeSet, 1, infra.Provider{Type: infra.Docker}, capabilityConfigs)
	require.NoError(t, err)
	return don
}

func testAptosBlockchain(chainID, chainSelector uint64) *aptoschain.Blockchain {
	return aptoschain.NewBlockchain(
		zerolog.Nop(),
		chainID,
		chainSelector,
		&blockchain.Output{Family: "aptos", ChainID: "4"},
	)
}

func workflowMap(t *testing.T, raw any) map[string]any {
	t.Helper()

	switch v := raw.(type) {
	case corechainlink.RawConfig:
		return map[string]any(v)
	case map[string]any:
		return v
	default:
		t.Fatalf("unexpected workflow type %T", raw)
		return nil
	}
}

type fakeChain struct {
	family  string
	chainID uint64
}

func (f fakeChain) ChainSelector() uint64 { return 0 }
func (f fakeChain) ChainID() uint64       { return f.chainID }
func (f fakeChain) ChainFamily() string   { return f.family }
func (f fakeChain) IsFamily(chainFamily string) bool {
	return f.family == chainFamily
}
func (f fakeChain) Fund(context.Context, string, uint64) error { return nil }
func (f fakeChain) CtfOutput() *blockchain.Output              { return nil }
func (f fakeChain) ToCldfChain() (cldfchain.BlockChain, error) { return nil, nil }
