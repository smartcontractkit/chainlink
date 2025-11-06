package vault

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"dario.cat/mergo"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	coretoml "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	vaultprotos "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	cre_jobs_seq "github.com/smartcontractkit/chainlink/deployment/cre/jobs/sequences"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	creseq "github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/sequences"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	depcontracts "github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
)

const flag = cre.VaultCapability

const (
	ContractQualifier = "vault"
)

type Vault struct{}

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

	cErr := don.ConfigureForGatewayAccess(registryChainID, *topology.GatewayConnectors)
	if cErr != nil {
		return nil, errors.Wrapf(cErr, "failed to add gateway connectors to node's TOML config in for don %s", don.Name)
	}

	workflowRegistryAddress, wfRegTypeVersion, wfErr := contracts.FindAddressesForChain(
		creEnv.CldfEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate
		creEnv.RegistryChainSelector,
		keystone_changeset.WorkflowRegistry.String(),
	)
	if wfErr != nil {
		return nil, errors.Wrap(wfErr, "failed to find WorkflowRegistry address")
	}

	// enable workflow registry syncer in node's TOML config
	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return nil, errors.Wrap(wErr, "failed to find worker nodes")
	}

	for _, workerNode := range workerNodes {
		currentConfig := don.NodeSets().NodeSpecs[workerNode.Index].Node.TestConfigOverrides
		updatedConfig, uErr := updateNodeConfig(workerNode, currentConfig, registryChainID, workflowRegistryAddress, wfRegTypeVersion)
		if uErr != nil {
			return nil, errors.Wrapf(uErr, "failed to update node config for node index %d", workerNode.Index)
		}
		don.NodeSets().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = *updatedConfig
	}

	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "vault",
			Version:        "1.0.0",
			CapabilityType: 1, // ACTION
		},
		Config: &capabilitiespb.CapabilityConfig{},
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

func updateNodeConfig(workerNode *cre.NodeMetadata, currentConfig string, registryChainID uint64, workflowRegistryAddress common.Address, wfRegTypeVersion cldf.TypeAndVersion) (*string, error) {
	var typedConfig corechainlink.Config
	unmarshallErr := toml.Unmarshal([]byte(currentConfig), &typedConfig)
	if unmarshallErr != nil {
		return nil, errors.Wrapf(unmarshallErr, "failed to unmarshal config for node index %d", workerNode.Index)
	}

	// enable workflow registry syncer
	typedConfig.Capabilities.WorkflowRegistry = coretoml.WorkflowRegistry{
		Address:         ptr.Ptr(workflowRegistryAddress.Hex()),
		NetworkID:       ptr.Ptr("evm"),
		ChainID:         ptr.Ptr(strconv.FormatUint(registryChainID, 10)),
		SyncStrategy:    ptr.Ptr("reconciliation"),
		ContractVersion: ptr.Ptr(wfRegTypeVersion.Version.String()),
	}

	stringifiedConfig, mErr := toml.Marshal(typedConfig)
	if mErr != nil {
		return nil, errors.Wrapf(mErr, "failed to marshal config for node index %d", workerNode.Index)
	}

	return ptr.Ptr(string(stringifiedConfig)), nil
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

	ocr3Config, ocr3confErr := contracts.DefaultOCR3_1Config(don.WorkersCount())
	if ocr3confErr != nil {
		return fmt.Errorf("failed to get default OCR3 config: %w", ocr3confErr)
	}

	dkgConfig, dErr := dkgReportingPluginConfig(don)
	if dErr != nil {
		return fmt.Errorf("failed to create DKG reporting plugin config: %w", dErr)
	}

	_, err = operations.ExecuteOperation(
		creEnv.CldfEnvironment.OperationsBundle,
		ks_contracts_op.ConfigureDKGOp,
		ks_contracts_op.ConfigureDKGOpDeps{
			Env: creEnv.CldfEnvironment,
		},
		ks_contracts_op.ConfigureDKGOpInput{
			ContractAddress:       vaultDKGOCR3Addr,
			ChainSelector:         creEnv.RegistryChainSelector,
			DON:                   don.KeystoneDONConfig(),
			Config:                ocr3Config,
			DryRun:                false,
			ReportingPluginConfig: *dkgConfig,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to configure DKG OCR3 contract")
	}

	cfgb, cErr := reportingPluginConfigOverride(vaultDKGOCR3Addr, creEnv, dons)
	if cErr != nil {
		return fmt.Errorf("failed to create Vault reporting plugin config override: %w", cErr)
	}

	_, err = operations.ExecuteOperation(
		creEnv.CldfEnvironment.OperationsBundle,
		depcontracts.ConfigureOCR3_1,
		depcontracts.ConfigureOCR3_1Deps{
			Env: creEnv.CldfEnvironment,
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
	bootstrap, isBootstrap := dons.Bootstrap()
	if !isBootstrap {
		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}

	bootInput := cre_jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: cre.EnvironmentName,
		DONName:     bootstrap.DON.Name,
		JobName:     "vault-bootstrap",
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: bootstrap.DON.Name},
		},
		Template: job_types.BootstrapVault,
		Inputs: job_types.JobSpecInput{
			"chainSelector":           creEnv.RegistryChainSelector,
			"contractQualifierPrefix": ContractQualifier,
		},
	}

	bootVerErr := cre_jobs.ProposeJobSpec{}.VerifyPreconditions(*creEnv.CldfEnvironment, bootInput)
	if bootVerErr != nil {
		return fmt.Errorf("precondition verification failed for Vault bootstrap job: %w", bootVerErr)
	}

	bootReport, bootErr := cre_jobs.ProposeJobSpec{}.Apply(*creEnv.CldfEnvironment, bootInput)
	if bootErr != nil {
		return fmt.Errorf("failed to propose Vault bootstrap job spec: %w", bootErr)
	}

	specs := make(map[string][]string)
	for _, r := range bootReport.Reports {
		out, ok := r.Output.(cre_jobs_seq.ProposeVaultBootstrapJobsOutput)
		if !ok {
			return fmt.Errorf("unable to cast to ProposeVaultBootstrapJobsOutput, actual type: %T", r.Output)
		}
		mErr := mergo.Merge(&specs, out.Specs, mergo.WithAppendSlice)
		if mErr != nil {
			return fmt.Errorf("failed to merge bootstrap job specs: %w", mErr)
		}
	}

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

func deployVaultContracts(testLogger zerolog.Logger, qualifier string, registryChainSelector uint64, env *cldf.Environment, contractVersions map[string]string) (*common.Address, *common.Address, error) {
	memoryDatastore := datastore.NewMemoryDataStore()

	// load all existing addresses into memory datastore
	mergeErr := memoryDatastore.Merge(env.DataStore)
	if mergeErr != nil {
		return nil, nil, fmt.Errorf("failed to merge existing datastore into memory datastore: %w", mergeErr)
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

	return ptr.Ptr(common.HexToAddress(vaultOCR3Addr)), ptr.Ptr(common.HexToAddress(vaultDKGOCR3Addr)), nil
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

func reportingPluginConfigOverride(vaultDKGOCR3Addr *common.Address, creEnv *cre.Environment, dons *cre.Dons) ([]byte, error) {
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

func EncryptSecret(secret, masterPublicKeyStr string) (string, error) {
	masterPublicKey := tdh2easy.PublicKey{}
	masterPublicKeyBytes, err := hex.DecodeString(masterPublicKeyStr)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode master public key")
	}
	err = masterPublicKey.Unmarshal(masterPublicKeyBytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal master public key")
	}
	cipher, err := tdh2easy.Encrypt(&masterPublicKey, []byte(secret))
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt secret")
	}
	cipherBytes, err := cipher.Marshal()
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal encrypted secrets to bytes")
	}
	return hex.EncodeToString(cipherBytes), nil
}
