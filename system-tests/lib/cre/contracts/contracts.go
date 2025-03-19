package contracts

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	"github.com/smartcontractkit/chainlink/deployment"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"

	workflow_registry_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	keystonenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
)

func ConfigureKeystone(input types.ConfigureKeystoneInput) error {
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
		if flags.HasFlag(donMetadata.Flags, types.CronCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "cron-trigger",
					Version:        "1.0.0",
					CapabilityType: 0, // TRIGGER
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		if flags.HasFlag(donMetadata.Flags, types.CustomComputeCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "custom-compute",
					Version:        "1.0.0",
					CapabilityType: 1, // ACTION
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		if flags.HasFlag(donMetadata.Flags, types.OCR3Capability) {
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

		if flags.HasFlag(donMetadata.Flags, types.WriteEVMCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "write_geth-testnet",
					Version:        "1.0.0",
					CapabilityType: 3, // TARGET
					ResponseType:   1, // OBSERVATION_IDENTICAL
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		// Add support for new capabilities here as needed

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

		donName := donMetadata.Name + "-don"
		donCapabilities = append(donCapabilities, keystone_changeset.DonCapabilities{
			Name:         donName,
			F:            1,
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

func DeployKeystone(testLogger zerolog.Logger, input *types.KeystoneContractsInput) (*types.KeystoneContractOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	if input.Out != nil && input.Out.UseCache {
		return input.Out, nil
	}

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	var err error
	forwarderAddress, err := DeployKeystoneForwarder(testLogger, input.CldEnv, input.ChainSelector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy Keystone Forwarder contract")
	}
	oCR3CapabilityAddress, err := DeployOCR3(testLogger, input.CldEnv, input.ChainSelector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy OCR3 contract")
	}
	capabilitiesRegistryAddress, err := DeployCapabilitiesRegistry(testLogger, input.CldEnv, input.ChainSelector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy Capabilities Registry contract")
	}
	workflowRegistryAddress, err := DeployWorkflowRegistry(testLogger, input.CldEnv, input.ChainSelector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy Workflow Registry contract")
	}

	out := &types.KeystoneContractOutput{
		ForwarderAddress:            forwarderAddress,
		OCR3CapabilityAddress:       oCR3CapabilityAddress,
		CapabilitiesRegistryAddress: capabilitiesRegistryAddress,
		WorkflowRegistryAddress:     workflowRegistryAddress,
	}

	input.Out = out
	return out, nil
}

func DeployOCR3(testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) (common.Address, error) {
	output, err := keystone_changeset.DeployOCR3(*ctfEnv, chainSelector)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to deploy OCR3 contract")
	}

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to merge address book")
	}

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	if err != nil {
		return common.Address{}, errors.Wrapf(err, "failed to get addresses for chain %d from the address book", chainSelector)
	}

	var ocr3capabilityAddr common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "OCR3Capability") {
			ocr3capabilityAddr = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed OCR3Capability contract at %s", ocr3capabilityAddr.Hex())
			break
		}
	}
	if ocr3capabilityAddr == (common.Address{}) {
		return common.Address{}, errors.New("failed to find OCR3Capability address in the address book")
	}

	return ocr3capabilityAddr, nil
}

func DeployCapabilitiesRegistry(testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) (common.Address, error) {
	output, err := keystone_changeset.DeployCapabilityRegistry(*ctfEnv, chainSelector)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to deploy Capabilities Registry contract")
	}

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to merge address book")
	}

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	if err != nil {
		return common.Address{}, errors.Wrapf(err, "failed to get addresses for chain %d from the address book", chainSelector)
	}

	var capabilitiesRegistryAddr common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "CapabilitiesRegistry") {
			capabilitiesRegistryAddr = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed Capabilities Registry contract at %s", capabilitiesRegistryAddr.Hex())
			break
		}
	}
	if capabilitiesRegistryAddr == (common.Address{}) {
		return common.Address{}, errors.New("failed to find Capabilities Registry address in the address book")
	}

	return capabilitiesRegistryAddr, nil
}

func DeployKeystoneForwarder(testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) (common.Address, error) {
	output, err := keystone_changeset.DeployForwarder(*ctfEnv, keystone_changeset.DeployForwarderRequest{
		ChainSelectors: []uint64{chainSelector},
	})
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to deploy forwarder contract")
	}

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to merge address book")
	}

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	if err != nil {
		return common.Address{}, errors.Wrapf(err, "failed to get addresses for chain %d from the address book", chainSelector)
	}

	var forwarderAddress common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "KeystoneForwarder") {
			forwarderAddress = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed KeystoneForwarder contract at %s", forwarderAddress.Hex())
			break
		}
	}
	if forwarderAddress == (common.Address{}) {
		return common.Address{}, errors.New("failed to find KeystoneForwarder address in the address book")
	}

	return forwarderAddress, nil
}

func DeployWorkflowRegistry(testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) (common.Address, error) {
	output, err := workflow_registry_changeset.Deploy(*ctfEnv, chainSelector)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to deploy workflow registry contract")
	}

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to merge address book")
	}

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	if err != nil {
		return common.Address{}, errors.Wrapf(err, "failed to get addresses for chain %d from the address book", chainSelector)
	}

	var workflowRegistryAddr common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "WorkflowRegistry") {
			workflowRegistryAddr = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed WorkflowRegistry contract at %s", workflowRegistryAddr.Hex())
		}
	}
	if workflowRegistryAddr == (common.Address{}) {
		return common.Address{}, errors.New("failed to find WorkflowRegistry address in the address book")
	}

	return workflowRegistryAddr, nil
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

func DeployFeedsConsumer(testLogger zerolog.Logger, input *types.DeployFeedConsumerInput) (*types.DeployFeedConsumerOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	if input.Out != nil && input.Out.UseCache {
		return input.Out, nil
	}

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	output, err := keystone_changeset.DeployFeedsConsumer(*input.CldEnv, &keystone_changeset.DeployFeedsConsumerRequest{
		ChainSelector: input.ChainSelector,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy feeds_consumer contract")
	}

	err = input.CldEnv.ExistingAddresses.Merge(output.AddressBook)
	if err != nil {
		return nil, errors.Wrap(err, "failed to merge address book")
	}

	addresses, err := input.CldEnv.ExistingAddresses.AddressesForChain(input.ChainSelector)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get addresses for chain %d from the address book", input.ChainSelector)
	}

	var feedsConsumerAddress common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "FeedConsumer") {
			feedsConsumerAddress = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed FeedConsumer contract at %s", feedsConsumerAddress.Hex())
			break
		}
	}

	if feedsConsumerAddress == (common.Address{}) {
		return nil, errors.New("failed to find FeedConsumer address in the address book")
	}

	out := &types.DeployFeedConsumerOutput{
		FeedConsumerAddress: feedsConsumerAddress,
	}

	input.Out = out
	return out, nil
}

func ConfigureFeedsConsumer(testLogger zerolog.Logger, input *types.ConfigureFeedConsumerInput) (*types.ConfigureFeedConsumerOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	if input.Out != nil && input.Out.UseCache {
		return input.Out, nil
	}

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(input.FeedConsumerAddress, input.SethClient.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create feeds consumer instance")
	}

	// Prepare hex-encoded and truncated workflow name

	var HashTruncateName = func(name string) string {
		// Compute SHA-256 hash of the input string
		hash := sha256.Sum256([]byte(name))

		// Encode as hex to ensure UTF8
		var hashBytes []byte = hash[:]
		resultHex := hex.EncodeToString(hashBytes)

		// Truncate to 10 bytes
		truncated := []byte(resultHex)[:10]
		return string(truncated)
	}

	truncatedNames := make([][10]byte, 0, len(input.AllowedWorkflowNames))
	for _, workflowName := range input.AllowedWorkflowNames {
		var workflowNameBytes [10]byte
		truncated := HashTruncateName(workflowName)
		copy(workflowNameBytes[:], []byte(truncated))

		truncatedNames = append(truncatedNames, workflowNameBytes)
	}

	_, decodeErr := input.SethClient.Decode(feedsConsumerInstance.SetConfig(
		input.SethClient.NewTXOpts(),
		input.AllowedSenders,        // forwarder contract!!!
		input.AllowedWorkflowOwners, // allowed workflow owners
		// here we need to use hex-encoded workflow name converted to []byte
		truncatedNames, // allowed workflow names
	))
	if decodeErr != nil {
		return nil, errors.Wrap(decodeErr, "failed to set config for feeds consumer")
	}

	out := &types.ConfigureFeedConsumerOutput{
		FeedConsumerAddress:   input.FeedConsumerAddress,
		AllowedSenders:        input.AllowedSenders,
		AllowedWorkflowOwners: input.AllowedWorkflowOwners,
		AllowedWorkflowNames:  input.AllowedWorkflowNames,
	}

	input.Out = out
	return out, nil
}
