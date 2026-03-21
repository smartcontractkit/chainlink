package standardcapabilities

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

type aliasedBaseCapability struct {
	capabilities.BaseCapability
	info capabilities.CapabilityInfo
}

func (a aliasedBaseCapability) Info(context.Context) (capabilities.CapabilityInfo, error) {
	return a.info, nil
}

type aliasedExecutableCapability struct {
	capabilities.ExecutableCapability
	info capabilities.CapabilityInfo
}

func (a aliasedExecutableCapability) Info(context.Context) (capabilities.CapabilityInfo, error) {
	return a.info, nil
}

type aliasedTriggerCapability struct {
	capabilities.TriggerCapability
	info capabilities.CapabilityInfo
}

func (a aliasedTriggerCapability) Info(context.Context) (capabilities.CapabilityInfo, error) {
	return a.info, nil
}

type aliasedExecutableAndTriggerCapability struct {
	capabilities.ExecutableAndTriggerCapability
	info capabilities.CapabilityInfo
}

func (a aliasedExecutableAndTriggerCapability) Info(context.Context) (capabilities.CapabilityInfo, error) {
	return a.info, nil
}

func aliasCapabilityID(baseCap capabilities.BaseCapability, overrideInfo capabilities.CapabilityInfo) (capabilities.BaseCapability, error) {
	switch c := baseCap.(type) {
	case capabilities.ExecutableAndTriggerCapability:
		return aliasedExecutableAndTriggerCapability{ExecutableAndTriggerCapability: c, info: overrideInfo}, nil
	case capabilities.TriggerCapability:
		return aliasedTriggerCapability{TriggerCapability: c, info: overrideInfo}, nil
	case capabilities.ExecutableCapability:
		return aliasedExecutableCapability{ExecutableCapability: c, info: overrideInfo}, nil
	case capabilities.BaseCapability:
		return aliasedBaseCapability{BaseCapability: c, info: overrideInfo}, nil
	default:
		return nil, fmt.Errorf("unsupported capability type %T for alias %q", baseCap, overrideInfo.ID)
	}
}
