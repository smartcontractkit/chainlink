package contracts

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	cre_contracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	df_changeset_types "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	cap_reg_v2_seq "github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
	syncer_v2 "github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer/v2"
)

const DonFamily = "test-don-family"

type donConfig struct {
	id uint32 // the DON id as registered in the capabilities registry
	keystone_changeset.DonCapabilities
	flags []cre.CapabilityFlag
}

type dons struct {
	c        map[string]donConfig
	offChain offchain.Client
}

func (d *dons) donsOrderedByID() []donConfig {
	out := make([]donConfig, 0, len(d.c))
	for _, don := range d.c {
		out = append(out, don)
	}

	// Use sort library to sort by ID
	sort.Slice(out, func(i, j int) bool {
		return out[i].id < out[j].id
	})

	return out
}

func (d *dons) allDonCapabilities() []keystone_changeset.DonCapabilities {
	out := make([]keystone_changeset.DonCapabilities, 0, len(d.c))
	for _, don := range d.donsOrderedByID() {
		out = append(out, don.DonCapabilities)
	}
	return out
}

func (d *dons) mustToV2ConfigureInput(chainSelector uint64, contractAddress string) cap_reg_v2_seq.ConfigureCapabilitiesRegistryInput {
	nops := make([]capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams, 0)
	nodes := make([]capabilities_registry_v2.CapabilitiesRegistryNodeParams, 0)
	capabilities := make([]capabilities_registry_v2.CapabilitiesRegistryCapability, 0)
	donParams := make([]capabilities_registry_v2.CapabilitiesRegistryNewDONParams, 0)

	// Collect unique capabilities and NOPs
	capabilityMap := make(map[string]capabilities_registry_v2.CapabilitiesRegistryCapability)
	nopMap := make(map[string]capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams)
	for _, don := range d.donsOrderedByID() {
		// Extract capabilities
		capIDs := make([]string, 0, len(don.Capabilities))
		for _, myCap := range don.Capabilities {
			capID := fmt.Sprintf("%s@%s", myCap.Capability.LabelledName, myCap.Capability.Version)
			capIDs = append(capIDs, capID)
			if _, exists := capabilityMap[capID]; !exists {
				metadataJSON, _ := json.Marshal(syncer_v2.CapabilityMetadata{
					CapabilityType: myCap.Capability.CapabilityType,
					ResponseType:   myCap.Capability.ResponseType,
				})
				capabilityMap[capID] = capabilities_registry_v2.CapabilitiesRegistryCapability{
					CapabilityId:          capID,
					ConfigurationContract: common.Address{},
					Metadata:              metadataJSON,
				}
			}
		}

		// Extract NOPs and nodes
		adminAddrs, err := generateAdminAddresses(len(don.Nops))
		if err != nil {
			panic(fmt.Sprintf("failed to generate admin addresses: %s", err))
		}
		for i, nop := range don.Nops {
			nopName := nop.Name
			if _, exists := nopMap[nopName]; !exists {
				nopMap[nopName] = capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams{
					Admin: adminAddrs[i],
					Name:  nopName,
				}

				ns, err := deployment.NodeInfo(nop.Nodes, d.offChain)
				if err != nil {
					panic(err)
				}

				// Add nodes for this NOP
				for _, n := range ns {
					ocrCfg, ok := n.OCRConfigForChainSelector(chainSelector)
					if !ok {
						continue
					}

					wfKey, err := hex.DecodeString(n.WorkflowKey)
					if err != nil {
						panic(err)
					}

					csKey, err := hex.DecodeString(n.CSAKey)
					if err != nil {
						panic(fmt.Errorf("failed to decode csa key: %w", err))
					}

					nodes = append(nodes, capabilities_registry_v2.CapabilitiesRegistryNodeParams{
						NodeOperatorId:      libc.MustSafeUint32(i + 1),
						P2pId:               n.PeerID,
						Signer:              ocrCfg.OffchainPublicKey,
						EncryptionPublicKey: [32]byte(wfKey),
						CsaKey:              [32]byte(csKey),
						CapabilityIds:       capIDs,
					})
				}
			}
		}

		// Create DON parameters
		var capConfigs []capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration
		for _, cap := range don.Capabilities {
			capID := fmt.Sprintf("%s@%s", cap.Capability.LabelledName, cap.Capability.Version)
			configBytes := []byte("{}")
			if cap.Config != nil {
				// Convert proto config to bytes if needed
				if protoBytes, err := proto.Marshal(cap.Config); err == nil {
					configBytes = protoBytes
				}
			}
			capConfigs = append(capConfigs, capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration{
				CapabilityId: capID,
				Config:       configBytes,
			})
		}

		var donNodes [][32]byte
		for _, nop := range don.Nops {
			for _, nodeID := range nop.Nodes {
				peerID, err := p2pkey.MakePeerID(nodeID)
				if err != nil {
					continue
				}
				donNodes = append(donNodes, peerID)
			}
		}

		donParams = append(donParams, capabilities_registry_v2.CapabilitiesRegistryNewDONParams{
			Name:                     don.Name,
			DonFamilies:              []string{DonFamily}, // Default empty
			Config:                   []byte("{}"),
			CapabilityConfigurations: capConfigs,
			Nodes:                    donNodes,
			F:                        don.F,
			IsPublic:                 true,
			AcceptsWorkflows:         true,
		})
	}

	// Convert maps to slices
	for _, cap := range capabilityMap {
		capabilities = append(capabilities, cap)
	}
	for _, nop := range nopMap {
		nops = append(nops, nop)
	}

	return cap_reg_v2_seq.ConfigureCapabilitiesRegistryInput{
		RegistryChainSel: chainSelector,
		ContractAddress:  contractAddress,
		Nops:             nops,
		Nodes:            nodes,
		Capabilities:     capabilities,
		DONs:             donParams,
	}
}

func generateAdminAddresses(count int) ([]common.Address, error) {
	if count <= 0 {
		return nil, errors.New("count must be a positive integer")
	}

	// Determine the number of hex digits needed for padding based on the count.
	// We use the count + 1 to account for the loop range and a safe margin.
	hexDigits := int(math.Ceil(math.Log10(float64(count+1)) / math.Log10(16)))
	if hexDigits < 1 {
		hexDigits = 1
	}

	// The total length of the address after the "0x" prefix must be 40.
	baseHexLen := 40 - hexDigits
	if baseHexLen <= 0 {
		return nil, errors.New("count is too large to generate unique addresses with this base")
	}

	// Create a base string of 'f' characters to ensure the addresses are not zero.
	baseString := strings.Repeat("f", baseHexLen)

	addresses := make([]common.Address, count)
	for i := 0; i < count; i++ {
		format := fmt.Sprintf("%s%%0%dx", baseString, hexDigits)
		fullAddress := fmt.Sprintf(format, i)
		addresses[i] = common.HexToAddress("0x" + fullAddress)
	}

	return addresses, nil
}

func toDons(input cre.ConfigureKeystoneInput) (*dons, error) {
	dons := &dons{
		c:        make(map[string]donConfig),
		offChain: input.CldEnv.Offchain,
	}

	for donIdx, donMetadata := range input.Topology.DonsMetadata.List() {
		// if it's only a gateway DON, we don't want to register it with the Capabilities Registry
		// since it doesn't have any capabilities
		if flags.HasOnlyOneFlag(donMetadata.Flags, cre.GatewayDON) {
			continue
		}

		var capabilities []keystone_changeset.DONCapabilityWithConfig

		// check what capabilities each DON has and register them with Capabilities Registry contract
		for _, configFn := range input.CapabilityRegistryConfigFns {
			if configFn == nil {
				continue
			}

			enabledCapabilities, err2 := configFn(donMetadata.Flags, input.NodeSets[donIdx])
			if err2 != nil {
				return nil, errors.Wrap(err2, "failed to get capabilities from config function")
			}

			capabilities = append(capabilities, enabledCapabilities...)
		}

		// add capabilities that were passed directly via the input (from the PostDONStartup of features)
		if input.DONCapabilityWithConfigs != nil && input.DONCapabilityWithConfigs[donMetadata.ID] != nil {
			capabilities = append(capabilities, input.DONCapabilityWithConfigs[donMetadata.ID]...)
		}

		workerNodes, wErr := donMetadata.Workers()
		if wErr != nil {
			return nil, errors.Wrap(wErr, "failed to find worker nodes")
		}

		donPeerIDs := make([]string, len(workerNodes))
		for i, node := range workerNodes {
			// we need to use p2pID here with the "p2p_" prefix
			donPeerIDs[i] = node.Keys.P2PKey.PeerID.String()
		}

		forwarderF := (len(workerNodes) - 1) / 3
		if forwarderF == 0 {
			if flags.HasFlag(donMetadata.Flags, cre.ConsensusCapability) || flags.HasFlag(donMetadata.Flags, cre.ConsensusCapabilityV2) {
				return nil, fmt.Errorf("incorrect number of worker nodes: %d. Resulting F must conform to formula: mod((N-1)/3) > 0", len(workerNodes))
			}
			// for other capabilities, we can use 1 as F
			forwarderF = 1
		}

		// we only need to assign P2P IDs to NOPs, since `ConfigureInitialContractsChangeset` method
		// will take care of creating DON to Nodes mapping
		nop := keystone_changeset.NOP{
			Name:  fmt.Sprintf("NOP for %s DON", donMetadata.Name),
			Nodes: donPeerIDs,
		}
		donName := donMetadata.Name + "-don"
		c := keystone_changeset.DonCapabilities{
			Name:         donName,
			F:            libc.MustSafeUint8(forwarderF),
			Nops:         []keystone_changeset.NOP{nop},
			Capabilities: capabilities,
		}

		dons.c[donName] = donConfig{
			id:              uint32(donMetadata.ID), //nolint:gosec // G115
			DonCapabilities: c,
			flags:           donMetadata.Flags,
		}
	}

	return dons, nil
}

func ConfigureCapabilityRegistry(input cre.ConfigureKeystoneInput, dons *dons) (CapabilitiesRegistry, error) {
	if !input.WithV2Registries {
		_, err := operations.ExecuteSequence(
			input.CldEnv.OperationsBundle,
			ks_contracts_op.ConfigureCapabilitiesRegistrySeq,
			ks_contracts_op.ConfigureCapabilitiesRegistrySeqDeps{
				Env:  input.CldEnv,
				Dons: dons.allDonCapabilities(),
			},
			ks_contracts_op.ConfigureCapabilitiesRegistrySeqInput{
				RegistryChainSel: input.ChainSelector,
				UseMCMS:          false,
				ContractAddress:  input.CapabilitiesRegistryAddress,
			},
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to configure capabilities registry")
		}

		capReg, err := cre_contracts.GetOwnedContractV2[*kcr.CapabilitiesRegistry](
			input.CldEnv.DataStore.Addresses(),
			input.CldEnv.BlockChains.EVMChains()[input.ChainSelector],
			input.CapabilitiesRegistryAddress.Hex(),
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get capabilities registry contract")
		}
		return &registryWrapper{V1: capReg.Contract}, nil
	}

	// Transform dons data to V2 sequence input format
	v2Input := dons.mustToV2ConfigureInput(input.ChainSelector, input.CapabilitiesRegistryAddress.Hex())
	_, err := operations.ExecuteSequence(
		input.CldEnv.OperationsBundle,
		cap_reg_v2_seq.ConfigureCapabilitiesRegistry,
		cap_reg_v2_seq.ConfigureCapabilitiesRegistryDeps{
			Env: input.CldEnv,
		},
		v2Input,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to configure capabilities registry")
	}

	capReg, err := cre_contracts.GetOwnedContractV2[*capabilities_registry_v2.CapabilitiesRegistry](
		input.CldEnv.DataStore.Addresses(),
		input.CldEnv.BlockChains.EVMChains()[input.ChainSelector],
		input.CapabilitiesRegistryAddress.Hex(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get capabilities registry contract")
	}

	return &registryWrapper{V2: capReg.Contract}, nil
}

func ConfigureKeystone(input cre.ConfigureKeystoneInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "input validation failed")
	}

	dons, err := toDons(input)
	if err != nil {
		return errors.Wrap(err, "failed to map input to dons")
	}

	_, err = ConfigureCapabilityRegistry(input, dons)
	if err != nil {
		return errors.Wrap(err, "failed to configure capability registry")
	}

	return nil
}

// values supplied by Alexandr Yepishev as the expected values for OCR3 config
func DefaultOCR3Config() (*keystone_changeset.OracleConfig, error) {
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
		MaxDurationQueryMillis:            1000,
		MaxDurationObservationMillis:      1000,
		MaxDurationShouldAcceptMillis:     1000,
		MaxDurationShouldTransmitMillis:   1000,
		MaxFaultyOracles:                  1,
		ConsensusCapOffchainConfig: &ocr3.ConsensusCapOffchainConfig{
			MaxQueryLengthBytes:       1000000,
			MaxObservationLengthBytes: 1000000,
			MaxOutcomeLengthBytes:     1000000,
			MaxReportLengthBytes:      1000000,
			MaxBatchSize:              1000,
		},
		UniqueReports: true,
	}

	return oracleConfig, nil
}

func DefaultChainCapabilityOCR3Config() (*keystone_changeset.OracleConfig, error) {
	cfg, err := DefaultOCR3Config()
	if err != nil {
		return nil, fmt.Errorf("failed to generate default OCR3 config: %w", err)
	}

	cfg.DeltaRoundMillis = 1000
	const kib = 1024
	const mib = 1024 * kib
	cfg.ConsensusCapOffchainConfig = nil
	cfg.ChainCapOffchainConfig = &ocr3.ChainCapOffchainConfig{
		MaxQueryLengthBytes:       mib,
		MaxObservationLengthBytes: 97 * kib,
		MaxReportLengthBytes:      mib,
		MaxOutcomeLengthBytes:     mib,
		MaxReportCount:            1000,
		MaxBatchSize:              200,
	}
	return cfg, nil
}

func FindAddressesForChain(addressBook cldf.AddressBook, chainSelector uint64, contractName string) (common.Address, cldf.TypeAndVersion, error) {
	addresses, err := addressBook.AddressesForChain(chainSelector)
	if err != nil {
		return common.Address{}, cldf.TypeAndVersion{}, errors.Wrap(err, "failed to get addresses for chain")
	}

	for addrStr, tv := range addresses {
		if !strings.Contains(tv.String(), contractName) {
			continue
		}

		return common.HexToAddress(addrStr), tv, nil
	}

	return common.Address{}, cldf.TypeAndVersion{}, fmt.Errorf("failed to find %s address in the address book for chain %d", contractName, chainSelector)
}

// TODO: CRE-742 use datastore
func MustFindAddressesForChain(addressBook cldf.AddressBook, chainSelector uint64, contractName string) common.Address {
	addr, _, err := FindAddressesForChain(addressBook, chainSelector, contractName)
	if err != nil {
		panic(fmt.Errorf("failed to find %s address in the address book for chain %d", contractName, chainSelector))
	}
	return addr
}

// MergeAllDataStores merges all DataStores (after contracts deployments)
func MergeAllDataStores(creEnvironment *cre.Environment, changesetOutputs ...cldf.ChangesetOutput) {
	framework.L.Info().Msg("Merging DataStores (after contracts deployments)...")
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

	creEnvironment.CldfEnvironment.DataStore = baseDataStore.Seal()
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

func DeployDataFeedsCacheContract(testLogger zerolog.Logger, chainSelector uint64, creEnvironment *cre.Environment) (common.Address, cldf.ChangesetOutput, error) {
	testLogger.Info().Msg("Deploying Data Feeds Cache contract...")
	deployDfConfig := df_changeset_types.DeployConfig{
		ChainsToDeploy: []uint64{chainSelector},
		Labels:         []string{"data-feeds"}, // label required by the changeset
	}

	dfOutput, dfErr := commonchangeset.RunChangeset(df_changeset.DeployCacheChangeset, *creEnvironment.CldfEnvironment, deployDfConfig)
	if dfErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrapf(dfErr, "failed to deploy Data Feeds Cache contract on chain %d", chainSelector)
	}

	mergeErr := creEnvironment.CldfEnvironment.ExistingAddresses.Merge(dfOutput.AddressBook) //nolint:staticcheck // won't migrate now
	if mergeErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrap(mergeErr, "failed to merge address book of Data Feeds Cache contract")
	}
	testLogger.Info().Msgf("Data Feeds Cache contract deployed to %d", chainSelector)

	dataFeedsCacheAddress, _, dataFeedsCacheErr := FindAddressesForChain(
		creEnvironment.CldfEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		chainSelector,
		df_changeset.DataFeedsCache.String(),
	)
	if dataFeedsCacheErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrapf(dataFeedsCacheErr, "failed to find Data Feeds Cache contract address on chain %d", chainSelector)
	}
	testLogger.Info().Msgf("Data Feeds Cache contract found on chain %d at address %s", chainSelector, dataFeedsCacheAddress)

	return dataFeedsCacheAddress, dfOutput, nil
}

func DeployReadBalancesContract(testLogger zerolog.Logger, chainSelector uint64, creEnvironment *cre.Environment) (common.Address, cldf.ChangesetOutput, error) {
	testLogger.Info().Msg("Deploying Read Balances contract...")
	deployReadBalanceRequest := &keystone_changeset.DeployRequestV2{ChainSel: chainSelector}
	rbOutput, rbErr := keystone_changeset.DeployBalanceReaderV2(*creEnvironment.CldfEnvironment, deployReadBalanceRequest)
	if rbErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrap(rbErr, "failed to deploy Read Balances contract")
	}

	mergeErr2 := creEnvironment.CldfEnvironment.ExistingAddresses.Merge(rbOutput.AddressBook) //nolint:staticcheck // won't migrate now
	if mergeErr2 != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrap(mergeErr2, "failed to merge address book of Read Balances contract")
	}
	testLogger.Info().Msgf("Read Balances contract deployed to %d", chainSelector)

	readBalancesAddress, _, readContractErr := FindAddressesForChain(
		creEnvironment.CldfEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		chainSelector,
		keystone_changeset.BalanceReader.String(),
	)
	if readContractErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrap(readContractErr, "failed to find Read Balances contract address")
	}
	testLogger.Info().Msgf("Read Balances contract found on chain %d at address %s", chainSelector, readBalancesAddress)

	return readBalancesAddress, rbOutput, nil
}

type DonInfo struct {
	F           uint8
	ConfigCount uint32
	NodeP2PIds  [][32]byte
}

type CapabilitiesRegistry interface {
	GetDON(opts *bind.CallOpts, donID uint32) (DonInfo, error)
}

type registryWrapper struct {
	V1 *kcr.CapabilitiesRegistry
	V2 *capabilities_registry_v2.CapabilitiesRegistry
}

func (rw *registryWrapper) GetDON(opts *bind.CallOpts, donID uint32) (DonInfo, error) {
	if rw.V1 == nil && rw.V2 == nil {
		return DonInfo{}, errors.New("nil capabilities registry contract")
	}

	if rw.V1 != nil && rw.V2 != nil {
		return DonInfo{}, errors.New("invalid registry wrapper state: two versions specified")
	}

	if rw.V1 != nil {
		d, err := rw.V1.GetDON(opts, donID)
		if err != nil {
			return DonInfo{}, err
		}

		return DonInfo{
			F:           d.F,
			ConfigCount: d.ConfigCount,
			NodeP2PIds:  d.NodeP2PIds,
		}, nil
	}

	if rw.V2 != nil {
		d, err := rw.V2.GetDON(opts, donID)
		if err != nil {
			return DonInfo{}, err
		}

		return DonInfo{
			F:           d.F,
			ConfigCount: d.ConfigCount,
			NodeP2PIds:  d.NodeP2PIds,
		}, nil
	}

	return DonInfo{}, errors.New("no valid capabilities registry contract")
}
