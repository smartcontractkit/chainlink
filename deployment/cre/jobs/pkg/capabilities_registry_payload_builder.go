package pkg

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

// BuildCapabilitiesRegistryPayload assembles the proto-JSON offchain capabilities registry payload
// (what CapabilitiesRegistryJob.OffchainConfig carries) from DON name -> capability ID -> config. DON
// names are the onchain DON names. Versions must strictly increase for an environment; the caller owns
// that. Chain capability entries are checked with ValidateChainCapabilityEntries, so a malformed one
// fails here rather than on the nodes. The payload is marshalled with the same plain protojson the CLD
// compiler uses.
func BuildCapabilitiesRegistryPayload(domain, env string, version uint64, dons map[string]map[string]*capabilitiespb.CapabilityConfig) (string, error) {
	reg := &capabilitiespb.OffchainCapabilitiesRegistry{
		Domain:  domain,
		Env:     env,
		Version: version,
		Dons:    make(map[string]*capabilitiespb.DONConfig, len(dons)),
	}
	for donName, caps := range dons {
		reg.Dons[donName] = &capabilitiespb.DONConfig{Capabilities: caps}
	}

	out, err := protojson.Marshal(reg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal offchain capabilities registry: %w", err)
	}
	if err = ValidateChainCapabilityEntries(string(out)); err != nil {
		return "", err
	}
	return string(out), nil
}
