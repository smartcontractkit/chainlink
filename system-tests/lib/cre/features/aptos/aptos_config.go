package aptos

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/durationpb"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
)

// BuildCapabilityConfig builds the Aptos capability config passed directly
// through the capability manager: method execution policy in MethodConfigs and
// Aptos-specific runtime inputs in SpecConfig.
func BuildCapabilityConfig(values map[string]any, p2pToTransmitterMap map[string]string, localOnly bool) (*capabilitiespb.CapabilityConfig, error) {
	methodSettings, err := resolveMethodConfigSettings(values)
	if err != nil {
		return nil, err
	}

	capConfig := &capabilitiespb.CapabilityConfig{
		MethodConfigs: methodConfigs(methodSettings),
		LocalOnly:     localOnly,
	}
	if err := setRuntimeSpecConfig(capConfig, methodSettings, p2pToTransmitterMap); err != nil {
		return nil, err
	}
	return capConfig, nil
}

func buildWorkerConfigJSON(chainID uint64, forwarderAddress string, settings methodConfigSettings, p2pToTransmitterMap map[string]string, isLocal bool) (string, error) {
	cfg := map[string]any{
		"chainId":             strconv.FormatUint(chainID, 10),
		"network":             "aptos",
		"creForwarderAddress": forwarderAddress,
		"isLocal":             isLocal,
		"deltaStage":          settings.DeltaStage,
	}
	if len(p2pToTransmitterMap) > 0 {
		cfg[specConfigP2PMapKey] = p2pToTransmitterMap
	}

	raw, err := json.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Aptos worker config: %w", err)
	}
	return string(raw), nil
}

func methodConfigs(settings methodConfigSettings) map[string]*capabilitiespb.CapabilityMethodConfig {
	return map[string]*capabilitiespb.CapabilityMethodConfig{
		"View": {
			RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
				RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
					TransmissionSchedule:      capabilitiespb.TransmissionSchedule_AllAtOnce,
					RequestTimeout:            durationpb.New(settings.RequestTimeout),
					ServerMaxParallelRequests: 10,
					RequestHasherType:         capabilitiespb.RequestHasherType_Simple,
				},
			},
		},
		"WriteReport": {
			RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
				RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
					TransmissionSchedule:      settings.TransmissionSchedule,
					DeltaStage:                durationpb.New(settings.DeltaStage),
					RequestTimeout:            durationpb.New(settings.RequestTimeout),
					ServerMaxParallelRequests: 10,
					RequestHasherType:         capabilitiespb.RequestHasherType_WriteReportExcludeSignatures,
				},
			},
		},
	}
}

func resolveMethodConfigSettings(values map[string]any) (methodConfigSettings, error) {
	settings := methodConfigSettings{
		RequestTimeout:       defaultRequestTimeout,
		DeltaStage:           defaultWriteDeltaStage,
		TransmissionSchedule: capabilitiespb.TransmissionSchedule_AllAtOnce,
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

	transmissionSchedule, ok, err := transmissionScheduleValue(values, transmissionScheduleKey)
	if err != nil {
		return methodConfigSettings{}, err
	}
	if ok {
		settings.TransmissionSchedule = transmissionSchedule
	}

	return settings, nil
}

func transmissionScheduleValue(values map[string]any, key string) (capabilitiespb.TransmissionSchedule, bool, error) {
	raw, ok := values[key]
	if !ok {
		return 0, false, nil
	}

	schedule, ok := raw.(string)
	if !ok {
		return 0, false, fmt.Errorf("%s must be a string, got %T", key, raw)
	}

	switch strings.TrimSpace(schedule) {
	case "allAtOnce":
		return capabilitiespb.TransmissionSchedule_AllAtOnce, true, nil
	case "oneAtATime":
		return capabilitiespb.TransmissionSchedule_OneAtATime, true, nil
	default:
		return 0, false, fmt.Errorf("%s must be allAtOnce or oneAtATime, got %q", key, schedule)
	}
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

func setRuntimeSpecConfig(capConfig *capabilitiespb.CapabilityConfig, settings methodConfigSettings, p2pToTransmitterMap map[string]string) error {
	if capConfig == nil {
		return pkgerrors.New("capability config is nil")
	}

	specConfig, err := values.FromMapValueProto(capConfig.SpecConfig)
	if err != nil {
		return fmt.Errorf("failed to decode existing spec config: %w", err)
	}
	if specConfig == nil {
		specConfig = values.EmptyMap()
	}

	delete(specConfig.Underlying, legacyTransmittersKey)

	scheduleValue, err := values.Wrap(remoteTransmissionScheduleString(settings.TransmissionSchedule))
	if err != nil {
		return fmt.Errorf("failed to wrap transmission schedule: %w", err)
	}
	specConfig.Underlying[specConfigScheduleKey] = scheduleValue

	deltaStageValue, err := values.Wrap(settings.DeltaStage)
	if err != nil {
		return fmt.Errorf("failed to wrap delta stage: %w", err)
	}
	specConfig.Underlying[specConfigDeltaStageKey] = deltaStageValue

	if len(p2pToTransmitterMap) > 0 {
		mapValue, err := values.Wrap(p2pToTransmitterMap)
		if err != nil {
			return fmt.Errorf("failed to wrap p2p transmitter map: %w", err)
		}
		specConfig.Underlying[specConfigP2PMapKey] = mapValue
	}

	capConfig.SpecConfig = values.ProtoMap(specConfig)
	return nil
}

func remoteTransmissionScheduleString(schedule capabilitiespb.TransmissionSchedule) string {
	switch schedule {
	case capabilitiespb.TransmissionSchedule_OneAtATime:
		return "oneAtATime"
	default:
		return "allAtOnce"
	}
}
