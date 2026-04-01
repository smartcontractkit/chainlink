package modifier

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"

	"google.golang.org/protobuf/encoding/protojson"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
)

// CapabilityConfigModifierParams carries shared inputs for per-DON capability config modification.
// Extend with new fields when additional modifiers need them; modifiers ignore unused fields.
type CapabilityConfigModifierParams struct {
	Env     *cldf.Environment
	DonName string
	P2PIDs  []p2pkey.PeerID
	// Configs is a per-DON clone of the caller's capability configs; modifiers mutate in place.
	Configs []contracts.CapabilityConfig
}

// CapabilityConfigModifier applies chain or capability-specific changes to Config (e.g. specConfig).
type CapabilityConfigModifier interface {
	Modify(params CapabilityConfigModifierParams) error
}

func DefaultCapabilityConfigModifiers() []CapabilityConfigModifier {
	return []CapabilityConfigModifier{
		aptosDonModifier{},
	}
}

// aptos

type aptosDonModifier struct{}

func (aptosDonModifier) Modify(params CapabilityConfigModifierParams) error {
	for i := range params.Configs {
		sel, isAptos, parseErr := parseAptosChainSelectorFromCapabilityID(params.Configs[i].Capability.CapabilityID)
		if parseErr != nil {
			return fmt.Errorf("capability %q: %w", params.Configs[i].Capability.CapabilityID, parseErr)
		}
		if !isAptos {
			continue
		}
		if params.Env == nil || params.Env.Offchain == nil {
			return errors.New("AddCapabilities: Aptos capabilities require Env.Offchain (Job Distributor client)")
		}
		if params.Configs[i].Config == nil {
			params.Configs[i].Config = make(map[string]any)
		}
		p2pMap, mapErr := buildAptosP2PToTransmitterMap(params.Env.Offchain, params.P2PIDs, sel)
		if mapErr != nil {
			return fmt.Errorf("capability %q: %w", params.Configs[i].Capability.CapabilityID, mapErr)
		}
		if mergeErr := mergeAptosP2PToTransmitterIntoConfig(params.Configs[i].Config, p2pMap); mergeErr != nil {
			return fmt.Errorf("capability %q: %w", params.Configs[i].Capability.CapabilityID, mergeErr)
		}
	}
	return nil
}

// aptosCapabilityIDPrefix is the capability id form used for Aptos chain capabilities
// (label before optional "@<version>"), e.g. aptos:ChainSelector:12345@1.0.0.
const aptosCapabilityIDPrefix = "aptos:ChainSelector:"

// parseAptosChainSelectorFromCapabilityID parses registry capability IDs of the form
// aptos:ChainSelector:<decimal>@<version>. The part after the last "@" is ignored for
// parsing so only the label matters (e.g. aptos:ChainSelector:12345@1.0.0 → 12345).
//
// Returns isAptos false and no error when the id does not start with aptosCapabilityIDPrefix
// (after stripping "@…"). Returns isAptos true and an error if the prefix is present but the
// selector is empty or not a base-10 uint64.
func parseAptosChainSelectorFromCapabilityID(capabilityID string) (selector uint64, isAptos bool, err error) {
	capID := capabilityID
	if i := strings.LastIndex(capabilityID, "@"); i >= 0 {
		capID = capabilityID[:i]
	}
	if !strings.HasPrefix(capID, aptosCapabilityIDPrefix) {
		return 0, false, nil
	}
	raw := strings.TrimPrefix(capID, aptosCapabilityIDPrefix)
	if raw == "" {
		return 0, true, fmt.Errorf("missing chain selector in capability id %q", capabilityID)
	}
	u, parseErr := strconv.ParseUint(raw, 10, 64)
	if parseErr != nil {
		return 0, true, fmt.Errorf("invalid chain selector in capability id %q: %w", capabilityID, parseErr)
	}
	return u, true, nil
}

// buildAptosP2PToTransmitterMap asks Job Distributor for node metadata for donPeerIDs
// Then builds a map used for CapabilityConfig spec:
// - lowercase hex of the 32-byte P2P id -> Aptos transmit account (OCR TransmitAccount)
// for the given aptosChainSelector.
//
// It walks only the nodes returned by NodeInfo. Each must have OCR config for
// aptosChainSelector and a non-empty transmit account after trim, or this returns an error.
func buildAptosP2PToTransmitterMap(
	offChainClient deployment.NodeChainConfigsLister,
	donPeerIDs []p2pkey.PeerID,
	aptosChainSelector uint64,
) (map[string]string, error) {
	if offChainClient == nil {
		return nil, errors.New("offchain client is nil")
	}
	if len(donPeerIDs) == 0 {
		return nil, errors.New("no DON peer IDs")
	}
	p2pStrs := make([]string, len(donPeerIDs))
	for i, pid := range donPeerIDs {
		p2pStrs[i] = pid.String()
	}
	nodes, nodeInfoErr := deployment.NodeInfo(p2pStrs, offChainClient)
	if nodeInfoErr != nil {
		return nil, fmt.Errorf("failed to get node info from JD: %w", nodeInfoErr)
	}
	out := make(map[string]string, len(nodes))
	for _, node := range nodes {
		ocrCfg, ok := node.OCRConfigForChainSelector(aptosChainSelector)
		if !ok {
			return nil, fmt.Errorf("node %s (%s) has no OCR2 config for chain selector %d",
				node.Name, node.PeerID.String(), aptosChainSelector)
		}
		transmitter := strings.TrimSpace(string(ocrCfg.TransmitAccount))
		if transmitter == "" {
			return nil, fmt.Errorf("empty Aptos transmit account for node %s (%s)", node.Name, node.PeerID.String())
		}
		out[hex.EncodeToString(node.PeerID[:])] = transmitter
	}
	return out, nil
}

// mergeAptosP2PToTransmitterIntoConfig sets cfg["specConfig"] to p2pMap (as p2pToTransmitterMap).
// Caller must omit specConfig or leave it empty; any non-empty specConfig returns an error for now.
// NOTE: we can make this smarter later if needed. Add overwriting / merging logic etc.
//
// specConfig is protobuf values.v1.Map JSON; we build it with values.Wrap so pkg.MarshalProto succeeds.
func mergeAptosP2PToTransmitterIntoConfig(cfg map[string]any, p2pMap map[string]string) error {
	if cfg == nil {
		return errors.New("nil capability config map")
	}
	if raw, ok := cfg["specConfig"]; ok && raw != nil {
		if !isEmptySpecConfigForAptosMerge(raw) {
			return errors.New("specConfig must be empty (omit or {}) for Aptos p2pToTransmitterMap injection")
		}
	}
	p2pVal, err := values.Wrap(p2pMap)
	if err != nil {
		return fmt.Errorf("wrap p2pToTransmitterMap: %w", err)
	}
	spec := values.EmptyMap()
	spec.Underlying["p2pToTransmitterMap"] = p2pVal
	out, err := protojson.Marshal(values.ProtoMap(spec))
	if err != nil {
		return fmt.Errorf("marshal specConfig: %w", err)
	}
	var specAsMap map[string]any
	if err := json.Unmarshal(out, &specAsMap); err != nil {
		return fmt.Errorf("specConfig map: %w", err)
	}
	cfg["specConfig"] = specAsMap
	return nil
}

// isEmptySpecConfigForAptosMerge reports whether user-provided specConfig is absent-equivalent:
// nil, {}, or values.v1.Map JSON with no entries ({ "fields": {} }).
func isEmptySpecConfigForAptosMerge(raw any) bool {
	if raw == nil {
		return true
	}
	m, ok := raw.(map[string]any)
	if !ok {
		return false
	}
	if len(m) == 0 {
		return true
	}
	if len(m) == 1 {
		fields, ok := m["fields"].(map[string]any)
		if ok && len(fields) == 0 {
			return true
		}
	}
	return false
}
