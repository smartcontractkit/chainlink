package solana

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/Masterminds/semver/v3"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
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
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/solana"
)

const flag = cre.WriteSolanaCapability

type Solana struct{}

func (o *Solana) Flag() cre.CapabilityFlag {
	return flag
}

func (o *Solana) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	var solChain *solana.Blockchain
	for _, bcOut := range creEnv.Blockchains {
		if bcOut.IsFamily(chainselectors.FamilySolana) {
			solChain = bcOut.(*solana.Blockchain)
			break
		}
	}

	// TODO check if not deployed yet
	programID, state, fErr := deployForwarder(testLogger, creEnv, solChain)
	if fErr != nil {
		return nil, errors.Wrap(fErr, "failed to deploy solana forwarder")
	}

	chainID, chErr := chainselectors.SolanaChainIdFromSelector(solChain.ChainSelector())
	if chErr != nil {
		return nil, errors.Wrapf(chErr, "failed to get Solana chain ID from selector %d", solChain.ChainSelector())
	}

	data := solanaInput{
		ChainSelector:    solChain.ChainSelector(),
		ForwarderAddress: *programID,
		ForwarderState:   *state,
	}

	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return nil, errors.Wrap(wErr, "failed to find worker nodes")
	}

	for _, workerNode := range workerNodes {
		currentConfig := don.NodeSets().NodeSpecs[workerNode.Index].Node.TestConfigOverrides
		updatedConfig, uErr := updateNodeConfig(workerNode, chainID, data, currentConfig, creEnv.CapabilityConfigs[cre.WriteSolanaCapability])
		if uErr != nil {
			return nil, errors.Wrapf(uErr, "failed to update node config for node index %d", workerNode.Index)
		}
		don.NodeSets().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = *updatedConfig
	}

	fullName := "write_solana_devnet@1.0.0"
	splitName := strings.Split(fullName, "@")

	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   splitName[0],
			Version:        splitName[1],
			CapabilityType: 3, // TARGET
			ResponseType:   1, // OBSERVATION_IDENTICAL
		},
		Config: &capabilitiespb.CapabilityConfig{},
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

func deployForwarder(testLogger zerolog.Logger, creEnv *cre.Environment, solChain *solana.Blockchain) (*string, *string, error) {
	memoryDatastore := datastore.NewMemoryDataStore()
	// load all existing addresses into memory datastore
	mergeErr := memoryDatastore.Merge(creEnv.CldfEnvironment.DataStore)
	if mergeErr != nil {
		return nil, nil, fmt.Errorf("failed to merge existing datastore into memory datastore: %w", mergeErr)
	}

	populateContracts := map[string]datastore.ContractType{
		deployment.KeystoneForwarderProgramName: ks_sol.ForwarderContract,
	}
	version := semver.MustParse(creEnv.ContractVersions[ks_sol.ForwarderContract.String()])

	// Forwarder for solana is predeployed on chain spin-up. We jus need to add it to memory datastore here
	errp := memory.PopulateDatastore(memoryDatastore.AddressRefStore, populateContracts,
		version, ks_sol.DefaultForwarderQualifier, solChain.ChainSelector())
	if errp != nil {
		return nil, nil, errors.Wrap(errp, "failed to populate datastore with predeployed contracts")
	}
	out, err := operations.ExecuteSequence(
		creEnv.CldfEnvironment.OperationsBundle,
		ks_sol_seq.DeployForwarderSeq,
		ks_sol_op.Deps{
			Env:       *creEnv.CldfEnvironment,
			Chain:     creEnv.CldfEnvironment.BlockChains.SolanaChains()[solChain.ChainSelector()],
			Datastore: memoryDatastore.Seal(),
		},
		ks_sol_seq.DeployForwarderSeqInput{
			ChainSel:     solChain.ChainSelector(),
			ProgramName:  deployment.KeystoneForwarderProgramName,
			Qualifier:    ks_sol.DefaultForwarderQualifier,
			ContractType: ks_sol.ForwarderContract,
			Version:      version,
		},
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to deploy sol forwarder")
	}

	err = memoryDatastore.AddressRefStore.Add(datastore.AddressRef{
		Address:       out.Output.State.String(),
		ChainSelector: solChain.ChainSelector(),
		Version:       semver.MustParse(creEnv.ContractVersions[ks_sol.ForwarderState.String()]),
		Qualifier:     ks_sol.DefaultForwarderQualifier,
		Type:          ks_sol.ForwarderState,
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to add address to the datastore for Solana Forwarder state")
	}

	testLogger.Info().Msgf("Deployed Forwarder %s contract on Solana chain chain %d programID: %s state: %s", creEnv.ContractVersions[ks_sol.ForwarderContract.String()], solChain.ChainSelector(), out.Output.ProgramID.String(), out.Output.State.String())

	creEnv.CldfEnvironment.DataStore = memoryDatastore.Seal()

	return ptr.Ptr(out.Output.ProgramID.String()), ptr.Ptr(out.Output.State.String()), nil
}

func updateNodeConfig(workerNode *cre.NodeMetadata, chainID string, data solanaInput, currentConfig string, capabilityConfig cre.CapabilityConfig) (*string, error) {
	key, ok := workerNode.Keys.Solana[chainID]
	if !ok {
		return nil, errors.Errorf("missing Solana key for chainID %s on node index %d", chainID, workerNode.Index)
	}
	data.FromAddress = key.PublicAddress

	mergedConfig := envconfig.ResolveCapabilityConfigForDON(
		cre.WriteSolanaCapability,
		capabilityConfig.Config,
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

	return ptr.Ptr(string(stringifiedConfig)), nil
}

type solanaInput struct {
	ChainSelector    uint64
	FromAddress      solanago.PublicKey
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
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	solChainsWithForwarder := make(map[uint64]struct{})
	solForwarders := creEnv.CldfEnvironment.DataStore.Addresses().Filter(datastore.AddressRefByQualifier(ks_sol.DefaultForwarderQualifier))
	for _, forwarder := range solForwarders {
		solChainsWithForwarder[forwarder.ChainSelector] = struct{}{}
	}

	// configure Solana forwarder only if we have some
	if len(solChainsWithForwarder) > 0 {
		cs := commonchangeset.Configure(ks_sol.ConfigureForwarders{},
			&ks_sol.ConfigureForwarderRequest{
				WFDonName:        don.Name,
				WFNodeIDs:        don.KeystoneDONConfig().NodeIDs,
				RegistryChainSel: creEnv.RegistryChainSelector,
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

	return nil
}
