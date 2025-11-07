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
	Domain                   string
	DONFilters               []offchain.TargetDONFilter
	DONs                     []DON             `yaml:"dons"`
	GatewayRequestTimeoutSec int               `yaml:"gatewayRequestTimeoutSec"`
	GatewayKeyChainSelector  pkg.ChainSelector `yaml:"gatewayKeyChainSelector"`
	JobLabels                map[string]string
}

type DON struct {
	Name     string
	F        int
	Handlers []string
}

type ProposeGatewayJobDeps struct {
	Env cldf.Environment
}

type ProposedSpec struct {
	NodeID   string
	Proposal *jobv1.Proposal
	Spec     string
}
type ProposeGatewayJobOutput struct {
	// jd node ID -> list of proposed job specs
	ProposedSpecs map[string]ProposedSpec
}

var ProposeGatewayJob = operations.NewOperation[ProposeGatewayJobInput, ProposeGatewayJobOutput, ProposeGatewayJobDeps](
	"propose-gateway-job-op",
	semver.MustParse("1.0.0"),
	"Propose Gateway Job",
	func(b operations.Bundle, deps ProposeGatewayJobDeps, input ProposeGatewayJobInput) (ProposeGatewayJobOutput, error) {
		targetDONs := make([]pkg.TargetDON, 0)
		for _, ad := range input.DONs {
			filter := offchain.TargetDONFilter{
				Key:   offchain.FilterKeyDONName,
				Value: ad.Name,
			}
			nodes, err := pkg.FetchNodeChainConfigsFromJD(deps.Env.GetContext(), deps.Env, filter)
			if err != nil {
				return ProposeGatewayJobOutput{}, fmt.Errorf("failed to fetch chain configs for DON %s filter %+v: %w", ad.Name, filter, err)
			}

			fam, chainID, err := parseSelector(uint64(input.GatewayKeyChainSelector))
			if err != nil {
				return ProposeGatewayJobOutput{}, fmt.Errorf("failed to parse chain selector %d: %w", input.GatewayKeyChainSelector, err)
			}

			req := pkg.FetchNodesRequest{
				Domain:  input.Domain,
				Filters: []offchain.TargetDONFilter{filter},
			}
			ns, err := pkg.FetchNodesFromJD(deps.Env.GetContext(), deps.Env, req)
			if err != nil {
				return ProposeGatewayJobOutput{}, fmt.Errorf("failed to fetch nodes from JD req %+v: %w", req, err)
			}

			// make map of node id to node
			m := make(map[string]*nodev1.Node, len(ns))
			for _, n := range ns {
				m[n.Id] = n
			}

			var members []pkg.TargetDONMember
			for _, n := range nodes {
				var found bool
				for _, cc := range n.ChainConfigs {
					if cc.Chain.Id == chainID && cc.Chain.Type == fam {
						nodeName := n.NodeID
						if matched, ok := m[n.NodeID]; ok {
							nodeName = matched.Name
						}
						members = append(members, pkg.TargetDONMember{
							Address: cc.AccountAddress,
							Name:    fmt.Sprintf("%s (DON %s)", nodeName, ad.Name),
						})
						found = true

						break
					}
				}

				if !found {
					return ProposeGatewayJobOutput{}, fmt.Errorf("could not find key belonging to chain id %s", chainID)
				}
			}

			td := pkg.TargetDON{
				ID:       ad.Name,
				F:        ad.F,
				Members:  members,
				Handlers: ad.Handlers,
			}

			targetDONs = append(targetDONs, td)
		}

		requestTimeoutSec := input.GatewayRequestTimeoutSec
		if requestTimeoutSec == 0 {
			requestTimeoutSec = defaultGatewayRequestTimeoutSec
		}

		gj := pkg.GatewayJob{
			JobName:           "CRE Gateway",
			TargetDONs:        targetDONs,
			RequestTimeoutSec: requestTimeoutSec,
		}

		err := gj.Validate()
		if err != nil {
			return ProposeGatewayJobOutput{}, err
		}

		nodes, err := pkg.FetchNodesFromJD(b.GetContext(), deps.Env, pkg.FetchNodesRequest{
			Domain:  input.Domain,
			Filters: input.DONFilters,
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
			ProposedSpecs: make(map[string]ProposedSpec),
		}
		for nodeIdx, n := range nodes {
			spec, err := gj.Resolve(nodeIdx)
			if err != nil {
				return ProposeGatewayJobOutput{}, fmt.Errorf("failed to resolve gateway job spec for node %s: %w", n.GetId(), err)
			}

			r, err := deps.Env.Offchain.ProposeJob(b.GetContext(), &jobv1.ProposeJobRequest{
				NodeId: n.GetId(),
				Spec:   spec,
				Labels: labels,
			})
			if err != nil {
				return ProposeGatewayJobOutput{}, fmt.Errorf("error proposing job to node %s spec %s : %w", n.GetId(), spec, err)
			}

			output.ProposedSpecs[n.GetId()] = ProposedSpec{
				NodeID:   n.GetId(),
				Proposal: r.Proposal,
				Spec:     spec,
			}
		}
		if len(output.ProposedSpecs) == 0 {
			return ProposeGatewayJobOutput{}, errors.New("no gateway jobs were proposed")
		}

		return output, nil
	},
)

func parseSelector(sel uint64) (nodev1.ChainType, string, error) {
	fam, err := chainsel.GetSelectorFamily(sel)
	if err != nil {
		return nodev1.ChainType_CHAIN_TYPE_UNSPECIFIED, "", fmt.Errorf("failed to get selector family: %w", err)
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
		return nodev1.ChainType_CHAIN_TYPE_UNSPECIFIED, "", fmt.Errorf("failed to get chain ID from selector %d: %w", sel, err)
	}

	return ct, chainID, nil
}
