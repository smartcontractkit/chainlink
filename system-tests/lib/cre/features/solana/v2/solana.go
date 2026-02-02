package v2

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"text/template"
	"time"

	"cosmossdk.io/errors"
	"dario.cat/mergo"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	solchain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/solana"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	flag           = cre.SolanaCapability
	configTemplate = `{
		"creForwarderAddress":"{{.CREForwarderAddress}}",
		"creForwarderState":"{{.CREForwarderState}}",
		"transmitter":"{{.NodeAddress}}",
		"isLocal":{{.IsLocal}},
		"chainId":"{{.ChainID}}",
		"network":"{{.Network}}"
	}`
	registrationRefresh = 20 * time.Second
	registrationExpiry  = 60 * time.Second
	deltaStage          = 14*time.Second + 2*time.Second // finalization time + 2 seconds delta
	requestTimeout      = 30 * time.Second
)

type SolChain interface {
	SolChainID() string
}

type Solana struct{}

func (s *Solana) Flag() cre.CapabilityFlag {
	return flag
}

func (s *Solana) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	// 1. Deploy forwarders to solana blockchains
	solChain := extractSolanaFromEnv(creEnv)
	programID, state, fErr := solana.DeployForwarder(testLogger, creEnv, solChain)
	if fErr != nil {
		return nil, errors.Wrapf(fErr, "failed to deploy forwarder for solana")
	}
	input := solana.SolanaInput{
		ForwarderAddress: *programID,
		ForwarderState:   *state,
	}
	// 2. Patch nodes TOML config to include workflow From Address
	cfgErr := updateNodeConfigs(creEnv, don, input, solChain.ChainSelector())
	if cfgErr != nil {
		return nil, errors.Wrapf(cfgErr, "failed to update node configs for solana")
	}

	// 3. Register Solana capability & its methods with Keystone
	capabilities, capErr := registerSolanaCapability(solChain.ChainSelector(), don)
	if capErr != nil {
		return nil, errors.Wrapf(capErr, "failed to register solana capability")
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

func (s *Solana) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	// 1. Deploy & Configure OCR3 Contracts (once solana consensus reads are supported)
	// 2. Create & Approve Solana Standard capability jobs for the DON
	jobErr := createJobs(ctx, don, dons, creEnv)
	if jobErr != nil {
		return errors.Wrapf(jobErr, "failed to create job for solana chain standard capability")
	}

	// 3. Configure Forwarders
	consensusDons := dons.DonsWithFlags(cre.ConsensusCapability, cre.ConsensusCapabilityV2)
	for _, don := range consensusDons {
		testLogger.Info().Msg(fmt.Sprintf("configure forwarder for: %s", don.Name))
		for _, n := range don.Nodes {
			testLogger.Info().Msg(fmt.Sprintf("solana keys: %s", n.Keys.OCR2BundleIDs["solana"]))
		}
		err := solana.ConfigureForwarders(ctx, testLogger, don, dons, creEnv)
		if err != nil {
			return err
		}
	}

	return nil
}

// post env
func createJobs(
	ctx context.Context,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	specs := make(map[string][]string)
	solChain := extractSolanaFromEnv(creEnv)
	capabilityConfig, ok := creEnv.CapabilityConfigs[flag]
	if !ok {
		return fmt.Errorf("%s config not found in capabilities config: %v", flag, creEnv.CapabilityConfigs)
	}

	command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryPath, creEnv.Provider)
	if cErr != nil {
		return errors.Wrap(cErr, "failed to get command for cron capability")
	}

	var nodeSet cre.NodeSetWithCapabilityConfigs
	for _, ns := range dons.AsNodeSetWithChainCapabilities() {
		if ns.GetName() == don.Name {
			nodeSet = ns
			break
		}
	}
	if nodeSet == nil {
		return fmt.Errorf("could not find node set for Don named '%s'", don.Name)
	}

	// propose bootstrap job once consensus reads are enabled
	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}

	chainID, chErr := chainselectors.SolanaChainIdFromSelector(solChain.ChainSelector())
	if chErr != nil {
		return errors.Wrapf(chErr, "failed to get Solana chain ID from selector %d", solChain.ChainSelector())
	}

	for _, workerNode := range workerNodes {
		key, ok := workerNode.Keys.Solana[chainID]
		if !ok {
			return fmt.Errorf("failed to get solana key (chainID %d, node index %d)", chainID, workerNode.Index)
		}

		version := creEnv.ContractVersions[ks_sol.ForwarderContract.String()]

		creForwarderKey := datastore.NewAddressRefKey(
			solChain.ChainSelector(),
			datastore.ContractType(ks_sol.ForwarderContract),
			version,
			ks_sol.DefaultForwarderQualifier,
		)
		creForwarderStateKey := datastore.NewAddressRefKey(
			solChain.ChainSelector(),
			datastore.ContractType(ks_sol.ForwarderState),
			version,
			ks_sol.DefaultForwarderQualifier,
		)
		creForwarderAddress, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(creForwarderKey)
		if err != nil {
			return errors.Wrap(err, "failed to get CRE Forwarder address")
		}
		creForwarderStateAddress, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(creForwarderStateKey)
		if err != nil {
			return errors.Wrap(err, "failed to get CRE Forwarder State address")
		}

		nodeAddress := key.PublicAddress.String()
		tmpl, err := template.New("solConfig").Parse(configTemplate)
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s config template", flag)
		}

		solChainID, err := solChain.SolClient.GetGenesisHash(ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to get sol genesis hash")
		}
		runtimeFallbacks := map[string]any{
			"CREForwarderAddress": creForwarderAddress.Address,
			"CREForwarderState":   creForwarderStateAddress.Address,
			"NodeAddress":         nodeAddress,
			"IsLocal":             true,
			"Network":             "solana",
			"ChainID":             solChainID.String(),
		}

		templateData, aErr := credon.ApplyRuntimeValues(capabilityConfig.Config, runtimeFallbacks)
		if aErr != nil {
			return errors.Wrap(aErr, "failed to apply runtime values")
		}

		var configBuffer bytes.Buffer
		if err := tmpl.Execute(&configBuffer, templateData); err != nil {
			return errors.Wrapf(err, "failed to execute %s config template", flag)
		}

		configStr := configBuffer.String()
		if err := credon.ValidateTemplateSubstitution(configStr, flag); err != nil {
			return errors.Wrapf(err, "%s template validation failed", flag)
		}

		workerInput := cre_jobs.ProposeJobSpecInput{
			Domain:      offchain.ProductLabel,
			Environment: cre.EnvironmentName,
			DONName:     don.Name,
			JobName:     fmt.Sprintf("sol-v2-worker-%s", chainID),
			ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: don.Name},
				{Key: "p2p_id", Value: workerNode.Keys.PeerID()}, // required since each node requires a different config (it contains its own from address)
			},
			Template: job_types.Solana,
			Inputs: job_types.JobSpecInput{
				"command": command,
				"config":  configStr,
			},
		}

		workerVerErr := cre_jobs.ProposeJobSpec{}.VerifyPreconditions(*creEnv.CldfEnvironment, workerInput)
		if workerVerErr != nil {
			return fmt.Errorf("precondition verification failed for Solana v2 worker job: %w", workerVerErr)
		}

		workerReport, workerErr := cre_jobs.ProposeJobSpec{}.Apply(*creEnv.CldfEnvironment, workerInput)
		if workerErr != nil {
			return fmt.Errorf("failed to propose Solana v2 worker job spec: %w", workerErr)
		}

		for _, r := range workerReport.Reports {
			out, ok := r.Output.(cre_jobs_ops.ProposeStandardCapabilityJobOutput)
			if !ok {
				return fmt.Errorf("unable to cast to ProposeStandardCapabilityJobOutput, actual type: %T", r.Output)
			}
			mErr := mergo.Merge(&specs, out.Specs, mergo.WithAppendSlice)
			if mErr != nil {
				return fmt.Errorf("failed to merge worker job specs: %w", mErr)
			}
		}
	}

	approveErr := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs)
	if approveErr != nil {
		return fmt.Errorf("failed to approve Solana v2 jobs: %w", approveErr)
	}

	return nil
}

// pre env
func registerSolanaCapability(selector uint64, don *cre.DonMetadata) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	var caps []keystone_changeset.DONCapabilityWithConfig
	methodConfigs, err := getMethodConfigs(don)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get method configs")
	}
	caps = append(caps, keystone_changeset.DONCapabilityWithConfig{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName: "solana" + ":ChainSelector:" + strconv.FormatUint(selector, 10),
			Version:      "1.0.0",
		},
		Config: &capabilitiespb.CapabilityConfig{
			MethodConfigs: methodConfigs,
		},
	})

	return caps, nil
}

func getMethodConfigs(don *cre.DonMetadata) (map[string]*capabilitiespb.CapabilityMethodConfig, error) {
	methodConfigs := make(map[string]*capabilitiespb.CapabilityMethodConfig)

	methodConfigs["WriteReport"] = writeReportActionConfig()

	triggerConfig, err := logTriggerConfig(don)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get LogTrigger config")
	}
	methodConfigs["LogTrigger"] = triggerConfig

	return methodConfigs, nil
}

func logTriggerConfig(don *cre.DonMetadata) (*capabilitiespb.CapabilityMethodConfig, error) {
	faultyNodes, faultyErr := don.NodeSets().MaxFaultyNodes()
	if faultyErr != nil {
		return nil, errors.Wrap(faultyErr, "failed to get faulty nodes")
	}

	return &capabilitiespb.CapabilityMethodConfig{
		RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh:     durationpb.New(registrationRefresh),
				RegistrationExpiry:      durationpb.New(registrationExpiry),
				MinResponsesToAggregate: faultyNodes + 1,
				MessageExpiry:           durationpb.New(2 * registrationExpiry),
				MaxBatchSize:            25,
				BatchCollectionPeriod:   durationpb.New(200 * time.Millisecond),
			},
		},
	}, nil
}

func writeReportActionConfig() *capabilitiespb.CapabilityMethodConfig {
	return &capabilitiespb.CapabilityMethodConfig{
		RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
			RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
				TransmissionSchedule:      capabilitiespb.TransmissionSchedule_OneAtATime,
				DeltaStage:                durationpb.New(deltaStage),
				RequestTimeout:            durationpb.New(requestTimeout),
				ServerMaxParallelRequests: 10,
				RequestHasherType:         capabilitiespb.RequestHasherType_WriteReportExcludeSignatures,
			},
		},
	}
}

func updateNodeConfigs(creEnv *cre.Environment, don *cre.DonMetadata, data solana.SolanaInput, selector uint64) error {
	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}
	chainID, chErr := chainselectors.SolanaChainIdFromSelector(selector)
	if chErr != nil {
		return chErr
	}

	for _, workerNode := range workerNodes {
		currentConfig := don.NodeSets().NodeSpecs[workerNode.Index].Node.TestConfigOverrides
		updatedConfig, updErr := solana.UpdateNodeConfig(workerNode, chainID, data, currentConfig, creEnv.CapabilityConfigs[cre.WriteSolanaCapability])
		if updErr != nil {
			return errors.Wrapf(updErr, "failed to update node config for node index %d", workerNode.Index)
		}
		don.NodeSets().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = *updatedConfig
	}

	return nil
}

func updateNodeConfig(workerNode *cre.NodeMetadata, nodeSet *cre.NodeSet, currentConfig string, selector uint64) (*string, error) {
	var typedConfig corechainlink.Config
	unmarshallErr := toml.Unmarshal([]byte(currentConfig), &typedConfig)
	if unmarshallErr != nil {
		return nil, errors.Wrapf(unmarshallErr, "failed to unmarshal config for node index %d", workerNode.Index)
	}

	// local CRE supports only 1 solana chain per environment
	if len(typedConfig.Solana) != 1 {
		return nil, fmt.Errorf("unexpected length of solana chain configs, expected 1, got %d", len(typedConfig.Solana))
	}

	skip := true
	typedConfig.Solana[0].Chain.SkipPreflight = &skip
	t, _ := config.NewDuration(time.Second * 30)
	typedConfig.Solana[0].Chain.TxRetentionTimeout = &t

	stringifiedConfig, mErr := toml.Marshal(typedConfig)
	if mErr != nil {
		return nil, errors.Wrapf(mErr, "failed to marshal config for node index %d", workerNode.Index)
	}
	fmt.Println("solcfg", string(stringifiedConfig))
	return ptr.Ptr(string(stringifiedConfig)), nil
}

func findSolFromAddress(workerNode *cre.NodeMetadata, selector uint64) (solanago.PublicKey, error) {
	chainID, chErr := chainselectors.SolanaChainIdFromSelector(selector)
	if chErr != nil {
		return solanago.PublicKey{}, errors.Wrapf(chErr, "failed to get Solana chain ID from selector %d", selector)
	}

	key, ok := workerNode.Keys.Solana[chainID]
	if !ok {
		return solanago.PublicKey{}, fmt.Errorf("missing Solana key for chainID %s on node index %d", chainID, workerNode.Index)
	}

	return key.PublicAddress, nil

}

func extractSolanaFromEnv(creEnv *cre.Environment) *solchain.Blockchain {
	var solChain *solchain.Blockchain
	for _, bcOut := range creEnv.Blockchains {
		if bcOut.IsFamily(chainselectors.FamilySolana) {
			solChain = bcOut.(*solchain.Blockchain)
			break
		}
	}

	return solChain
}
