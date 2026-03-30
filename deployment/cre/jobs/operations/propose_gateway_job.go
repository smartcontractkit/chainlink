package operations

import (
	"errors"
	"fmt"

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
	GatewayKeyChainSelector     pkg.ChainSelector `yaml:"gatewayKeyChainSelector"`
	JobLabels                   map[string]string
}

type DON struct {
	Name     string   `yaml:"name"`
	F        pkg.Int  `yaml:"f"`
	Handlers []string `yaml:"handlers"`
}

type GatewayService struct {
	ServiceName string   `yaml:"servicename"`
	Handlers    []string `yaml:"handlers"`
	DONs        []string `yaml:"dons"`
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
// When ServiceCentricFormatEnabled is true, it derives the set of unique DON names from
// input.Services; otherwise it uses the don-centric input.DONs list.
func proposeGatewayJob(b operations.Bundle, deps ProposeGatewayJobDeps, input ProposeGatewayJobInput) (ProposeGatewayJobOutput, error) {
	requestTimeoutSec := int(input.GatewayRequestTimeoutSec)
	if requestTimeoutSec == 0 {
		requestTimeoutSec = defaultGatewayRequestTimeoutSec
	}

	var gj pkg.GatewayJob
	if input.ServiceCentricFormatEnabled {
		built, err := buildServiceCentricJob(deps, input, requestTimeoutSec)
		if err != nil {
			return ProposeGatewayJobOutput{}, err
		}
		gj = built
	} else {
		built, err := buildLegacyFormatJob(deps, input, requestTimeoutSec)
		if err != nil {
			return ProposeGatewayJobOutput{}, err
		}
		gj = built
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
	donNameSet := make(map[string]struct{})
	for _, svc := range input.Services {
		for _, donName := range svc.DONs {
			donNameSet[donName] = struct{}{}
		}
	}

	dons := make([]pkg.TargetDON, 0, len(donNameSet))
	for donName := range donNameSet {
		members, f, err := resolveDONMembers(deps, input, donName)
		if err != nil {
			return pkg.GatewayJob{}, err
		}
		dons = append(dons, pkg.TargetDON{
			ID:      donName,
			F:       f,
			Members: members,
		})
	}

	services := make([]pkg.GatewayServiceConfig, len(input.Services))
	for i, svc := range input.Services {
		services[i] = pkg.GatewayServiceConfig{
			ServiceName: svc.ServiceName,
			Handlers:    svc.Handlers,
			DONs:        svc.DONs,
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
	}, nil
}

func buildLegacyFormatJob(deps ProposeGatewayJobDeps, input ProposeGatewayJobInput, requestTimeoutSec int) (pkg.GatewayJob, error) {
	targetDONs := make([]pkg.TargetDON, 0, len(input.DONs))
	for _, ad := range input.DONs {
		members, _, err := resolveDONMembers(deps, input, ad.Name)
		if err != nil {
			return pkg.GatewayJob{}, err
		}
		targetDONs = append(targetDONs, pkg.TargetDON{
			ID:       ad.Name,
			F:        int(ad.F),
			Members:  members,
			Handlers: ad.Handlers,
		})
	}

	return pkg.GatewayJob{
		JobName:           "CRE Gateway",
		TargetDONs:        targetDONs,
		RequestTimeoutSec: requestTimeoutSec,
		AllowedPorts:      toIntSlice(input.AllowedPorts),
		AllowedSchemes:    input.AllowedSchemes,
		AllowedIPsCIDR:    input.AllowedIPsCIDR,
		AuthGatewayID:     input.AuthGatewayID,
	}, nil
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
