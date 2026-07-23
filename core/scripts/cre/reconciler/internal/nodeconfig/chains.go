package nodeconfig

import (
	"fmt"
	"os"
	"sort"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

// managedLayerName is the configuration layer reconciler itself writes
// (see infra.PatchChartValues). It never contains user-authored [[EVM]] chain
// config, so it's skipped when discovering chains from the chart.
const managedLayerName = "30-cre"

// DiscoverChains reads each chart node's existing "configuration" TOML layers
// (chainlink-node.instances.<name>.configuration in the chart YAML) and
// extracts declared EVM chains, for prepopulating the UI's Chains tab.
//
// A chain ID can legitimately appear with more than one (ws_url, http_url)
// pair across nodes (e.g. a gateway node pointed at a different, sometimes
// unusable, RPC than worker nodes) — every distinct variant is returned so
// the user can see and prune them in the UI, rather than silently collapsing
// to whichever one happened to be read last. The Registry flag is left unset
// on every result — the user designates it explicitly. This is prepopulation
// only; the authoritative chain list is whatever is saved into desired.toml.
func DiscoverChains(cv *domain.ChartValues) ([]domain.Chain, error) {
	byFile := make(map[string][]domain.ChartNodeInfo)
	for _, n := range cv.Nodes {
		if n.ConfigFile == "" {
			continue
		}
		byFile[n.ConfigFile] = append(byFile[n.ConfigFile], n)
	}

	seen := make(map[string]bool)
	var chains []domain.Chain
	for file, nodes := range byFile {
		if err := discoverChainsFromFile(file, nodes, &chains, seen); err != nil {
			return nil, errors.Wrapf(err, "failed to discover chains from %s", file)
		}
	}

	sort.Slice(chains, func(i, j int) bool {
		if chains[i].ChainID != chains[j].ChainID {
			return chains[i].ChainID < chains[j].ChainID
		}
		if chains[i].WSURL != chains[j].WSURL {
			return chains[i].WSURL < chains[j].WSURL
		}
		return chains[i].HTTPURL < chains[j].HTTPURL
	})
	return chains, nil
}

func discoverChainsFromFile(path string, nodes []domain.ChartNodeInfo, chains *[]domain.Chain, seen map[string]bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return err
	}
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return nil
	}
	docMap := root.Content[0]
	if docMap.Kind != yaml.MappingNode {
		return nil
	}

	clNode := infra.FindKey(docMap, "chainlink-node")
	if clNode == nil {
		return nil
	}
	instances := infra.FindKey(clNode, "instances")
	if instances == nil {
		return nil
	}

	for _, n := range nodes {
		nodeEntry := infra.FindKey(instances, n.Name)
		if nodeEntry == nil {
			continue
		}
		configList := infra.FindKey(nodeEntry, "configuration")
		if configList == nil || configList.Kind != yaml.SequenceNode {
			continue
		}
		for _, item := range configList.Content {
			if item.Kind != yaml.MappingNode || len(item.Content) < 2 {
				continue
			}
			if item.Content[0].Value == managedLayerName {
				continue
			}
			if err := mergeEVMChains(item.Content[1].Value, chains, seen); err != nil {
				return err
			}
		}
	}
	return nil
}

func mergeEVMChains(tomlContent string, chains *[]domain.Chain, seen map[string]bool) error {
	var cfg corechainlink.Config
	if err := toml.Unmarshal([]byte(tomlContent), &cfg); err != nil {
		return err
	}
	for _, evmCfg := range cfg.EVM {
		if evmCfg == nil || evmCfg.ChainID == nil {
			continue
		}
		chainID := evmCfg.ChainID.ToInt().Uint64()

		ch := domain.Chain{ChainID: chainID}
		if len(evmCfg.Nodes) > 0 {
			node := evmCfg.Nodes[0]
			if node.WSURL != nil {
				ch.WSURL = node.WSURL.String()
			}
			if node.HTTPURL != nil {
				ch.HTTPURL = node.HTTPURL.String()
			}
		}

		key := fmt.Sprintf("%d|%s|%s", ch.ChainID, ch.WSURL, ch.HTTPURL)
		if seen[key] {
			continue
		}
		seen[key] = true
		*chains = append(*chains, ch)
	}
	return nil
}
