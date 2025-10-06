package config

import (
	"context"
	"fmt"
	"maps"
	"math/big"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/chaintype"
	evmconfigtoml "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	chainlinkbig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	coretoml "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

const TronEVMChainID = 3360022319

func Generate(input cre.GenerateConfigsInput, nodeConfigTransformers []cre.NodeConfigTransformerFn) (cre.NodeIndexToConfigOverride, error) {
	configOverrides := make(cre.NodeIndexToConfigOverride)

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	commonInputs, inputsErr := gatherCommonInputs(input)
	if inputsErr != nil {
		return nil, errors.Wrap(inputsErr, "failed to gather common inputs")
	}

	for nodeIdx, nodeMetadata := range input.DonMetadata.NodesMetadata {
		nodeConfig := baseNodeConfig()
		for _, role := range nodeMetadata.Roles {
			switch role {
			case cre.BootstrapNode:
				var cErr error
				nodeConfig, cErr = addBootstrapNodeConfig(nodeConfig, input.OCRPeeringData, input.CapabilitiesPeeringData, commonInputs)
				if cErr != nil {
					return nil, errors.Wrapf(cErr, "failed to add bootstrap node config for node at index %d in DON %s", nodeIdx, input.DonMetadata.Name)
				}
			case cre.WorkerNode:
				var cErr error
				nodeConfig, cErr = addWorkerNodeConfig(nodeConfig, input.OCRPeeringData, input.CapabilitiesPeeringData, commonInputs, input.GatewayConnectorOutput, input.DonMetadata.Name, input.DonMetadata.Flags, nodeMetadata)
				if cErr != nil {
					return nil, errors.Wrapf(cErr, "failed to add worker node config for node at index %d in DON %s", nodeIdx, input.DonMetadata.Name)
				}
			case cre.GatewayNode:
				var cErr error
				nodeConfig, cErr = addGatewayNodeConfig(nodeConfig, commonInputs)
				if cErr != nil {
					return nil, errors.Wrapf(cErr, "failed to add gateway node config for node at index %d in DON %s", nodeIdx, input.DonMetadata.Name)
				}
			default:
				supportedRoles := []string{cre.BootstrapNode, cre.WorkerNode, cre.GatewayNode}
				return nil, fmt.Errorf("unsupported node type %s found for node at index %d in DON %s. Supported roles: %s", role, nodeIdx, input.DonMetadata.Name, strings.Join(supportedRoles, ", "))
			}
		}

		marshalledConfig, mErr := toml.Marshal(nodeConfig)
		if mErr != nil {
			return nil, errors.Wrapf(mErr, "failed to marshal node config for node at index %d in DON %s", nodeIdx, input.DonMetadata.Name)
		}

		configOverrides[nodeIdx] = string(marshalledConfig)
	}

	// execute capability-provided functions that transform the node config (currently: write-evm, write-solana)
	// these functions must return whole node configs after transforming them, instead of just returning configuration parts
	// that need to be merged into the existing config
	for _, transformer := range nodeConfigTransformers {
		if transformer == nil {
			continue
		}

		modifiedConfigs, err := transformer(input, configOverrides)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate nodeset configs")
		}

		maps.Copy(configOverrides, modifiedConfigs)
	}

	return configOverrides, nil
}

func baseNodeConfig() corechainlink.Config {
	return corechainlink.Config{
		Core: coretoml.Core{
			Feature: coretoml.Feature{
				LogPoller: ptr.Ptr(true),
			},
			Log: coretoml.Log{
				JSONConsole: ptr.Ptr(true),
				Level:       ptr.Ptr(coretoml.LogLevel(zapcore.DebugLevel)),
			},
			OCR2: coretoml.OCR2{
				Enabled:              ptr.Ptr(true),
				DatabaseTimeout:      commonconfig.MustNewDuration(1 * time.Second),
				ContractPollInterval: commonconfig.MustNewDuration(1 * time.Second),
			},
		},
	}
}

func addBootstrapNodeConfig(
	existingConfig corechainlink.Config,
	ocrPeeringData cre.OCRPeeringData,
	capabilitiesPeeringData cre.CapabilitiesPeeringData,
	commonInputs *commonInputs,
) (corechainlink.Config, error) {
	ocrBoostrapperLocator, ocrBErr := commontypes.NewBootstrapperLocator(ocrPeeringData.OCRBootstraperPeerID, []string{"localhost:" + strconv.Itoa(ocrPeeringData.Port)})
	if ocrBErr != nil {
		return existingConfig, errors.Wrap(ocrBErr, "failed to create OCR bootstrapper locator")
	}

	existingConfig.P2P = coretoml.P2P{
		V2: coretoml.P2PV2{
			Enabled:              ptr.Ptr(true),
			ListenAddresses:      ptr.Ptr([]string{"0.0.0.0:" + strconv.Itoa(ocrPeeringData.Port)}),
			DefaultBootstrappers: ptr.Ptr([]commontypes.BootstrapperLocator{*ocrBoostrapperLocator}),
		},
	}

	existingConfig.Capabilities = coretoml.Capabilities{
		Peering: coretoml.P2P{
			V2: coretoml.P2PV2{
				Enabled: ptr.Ptr(false),
			},
		},
		SharedPeering: coretoml.SharedPeering{
			Enabled: ptr.Ptr(true),
		},
		Dispatcher: coretoml.Dispatcher{
			SendToSharedPeer: ptr.Ptr(true),
		},
	}

	for _, evmChain := range commonInputs.evmChains {
		appendEVMChain(&existingConfig.EVM, evmChain)
	}

	if commonInputs.solanaChain != nil {
		existingConfig.Solana = append(existingConfig.Solana, &solcfg.TOMLConfig{
			Enabled: ptr.Ptr(true),
			ChainID: ptr.Ptr(commonInputs.solanaChain.ChainID),
			Nodes: []*solcfg.Node{
				{
					Name: &commonInputs.solanaChain.Name,
					URL:  commonconfig.MustParseURL(commonInputs.solanaChain.NodeURL),
				},
			},
		})
	}

	existingConfig.Capabilities.ExternalRegistry = coretoml.ExternalRegistry{
		Address:         ptr.Ptr(commonInputs.capabilityRegistry.address.Hex()),
		NetworkID:       ptr.Ptr("evm"),
		ChainID:         ptr.Ptr(strconv.FormatUint(commonInputs.registryChainID, 10)),
		ContractVersion: ptr.Ptr(commonInputs.capabilityRegistry.versionType.Version.String()),
	}

	return existingConfig, nil
}

func addWorkerNodeConfig(
	existingConfig corechainlink.Config,
	ocrPeeringData cre.OCRPeeringData,
	capabilitiesPeeringData cre.CapabilitiesPeeringData,
	commonInputs *commonInputs,
	gatewayConnector *cre.GatewayConnectorOutput,
	donName string,
	donFlags []string,
	m *cre.NodeMetadata,
) (corechainlink.Config, error) {
	ocrBoostrapperLocator, ocrBErr := commontypes.NewBootstrapperLocator(ocrPeeringData.OCRBootstraperPeerID, []string{ocrPeeringData.OCRBootstraperHost + ":" + strconv.Itoa(ocrPeeringData.Port)})
	if ocrBErr != nil {
		return existingConfig, errors.Wrap(ocrBErr, "failed to create OCR bootstrapper locator")
	}

	existingConfig.P2P = coretoml.P2P{
		V2: coretoml.P2PV2{
			Enabled:              ptr.Ptr(true),
			ListenAddresses:      ptr.Ptr([]string{"0.0.0.0:" + strconv.Itoa(ocrPeeringData.Port)}),
			DefaultBootstrappers: ptr.Ptr([]commontypes.BootstrapperLocator{*ocrBoostrapperLocator}),
		},
	}

	existingConfig.Capabilities = coretoml.Capabilities{
		Peering: coretoml.P2P{
			V2: coretoml.P2PV2{
				Enabled: ptr.Ptr(false),
			},
		},
		SharedPeering: coretoml.SharedPeering{
			Enabled: ptr.Ptr(true),
		},
		Dispatcher: coretoml.Dispatcher{
			SendToSharedPeer: ptr.Ptr(true),
		},
	}

	for _, evmChain := range commonInputs.evmChains {
		appendEVMChain(&existingConfig.EVM, evmChain)
	}

	if commonInputs.solanaChain != nil {
		existingConfig.Solana = append(existingConfig.Solana, &solcfg.TOMLConfig{
			ChainID: ptr.Ptr(commonInputs.solanaChain.ChainID),
			Enabled: ptr.Ptr(true),
			Nodes: solcfg.Nodes{
				{
					Name: ptr.Ptr(commonInputs.solanaChain.Name),
					URL:  commonconfig.MustParseURL(commonInputs.solanaChain.NodeURL),
				},
			},
		})
	}

	existingConfig.Capabilities.ExternalRegistry = coretoml.ExternalRegistry{
		Address:         ptr.Ptr(commonInputs.capabilityRegistry.address.Hex()),
		NetworkID:       ptr.Ptr("evm"),
		ChainID:         ptr.Ptr(strconv.FormatUint(commonInputs.registryChainID, 10)),
		ContractVersion: ptr.Ptr(commonInputs.capabilityRegistry.versionType.Version.String()),
	}

	if flags.HasFlag(donFlags, cre.WorkflowDON) {
		existingConfig.Capabilities.WorkflowRegistry = coretoml.WorkflowRegistry{
			Address:         ptr.Ptr(commonInputs.workflowRegistry.address.Hex()),
			NetworkID:       ptr.Ptr("evm"),
			ChainID:         ptr.Ptr(strconv.FormatUint(commonInputs.registryChainID, 10)),
			SyncStrategy:    ptr.Ptr("reconciliation"),
			ContractVersion: ptr.Ptr(commonInputs.workflowRegistry.versionType.Version.String()),
		}
	}

	if flags.HasFlag(donFlags, cre.WorkflowDON) || don.NodeNeedsAnyGateway(donFlags) {
		evmKey, ok := m.Keys.EVM[commonInputs.registryChainID]
		if !ok {
			return existingConfig, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", commonInputs.registryChainID, m.Index)
		}

		gateways := []coretoml.ConnectorGateway{}
		if gatewayConnector != nil && len(gatewayConnector.Configurations) > 0 {
			for _, gateway := range gatewayConnector.Configurations {
				gateways = append(gateways, coretoml.ConnectorGateway{
					ID: ptr.Ptr(gateway.AuthGatewayID),
					URL: ptr.Ptr(fmt.Sprintf("ws://%s:%d%s",
						gateway.Outgoing.Host,
						gateway.Outgoing.Port,
						gateway.Outgoing.Path)),
				})
			}

			existingConfig.Capabilities.GatewayConnector = coretoml.GatewayConnector{
				DonID:             ptr.Ptr(donName),
				ChainIDForNodeKey: ptr.Ptr(strconv.FormatUint(commonInputs.registryChainID, 10)),
				NodeAddress:       ptr.Ptr(evmKey.PublicAddress.Hex()),
				Gateways:          gateways,
			}
		}
	}

	return existingConfig, nil
}

func addGatewayNodeConfig(
	existingConfig corechainlink.Config,
	commonInputs *commonInputs,
) (corechainlink.Config, error) {
OUTER:
	for _, evmChain := range commonInputs.evmChains {
		// add only unconfigured chains, since other roles might have already added some chains
		for _, existingEVM := range existingConfig.EVM {
			if existingEVM.ChainID.Cmp(chainlinkbig.New(big.NewInt(libc.MustSafeInt64(evmChain.ChainID)))) == 0 {
				continue OUTER
			}
		}
		appendEVMChain(&existingConfig.EVM, evmChain)
	}

	existingConfig.Capabilities.ExternalRegistry = coretoml.ExternalRegistry{
		Address:         ptr.Ptr(commonInputs.capabilityRegistry.address.Hex()),
		NetworkID:       ptr.Ptr("evm"),
		ChainID:         ptr.Ptr(strconv.FormatUint(commonInputs.registryChainID, 10)),
		ContractVersion: ptr.Ptr(commonInputs.capabilityRegistry.versionType.Version.String()),
	}

	existingConfig.Capabilities.WorkflowRegistry = coretoml.WorkflowRegistry{
		Address:         ptr.Ptr(commonInputs.workflowRegistry.address.Hex()),
		NetworkID:       ptr.Ptr("evm"),
		ChainID:         ptr.Ptr(strconv.FormatUint(commonInputs.registryChainID, 10)),
		ContractVersion: ptr.Ptr(commonInputs.workflowRegistry.versionType.Version.String()),
		SyncStrategy:    ptr.Ptr("reconciliation"),
	}

	return existingConfig, nil
}

type addressTypeVersion struct {
	address     common.Address
	versionType cldf.TypeAndVersion
}

type commonInputs struct {
	registryChainID       uint64
	registryChainSelector uint64

	workflowRegistry   addressTypeVersion
	capabilityRegistry addressTypeVersion

	evmChains   []*evmChain
	solanaChain *solanaChain
}

func gatherCommonInputs(input cre.GenerateConfigsInput) (*commonInputs, error) {
	registryChainID, homeErr := chain_selectors.ChainIdFromSelector(input.HomeChainSelector)
	if homeErr != nil {
		return nil, errors.Wrap(homeErr, "failed to get home chain ID")
	}

	// prepare chains, we need chainIDs and URLs
	evmChains := findEVMChains(input)
	solanaChain, solErr := findOneSolanaChain(input)
	if solErr != nil {
		return nil, errors.Wrap(solErr, "failed to find Solana chain in the environment configuration")
	}

	// find contract addresses
	capabilitiesRegistryAddress, capRegTypeVersion, capErr := crecontracts.FindAddressesForChain(input.AddressBook, input.HomeChainSelector, keystone_changeset.CapabilitiesRegistry.String())
	if capErr != nil {
		return nil, errors.Wrap(capErr, "failed to find CapabilitiesRegistry address")
	}

	workflowRegistryAddress, wfRegTypeVersion, wfErr := crecontracts.FindAddressesForChain(input.AddressBook, input.HomeChainSelector, keystone_changeset.WorkflowRegistry.String())
	if wfErr != nil {
		return nil, errors.Wrap(wfErr, "failed to find WorkflowRegistry address")
	}

	return &commonInputs{
		registryChainID:       registryChainID,
		registryChainSelector: input.HomeChainSelector,
		workflowRegistry: addressTypeVersion{
			address:     workflowRegistryAddress,
			versionType: wfRegTypeVersion,
		},
		evmChains:   evmChains,
		solanaChain: solanaChain,
		capabilityRegistry: addressTypeVersion{
			address:     capabilitiesRegistryAddress,
			versionType: capRegTypeVersion,
		},
	}, nil
}

type evmChain struct {
	Name    string
	ChainID uint64
	HTTPRPC string
	WSRPC   string
}

func findEVMChains(input cre.GenerateConfigsInput) []*evmChain {
	evmChains := make([]*evmChain, 0)
	for chainSelector, bcOut := range input.BlockchainOutput {
		if bcOut.SolChain != nil {
			continue
		}

		// if the DON doesn't support the chain, we skip it; if slice is empty, it means that the DON supports all chains
		// TODO: review if we really need this SupportedChains functionality
		if len(input.DonMetadata.CapabilitiesAwareNodeSet().EVMChains()) > 0 && !slices.Contains(input.DonMetadata.CapabilitiesAwareNodeSet().EVMChains(), bcOut.ChainID) {
			continue
		}

		evmChains = append(evmChains, &evmChain{
			Name:    fmt.Sprintf("node-%d", chainSelector),
			ChainID: bcOut.ChainID,
			HTTPRPC: bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
			WSRPC:   bcOut.BlockchainOutput.Nodes[0].InternalWSUrl,
		})
	}
	return evmChains
}

type solanaChain struct {
	Name    string
	ChainID string
	NodeURL string
}

func findOneSolanaChain(input cre.GenerateConfigsInput) (*solanaChain, error) {
	var solChain *solanaChain
	chainsFound := 0

	for _, bcOut := range input.BlockchainOutput {
		if bcOut.SolChain == nil {
			continue
		}

		chainsFound++
		if chainsFound > 1 {
			return nil, errors.New("multiple Solana chains found, expected only one")
		}

		ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
		chainID, err := bcOut.SolClient.GetGenesisHash(ctx)
		if err != nil {
			cancelFn()
			return nil, errors.Wrap(err, "failed to get chainID for Solana")
		}
		cancelFn()

		solChain = &solanaChain{
			Name:    fmt.Sprintf("node-%d", bcOut.SolChain.ChainSelector),
			ChainID: chainID.String(),
			NodeURL: bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
		}
	}

	return solChain, nil
}

func buildTronEVMConfig(evmChain *evmChain) evmconfigtoml.EVMConfig {
	tronRPC := strings.Replace(evmChain.HTTPRPC, "jsonrpc", "wallet", 1)
	return evmconfigtoml.EVMConfig{
		ChainID: chainlinkbig.New(big.NewInt(libc.MustSafeInt64(evmChain.ChainID))),
		Chain: evmconfigtoml.Chain{
			AutoCreateKey:         ptr.Ptr(false),
			ChainType:             chaintype.NewConfig("tron"),
			LogBroadcasterEnabled: ptr.Ptr(false),
			NodePool: evmconfigtoml.NodePool{
				NewHeadsPollInterval: commonconfig.MustNewDuration(10 * time.Second),
			},
		},
		Nodes: []*evmconfigtoml.Node{
			{
				Name:              ptr.Ptr(evmChain.Name),
				HTTPURL:           commonconfig.MustParseURL(evmChain.HTTPRPC),
				HTTPURLExtraWrite: commonconfig.MustParseURL(tronRPC),
			},
		},
	}
}

func buildEVMConfig(evmChain *evmChain) evmconfigtoml.EVMConfig {
	return evmconfigtoml.EVMConfig{
		ChainID: chainlinkbig.New(big.NewInt(libc.MustSafeInt64(evmChain.ChainID))),
		Chain: evmconfigtoml.Chain{
			AutoCreateKey: ptr.Ptr(false),
		},
		Nodes: []*evmconfigtoml.Node{
			{
				Name:    ptr.Ptr(evmChain.Name),
				WSURL:   commonconfig.MustParseURL(evmChain.WSRPC),
				HTTPURL: commonconfig.MustParseURL(evmChain.HTTPRPC),
			},
		},
	}
}

func appendEVMChain(existingConfig *evmconfigtoml.EVMConfigs, evmChain *evmChain) {
	var cfg evmconfigtoml.EVMConfig
	if evmChain.ChainID == TronEVMChainID {
		cfg = buildTronEVMConfig(evmChain)
	} else {
		cfg = buildEVMConfig(evmChain)
	}
	*existingConfig = append(*existingConfig, &cfg)
}
