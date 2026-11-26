package pkg

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

// Chain capabilities (EVM, Aptos, Stellar) launch from the offchain capabilities registry payload
// delivered by CapabilitiesRegistryJob. The checks below are for the author of that payload to run
// before delivery; they are not applied by ResolveJob so existing callers keep their behavior, and
// they never touch onchain data. They only parse the payload strictly; the node's structural
// validation (globalconfig.Validate in core) is not imported here, to keep chainlink/v2 out of
// this package's dependencies. Run the caller's own structural validation alongside them.

const solanaFamily = "solana"

// chainCapabilityIDPattern matches the registry launch chain capability IDs, for example
// "evm:ChainSelector:3379446385462418246@1.0.0".
var chainCapabilityIDPattern = regexp.MustCompile(`^(evm|aptos|stellar):ChainSelector:[0-9]+@[^@\s]+$`)

// requiredChainSpecConfigKeys are the spec_config keys every chain capability needs to launch from
// the registry alone: the chain, the network and the CRE forwarder it writes reports to.
var requiredChainSpecConfigKeys = []string{"chainId", "network", "creForwarderAddress"}

func isChainFamilyID(capID string) bool {
	for _, family := range []string{"evm", "aptos", "stellar", solanaFamily} {
		if strings.HasPrefix(capID, family+":") {
			return true
		}
	}
	return false
}

// ValidateChainCapabilityEntries checks the chain capability entries of an offchain capabilities
// registry payload (proto-JSON). For every evm:/aptos:/stellar: capability in every DON it
// requires:
//   - a well-formed capability ID,
//   - localOnly not true: the onchain value is false, and the payload must match it. Omitted means
//     false (protojson drops false when it marshals), so an explicit "localOnly": false is not
//     required,
//   - a spec_config with chainId, network and creForwarderAddress.
//
// solana: entries are refused: Solana stays on job specs. Entries of other capabilities are not
// checked. All problems are reported together.
func ValidateChainCapabilityEntries(raw string) error {
	reg, err := parsePayload(raw)
	if err != nil {
		return err
	}

	var problems []error
	for _, donName := range sortedKeys(reg.GetDons()) {
		don := reg.GetDons()[donName]
		for _, capID := range sortedKeys(don.GetCapabilities()) {
			if !isChainFamilyID(capID) {
				continue
			}
			where := fmt.Sprintf("dons[%q].capabilities[%q]", donName, capID)
			if strings.HasPrefix(capID, solanaFamily+":") {
				problems = append(problems, fmt.Errorf("%s: solana stays on job specs and cannot be launched from the registry", where))
				continue
			}
			if !chainCapabilityIDPattern.MatchString(capID) {
				problems = append(problems, fmt.Errorf("%s: malformed chain capability ID, want <evm|aptos|stellar>:ChainSelector:<n>@<version>", where))
				continue
			}
			capCfg := don.GetCapabilities()[capID]
			if capCfg.GetLocalOnly() {
				problems = append(problems, fmt.Errorf("%s: localOnly must be false to match the onchain value", where))
			}
			fields := capCfg.GetSpecConfig().GetFields()
			for _, key := range requiredChainSpecConfigKeys {
				if _, ok := fields[key]; !ok {
					problems = append(problems, fmt.Errorf("%s: spec_config is missing %q", where, key))
				}
			}
		}
	}
	return errors.Join(problems...)
}

// parsePayload strictly decodes a payload: an unknown field is an error, as on the node.
func parsePayload(raw string) (*capabilitiespb.OffchainCapabilitiesRegistry, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("offchain config: empty payload")
	}
	var reg capabilitiespb.OffchainCapabilitiesRegistry
	if err := protojson.Unmarshal([]byte(raw), &reg); err != nil {
		return nil, fmt.Errorf("offchain config: invalid payload: %w", err)
	}
	return &reg, nil
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
