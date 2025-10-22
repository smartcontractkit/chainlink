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
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr/donlevel"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/consensus"
)

const flag = cre.ConsensusCapabilityV2

type Consensus struct{}

func (c *Consensus) Flag() cre.CapabilityFlag {
	return flag
}

func (c *Consensus) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "consensus",
			Version:        "1.0.0-alpha",
			CapabilityType: 2, // CONSENSUS
			ResponseType:   0, // REPORT
		},
		Config: &capabilitiespb.CapabilityConfig{},
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

const ContractQualifier = "capability_consensus"

func (c *Consensus) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	_, ocr3ContractAddr, ocrErr := contracts.DeployOCR3Contract(testLogger, ContractQualifier, creEnv.RegistryChainSelector, creEnv.CldfEnvironment, creEnv.ContractVersions)
	if ocrErr != nil {
		return fmt.Errorf("failed to deploy OCR3 (consensus v2) contract %w", ocrErr)
	}

	jobsErr := createJobs(
		ctx,
		don,
		dons,
		creEnv,
	)
	if jobsErr != nil {
		return fmt.Errorf("failed to create OCR3 jobs: %w", jobsErr)
	}

	// wait for LP to be started (otherwise it won't pick up contract's configuration events)
	if err := consensus.WaitForLogPollerToBeHealthy(don); err != nil {
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
			ChainSelector:   creEnv.RegistryChainSelector,
			DON:             don.KeystoneDONConfig(),
			Config:          don.ResolveORC3Config(ocr3Config),
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
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	var generateJobSpec = func(logger zerolog.Logger, chainID uint64, nodeAddress string, mergedConfig map[string]any) (string, error) {
		runtimeFallbacks := buildRuntimeValues(chainID, "evm", nodeAddress)

		templateData, aErr := credon.ApplyRuntimeValues(mergedConfig, runtimeFallbacks)
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

	jobSpecs, jErr := ocr.GenerateJobSpecsForStandardCapabilityWithOCR(
		don,
		dons,
		creEnv,
		flag,
		func(_ uint64) string {
			return ContractQualifier
		},
		dataStoreOCR3ContractKeyProvider,
		donlevel.CapabilityEnabler,
		donlevel.EnabledChainsProvider,
		generateJobSpec,
		donlevel.ConfigMerger,
	)
	if jErr != nil {
		return errors.Wrap(jErr, "failed to generate EVM OCR3 job specs")
	}

	jobErr := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, dons, jobSpecs)
	if jobErr != nil {
		return fmt.Errorf("failed to create EVM OCR3 jobs for don %s: %w", don.Name, jobErr)
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
