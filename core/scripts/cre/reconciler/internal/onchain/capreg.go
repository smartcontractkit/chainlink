package onchain

import (
	"context"
	"fmt"
	"maps"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-deployments-framework/offchain"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	feature_set "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/sets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

func (d *Deployer) configureCapReg(
	topology *cre.Topology,
	creEnv *cre.Environment,
	chainSelector uint64,
	state *domain.StateFile,
) (crecontracts.CapabilityRegistry, error) {
	d.log.Info().Msg("P4: Configuring Capabilities Registry")

	if creEnv.CldfEnvironment.Offchain == nil {
		return nil, errors.New("JD client is required to configure Capabilities Registry")
	}

	capRegAddrHex := state.GetAddress(keystone_changeset.CapabilitiesRegistry.String())
	if capRegAddrHex == "" {
		return nil, errors.New("capabilities registry address not found in state")
	}
	capRegAddr := common.HexToAddress(capRegAddrHex)

	input := cre.ConfigureCapabilityRegistryInput{
		ChainSelector:                   chainSelector,
		Topology:                        topology,
		CldEnv:                          creEnv.CldfEnvironment,
		NodeSets:                        topology.NodeSets(),
		Blockchains:                     creEnv.Blockchains,
		Provider:                        creEnv.Provider,
		CapabilitiesRegistryAddress:     &capRegAddr,
		DONCapabilityWithConfigs:        d.donsCapabilities,
		CapabilityToOCR3Config:          d.capabilityToOCR3Config,
		CapabilityToExtraSignerFamilies: d.capabilityToExtraSignerFamilies,
	}

	capReg, err := crecontracts.ExecuteConfigureCapabilitiesRegistry(input)
	if err != nil {
		return nil, err
	}

	d.log.Info().Str("capReg", capRegAddrHex).Msg("Capabilities Registry configured")
	return capReg, nil
}

func (d *Deployer) resolveDONIDs(
	capReg crecontracts.CapabilityRegistry,
	desired *domain.DesiredState,
	cv *domain.ChartValues,
	state *domain.StateFile,
) error {
	d.log.Info().Msg("P5: Resolving DON IDs")

	donNames := make([]string, 0, len(desired.DONs))
	for _, don := range desired.DONs {
		if don.IsBootstrapOnly(cv) || don.IsGatewayDon() {
			continue
		}
		donNames = append(donNames, don.Name)
	}

	ids, err := crecontracts.ResolveContractDonIDs(capReg, donNames)
	if err != nil {
		return errors.Wrap(err, "failed to resolve DON IDs from contract")
	}

	if state.DONIDs == nil {
		state.DONIDs = make(map[string]uint64)
	}
	for name, id := range ids {
		state.DONIDs[name] = uint64(id)
		d.log.Info().Str("don", name).Uint64("id", uint64(id)).Msg("Resolved DON ID")
	}

	return nil
}

func (d *Deployer) runPreEnvStartup(ctx context.Context, topology *cre.Topology, creEnv *cre.Environment) error {
	features := feature_set.New()

	d.donsCapabilities = make(map[uint64][]keystone_changeset.DONCapabilityWithConfig)
	d.capabilityToOCR3Config = make(map[string]*ocr3.OracleConfig)
	d.capabilityToExtraSignerFamilies = make(map[string][]string)

	for _, feature := range features.List() {
		for _, donMetadata := range topology.DonsMetadataWithFlag(feature.Flag()) {
			d.log.Info().Str("feature", feature.Flag()).Str("don", donMetadata.Name).Msg("Running PreEnvStartup")

			output, err := feature.PreEnvStartup(ctx, d.log, donMetadata, topology, creEnv)
			if err != nil {
				return errors.Wrapf(err, "PreEnvStartup failed for %s on DON %s", feature.Flag(), donMetadata.Name)
			}
			if output == nil {
				continue
			}

			d.donsCapabilities[donMetadata.ID] = append(d.donsCapabilities[donMetadata.ID], output.DONCapabilityWithConfig...)
			maps.Copy(d.capabilityToOCR3Config, output.CapabilityToOCR3Config)
			for cap, families := range output.CapabilityToExtraSignerFamilies {
				d.capabilityToExtraSignerFamilies[cap] = append([]string(nil), families...)
			}

			d.log.Info().
				Str("feature", feature.Flag()).
				Str("don", donMetadata.Name).
				Int("capabilities", len(output.DONCapabilityWithConfig)).
				Int("ocr3Configs", len(output.CapabilityToOCR3Config)).
				Msg("PreEnvStartup complete")
		}
	}

	return nil
}

type capRegWorker struct {
	donName  string
	nodeName string
	worker   *cre.NodeMetadata
	donMeta  *cre.DonMetadata
}

func capRegWorkerNodes(topology *cre.Topology, runtime map[string]domain.NodeRuntimeInfo) []capRegWorker {
	if topology == nil || topology.DonsMetadata == nil {
		return nil
	}

	var workers []capRegWorker
	for _, donMeta := range topology.DonsMetadata.List() {
		if flags.HasNoOtherFlags(donMeta.Flags, []string{cre.GatewayDON, cre.BootstrapDON}) {
			continue
		}
		donWorkers, err := donMeta.Workers()
		if err != nil {
			continue
		}
		for _, worker := range donWorkers {
			workers = append(workers, capRegWorker{
				donName:  donMeta.Name,
				worker:   worker,
				donMeta:  donMeta,
				nodeName: chartNodeNameForWorker(worker, runtime),
			})
		}
	}
	return workers
}

func verifyCapRegNodeInfo(
	offchainClient offchain.Client,
	chainSelector uint64,
	topology *cre.Topology,
	runtime map[string]domain.NodeRuntimeInfo,
) error {
	if offchainClient == nil {
		return errors.New("JD client is required to verify node chain configs for CapReg")
	}

	lookupIDs := make([]string, 0)
	workerLabels := make([]string, 0)
	for _, w := range capRegWorkerNodes(topology, runtime) {
		if w.worker == nil || w.worker.Keys == nil {
			return fmt.Errorf("DON %s: worker node has no keys for CapReg verification", w.donName)
		}
		p2pID := w.worker.Keys.PeerID()
		if p2pID == "" {
			return fmt.Errorf("DON %s: worker node missing P2P ID for CapReg verification", w.donName)
		}
		lookupIDs = append(lookupIDs, p2pID)
		workerLabels = append(workerLabels, fmt.Sprintf("%s/%s", w.donName, w.nodeName))
	}

	nodes, err := deployment.NodeInfo(lookupIDs, offchainClient)
	if err != nil {
		return errors.Wrapf(err, "chain selector %d: failed to load node info for CapReg preflight", chainSelector)
	}
	if len(nodes) != len(lookupIDs) {
		return fmt.Errorf(
			"chain selector %d: expected %d CapReg worker nodes in JD, got %d",
			chainSelector, len(lookupIDs), len(nodes),
		)
	}

	for i, node := range nodes {
		if _, ok := node.OCRConfigForChainSelector(chainSelector); !ok {
			return fmt.Errorf(
				"chain selector %d: node %s (lookup %s) missing registry-chain OCR config in JD — JD chain config preparation may have failed",
				chainSelector, workerLabels[i], lookupIDs[i],
			)
		}
	}
	return nil
}
