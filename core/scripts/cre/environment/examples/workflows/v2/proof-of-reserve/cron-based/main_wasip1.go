package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"gopkg.in/yaml.v3"
	"main/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
	workflowpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
)

func RunSimpleCronWorkflow(
	config types.WorkflowConfig,
	logger *slog.Logger,
	secretsProvider cre.SecretsProvider,
) (cre.Workflow[types.WorkflowConfig], error) {
	return cre.Workflow[types.WorkflowConfig]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}), // every 30 seconds
			onTrigger,
		),
	}, nil
}

func onTrigger(config types.WorkflowConfig, runtime cre.Runtime, outputs *cron.Payload) (string, error) {
	runtime.Logger().With().Info(fmt.Sprintf("OnTrigger started with config %+v", config))
	evmClient := evm.Client{ChainSelector: chain_selectors.GETH_TESTNET.Selector}

	pricePromise := cre.RunInNodeMode(config, runtime,
		func(config types.WorkflowConfig, nodeRuntime cre.NodeRuntime) (priceOutput, error) {
			return getPrice(config, nodeRuntime)
		},
		cre.ConsensusIdenticalAggregation[priceOutput](),
	)

	price, err := pricePromise.Await()
	if err != nil {
		return "", fmt.Errorf("failed to get price: %w", err)
	}

	receivedFeedReportType, err := abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{Name: "FeedId", Type: "bytes32"},
		{Name: "Timestamp", Type: "uint32"},
		{Name: "Price", Type: "uint224"},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create receivedFeedReportType: %w", err)
	}

	args := abi.Arguments{
		{Type: receivedFeedReportType},
	}

	encodedPrice, err := args.Pack([]priceOutput{price})
	if err != nil {
		return "", fmt.Errorf("failed to pack price report: %w", err)
	}

	promise := runtime.GenerateReport(&workflowpb.ReportRequest{
		EncodedPayload: encodedPrice,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	})

	report, err := promise.Await()
	if err != nil {
		fmt.Println(">>> Error in generateReport")
		return "", fmt.Errorf("failed to generate report: %w", err)
	}

	wrPromise := evmClient.WriteReport(runtime, &evm.WriteCreReportRequest{
		Receiver:  common.Hex2Bytes(config.DataFeedsCacheAddress),
		Report:    report,
		GasConfig: &evm.GasConfig{GasLimit: 60000},
	})

	wrRep, err := wrPromise.Await()
	if err != nil {
		return "", err
	}

	var message = "empty message"
	if wrRep.ErrorMessage != nil {
		message = *wrRep.ErrorMessage
	}

	return message, nil
}

func main() {
	wasm.NewRunner(func(configBytes []byte) (types.WorkflowConfig, error) {
		cfg := types.WorkflowConfig{}
		if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
			return types.WorkflowConfig{}, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		if cfg.AuthKeySecretName != "" {
			cfg.AuthKey = sdk.SecretValue(cfg.AuthKeySecretName)
		}

		return cfg, nil
	}).Run(RunSimpleCronWorkflow)
}

// GenerateRandomBytes returns a slice of securely generated random bytes.
// It will return an error if the system's secure random number generator fails.
func GenerateRandomBytes(length int) []byte {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return bytes
}

func recoverSigner(rawReport, reportContext, signature []byte) (common.Address, error) {
	if len(signature) != 65 {
		return common.Address{}, fmt.Errorf("invalid signature length: %d", len(signature))
	}

	// Step 1: hash the raw report
	reportHash := crypto.Keccak256(rawReport)

	// Step 2: concat with reportContext
	message := append(reportHash, reportContext...)

	// Step 3: hash the full payload
	finalHash := crypto.Keccak256(message)

	// Step 4: recover public key from sig
	// Note: ECDSA signatures expect V in {27,28}
	v := signature[64]
	if v < 27 {
		v += 27
	}

	sig := append(signature[:64], v-27) // geth expects v ∈ {0,1}

	pubKey, err := crypto.SigToPub(finalHash, sig)
	if err != nil {
		return common.Address{}, err
	}

	// Step 5: convert to Ethereum address
	return crypto.PubkeyToAddress(*pubKey), nil
}

type priceOutput struct {
	FeedID    [32]byte
	Timestamp int64
	Price     int
}

type trueUSDResponse struct {
	AccountName string    `json:"accountName"`
	TotalTrust  float64   `json:"totalTrust"`
	Ripcord     bool      `json:"ripcord"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func getPrice(config types.WorkflowConfig, runtime cre.NodeRuntime) (priceOutput, error) {
	httpClient := &http.Client{}

	feedID, err := convertFeedIDtoBytes(config.FeedID)
	if err != nil {
		return priceOutput{}, fmt.Errorf("cannot convert feedID to bytes : %w : %b", err, feedID)
	}

	fetchRequest := http.Request{
		Url:       config.URL + "?feedID=" + config.FeedID,
		Method:    "GET",
		TimeoutMs: 5000,
	}

	runtime.Logger().With().Info(fmt.Sprintf("getPrice config authkey is %s nad athkey: %s", config.AuthKey, config.AuthKeySecretName))
	if string(config.AuthKey) != "" {
		runtime.Logger().With().Info(fmt.Sprintf("getPrice config set headers %s", config.AuthKey))

		fetchRequest.Headers = map[string]string{
			"Authorization": string(config.AuthKey),
		}
	}

	promiseResp := httpClient.SendRequest(runtime, &fetchRequest)
	if err != nil {
		return priceOutput{}, fmt.Errorf("failed to send price request to %s, err: %w", fetchRequest.String(), err)
	}

	r, err := promiseResp.Await()
	if err != nil {
		return priceOutput{}, fmt.Errorf("failed to await price response from %s and %v err: %w", fetchRequest.String(), fetchRequest.Headers, err)
	}

	var resp trueUSDResponse
	err = json.Unmarshal(r.Body, &resp)
	if err != nil {
		return priceOutput{}, fmt.Errorf("failed to unmarshal price response: %w", err)
	}

	runtime.Logger().With(
		"feedID", config.FeedID,
	).Info(fmt.Sprintf("TrueUSD price found: %.2f", resp.TotalTrust))

	if resp.Ripcord {
		runtime.Logger().With(
			"feedID", config.FeedID,
		).Info(fmt.Sprintf("ripcord flag set for feed ID %s", config.FeedID))
		return priceOutput{}, sdk.BreakErr
	}

	return priceOutput{
		FeedID:    feedID, // TrueUSD
		Timestamp: resp.UpdatedAt.Unix(),
		Price:     int(resp.TotalTrust * 100), // Convert to integer cents
	}, nil
}

func convertFeedIDtoBytes(feedIDStr string) ([32]byte, error) {
	if feedIDStr == "" {
		return [32]byte{}, fmt.Errorf("feedID string is empty")
	}

	if len(feedIDStr) < 2 {
		return [32]byte{}, fmt.Errorf("feedID string too short: %q", feedIDStr)
	}

	b, err := hex.DecodeString(feedIDStr[2:])
	if err != nil {
		return [32]byte{}, err
	}

	if len(b) < 32 {
		nb := [32]byte{}
		copy(nb[:], b[:])
		return nb, err
	}

	return [32]byte(b), nil
}
