package operations

import (
	"fmt"
	"maps"

	"github.com/Masterminds/semver/v3"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

type ProposeStandardCapabilityJobDeps struct {
	Env cldf.Environment
}

type ProposeStandardCapabilityJobInput struct {
	Domain  string
	DONName string

	// Job is the standard capability job to propose.
	// If GenerateOracleFactory is true, the OracleFactory field will be ignored and generated.
	// If false, the OracleFactory field will be used as-is.
	Job pkg.StandardCapabilityJob

	DONFilters  []offchain.TargetDONFilter
	ExtraLabels map[string]string
}

type ProposeStandardCapabilityJobOutput struct {
	Specs map[string][]string
}

var ProposeStandardCapabilityJob = operations.NewSequence[
	ProposeStandardCapabilityJobInput,
	ProposeStandardCapabilityJobOutput,
	ProposeStandardCapabilityJobDeps,
](
	"propose-standard-capability-job-seq",
	semver.MustParse("1.0.0"),
	"Propose Standard Capability Job",
	func(b operations.Bundle, deps ProposeStandardCapabilityJobDeps, input ProposeStandardCapabilityJobInput) (ProposeStandardCapabilityJobOutput, error) {
		if err := input.Job.Validate(); err != nil {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("invalid job: %w", err)
		}

		filter := &node.ListNodesRequest_Filter{
			Selectors: []*ptypes.Selector{
				{
					Key: "don-" + input.DONName,
					Op:  ptypes.SelectorOp_EXIST,
				},
				{
					Key:   "environment",
					Op:    ptypes.SelectorOp_EQ,
					Value: &deps.Env.Name,
				},
				{
					Key:   "product",
					Op:    ptypes.SelectorOp_EQ,
					Value: &input.Domain,
				},
			},
		}

		for _, f := range input.DONFilters {
			filter = f.AddToFilterIfNotPresent(filter)
		}

		nodes, err := offchain.FetchNodesFromJD(b.GetContext(), deps.Env.Offchain, filter)
		if err != nil {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to fetch nodes from JD: %w", err)
		}

		nodeIDs := make([]string, len(nodes))
		for i, n := range nodes {
			nodeIDs[i] = n.Id
		}

		nodeInfos, err := deployment.NodeInfo(nodeIDs, deps.Env.Offchain)
		if err != nil {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to fetch node infos: %w", err)
		}

		if !input.Job.GenerateOracleFactory {
			specs := make(map[string][]string)

			for _, ni := range nodeInfos {
				spec, err := input.Job.Resolve()
				if err != nil {
					return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to resolve consensus job for node %s: %w", ni.NodeID, err)
				}

				jobLabels := map[string]string{
					offchain.CapabilityLabel: input.Job.JobName,
				}
				maps.Copy(jobLabels, input.ExtraLabels)

				// 1 spec per node, each spec is unique to the node due to the oracle factory config
				report, err := operations.ExecuteOperation(b, ProposeJobSpec, ProposeJobSpecDeps(deps), ProposeJobSpecInput{
					Domain:    input.Domain,
					DONName:   input.DONName,
					Spec:      spec,
					JobLabels: jobLabels,
					DONFilters: []offchain.TargetDONFilter{
						{Key: "p2p_id", Value: ni.PeerID.String()},
					},
				})
				if err != nil {
					return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to propose consensus job: %w", err)
				}

				maps.Copy(specs, report.Output.Specs)
			}

			return ProposeStandardCapabilityJobOutput{Specs: specs}, nil
		}

		// If no oracle factory is provided, we have to build it

		addrRefKey := pkg.GetOCR3CapabilityAddressRefKey(uint64(input.Job.ChainSelectorEVM), input.Job.ContractQualifier)
		contractAddrRef, err := deps.Env.DataStore.Addresses().Get(addrRefKey)
		if err != nil {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to get OCR3 contract address for chain selector %d and qualifier %s: %w", input.Job.ChainSelectorEVM, input.Job.ContractQualifier, err)
		}

		chainID, err := chainsel.GetChainIDFromSelector(uint64(input.Job.ChainSelectorEVM))
		if err != nil {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to get chain ID from selector: %w", err)
		}

		specs := make(map[string][]string)

		for _, ni := range nodeInfos {
			evmConfig, ok := ni.OCRConfigForChainSelector(uint64(input.Job.ChainSelectorEVM))
			if !ok {
				return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("no evm ocr2 config for node %s", ni.NodeID)
			}
			aptosConfig, ok := ni.OCRConfigForChainSelector(uint64(input.Job.ChainSelectorAptos))
			if !ok {
				return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("no aptos ocr2 config for node %s", ni.NodeID)
			}

			oracleFactory := &pkg.OracleFactory{
				Enabled:            true,
				BootstrapPeers:     input.Job.BootstrapPeers,
				OCRContractAddress: contractAddrRef.Address,
				OCRKeyBundleID:     evmConfig.KeyBundleID,
				ChainID:            chainID,
				TransmitterID:      string(evmConfig.TransmitAccount),
				OnchainSigningStrategy: pkg.OnchainSigningStrategy{
					StrategyName: "multi-chain",
					Config: map[string]string{"evm": evmConfig.KeyBundleID,
						"aptos": aptosConfig.KeyBundleID},
				},
			}

			input.Job.OracleFactory = oracleFactory

			spec, err := input.Job.Resolve()
			if err != nil {
				return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to resolve consensus job for node %s: %w", ni.NodeID, err)
			}

			jobLabels := map[string]string{
				offchain.CapabilityLabel: input.Job.JobName,
			}
			maps.Copy(jobLabels, input.ExtraLabels)

			// 1 spec per node, each spec is unique to the node due to the oracle factory config
			report, err := operations.ExecuteOperation(b, ProposeJobSpec, ProposeJobSpecDeps(deps), ProposeJobSpecInput{
				Domain:    input.Domain,
				DONName:   input.DONName,
				Spec:      spec,
				JobLabels: jobLabels,
				DONFilters: []offchain.TargetDONFilter{
					{Key: "p2p_id", Value: ni.PeerID.String()},
				},
			})
			if err != nil {
				return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to propose consensus job: %w", err)
			}

			maps.Copy(specs, report.Output.Specs)
		}

		return ProposeStandardCapabilityJobOutput{Specs: specs}, nil
	})
