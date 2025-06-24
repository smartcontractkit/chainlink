package cre

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/goccy/go-json"
	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	computecap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/compute"
	consensuscap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/consensus"
	croncap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/cron"
	gatewayconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/gateway"
	crecompute "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/compute"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	cregateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
)

const (
	ConfigTOML               = "environment-one-don-single-chain.toml"
	CreConfigDefaultFileName = "cre.yaml"
	DefaultBinaryFilename    = "mytestworkflow.wasm.br"
	WorkflowName             = "mytestworkflow"

	SomeAddr = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	AnvilPk  = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	FakeFeedID     = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
	FakeFeedURL    = "http://example.com"
	FakeTargetName = "mytestprofile"
	FakeGhToke     = "ghp_QwJ5CZr8SL0x9vKSCaMxOHcDkT6QOPva3OpR"
	FakeProfile    = "mytestprofile"
)

func prepEnv(t *testing.T) {
	err := os.Setenv("CRE_ETH_PRIVATE_KEY", AnvilPk)
	require.NoError(t, err)

	err = os.Setenv("CRE_GITHUB_API_TOKEN", FakeGhToke)
	require.NoError(t, err)

	err = os.Setenv("CRE_PROFILE", FakeProfile)
	require.NoError(t, err)

	err = os.Setenv("CTF_CONFIGS", ConfigTOML)
	require.NoError(t, err)

	err = os.Setenv("PRIVATE_KEY", AnvilPk)
	require.NoError(t, err)
}

type mockMinioStorageSettings struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	SessionToken    string `yaml:"session_token"`
	UseSSL          bool   `yaml:"use_ssl"`
	Region          string `yaml:"region"`
}

type mockWorkflowStorageSettings struct {
	Minio mockMinioStorageSettings `yaml:"minio"`
}

type mockSettings struct {
	DevPlatform     crecli.DevPlatform          `yaml:"dev-platform"`
	StorageSettings mockWorkflowStorageSettings `yaml:"workflow_storage"`
	UserWorkflow    crecli.UserWorkflow         `yaml:"user-workflow"`
	Contracts       crecli.Contracts            `yaml:"contracts"`
	RPCS            []crecli.RPC                `json:"rpcs"`
}

type mockProfileSettings struct {
	Settings mockSettings `yaml:"mytestprofile"`
}

type PoRWorkflowConfig struct {
	FeedID            string  `yaml:"feed_id"`
	URL               string  `yaml:"url"`
	ConsumerAddress   string  `yaml:"consumer_address"`
	WriteTargetName   string  `yaml:"write_target_name"`
	AuthKeySecretName *string `yaml:"auth_key_secret_name,omitempty"`
}

func mockWorkflowConfig(t *testing.T, workflowConfigPath string) {
	workflowConfig := PoRWorkflowConfig{
		FeedID:            FakeFeedID,
		URL:               FakeFeedURL,
		ConsumerAddress:   SomeAddr,
		WriteTargetName:   FakeTargetName,
		AuthKeySecretName: nil,
	}

	config, err := json.Marshal(workflowConfig)
	require.NoError(t, err)
	err = os.WriteFile(workflowConfigPath, config, 0600)
	require.NoError(t, err)

	// another hack to get CRE CLI to work; same configuration as JSON to get verification to work
	config, err = yaml.Marshal(workflowConfig)
	require.NoError(t, err)
	err = os.WriteFile(path.Join(filepath.Dir(workflowConfigPath), "workflow.yaml"), config, 0600)
	require.NoError(t, err)
}

func mockCreConfig(t *testing.T, configPath string, provider s3provider.Provider) {
	settings := mockProfileSettings{
		Settings: mockSettings{
			DevPlatform: crecli.DevPlatform{
				DonID: 1,
			},
			StorageSettings: mockWorkflowStorageSettings{
				Minio: mockMinioStorageSettings{
					Endpoint:        provider.GetEndpoint(),
					AccessKeyID:     provider.GetAccessKey(),
					SecretAccessKey: provider.GetSecretKey(),
					SessionToken:    uuid.NewString(),
					UseSSL:          false,
					Region:          provider.GetRegion(),
				},
			},
			UserWorkflow: crecli.UserWorkflow{
				WorkflowOwnerAddress: SomeAddr,
				WorkflowName:         WorkflowName,
			},
			Contracts: crecli.Contracts{
				ContractRegistry: []crecli.ContractRegistry{
					{
						Name:          "WorkflowRegistry",
						Address:       "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0",
						ChainSelector: 3379446385462418246,
					},
					{
						Name:          "CapabilitiesRegistry",
						Address:       "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512",
						ChainSelector: 3379446385462418246,
					},
				},
				DataFeeds: nil,
				Keystone:  nil,
			},

			RPCS: []crecli.RPC{
				{
					ChainSelector: 3379446385462418246,
					URL:           "http://127.0.0.1:8545",
				},
			},
		},
	}

	creConfig, err := yaml.Marshal(settings)
	require.NoError(t, err)

	creConfigFilePath := path.Join(configPath, CreConfigDefaultFileName)
	err = os.WriteFile(creConfigFilePath, creConfig, 0600)
	require.NoError(t, err)
}

func TestWorkflowS3(t *testing.T) {
	configErr := setCICtfConfigIfMissing(ConfigTOML)
	require.NoError(t, configErr, "failed to set CTF config")

	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err)

	// TODO: PoR repo path coming from the framework runner
	porRepoPath := crecli.DerefString(in.WorkflowConfigs[0].WorkflowFolderLocation)

	workflowFilePath := path.Join(porRepoPath, "main.go")
	workflowConfigPath := path.Join(porRepoPath, "config.json")

	prepEnv(t)

	mockWorkflowConfig(t, workflowConfigPath)

	cli := crecli.NewCreCli(in.DependenciesConfig.CRECLIBinaryPath)
	err = cli.Compile(
		workflowFilePath,
		workflowConfigPath,
		DefaultBinaryFilename,
	)
	require.NoError(t, err)

	wasmFilePath := path.Join(filepath.Dir(workflowFilePath), DefaultBinaryFilename)
	// hack to address CRE CLI broken behaviour: `-o` fileName.wasm.br produces `fileName.wasm.br.b64`
	// and its `upload` command returns:
	// `Error: failed to create or update object: ... supported extensions: .wasm.br, .json, .yaml, .yml
	err = os.Rename(wasmFilePath+".b64", wasmFilePath)
	require.NoError(t, err)

	stats, err := os.Stat(wasmFilePath)
	require.NoError(t, err)
	require.NotZero(t, stats.Size())

	provider, err := s3provider.NewMinioFactory().New()
	require.NoError(t, err)

	mockCreConfig(t, filepath.Dir(workflowFilePath), provider)

	output, err := cli.Upload(
		crecli.MINIO,
		wasmFilePath,
		workflowConfigPath,
	)
	require.NoError(t, err)

	fmt.Printf("%#v\n", output)

	require.NoError(t, err, "couldn't load test config")
	require.Len(t, in.NodeSets, 1, "expected 1 node set in the test config")

	firstBlockchain := in.Blockchains[0]

	chainIDInt, err := strconv.Atoi(firstBlockchain.ChainID)
	require.NoError(t, err, "failed to convert chain ID to int")
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       SinglePoRDonCapabilitiesFlags,
				DONTypes:           []string{keystonetypes.WorkflowDON, keystonetypes.GatewayDON},
				BootstrapNodeIndex: 0, // not required, but set to make the configuration explicit
				GatewayNodeIndex:   0, // not required, but set to make the configuration explicit
			},
		}
	}

	capabilityFactoryFns := []keystonetypes.DONCapabilityWithConfigFactoryFn{
		computecap.ComputeCapabilityFactoryFn,
		consensuscap.OCR3CapabilityFactoryFn,
		croncap.CronCapabilityFactoryFn,
	}

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:            mustSetCapabilitiesFn(in.NodeSets),
		CapabilitiesContractFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:                     in.Blockchains,
		JdInput:                              *in.JD,
		InfraInput:                           *in.Infra,
		JobSpecFactoryFunctions: []keystonetypes.JobSpecFactoryFn{
			creconsensus.ConsensusJobSpecFactoryFn(chainIDUint64),
			cregateway.GatewayJobSpecFactoryFn([]int{}, []string{}, []string{"0.0.0.0/0"}),
			crecompute.ComputeJobSpecFactoryFn,
		},
		ConfigFactoryFunctions: []keystonetypes.ConfigFactoryFn{
			gatewayconfig.GenerateConfig,
		},
	}

	_, err = creenv.SetupTestEnvironment(
		t.Context(),
		framework.L,
		cldlogger.NewSingleFileLogger(t),
		universalSetupInput,
	)
	require.NoError(t, err, "failed to setup test environment")

	creCofigHandle, err := os.Open(path.Join(porRepoPath, CreConfigDefaultFileName))
	require.NoError(t, err)

	err = crecli.DeployWorkflow(
		cli.CreCLICommandPath,
		output.BinaryURL,
		&output.ConfigURL,
		nil,
		creCofigHandle,
		&porRepoPath,
	)
	require.NoError(t, err)
}
