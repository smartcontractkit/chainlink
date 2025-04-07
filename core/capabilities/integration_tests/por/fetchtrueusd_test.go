package por

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/compute"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/framework"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/keystone"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

const commandOverrideForCustomComputeAction = "__builtin_custom-compute-action"

var defaultConfig = compute.Config{
	ServiceConfig: webapi.ServiceConfig{
		RateLimiter: common.RateLimiterConfig{
			GlobalRPS:      100.0,
			GlobalBurst:    100,
			PerSenderRPS:   100.0,
			PerSenderBurst: 100,
		},
	},
}

type fetchTrueUSDConfig struct {
	ConsumerAddress         string
	WriteTargetCapabilityID string
	CronSchedule            string
}

func Test_FetchTrueUSDWorkflow(t *testing.T) {
	ctx := t.Context()

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.InfoLevel)

	fetchTrueUSDPath, err := filepath.Abs("./workflows/fetchtrueusd")
	require.NoError(t, err)

	wasmFile := filepath.Join(fetchTrueUSDPath, "fetchtrueusd.wasm")
	mainFile := filepath.Join(fetchTrueUSDPath, "main.go")
	framework.CreateWasmBinary(t, mainFile, wasmFile)

	triggerSink := framework.NewTriggerSink(t, "cron-trigger", "1.0.0")
	dataFeedsCacheContract := setupDons(ctx, t, lggr, wasmFile, "*/5 * * * * *", triggerSink)

	bundleReceived := make(chan *data_feeds_cache.DataFeedsCacheBundleReportUpdated, 1000)
	bundleSub, err := dataFeedsCacheContract.WatchBundleReportUpdated(&bind.WatchOpts{}, bundleReceived, nil, nil)
	require.NoError(t, err)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Simulate cron trigger firing
	wrappedMap, err := values.WrapMap(map[string]any{})
	require.NoError(t, err)
	triggerSink.SendOutput(wrappedMap, uuid.New().String())

loop:
	for {
		select {
		case <-ctxWithTimeout.Done():
			t.Fatalf("timed out waiting for bundle")
		case err := <-bundleSub.Err():
			require.NoError(t, err)
		case <-bundleReceived:
			break loop
		}
	}
}

func generateRandomReservesResponse() string {
	rootTotalTrust := 502939900.88
	rootTotalToken := 495516082.75

	deviation := func(value float64) float64 {
		return value * (1 + (rand.Float64()*0.1 - 0.05))
	}

	response := fmt.Sprintf(`{
		"accountName": "TrueUSD",
		"totalTrust": %.2f,
		"totalToken": %.2f,
		"updatedAt": "%s",
		"token": [
			{
				"tokenName": "TUSD (AVAX)",
				"totalTokenByChain": 2327025.78
			},
			{
				"tokenName": "TUSD (ETH)",
				"totalTokenByChain": 316126540.9523497
			},
			{
				"tokenName": "TUSD (TRON)",
				"totalTokenByChain": 167032152.1376503
			},
			{
				"tokenName": "TUSD (BSC)",
				"totalTokenByChain": 10030363.88
			}
		],
		"ripcord": false,
		"ripcordDetails": []
	}`, deviation(rootTotalTrust), deviation(rootTotalToken), time.Now().Format(time.RFC3339))

	return response
}

type ComputeFetcherFactory struct{}

func (n ComputeFetcherFactory) NewFetcher(log commonlogger.Logger, emitter custmsg.MessageEmitter) compute.FetcherFn {
	return func(ctx context.Context, req *host.FetchRequest) (*host.FetchResponse, error) {
		if req.URL == "https://api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD" {
			return &host.FetchResponse{
				Body: []byte(generateRandomReservesResponse()),
			}, nil
		}

		return nil, errors.New("no fetcher configured for url")
	}
}

func setupDons(ctx context.Context, t *testing.T, lggr logger.SugaredLogger, workflowURL string, cronSchedule string,
	triggerFactory framework.TriggerFactory) *data_feeds_cache.DataFeedsCache {
	configURL := "workflow-config.json"
	workflowConfig := fetchTrueUSDConfig{
		CronSchedule: cronSchedule,
	}

	compressedBinary, base64EncodedCompressedBinary := framework.GetCompressedWorkflowWasm(t, workflowURL)

	syncerFetcherFunc := func(ctx context.Context, messageID string, req capabilities.Request) ([]byte, error) {
		url := req.URL
		switch url {
		case workflowURL:
			return []byte(base64EncodedCompressedBinary), nil
		case configURL:
			configBytes, err := json.Marshal(workflowConfig)
			require.NoError(t, err)
			return configBytes, nil
		}

		return nil, fmt.Errorf("unknown  url: %s", url)
	}

	donContext := framework.CreateDonContextWithWorkflowRegistry(ctx, t, syncerFetcherFunc, ComputeFetcherFactory{})

	workflowDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "FetchTrueUSDWorkflow", NumNodes: 4, F: 1, AcceptsWorkflows: true})
	require.NoError(t, err)

	writeCapabilityDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "FetchTrueUSDWriteCapability", NumNodes: 4, F: 1, AcceptsWorkflows: false})
	require.NoError(t, err)

	// Setup DONs
	workflowDon := framework.NewDON(ctx, t, lggr, workflowDonConfiguration,
		[]commoncap.DON{writeCapabilityDonConfiguration.DON},
		donContext, true, 1*time.Second)

	writeCapabilityDon := framework.NewDON(ctx, t, lggr, writeCapabilityDonConfiguration,
		[]commoncap.DON{},
		donContext, true, 1*time.Second)

	// Setup workflow DON
	computeConfig, err := toml.Marshal(defaultConfig)
	require.NoError(t, err)
	workflowDon.AddStandardCapability("compute-capability", commandOverrideForCustomComputeAction, "'''"+string(computeConfig)+"'''")
	workflowDon.AddOCR3NonStandardCapability()
	workflowDon.AddExternalTriggerCapability(triggerFactory)

	workflowDon.Initialise()

	forwarderAddr, _ := keystone.SetupForwarderContract(t, workflowDon, donContext.EthBlockchain)

	workflowName := "TestWf"
	workflowOwner := donContext.EthBlockchain.TransactionOpts().From.String()

	dataFeedsCacheAddr, dataFeedsCache := SetupDataFeedsCacheContract(t, donContext.EthBlockchain, forwarderAddr, workflowOwner,
		workflows.HashTruncateName(workflowName))

	workflowConfig.ConsumerAddress = dataFeedsCacheAddr.String()

	// Setup Write capability DON
	writeTargetCapabilityID, err := writeCapabilityDon.AddPublishedEthereumWriteTargetNonStandardCapability(forwarderAddr)
	require.NoError(t, err)
	workflowConfig.WriteTargetCapabilityID = writeTargetCapabilityID

	writeCapabilityDon.Initialise()
	servicetest.Run(t, writeCapabilityDon)
	servicetest.Run(t, workflowDon)

	donContext.WaitForCapabilitiesToBeExposed(t, writeCapabilityDon, workflowDon)

	workflowConfigBytes, err := json.Marshal(workflowConfig)
	require.NoError(t, err)

	registerWorkflow(t, donContext, workflowName, compressedBinary, "", workflowDon,
		workflowURL, configURL, workflowConfigBytes)
	return dataFeedsCache
}

func SetupDataFeedsCacheContract(t *testing.T, backend *framework.EthBlockchain,
	forwarderAddress ethcommon.Address, workflowOwner string, workflowName string) (ethcommon.Address, *data_feeds_cache.DataFeedsCache) {
	addr, _, dataFeedsCache, err := data_feeds_cache.DeployDataFeedsCache(backend.TransactionOpts(), backend.Client())
	require.NoError(t, err)
	backend.Commit()

	var nameBytes [10]byte
	copy(nameBytes[:], workflowName)

	ownerAddr := ethcommon.HexToAddress(workflowOwner)

	_, err = dataFeedsCache.SetFeedAdmin(backend.TransactionOpts(), ownerAddr, true)
	require.NoError(t, err)

	backend.Commit()

	feedIDBytes := [16]byte{}
	copy(feedIDBytes[:], ethcommon.FromHex("0xA1B2C3D4E5F600010203040506070809"))

	_, err = dataFeedsCache.SetBundleFeedConfigs(backend.TransactionOpts(), [][16]byte{feedIDBytes}, []string{"fetchtrueusd"},
		[][]uint8{{2, 2, 2, 2}}, []data_feeds_cache.DataFeedsCacheWorkflowMetadata{
			{
				AllowedSender:        forwarderAddress,
				AllowedWorkflowOwner: ownerAddr,
				AllowedWorkflowName:  nameBytes,
			},
		})

	require.NoError(t, err)
	backend.Commit()

	return addr, dataFeedsCache
}

func registerWorkflow(t *testing.T, donContext framework.DonContext, workflowName string, compressedBinary []byte,
	secretsURL string, workflowDon *framework.DON, binaryURL string, configURL string, configBytes []byte) {
	workflowID, err := workflows.GenerateWorkflowID(donContext.EthBlockchain.TransactionOpts().From[:], workflowName, compressedBinary, configBytes, secretsURL)
	require.NoError(t, err)

	err = workflowDon.AddWorkflow(framework.Workflow{
		Name:       workflowName,
		ID:         workflowID,
		Status:     0,
		BinaryURL:  binaryURL,
		ConfigURL:  configURL,
		SecretsURL: secretsURL,
	})
	require.NoError(t, err)
}
