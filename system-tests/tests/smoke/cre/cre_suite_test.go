package cre

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"dario.cat/mergo"
	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	cre_jobs_pkg "github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

//////////// SMOKE TESTS /////////////
// target happy path and sanity checks
// all other tests (e.g. edge cases, negative conditions)
// should go to a `regression` package
/////////////////////////////////////

var v2RegistriesFlags = []string{"--with-contracts-version", "v2"}

/*
To execute tests locally start the local CRE first:
Inside `core/scripts/cre/environment` directory
 1. Ensure the necessary capabilities (i.e. readcontract, http-trigger, http-action) are listed in the environment configuration
 2. Identify the appropriate topology that you want to test
 3. Stop and clear any existing environment: `go run . env stop -a`
 4. Run: `CTF_CONFIGS=<path-to-your-topology-config> go run . env start && ./bin/ctf obs up` to start env + observability
 5. Optionally run blockscout `./bin/ctf bs up`
 6. Execute the tests in `system-tests/tests/smoke/cre`: `go test -timeout 15m -run "^Test_CRE_V2"`.
*/
func Test_CRE_V1_Proof_Of_Reserve(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))
	// WARNING: currently we can't run these tests in parallel, because each test rebuilds environment structs and that includes
	// logging into CL node with GraphQL API, which allows only 1 session per user at a time.

	// requires `readcontract`, `cron`
	priceProvider, porWfCfg := beforePoRTest(t, testEnv, "por-workflowV1", PoRWFV1Location)
	ExecutePoRTest(t, testEnv, priceProvider, porWfCfg, false)
}

func Test_CRE_V1_Tron(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-don-tron.toml"))

	priceProvider, porWfCfg := beforePoRTest(t, testEnv, "por-workflowV1", PoRWFV1Location)
	ExecutePoRTest(t, testEnv, priceProvider, porWfCfg, false)
}

func Test_CRE_V1_SecureMint(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-don-solana.toml"))

	ExecuteSecureMintTest(t, testEnv)
}

/*
// TODO: Move Billing tests to v2 Registries
func Test_CRE_V1_Billing_EVM_Write(t *testing.T) {
	quarantine.Flaky(t, "DX-1911")
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))

	require.NoError(
		t,
		startBillingStackIfIsNotRunning(t, testEnv.TestConfig.RelativePathToRepoRoot, testEnv.TestConfig.EnvironmentDirPath, testEnv),
		"failed to start Billing stack",
	)

	priceProvider, porWfCfg := beforePoRTest(t, testEnv, "por-workflowV2-billing", PoRWFV2Location)
	porWfCfg.FeedIDs = []string{porWfCfg.FeedIDs[0]}
	ExecutePoRTest(t, testEnv, priceProvider, porWfCfg, true)
}
*/

func Test_CRE_V1_Billing_Cron_Beholder(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))

	require.NoError(
		t,
		startBillingStackIfIsNotRunning(t, testEnv.TestConfig.RelativePathToRepoRoot, testEnv.TestConfig.EnvironmentDirPath, testEnv),
		"failed to start Billing stack",
	)

	ExecuteBillingTest(t, testEnv)
}

// buildRuntimeValues creates runtime-generated  values for any keys not specified in TOML
func buildRuntimeValues(chainID uint64, networkFamily, creForwarderAddress, nodeAddress string) map[string]any {
	return map[string]any{
		"ChainID":             chainID,
		"NetworkFamily":       networkFamily,
		"CreForwarderAddress": creForwarderAddress,
		"NodeAddress":         nodeAddress,
	}
}

func Test_EVM_Job_Update(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)
	flag := cre.EVMCapability
	creEnv := testEnv.CreEnvironment
	dons := testEnv.Dons

	var f = func(don *cre.Don) error {
		// horrible copy & paste of `createJobs` function from /Users/bartektofel/Desktop/repos/chainlink/system-tests/lib/cre/features/evm/v2/evm.go

		configTemplate := `{"chainId":{{.ChainID}}, "network":"{{.NetworkFamily}}", "logTriggerPollInterval":{{printf "%0.f" .LogTriggerPollInterval}}, "creForwarderAddress":"{{.CreForwarderAddress}}", "receiverGasMinimum":{{.ReceiverGasMinimum}}, "nodeAddress":"{{.NodeAddress}}"{{with .LogTriggerSendChannelBufferSize}},"logTriggerSendChannelBufferSize":{{printf "%d" .}}{{end}}{{with .LogTriggerLimitQueryLogSize}},"logTriggerLimitQueryLogSize":{{printf "%d" .}}{{end}}}`
		specs := make(map[string][]string)
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

		bootstrap, isBootstrap := dons.Bootstrap()
		if !isBootstrap {
			return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
		}

		workerNodes, wErr := don.Workers()
		if wErr != nil {
			return errors.Wrap(wErr, "failed to find worker nodes")
		}

		chainConfig, ok := nodeSet.GetChainCapabilityConfigs()[flag]
		if !ok {
			return fmt.Errorf("could not find capability config for capability %s in node set %s", flag, nodeSet.GetName())
		}

		for _, chainID := range chainConfig.EnabledChains {
			chainSelector, selErr := chainselectors.SelectorFromChainId(chainID)
			if selErr != nil {
				return errors.Wrapf(selErr, "failed to get chain selector from chainID %d", chainID)
			}
			chainIDStr := strconv.FormatUint(chainID, 10)
			qualifier := ks_contracts_op.CapabilityContractIdentifier(chainID)

			ocr3Key := datastore.NewAddressRefKey(
				chainSelector,
				datastore.ContractType(keystone_changeset.OCR3Capability.String()),
				semver.MustParse("1.0.0"),
				qualifier,
			)
			ocr3ConfigContractAddress, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(ocr3Key)
			if err != nil {
				return errors.Wrapf(err, "failed to get contract address for key %s and chainID %d", ocr3Key, chainID)
			}

			for _, workerNode := range workerNodes {
				evmKey, ok := workerNode.Keys.EVM[chainID]
				if !ok {
					return fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
				}
				nodeAddress := evmKey.PublicAddress.Hex()

				creForwarderKey := datastore.NewAddressRefKey(
					chainSelector,
					datastore.ContractType(keystone_changeset.KeystoneForwarder.String()),
					semver.MustParse("1.0.0"),
					"",
				)
				creForwarderAddress, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(creForwarderKey)
				if err != nil {
					return errors.Wrap(err, "failed to get CRE Forwarder address")
				}

				runtimeFallbacks := buildRuntimeValues(chainID, "evm", creForwarderAddress.Address, nodeAddress)

				capabilityConfig.Config["ReceiverGasMinimum"] = 521
				_, templateData, rErr := envconfig.ResolveCapabilityForChain(flag, nodeSet.GetChainCapabilityConfigs(), capabilityConfig.Config, chainID)
				if rErr != nil {
					return errors.Wrap(rErr, "failed to resolve capability config for chain")
				}

				var aErr error
				templateData, aErr = credon.ApplyRuntimeValues(templateData, runtimeFallbacks)
				if aErr != nil {
					return errors.Wrap(aErr, "failed to apply runtime values")
				}

				tmpl, err := template.New("evmConfig").Parse(configTemplate)
				if err != nil {
					return errors.Wrapf(err, "failed to parse %s config template", flag)
				}

				var configBuffer bytes.Buffer
				if err := tmpl.Execute(&configBuffer, templateData); err != nil {
					return errors.Wrapf(err, "failed to execute %s config template", flag)
				}

				configStr := configBuffer.String()

				if err := credon.ValidateTemplateSubstitution(configStr, flag); err != nil {
					return errors.Wrapf(err, "%s template validation failed", flag)
				}

				evmKeyBundle, ok := workerNode.Keys.OCR2BundleIDs[chainselectors.FamilyEVM] // we can always expect evm bundle key id present since evm is the registry chain
				if !ok {
					return errors.New("failed to get key bundle id for evm family")
				}

				bootstrapPeers := []string{fmt.Sprintf("%s@%s:%d", strings.TrimPrefix(bootstrap.Keys.PeerID(), "p2p_"), bootstrap.Host, cre.OCRPeeringPort)}

				strategyName := "single-chain"
				if len(workerNode.Keys.OCR2BundleIDs) > 1 {
					strategyName = "multi-chain"
				}

				workerInput := cre_jobs.ProposeJobSpecInput{
					Domain:      offchain.ProductLabel,
					Environment: cre.EnvironmentName,
					DONName:     don.Name,
					JobName:     fmt.Sprintf("evm-capabilities-v2-%d", chainID),
					ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
					DONFilters: []offchain.TargetDONFilter{
						{Key: offchain.FilterKeyDONName, Value: don.Name},
						{Key: "p2p_id", Value: workerNode.Keys.PeerID()}, // required since each node requires a different config (it contains its own from address)
					},
					Template: job_types.EVM,
					Inputs: job_types.JobSpecInput{
						"command": command,
						"config":  configStr,
						"oracleFactory": cre_jobs_pkg.OracleFactory{
							Enabled:            true,
							ChainID:            chainIDStr,
							BootstrapPeers:     bootstrapPeers,
							OCRContractAddress: ocr3ConfigContractAddress.Address,
							OCRKeyBundleID:     evmKeyBundle,
							TransmitterID:      nodeAddress,
							OnchainSigningStrategy: cre_jobs_pkg.OnchainSigningStrategy{
								StrategyName: strategyName,
								Config:       workerNode.Keys.OCR2BundleIDs,
							},
						},
					},
				}

				workerVerErr := cre_jobs.ProposeJobSpec{}.VerifyPreconditions(*creEnv.CldfEnvironment, workerInput)
				if workerVerErr != nil {
					return fmt.Errorf("precondition verification failed for EVM v2 worker job: %w", workerVerErr)
				}

				workerReport, workerErr := cre_jobs.ProposeJobSpec{}.Apply(*creEnv.CldfEnvironment, workerInput)
				if workerErr != nil {
					return fmt.Errorf("failed to propose EVM v2 worker job spec: %w", workerErr)
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
		}

		approveErr := jobs.Approve(t.Context(), creEnv.CldfEnvironment.Offchain, dons, specs)
		if approveErr != nil {
			return fmt.Errorf("failed to approve EVM v2 jobs: %w", approveErr)
		}

		return nil
	}

	dd := dons.DonsWithFlags(flag)

	for _, don := range dd {
		err := f(don)
		require.NoError(t, err)
	}
}

//////////// V2 TESTS /////////////
/*
To execute tests with v2 contracts start the local CRE first:
 1. Inside `core/scripts/cre/environment` directory: `go run . env restart --with-beholder --with-contracts-version v2`
 2. Execute the tests in `system-tests/tests/smoke/cre`: `go test -timeout 15m -run "^Test_CRE_V2"`.
*/
func Test_CRE_V2_Suite(t *testing.T) {
	topology := os.Getenv("TOPOLOGY_NAME")
	t.Run("[v2] Proof Of Reserve - "+topology, func(t *testing.T) {
		testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)
		priceProvider, wfConfig := beforePoRTest(t, testEnv, "por-workflow-v2", PoRWFV2Location)
		ExecutePoRTest(t, testEnv, priceProvider, wfConfig, false)
	})

	t.Run("[v2] Vault DON - "+topology, func(t *testing.T) {
		testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

		ExecuteVaultTest(t, testEnv)
	})

	t.Run("[v2] Cron Beholder - "+topology, func(t *testing.T) {
		testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

		ExecuteCronBeholderTest(t, testEnv)
	})

	t.Run("[v2] HTTP Trigger Action - "+topology, func(t *testing.T) {
		testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

		ExecuteHTTPTriggerActionTest(t, testEnv)
	})

	t.Run("[v2] HTTP Action CRUD Success - "+topology, func(t *testing.T) {
		testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

		ExecuteHTTPActionCRUDSuccessTest(t, testEnv)
	})

	t.Run("[v2] DON Time - "+topology, func(t *testing.T) {
		testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

		ExecuteDonTimeTest(t, testEnv)
	})
	t.Run("[v2] Consensus - "+topology, func(t *testing.T) {
		testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

		ExecuteConsensusTest(t, testEnv)
	})
}

func Test_CRE_V2_EVM_Suite(t *testing.T) {
	topology := os.Getenv("TOPOLOGY_NAME")
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

	t.Run("[v2] EVM Write - "+topology, func(t *testing.T) {
		priceProvider, porWfCfg := beforePoRTest(t, testEnv, "por-workflowV2", PoRWFV2Location)
		ExecutePoRTest(t, testEnv, priceProvider, porWfCfg, false)
	})

	t.Run("[v2] EVM Read - "+topology, func(t *testing.T) {
		ExecuteEVMReadTest(t, testEnv)
	})

	t.Run("[v2] EVM LogTrigger - "+topology, func(t *testing.T) {
		ExecuteEVMLogTriggerTest(t, testEnv)
	})
}

func Test_CRE_V2_HTTP_Action_Suite(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

	ExecuteHTTPActionCRUDSuccessTest(t, testEnv)
}

func Test_CRE_V2_Beholder_Suite(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), append(v2RegistriesFlags, "--with-dashboards")...)

	ExecuteLogStreamingTest(t, testEnv)
}
