package contracts

import (
	"context"
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

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	cap_reg_v2_seq "github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
	cre_contracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	syncer_v2 "github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer/v2"
)

type DeployKeystoneContractsInput struct {
	CldfEnvironment  *cldf.Environment
	CtfBlockchains   []blockchains.Blockchain
	ContractVersions map[string]string
	WithV2Registries bool
}

type DeployKeystoneContractsOutput struct {
	Env             *cldf.Environment
	MemoryDataStore *datastore.MemoryDataStore
}

func DeployKeystoneContracts(
	ctx context.Context,
	testLogger zerolog.Logger,
	singleFileLogger logger.Logger,
	input DeployKeystoneContractsInput,
) (*DeployKeystoneContractsOutput, error) {
	memoryDatastore := datastore.NewMemoryDataStore()

	registryChainOutput := input.CtfBlockchains[0]
	registryChainSelector := registryChainOutput.ChainSelector()
	deployRegistrySeq := ks_contracts_op.DeployRegistryContractsSequence
	if input.WithV2Registries {
		deployRegistrySeq = ks_contracts_op.DeployV2RegistryContractsSequence
	}

	registryContractsReport, seqErr := operations.ExecuteSequence(
		input.CldfEnvironment.OperationsBundle,
		deployRegistrySeq,
		ks_contracts_op.DeployContractsSequenceDeps{
			Env: input.CldfEnvironment,
		},
		ks_contracts_op.DeployRegistryContractsSequenceInput{
			RegistryChainSelector: registryChainSelector,
		},
	)
	if seqErr != nil {
		return nil, errors.Wrap(seqErr, "failed to deploy Keystone contracts")
	}

	if err := input.CldfEnvironment.ExistingAddresses.Merge(registryContractsReport.Output.AddressBook); err != nil { //nolint:staticcheck // won't migrate now
		return nil, errors.Wrap(err, "failed to merge address book with Keystone contracts addresses")
	}

	if err := memoryDatastore.Merge(registryContractsReport.Output.Datastore); err != nil {
		return nil, errors.Wrap(err, "failed to merge datastore with Keystone contracts addresses")
	}

	wfRegAddr := MustGetAddressFromMemoryDataStore(memoryDatastore, registryChainSelector, keystone_changeset.WorkflowRegistry.String(), input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")
	testLogger.Info().Msgf("Deployed Workflow Registry %s contract on chain %d at %s", input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], registryChainSelector, wfRegAddr)

	capRegAddr := MustGetAddressFromMemoryDataStore(memoryDatastore, registryChainSelector, keystone_changeset.CapabilitiesRegistry.String(), input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()], "")
	testLogger.Info().Msgf("Deployed Capabilities Registry %s contract on chain %d at %s", input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()], registryChainSelector, capRegAddr)

	input.CldfEnvironment.DataStore = memoryDatastore.Seal()

	return &DeployKeystoneContractsOutput{
		Env:             input.CldfEnvironment,
		MemoryDataStore: memoryDatastore,
	}, nil
}

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
	nodes := make([]contracts.NodesInput, 0)
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

					nodes = append(nodes, contracts.NodesInput{
						NOP:                 nopName,
						P2pID:               n.PeerID,
						Signer:              ocrCfg.OffchainPublicKey,
						EncryptionPublicKey: [32]byte(wfKey),
						CsaKey:              [32]byte(csKey),
						CapabilityIDs:       capIDs,
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
	hexDigits := max(int(math.Ceil(math.Log10(float64(count+1))/math.Log10(16))), 1)

	// The total length of the address after the "0x" prefix must be 40.
	baseHexLen := 40 - hexDigits
	if baseHexLen <= 0 {
		return nil, errors.New("count is too large to generate unique addresses with this base")
	}

	// Create a base string of 'f' characters to ensure the addresses are not zero.
	baseString := strings.Repeat("f", baseHexLen)

	addresses := make([]common.Address, count)
	for i := range count {
		format := fmt.Sprintf("%s%%0%dx", baseString, hexDigits)
		fullAddress := fmt.Sprintf(format, i)
		addresses[i] = common.HexToAddress("0x" + fullAddress)
	}

	return addresses, nil
}

func toDons(input cre.ConfigureCapabilityRegistryInput) (*dons, error) {
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

func ConfigureCapabilityRegistry(input cre.ConfigureCapabilityRegistryInput) (CapabilitiesRegistry, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	dons, dErr := toDons(input)
	if dErr != nil {
		return nil, errors.Wrap(dErr, "failed to map input to dons")
	}
	if !input.WithV2Registries {
		_, seqErr := operations.ExecuteSequence(
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
		if seqErr != nil {
			return nil, errors.Wrap(seqErr, "failed to configure capabilities registry")
		}

		capReg, cErr := cre_contracts.GetOwnedContractV2[*kcr.CapabilitiesRegistry](
			input.CldEnv.DataStore.Addresses(),
			input.CldEnv.BlockChains.EVMChains()[input.ChainSelector],
			input.CapabilitiesRegistryAddress.Hex(),
		)
		if cErr != nil {
			return nil, errors.Wrap(cErr, "failed to get capabilities registry contract")
		}
		return &registryWrapper{V1: capReg.Contract}, nil
	}

	// Transform dons data to V2 sequence input format
	v2Input := dons.mustToV2ConfigureInput(input.ChainSelector, input.CapabilitiesRegistryAddress.Hex())
	_, seqErr := operations.ExecuteSequence(
		input.CldEnv.OperationsBundle,
		cap_reg_v2_seq.ConfigureCapabilitiesRegistry,
		cap_reg_v2_seq.ConfigureCapabilitiesRegistryDeps{
			Env: input.CldEnv,
		},
		v2Input,
	)
	if seqErr != nil {
		return nil, errors.Wrap(seqErr, "failed to configure capabilities registry")
	}

	capReg, cErr := cre_contracts.GetOwnedContractV2[*capabilities_registry_v2.CapabilitiesRegistry](
		input.CldEnv.DataStore.Addresses(),
		input.CldEnv.BlockChains.EVMChains()[input.ChainSelector],
		input.CapabilitiesRegistryAddress.Hex(),
	)
	if cErr != nil {
		return nil, errors.Wrap(cErr, "failed to get capabilities registry contract")
	}

	return &registryWrapper{V2: capReg.Contract}, nil
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
