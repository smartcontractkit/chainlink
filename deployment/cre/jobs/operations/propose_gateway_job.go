package operations

import (
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

type ProposeGatewayJobInput struct {
	Domain                  string
	DONFilters              []offchain.TargetDONFilter
	DONs                    []DON             `yaml:"dons"`
	GatewayKeyChainSelector pkg.ChainSelector `yaml:"gatewayKeyChainSelector"`
	JobLabels               map[string]string
}

type DON struct {
	Name     string
	Handlers []string
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
	func(b operations.Bundle, deps ProposeGatewayJobDeps, input ProposeGatewayJobInput) (ProposeGatewayJobOutput, error) {
		targetDONs := make([]pkg.TargetDON, 0)
		for _, ad := range input.DONs {
			filter := offchain.TargetDONFilter{
				Key:   offchain.FilterKeyDONName,
				Value: ad.Name,
			}
			nodes, err := pkg.FetchNodeChainConfigsFromJD(deps.Env.GetContext(), deps.Env, filter)
			if err != nil {
				return ProposeGatewayJobOutput{}, err
			}

			fam, chainID, err := parseSelector(uint64(input.GatewayKeyChainSelector))
			if err != nil {
				return ProposeGatewayJobOutput{}, err
			}

			var members []pkg.TargetDONMember
			for _, n := range nodes {
				var found bool
				for _, cc := range n.ChainConfigs {
					if cc.Chain.Id == chainID && cc.Chain.Type == fam {
						members = append(members, pkg.TargetDONMember{
							Address: cc.AccountAddress,
							Name:    fmt.Sprintf("DON %s - Node %s", ad.Name, n.NodeID),
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
				Members:  members,
				Handlers: ad.Handlers,
			}

			targetDONs = append(targetDONs, td)
		}

		gj := pkg.GatewayJob{
			JobName:    "CRE Gateway",
			TargetDONs: targetDONs,
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
			Specs: make(map[string][]string),
		}
		for nodeIdx, n := range nodes {
			spec, err := gj.Resolve(nodeIdx)
			if err != nil {
				return ProposeGatewayJobOutput{}, err
			}

			_, err = deps.Env.Offchain.ProposeJob(b.GetContext(), &jobv1.ProposeJobRequest{
				NodeId: n.GetId(),
				Spec:   spec,
				Labels: labels,
			})
			if err != nil {
				return ProposeGatewayJobOutput{}, fmt.Errorf("error proposing job to node %s spec %s : %w", n.GetId(), spec, err)
			}

			output.Specs[n.GetId()] = append(output.Specs[n.GetId()], spec)
		}

		return output, nil
	},
)

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
