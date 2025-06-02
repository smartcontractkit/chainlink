package v1_5

import (
	"errors"
	"fmt"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	integrationtesthelpers "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers/integration"
)

var _ cldf.ChangeSet[JobSpecsForLanesConfig] = JobSpecsForLanesChangeset

type JobSpecsForLanesConfig struct {
	Configs []JobSpecInput
}

func (c JobSpecsForLanesConfig) Validate() error {
	for _, cfg := range c.Configs {
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid JobSpecInput: %w", err)
		}
	}
	return nil
}

type JobSpecInput struct {
	SourceChainSelector      uint64
	DestinationChainSelector uint64
	DestEVMChainID           uint64
	DestinationStartBlock    uint64
	TokenPricesUSDPipeline   string
	PriceGetterConfigJson    string
	USDCAttestationAPI       string
	USDCCfg                  *config.USDCConfig
}

func (j JobSpecInput) Validate() error {
	if err := cldf.IsValidChainSelector(j.SourceChainSelector); err != nil {
		return fmt.Errorf("SourceChainSelector is invalid: %w", err)
	}
	if err := cldf.IsValidChainSelector(j.DestinationChainSelector); err != nil {
		return fmt.Errorf("DestinationChainSelector is invalid: %w", err)
	}
	if j.TokenPricesUSDPipeline == "" && j.PriceGetterConfigJson == "" {
		return errors.New("TokenPricesUSDPipeline or PriceGetterConfigJson is required")
	}
	if j.USDCCfg != nil {
		if err := j.USDCCfg.ValidateUSDCConfig(); err != nil {
			return fmt.Errorf("USDCCfg is invalid: %w", err)
		}
		if j.USDCAttestationAPI == "" {
			return errors.New("USDCAttestationAPI is required")
		}
	}
	return nil
}

func JobSpecsForLanesChangeset(env cldf.Environment, c JobSpecsForLanesConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid JobSpecsForLanesConfig: %w", err)
	}
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	nodesToJobSpecs, err := jobSpecsForLane(env, state, c)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	// Now we propose the job specs to the offchain system.
	var Jobs []cldf.ProposedJob
	for nodeID, jobs := range nodesToJobSpecs {
		for _, job := range jobs {
			Jobs = append(Jobs, cldf.ProposedJob{
				Node: nodeID,
				Spec: job,
			})
			res, err := env.Offchain.ProposeJob(env.GetContext(),
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			if err != nil {
				// If we fail to propose a job, we should return an error and the jobs we've already proposed.
				// This is so that we can retry the proposal with manual intervention.
				// JOBID will be empty if the proposal failed.
				return cldf.ChangesetOutput{
					Jobs: Jobs,
				}, fmt.Errorf("failed to propose job %s: %w", job, err)
			}
			Jobs[len(Jobs)-1].JobID = res.Proposal.JobId
		}
	}
	return cldf.ChangesetOutput{
		Jobs: Jobs,
	}, nil
}

func jobSpecsForLane(
	env cldf.Environment,
	state stateview.CCIPOnChainState,
	lanesCfg JobSpecsForLanesConfig,
) (map[string][]string, error) {
	nodes, err := deployment.NodeInfo(env.NodeIDs, env.Offchain)
	if err != nil {
		return nil, err
	}
	nodesToJobSpecs := make(map[string][]string)
	evmChains := env.BlockChains.EVMChains()
	for _, node := range nodes {
		var specs []string
		for _, cfg := range lanesCfg.Configs {
			destChainState := state.MustGetEVMChainState(cfg.DestinationChainSelector)
			sourceChain := evmChains[cfg.SourceChainSelector]
			destChain := evmChains[cfg.DestinationChainSelector]

			ccipJobParam := integrationtesthelpers.CCIPJobSpecParams{
				OffRamp:                destChainState.EVM2EVMOffRamp[cfg.SourceChainSelector].Address(),
				CommitStore:            destChainState.CommitStore[cfg.SourceChainSelector].Address(),
				SourceChainName:        sourceChain.Name(),
				DestChainName:          destChain.Name(),
				DestEvmChainId:         cfg.DestEVMChainID,
				TokenPricesUSDPipeline: cfg.TokenPricesUSDPipeline,
				PriceGetterConfig:      cfg.PriceGetterConfigJson,
				DestStartBlock:         cfg.DestinationStartBlock,
				USDCAttestationAPI:     cfg.USDCAttestationAPI,
				USDCConfig:             cfg.USDCCfg,
				P2PV2Bootstrappers:     nodes.BootstrapLocators(),
			}
			if !node.IsBootstrap {
				ocrCfg, found := node.OCRConfigForChainSelector(cfg.DestinationChainSelector)
				if !found {
					return nil, fmt.Errorf("OCR config not found for chain %s", destChain.String())
				}
				ocrKeyBundleID := ocrCfg.KeyBundleID
				transmitterID := ocrCfg.TransmitAccount
				commitSpec, err := ccipJobParam.CommitJobSpec()
				if err != nil {
					return nil, fmt.Errorf("failed to generate commit job spec for source %s and destination %s: %w",
						sourceChain.String(), destChain.String(), err)
				}
				commitSpec.OCR2OracleSpec.OCRKeyBundleID.SetValid(ocrKeyBundleID)
				commitSpec.OCR2OracleSpec.TransmitterID.SetValid(string(transmitterID))
				commitSpecStr, err := commitSpec.String()
				if err != nil {
					return nil, fmt.Errorf("failed to convert commit job spec to string for source %s and destination %s: %w",
						sourceChain.String(), destChain.String(), err)
				}
				execSpec, err := ccipJobParam.ExecutionJobSpec()
				if err != nil {
					return nil, fmt.Errorf("failed to generate execution job spec for source %s and destination %s: %w",
						sourceChain.String(), destChain.String(), err)
				}
				execSpec.OCR2OracleSpec.OCRKeyBundleID.SetValid(ocrKeyBundleID)
				execSpec.OCR2OracleSpec.TransmitterID.SetValid(string(transmitterID))
				execSpecStr, err := execSpec.String()
				if err != nil {
					return nil, fmt.Errorf("failed to convert execution job spec to string for source %s and destination %s: %w",
						sourceChain.String(), destChain.String(), err)
				}
				specs = append(specs, commitSpecStr, execSpecStr)
			} else {
				bootstrapSpec := ccipJobParam.BootstrapJob(destChainState.CommitStore[cfg.SourceChainSelector].Address().String())
				bootstrapSpecStr, err := bootstrapSpec.String()
				if err != nil {
					return nil, fmt.Errorf("failed to convert bootstrap job spec to string for source %s and destination %s: %w",
						sourceChain.String(), destChain.String(), err)
				}
				specs = append(specs, bootstrapSpecStr)
			}
		}
		nodesToJobSpecs[node.NodeID] = append(nodesToJobSpecs[node.NodeID], specs...)
	}
	return nodesToJobSpecs, nil
}
