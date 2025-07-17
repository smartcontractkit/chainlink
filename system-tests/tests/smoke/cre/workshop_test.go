package cre

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	dfilter "github.com/docker/docker/api/types/filters"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	consensuscap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/consensus"
	croncap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/cron"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	gatewayconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/gateway"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	crecron "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/cron"
	cregateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func Test_V2_Workflow_Workshop(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping workshop test in CI")
	}

	testLogger := framework.L

	/*
		TEST SETUP:
		- set required env vars
		- load test config
		- start DON
		- deploy contracts
		- create jobs
		- register workflow
	*/

	// set required env vars
	setPkErr := os.Setenv("PRIVATE_KEY", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80") // not a secret, it's a known developer key used by Anvil
	require.NoError(t, setPkErr, "failed to set PRIVATE_KEY")

	setCtfConfigsErr := os.Setenv("CTF_CONFIGS", "workshop_test.toml")
	require.NoError(t, setCtfConfigsErr, "failed to set CTF_CONFIGS")

	// load test config
	in, err := framework.Load[V2WorkflowTestConfig](t)
	require.NoError(t, err, "couldn't load test config")

	// setup test environment
	containerPath, pathErr := capabilities.DefaultContainerDirectory(in.Infra.InfraType)
	require.NoError(t, pathErr, "failed to get default container directory")

	chainIDInt, err := strconv.Atoi(in.Blockchain.ChainID)
	require.NoError(t, err, "failed to convert chain ID to int")
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets: []*types.CapabilitiesAwareNodeSet{
			{
				Input:              in.NodeSets[0],
				Capabilities:       []string{types.CronCapability, types.OCR3Capability, types.CustomComputeCapability, types.WriteEVMCapability},
				DONTypes:           []string{types.WorkflowDON, types.GatewayDON},
				BootstrapNodeIndex: 0,
			},
		},
		CapabilitiesContractFactoryFunctions: []func([]string) []keystone_changeset.DONCapabilityWithConfig{
			croncap.CronCapabilityFactoryFn,
			consensuscap.OCR3CapabilityFactoryFn,
		},
		BlockchainsInput: []*types.WrappedBlockchainInput{in.Blockchain},
		JdInput:          *in.JD,
		InfraInput:       *in.Infra,
		CustomBinariesPaths: map[string]string{
			types.CronCapability: in.Dependencies.CronBinaryPath,
		},
		JobSpecFactoryFunctions: []types.JobSpecFactoryFn{
			creconsensus.ConsensusJobSpecFactoryFn(chainIDUint64),
			crecron.CronJobSpecFactoryFn(filepath.Join(containerPath, filepath.Base(in.Dependencies.CronBinaryPath))),
			cregateway.GatewayJobSpecFactoryFn([]int{}, []string{}, []string{"0.0.0.0/0"}),
		},
		ConfigFactoryFunctions: []types.ConfigFactoryFn{
			gatewayconfig.GenerateConfig,
		},
		CustomAnvilMiner: in.CustomAnvilMiner,
	}

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(t.Context(), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
	require.NoError(t, setupErr, "failed to setup test environment")

	// compile and upload workflow
	containerTargetDir := "/home/chainlink/workflows"
	testLogger.Info().Msg("Proceeding to register test workflow...")
	workflowRegistryAddress, workflowRegistryErr := crecontracts.FindAddressesForChain(universalSetupOutput.CldEnvironment.ExistingAddresses, universalSetupOutput.BlockchainOutput[0].ChainSelector, keystone_changeset.WorkflowRegistry.String()) //nolint:staticcheck // won't migrate now
	require.NoError(t, workflowRegistryErr, "failed to find workflow registry address for chain %d", universalSetupOutput.BlockchainOutput[0].ChainID)

	// TODO: add code that will compile your workflow
	// Hint: use the creworkflow.CompileWorkflow() function to compile your workflow and get compressedWorkflowWasmPath variable

	copyErr := creworkflow.CopyWorkflowToDockerContainers(compressedWorkflowWasmPath, "workflow-node", containerTargetDir)
	require.NoError(t, copyErr, "failed to copy workflow to docker containers")

	t.Cleanup(func() {
		_ = os.Remove(compressedWorkflowWasmPath)
	})

	// register workflow
	registerErr := creworkflow.RegisterWithContract(
		t.Context(),
		universalSetupOutput.BlockchainOutput[0].SethClient,
		workflowRegistryAddress,
		universalSetupOutput.DonTopology.WorkflowDonID,
		"test-workflow",
		"file://"+compressedWorkflowWasmPath,
		nil, // no config URL
		nil, // no secrets URL
		&containerTargetDir,
	)
	require.NoError(t, registerErr, "failed to register workflow")

	/*
		TEST EXECUTION:
		- check container logs in loop until workflow execution is detected
		- or until we detect workflow engine initialization failure
	*/

	// wait for workflow to execute at least once
	err = waitForLog(testLogger, *regexp.MustCompile(".*Workflow execution finished successfully"), 1, 2*time.Minute)
	require.NoError(t, err, "failed to wait for log")

	testLogger.Info().Msg("Workflow executed successfully")
}

func waitForLog(testLogger zerolog.Logger, pattern regexp.Regexp, minOccurencePerNode int, timeout time.Duration) error {
	provider, err := tc.NewDockerProvider()
	if err != nil {
		return fmt.Errorf("failed to create Docker provider: %w", err)
	}
	containers, err := provider.Client().ContainerList(context.Background(), container.ListOptions{
		All: true,
		Filters: dfilter.NewArgs(dfilter.KeyValuePair{
			Key:   "label",
			Value: "framework=ctf",
		}),
	})
	if err != nil {
		return fmt.Errorf("failed to list Docker containers: %w", err)
	}
	workflowNodeContainers := make([]*container.Summary, 0)
	for _, containerInfo := range containers {
		if strings.Contains(containerInfo.Names[0], "workflow-node") && !strings.Contains(containerInfo.Names[0], "workflow-node0") {
			workflowNodeContainers = append(workflowNodeContainers, &containerInfo)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	eg := errgroup.Group{}

	for _, containerInfo := range workflowNodeContainers {
		eg.Go(func() error {
			ticker := time.NewTicker(10 * time.Second)
			occurrences := 0

			for {
				select {
				case <-ctx.Done():
					return fmt.Errorf("expected %d occurrences of %s, got %d for node %s", minOccurencePerNode, pattern.String(), occurrences, containerInfo.Names[0])
				case <-ticker.C:
					testLogger.Info().Msgf("Checking logs for %s", strings.TrimPrefix(containerInfo.Names[0], "/"))
					occurrences = 0
					logOptions := container.LogsOptions{ShowStdout: true, ShowStderr: true}
					logs, err := provider.Client().ContainerLogs(context.Background(), containerInfo.ID, logOptions)
					if err != nil {
						return errors.Wrap(err, "failed to read logs for container "+containerInfo.Names[0])
					}

					header := make([]byte, 8) // Docker stream header is 8 bytes
					for {
						_, err := io.ReadFull(logs, header)
						if err == io.EOF {
							break
						}
						if err != nil {
							break
						}

						msgSize := binary.BigEndian.Uint32(header[4:8])

						msg := make([]byte, msgSize)
						_, err = io.ReadFull(logs, msg)
						if err != nil {
							break
						}

						if pattern.Match(msg) {
							occurrences++
							if occurrences >= minOccurencePerNode {
								testLogger.Info().Msgf("Found expected occurrences of '%s' pattern in %s's logs", pattern.String(), strings.TrimPrefix(containerInfo.Names[0], "/"))
								return nil
							}
						}

						failedInitializationPattern := regexp.MustCompile(".*Workflow Engine initialization failed.*")
						if failedInitializationPattern.Match(msg) {
							return fmt.Errorf("workflow engine initialization failed for node %s", containerInfo.Names[0])
						}
					}
				}
			}
		})
	}

	return eg.Wait()
}

type WorkflowV2DependenciesConfig struct {
	CronBinaryPath string `toml:"cron_capability_binary_path" validate:"required"`
}

type V2WorkflowTestConfig struct {
	Blockchain       *types.WrappedBlockchainInput `toml:"blockchain" validate:"required"`
	NodeSets         []*ns.Input                   `toml:"nodesets" validate:"required"`
	JD               *jd.Input                     `toml:"jd" validate:"required"`
	Infra            *libtypes.InfraInput          `toml:"infra" validate:"required"`
	Dependencies     *WorkflowV2DependenciesConfig `toml:"dependencies" validate:"required"`
	CustomAnvilMiner *types.CustomAnvilMiner       `toml:"custom_anvil_miner"`
}
