package operations

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
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
	ServiceName string           `yaml:"servicename"`
	Handlers    []string         `yaml:"handlers"`
	DONs        []string         `yaml:"dons"`
	Auth0       *pkg.Auth0Config `yaml:"auth0,omitempty"`
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
	var donNames []string
	for _, svc := range input.Services {
		donNames = append(donNames, svc.DONs...)
	}
	shardedDONs, err := groupShardedDONs(donNames)
	if err != nil {
		return pkg.GatewayJob{}, err
	}

	dons := make([]pkg.TargetDON, 0, len(shardedDONs))
	for _, sd := range shardedDONs {
		shards := make([][]pkg.TargetDONMember, len(sd.shardDONNames))
		var f int
		for shardIdx, shardDONName := range sd.shardDONNames {
			members, shardF, err := resolveDONMembers(deps, input, shardDONName)
			if err != nil {
				return pkg.GatewayJob{}, err
			}
			if shardIdx == 0 {
				f = shardF
			} else if shardF != f {
				return pkg.GatewayJob{}, fmt.Errorf("DON %s: shard %d (%s) has F=%d, but shard 0 has F=%d; all shards must have the same size", sd.donName, shardIdx, shardDONName, shardF, f)
			}
			shards[shardIdx] = members
		}
		dons = append(dons, pkg.TargetDON{
			ID:     sd.donName,
			F:      f,
			Shards: shards,
		})
	}

	services := make([]pkg.GatewayServiceConfig, len(input.Services))
	for i, svc := range input.Services {
		services[i] = pkg.GatewayServiceConfig{
			ServiceName: svc.ServiceName,
			Handlers:    svc.Handlers,
			DONs:        baseDONNames(svc.DONs),
			Auth0:       svc.Auth0,
		}
	}

	return pkg.GatewayJob{
		ServiceCentricFormatEnabled: true,
		JobName:                     "CRE Gateway",
		DONs:                        dons,
		Services:                    services,
		RequestTimeoutSec:           requestTimeoutSec,
		AllowedPorts:                toIntSlice(input.AllowedPorts),
		AllowedSchemes:              input.AllowedSchemes,
		AllowedIPsCIDR:              input.AllowedIPsCIDR,
		AuthGatewayID:               input.AuthGatewayID,
		ExternalJobID:               input.ExternalJobID,
		AuthGatewayIDPrefix:         input.AuthGatewayIDPrefix,
	}, nil
}

// shardDONNameRe matches the DON name of shard N>0 of a sharded DON. It mirrors
// gateway config.GatewayDONIDForShard: shard 0 uses the bare DON name and shard N>0 uses "<base>_shard-N".
var shardDONNameRe = regexp.MustCompile(`^(.+)_shard-([0-9]+)$`)

// splitShardDONName returns the base DON name and shard index for a DON name.
// Names without a "_shard-N" suffix are shard 0 of a DON with that name.
func splitShardDONName(donName string) (string, int, error) {
	m := shardDONNameRe.FindStringSubmatch(donName)
	if m == nil {
		return donName, 0, nil
	}
	idx, err := strconv.Atoi(m[2])
	if err != nil {
		return "", 0, fmt.Errorf("invalid shard index in DON name %s: %w", donName, err)
	}
	if idx == 0 {
		return "", 0, fmt.Errorf("invalid DON name %s: shard 0 must use the bare DON name %s", donName, m[1])
	}
	return m[1], idx, nil
}

type shardedDON struct {
	donName string
	// shardDONNames holds the JD DON name of each shard, indexed by shard index.
	shardDONNames []string
}

// groupShardedDONs groups DON names into sharded DONs keyed by base name, sorted by name.
// Every sharded DON must have contiguous shards 0..N-1.
func groupShardedDONs(donNames []string) ([]shardedDON, error) {
	byBase := make(map[string]map[int]string)
	for _, name := range donNames {
		base, idx, err := splitShardDONName(name)
		if err != nil {
			return nil, err
		}
		if byBase[base] == nil {
			byBase[base] = make(map[int]string)
		}
		byBase[base][idx] = name
	}

	shardedDONs := make([]shardedDON, 0, len(byBase))
	for base, shards := range byBase {
		shardDONNames := make([]string, len(shards))
		for idx, name := range shards {
			if idx >= len(shards) {
				return nil, fmt.Errorf("DON %s: shards must be contiguous starting from 0, got shard %d (%s) with only %d shards referenced", base, idx, name, len(shards))
			}
			shardDONNames[idx] = name
		}
		shardedDONs = append(shardedDONs, shardedDON{donName: base, shardDONNames: shardDONNames})
	}
	slices.SortFunc(shardedDONs, func(a, b shardedDON) int { return strings.Compare(a.donName, b.donName) })
	return shardedDONs, nil
}

// baseDONNames maps DON names to their base (sharded DON) names, deduplicated and order-preserving.
// Callers must have validated the names with groupShardedDONs.
func baseDONNames(donNames []string) []string {
	out := make([]string, 0, len(donNames))
	for _, name := range donNames {
		base, _, _ := splitShardDONName(name)
		if !slices.Contains(out, base) {
			out = append(out, base)
		}
	}
	return out
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
