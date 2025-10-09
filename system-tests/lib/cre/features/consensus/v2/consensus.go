package v2

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr/donlevel"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/consensus"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const flag = cre.ConsensusCapabilityV2

type Consensus struct{}

func (o *Consensus) Flag() cre.CapabilityFlag {
	return flag
}

func (o *Consensus) PreEnvStartup(
	testLogger zerolog.Logger,
	registryChainSelector uint64,
	cldfEnv *cldf.Environment,
	provider infra.Provider,
	topology *cre.Topology,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	capabilityConfigs cre.CapabilityConfigs,
	contractVersions map[string]string,
	gatewayConfigs map[cre.NodeUUID]*config.GatewayConfig,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := make(map[uint64][]keystone_changeset.DONCapabilityWithConfig)
	for _, donMetadata := range topology.DonsMetadataWithFlag(flag) {
		if capabilities[donMetadata.ID] == nil {
			capabilities[donMetadata.ID] = []keystone_changeset.DONCapabilityWithConfig{}
		}
		capabilities[donMetadata.ID] = append(capabilities[donMetadata.ID], keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "consensus",
				Version:        "1.0.0-alpha",
				CapabilityType: 2, // CONSENSUS
				ResponseType:   0, // REPORT
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfigs: capabilities,
	}, nil
}

const ConsensusV2ContractQualifier = "capability_consensus"

func (o *Consensus) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	nodeSetOutput []*cre.WrappedNodeOutput,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	contractVersions map[string]string,
	provider infra.Provider,
	capabilityConfigs map[string]cre.CapabilityConfig,
) error {
	// should we support more than one DON with OCR3 capability? Could there be 0? I guess as long as there's 1 with consensus v2?
	consensusV2DON, oneErr := creEnv.DonTopology.OneDonWithFlag(flag)
	if oneErr != nil {
		return oneErr
	}

	_, ocr3ContractAddr, ocrErr := contracts.DeployOCR3Contract(testLogger, ConsensusV2ContractQualifier, creEnv.DonTopology.HomeChainSelector, creEnv.CldfEnvironment, contractVersions)
	if ocrErr != nil {
		return fmt.Errorf("failed to deploy OCR3 (consensus v2) contract %w", ocrErr)
	}

	jobsErr := createJobs(
		ctx,
		creEnv.CldfEnvironment,
		// creEnv.DonTopology.HomeChainSelector,
		creEnv.DonTopology,
		provider,
		capabilityConfigs,
	)
	if jobsErr != nil {
		return fmt.Errorf("failed to create OCR3 jobs: %w", jobsErr)
	}

	// wait for LP to be started (otherwise it won't pick up contract's configuration events)
	if err := consensus.WaitForLogPollerToBeHealthy(consensusV2DON); err != nil {
		return errors.Wrap(err, "failed while waiting for Log Poller to become healthy")
	}

	ocr3Config, ocr3confErr := contracts.DefaultOCR3Config()
	if ocr3confErr != nil {
		return fmt.Errorf("failed to get default OCR3 config: %w", ocr3confErr)
	}

	_, ocr3Err := operations.ExecuteOperation(
		creEnv.CldfEnvironment.OperationsBundle,
		ks_contracts_op.ConfigureOCR3Op,
		ks_contracts_op.ConfigureOCR3OpDeps{
			Env: creEnv.CldfEnvironment,
		},
		ks_contracts_op.ConfigureOCR3OpInput{
			ContractAddress: ocr3ContractAddr,
			ChainSelector:   creEnv.DonTopology.HomeChainSelector,
			DON:             consensusV2DON.KeystoneDONConfig(),
			Config:          consensusV2DON.ResolveORC3Config(ocr3Config),
			DryRun:          false,
		},
	)

	if ocr3Err != nil {
		return errors.Wrap(ocr3Err, "failed to configure OCR3 contract")
	}

	return nil
}

const configTemplate = `'{"chainId":{{.ChainID}},"network":"{{.NetworkFamily}}","nodeAddress":"{{.NodeAddress}}"}'`

func createJobs(
	ctx context.Context,
	cldfEnv *cldf.Environment,
	// registryChainSelector uint64,
	donTopology *cre.DonTopology,
	provider infra.Provider,
	capabilityConfigs map[string]cre.CapabilityConfig,
) error {
	var generateJobSpec = func(logger zerolog.Logger, chainID uint64, nodeAddress string, mergedConfig map[string]any) (string, error) {
		runtimeFallbacks := buildRuntimeValues(chainID, "evm", nodeAddress)

		templateData, aErr := don.ApplyRuntimeValues(mergedConfig, runtimeFallbacks)
		if aErr != nil {
			return "", errors.Wrap(aErr, "failed to apply runtime values")
		}

		tmpl, err := template.New("consensusConfig").Parse(configTemplate)
		if err != nil {
			return "", errors.Wrap(err, "failed to parse consensus config template")
		}

		var configBuffer bytes.Buffer
		if err := tmpl.Execute(&configBuffer, templateData); err != nil {
			return "", errors.Wrap(err, "failed to execute consensus config template")
		}

		return configBuffer.String(), nil
	}

	var dataStoreOCR3ContractKeyProvider = func(contractName string, chainSelector uint64) datastore.AddressRefKey {
		return datastore.NewAddressRefKey(
			chainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			contractName,
		)
	}

	donsToJobSpecs, jErr := ocr.GenerateJobSpecsForStandardCapabilityWithOCR(
		donTopology,
		cldfEnv.DataStore,
		donTopology.Dons.AsNodeSetWithChainCapabilities(),
		provider,
		flag,
		func(_ uint64) string {
			return ConsensusV2ContractQualifier
		},
		dataStoreOCR3ContractKeyProvider,
		donlevel.CapabilityEnabler,
		donlevel.EnabledChainsProvider,
		generateJobSpec,
		donlevel.ConfigMerger,
		capabilityConfigs,
	)

	if jErr != nil {
		return errors.Wrap(jErr, "failed to generate EVM OCR3 job specs")
	}

	for _, don := range donTopology.Dons.List() {
		jobSpecs, ok := donsToJobSpecs[don.ID]
		if !ok {
			continue
		}
		jobErr := jobs.Create(ctx, cldfEnv.Offchain, donTopology, jobSpecs)

		if jobErr != nil {
			return fmt.Errorf("failed to create EVM OCR3 jobs for don %s: %w", don.Name, jobErr)
		}
	}

	return nil
}

func buildRuntimeValues(chainID uint64, networkFamily, nodeAddress string) map[string]any {
	return map[string]any{
		"ChainID":       chainID,
		"NetworkFamily": networkFamily,
		"NodeAddress":   nodeAddress,
	}
}
