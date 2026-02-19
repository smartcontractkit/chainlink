package vault

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"dario.cat/mergo"
	"github.com/Masterminds/semver/v3"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/ocr3_1"
	coretoml "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr/capregconfig"

	vaultprotos "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

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

	workflowRegistryAddress := contracts.MustGetAddressFromDataStore(creEnv.CldfEnvironment.DataStore, creEnv.RegistryChainSelector, keystone_changeset.WorkflowRegistry.String(), creEnv.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")

	// enable workflow registry syncer in node's TOML config
	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return nil, errors.Wrap(wErr, "failed to find worker nodes")
	}

	for _, workerNode := range workerNodes {
		currentConfig := don.MustNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides
		updatedConfig, uErr := updateNodeConfig(workerNode, currentConfig, registryChainID, common.HexToAddress(workflowRegistryAddress), creEnv.ContractVersions[keystone_changeset.WorkflowRegistry.String()])
		if uErr != nil {
			return nil, errors.Wrapf(uErr, "failed to update node config for node index %d", workerNode.Index)
		}
		don.MustNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = *updatedConfig
	}

	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "vault",
			Version:        "1.0.0",
			CapabilityType: 1, // ACTION
		},
		Config: &capabilitiespb.CapabilityConfig{
			LocalOnly: don.HasOnlyLocalCapabilities(),
		},
		UseCapRegOCRConfig: true,
	}}

	workers, wErr := don.Workers()
	if wErr != nil {
		return nil, errors.Wrap(wErr, "failed to find worker nodes")
	}
	numWorkers := len(workers)

	dkgRPC, dErr := dkgReportingPluginConfigFromMetadata(don)
	if dErr != nil {
		return nil, errors.Wrap(dErr, "failed to create DKG reporting plugin config")
	}
	dkgRPCBytes, mErr := dkgRPC.MarshalBinary()
	if mErr != nil {
		return nil, errors.Wrap(mErr, "failed to marshal DKG reporting plugin config")
	}

	vaultRPCFactory := func(capRegAddress string, chainID uint64, donID uint32, capabilityID string, generatedConfigs map[string]*capabilitiespb.OCR3Config) ([]byte, error) {
		dkgCfg, ok := generatedConfigs["dkg"]
		if !ok {
			return nil, fmt.Errorf("DKG config not found in generated configs (needed to compute DKG instance ID)")
		}

		digest, err := capregconfig.ComputeConfigDigest(chainID, capRegAddress, capabilityID, donID, "dkg", dkgCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to compute DKG config digest: %w", err)
		}

		instanceID := string(dkgocrtypes.MakeInstanceID(common.HexToAddress(capRegAddress), digest))
		cfg := vaultprotos.ReportingPluginConfig{
			DKGInstanceID:                   &instanceID,
			EnableDeterministicPendingQueue: pendingQueueEnabledFromMetadata(don),
		}
		return proto.Marshal(&cfg)
	}

	ocr3Config := contracts.DefaultOCR3_1Config(numWorkers)

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
		CapabilityToOCR3_1Config: map[string]*cre.OCR3_1CapabilityConfig{
			"vault": {
				Configs: map[string]*cre.OCR3_1ConfigEntry{
					"dkg": {
						Config:                copyOCR3_1Config(ocr3Config),
						ReportingPluginConfig: dkgRPCBytes,
					},
					"vault": {
						Config:                       copyOCR3_1Config(ocr3Config),
						ReportingPluginConfigFactory: vaultRPCFactory,
					},
				},
			},
		},
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
		Address:         ptr.Ptr(workflowRegistryAddress.Hex()),
		NetworkID:       ptr.Ptr("evm"),
		ChainID:         ptr.Ptr(strconv.FormatUint(registryChainID, 10)),
		SyncStrategy:    ptr.Ptr("reconciliation"),
		ContractVersion: ptr.Ptr(wfRegVersion.String()),
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
	return createJobs(ctx, creEnv, don, dons)
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

	specs := make(map[string][]string)

	_, ocrPeeringCfg, err := cre.PeeringCfgs(bootstrap)
	if err != nil {
		return errors.Wrap(err, "failed to get peering configs")
	}

	capRegVersion := creEnv.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()].String()

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
			"capRegVersion":        capRegVersion,
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

func dkgReportingPluginConfigFromMetadata(don *cre.DonMetadata) (*dkgocrtypes.ReportingPluginConfig, error) {
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

func pendingQueueEnabledFromMetadata(don *cre.DonMetadata) bool {
	cc, ok := don.CapabilityConfigs[flag]
	if !ok {
		return false
	}
	setting, ok := cc.Values["EnableDeterministicPendingQueue"]
	if !ok {
		return false
	}
	enabled, ok := setting.(bool)
	if !ok {
		return false
	}
	return enabled
}

func copyOCR3_1Config(src *ocr3_1.V3_1OracleConfig) *ocr3_1.V3_1OracleConfig {
	cp := *src
	cp.TransmissionSchedule = make([]int, len(src.TransmissionSchedule))
	copy(cp.TransmissionSchedule, src.TransmissionSchedule)
	return &cp
}

func EncryptSecret(secret, masterPublicKeyStr string, owner common.Address) (string, error) {
	masterPublicKey := tdh2easy.PublicKey{}
	masterPublicKeyBytes, err := hex.DecodeString(masterPublicKeyStr)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode master public key")
	}
	err = masterPublicKey.Unmarshal(masterPublicKeyBytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal master public key")
	}
	var label [32]byte
	copy(label[12:], owner.Bytes()) // left-pad with 12 zero
	cipher, err := tdh2easy.EncryptWithLabel(&masterPublicKey, []byte(secret), label)
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt secret")
	}
	cipherBytes, err := cipher.Marshal()
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal encrypted secrets to bytes")
	}
	return hex.EncodeToString(cipherBytes), nil
}
