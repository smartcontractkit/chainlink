package conversions

import (
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"

	chainselectors "github.com/smartcontractkit/chain-selectors"
)

// WARNING: Hacky and brittle - used during the current std-capability transition to map
// job commands back to capability IDs. The standard-capability delegate still needs this
// for registry-based launch allowlisting and OCR config wiring, including remote caps like Aptos.
func GetCapabilityIDFromCommand(command string, config string) string {
	switch filepath.Base(command) {
	case "evm":
		var cfg struct {
			ChainID uint64 `json:"chainId"`
		}
		if err := json.Unmarshal([]byte(config), &cfg); err != nil {
			return ""
		}
		selector, err := chainselectors.SelectorFromChainId(cfg.ChainID)
		if err != nil {
			return ""
		}
		return "evm:ChainSelector:" + strconv.FormatUint(selector, 10) + "@1.0.0"
	case "aptos":
		var cfg struct {
			ChainID string `json:"chainId"`
		}
		if err := json.Unmarshal([]byte(config), &cfg); err != nil {
			return ""
		}
		chainID, err := strconv.ParseUint(cfg.ChainID, 10, 64)
		if err != nil {
			return ""
		}
		selector, ok := chainselectors.AptosChainIdToChainSelector()[chainID]
		if !ok {
			return ""
		}
		return "aptos:ChainSelector:" + strconv.FormatUint(selector, 10) + "@1.0.0"
	case "consensus":
		return "consensus@1.0.0-alpha"
	case "cron":
		return "cron-trigger@1.0.0"
	case "http_trigger":
		return "http-trigger@1.0.0-alpha"
	case "http_action":
		return "http-actions@1.0.0-alpha" // plural "actions"
	default:
		return ""
	}
}

// GetCommandFromCapabilityID is the inverse of GetCapabilityIDFromCommand: it maps a capability ID
// back to the base command name (e.g. "evm", "consensus"). Returns "" for unrecognized IDs.
func GetCommandFromCapabilityID(capabilityID string) string {
	switch {
	case strings.HasPrefix(capabilityID, "evm"):
		return "evm"
	case strings.HasPrefix(capabilityID, "aptos:ChainSelector:"):
		return "aptos"
	case strings.HasPrefix(capabilityID, "consensus"):
		return "consensus"
	case strings.HasPrefix(capabilityID, "cron-trigger"):
		return "cron"
	case strings.HasPrefix(capabilityID, "http-trigger"):
		return "http_trigger"
	case strings.HasPrefix(capabilityID, "http-actions"):
		return "http_action"
	default:
		return ""
	}
}
