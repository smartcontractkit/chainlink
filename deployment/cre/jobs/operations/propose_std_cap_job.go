package operations

import (
	"fmt"
	"maps"

	"github.com/Masterminds/semver/v3"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

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

	JobNodesSpecifier offchain.NodesSpecifier
	ExtraLabels       map[string]string
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
		err := input.JobNodesSpecifier.Validate()
		if err != nil {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("invalid job nodes specifier: %w", err)
		}

		filter := input.JobNodesSpecifier.Filter(input.DONName, deps.Env.Name, input.Domain)
		nodes, err := offchain.FetchNodesFromJD(b.GetContext(), deps.Env.Offchain, filter)
		if err != nil {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to fetch nodes from JD: %w", err)
		}
		if len(nodes) == 0 {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("no nodes found on JD for DON `%s` for specifier %+v, filter %+v", input.DONName, input.JobNodesSpecifier, filter)
		}

		nodeIDs := make([]string, len(nodes))
		for i, n := range nodes {
			nodeIDs[i] = n.Id
		}
		nodeInfos, err := deployment.NodeInfo(nodeIDs, deps.Env.Offchain)
		if err != nil {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to fetch node infos: %w", err)
		}
		if len(nodeInfos) == 0 {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("no nodes info found for DON `%s` with filters %+v and node IDs %v", input.DONName, input.JobNodesSpecifier, nodeIDs)
		}
		// remove bootstrap nodes from the list
		var nonBootstrapNodeInfos deployment.Nodes
		for _, ni := range nodeInfos {
			if ni.IsBootstrap {
				deps.Env.Logger.Warnw("Skipping bootstrap node for standard capability job proposal", "node_id", ni.NodeID)
				continue
			}
			nonBootstrapNodeInfos = append(nonBootstrapNodeInfos, ni)
		}

		generateOracleFactory := input.Job.GenerateOracleFactory && input.Job.OracleFactory == nil
		if !generateOracleFactory {
			specs := make(map[string][]string)

			for _, ni := range nonBootstrapNodeInfos {
				spec, err := input.Job.Resolve()
				if err != nil {
					return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to resolve standard capability job for node %s: %w", ni.NodeID, err)
				}

				jobLabels := map[string]string{
					offchain.CapabilityLabel: input.Job.JobName,
				}
				maps.Copy(jobLabels, input.ExtraLabels)

				// 1 spec per node, each spec is unique to the node due to the oracle factory config
				specifier := offchain.NodesSpecifier{
					NodeIDs: []string{ni.NodeID},
				}
				report, err := operations.ExecuteOperation(b, ProposeJobSpec, ProposeJobSpecDeps(deps), ProposeJobSpecInput{
					Domain:            input.Domain,
					DONName:           input.DONName,
					Spec:              spec,
					JobLabels:         jobLabels,
					JobNodesSpecifier: specifier,
				})
				if err != nil {
					return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to propose standard capability job to node %s, %+v: %w", ni.NodeID, ni, err)
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

		for _, ni := range nonBootstrapNodeInfos {
			evmConfig, ok := ni.OCRConfigForChainSelector(uint64(input.Job.ChainSelectorEVM))
			if !ok {
				return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("no evm ocr2 config for node %s", ni.NodeID)
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
					Config:       map[string]string{"evm": evmConfig.KeyBundleID},
				},
			}

			if input.Job.ChainSelectorAptos > 0 {
				aptosConfig, ok := ni.OCRConfigForChainSelector(uint64(input.Job.ChainSelectorAptos))
				if !ok {
					return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("no aptos ocr2 config for node %s", ni.NodeID)
				}

				oracleFactory.OnchainSigningStrategy.Config["aptos"] = aptosConfig.KeyBundleID
			}

			input.Job.OracleFactory = oracleFactory

			spec, err := input.Job.Resolve()
			if err != nil {
				return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to resolve standard capability job for node %s: %w", ni.NodeID, err)
			}

			jobLabels := map[string]string{
				offchain.CapabilityLabel: input.Job.JobName,
			}
			maps.Copy(jobLabels, input.ExtraLabels)

			// 1 spec per node, each spec is unique to the node due to the oracle factory config
			specifier := offchain.NodesSpecifier{
				NodeIDs: []string{ni.NodeID},
			}
			report, err := operations.ExecuteOperation(b, ProposeJobSpec, ProposeJobSpecDeps(deps), ProposeJobSpecInput{
				Domain:            input.Domain,
				DONName:           input.DONName,
				Spec:              spec,
				JobLabels:         jobLabels,
				JobNodesSpecifier: specifier,
			})
			if err != nil {
				return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to propose standard capability job: %w", err)
			}

			maps.Copy(specs, report.Output.Specs)
		}

		return ProposeStandardCapabilityJobOutput{Specs: specs}, nil
	})
