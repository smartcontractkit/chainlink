package v2

import (
	"bytes"
	"context"
	"fmt"
	"maps"
	"runtime"
	"strconv"
	"text/template"
	"time"

	"dario.cat/mergo"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/durationpb"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	cre_sol "github.com/smartcontractkit/chainlink/deployment/cre/forwarder/solana"
	cre_sol_seq "github.com/smartcontractkit/chainlink/deployment/cre/forwarder/solana/sequence"
	cre_sol_op "github.com/smartcontractkit/chainlink/deployment/cre/forwarder/solana/sequence/operation"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/utils/solutils"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	solchain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/jobhelpers"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
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
	deltaStage     = 14*time.Second + 2*time.Second // finalization time + 2 seconds delta
	requestTimeout = 30 * time.Second
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
	_ *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	// 1. Deploy forwarders to solana blockchains
	solChain := extractSolanaFromEnv(creEnv)
	programID, state, fErr := deployForwarder(testLogger, creEnv, solChain)
	if fErr != nil {
		return nil, errors.Wrapf(fErr, "failed to deploy forwarder for solana")
	}
	input := input{
		ForwarderAddress: *programID,
		ForwarderState:   *state,
	}
	// 2. Patch nodes TOML config to include workflow From Address
	cfgErr := patchNodeTOML(creEnv, don, input, solChain.ChainSelector())
	if cfgErr != nil {
		return nil, errors.Wrapf(cfgErr, "failed to update node configs for solana")
	}

	// 3. Register Solana capability & its methods with Keystone
	capabilities := registerSolanaCapability(solChain.ChainSelector())
	capabilityToExtraSignerFamilies := make(map[string][]string, len(capabilities))
	for _, capability := range capabilities {
		capabilityToExtraSignerFamilies[capability.Capability.LabelledName] = []string{chainselectors.FamilySolana}
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig:         capabilities,
		CapabilityToExtraSignerFamilies: capabilityToExtraSignerFamilies,
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
		err := configureForwarders(ctx, testLogger, don, dons, creEnv)
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
	solChain := extractSolanaFromEnv(creEnv)

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
	config, resolveErr := cre.ResolveCapabilityConfig(nodeSet, flag, cre.CapabilityScope{})
	if resolveErr != nil {
		return errors.Wrap(resolveErr, "unable to find solana capability config")
	}

	command, cErr := standardcapability.GetCommand(config.BinaryName)
	if cErr != nil {
		return errors.Wrap(cErr, "failed to get command for cron capability")
	}

	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}

	chainID, chErr := chainselectors.SolanaChainIdFromSelector(solChain.ChainSelector())
	if chErr != nil {
		return errors.Wrapf(chErr, "failed to get Solana chain ID from selector %d", solChain.ChainSelector())
	}

	solChainID, err := solChain.GenesisHash(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to get sol genesis hash")
	}
	version := creEnv.ContractVersions[cre_sol.ForwarderContract.String()]
	creForwarderKey := datastore.NewAddressRefKey(
		solChain.ChainSelector(),
		cre_sol.ForwarderContract,
		version,
		cre_sol.DefaultForwarderQualifier,
	)
	creForwarderStateKey := datastore.NewAddressRefKey(
		solChain.ChainSelector(),
		cre_sol.ForwarderState,
		version,
		cre_sol.DefaultForwarderQualifier,
	)
	creForwarderAddress, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(creForwarderKey)
	if err != nil {
		return errors.Wrap(err, "failed to get CRE Forwarder address")
	}
	creForwarderStateAddress, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(creForwarderStateKey)
	if err != nil {
		return errors.Wrap(err, "failed to get CRE Forwarder State address")
	}
	tmpl, err := template.New("solConfig").Parse(configTemplate)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %s config template", flag)
	}

	results := make([]map[string][]string, len(workerNodes))
	group, groupCtx := errgroup.WithContext(ctx)
	group.SetLimit(min(len(workerNodes), runtime.GOMAXPROCS(0)))
	for i, workerNode := range workerNodes {
		group.Go(func() error {
			key, ok := workerNode.Keys.Solana[chainID]
			if !ok {
				return fmt.Errorf("failed to get solana key (chainID %s, node index %d)", chainID, workerNode.Index)
			}

			nodeAddress := key.PublicAddress.String()
			runtimeFallbacks := map[string]any{
				"CREForwarderAddress": creForwarderAddress.Address,
				"CREForwarderState":   creForwarderStateAddress.Address,
				"NodeAddress":         nodeAddress,
				"IsLocal":             true,
				"Network":             "solana",
				"ChainID":             solChainID,
			}

			templateData, aErr := credon.ApplyRuntimeValues(maps.Clone(config.Values), runtimeFallbacks)
			if aErr != nil {
				return errors.Wrap(aErr, "failed to apply runtime values")
			}

			var configBuffer bytes.Buffer
			if executeErr := tmpl.Execute(&configBuffer, templateData); executeErr != nil {
				return errors.Wrapf(executeErr, "failed to execute %s config template", flag)
			}

			configStr := configBuffer.String()
			if validateErr := credon.ValidateTemplateSubstitution(configStr, flag); validateErr != nil {
				return errors.Wrapf(validateErr, "%s template validation failed", flag)
			}

			workerInput := cre_jobs.ProposeJobSpecInput{
				Domain:      offchain.ProductLabel,
				Environment: cre.EnvironmentName,
				DONName:     don.Name,
				JobName:     "sol-v2-worker-" + chainID,
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

			specs := make(map[string][]string)
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

			select {
			case <-groupCtx.Done():
				return groupCtx.Err()
			default:
			}

			results[i] = specs
			return nil
		})
	}

	if wErr := group.Wait(); wErr != nil {
		return wErr
	}

	specs, mErr := jobhelpers.MergeSpecsByIndex(results)
	if mErr != nil {
		return mErr
	}

	approveErr := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs)
	if approveErr != nil {
		return fmt.Errorf("failed to approve Solana v2 jobs: %w", approveErr)
	}

	return nil
}

// pre env
func registerSolanaCapability(selector uint64) []keystone_changeset.DONCapabilityWithConfig {
	var caps []keystone_changeset.DONCapabilityWithConfig
	methodConfigs := getMethodConfigs()
	caps = append(caps, keystone_changeset.DONCapabilityWithConfig{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName: "solana" + ":ChainSelector:" + strconv.FormatUint(selector, 10),
			Version:      "1.0.0",
		},
		Config: &capabilitiespb.CapabilityConfig{
			MethodConfigs: methodConfigs,
		},
	})

	return caps
}

func getMethodConfigs() map[string]*capabilitiespb.CapabilityMethodConfig {
	methodConfigs := make(map[string]*capabilitiespb.CapabilityMethodConfig)

	methodConfigs["WriteReport"] = writeReportActionConfig()
	// PLEX-1828
	// PLEX-1918 Add the rest of solana methods here

	return methodConfigs
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

func patchNodeTOML(creEnv *cre.Environment, don *cre.DonMetadata, data input, selector uint64) error {
	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}
	chainID, chErr := chainselectors.SolanaChainIdFromSelector(selector)
	if chErr != nil {
		return chErr
	}

	for _, workerNode := range workerNodes {
		currentConfig := don.MustNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides
		updatedConfig, updErr := updateNodeConfig(workerNode, chainID, data, currentConfig, don.CapabilityConfigs[cre.SolanaCapability])
		if updErr != nil {
			return errors.Wrapf(updErr, "failed to update node config for node index %d", workerNode.Index)
		}
		don.MustNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = *updatedConfig
	}

	return nil
}

func deployForwarder(testLogger zerolog.Logger, creEnv *cre.Environment, solChain *solchain.Blockchain) (*string, *string, error) {
	memoryDatastore, err := contracts.NewDataStoreFromExisting(creEnv.CldfEnvironment.DataStore)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create memory datastore: %w", err)
	}

	version := creEnv.ContractVersions[cre_sol.ForwarderContract.String()]

	// Forwarder for solana is predeployed on chain spin-up. We just need to add it to memory datastore here
	err = memoryDatastore.Addresses().Add(datastore.AddressRef{
		Address:       solutils.GetProgramID(solutils.ProgKeystoneForwarder),
		ChainSelector: solChain.ChainSelector(),
		Type:          cre_sol.ForwarderContract,
		Version:       version,
		Qualifier:     cre_sol.DefaultForwarderQualifier,
	})
	if err != nil && !errors.Is(err, datastore.ErrAddressRefExists) {
		return nil, nil, fmt.Errorf("failed to add address to the datastore for Solana Forwarder: %w", err)
	}

	out, err := operations.ExecuteSequence(
		creEnv.CldfEnvironment.OperationsBundle,
		cre_sol_seq.DeployForwarderSeq,
		cre_sol_op.Deps{
			Env:       *creEnv.CldfEnvironment,
			Chain:     creEnv.CldfEnvironment.BlockChains.SolanaChains()[solChain.ChainSelector()],
			Datastore: memoryDatastore.Seal(),
		},
		cre_sol_seq.DeployForwarderSeqInput{
			ChainSel:     solChain.ChainSelector(),
			ProgramName:  solutils.ProgKeystoneForwarder,
			Qualifier:    cre_sol.DefaultForwarderQualifier,
			ContractType: cre_sol.ForwarderContract,
			Version:      version,
		},
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to deploy sol forwarder")
	}

	err = memoryDatastore.AddressRefStore.Add(datastore.AddressRef{
		Address:       out.Output.State.String(),
		ChainSelector: solChain.ChainSelector(),
		Version:       creEnv.ContractVersions[cre_sol.ForwarderState.String()],
		Qualifier:     cre_sol.DefaultForwarderQualifier,
		Type:          cre_sol.ForwarderState,
	})

	if err != nil && !errors.Is(err, datastore.ErrAddressRefExists) {
		return nil, nil, errors.Wrap(err, "failed to add address to the datastore for Solana Forwarder state")
	}

	testLogger.Info().Msgf("Deployed Forwarder %s contract on Solana chain chain %d programID: %s state: %s", creEnv.ContractVersions[cre_sol.ForwarderContract.String()], solChain.ChainSelector(), out.Output.ProgramID.String(), out.Output.State.String())

	creEnv.CldfEnvironment.DataStore = memoryDatastore.Seal()

	return ptr.Ptr(out.Output.ProgramID.String()), ptr.Ptr(out.Output.State.String()), nil
}

func updateNodeConfig(workerNode *cre.NodeMetadata, chainID string, data input, currentConfig string, capabilityConfig cre.CapabilityConfig) (*string, error) {
	key, ok := workerNode.Keys.Solana[chainID]
	if !ok {
		return nil, errors.Errorf("missing Solana key for chainID %s on node index %d", chainID, workerNode.Index)
	}
	data.FromAddress = key.PublicAddress

	runtimeValues := map[string]any{
		"FromAddress":      data.FromAddress.String(),
		"ForwarderAddress": data.ForwarderAddress,
		"ForwarderState":   data.ForwarderState,
	}

	var mErr error
	data.WorkflowConfig, mErr = credon.ApplyRuntimeValues(capabilityConfig.Values, runtimeValues)
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

	if err := credon.ValidateTemplateSubstitution(configStr, flag); err != nil {
		return nil, fmt.Errorf("%s template validation failed: %w\nRendered template: %s", flag, err, configStr)
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

type input struct {
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
		TxAcceptanceState = {{printf "%d" .TxAcceptanceState}}
		Local = {{.Local}}
	`

func configureForwarders(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	solChainsWithForwarder := make(map[uint64]struct{})
	solForwarders := creEnv.CldfEnvironment.DataStore.Addresses().Filter(datastore.AddressRefByQualifier(cre_sol.DefaultForwarderQualifier))
	for _, forwarder := range solForwarders {
		solChainsWithForwarder[forwarder.ChainSelector] = struct{}{}
	}

	// configure Solana forwarder only if we have some
	if len(solChainsWithForwarder) > 0 {
		cs := commonchangeset.Configure(cre_sol.ConfigureForwarders{},
			&cre_sol.ConfigureForwarderRequest{
				DON: forwarder.DonConfiguration{
					Name:    don.Name,
					F:       don.F,
					ID:      libc.MustSafeUint32FromUint64(don.ID),
					NodeIDs: don.KeystoneDONConfig().NodeIDs,
					Version: 1,
				},
				Chains:    solChainsWithForwarder,
				Qualifier: cre_sol.DefaultForwarderQualifier,
				Version:   "1.0.0",
			},
		)

		_, err := cs.Apply(*creEnv.CldfEnvironment)
		if err != nil {
			return errors.Wrap(err, "failed to configure Solana forwarders")
		}
	}

	return nil
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
