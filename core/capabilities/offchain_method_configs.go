package capabilities

import (
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
)

// offchainMethodConfigs is a per-pass snapshot of the offchain registry's method_configs,
// keyed by on-chain DON name. Nil when the cutover gate is off or no payload is applied, in
// which case selection retains the on-chain config and behavior is byte-for-byte the legacy one.
type offchainMethodConfigs struct {
	reg      *capabilitiespb.OffchainCapabilitiesRegistry
	donNames map[uint32]string // on-chain DON ID -> offchain DON name
}

// newOffchainMethodConfigs builds the snapshot from the already-loaded offchain registry and
// the DON set. DON names are resolved the same way the LocalCapabilityManager does (unnamed or
// name-colliding DONs are left out, so they keep their on-chain method configs).
func newOffchainMethodConfigs(reg *capabilitiespb.OffchainCapabilitiesRegistry, dons []registry.DON) *offchainMethodConfigs {
	if reg == nil {
		return nil
	}
	byName := make(map[string][]uint32, len(dons))
	for _, don := range dons {
		if don.Name == "" {
			continue
		}
		byName[don.Name] = append(byName[don.Name], don.ID)
	}
	donNames := make(map[uint32]string, len(dons))
	for name, ids := range byName {
		if len(ids) == 1 {
			donNames[ids[0]] = name
		}
	}
	return &offchainMethodConfigs{reg: reg, donNames: donNames}
}

// configuration selects the complete capability config before decoding it. A missing
// offchain capability retains the on-chain config, including its LocalOnly setting.
func (o *offchainMethodConfigs) configuration(capID string, don registry.DON, onchain registry.CapabilityConfiguration) (capabilities.CapabilityConfiguration, error) {
	if o == nil || o.reg == nil {
		return onchain.Unmarshal()
	}
	donName, ok := o.donNames[don.ID]
	if !ok {
		return onchain.Unmarshal()
	}
	offCap, ok := o.reg.GetDons()[donName].GetCapabilities()[capID]
	if !ok {
		return onchain.Unmarshal()
	}
	raw, err := proto.Marshal(offCap)
	if err != nil {
		return capabilities.CapabilityConfiguration{}, err
	}
	config, err := (registry.CapabilityConfiguration{Config: raw}).Unmarshal()
	// A present offchain capability without methods explicitly selects no methods.
	if err == nil && config.CapabilityMethodConfig == nil {
		config.CapabilityMethodConfig = map[string]capabilities.CapabilityMethodConfig{}
	}
	return config, err
}
