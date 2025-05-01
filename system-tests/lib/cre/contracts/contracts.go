package contracts

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/deployment"
	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	df_changeset_types "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	workflow_registry_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	keystonenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
)

// deprecated, use ComputeCapabilityFactoryFn, OCR3CapabilityFactoryFn, CronCapabilityFactoryFn instead
var DefaultCapabilityFactoryFn = func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
	var capabilities []keystone_changeset.DONCapabilityWithConfig

	if flags.HasFlag(donFlags, types.CronCapability) {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "cron-trigger",
				Version:        "1.0.0",
				CapabilityType: 0, // TRIGGER
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	if flags.HasFlag(donFlags, types.CustomComputeCapability) {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "custom-compute",
				Version:        "1.0.0",
				CapabilityType: 1, // ACTION
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	if flags.HasFlag(donFlags, types.OCR3Capability) {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "offchain_reporting",
				Version:        "1.0.0",
				CapabilityType: 2, // CONSENSUS
				ResponseType:   0, // REPORT
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities
}

// deprecated, use capabilities.webapi.WebAPICapabilityFactoryFn instead
var WebAPICapabilityFactoryFn = func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
	var capabilities []keystone_changeset.DONCapabilityWithConfig

	if flags.HasFlag(donFlags, types.WebAPITriggerCapability) {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "web-api-trigger",
				Version:        "1.0.0",
				CapabilityType: 0, // TRIGGER
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	if flags.HasFlag(donFlags, types.WebAPITargetCapability) {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "web-api-target",
				Version:        "1.0.0",
				CapabilityType: 3, // TARGET
				ResponseType:   1, // OBSERVATION_IDENTICAL
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities
}

// deprecated, use capabilities.chainwriter.ChainWriterCapabilityFactory instead
var ChainWriterCapabilityFactory = func(chainID uint64) func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
	return func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
		var capabilities []keystone_changeset.DONCapabilityWithConfig

		fullName := corevm.GenerateWriteTargetName(chainID)
		splitName := strings.Split(fullName, "@")

		if flags.HasFlag(donFlags, types.WriteEVMCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   splitName[0],
					Version:        splitName[1],
					CapabilityType: 3, // TARGET
					ResponseType:   1, // OBSERVATION_IDENTICAL
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		return capabilities
	}
}

// deprecated, use capabilities.chainreader.ChainReaderCapabilityFactory instead
var ChainReaderCapabilityFactory = func(chainID uint64, chainFamily string) func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
	return func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
		var capabilities []keystone_changeset.DONCapabilityWithConfig

		if flags.HasFlag(donFlags, types.LogTriggerCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   fmt.Sprintf("log-event-trigger-%s-%d", chainFamily, chainID),
					Version:        "1.0.0",
					CapabilityType: 0, // TRIGGER
					ResponseType:   0, // REPORT
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		if flags.HasFlag(donFlags, types.ReadContractCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   fmt.Sprintf("read-contract-%s-%d", chainFamily, chainID),
					Version:        "1.0.0",
					CapabilityType: 1, // ACTION
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		return capabilities
	}
}

func ConfigureKeystone(input types.ConfigureKeystoneInput, capabilityFactoryFns []types.DONCapabilityWithConfigFactoryFn) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "input validation failed")
	}

	donCapabilities := make([]keystone_changeset.DonCapabilities, 0, len(input.Topology.DonsMetadata))

	for _, donMetadata := range input.Topology.DonsMetadata {
		// if it's only a gateway DON, we don't want to register it with the Capabilities Registry
		// since it doesn't have any capabilities
		if flags.HasOnlyOneFlag(donMetadata.Flags, types.GatewayDON) {
			continue
		}

		var capabilities []keystone_changeset.DONCapabilityWithConfig

		// check what capabilities each DON has and register them with Capabilities Registry contract
		for _, factoryFn := range capabilityFactoryFns {
			capabilities = append(capabilities, factoryFn(donMetadata.Flags)...)
		}

		workerNodes, workerNodesErr := node.FindManyWithLabel(donMetadata.NodesMetadata, &types.Label{
			Key:   node.NodeTypeKey,
			Value: types.WorkerNode,
		}, node.EqualLabels)

		if workerNodesErr != nil {
			return errors.Wrap(workerNodesErr, "failed to find worker nodes")
		}

		donPeerIDs := make([]string, len(workerNodes))
		for i, node := range workerNodes {
			p2pID, err := keystonenode.ToP2PID(node, keystonenode.NoOpTransformFn)
			if err != nil {
				return errors.Wrapf(err, "failed to get p2p id for node %d", i)
			}

			donPeerIDs[i] = p2pID
		}

		// we only need to assign P2P IDs to NOPs, since `ConfigureInitialContractsChangeset` method
		// will take care of creating DON to Nodes mapping
		nop := keystone_changeset.NOP{
			Name:  fmt.Sprintf("NOP for %s DON", donMetadata.Name),
			Nodes: donPeerIDs,
		}

		forwarderF := (len(workerNodes) - 1) / 3

		if forwarderF == 0 {
			if flags.HasFlag(donMetadata.Flags, types.OCR3Capability) {
				return fmt.Errorf("incorrect number of worker nodes: %d. Resulting F must conform to formula: mod((N-1)/3) = 0", len(workerNodes))
			}
			// for other capabilities, we can use 1 as F
			forwarderF = 1
		}

		donName := donMetadata.Name + "-don"
		donCapabilities = append(donCapabilities, keystone_changeset.DonCapabilities{
			Name:         donName,
			F:            libc.MustSafeUint8(forwarderF),
			Nops:         []keystone_changeset.NOP{nop},
			Capabilities: capabilities,
		})
	}

	var transmissionSchedule []int

	for _, metaDon := range input.Topology.DonsMetadata {
		if flags.HasFlag(metaDon.Flags, types.OCR3Capability) {
			workerNodes, workerNodesErr := node.FindManyWithLabel(metaDon.NodesMetadata, &types.Label{
				Key:   node.NodeTypeKey,
				Value: types.WorkerNode,
			}, node.EqualLabels)

			if workerNodesErr != nil {
				return errors.Wrap(workerNodesErr, "failed to find worker nodes")
			}

			// this schedule makes sure that all worker nodes are transmitting OCR3 reports
			transmissionSchedule = []int{len(workerNodes)}
			break
		}
	}

	if len(transmissionSchedule) == 0 {
		return errors.New("no OCR3-capable DON found in the topology")
	}

	// values supplied by Alexandr Yepishev as the expected values for OCR3 config
	oracleConfig := keystone_changeset.OracleConfig{
		DeltaProgressMillis:               5000,
		DeltaResendMillis:                 5000,
		DeltaInitialMillis:                5000,
		DeltaRoundMillis:                  2000,
		DeltaGraceMillis:                  500,
		DeltaCertifiedCommitRequestMillis: 1000,
		DeltaStageMillis:                  30000,
		MaxRoundsPerEpoch:                 10,
		TransmissionSchedule:              transmissionSchedule,
		MaxDurationQueryMillis:            1000,
		MaxDurationObservationMillis:      1000,
		MaxDurationShouldAcceptMillis:     1000,
		MaxDurationShouldTransmitMillis:   1000,
		MaxFaultyOracles:                  1,
		MaxQueryLengthBytes:               1000000,
		MaxObservationLengthBytes:         1000000,
		MaxReportLengthBytes:              1000000,
		MaxBatchSize:                      1000,
		UniqueReports:                     true,
	}

	cfg := keystone_changeset.InitialContractsCfg{
		RegistryChainSel: input.ChainSelector,
		Dons:             donCapabilities,
		OCR3Config:       &oracleConfig,
	}

	_, err := keystone_changeset.ConfigureInitialContractsChangeset(*input.CldEnv, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to configure initial contracts")
	}

	return nil
}

// values supplied by Alexandr Yepishev as the expected values for OCR3 config
func DefaultOCR3Config(topology *types.Topology) (*keystone_changeset.OracleConfig, error) {
	var transmissionSchedule []int

	for _, metaDon := range topology.DonsMetadata {
		if flags.HasFlag(metaDon.Flags, types.OCR3Capability) {
			workerNodes, workerNodesErr := node.FindManyWithLabel(metaDon.NodesMetadata, &types.Label{
				Key:   node.NodeTypeKey,
				Value: types.WorkerNode,
			}, node.EqualLabels)

			if workerNodesErr != nil {
				return nil, errors.Wrap(workerNodesErr, "failed to find worker nodes")
			}

			// this schedule makes sure that all worker nodes are transmitting OCR3 reports
			transmissionSchedule = []int{len(workerNodes)}
			break
		}
	}

	if len(transmissionSchedule) == 0 {
		return nil, errors.New("no OCR3-capable DON found in the topology")
	}

	oracleConfig := &keystone_changeset.OracleConfig{
		DeltaProgressMillis:               5000,
		DeltaResendMillis:                 5000,
		DeltaInitialMillis:                5000,
		DeltaRoundMillis:                  2000,
		DeltaGraceMillis:                  500,
		DeltaCertifiedCommitRequestMillis: 1000,
		DeltaStageMillis:                  30000,
		MaxRoundsPerEpoch:                 10,
		TransmissionSchedule:              transmissionSchedule,
		MaxDurationQueryMillis:            1000,
		MaxDurationObservationMillis:      1000,
		MaxDurationShouldAcceptMillis:     1000,
		MaxDurationShouldTransmitMillis:   1000,
		MaxFaultyOracles:                  1,
		MaxQueryLengthBytes:               1000000,
		MaxObservationLengthBytes:         1000000,
		MaxReportLengthBytes:              1000000,
		MaxBatchSize:                      1000,
		UniqueReports:                     true,
	}

	return oracleConfig, nil
}

func FindAddressesForChain(addressBook deployment.AddressBook, chainSelector uint64, contractName string) (common.Address, error) {
	addresses, err := addressBook.AddressesForChain(chainSelector)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to get addresses for chain")
	}

	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), contractName) {
			return common.HexToAddress(addrStr), nil
		}
	}

	return common.Address{}, fmt.Errorf("failed to find %s address in the address book for chain %d", contractName, chainSelector)
}

func MustFindAddressesForChain(addressBook deployment.AddressBook, chainSelector uint64, contractName string) common.Address {
	addr, err := FindAddressesForChain(addressBook, chainSelector, contractName)
	if err != nil {
		panic(fmt.Errorf("failed to find %s address in the address book for chain %d", contractName, chainSelector))
	}
	return addr
}

func ConfigureWorkflowRegistry(testLogger zerolog.Logger, input *types.WorkflowRegistryInput) (*types.WorkflowRegistryOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	if input.Out != nil && input.Out.UseCache {
		return input.Out, nil
	}

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	_, err := workflow_registry_changeset.UpdateAllowedDons(*input.CldEnv, &workflow_registry_changeset.UpdateAllowedDonsRequest{
		RegistryChainSel: input.ChainSelector,
		DonIDs:           input.AllowedDonIDs,
		Allowed:          true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to update allowed Dons")
	}

	addresses := make([]string, 0, len(input.WorkflowOwners))
	for _, owner := range input.WorkflowOwners {
		addresses = append(addresses, owner.Hex())
	}

	_, err = workflow_registry_changeset.UpdateAuthorizedAddresses(*input.CldEnv, &workflow_registry_changeset.UpdateAuthorizedAddressesRequest{
		RegistryChainSel: input.ChainSelector,
		Addresses:        addresses,
		Allowed:          true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to update authorized addresses")
	}

	out := &types.WorkflowRegistryOutput{
		ChainSelector:  input.ChainSelector,
		AllowedDonIDs:  input.AllowedDonIDs,
		WorkflowOwners: input.WorkflowOwners,
	}

	input.Out = out
	return out, nil
}

func ConfigureDataFeedsCache(testLogger zerolog.Logger, input *types.ConfigureDataFeedsCacheInput) (*types.ConfigureDataFeedsCacheOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	if input.Out != nil && input.Out.UseCache {
		return input.Out, nil
	}

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	if input.AdminAddress != (common.Address{}) {
		setAdminConfig := df_changeset_types.SetFeedAdminConfig{
			ChainSelector: input.ChainSelector,
			CacheAddress:  input.DataFeedsCacheAddress,
			AdminAddress:  input.AdminAddress,
			IsAdmin:       true,
		}
		_, setAdminErr := df_changeset.RunChangeset(df_changeset.SetFeedAdminChangeset, *input.CldEnv, setAdminConfig)
		if setAdminErr != nil {
			return nil, errors.Wrap(setAdminErr, "failed to set feed admin")
		}
	}

	metadatas := []data_feeds_cache.DataFeedsCacheWorkflowMetadata{}
	for idx := range input.AllowedWorkflowNames {
		metadatas = append(metadatas, data_feeds_cache.DataFeedsCacheWorkflowMetadata{
			AllowedWorkflowName:  df_changeset.HashedWorkflowName(input.AllowedWorkflowNames[idx]),
			AllowedSender:        input.AllowedSenders[idx],
			AllowedWorkflowOwner: input.AllowedWorkflowOwners[idx],
		})
	}

	feeIDs := []string{}
	for _, feedID := range input.FeedIDs {
		feeIDs = append(feeIDs, feedID[:32])
	}

	_, setFeedConfigErr := df_changeset.RunChangeset(df_changeset.SetFeedConfigChangeset, *input.CldEnv, df_changeset_types.SetFeedDecimalConfig{
		ChainSelector:    input.ChainSelector,
		CacheAddress:     input.DataFeedsCacheAddress,
		DataIDs:          feeIDs,
		Descriptions:     input.Descriptions,
		WorkflowMetadata: metadatas,
	})

	if setFeedConfigErr != nil {
		return nil, errors.Wrap(setFeedConfigErr, "failed to set feed config")
	}

	out := &types.ConfigureDataFeedsCacheOutput{
		DataFeedsCacheAddress: input.DataFeedsCacheAddress,
		FeedIDs:               input.FeedIDs,
		AllowedSenders:        input.AllowedSenders,
		AllowedWorkflowOwners: input.AllowedWorkflowOwners,
		AllowedWorkflowNames:  input.AllowedWorkflowNames,
	}

	if input.AdminAddress != (common.Address{}) {
		out.AdminAddress = input.AdminAddress
	}

	input.Out = out

	return out, nil
}
