package stellar

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
)

const (
	requestTimeoutKey      = "RequestTimeout"
	deltaStageKey          = "DeltaStage"
	defaultRequestTimeout  = 30 * time.Second
	defaultWriteDeltaStage = 15 * time.Second
)

type methodConfigSettings struct {
	RequestTimeout time.Duration
	DeltaStage     time.Duration
}

func CapabilityLabel(chainSelector uint64) string {
	return "stellar:ChainSelector:" + strconv.FormatUint(chainSelector, 10)
}

func methodConfigs(settings methodConfigSettings) map[string]*capabilitiespb.CapabilityMethodConfig {
	m := make(map[string]*capabilitiespb.CapabilityMethodConfig)
	for _, action := range []string{"GetLatestLedger", "ReadContract"} {
		m[action] = readActionConfig(settings.RequestTimeout)
	}
	m["WriteReport"] = writeReportActionConfig(settings)
	return m
}

func readActionConfig(requestTimeout time.Duration) *capabilitiespb.CapabilityMethodConfig {
	return &capabilitiespb.CapabilityMethodConfig{
		RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
			RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
				RequestTimeout:            durationpb.New(requestTimeout),
				ServerMaxParallelRequests: 10,
				RequestHasherType:         capabilitiespb.RequestHasherType_Simple,
			},
		},
	}
}

func writeReportActionConfig(settings methodConfigSettings) *capabilitiespb.CapabilityMethodConfig {
	return &capabilitiespb.CapabilityMethodConfig{
		RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
			RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
				TransmissionSchedule:      capabilitiespb.TransmissionSchedule_OneAtATime,
				DeltaStage:                durationpb.New(settings.DeltaStage),
				RequestTimeout:            durationpb.New(settings.RequestTimeout),
				ServerMaxParallelRequests: 10,
				RequestHasherType:         capabilitiespb.RequestHasherType_WriteReportExcludeSignatures,
			},
		},
	}
}

func resolveMethodConfigSettings(values map[string]any) (methodConfigSettings, error) {
	settings := methodConfigSettings{
		RequestTimeout: defaultRequestTimeout,
		DeltaStage:     defaultWriteDeltaStage,
	}
	if values == nil {
		return settings, nil
	}

	requestTimeout, ok, err := durationValue(values, requestTimeoutKey)
	if err != nil {
		return methodConfigSettings{}, err
	}
	if ok {
		settings.RequestTimeout = requestTimeout
	}

	deltaStage, ok, err := durationValue(values, deltaStageKey)
	if err != nil {
		return methodConfigSettings{}, err
	}
	if ok {
		settings.DeltaStage = deltaStage
	}

	return settings, nil
}

func buildJobConfigJSON(
	chainID string,
	forwarderAddress string,
	settings methodConfigSettings,
	isLocal bool,
) (string, error) {
	// Reuse the deployment changeset's typed config schema so the Local CRE job
	// config stays a single source of truth with the production
	// ProposeStellarCapJobSpec, instead of a hand-rolled map[string]any.
	cfg := jobs.StellarOverrideDefaultCfg{
		Network:             network,
		ChainID:             chainID,
		CREForwarderAddress: forwarderAddress,
		IsLocal:             isLocal,
		DeltaStage:          settings.DeltaStage,
	}

	raw, err := json.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Stellar standard capability job config: %w", err)
	}
	return string(raw), nil
}

// buildSpecConfig returns the Stellar capability config as registry spec_config. It carries
// the same keys the worker job spec used to supply, so the node can launch the capability
// from the registry alone.
func buildSpecConfig(
	chainID string,
	forwarderAddress string,
	settings methodConfigSettings,
	isLocal bool,
) (*valuespb.Map, error) {
	raw, err := buildJobConfigJSON(chainID, forwarderAddress, settings, isLocal)
	if err != nil {
		return nil, err
	}

	var cfg map[string]any
	if err = json.Unmarshal([]byte(raw), &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode Stellar config: %w", err)
	}
	for key, v := range cfg {
		// encoding/json decodes every number as float64; keep integral ones integral.
		if f, ok := v.(float64); ok && f == math.Trunc(f) {
			cfg[key] = int64(f)
		}
	}

	wrapped, err := values.NewMap(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap Stellar spec config: %w", err)
	}
	return values.ProtoMap(wrapped), nil
}

func durationValue(values map[string]any, key string) (time.Duration, bool, error) {
	raw, ok := values[key]
	if !ok {
		return 0, false, nil
	}

	switch v := raw.(type) {
	case string:
		parsed, err := time.ParseDuration(strings.TrimSpace(v))
		if err != nil {
			return 0, false, fmt.Errorf("%s must be a valid duration string: %w", key, err)
		}
		return parsed, true, nil
	case time.Duration:
		return v, true, nil
	default:
		return 0, false, fmt.Errorf("%s must be a duration string, got %T", key, raw)
	}
}
