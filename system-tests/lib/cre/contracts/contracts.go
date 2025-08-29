package contracts

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	ks_solana "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"

	cre_contracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	df_changeset_types "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"

	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
)

func ConfigureKeystone(input cre.ConfigureKeystoneInput, capabilityRegistryConfigFns []cre.CapabilityRegistryConfigFn) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "input validation failed")
	}

	donCapabilities := make([]keystone_changeset.DonCapabilities, 0, len(input.Topology.DonsMetadata))

	for donIdx, donMetadata := range input.Topology.DonsMetadata {
		// if it's only a gateway DON, we don't want to register it with the Capabilities Registry
		// since it doesn't have any capabilities
		if flags.HasOnlyOneFlag(donMetadata.Flags, cre.GatewayDON) {
			continue
		}

		var capabilities []keystone_changeset.DONCapabilityWithConfig

		// check what capabilities each DON has and register them with Capabilities Registry contract
		for _, configFn := range capabilityRegistryConfigFns {
			if configFn == nil {
				continue
			}

			capabilitiesFn, configFnErr := configFn(donMetadata.Flags, input.NodeSets[donIdx])
			if configFnErr != nil {
				return errors.Wrap(configFnErr, "failed to get capabilities from config function")
			}

			capabilities = append(capabilities, capabilitiesFn...)
		}

		workerNodes, workerNodesErr := crenode.FindManyWithLabel(donMetadata.NodesMetadata, &cre.Label{
			Key:   crenode.NodeTypeKey,
			Value: cre.WorkerNode,
		}, crenode.EqualLabels)

		if workerNodesErr != nil {
			return errors.Wrap(workerNodesErr, "failed to find worker nodes")
		}

		donPeerIDs := make([]string, len(workerNodes))
		for i, node := range workerNodes {
			p2pID, err := crenode.ToP2PID(node, crenode.NoOpTransformFn)
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
			if flags.HasFlag(donMetadata.Flags, cre.ConsensusCapability) || flags.HasFlag(donMetadata.Flags, cre.ConsensusCapabilityV2) {
				return fmt.Errorf("incorrect number of worker nodes: %d. Resulting F must conform to formula: mod((N-1)/3) > 0", len(workerNodes))
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
		if flags.HasFlag(metaDon.Flags, cre.ConsensusCapability) || flags.HasFlag(metaDon.Flags, cre.ConsensusCapabilityV2) {
			workerNodes, workerNodesErr := crenode.FindManyWithLabel(metaDon.NodesMetadata, &cre.Label{
				Key:   crenode.NodeTypeKey,
				Value: cre.WorkerNode,
			}, crenode.EqualLabels)

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

	_, err := operations.ExecuteSequence(
		input.CldEnv.OperationsBundle,
		ks_contracts_op.ConfigureCapabilitiesRegistrySeq,
		ks_contracts_op.ConfigureCapabilitiesRegistrySeqDeps{
			Env:  input.CldEnv,
			Dons: donCapabilities,
		},
		ks_contracts_op.ConfigureCapabilitiesRegistrySeqInput{
			RegistryChainSel: input.ChainSelector,
			UseMCMS:          false,
			ContractAddress:  input.CapabilitiesRegistryAddress,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to configure capabilities registry")
	}

	capReg, err := cre_contracts.GetOwnedContractV2[*kcr.CapabilitiesRegistry](
		input.CldEnv.DataStore.Addresses(),
		input.CldEnv.BlockChains.EVMChains()[input.ChainSelector],
		input.CapabilitiesRegistryAddress.Hex(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to get capabilities registry contract")
	}

	configDONs := make([]ks_contracts_op.ConfigureKeystoneDON, 0)
	for _, donCap := range donCapabilities {
		don := ks_contracts_op.ConfigureKeystoneDON{
			Name: donCap.Name,
		}
		for _, nop := range donCap.Nops {
			don.NodeIDs = append(don.NodeIDs, nop.Nodes...)
		}
		configDONs = append(configDONs, don)
	}

	// remove chains that do not require any configurations ('read-only' chains that do not have forwarders deployed)
	allAddresses, addrErr := input.CldEnv.ExistingAddresses.Addresses() //nolint:staticcheck // ignore SA1019 as ExistingAddresses is deprecated but still used
	if addrErr != nil {
		return errors.Wrap(addrErr, "failed to get addresses from address book")
	}

	evmChainsWithForwarders := make(map[uint64]struct{})
	for chainSelector, addresses := range allAddresses {
		for _, typeAndVersion := range addresses {
			if typeAndVersion.Type == keystone_changeset.KeystoneForwarder {
				evmChainsWithForwarders[chainSelector] = struct{}{}
			}
		}
	}

	solChainsWithForwarder := make(map[uint64]struct{})
	solForwarders := input.CldEnv.DataStore.Addresses().Filter(datastore.AddressRefByQualifier(ks_solana.DefaultForwarderQualifier))
	for _, forwarder := range solForwarders {
		solChainsWithForwarder[forwarder.ChainSelector] = struct{}{}
	}

	// configure Solana forwarder only if we have some
	if len(solChainsWithForwarder) > 0 {
		for _, don := range configDONs {
			cs := commonchangeset.Configure(ks_solana.ConfigureForwarders{},
				&ks_solana.ConfigureForwarderRequest{
					WFDonName:        don.Name,
					WFNodeIDs:        don.NodeIDs,
					RegistryChainSel: input.ChainSelector,
					Chains:           solChainsWithForwarder,
					Qualifier:        ks_solana.DefaultForwarderQualifier,
					Version:          "1.0.0",
				},
			)

			_, err = cs.Apply(*input.CldEnv)
			if err != nil {
				return errors.Wrap(err, "failed to configure Solana forwarders")
			}
		}
	}

	// configure EVM forwarders only if we have some
	if len(evmChainsWithForwarders) > 0 {
		_, err = operations.ExecuteSequence(
			input.CldEnv.OperationsBundle,
			ks_contracts_op.ConfigureForwardersSeq,
			ks_contracts_op.ConfigureForwardersSeqDeps{
				Env:      input.CldEnv,
				Registry: capReg.Contract,
			},
			ks_contracts_op.ConfigureForwardersSeqInput{
				RegistryChainSel: input.ChainSelector,
				DONs:             configDONs,
				Chains:           evmChainsWithForwarders,
			},
		)
		if err != nil {
			return errors.Wrap(err, "failed to configure forwarders")
		}
	}

	_, err = operations.ExecuteOperation(
		input.CldEnv.OperationsBundle,
		ks_contracts_op.ConfigureOCR3Op,
		ks_contracts_op.ConfigureOCR3OpDeps{
			Env:      input.CldEnv,
			Registry: capReg.Contract,
		},
		ks_contracts_op.ConfigureOCR3OpInput{
			ContractAddress:  input.OCR3Address,
			RegistryChainSel: input.ChainSelector,
			DONs:             configDONs,
			Config:           &input.OCR3Config,
			DryRun:           false,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to configure OCR3 contract")
	}

	_, err = operations.ExecuteOperation(
		input.CldEnv.OperationsBundle,
		ks_contracts_op.ConfigureOCR3Op,
		ks_contracts_op.ConfigureOCR3OpDeps{
			Env:      input.CldEnv,
			Registry: capReg.Contract,
		},
		ks_contracts_op.ConfigureOCR3OpInput{
			ContractAddress:  input.DONTimeAddress,
			RegistryChainSel: input.ChainSelector,
			DONs:             configDONs,
			Config:           &input.DONTimeConfig,
			DryRun:           false,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to configure DON Time contract")
	}

	if input.VaultOCR3Address.Cmp(common.Address{}) != 0 {
		_, err = operations.ExecuteOperation(
			input.CldEnv.OperationsBundle,
			ks_contracts_op.ConfigureOCR3Op,
			ks_contracts_op.ConfigureOCR3OpDeps{
				Env:      input.CldEnv,
				Registry: capReg.Contract,
			},
			ks_contracts_op.ConfigureOCR3OpInput{
				ContractAddress:  input.VaultOCR3Address,
				RegistryChainSel: input.ChainSelector,
				DONs:             configDONs,
				Config:           &input.VaultOCR3Config,
				DryRun:           false,
			},
		)
		if err != nil {
			return errors.Wrap(err, "failed to configure Vault OCR3 contract")
		}
	}

	for chainSelector, evmOCR3Address := range *input.EVMOCR3Addresses {
		if evmOCR3Address.Cmp(common.Address{}) != 0 {
			_, err = operations.ExecuteOperation(
				input.CldEnv.OperationsBundle,
				ks_contracts_op.ConfigureOCR3Op,
				ks_contracts_op.ConfigureOCR3OpDeps{
					Env:      input.CldEnv,
					Registry: capReg.Contract,
				},
				ks_contracts_op.ConfigureOCR3OpInput{
					ContractAddress:  &evmOCR3Address,
					RegistryChainSel: chainSelector,
					DONs:             configDONs,
					Config:           &input.EVMOCR3Config,
					DryRun:           false,
				},
			)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to configure EVM OCR3 contract for chain selector: %d, address:%s", chainSelector, evmOCR3Address.Hex()))
			}
		}
	}

	if input.ConsensusV2OCR3Address.Cmp(common.Address{}) != 0 {
		_, err = operations.ExecuteOperation(
			input.CldEnv.OperationsBundle,
			ks_contracts_op.ConfigureOCR3Op,
			ks_contracts_op.ConfigureOCR3OpDeps{
				Env:      input.CldEnv,
				Registry: capReg.Contract,
			},
			ks_contracts_op.ConfigureOCR3OpInput{
				ContractAddress:  input.ConsensusV2OCR3Address,
				RegistryChainSel: input.ChainSelector,
				DONs:             configDONs,
				Config:           &input.ConsensusV2OCR3Config,
				DryRun:           false,
			},
		)
		if err != nil {
			return errors.Wrap(err, "failed to configure Consensus OCR3 contract")
		}
	}
	return nil
}

// values supplied by Alexandr Yepishev as the expected values for OCR3 config
func DefaultOCR3Config(topology *cre.Topology) (*keystone_changeset.OracleConfig, error) {
	var transmissionSchedule []int

	for _, metaDon := range topology.DonsMetadata {
		if flags.HasFlag(metaDon.Flags, cre.ConsensusCapability) || flags.HasFlag(metaDon.Flags, cre.ConsensusCapabilityV2) {
			workerNodes, workerNodesErr := crenode.FindManyWithLabel(metaDon.NodesMetadata, &cre.Label{
				Key:   crenode.NodeTypeKey,
				Value: cre.WorkerNode,
			}, crenode.EqualLabels)

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

	// values supplied by Alexandr Yepishev as the expected values for OCR3 config
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

func FindAddressesForChain(addressBook cldf.AddressBook, chainSelector uint64, contractName string) (common.Address, error) {
	addresses, err := addressBook.AddressesForChain(chainSelector)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to get addresses for chain")
	}

	for addrStr, tv := range addresses {
		if !strings.Contains(tv.String(), contractName) {
			continue
		}

		return common.HexToAddress(addrStr), nil
	}

	return common.Address{}, fmt.Errorf("failed to find %s address in the address book for chain %d", contractName, chainSelector)
}

func MustFindAddressesForChain(addressBook cldf.AddressBook, chainSelector uint64, contractName string) common.Address {
	addr, err := FindAddressesForChain(addressBook, chainSelector, contractName)
	if err != nil {
		panic(fmt.Errorf("failed to find %s address in the address book for chain %d", contractName, chainSelector))
	}
	return addr
}

// MergeAllDataStores merges all DataStores (after contracts deployments)
func MergeAllDataStores(fullCldEnvOutput *cre.FullCLDEnvironmentOutput, changesetOutputs ...cldf.ChangesetOutput) {
	fmt.Print("Merging DataStores (after contracts deployments)...")
	minChangesetsCap := 2
	if len(changesetOutputs) < minChangesetsCap {
		panic(fmt.Errorf("DataStores merging failed: at least %d changesets required", minChangesetsCap))
	}

	// Start with the first changeset's data store
	baseDataStore := changesetOutputs[0].DataStore

	// Merge all subsequent changesets into the base data store
	for i := 1; i < len(changesetOutputs); i++ {
		otherDataStore := changesetOutputs[i].DataStore
		mergeErr := baseDataStore.Merge(otherDataStore.Seal())
		if mergeErr != nil {
			panic(errors.Wrap(mergeErr, "DataStores merging failed"))
		}
	}

	fullCldEnvOutput.Environment.DataStore = baseDataStore.Seal()
}

func ConfigureWorkflowRegistry(testLogger zerolog.Logger, input *cre.WorkflowRegistryInput) (*cre.WorkflowRegistryOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	if input.Out != nil && input.Out.UseCache {
		return input.Out, nil
	}

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	allowedDonIDs := make([]uint32, len(input.AllowedDonIDs))
	for i, donID := range input.AllowedDonIDs {
		allowedDonIDs[i] = libc.MustSafeUint32FromUint64(donID)
	}

	report, err := operations.ExecuteSequence(
		input.CldEnv.OperationsBundle,
		ks_contracts_op.ConfigWorkflowRegistrySeq,
		ks_contracts_op.ConfigWorkflowRegistrySeqDeps{
			Env: input.CldEnv,
		},
		ks_contracts_op.ConfigWorkflowRegistrySeqInput{
			ContractAddress:       input.ContractAddress,
			RegistryChainSelector: input.ChainSelector,
			AllowedDonIDs:         allowedDonIDs,
			WorkflowOwners:        input.WorkflowOwners,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to configure workflow registry")
	}

	input.Out = &cre.WorkflowRegistryOutput{
		ChainSelector:  report.Output.RegistryChainSelector,
		AllowedDonIDs:  report.Output.AllowedDonIDs,
		WorkflowOwners: report.Output.WorkflowOwners,
	}
	return input.Out, nil
}

func ConfigureDataFeedsCache(testLogger zerolog.Logger, input *cre.ConfigureDataFeedsCacheInput) (*cre.ConfigureDataFeedsCacheOutput, error) {
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
		_, setAdminErr := commonchangeset.RunChangeset(df_changeset.SetFeedAdminChangeset, *input.CldEnv, setAdminConfig)
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

	_, setFeedConfigErr := commonchangeset.RunChangeset(df_changeset.SetFeedConfigChangeset, *input.CldEnv, df_changeset_types.SetFeedDecimalConfig{
		ChainSelector:    input.ChainSelector,
		CacheAddress:     input.DataFeedsCacheAddress,
		DataIDs:          feeIDs,
		Descriptions:     input.Descriptions,
		WorkflowMetadata: metadatas,
	})

	if setFeedConfigErr != nil {
		return nil, errors.Wrap(setFeedConfigErr, "failed to set feed config")
	}

	out := &cre.ConfigureDataFeedsCacheOutput{
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

func DeployDataFeedsCacheContract(testLogger zerolog.Logger, chainSelector uint64, fullCldEnvOutput *cre.FullCLDEnvironmentOutput) (common.Address, cldf.ChangesetOutput, error) {
	testLogger.Info().Msg("Deploying Data Feeds Cache contract...")
	deployDfConfig := df_changeset_types.DeployConfig{
		ChainsToDeploy: []uint64{chainSelector},
		Labels:         []string{"data-feeds"}, // label required by the changeset
	}

	dfOutput, dfErr := commonchangeset.RunChangeset(df_changeset.DeployCacheChangeset, *fullCldEnvOutput.Environment, deployDfConfig)
	if dfErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrapf(dfErr, "failed to deploy Data Feeds Cache contract on chain %d", chainSelector)
	}

	mergeErr := fullCldEnvOutput.Environment.ExistingAddresses.Merge(dfOutput.AddressBook) //nolint:staticcheck // won't migrate now
	if mergeErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrap(mergeErr, "failed to merge address book of Data Feeds Cache contract")
	}
	testLogger.Info().Msgf("Data Feeds Cache contract deployed to %d", chainSelector)

	dataFeedsCacheAddress, dataFeedsCacheErr := FindAddressesForChain(
		fullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		chainSelector,
		df_changeset.DataFeedsCache.String(),
	)
	if dataFeedsCacheErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrapf(dataFeedsCacheErr, "failed to find Data Feeds Cache contract address on chain %d", chainSelector)
	}
	testLogger.Info().Msgf("Data Feeds Cache contract found on chain %d at address %s", chainSelector, dataFeedsCacheAddress)

	return dataFeedsCacheAddress, dfOutput, nil
}

func DeployReadBalancesContract(testLogger zerolog.Logger, chainSelector uint64, fullCldEnvOutput *cre.FullCLDEnvironmentOutput) (common.Address, cldf.ChangesetOutput, error) {
	testLogger.Info().Msg("Deploying Read Balances contract...")
	deployReadBalanceRequest := &keystone_changeset.DeployRequestV2{ChainSel: chainSelector}
	rbOutput, rbErr := keystone_changeset.DeployBalanceReaderV2(*fullCldEnvOutput.Environment, deployReadBalanceRequest)
	if rbErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrap(rbErr, "failed to deploy Read Balances contract")
	}

	mergeErr2 := fullCldEnvOutput.Environment.ExistingAddresses.Merge(rbOutput.AddressBook) //nolint:staticcheck // won't migrate now
	if mergeErr2 != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrap(mergeErr2, "failed to merge address book of Read Balances contract")
	}
	testLogger.Info().Msgf("Read Balances contract deployed to %d", chainSelector)

	readBalancesAddress, readContractErr := FindAddressesForChain(
		fullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		chainSelector,
		keystone_changeset.BalanceReader.String(),
	)
	if readContractErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrap(readContractErr, "failed to find Read Balances contract address")
	}
	testLogger.Info().Msgf("Read Balances contract found on chain %d at address %s", chainSelector, readBalancesAddress)

	return readBalancesAddress, rbOutput, nil
}
