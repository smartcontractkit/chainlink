package onchain

import (
	stderrors "errors"
	"fmt"
	"maps"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
)

func validateDiscoveredEVMAddresses(
	donName string,
	nodeNames []string,
	requiredChainIDs []uint64,
	runtime map[string]domain.NodeRuntimeInfo,
) error {
	if len(requiredChainIDs) == 0 {
		return nil
	}

	var errs []error
	for _, nodeName := range nodeNames {
		info, ok := runtime[nodeName]
		if !ok {
			errs = append(errs, fmt.Errorf("DON %s node %s: no discovered runtime info", donName, nodeName))
			continue
		}
		for _, chainID := range requiredChainIDs {
			chainKey := strconv.FormatUint(chainID, 10)
			if info.EVMAddress == nil {
				errs = append(errs, fmt.Errorf("DON %s node %s: missing EVM address for chain %d", donName, nodeName, chainID))
				continue
			}
			addr, ok := info.EVMAddress[chainKey]
			if !ok || strings.TrimSpace(addr) == "" {
				errs = append(errs, fmt.Errorf("DON %s node %s: missing EVM address for chain %d", donName, nodeName, chainID))
			}
		}
	}
	return stderrors.Join(errs...)
}

func hydrateDiscoveredEVMAddresses(
	donMeta *cre.DonMetadata,
	nodeNames []string,
	requiredChainIDs []uint64,
	runtime map[string]domain.NodeRuntimeInfo,
) error {
	if len(requiredChainIDs) == 0 {
		return nil
	}

	var errs []error
	for i, nodeMeta := range donMeta.NodesMetadata {
		if i >= len(nodeNames) {
			errs = append(errs, fmt.Errorf("DON %s: missing node name for metadata index %d", donMeta.Name, i))
			continue
		}
		nodeName := nodeNames[i]
		if nodeMeta.Keys == nil {
			errs = append(errs, fmt.Errorf("DON %s node %s: node keys not initialized", donMeta.Name, nodeName))
			continue
		}

		nodeRuntime, ok := runtime[nodeName]
		if !ok {
			errs = append(errs, fmt.Errorf("DON %s node %s: no discovered runtime info", donMeta.Name, nodeName))
			continue
		}

		for _, chainID := range requiredChainIDs {
			chainKey := strconv.FormatUint(chainID, 10)
			if nodeRuntime.EVMAddress == nil {
				errs = append(errs, fmt.Errorf("DON %s node %s: missing discovered EVM address for chain %d", donMeta.Name, nodeName, chainID))
				continue
			}
			addr, ok := nodeRuntime.EVMAddress[chainKey]
			if !ok || strings.TrimSpace(addr) == "" {
				errs = append(errs, fmt.Errorf("DON %s node %s: missing discovered EVM address for chain %d", donMeta.Name, nodeName, chainID))
				continue
			}

			evmKey := nodeMeta.Keys.EVM[chainID]
			if evmKey == nil {
				evmKey = &crypto.EVMKey{}
				nodeMeta.Keys.EVM[chainID] = evmKey
			}
			evmKey.PublicAddress = common.HexToAddress(addr)
		}
	}
	return stderrors.Join(errs...)
}

// hydrateDiscoveredOCR2BundleIDs copies OCR2 key bundle IDs discovered during
// D1 (per chain family, e.g. "evm") into each node's Keys.OCR2BundleIDs, which
// Features look up directly (e.g. evm.go's per-node OCR3 job proposal).
func hydrateDiscoveredOCR2BundleIDs(
	donMeta *cre.DonMetadata,
	nodeNames []string,
	runtime map[string]domain.NodeRuntimeInfo,
) error {
	var errs []error
	for i, nodeMeta := range donMeta.NodesMetadata {
		if i >= len(nodeNames) {
			errs = append(errs, fmt.Errorf("DON %s: missing node name for metadata index %d", donMeta.Name, i))
			continue
		}
		nodeName := nodeNames[i] //nolint:gosec // bounds-checked by the `i >= len(nodeNames)` guard above
		if nodeMeta.Keys == nil {
			errs = append(errs, fmt.Errorf("DON %s node %s: node keys not initialized", donMeta.Name, nodeName))
			continue
		}

		nodeRuntime, ok := runtime[nodeName]
		if !ok {
			errs = append(errs, fmt.Errorf("DON %s node %s: no discovered runtime info", donMeta.Name, nodeName))
			continue
		}
		if len(nodeRuntime.OCR2BundleIDs) == 0 {
			errs = append(errs, fmt.Errorf("DON %s node %s: no discovered OCR2 key bundles", donMeta.Name, nodeName))
			continue
		}

		if nodeMeta.Keys.OCR2BundleIDs == nil {
			nodeMeta.Keys.OCR2BundleIDs = make(map[string]string)
		}
		maps.Copy(nodeMeta.Keys.OCR2BundleIDs, nodeRuntime.OCR2BundleIDs)
	}
	return stderrors.Join(errs...)
}

// hydrateDiscoveredHosts overwrites each node's Keys-adjacent Host field
// (system-tests/lib's NewDonMetadata sets it to a synthetic
// "<donName>-bt-<index>" convention meant for CTFv2's own docker-compose/k8s
// provisioning) with the real in-cluster DNS name for the Griddle-deployed
// chart node. cre.PeeringCfgs (used to build OCR2/OCR3 job specs'
// bootstrapperOCR3Urls) reads Host directly, so leaving the synthetic value in
// place makes nodes dial a hostname that was never actually deployed.
func hydrateDiscoveredHosts(
	donMeta *cre.DonMetadata,
	nodeNames []string,
	cv *domain.ChartValues,
) error {
	var errs []error
	for i, nodeMeta := range donMeta.NodesMetadata {
		if i >= len(nodeNames) {
			errs = append(errs, fmt.Errorf("DON %s: missing node name for metadata index %d", donMeta.Name, i))
			continue
		}
		nodeMeta.Host = cv.NodeInternalHost(nodeNames[i]) //nolint:gosec // bounds-checked by the `i >= len(nodeNames)` guard above
	}
	return stderrors.Join(errs...)
}
