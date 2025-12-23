package operations

import (
	"context"
	"fmt"
	"maps"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pelletier/go-toml/v2"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/view"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/deployment/helpers/pointer"
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

	// NodeIDToConfig is a map of node IDs to custom per node configs,
	// throws an error if nodes from the map keys aren't an exact match with the DON nodes.
	NodeIDToConfig map[string]string

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
				{
					Key:   "type",
					Op:    ptypes.SelectorOp_EQ,
					Value: pointer.To(PluginNodeType),
				},
			},
		}
		for _, f := range input.DONFilters {
			filter = f.AddToFilter(filter)
		}

		for _, f := range input.DONFilters {
			filter = f.AddToFilterIfNotPresent(filter)
		}

		nodes, err := offchain.FetchNodesFromJD(b.GetContext(), deps.Env.Offchain, filter)
		if err != nil {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to fetch nodes from JD: %w", err)
		}
		if len(nodes) == 0 {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("no nodes found on JD for DON `%s` with filters %+v", input.DONName, filter)
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
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("no nodes info found for DON `%s` with filters %+v and node IDs %v", input.DONName, input.DONFilters, nodeIDs)
		}

		setPerNodeCfg := len(input.NodeIDToConfig) > 0
		if setPerNodeCfg {
			if len(input.NodeIDToConfig) != len(nodeInfos) {
				return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("number of nodes found (%d) does not match number of configs provided (%d)", len(nodeInfos), len(input.NodeIDToConfig))
			}
			for _, n := range nodeInfos {
				if _, ok := input.NodeIDToConfig[n.NodeID]; !ok {
					return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("node ID %s found in DON nodes but not in provided configs", n.NodeID)
				}
			}
		}

		shouldGenerateOracleFactory := input.Job.GenerateOracleFactory && input.Job.OracleFactory == nil
		if !shouldGenerateOracleFactory {
			specs := make(map[string][]string)

			for _, ni := range nodeInfos {
				spec, err := resolveJob(b.GetContext(), deps.Env.Logger, input.Job, setPerNodeCfg, ni.NodeID, input.NodeIDToConfig, deps.Env.Offchain)
				if err != nil {
					return ProposeStandardCapabilityJobOutput{}, err
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
					return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to propose standard capability job: %w", err)
				}

				maps.Copy(specs, report.Output.Specs)
			}

			return ProposeStandardCapabilityJobOutput{Specs: specs}, nil
		}

		// If no oracle factory is provided, we have to build it
		specs := make(map[string][]string)

		for _, ni := range nodeInfos {
			oracleFactory, err := generateOracleFactory(deps.Env, ni, input.Job)
			if err != nil {
				return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to generate oracle factory for node %s: %w", ni.NodeID, err)
			}

			input.Job.OracleFactory = oracleFactory

			spec, err := resolveJob(b.GetContext(), deps.Env.Logger, input.Job, setPerNodeCfg, ni.NodeID, input.NodeIDToConfig, deps.Env.Offchain)
			if err != nil {
				return ProposeStandardCapabilityJobOutput{}, err
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
				return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to propose standard capability job: %w", err)
			}

			maps.Copy(specs, report.Output.Specs)
		}

		return ProposeStandardCapabilityJobOutput{Specs: specs}, nil
	})

func resolveJob(ctx context.Context, lggr logger.Logger, job pkg.StandardCapabilityJob, setPerNodeCfg bool, nodeID string, nodeIDToConfig map[string]string, oc cldf_offchain.Client) (string, error) {
	if setPerNodeCfg {
		customCfg, ok := nodeIDToConfig[nodeID]
		if !ok {
			return "", fmt.Errorf("no custom config found for node ID %s", nodeID)
		}
		job.Config = customCfg
	}

	externalJobID, err := getEVMExternalJobIDByName(ctx, lggr, job.JobName, nodeID, oc)
	if err != nil {
		return "", err
	}
	if externalJobID != "" {
		job.ExternalJobID = externalJobID
	}

	spec, err := job.Resolve()
	if err != nil {
		return "", fmt.Errorf("failed to resolve standard capability job for node %s: %w", nodeID, err)
	}

	return spec, nil
}

// getEVMExternalJobIDByName returns an empty string if id is not found.
func getEVMExternalJobIDByName(ctx context.Context, lggr logger.Logger, jobName, nodeID string, oc cldf_offchain.Client) (string, error) {
	if !strings.Contains(jobName, "evm-capabilities-v2") {
		return "", nil
	}

	nodesJobs, _, err := view.ApprovedJobspecs(ctx, lggr, []string{nodeID}, oc)
	if err != nil {
		return "", fmt.Errorf("failed to fetch approved jobs for node %s: %w", nodeID, err)
	}

	nodeJobs, ok := nodesJobs[nodeID]
	if !ok || len(nodeJobs) == 0 {
		return "", nil
	}

	specFormattedJobName := `name = "` + jobName + `"`
	for _, j := range nodeJobs {
		if !strings.Contains(j.Spec, specFormattedJobName) {
			continue
		}

		ji := make(job_types.JobSpecInput)
		if err = toml.Unmarshal([]byte(j.Spec), &ji); err != nil {
			return "", fmt.Errorf("failed to unmarshal job spec toml for job %s on node %s: %w", jobName, nodeID, err)
		}

		if s, _ := ji["externalJobID"].(string); s != "" {
			return s, nil
		}
	}

	return "", nil
}

func generateOracleFactory(cldEnv cldf.Environment, nodeInfo deployment.Node, job pkg.StandardCapabilityJob) (*pkg.OracleFactory, error) {
	contractChainSelector := job.ChainSelectorEVM
	if job.OCRChainSelector != 0 {
		contractChainSelector = job.OCRChainSelector
	}

	addrRefKey := pkg.GetOCR3CapabilityAddressRefKey(uint64(contractChainSelector), job.ContractQualifier)
	contractAddrRef, err := cldEnv.DataStore.Addresses().Get(addrRefKey)
	if err != nil {
		return &pkg.OracleFactory{}, fmt.Errorf("failed to get OCR3 contract address for chain selector %d and qualifier %s: %w", contractChainSelector, job.ContractQualifier, err)
	}

	if addrRefKey.ChainSelector() != uint64(contractChainSelector) {
		return &pkg.OracleFactory{}, fmt.Errorf(
			"mismatched chain selector in address ref key for OCR3 contract %s: expected %d, got %d",
			addrRefKey.String(),
			contractChainSelector,
			addrRefKey.ChainSelector(),
		)
	}

	contractChainID, err := chainsel.GetChainIDFromSelector(addrRefKey.ChainSelector())
	if err != nil {
		return &pkg.OracleFactory{}, fmt.Errorf("failed to get chainID for chain selector %d and qualifier %s: %w", contractChainSelector, job.ContractQualifier, err)
	}

	evmOCRConfig, ok := nodeInfo.OCRConfigForChainSelector(uint64(contractChainSelector))
	if !ok {
		return &pkg.OracleFactory{}, fmt.Errorf("no evm ocr2 config for node %s", nodeInfo.NodeID)
	}

	if job.OCRSigningStrategy == "" {
		job.OCRSigningStrategy = "multi-chain"
	}

	oracleFactory := &pkg.OracleFactory{
		Enabled:            true,
		BootstrapPeers:     job.BootstrapPeers,
		OCRContractAddress: contractAddrRef.Address,
		OCRKeyBundleID:     evmOCRConfig.KeyBundleID,
		ChainID:            contractChainID,
		TransmitterID:      string(evmOCRConfig.TransmitAccount),
		OnchainSigningStrategy: pkg.OnchainSigningStrategy{
			StrategyName: job.OCRSigningStrategy,
			Config:       map[string]string{"evm": evmOCRConfig.KeyBundleID},
		},
	}

	if job.ChainSelectorAptos > 0 {
		aptosConfig, ok := nodeInfo.OCRConfigForChainSelector(uint64(job.ChainSelectorAptos))
		if !ok {
			return &pkg.OracleFactory{}, fmt.Errorf("no aptos ocr2 config for node %s", nodeInfo.NodeID)
		}

		oracleFactory.OnchainSigningStrategy.Config["aptos"] = aptosConfig.KeyBundleID
	}

	if job.ChainSelectorSolana > 0 {
		solanaConfig, ok := nodeInfo.OCRConfigForChainSelector(uint64(job.ChainSelectorSolana))
		if !ok {
			return &pkg.OracleFactory{}, fmt.Errorf("no solana ocr2 config for node %s", nodeInfo.NodeID)
		}

		oracleFactory.OnchainSigningStrategy.Config["solana"] = solanaConfig.KeyBundleID
	}

	return oracleFactory, nil
}
