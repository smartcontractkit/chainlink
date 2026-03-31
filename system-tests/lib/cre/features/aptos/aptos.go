package aptos

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	crejobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

const (
	flag                    = cre.AptosCapability
	forwarderContractType   = "AptosForwarder"
	forwarderConfigVersion  = 1
	capabilityVersion       = "1.0.0"
	capabilityLabelPrefix   = "aptos:ChainSelector:"
	specConfigP2PMapKey     = "p2pToTransmitterMap"
	specConfigScheduleKey   = "transmissionSchedule"
	specConfigDeltaStageKey = "deltaStage"
	legacyTransmittersKey   = "aptosTransmitters"
	requestTimeoutKey       = "RequestTimeout"
	deltaStageKey           = "DeltaStage"
	transmissionScheduleKey = "TransmissionSchedule"
	forwarderQualifier      = ""
	defaultWriteDeltaStage  = 500*time.Millisecond + 1*time.Second
	defaultRequestTimeout   = 30 * time.Second
)

var forwarderContractVersion = semver.MustParse("1.0.0")

type Aptos struct{}

type methodConfigSettings struct {
	RequestTimeout       time.Duration
	DeltaStage           time.Duration
	TransmissionSchedule capabilitiespb.TransmissionSchedule
}

func (a *Aptos) Flag() cre.CapabilityFlag {
	return flag
}

func CapabilityLabel(chainSelector uint64) string {
	return capabilityLabelPrefix + strconv.FormatUint(chainSelector, 10)
}

func (a *Aptos) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	_ *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	enabledChainIDs, err := don.MustNodeSet().GetEnabledChainIDsForCapability(flag)
	if err != nil {
		return nil, fmt.Errorf("could not find enabled chainIDs for '%s' in don '%s': %w", flag, don.Name, err)
	}
	if len(enabledChainIDs) == 0 {
		return nil, nil
	}

	forwardersByChainID, err := ensureForwardersForChains(ctx, testLogger, creEnv, enabledChainIDs)
	if err != nil {
		return nil, err
	}
	err = patchNodeTOML(don, forwardersByChainID)
	if err != nil {
		return nil, err
	}

	workers, err := don.Workers()
	if err != nil {
		return nil, err
	}
	p2pToTransmitterMap, err := p2pToTransmitterMapForWorkers(workers)
	if err != nil {
		return nil, fmt.Errorf("failed to collect Aptos worker transmitters for DON %q from metadata: %w", don.Name, err)
	}

	caps, capabilityToOCR3Config, capabilityLabels, err := buildCapabilityRegistrations(don, creEnv.Blockchains, enabledChainIDs, p2pToTransmitterMap)
	if err != nil {
		return nil, err
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: caps,
		CapabilityToOCR3Config:  capabilityToOCR3Config,
		CapabilityToExtraSignerFamilies: cre.CapabilityToExtraSignerFamilies(
			cre.OCRExtraSignerFamiliesForFamily(chainselectors.FamilyAptos),
			capabilityLabels...,
		),
	}, nil
}

func (a *Aptos) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	nodeSet, err := nodeSetForDON(dons, don.Name)
	if err != nil {
		return err
	}

	enabledChainIDs, err := nodeSet.GetEnabledChainIDsForCapability(flag)
	if err != nil {
		return fmt.Errorf("could not find enabled chainIDs for '%s' in don '%s': %w", flag, don.Name, err)
	}
	if len(enabledChainIDs) == 0 {
		return nil
	}

	err = configureForwarders(ctx, testLogger, don, creEnv, enabledChainIDs)
	if err != nil {
		return err
	}

	specs, err := proposeAptosWorkerSpecs(ctx, don, dons, creEnv, nodeSet, enabledChainIDs)
	if err != nil {
		return err
	}
	if len(specs) == 0 {
		return nil
	}
	err = crejobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs)
	if err != nil {
		return fmt.Errorf("failed to approve Aptos jobs: %w", err)
	}
	return nil
}

func buildCapabilityRegistrations(
	don *cre.DonMetadata,
	blockchains []creblockchains.Blockchain,
	enabledChainIDs []uint64,
	p2pToTransmitterMap map[string]string,
) ([]keystone_changeset.DONCapabilityWithConfig, map[string]*ocr3.OracleConfig, []string, error) {
	caps := make([]keystone_changeset.DONCapabilityWithConfig, 0, len(enabledChainIDs))
	capabilityToOCR3Config := make(map[string]*ocr3.OracleConfig, len(enabledChainIDs))
	capabilityLabels := make([]string, 0, len(enabledChainIDs))

	for _, chainID := range enabledChainIDs {
		aptosChain, err := findAptosChainByChainID(blockchains, chainID)
		if err != nil {
			return nil, nil, nil, err
		}

		labelledName := CapabilityLabel(aptosChain.ChainSelector())
		capabilityConfig, err := cre.ResolveCapabilityConfig(don.MustNodeSet(), flag, cre.ChainCapabilityScope(chainID))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("could not resolve capability config for '%s' on chain %d: %w", flag, chainID, err)
		}
		capConfig, err := BuildCapabilityConfig(capabilityConfig.Values, p2pToTransmitterMap, don.HasOnlyLocalCapabilities())
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to build Aptos capability config for capability %s: %w", labelledName, err)
		}

		caps = append(caps, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   labelledName,
				Version:        capabilityVersion,
				CapabilityType: 1,
			},
			Config:             capConfig,
			UseCapRegOCRConfig: true,
		})
		capabilityLabels = append(capabilityLabels, labelledName)
		capabilityToOCR3Config[labelledName] = crecontracts.DefaultChainCapabilityOCR3Config()
	}

	return caps, capabilityToOCR3Config, capabilityLabels, nil
}

func nodeSetForDON(dons *cre.Dons, donName string) (cre.NodeSetWithCapabilityConfigs, error) {
	for _, ns := range dons.AsNodeSetWithChainCapabilities() {
		if ns.GetName() == donName {
			return ns, nil
		}
	}
	return nil, fmt.Errorf("could not find node set for Don named '%s'", donName)
}

func bootstrapPeersForDons(dons *cre.Dons) ([]string, error) {
	bootstrapNode, ok := dons.Bootstrap()
	if !ok {
		return nil, pkgerrors.New("bootstrap node not found; required for Aptos OCR bootstrap peers")
	}

	return []string{
		fmt.Sprintf("%s@%s:%d", strings.TrimPrefix(bootstrapNode.Keys.PeerID(), "p2p_"), bootstrapNode.Host, cre.OCRPeeringPort),
	}, nil
}

var (
	_ cre.Feature = (*Aptos)(nil)
)
