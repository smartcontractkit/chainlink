package vault

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/Masterminds/semver/v3"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/durationpb"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/ocr3_1"
	depcontracts "github.com/smartcontractkit/chainlink/deployment/cre/ocr3/ocr3_1/changeset/operations/contracts"
	coretoml "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	vaultprotos "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	creseq "github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/sequences"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
)

const flag = cre.VaultCapability

const (
	ContractQualifier = "vault"
)

type Vault struct{}

type runtimeConfig struct {
	Auth0 *cre.GatewayServiceAuth0Config `json:"auth0"`
}

func (o *Vault) Flag() cre.CapabilityFlag {
	return flag
}

func (o *Vault) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	auth0Config, cfgErr := resolveRuntimeConfig(don.MustNodeSet())
	if cfgErr != nil {
		return nil, errors.Wrap(cfgErr, "failed to resolve vault runtime config")
	}

	// use registry chain, because that is the chain we used when generating gateway connector part of node config (check below)
	registryChainID, chErr := chainselectors.ChainIdFromSelector(creEnv.RegistryChainSelector)
	if chErr != nil {
		return nil, errors.Wrapf(chErr, "failed to get chain ID from selector %d", creEnv.RegistryChainSelector)
	}

	// add 'vault' handler to gateway config
	// add gateway connector to to node TOML config, so that node can route vault requests to the gateway
	hErr := topology.AddGatewayHandlers(*don, []string{pkg.GatewayHandlerTypeVault})
	if hErr != nil {
		return nil, errors.Wrapf(hErr, "failed to add gateway handlers to gateway config for don %s ", don.Name)
	}
	if auth0Config.Auth0 != nil {
		if err := applyGatewayAuth0Config(topology, don.Name, auth0Config.Auth0); err != nil {
			return nil, errors.Wrapf(err, "failed to apply auth0 gateway config for don %s", don.Name)
		}
	}

	cErr := don.ConfigureForGatewayAccess(registryChainID, *topology.GatewayConnectors)
	if cErr != nil {
		return nil, errors.Wrapf(cErr, "failed to add gateway connectors to node's TOML config in for don %s", don.Name)
	}

	workflowRegistryAddress := contracts.MustGetAddressFromDataStore(creEnv.CldfEnvironment.DataStore, creEnv.RegistryChainSelector, keystone_changeset.WorkflowRegistry.String(), creEnv.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")

	donsToConfigure := []*cre.DonMetadata{don}
	workflowDONs, wfErr := topology.DonsMetadata.WorkflowDONs()
	if wfErr == nil {
		for _, workflowDON := range workflowDONs {
			if workflowDON.ID == don.ID {
				continue
			}
			donsToConfigure = append(donsToConfigure, workflowDON)
		}
	}

	for _, donToConfigure := range donsToConfigure {
		if err := configureWorkersNodeConfig(donToConfigure, registryChainID, common.HexToAddress(workflowRegistryAddress), creEnv.ContractVersions[keystone_changeset.WorkflowRegistry.String()]); err != nil {
			return nil, err
		}
	}

	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "vault",
			Version:        "1.0.0",
			CapabilityType: 1, // ACTION
		},
		Config: &capabilitiespb.CapabilityConfig{
			LocalOnly:     don.HasOnlyLocalCapabilities(),
			MethodConfigs: vaultMethodConfigs(),
		},
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

func updateNodeConfig(workerNode *cre.NodeMetadata, currentConfig string, registryChainID uint64, workflowRegistryAddress common.Address, wfRegVersion *semver.Version) (*string, error) {
	var typedConfig corechainlink.Config
	unmarshallErr := toml.Unmarshal([]byte(currentConfig), &typedConfig)
	if unmarshallErr != nil {
		return nil, errors.Wrapf(unmarshallErr, "failed to unmarshal config for node index %d", workerNode.Index)
	}

	// enable workflow registry syncer
	typedConfig.Capabilities.WorkflowRegistry = coretoml.WorkflowRegistry{
		Address:         new(workflowRegistryAddress.Hex()),
		NetworkID:       new("evm"),
		ChainID:         new(strconv.FormatUint(registryChainID, 10)),
		SyncStrategy:    new("reconciliation"),
		ContractVersion: new(wfRegVersion.String()),
	}
	typedConfig.CRE.Linking = &coretoml.LinkingConfig{
		URL:        new(strings.TrimPrefix(framework.HostDockerInternal(), "http://") + ":18124"),
		TLSEnabled: new(false),
	}

	stringifiedConfig, mErr := toml.Marshal(typedConfig)
	if mErr != nil {
		return nil, errors.Wrapf(mErr, "failed to marshal config for node index %d", workerNode.Index)
	}

	return new(string(stringifiedConfig)), nil
}

func configureWorkersNodeConfig(don *cre.DonMetadata, registryChainID uint64, workflowRegistryAddress common.Address, wfRegVersion *semver.Version) error {
	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrapf(wErr, "failed to find worker nodes for don %s", don.Name)
	}

	for _, workerNode := range workerNodes {
		currentConfig := don.MustNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides
		updatedConfig, uErr := updateNodeConfig(workerNode, currentConfig, registryChainID, workflowRegistryAddress, wfRegVersion)
		if uErr != nil {
			return errors.Wrapf(uErr, "failed to update node config for don %s node index %d", don.Name, workerNode.Index)
		}
		don.MustNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = *updatedConfig
	}

	return nil
}

func (o *Vault) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	vaultOCR3Addr, vaultDKGOCR3Addr, err := deployVaultContracts(testLogger, ContractQualifier, creEnv.RegistryChainSelector, creEnv.CldfEnvironment, creEnv.ContractVersions)
	if err != nil {
		return fmt.Errorf("failed to deploy Vault OCR3 contract %w", err)
	}

	jobErr := createJobs(
		ctx,
		creEnv,
		don,
		dons,
	)
	if jobErr != nil {
		return fmt.Errorf("failed to create OCR3 jobs: %w", jobErr)
	}

	ocr3Config := contracts.DefaultOCR3_1Config(don.WorkersCount())
	dkgOffchainConfig := &ocr3_1.DKGOffchainConfig{
		T: 1,
	}

	workers, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}

	for _, workerNode := range workers {
		pubKey := hex.EncodeToString(workerNode.Keys.DKGKey.PubKey)
		dkgOffchainConfig.DealerPublicKeys = append(dkgOffchainConfig.DealerPublicKeys, pubKey)
		dkgOffchainConfig.RecipientPublicKeys = append(dkgOffchainConfig.RecipientPublicKeys, pubKey)
	}
	ocr3Config.DKGOffchainConfig = dkgOffchainConfig

	chain, ok := creEnv.CldfEnvironment.BlockChains.EVMChains()[creEnv.RegistryChainSelector]
	if !ok {
		return fmt.Errorf("chain with selector %d not found in environment", creEnv.RegistryChainSelector)
	}

	strategy, err := strategies.CreateStrategy(
		chain,
		*creEnv.CldfEnvironment,
		nil,
		nil,
		*vaultDKGOCR3Addr,
		"PostEnvStartup - Configure OCR3 Contract - Vault DKG",
	)
	if err != nil {
		return fmt.Errorf("failed to create strategy: %w", err)
	}

	_, err = operations.ExecuteOperation(
		creEnv.CldfEnvironment.OperationsBundle,
		ks_contracts_op.ConfigureDKGOp,
		ks_contracts_op.ConfigureDKGOpDeps{
			Env:      creEnv.CldfEnvironment,
			Strategy: strategy,
		},
		ks_contracts_op.ConfigureDKGOpInput{
			ContractAddress: vaultDKGOCR3Addr,
			ChainSelector:   creEnv.RegistryChainSelector,
			DON:             don.KeystoneDONConfig(),
			Config:          ocr3Config,
			DryRun:          false,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to configure DKG OCR3 contract")
	}

	cfgb, cErr := reportingPluginConfigOverride(vaultDKGOCR3Addr, creEnv)
	if cErr != nil {
		return fmt.Errorf("failed to create Vault reporting plugin config override: %w", cErr)
	}

	strategy, err = strategies.CreateStrategy(
		chain,
		*creEnv.CldfEnvironment,
		nil,
		nil,
		*vaultOCR3Addr,
		"PostEnvStartup - Configure OCR3 Contract - Vault",
	)
	if err != nil {
		return fmt.Errorf("failed to create strategy: %w", err)
	}

	_, err = operations.ExecuteOperation(
		creEnv.CldfEnvironment.OperationsBundle,
		depcontracts.ConfigureOCR3_1,
		depcontracts.ConfigureOCR3_1Deps{
			Env:      creEnv.CldfEnvironment,
			Strategy: strategy,
		},
		depcontracts.ConfigureOCR3_1Input{
			ContractAddress:               vaultOCR3Addr,
			ChainSelector:                 creEnv.RegistryChainSelector,
			DON:                           don.KeystoneDONConfig(),
			Config:                        ocr3Config,
			DryRun:                        false,
			ReportingPluginConfigOverride: cfgb,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to configure Vault OCR3 contract")
	}

	return nil
}

func createJobs(
	ctx context.Context,
	creEnv *cre.Environment,
	don *cre.Don,
	dons *cre.Dons,
) error {
	auth0Config := &runtimeConfig{}
	if capConfig, ok := don.GetCapabilityConfig(flag); ok {
		var err error
		auth0Config, err = decodeRuntimeConfig(capConfig.Values)
		if err != nil {
			return fmt.Errorf("failed to resolve vault runtime config: %w", err)
		}
	}

	bootstrap, isBootstrap := dons.Bootstrap()
	if !isBootstrap {
		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}

	specs := make(map[string][]string)

	_, ocrPeeringCfg, err := cre.PeeringCfgs(bootstrap)
	if err != nil {
		return errors.Wrap(err, "failed to get peering configs")
	}

	workerInput := cre_jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: cre.EnvironmentName,
		DONName:     don.Name,
		JobName:     "vault-worker",
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: don.Name},
		},
		Template: job_types.OCR3,
		Inputs: job_types.JobSpecInput{
			"chainSelectorEVM":     creEnv.RegistryChainSelector,
			"contractQualifier":    ContractQualifier + "_plugin",
			"dkgContractQualifier": ContractQualifier + "_dkg",
			"templateName":         "worker-vault",
			"bootstrapperOCR3Urls": []string{ocrPeeringCfg.OCRBootstraperPeerID + "@" + ocrPeeringCfg.OCRBootstraperHost + ":" + strconv.Itoa(ocrPeeringCfg.Port)},
		},
	}
	if auth0Config.Auth0 != nil {
		workerInput.Inputs["auth0"] = auth0Config.Auth0
	}

	workerVerErr := cre_jobs.ProposeJobSpec{}.VerifyPreconditions(*creEnv.CldfEnvironment, workerInput)
	if workerVerErr != nil {
		return fmt.Errorf("precondition verification failed for Vault worker job: %w", workerVerErr)
	}

	workerReport, workerErr := cre_jobs.ProposeJobSpec{}.Apply(*creEnv.CldfEnvironment, workerInput)
	if workerErr != nil {
		return fmt.Errorf("failed to propose Vault worker job spec: %w", workerErr)
	}

	for _, r := range workerReport.Reports {
		out, ok := r.Output.(cre_jobs_ops.ProposeOCR3JobOutput)
		if !ok {
			return fmt.Errorf("unable to cast to ProposeOCR3JobOutput, actual type: %T", r.Output)
		}
		mErr := mergo.Merge(&specs, out.Specs, mergo.WithAppendSlice)
		if mErr != nil {
			return fmt.Errorf("failed to merge worker job specs: %w", mErr)
		}
	}

	approveErr := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs)
	if approveErr != nil {
		return fmt.Errorf("failed to approve Vault jobs: %w", approveErr)
	}

	return nil
}

func resolveRuntimeConfig(nodeSet *cre.NodeSet) (*runtimeConfig, error) {
	if nodeSet == nil {
		return &runtimeConfig{}, nil
	}

	capConfig, ok := nodeSet.GetCapabilityConfig(flag)
	if !ok || len(capConfig.Values) == 0 {
		return &runtimeConfig{}, nil
	}

	return decodeRuntimeConfig(capConfig.Values)
}

func decodeRuntimeConfig(values map[string]any) (*runtimeConfig, error) {
	if len(values) == 0 {
		return &runtimeConfig{}, nil
	}

	b, err := json.Marshal(values)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vault capability values: %w", err)
	}

	cfg := &runtimeConfig{}
	if err := json.Unmarshal(b, cfg); err != nil {
		return nil, fmt.Errorf("failed to decode vault capability values: %w", err)
	}

	return cfg, nil
}

func applyGatewayAuth0Config(topology *cre.Topology, donName string, auth0 *cre.GatewayServiceAuth0Config) error {
	if topology == nil || auth0 == nil {
		return nil
	}

	for idx := range topology.GatewayServiceConfigs {
		svc := &topology.GatewayServiceConfigs[idx]
		if svc.ServiceName != pkg.ServiceNameVault || !slices.Contains(svc.DONs, donName) {
			continue
		}
		if svc.Auth0 != nil && (svc.Auth0.IssuerURL != auth0.IssuerURL || svc.Auth0.Audience != auth0.Audience) {
			return fmt.Errorf("vault gateway service %q already has conflicting auth0 config", svc.ServiceName)
		}

		svc.Auth0 = &cre.GatewayServiceAuth0Config{
			IssuerURL: auth0.IssuerURL,
			Audience:  auth0.Audience,
		}
		return nil
	}

	return fmt.Errorf("vault gateway service config not found for DON %s", donName)
}

func deployVaultContracts(testLogger zerolog.Logger, qualifier string, registryChainSelector uint64, env *cldf.Environment, contractVersions map[cre.ContractType]*semver.Version) (*common.Address, *common.Address, error) {
	memoryDatastore, mErr := contracts.NewDataStoreFromExisting(env.DataStore)
	if mErr != nil {
		return nil, nil, fmt.Errorf("failed to create memory datastore: %w", mErr)
	}

	report, err := operations.ExecuteSequence(
		env.OperationsBundle,
		creseq.DeployVault,
		creseq.DeployVaultDeps{
			Env: env,
		},
		creseq.DeployVaultInput{
			ChainSelector: registryChainSelector,
			Qualifier:     qualifier,
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy OCR3 contract '%s' on chain %d: %w", qualifier, registryChainSelector, err)
	}
	if err = memoryDatastore.Merge(report.Output.Datastore); err != nil {
		return nil, nil, fmt.Errorf("failed to merge datastore with OCR3 contract address for '%s' on chain %d: %w", qualifier, registryChainSelector, err)
	}

	vaultOCR3Addr := report.Output.PluginAddress
	testLogger.Info().Msgf("Deployed OCR3 %s (Vault) contract on chain %d at %s", contractVersions[keystone_changeset.OCR3Capability.String()], registryChainSelector, vaultOCR3Addr)
	vaultDKGOCR3Addr := report.Output.DKGAddress
	testLogger.Info().Msgf("Deployed OCR3 %s (DKG) contract on chain %d at %s", contractVersions[keystone_changeset.OCR3Capability.String()], registryChainSelector, vaultDKGOCR3Addr)

	env.DataStore = memoryDatastore.Seal()

	return new(common.HexToAddress(vaultOCR3Addr)), new(common.HexToAddress(vaultDKGOCR3Addr)), nil
}

func dkgReportingPluginConfig(don *cre.Don) (*dkgocrtypes.ReportingPluginConfig, error) {
	cfg := &dkgocrtypes.ReportingPluginConfig{
		T: 1,
	}

	workers, wErr := don.Workers()
	if wErr != nil {
		return nil, errors.Wrap(wErr, "failed to find worker nodes")
	}

	for _, workerNode := range workers {
		pubKey := workerNode.Keys.DKGKey.PubKey
		cfg.DealerPublicKeys = append(cfg.DealerPublicKeys, pubKey)
		cfg.RecipientPublicKeys = append(cfg.RecipientPublicKeys, pubKey)
	}

	return cfg, nil
}

func reportingPluginConfigOverride(vaultDKGOCR3Addr *common.Address, creEnv *cre.Environment) ([]byte, error) {
	client := creEnv.CldfEnvironment.BlockChains.EVMChains()[creEnv.RegistryChainSelector].Client
	dkgContract, err := ocr3_capability.NewOCR3Capability(*vaultDKGOCR3Addr, client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create OCR3 capability contract")
	}
	details, err := dkgContract.LatestConfigDetails(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest config details from OCR3 capability contract")
	}
	instanceID := string(dkgocrtypes.MakeInstanceID(dkgContract.Address(), details.ConfigDigest))
	cfg := vaultprotos.ReportingPluginConfig{
		DKGInstanceID: &instanceID,
	}
	cfgb, err := proto.Marshal(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal vault reporting plugin config")
	}

	return cfgb, nil
}

func vaultMethodConfigs() map[string]*capabilitiespb.CapabilityMethodConfig {
	return map[string]*capabilitiespb.CapabilityMethodConfig{
		vaultprotos.MethodGetSecrets: {
			RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
				RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
					RequestTimeout:            durationpb.New(2 * time.Minute),
					ServerMaxParallelRequests: 10,
					RequestHasherType:         capabilitiespb.RequestHasherType_Simple,
				},
			},
		},
	}
}
