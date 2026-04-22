package contracts

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
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
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	syncer_v2 "github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer/v2"
)

type DeployKeystoneContractsInput struct {
	CldfEnvironment  *cldf.Environment
	CtfBlockchains   []blockchains.Blockchain
	ContractVersions map[cre.ContractType]*semver.Version
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
	c                     map[string]donConfig
	offChain              offchain.Client
	env                   *cldf.Environment
	registryChainSelector uint64
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

// embedOCR3Config computes the full OCR3 configuration for a consensus V2 DON
// and embeds it in the capability config proto's Ocr3Configs map.
func (d *dons) embedOCR3Config(capConfig *capabilitiespb.CapabilityConfig, don donConfig, registryChainSelector uint64, oracleConfig *ocr3.OracleConfig, extraSignerFamilies []string) error {
	oracleConfig.TransmissionSchedule = []int{len(don.Nops[0].Nodes)}

	var allNodeIDs []string
	for _, nop := range don.Nops {
		allNodeIDs = append(allNodeIDs, nop.Nodes...)
	}

	nodes, err := deployment.NodeInfo(allNodeIDs, d.offChain)
	if err != nil {
		return fmt.Errorf("failed to get node info: %w", err)
	}

	ocrConfig, err := ocr3.GenerateOCR3ConfigFromNodes(*oracleConfig, nodes, registryChainSelector, d.env.OCRSecrets, nil, extraSignerFamilies)
	if err != nil {
		return fmt.Errorf("failed to generate OCR3 config: %w", err)
	}

	transmitterBytes := make([][]byte, len(ocrConfig.Transmitters))
	for i, t := range ocrConfig.Transmitters {
		transmitterBytes[i] = t.Bytes()
	}

	ocr3Proto := &capabilitiespb.OCR3Config{
		Signers:               ocrConfig.Signers,
		Transmitters:          transmitterBytes,
		F:                     uint32(ocrConfig.F),
		OnchainConfig:         ocrConfig.OnchainConfig,
		OffchainConfigVersion: ocrConfig.OffchainConfigVersion,
		OffchainConfig:        ocrConfig.OffchainConfig,
		ConfigCount:           1,
	}

	if capConfig.Ocr3Configs == nil {
		capConfig.Ocr3Configs = make(map[string]*capabilitiespb.OCR3Config)
	}
	capConfig.Ocr3Configs[capabilitiespb.OCR3ConfigDefaultKey] = ocr3Proto

	return nil
}

func (d *dons) mustToV2ConfigureInput(chainSelector uint64, contractAddress string, capabilityToOCR3Config map[string]*ocr3.OracleConfig, capabilityToExtraSignerFamilies map[string][]string) cap_reg_v2_seq.ConfigureCapabilitiesRegistryInput {
	nops := make([]capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams, 0)
	nodes := make([]contracts.NodesInput, 0)
	capabilities := make([]contracts.RegisterableCapability, 0)
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
				ns, err := deployment.NodeInfo(nop.Nodes, d.offChain)
				if err != nil {
					panic(err)
				}
				nopMap[nopName] = capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams{
					Admin: adminAddrs[i],
					Name:  nopName,
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

			capConfig := cap.Config
			shouldMarshalProtoConfig := capConfig != nil
			if cap.UseCapRegOCRConfig {
				if capConfig == nil {
					capConfig = &capabilitiespb.CapabilityConfig{}
				}
				shouldMarshalProtoConfig = true

				ocrConfig := capabilityToOCR3Config[cap.Capability.LabelledName]
				if ocrConfig == nil {
					panic("no OCR3 config found for capability " + cap.Capability.LabelledName)
				}
				if err := d.embedOCR3Config(capConfig, don, chainSelector, ocrConfig, capabilityToExtraSignerFamilies[cap.Capability.LabelledName]); err != nil {
					panic(fmt.Sprintf("failed to embed OCR3 config for capability %s: %s", cap.Capability.LabelledName, err))
				}
			}

			if shouldMarshalProtoConfig {
				if protoBytes, err := proto.Marshal(capConfig); err == nil {
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
	for _, cp := range capabilityMap {
		var metadata map[string]any
		err := json.Unmarshal(cp.Metadata, &metadata)
		if err != nil {
			panic(fmt.Sprintf("failed to unmarshal capability metadata: %s", err))
		}
		capabilities = append(capabilities, contracts.RegisterableCapability{
			Metadata:              metadata,
			CapabilityID:          cp.CapabilityId,
			ConfigurationContract: cp.ConfigurationContract,
		})
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
		c:                     make(map[string]donConfig),
		offChain:              input.CldEnv.Offchain,
		env:                   input.CldEnv,
		registryChainSelector: input.ChainSelector,
	}

	for donIdx, donMetadata := range input.Topology.DonsMetadata.List() {
		// if it's only a gateway or bootstrapDON, we don't want to register it with the Capabilities Registry
		// since it doesn't have any capabilities
		if flags.HasNoOtherFlags(donMetadata.Flags, []string{cre.GatewayDON, cre.BootstrapDON}) {
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

		// add capabilities that were passed directly via feature startup hooks
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

func ConfigureCapabilityRegistry(ctx context.Context, input cre.ConfigureCapabilityRegistryInput) (CapabilitiesRegistry, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	dons, dErr := toDons(input)
	if dErr != nil {
		return nil, errors.Wrap(dErr, "failed to map input to dons")
	}
	var capReg CapabilitiesRegistry
	if !input.WithV2Registries {
		for _, don := range dons.donsOrderedByID() {
			for i, cap := range don.Capabilities {
				if !cap.UseCapRegOCRConfig || cap.Config == nil {
					continue
				}
				ocrConfig := input.CapabilityToOCR3Config[cap.Capability.LabelledName]
				if ocrConfig == nil {
					return nil, fmt.Errorf("no OCR3 config found for capability %s", cap.Capability.LabelledName)
				}
				if err := dons.embedOCR3Config(don.Capabilities[i].Config, don, input.ChainSelector, ocrConfig, input.CapabilityToExtraSignerFamilies[cap.Capability.LabelledName]); err != nil {
					return nil, fmt.Errorf("failed to embed OCR3 config for capability %s: %w", cap.Capability.LabelledName, err)
				}
			}
		}

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

		capRegContract, cErr := cre_contracts.GetOwnedContractV2[*kcr.CapabilitiesRegistry](
			input.CldEnv.DataStore.Addresses(),
			input.CldEnv.BlockChains.EVMChains()[input.ChainSelector],
			input.CapabilitiesRegistryAddress.Hex(),
			"",
		)
		if cErr != nil {
			return nil, errors.Wrap(cErr, "failed to get capabilities registry contract")
		}
		capReg = &registryWrapper{V1: capRegContract.Contract}
	} else {
		// Transform dons data to V2 sequence input format
		v2Input := dons.mustToV2ConfigureInput(input.ChainSelector, input.CapabilitiesRegistryAddress.Hex(), input.CapabilityToOCR3Config, input.CapabilityToExtraSignerFamilies)
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

		capRegContract, cErr := cre_contracts.GetOwnedContractV2[*capabilities_registry_v2.CapabilitiesRegistry](
			input.CldEnv.DataStore.Addresses(),
			input.CldEnv.BlockChains.EVMChains()[input.ChainSelector],
			input.CapabilitiesRegistryAddress.Hex(),
			"",
		)
		if cErr != nil {
			return nil, errors.Wrap(cErr, "failed to get capabilities registry contract")
		}

		capReg = &registryWrapper{V2: capRegContract.Contract}
	}

	// TODO: remove this once the race condition is fixed (CRE-2684)
	if waitErr := waitForWorkflowWorkersCapabilityRegistrySync(ctx, input); waitErr != nil {
		return nil, errors.Wrap(waitErr, "failed waiting for workflow nodes to sync capability registry state")
	}

	return capReg, nil
}

type DonInfo struct {
	ID          uint32
	F           uint8
	ConfigCount uint32
	NodeP2PIds  [][32]byte
}

// DonForResolution holds the data needed to resolve a DON's contract ID from the Capabilities Registry.
type DonForResolution struct {
	Name       string
	NodeP2PIds [][32]byte
}

type CapabilitiesRegistry interface {
	GetDON(opts *bind.CallOpts, donID uint32) (DonInfo, error)
	GetDONByName(opts *bind.CallOpts, donName string) (DonInfo, error)
	GetDONs(opts *bind.CallOpts) ([]DonInfo, error)
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
			ID:          d.Id,
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
			ID:          d.Id,
			F:           d.F,
			ConfigCount: d.ConfigCount,
			NodeP2PIds:  d.NodeP2PIds,
		}, nil
	}

	return DonInfo{}, errors.New("no valid capabilities registry contract")
}

func (rw *registryWrapper) GetDONByName(opts *bind.CallOpts, donName string) (DonInfo, error) {
	if rw.V1 != nil {
		return DonInfo{}, errors.New("GetDONByName not supported for V1 capabilities registry")
	}
	if rw.V2 == nil {
		return DonInfo{}, errors.New("nil capabilities registry contract")
	}

	d, err := rw.V2.GetDONByName(opts, donName)
	if err != nil {
		return DonInfo{}, err
	}

	return DonInfo{
		ID:          d.Id,
		F:           d.F,
		ConfigCount: d.ConfigCount,
		NodeP2PIds:  d.NodeP2PIds,
	}, nil
}

func (rw *registryWrapper) GetDONs(opts *bind.CallOpts) ([]DonInfo, error) {
	if rw.V2 != nil {
		return nil, errors.New("GetDONs for V2 should use GetDONs with pagination; use GetDONByName for resolution")
	}
	if rw.V1 == nil {
		return nil, errors.New("nil capabilities registry contract")
	}

	dons, err := rw.V1.GetDONs(opts)
	if err != nil {
		return nil, err
	}

	result := make([]DonInfo, len(dons))
	for i, d := range dons {
		result[i] = DonInfo{
			ID:          d.Id,
			F:           d.F,
			ConfigCount: d.ConfigCount,
			NodeP2PIds:  d.NodeP2PIds,
		}
	}
	return result, nil
}

// ResolveContractDonIDs retrieves the actual contract donIDs for the given DONs.
// For V2: uses GetDONByName with donName = don.Name + "-don"
// For V1: uses GetDONs and matches by node P2P IDs (one match is enough; P2P IDs are unique)
func ResolveContractDonIDs(
	capReg CapabilitiesRegistry,
	dons []DonForResolution,
	withV2Registries bool,
) (map[string]uint32, error) {
	result := make(map[string]uint32)

	if withV2Registries {
		for _, don := range dons {
			donName := don.Name + "-don"
			info, err := capReg.GetDONByName(nil, donName)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get DON by name %s", donName)
			}
			result[don.Name] = info.ID
		}
		return result, nil
	}

	// V1: match by P2P IDs
	contractDons, err := capReg.GetDONs(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get DONs from capabilities registry")
	}

	// Build an index once so resolution is O(#ourP2P) per DON.
	contractDonIDByP2P := make(map[[32]byte]uint32)
	for _, cd := range contractDons {
		for _, contractP2P := range cd.NodeP2PIds {
			if existingDonID, exists := contractDonIDByP2P[contractP2P]; exists && existingDonID != cd.ID {
				return nil, fmt.Errorf(
					"duplicate contract P2P ID found across DONs: p2pID=%x, donIDs=%d and %d",
					contractP2P,
					existingDonID,
					cd.ID,
				)
			}
			contractDonIDByP2P[contractP2P] = cd.ID
		}
	}

	for _, don := range dons {
		found := false
		for _, ourP2P := range don.NodeP2PIds {
			if donID, ok := contractDonIDByP2P[ourP2P]; ok {
				result[don.Name] = donID
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("failed to resolve contract donID for DON %s: no matching P2P IDs in registry", don.Name)
		}
	}
	return result, nil
}

// ResolveAndApplyContractDonIDs resolves contract donIDs from the Capabilities Registry and applies them
// to topology, dons, and nodeSets. Uses GetDONByName (V2) or GetDONs+P2P match (V1).
func ResolveAndApplyContractDonIDs(
	capReg CapabilitiesRegistry,
	dons *cre.Dons,
	topology *cre.Topology,
	nodeSets []*cre.NodeSet,
	withV2Registries bool,
) error {
	resolvedDonIDs, err := resolveContractDonIDsFromDons(capReg, dons, withV2Registries)
	if err != nil {
		return err
	}
	if len(resolvedDonIDs) == 0 {
		return nil
	}

	return applyResolvedContractDonIDs(resolvedDonIDs, nodeSets, dons, topology)
}

func resolveContractDonIDsFromDons(
	capReg CapabilitiesRegistry,
	dons *cre.Dons,
	withV2Registries bool,
) (map[string]uint32, error) {
	registeredDons := make([]*cre.Don, 0)
	for _, don := range dons.List() {
		if !flags.HasNoOtherFlags(don.Flags, []string{cre.GatewayDON, cre.BootstrapDON}) {
			registeredDons = append(registeredDons, don)
		}
	}
	if len(registeredDons) == 0 {
		return nil, nil
	}

	donsForResolution := make([]DonForResolution, 0, len(registeredDons))
	for _, don := range registeredDons {
		workers, err := don.Workers()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get workers for DON %s", don.Name)
		}

		nodeP2PIds := make([][32]byte, 0, len(workers))
		for _, worker := range workers {
			if worker.Keys.P2PKey == nil {
				continue
			}
			peerID, err := p2pkey.MakePeerID(worker.Keys.P2PKey.PeerID.String())
			if err != nil {
				return nil, errors.Wrapf(err, "failed to make peer ID for node in DON %s", don.Name)
			}
			nodeP2PIds = append(nodeP2PIds, peerID)
		}
		if len(nodeP2PIds) == 0 {
			return nil, fmt.Errorf("DON %s has no valid worker P2P IDs", don.Name)
		}

		donsForResolution = append(donsForResolution, DonForResolution{
			Name:       don.Name,
			NodeP2PIds: nodeP2PIds,
		})
	}

	return ResolveContractDonIDs(capReg, donsForResolution, withV2Registries)
}

func applyResolvedContractDonIDs(
	resolvedDonIDs map[string]uint32,
	nodeSets []*cre.NodeSet,
	dons *cre.Dons,
	topology *cre.Topology,
) error {
	workflowDonsMetadata, wErr := topology.DonsMetadata.WorkflowDONs()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to get workflow DONs metadata")
	}

	topology.WorkflowDONIDs = make([]uint64, 0, len(workflowDonsMetadata))
	for _, donMeta := range workflowDonsMetadata {
		if id, ok := resolvedDonIDs[donMeta.Name]; ok {
			topology.WorkflowDONIDs = append(topology.WorkflowDONIDs, uint64(id))
			donMeta.ID = uint64(id)
		}
	}
	for _, donMeta := range topology.DonsMetadata.List() {
		if !flags.HasNoOtherFlags(donMeta.Flags, []string{cre.GatewayDON, cre.BootstrapDON}) {
			if id, ok := resolvedDonIDs[donMeta.Name]; ok {
				donMeta.ID = uint64(id)
			}
		}
	}
	for _, don := range dons.List() {
		if !flags.HasNoOtherFlags(don.Flags, []string{cre.GatewayDON, cre.BootstrapDON}) {
			if id, ok := resolvedDonIDs[don.Name]; ok {
				don.ID = uint64(id)
			}
		}
	}
	for _, ns := range nodeSets {
		if !flags.HasNoOtherFlags(ns.Flags(), []string{cre.GatewayDON, cre.BootstrapDON}) {
			if id, ok := resolvedDonIDs[ns.Name]; ok {
				ns.ContractDonID = uint64(id)
			}
		}
	}

	return nil
}
