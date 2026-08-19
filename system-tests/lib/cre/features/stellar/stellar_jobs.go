package stellar

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"dario.cat/mergo"
	pkgerrors "github.com/pkg/errors"

	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
)

// proposeStellarStandardCapabilityJobs proposes the standard-capability job that
// runs the Stellar chain capability.
func proposeStellarStandardCapabilityJobs(
	_ context.Context,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
	nodeSet cre.NodeSetWithCapabilityConfigs,
) (map[string][]string, error) {
	bootstrapPeers, err := bootstrapPeersForDons(dons)
	if err != nil {
		return nil, err
	}

	stellarChain := extractStellarFromEnv(creEnv)
	if stellarChain == nil {
		return nil, pkgerrors.New("no stellar blockchain found in environment")
	}

	return proposeStellarStandardCapabilityJobsForChain(
		creEnv,
		don.Name,
		nodeSet,
		stellarChain.StellarChainID(),
		stellarChain.ChainSelector(),
		bootstrapPeers,
		don.HasFlag(cre.WorkflowDON),
	)
}

func proposeStellarStandardCapabilityJobsForChain(
	creEnv *cre.Environment,
	donName string,
	nodeSet cre.NodeSetWithCapabilityConfigs,
	chainID string,
	chainSelector uint64,
	bootstrapPeers []string,
	isLocal bool,
) (map[string][]string, error) {
	capabilityConfig, err := cre.ResolveCapabilityConfig(nodeSet, flag, cre.CapabilityScope{})
	if err != nil {
		return nil, fmt.Errorf("could not resolve capability config for '%s': %w", flag, err)
	}

	command, err := standardcapability.GetCommand(capabilityConfig.BinaryName)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to get command for Stellar capability")
	}

	methodSettings, err := resolveMethodConfigSettings(capabilityConfig.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve Stellar method config settings: %w", err)
	}

	forwarderAddress := mustForwarderAddress(creEnv.CldfEnvironment.DataStore, chainSelector)

	configStr, err := buildJobConfigJSON(chainID, forwarderAddress, methodSettings, isLocal)
	if err != nil {
		return nil, fmt.Errorf("failed to build Stellar standard capability job config: %w", err)
	}

	jobInput, err := newStellarStandardCapabilityJobInput(
		creEnv,
		donName,
		command,
		configStr,
		bootstrapPeers,
		chainSelector,
	)
	if err != nil {
		return nil, err
	}

	proposer := cre_jobs.ProposeJobSpec{}
	if verifyErr := proposer.VerifyPreconditions(*creEnv.CldfEnvironment, jobInput); verifyErr != nil {
		return nil, fmt.Errorf("precondition verification failed for Stellar standard capability job: %w", verifyErr)
	}

	jobReport, err := proposer.Apply(*creEnv.CldfEnvironment, jobInput)
	if err != nil {
		return nil, fmt.Errorf("failed to propose Stellar standard capability job spec: %w", err)
	}

	mergedSpecs := make(map[string][]string)
	for _, report := range jobReport.Reports {
		out, ok := report.Output.(cre_jobs_ops.ProposeStandardCapabilityJobOutput)
		if !ok {
			return nil, fmt.Errorf("unable to cast to ProposeStandardCapabilityJobOutput, actual type: %T", report.Output)
		}
		if mergeErr := mergo.Merge(&mergedSpecs, out.Specs, mergo.WithAppendSlice); mergeErr != nil {
			return nil, fmt.Errorf("failed to merge Stellar standard capability job specs: %w", mergeErr)
		}
	}

	return mergedSpecs, nil
}

func newStellarStandardCapabilityJobInput(
	creEnv *cre.Environment,
	donName string,
	command string,
	configStr string,
	bootstrapPeers []string,
	stellarChainSelector uint64,
) (cre_jobs.ProposeJobSpecInput, error) {
	capRegVersion, ok := creEnv.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()]
	if !ok {
		return cre_jobs.ProposeJobSpecInput{}, pkgerrors.New("CapabilitiesRegistry version not found in contract versions")
	}

	return cre_jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: cre.EnvironmentName,
		DONName:     donName,
		JobName:     "stellar-worker-" + strconv.FormatUint(stellarChainSelector, 10),
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: donName},
		},
		Template: job_types.Stellar,
		Inputs: job_types.JobSpecInput{
			"command":              command,
			"config":               configStr,
			"chainSelectorEVM":     creEnv.RegistryChainSelector,
			"chainSelectorStellar": stellarChainSelector,
			"bootstrapPeers":       bootstrapPeers,
			"useCapRegOCRConfig":   true,
			"capRegVersion":        capRegVersion.String(),
		},
	}, nil
}

func bootstrapPeersForDons(dons *cre.Dons) ([]string, error) {
	bootstrapNode, ok := dons.Bootstrap()
	if !ok {
		return nil, pkgerrors.New("bootstrap node not found; required for Stellar OCR bootstrap peers")
	}

	return []string{
		fmt.Sprintf("%s@%s:%d", strings.TrimPrefix(bootstrapNode.Keys.PeerID(), "p2p_"), bootstrapNode.Host, cre.OCRPeeringPort),
	}, nil
}

func nodeSetForDON(dons *cre.Dons, donName string) (cre.NodeSetWithCapabilityConfigs, error) {
	for _, ns := range dons.AsNodeSetWithChainCapabilities() {
		if ns.GetName() == donName {
			return ns, nil
		}
	}
	return nil, fmt.Errorf("could not find node set for Don named '%s'", donName)
}
