package operations

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/Masterminds/semver/v3"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

const defaultGatewayRequestTimeoutSec = 12

type ProposeGatewayJobInput struct {
	Domain                      string
	DONFilters                  []offchain.TargetDONFilter
	ServiceCentricFormatEnabled bool              `yaml:"serviceCentricFormatEnabled"`
	DONs                        []DON             `yaml:"dons"`
	Services                    []GatewayService  `yaml:"services"`
	GatewayRequestTimeoutSec    pkg.Int           `yaml:"gatewayRequestTimeoutSec"`
	PongTimeoutSec              pkg.Int           `yaml:"pongTimeoutSec"`
	AllowedPorts                []pkg.Int         `yaml:"allowedPorts"`
	AllowedSchemes              []string          `yaml:"allowedSchemes"`
	AllowedIPsCIDR              []string          `yaml:"allowedIPsCIDR"`
	AuthGatewayID               string            `yaml:"authGatewayID"`
	AuthGatewayIDPrefix         string            `yaml:"authGatewayIDPrefix"`
	GatewayKeyChainSelector     pkg.ChainSelector `yaml:"gatewayKeyChainSelector"`
	JobLabels                   map[string]string
	// ExternalJobID, when set, overrides the deterministic externalJobID
	// GatewayJob.Resolve would otherwise derive from JobName.
	ExternalJobID string `yaml:"externalJobID"`
}

type DON struct {
	Name     string   `yaml:"name"`
	F        pkg.Int  `yaml:"f"`
	Handlers []string `yaml:"handlers"`
}

type GatewayService struct {
	ServiceName string   `yaml:"servicename"`
	Handlers    []string `yaml:"handlers"`
	// Deprecated: use ShardedDONs. Each entry is treated as a single-shard DON.
	DONs []string `yaml:"dons"`
	// ShardedDONs lists the DONs served by this service. Each entry holds the JD DON names
	// of the DON's shards, ordered by shard index, and must follow the format
	// ["<don-name>", "<don-name>_shard-1", "<don-name>_shard-2", ...]. The first one is used
	// as the DonName in the gateway job spec, and the gateway derives the DON ID of shard N
	// from it with config.GatewayDONIDForShard; the shard's nodes can only connect if that ID
	// matches the DonID in their GatewayConnector config. A non-sharded DON has a single shard.
	ShardedDONs [][]string       `yaml:"shardedDONs"`
	Auth0       *pkg.Auth0Config `yaml:"auth0,omitempty"`
}

// allShardedDONs returns the service's ShardedDONs followed by its deprecated DONs,
// each converted to a single-shard DON.
func (s GatewayService) allShardedDONs() [][]string {
	out := slices.Clone(s.ShardedDONs)
	for _, donName := range s.DONs {
		out = append(out, []string{donName})
	}
	return out
}

type ProposeGatewayJobDeps struct {
	Env cldf.Environment
}

type ProposeGatewayJobOutput struct {
	Specs map[string][]string
}

var ProposeGatewayJob = operations.NewOperation[ProposeGatewayJobInput, ProposeGatewayJobOutput, ProposeGatewayJobDeps](
	"propose-gateway-job-op",
	semver.MustParse("1.0.0"),
	"Propose Gateway Job",
	proposeGatewayJob,
)

// proposeGatewayJob builds a gateway job spec and then proposes it to the nodes of a DON.
// It requires ServiceCentricFormatEnabled to be true and derives the set of unique DON
// names from input.Services; the legacy don-centric input.DONs format is no longer supported.
func proposeGatewayJob(b operations.Bundle, deps ProposeGatewayJobDeps, input ProposeGatewayJobInput) (ProposeGatewayJobOutput, error) {
	requestTimeoutSec := int(input.GatewayRequestTimeoutSec)
	if requestTimeoutSec == 0 {
		requestTimeoutSec = defaultGatewayRequestTimeoutSec
	}

	if !input.ServiceCentricFormatEnabled {
		return ProposeGatewayJobOutput{}, errors.New("ServiceCentricFormatEnabled has to be true - legacy format is no longer supported")
	}
	gj, err := buildServiceCentricJob(deps, input, requestTimeoutSec)
	if err != nil {
		return ProposeGatewayJobOutput{}, err
	}

	if err := gj.Validate(); err != nil {
		return ProposeGatewayJobOutput{}, err
	}

	filters := &nodev1.ListNodesRequest_Filter{}
	for _, f := range input.DONFilters {
		filters = offchain.TargetDONFilter{
			Key:   f.Key,
			Value: f.Value,
		}.AddToFilter(filters)
	}

	nodes, err := pkg.FetchNodesFromJD(b.GetContext(), deps.Env, pkg.FetchNodesRequest{
		Domain:  input.Domain,
		Filters: filters,
	})
	if err != nil {
		return ProposeGatewayJobOutput{}, fmt.Errorf("failed to fetch nodes from JD: %w", err)
	}
	if len(nodes) == 0 {
		return ProposeGatewayJobOutput{}, fmt.Errorf("no nodes found for domain %s with filters %+v", input.Domain, input.DONFilters)
	}

	labels := make([]*ptypes.Label, 0, len(input.JobLabels))
	for k, v := range input.JobLabels {
		newVal := v
		labels = append(labels, &ptypes.Label{
			Key:   k,
			Value: &newVal,
		})
	}

	output := ProposeGatewayJobOutput{
		Specs: make(map[string][]string),
	}
	for nodeIdx, n := range nodes {
		spec, specErr := gj.Resolve(nodeIdx)
		if specErr != nil {
			return ProposeGatewayJobOutput{}, specErr
		}

		_, propErr := deps.Env.Offchain.ProposeJob(b.GetContext(), &jobv1.ProposeJobRequest{
			NodeId: n.GetId(),
			Spec:   spec,
			Labels: labels,
		})
		if propErr != nil {
			return ProposeGatewayJobOutput{}, fmt.Errorf("error proposing job to node %s spec %s : %w", n.GetId(), spec, propErr)
		}

		output.Specs[n.GetId()] = append(output.Specs[n.GetId()], spec)
	}
	if len(output.Specs) == 0 {
		return ProposeGatewayJobOutput{}, errors.New("no gateway jobs were proposed")
	}

	return output, nil
}

func buildServiceCentricJob(deps ProposeGatewayJobDeps, input ProposeGatewayJobInput, requestTimeoutSec int) (pkg.GatewayJob, error) {
	shardedDONs, err := collectShardedDONs(input.Services)
	if err != nil {
		return pkg.GatewayJob{}, err
	}

	dons := make([]pkg.TargetDON, 0, len(shardedDONs))
	for _, shardDONNames := range shardedDONs {
		donName := shardDONNames[0]
		shards := make([][]pkg.TargetDONMember, len(shardDONNames))
		var f int
		for shardIdx, shardDONName := range shardDONNames {
			members, shardF, err := resolveDONMembers(deps, input, shardDONName)
			if err != nil {
				return pkg.GatewayJob{}, err
			}
			if shardIdx == 0 {
				f = shardF
			} else if shardF != f {
				return pkg.GatewayJob{}, fmt.Errorf("DON %s: shard %d (%s) has F=%d, but shard 0 has F=%d; all shards must have the same size", donName, shardIdx, shardDONName, shardF, f)
			}
			shards[shardIdx] = members
		}
		dons = append(dons, pkg.TargetDON{
			ID:     donName,
			F:      f,
			Shards: shards,
		})
	}

	services := make([]pkg.GatewayServiceConfig, len(input.Services))
	for i, svc := range input.Services {
		var donNames []string
		for _, shardDONNames := range svc.allShardedDONs() {
			donNames = append(donNames, shardDONNames[0])
		}
		services[i] = pkg.GatewayServiceConfig{
			ServiceName: svc.ServiceName,
			Handlers:    svc.Handlers,
			DONs:        donNames,
			Auth0:       svc.Auth0,
		}
	}

	return pkg.GatewayJob{
		ServiceCentricFormatEnabled: true,
		JobName:                     "CRE Gateway",
		DONs:                        dons,
		Services:                    services,
		RequestTimeoutSec:           requestTimeoutSec,
		PongTimeoutSec:              int(input.PongTimeoutSec),
		AllowedPorts:                toIntSlice(input.AllowedPorts),
		AllowedSchemes:              input.AllowedSchemes,
		AllowedIPsCIDR:              input.AllowedIPsCIDR,
		AuthGatewayID:               input.AuthGatewayID,
		ExternalJobID:               input.ExternalJobID,
		AuthGatewayIDPrefix:         input.AuthGatewayIDPrefix,
	}, nil
}

// collectShardedDONs returns the unique DONs referenced by services as lists of shard
// JD DON names, sorted by DON name (the first shard's name). A DON referenced by several
// services must have the same shards everywhere, and a JD DON can be a shard of only one DON.
func collectShardedDONs(services []GatewayService) ([][]string, error) {
	byName := make(map[string][]string)
	donByShard := make(map[string]string)
	for _, svc := range services {
		for _, shardDONNames := range svc.allShardedDONs() {
			if len(shardDONNames) == 0 {
				return nil, fmt.Errorf("service %s: sharded DON must have at least one shard", svc.ServiceName)
			}
			donName := shardDONNames[0]
			if existing, ok := byName[donName]; ok {
				if !slices.Equal(existing, shardDONNames) {
					return nil, fmt.Errorf("DON %s is referenced with different shards: %v and %v", donName, existing, shardDONNames)
				}
				continue
			}
			for _, shardDONName := range shardDONNames {
				if shardDONName == "" {
					return nil, fmt.Errorf("service %s: shard DON name must not be empty", svc.ServiceName)
				}
				if otherDON, ok := donByShard[shardDONName]; ok {
					return nil, fmt.Errorf("JD DON %s is used as a shard of both DON %s and DON %s", shardDONName, otherDON, donName)
				}
				donByShard[shardDONName] = donName
			}
			byName[donName] = shardDONNames
		}
	}

	out := slices.Collect(maps.Values(byName))
	slices.SortFunc(out, func(a, b []string) int { return strings.Compare(a[0], b[0]) })
	return out, nil
}

func resolveDONMembers(deps ProposeGatewayJobDeps, input ProposeGatewayJobInput, donName string) ([]pkg.TargetDONMember, int, error) {
	filters := &nodev1.ListNodesRequest_Filter{}
	for _, f := range input.DONFilters {
		if f.Key == offchain.FilterKeyDONName {
			continue
		}
		filters = offchain.TargetDONFilter{
			Key:   f.Key,
			Value: f.Value,
		}.AddToFilter(filters)
	}
	filtersWithTargetDONName := offchain.TargetDONFilter{
		Key:   offchain.FilterKeyDONName,
		Value: donName,
	}.AddToFilter(filters)

	ns, err := pkg.FetchNodesFromJD(deps.Env.GetContext(), deps.Env, pkg.FetchNodesRequest{
		Domain:  input.Domain,
		Filters: filtersWithTargetDONName,
	})
	if err != nil {
		return nil, 0, err
	}
	if len(ns) == 0 {
		return nil, 0, fmt.Errorf("no nodes with filters %s", input.DONFilters)
	}

	nodeChainConfigs, err := pkg.FetchNodeChainConfigsFromJD(deps.Env.GetContext(), deps.Env, pkg.FetchNodesRequest{
		Domain:  input.Domain,
		Filters: filtersWithTargetDONName,
	})
	if err != nil {
		return nil, 0, err
	}
	if len(nodeChainConfigs) == 0 {
		return nil, 0, fmt.Errorf("no chain configs with filters %s", input.DONFilters)
	}

	fam, chainID, err := parseSelector(uint64(input.GatewayKeyChainSelector))
	if err != nil {
		return nil, 0, err
	}

	m := make(map[string]*nodev1.Node, len(ns))
	for _, n := range ns {
		m[n.Id] = n
	}

	var members []pkg.TargetDONMember
	for _, n := range nodeChainConfigs {
		var found bool
		for _, cc := range n.ChainConfigs {
			if cc.Chain.Id == chainID && cc.Chain.Type == fam {
				nodeName := n.NodeID
				if matched, ok := m[n.NodeID]; ok {
					nodeName = matched.Name
				}
				members = append(members, pkg.TargetDONMember{
					Address: cc.AccountAddress,
					Name:    fmt.Sprintf("%s (DON %s)", nodeName, donName),
				})
				found = true
				break
			}
		}
		if !found {
			return nil, 0, fmt.Errorf("could not find key belonging to chain id %s on node %s", chainID, n.NodeID)
		}
	}

	f := (len(members) - 1) / 3
	return members, f, nil
}

func toIntSlice(vs []pkg.Int) []int {
	out := make([]int, len(vs))
	for i, v := range vs {
		out[i] = int(v)
	}
	return out
}

func parseSelector(sel uint64) (nodev1.ChainType, string, error) {
	fam, err := chainsel.GetSelectorFamily(sel)
	if err != nil {
		return nodev1.ChainType_CHAIN_TYPE_UNSPECIFIED, "", err
	}

	var ct nodev1.ChainType
	switch fam {
	case chainsel.FamilyEVM:
		ct = nodev1.ChainType_CHAIN_TYPE_EVM
	case chainsel.FamilySolana:
		ct = nodev1.ChainType_CHAIN_TYPE_SOLANA
	case chainsel.FamilyStarknet:
		ct = nodev1.ChainType_CHAIN_TYPE_STARKNET
	case chainsel.FamilyAptos:
		ct = nodev1.ChainType_CHAIN_TYPE_APTOS
	default:
		return nodev1.ChainType_CHAIN_TYPE_UNSPECIFIED, "", fmt.Errorf("unsupported chain type: %s", fam)
	}

	chainID, err := chainsel.GetChainIDFromSelector(sel)
	if err != nil {
		return nodev1.ChainType_CHAIN_TYPE_UNSPECIFIED, "", err
	}

	return ct, chainID, nil
}
