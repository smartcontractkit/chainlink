package aptos

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"maps"
	"strconv"
	"strings"

	"dario.cat/mergo"
	pkgerrors "github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	crejobops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	jobtypes "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/jobhelpers"
)

func proposeAptosWorkerSpecs(
	ctx context.Context,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
	nodeSet cre.NodeSetWithCapabilityConfigs,
	enabledChainIDs []uint64,
) (map[string][]string, error) {
	bootstrapPeers, err := bootstrapPeersForDons(dons)
	if err != nil {
		return nil, err
	}

	workerMetadata, err := don.Metadata().Workers()
	if err != nil {
		return nil, fmt.Errorf("failed to collect Aptos worker metadata for DON %q: %w", don.Name, err)
	}
	p2pToTransmitterMap, err := p2pToTransmitterMapForWorkers(workerMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to collect Aptos worker transmitters for DON %q: %w", don.Name, err)
	}

	results := make([]map[string][]string, len(enabledChainIDs))
	group, _ := errgroup.WithContext(ctx)
	group.SetLimit(jobhelpers.Parallelism(len(enabledChainIDs)))

	for i, chainID := range enabledChainIDs {
		group.Go(func() error {
			mergedSpecs, buildErr := proposeAptosWorkerSpecsForChain(
				creEnv,
				don.Name,
				nodeSet,
				flag,
				chainID,
				bootstrapPeers,
				p2pToTransmitterMap,
			)
			if buildErr != nil {
				return buildErr
			}
			results[i] = maps.Clone(mergedSpecs)
			return nil
		})
	}

	waitErr := group.Wait()
	if waitErr != nil {
		return nil, waitErr
	}

	specs, err := jobhelpers.MergeSpecsByIndex(results)
	if err != nil {
		return nil, err
	}
	return specs, nil
}

func proposeAptosWorkerSpecsForChain(
	creEnv *cre.Environment,
	donName string,
	nodeSet cre.NodeSetWithCapabilityConfigs,
	flag string,
	chainID uint64,
	bootstrapPeers []string,
	p2pToTransmitterMap map[string]string,
) (map[string][]string, error) {
	aptosChain, err := findAptosChainByChainID(creEnv.Blockchains, chainID)
	if err != nil {
		return nil, err
	}

	capabilityConfig, err := cre.ResolveCapabilityConfig(nodeSet, flag, cre.ChainCapabilityScope(chainID))
	if err != nil {
		return nil, fmt.Errorf("could not resolve capability config for '%s' on chain %d: %w", flag, chainID, err)
	}
	command, err := standardcapability.GetCommand(capabilityConfig.BinaryName)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to get command for Aptos capability")
	}

	forwarderAddress := mustForwarderAddress(creEnv.CldfEnvironment.DataStore, aptosChain.ChainSelector())
	methodSettings, err := resolveMethodConfigSettings(capabilityConfig.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve Aptos method config settings for chain %d: %w", chainID, err)
	}
	configStr, err := buildWorkerConfigJSON(chainID, forwarderAddress, methodSettings, p2pToTransmitterMap, true)
	if err != nil {
		return nil, fmt.Errorf("failed to build Aptos worker config: %w", err)
	}

	workerInput, err := newAptosWorkerJobInput(creEnv, donName, command, configStr, bootstrapPeers, aptosChain.ChainSelector(), chainID)
	if err != nil {
		return nil, err
	}

	proposer := jobs.ProposeJobSpec{}
	verifyErr := proposer.VerifyPreconditions(*creEnv.CldfEnvironment, workerInput)
	if verifyErr != nil {
		return nil, fmt.Errorf("precondition verification failed for Aptos worker job: %w", verifyErr)
	}
	workerReport, err := proposer.Apply(*creEnv.CldfEnvironment, workerInput)
	if err != nil {
		return nil, fmt.Errorf("failed to propose Aptos worker job spec: %w", err)
	}

	mergedSpecs := make(map[string][]string)
	for _, report := range workerReport.Reports {
		out, ok := report.Output.(crejobops.ProposeStandardCapabilityJobOutput)
		if !ok {
			return nil, fmt.Errorf("unable to cast to ProposeStandardCapabilityJobOutput, actual type: %T", report.Output)
		}
		if err := mergo.Merge(&mergedSpecs, out.Specs, mergo.WithAppendSlice); err != nil {
			return nil, fmt.Errorf("failed to merge Aptos worker job specs: %w", err)
		}
	}

	return mergedSpecs, nil
}

func newAptosWorkerJobInput(
	creEnv *cre.Environment,
	donName string,
	command string,
	configStr string,
	bootstrapPeers []string,
	aptosChainSelector uint64,
	chainID uint64,
) (jobs.ProposeJobSpecInput, error) {
	capRegVersion, ok := creEnv.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()]
	if !ok {
		return jobs.ProposeJobSpecInput{}, errors.New("CapabilitiesRegistry version not found in contract versions")
	}

	return jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: cre.EnvironmentName,
		DONName:     donName,
		JobName:     "aptos-worker-" + strconv.FormatUint(chainID, 10),
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: donName},
		},
		Template: jobtypes.Aptos,
		Inputs: jobtypes.JobSpecInput{
			"command":            command,
			"config":             configStr,
			"chainSelectorEVM":   creEnv.RegistryChainSelector,
			"chainSelectorAptos": aptosChainSelector,
			"bootstrapPeers":     bootstrapPeers,
			"useCapRegOCRConfig": true,
			"capRegVersion":      capRegVersion.String(),
		},
	}, nil
}

func donOraclePublicKeys(ctx context.Context, don *cre.Don) ([][]byte, error) {
	workers, err := don.Workers()
	if err != nil {
		return nil, fmt.Errorf("failed to list worker nodes for DON %q: %w", don.Name, err)
	}

	oracles := make([][]byte, 0, len(workers))
	for _, worker := range workers {
		ocr2ID := ""
		if worker.Keys != nil && worker.Keys.OCR2BundleIDs != nil {
			ocr2ID = worker.Keys.OCR2BundleIDs[chainselectors.FamilyAptos]
		}
		if ocr2ID == "" {
			fetchedID, err := worker.Clients.GQLClient.FetchOCR2KeyBundleID(ctx, strings.ToUpper(chainselectors.FamilyAptos))
			if err != nil {
				return nil, fmt.Errorf("missing Aptos OCR2 bundle id for worker %q in DON %q and fallback fetch failed: %w", worker.Name, don.Name, err)
			}
			if fetchedID == "" {
				return nil, fmt.Errorf("missing Aptos OCR2 bundle id for worker %q in DON %q", worker.Name, don.Name)
			}
			ocr2ID = fetchedID
			if worker.Keys != nil {
				if worker.Keys.OCR2BundleIDs == nil {
					worker.Keys.OCR2BundleIDs = make(map[string]string)
				}
				worker.Keys.OCR2BundleIDs[chainselectors.FamilyAptos] = ocr2ID
			}
		}

		exported, err := worker.ExportOCR2Keys(ocr2ID)
		if err != nil {
			return nil, fmt.Errorf("failed to export Aptos OCR2 key for worker %q (bundle %s): %w", worker.Name, ocr2ID, err)
		}
		pubkey, err := parseOCR2OnchainPublicKey(exported.OnchainPublicKey)
		if err != nil {
			return nil, fmt.Errorf("invalid Aptos OCR2 onchain public key for worker %q: %w", worker.Name, err)
		}
		oracles = append(oracles, pubkey)
	}

	return oracles, nil
}

func p2pToTransmitterMapForWorkers(workers []*cre.NodeMetadata) (map[string]string, error) {
	if len(workers) == 0 {
		return nil, pkgerrors.New("no DON worker nodes provided")
	}

	p2pToTransmitterMap := make(map[string]string)
	for _, worker := range workers {
		if worker.Keys == nil || worker.Keys.P2PKey == nil {
			return nil, fmt.Errorf("missing P2P key for worker index %d", worker.Index)
		}

		account := ""
		if worker.Keys.Aptos != nil {
			account = worker.Keys.Aptos.Account
		}
		if account == "" {
			return nil, fmt.Errorf("missing Aptos account for worker index %d", worker.Index)
		}

		transmitter, err := normalizeTransmitter(account)
		if err != nil {
			return nil, fmt.Errorf("invalid Aptos transmitter for worker index %d: %w", worker.Index, err)
		}

		peerKey := hex.EncodeToString(worker.Keys.P2PKey.PeerID[:])
		p2pToTransmitterMap[peerKey] = transmitter
	}

	if len(p2pToTransmitterMap) == 0 {
		return nil, pkgerrors.New("no Aptos transmitters found for DON workers")
	}

	return p2pToTransmitterMap, nil
}

func aptosDonIDUint32(donID uint64) (uint32, error) {
	if donID > uint64(^uint32(0)) {
		return 0, fmt.Errorf("don id %d exceeds u32", donID)
	}
	return uint32(donID), nil
}

func parseOCR2OnchainPublicKey(hexValue string) ([]byte, error) {
	trimmed := strings.TrimPrefix(strings.TrimSpace(hexValue), "ocr2on_aptos_")
	decoded, err := hex.DecodeString(trimmed)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}
