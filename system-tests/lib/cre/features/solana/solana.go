package solana

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	ks_sol_seq "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana/sequence"
	ks_sol_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana/sequence/operation"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"

	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

const flag = cre.WriteSolanaCapability

type Solana struct{}

func (o *Solana) Flag() cre.CapabilityFlag {
	return flag
}

func (o *Solana) PreEnvStartup(
	testLogger zerolog.Logger,
	registryChainSelector uint64,
	cldfEnv *cldf.Environment,
	provider infra.Provider,
	topology *cre.Topology,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	capabilityConfigs cre.CapabilityConfigs,
	contractVersions map[string]string,
	gatewayJobConfigs map[cre.NodeUUID]*config.GatewayConfig,
) (*cre.PreEnvStartupOutput, error) {
	donsMetadata := topology.DonsMetadataWithFlag(flag)
	if len(donsMetadata) == 0 {
		return nil, nil
	}

	if len(donsMetadata) > 1 {
		// there's only one Solana chain out there in reality, hence the limitation
		return nil, errors.New("only one DON with write-solana capability is supported")
	}

	solanaDON := donsMetadata[0]

	// deploy forwarders
	solForwardersSelectors := make([]uint64, 0)
	for _, bcOut := range blockchainOutputs {
		// consider we have just 1 solana chain
		if bcOut.SolChain != nil {
			solForwardersSelectors = append(solForwardersSelectors, bcOut.SolChain.ChainSelector)
			continue
		}
	}

	memoryDatastore := datastore.NewMemoryDataStore()
	// load all existing addresses into memory datastore
	mergeErr := memoryDatastore.Merge(cldfEnv.DataStore)
	if mergeErr != nil {
		return nil, fmt.Errorf("failed to merge existing datastore into memory datastore: %w", mergeErr)
	}

	for _, sel := range solForwardersSelectors {
		populateContracts := map[string]datastore.ContractType{
			deployment.KeystoneForwarderProgramName: ks_sol.ForwarderContract,
		}
		version := semver.MustParse(contractVersions[ks_sol.ForwarderContract.String()])

		// Forwarder for solana is predeployed on chain spin-up. We jus need to add it to memory datastore here
		errp := memory.PopulateDatastore(memoryDatastore.AddressRefStore, populateContracts,
			version, ks_sol.DefaultForwarderQualifier, sel)
		if errp != nil {
			return nil, errors.Wrap(errp, "failed to populate datastore with predeployed contracts")
		}
		out, err := operations.ExecuteSequence(
			cldfEnv.OperationsBundle,
			ks_sol_seq.DeployForwarderSeq,
			ks_sol_op.Deps{
				Env:       *cldfEnv,
				Chain:     cldfEnv.BlockChains.SolanaChains()[sel],
				Datastore: memoryDatastore.Seal(),
			},
			ks_sol_seq.DeployForwarderSeqInput{
				ChainSel:     sel,
				ProgramName:  deployment.KeystoneForwarderProgramName,
				Qualifier:    ks_sol.DefaultForwarderQualifier,
				ContractType: ks_sol.ForwarderContract,
				Version:      version,
			},
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to deploy sol forwarder")
		}

		err = memoryDatastore.AddressRefStore.Add(datastore.AddressRef{
			Address:       out.Output.State.String(),
			ChainSelector: sel,
			Version:       semver.MustParse(contractVersions[ks_sol.ForwarderState.String()]),
			Qualifier:     ks_sol.DefaultForwarderQualifier,
			Type:          ks_sol.ForwarderState,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to add address to the datastore for Solana Forwarder state")
		}

		testLogger.Info().Msgf("Deployed Forwarder %s contract on Solana chain chain %d programID: %s state: %s", contractVersions[ks_sol.ForwarderContract.String()], sel, out.Output.ProgramID.String(), out.Output.State.String())
	}

	cldfEnv.DataStore = memoryDatastore.Seal()

	// modify node config
	data := solanaInput{}
	for _, bcOut := range blockchainOutputs {
		if bcOut.SolChain == nil {
			continue
		}
		data.ChainSelector = bcOut.SolChain.ChainSelector
		// find Solana forwarder address
		forwarders := cldfEnv.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(data.ChainSelector))
		for _, addr := range forwarders {
			if addr.Type == ks_sol.ForwarderState {
				data.ForwarderState = addr.Address
				continue
			}
			data.ForwarderAddress = addr.Address
		}

		break
	}

	workerNodes, wErr := solanaDON.Workers()
	if wErr != nil {
		return nil, errors.Wrap(wErr, "failed to find worker nodes")
	}

	for _, workerNode := range workerNodes {
		chainID, chErr := chainselectors.SolanaChainIdFromSelector(data.ChainSelector)
		if chErr != nil {
			return nil, errors.Wrapf(chErr, "failed to get Solana chain ID from selector %d", data.ChainSelector)
		}

		key, ok := workerNode.Keys.Solana[chainID]
		if !ok {
			return nil, errors.Errorf("missing Solana key for chainID %s on node index %d", chainID, workerNode.Index)
		}
		data.FromAddress = key.PublicAddress

		if capabilityConfigs == nil {
			return nil, errors.New("additional capabilities configs are nil, but are required to configure the write-solana capability")
		}

		if writeSolConfig, ok := capabilityConfigs[cre.WriteSolanaCapability]; ok {
			mergedConfig := envconfig.ResolveCapabilityConfigForDON(
				cre.WriteSolanaCapability,
				writeSolConfig.Config,
				nil,
			)

			runtimeValues := map[string]any{
				"FromAddress":      data.FromAddress.String(),
				"ForwarderAddress": data.ForwarderAddress,
				"ForwarderState":   data.ForwarderState,
			}

			var mErr error
			data.WorkflowConfig, mErr = don.ApplyRuntimeValues(mergedConfig, runtimeValues)
			if mErr != nil {
				return nil, errors.Wrap(mErr, "failed to apply runtime values")
			}
		} else {
			fmt.Println("sol config not found")
		}

		currentConfig := solanaDON.CapabilitiesAwareNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides

		var typedConfig corechainlink.Config
		unmarshallErr := toml.Unmarshal([]byte(currentConfig), &typedConfig)
		if unmarshallErr != nil {
			return nil, errors.Wrapf(unmarshallErr, "failed to unmarshal config for node index %d", workerNode.Index)
		}

		if len(typedConfig.Solana) != 1 {
			return nil, fmt.Errorf("only 1 Solana chain is supported, but found %d for node at index %d", len(typedConfig.Solana), workerNode.Index)
		}

		if typedConfig.Solana[0].ChainID == nil {
			return nil, fmt.Errorf("solana chainID is nil for node at index %d", workerNode.Index)
		}

		var solCfg solcfg.WorkflowConfig

		// Execute template with chain's workflow configuration
		tmpl, err := template.New("solanaWorkflowConfig").Parse(solWorkflowConfigTemplate)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse Solana workflow config template")
		}
		var configBuffer bytes.Buffer
		if executeErr := tmpl.Execute(&configBuffer, data.WorkflowConfig); executeErr != nil {
			return nil, errors.Wrap(executeErr, "failed to execute Solana workflow config template")
		}

		configStr := configBuffer.String()

		if err := don.ValidateTemplateSubstitution(configStr, flag); err != nil {
			return nil, errors.Wrapf(err, "%s template validation failed", flag)
		}

		unmarshallErr = toml.Unmarshal([]byte(configStr), &solCfg)
		if unmarshallErr != nil {
			return nil, errors.Wrap(unmarshallErr, "failed to unmarshal Solana.Workflow config")
		}

		typedConfig.Solana[0].Workflow = solCfg

		stringifiedConfig, mErr := toml.Marshal(typedConfig)
		if mErr != nil {
			return nil, errors.Wrapf(mErr, "failed to marshal config for node index %d", workerNode.Index)
		}

		solanaDON.CapabilitiesAwareNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = string(stringifiedConfig)
	}

	capabilities := make(map[uint64][]keystone_changeset.DONCapabilityWithConfig)
	fullName := "write_solana_devnet@1.0.0"
	splitName := strings.Split(fullName, "@")

	capabilities[solanaDON.ID] = append(capabilities[solanaDON.ID], keystone_changeset.DONCapabilityWithConfig{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   splitName[0],
			Version:        splitName[1],
			CapabilityType: 3, // TARGET
			ResponseType:   1, // OBSERVATION_IDENTICAL
		},
		Config: &capabilitiespb.CapabilityConfig{},
	})

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfigs: capabilities,
	}, nil
}

type solanaInput struct {
	ChainSelector    uint64
	FromAddress      solana.PublicKey
	ForwarderAddress string
	ForwarderState   string
	HasWrite         bool
	WorkflowConfig   map[string]any // Configuration for Solana.Workflow section
}

const solWorkflowConfigTemplate = `
		ForwarderAddress = '{{.ForwarderAddress}}'
		FromAddress      = '{{.FromAddress}}'
		ForwarderState   = '{{.ForwarderState}}'
		PollPeriod = '{{.PollPeriod}}'
		AcceptanceTimeout = '{{.AcceptanceTimeout}}'
		TxAcceptanceState = {{.TxAcceptanceState}}
		Local = {{.Local}}
	`

func (o *Solana) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	nodeSetOutput []*cre.WrappedNodeOutput,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	contractVersions map[string]string,
	provider infra.Provider,
	capabilityConfigs map[string]cre.CapabilityConfig,
) error {
	dons := creEnv.DonTopology.DonsWithFlag(flag)
	if len(dons) == 0 {
		return nil
	}

	solChainsWithForwarder := make(map[uint64]struct{})
	solForwarders := creEnv.CldfEnvironment.DataStore.Addresses().Filter(datastore.AddressRefByQualifier(ks_sol.DefaultForwarderQualifier))
	for _, forwarder := range solForwarders {
		solChainsWithForwarder[forwarder.ChainSelector] = struct{}{}
	}

	// configure Solana forwarder only if we have some
	if len(solChainsWithForwarder) > 0 {
		for _, don := range dons {
			cs := commonchangeset.Configure(ks_sol.ConfigureForwarders{},
				&ks_sol.ConfigureForwarderRequest{
					WFDonName:        don.Name,
					WFNodeIDs:        don.KeystoneDONConfig().NodeIDs,
					RegistryChainSel: creEnv.DonTopology.HomeChainSelector,
					Chains:           solChainsWithForwarder,
					Qualifier:        ks_sol.DefaultForwarderQualifier,
					Version:          "1.0.0",
				},
			)

			_, err := cs.Apply(*creEnv.CldfEnvironment)
			if err != nil {
				return errors.Wrap(err, "failed to configure Solana forwarders")
			}
		}
	}

	return nil
}
