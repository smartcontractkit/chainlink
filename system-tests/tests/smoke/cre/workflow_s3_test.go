package cre

import (
	"fmt"
	"github.com/goccy/go-json"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	"github.com/stretchr/testify/require"
)

func prepEnv(t *testing.T) {
	const (
		AnvilPk     = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		FakeGhToke  = "ghp_QwJ5CZr8SL0x9vKSCaMxOHcDkT6QOPva3OpR"
		FakeProfile = "mytestprofile"
	)

	err := os.Setenv("CRE_ETH_PRIVATE_KEY", AnvilPk)
	require.NoError(t, err)

	err = os.Setenv("CRE_GITHUB_API_TOKEN", FakeGhToke)
	require.NoError(t, err)

	err = os.Setenv("CRE_PROFILE", FakeProfile)
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
	StorageSettings mockWorkflowStorageSettings `yaml:"workflow_storage"`
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
	const (
		SomeAddr       = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
		FakeFeedID     = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
		FakeFeedURL    = "http://example.com"
		FakeTargetName = "mytestprofile"
	)

	workflowConfig := PoRWorkflowConfig{
		FeedID:            FakeFeedID,
		URL:               FakeFeedURL,
		ConsumerAddress:   SomeAddr,
		WriteTargetName:   FakeTargetName,
		AuthKeySecretName: nil,
	}

	config, err := json.Marshal(workflowConfig)
	require.NoError(t, err)
	err = os.WriteFile(workflowConfigPath, config, 0644)
	require.NoError(t, err)

	// another hack to get CRE CLI to work; same configuration as JSON to get verification to work
	config, err = yaml.Marshal(workflowConfig)
	require.NoError(t, err)
	err = os.WriteFile(path.Join(filepath.Dir(workflowConfigPath), "workflow.yaml"), config, 0644)
	require.NoError(t, err)
}

func mockCreConfig(t *testing.T, configPath string, provider s3provider.Provider) {
	const CreConfigDefaultFileName = "cre.yaml"

	settings := mockProfileSettings{
		Settings: mockSettings{
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
		},
	}

	creConfig, err := yaml.Marshal(settings)
	require.NoError(t, err)

	creConfigFilePath := path.Join(configPath, CreConfigDefaultFileName)
	err = os.WriteFile(creConfigFilePath, creConfig, 0644)
	require.NoError(t, err)
}

func TestWorkflowS3(t *testing.T) {
	const DefaultBinaryFilename = "mytestworkflow.wasm.br"

	// TODO: PoR repo path coming from the framework runner
	porRepoPath := ""

	workflowFilePath := path.Join(porRepoPath, "/cron-based/main.go")
	workflowConfigPath := path.Join(porRepoPath, "/cron-based/config.json")

	prepEnv(t)

	mockWorkflowConfig(t, workflowConfigPath)

	cli := crecli.NewCreCli()
	err := cli.Compile(
		workflowFilePath,
		workflowConfigPath,
		DefaultBinaryFilename,
	)
	require.NoError(t, err)

	wasmFilePath := path.Join(filepath.Dir(workflowFilePath), DefaultBinaryFilename)
	// hack to address CRE CLI broken behaviour: `-o` fileName.wasm.br produces `fileName.wasm.br.b64`
	// and its `upload` command returns:
	// `Error: failed to create or update object: ... supported extensions: .wasm.br, .json, .yaml, .yml
	err = os.Rename(fmt.Sprintf("%s.b64", wasmFilePath), wasmFilePath)
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

	// hack endpoints to be relative for the test-containers inside routing
	output.BinaryURL = strings.Replace(output.BinaryURL, provider.GetEndpoint(), provider.GetBaseEndpoint(), -1)
	output.ConfigURL = strings.Replace(output.ConfigURL, provider.GetEndpoint(), provider.GetBaseEndpoint(), -1)

	fmt.Printf("%#v\n", output)

	// TODO: bootup simple DON and deploy workflow from S3 using URLs from `output`
	//		 implement MinimumCLI.Deploy()
}
