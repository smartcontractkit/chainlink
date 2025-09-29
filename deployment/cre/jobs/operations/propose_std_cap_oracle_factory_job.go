package operations

import (
	"fmt"

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

type ProposeStandardCapabilityWithOracleFactoryJobDeps struct {
	Env cldf.Environment
}

type ProposeStandardCapabilityWithOracleFactoryJobInput struct {
	Domain      string
	DONName     string
	Job         pkg.StandardCapabilityJobWithOracleFactory
	DONFilters  []offchain.TargetDONFilter
	ExtraLabels map[string]string
}

type ProposeStandardCapabilityWithOracleFactoryJobOutput struct {
	Specs map[string][]string
}

var ProposeStandardCapabilityWithOracleFactoryJob = operations.NewSequence[
	ProposeStandardCapabilityWithOracleFactoryJobInput,
	ProposeStandardCapabilityWithOracleFactoryJobOutput,
	ProposeStandardCapabilityWithOracleFactoryJobDeps,
](
	"propose-standard-capability-oracle-factory-job-seq",
	semver.MustParse("1.0.0"),
	"Propose Standard Capability w/ Oracle Factory Job",
	func(b operations.Bundle, deps ProposeStandardCapabilityWithOracleFactoryJobDeps, input ProposeStandardCapabilityWithOracleFactoryJobInput) (ProposeStandardCapabilityWithOracleFactoryJobOutput, error) {
		if err := input.Job.Validate(); err != nil {
			return ProposeStandardCapabilityWithOracleFactoryJobOutput{}, fmt.Errorf("invalid job: %w", err)
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
		nodes, err := offchain.FetchNodesFromJD(b.GetContext(), deps.Env.Offchain, filter)
		if err != nil {
			return ProposeStandardCapabilityWithOracleFactoryJobOutput{}, fmt.Errorf("failed to fetch nodes from JD: %w", err)
		}

		nodeIDs := make([]string, len(nodes))
		for i, n := range nodes {
			nodeIDs[i] = n.Id
		}

		nodeInfos, err := deployment.NodeInfo(nodeIDs, deps.Env.Offchain)
		if err != nil {
			return ProposeStandardCapabilityWithOracleFactoryJobOutput{}, fmt.Errorf("failed to fetch node infos: %w", err)
		}

		addrRefKey := pkg.GetOCR3CapabilityAddressRefKey(uint64(input.Job.ChainSelectorEVM), input.Job.ContractQualifier)
		contractAddrRef, err := deps.Env.DataStore.Addresses().Get(addrRefKey)
		if err != nil {
			return ProposeStandardCapabilityWithOracleFactoryJobOutput{}, fmt.Errorf("failed to get OCR3 contract address for chain selector %d and qualifier %s: %w", input.Job.ChainSelectorEVM, input.Job.ContractQualifier, err)
		}

		chainID, err := chainsel.GetChainIDFromSelector(uint64(input.Job.ChainSelectorEVM))
		if err != nil {
			return ProposeStandardCapabilityWithOracleFactoryJobOutput{}, fmt.Errorf("failed to get chain ID from selector: %w", err)
		}

		specs := make(map[string][]string)

		for _, ni := range nodeInfos {
			evmConfig, ok := ni.OCRConfigForChainSelector(uint64(input.Job.ChainSelectorEVM))
			if !ok {
				return ProposeStandardCapabilityWithOracleFactoryJobOutput{}, fmt.Errorf("no evm ocr2 config for node %s", ni.NodeID)
			}
			aptosConfig, ok := ni.OCRConfigForChainSelector(uint64(input.Job.ChainSelectorAptos))
			if !ok {
				return ProposeStandardCapabilityWithOracleFactoryJobOutput{}, fmt.Errorf("no aptos ocr2 config for node %s", ni.NodeID)
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
				return ProposeStandardCapabilityWithOracleFactoryJobOutput{}, fmt.Errorf("failed to resolve consensus job for node %s: %w", ni.NodeID, err)
			}

			jobLabels := map[string]string{
				offchain.CapabilityLabel: input.Job.JobName,
			}
			for k, v := range input.ExtraLabels {
				jobLabels[k] = v
			}

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
				return ProposeStandardCapabilityWithOracleFactoryJobOutput{}, fmt.Errorf("failed to propose consensus job: %w", err)
			}

			for k, v := range report.Output.Specs {
				specs[k] = v
			}
		}

		return ProposeStandardCapabilityWithOracleFactoryJobOutput{Specs: specs}, nil
	})
