//go:build wasip1

package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/types/known/durationpb"
	"gopkg.in/yaml.v3"

	protos "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	pb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/confidentialhttp"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/maxlimits/config"
)

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		var cfg config.Config
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
		}
		return cfg, nil
	}).Run(RunWorkflow)
}

// ---------------------------------------------------------------------------
// Workflow registration: 10 trigger subscriptions (max TriggerSubscriptionLimit)
//   1 cron + 1 HTTP + 8 EVM log triggers
// ---------------------------------------------------------------------------

func RunWorkflow(cfg config.Config, logger *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	wf := cre.Workflow[config.Config]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: cfg.CronSchedule}),
			onCronTrigger,
		),
		cre.Handler(
			http.Trigger(buildHTTPTriggerConfig(cfg)),
			onHTTPTrigger,
		),
	}

	for _, lt := range cfg.LogTriggers {
		wf = append(wf, cre.Handler(
			evm.LogTrigger(lt.ChainSelector, buildLogFilter(lt)),
			onLogTrigger,
		))
	}

	logger.Info("Max-limits workflow registered",
		"triggers", len(wf),
		"chainReads", len(cfg.ChainReads),
		"chainWrites", len(cfg.ChainWrites),
		"httpEndpoints", len(cfg.HTTPEndpoints),
		"confHTTPEndpoints", len(cfg.ConfHTTPEndpoints),
		"secrets", len(cfg.Secrets),
		"consensusRounds", cfg.ConsensusRounds,
		"logEventCount", cfg.LogEventCount,
	)

	return wf, nil
}

// ---------------------------------------------------------------------------
// Cron handler: the main stress path exercising all per-execution call limits
// ---------------------------------------------------------------------------

func onCronTrigger(cfg config.Config, runtime cre.Runtime, _ *cron.Payload) (_ string, err error) {
	logger := runtime.Logger()
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in onCronTrigger: %v", r)
		}
	}()

	logger.Info("[maxlimits] cron handler started")

	// Phase 1: Secrets (up to 5 calls — Secrets.CallLimit)
	secretValues, sErr := runSecretsPhase(cfg, runtime)
	if sErr != nil {
		logger.Warn("[maxlimits] secrets phase error", "error", sErr)
	} else {
		logger.Info("[maxlimits] secrets phase done", "count", len(secretValues))
	}

	// Phase 2: HTTP actions via consensus (up to 5 HTTP calls + 5 consensus rounds)
	httpErr := runHTTPPhase(cfg, runtime)
	if httpErr != nil {
		logger.Warn("[maxlimits] HTTP phase error", "error", httpErr)
	} else {
		logger.Info("[maxlimits] HTTP phase done")
	}

	// Phase 3: Confidential HTTP (up to 5 calls — ConfidentialHTTP.CallLimit)
	confErr := runConfHTTPPhase(cfg, runtime)
	if confErr != nil {
		logger.Warn("[maxlimits] confHTTP phase error", "error", confErr)
	} else {
		logger.Info("[maxlimits] confHTTP phase done")
	}

	// Phase 4: Chain reads (up to 15 calls — ChainRead.CallLimit)
	readErr := runChainReadPhase(cfg, runtime)
	if readErr != nil {
		logger.Warn("[maxlimits] chain-read phase error", "error", readErr)
	} else {
		logger.Info("[maxlimits] chain-read phase done")
	}

	// Phase 5: Additional consensus rounds (up to 15 remaining — Consensus.CallLimit)
	consErr := runConsensusPhase(cfg, runtime)
	if consErr != nil {
		logger.Warn("[maxlimits] consensus phase error", "error", consErr)
	} else {
		logger.Info("[maxlimits] consensus phase done")
	}

	// Phase 6: Chain writes (up to 10 targets — ChainWrite.TargetsLimit)
	writeErr := runChainWritePhase(cfg, runtime)
	if writeErr != nil {
		logger.Warn("[maxlimits] chain-write phase error", "error", writeErr)
	} else {
		logger.Info("[maxlimits] chain-write phase done")
	}

	// Phase 7: Logging (up to 999 events — LogEventLimit minus overhead)
	runLoggingPhase(cfg, runtime)

	logger.Info("[maxlimits] cron handler completed all phases")

	return buildMaxResponse(), nil
}

// ---------------------------------------------------------------------------
// Phase 1: Secrets — up to Secrets.CallLimit (5)
// ---------------------------------------------------------------------------

func runSecretsPhase(cfg config.Config, runtime cre.Runtime) ([]string, error) {
	var results []string
	for _, s := range cfg.Secrets {
		resp, err := runtime.GetSecret(&protos.SecretRequest{
			Id:        s.ID,
			Namespace: s.Namespace,
		}).Await()
		if err != nil {
			return results, fmt.Errorf("secret %s/%s: %w", s.Namespace, s.ID, err)
		}
		results = append(results, resp.Value)
	}
	return results, nil
}

// ---------------------------------------------------------------------------
// Phase 2: HTTP actions — up to HTTPAction.CallLimit (5)
// Each http.SendRequest wraps RunInNodeMode (consumes 1 consensus round each)
// ---------------------------------------------------------------------------

func runHTTPPhase(cfg config.Config, runtime cre.Runtime) error {
	for i, ep := range cfg.HTTPEndpoints {
		client := &http.Client{}
		_, err := http.SendRequest(cfg, runtime, client,
			func(_ config.Config, _ *slog.Logger, req *http.SendRequester) (string, error) {
				resp, reqErr := req.SendRequest(&http.Request{
					Url:    ep.URL,
					Method: ep.Method,
					Body:   []byte(ep.Body),
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Timeout: &durationpb.Duration{Seconds: 10},
					CacheSettings: &http.CacheSettings{
						Store:  true,
						MaxAge: &durationpb.Duration{Seconds: 300},
					},
				}).Await()
				if reqErr != nil {
					return "", reqErr
				}
				return string(resp.Body), nil
			},
			cre.ConsensusIdenticalAggregation[string](),
		).Await()
		if err != nil {
			return fmt.Errorf("HTTP endpoint %d (%s): %w", i, ep.URL, err)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Phase 3: Confidential HTTP — up to ConfidentialHTTP.CallLimit (5)
// Runs at DON level, no separate consensus round needed
// ---------------------------------------------------------------------------

func runConfHTTPPhase(cfg config.Config, runtime cre.Runtime) error {
	client := confidentialhttp.Client{}
	for i, ep := range cfg.ConfHTTPEndpoints {
		_, err := client.SendRequest(runtime, &confidentialhttp.ConfidentialHTTPRequest{
			Request: &confidentialhttp.HTTPRequest{
				Url:    ep.URL,
				Method: ep.Method,
				Body:   &confidentialhttp.HTTPRequest_BodyString{BodyString: ep.Body},
				MultiHeaders: map[string]*confidentialhttp.HeaderValues{
					"Content-Type": {Values: []string{"application/json"}},
				},
			},
			VaultDonSecrets: []*confidentialhttp.SecretIdentifier{
				{Key: ep.SecretKey},
			},
		}).Await()
		if err != nil {
			return fmt.Errorf("confHTTP endpoint %d (%s): %w", i, ep.URL, err)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Phase 4: Chain reads — up to ChainRead.CallLimit (15)
// Mix of BalanceAt, CallContract, FilterLogs, HeaderByNumber,
// GetTransactionReceipt, GetTransactionByHash, EstimateGas
// ---------------------------------------------------------------------------

func runChainReadPhase(cfg config.Config, runtime cre.Runtime) error {
	for i, rc := range cfg.ChainReads {
		client := evm.Client{ChainSelector: rc.ChainSelector}
		var err error
		switch rc.Method {
		case "BalanceAt":
			_, err = client.BalanceAt(runtime, &evm.BalanceAtRequest{
				Account:     rc.AccountAddress,
				BlockNumber: nil,
			}).Await()
		case "CallContract":
			_, err = client.CallContract(runtime, &evm.CallContractRequest{
				Call: &evm.CallMsg{
					To:   rc.ContractAddress,
					Data: rc.CallData,
				},
			}).Await()
		case "FilterLogs":
			_, err = client.FilterLogs(runtime, &evm.FilterLogsRequest{
				FilterQuery: &evm.FilterQuery{
					FromBlock: pb.NewBigIntFromInt(big.NewInt(rc.FromBlock)),
					ToBlock:   pb.NewBigIntFromInt(big.NewInt(rc.ToBlock)),
					Addresses: [][]byte{rc.ContractAddress},
				},
			}).Await()
		case "HeaderByNumber":
			_, err = client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
				BlockNumber: pb.NewBigIntFromInt(big.NewInt(rc.BlockNumber)),
			}).Await()
		case "GetTxReceipt":
			_, err = client.GetTransactionReceipt(runtime, &evm.GetTransactionReceiptRequest{
				Hash: rc.TxHash,
			}).Await()
		case "GetTxByHash":
			_, err = client.GetTransactionByHash(runtime, &evm.GetTransactionByHashRequest{
				Hash: rc.TxHash,
			}).Await()
		case "EstimateGas":
			_, err = client.EstimateGas(runtime, &evm.EstimateGasRequest{
				Msg: &evm.CallMsg{
					To:   rc.ContractAddress,
					Data: rc.CallData,
				},
			}).Await()
		default:
			err = fmt.Errorf("unknown chain-read method: %s", rc.Method)
		}
		if err != nil {
			return fmt.Errorf("chain-read %d (%s): %w", i, rc.Method, err)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Phase 5: Additional consensus rounds — fills remaining Consensus.CallLimit (20)
// Phase 2 consumed up to 5, so we run up to 15 more here.
// Alternates between IdenticalAggregation and MedianAggregation.
// ---------------------------------------------------------------------------

func runConsensusPhase(cfg config.Config, runtime cre.Runtime) error {
	rounds := cfg.ConsensusRounds
	if rounds <= 0 {
		rounds = 15
	}
	payloadSize := cfg.ConsensusPayloadSize
	if payloadSize <= 0 {
		payloadSize = 99 * 1024 // just under 100 KB
	}

	for i := 0; i < rounds; i++ {
		if i%2 == 0 {
			_, err := cre.RunInNodeMode(cfg, runtime,
				func(_ config.Config, _ cre.NodeRuntime) ([]byte, error) {
					data := make([]byte, payloadSize)
					for j := range data {
						data[j] = byte((i + j) % 256)
					}
					return data, nil
				},
				cre.ConsensusIdenticalAggregation[[]byte](),
			).Await()
			if err != nil {
				return fmt.Errorf("consensus round %d (identical): %w", i, err)
			}
		} else {
			_, err := cre.RunInNodeMode(cfg, runtime,
				func(_ config.Config, _ cre.NodeRuntime) (int, error) {
					return 42 + i, nil
				},
				cre.ConsensusMedianAggregation[int](),
			).Await()
			if err != nil {
				return fmt.Errorf("consensus round %d (median): %w", i, err)
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Phase 6: Chain writes — up to ChainWrite.TargetsLimit (10)
// Each write: GenerateReport + WriteReport to a distinct receiver
// Reports approach ReportSizeLimit (5 KB), gas up to TransactionGasLimit (5M)
// ---------------------------------------------------------------------------

func runChainWritePhase(cfg config.Config, runtime cre.Runtime) error {
	payload := cfg.ReportPayload
	if len(payload) == 0 {
		payload = buildDefaultReportPayload()
	}

	for i, wt := range cfg.ChainWrites {
		report, err := runtime.GenerateReport(&cre.ReportRequest{
			EncodedPayload: payload,
			EncoderName:    "evm",
			SigningAlgo:    "ecdsa",
			HashingAlgo:    "keccak256",
		}).Await()
		if err != nil {
			return fmt.Errorf("generate report %d: %w", i, err)
		}

		client := evm.Client{ChainSelector: wt.ChainSelector}
		gasLimit := wt.GasLimit
		if gasLimit == 0 {
			gasLimit = 4_999_000 // just under 5M
		}
		_, err = client.WriteReport(runtime, &evm.WriteCreReportRequest{
			Receiver:  wt.Receiver,
			Report:    report,
			GasConfig: &evm.GasConfig{GasLimit: gasLimit},
		}).Await()
		if err != nil {
			return fmt.Errorf("write report %d: %w", i, err)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Phase 7: Logging — up to LogEventLimit (1000), each line under LogLineLimit (1 KB)
// We emit cfg.LogEventCount lines (default 999).
// ---------------------------------------------------------------------------

func runLoggingPhase(cfg config.Config, runtime cre.Runtime) {
	count := cfg.LogEventCount
	if count <= 0 {
		count = 999
	}
	for i := 0; i < count; i++ {
		runtime.Logger().Info(fmt.Sprintf("[maxlimits] telemetry line %d/%d: phase=logging ts=%d",
			i+1, count, i))
	}
}

// ---------------------------------------------------------------------------
// HTTP trigger handler (light path)
// ---------------------------------------------------------------------------

func onHTTPTrigger(cfg config.Config, runtime cre.Runtime, payload *http.Payload) (string, error) {
	runtime.Logger().Info("[maxlimits] HTTP trigger fired",
		"inputLen", len(payload.Input),
		"hasKey", payload.Key != nil,
	)
	return "maxlimits-http-ack", nil
}

// ---------------------------------------------------------------------------
// EVM log trigger handler (light path, shared by all 8 log triggers)
// ---------------------------------------------------------------------------

func onLogTrigger(cfg config.Config, runtime cre.Runtime, log *evm.Log) (string, error) {
	runtime.Logger().Info("[maxlimits] log trigger fired",
		"txHash", hex.EncodeToString(log.TxHash),
		"logIndex", log.Index,
		"dataLen", len(log.Data),
	)
	return "maxlimits-log-ack", nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func buildHTTPTriggerConfig(cfg config.Config) *http.Config {
	keys := make([]*http.AuthorizedKey, len(cfg.HTTPAuthorizedKeys))
	for i, k := range cfg.HTTPAuthorizedKeys {
		var kt http.KeyType
		switch k.Type {
		case "ECDSA_EVM":
			kt = http.KeyType_KEY_TYPE_ECDSA_EVM
		default:
			kt = http.KeyType_KEY_TYPE_ECDSA_EVM
		}
		keys[i] = &http.AuthorizedKey{
			Type:      kt,
			PublicKey: k.PublicKey,
		}
	}
	return &http.Config{AuthorizedKeys: keys}
}

func buildLogFilter(lt config.LogTriggerConfig) *evm.FilterLogTriggerRequest {
	addresses := make([][]byte, len(lt.Addresses))
	for i, a := range lt.Addresses {
		b, _ := hex.DecodeString(strings.TrimPrefix(a, "0x"))
		addresses[i] = b
	}

	topics := make([]*evm.TopicValues, len(lt.TopicSlots))
	for i, slot := range lt.TopicSlots {
		vals := make([][]byte, len(slot.Values))
		for j, v := range slot.Values {
			b, _ := hex.DecodeString(strings.TrimPrefix(v, "0x"))
			vals[j] = b
		}
		topics[i] = &evm.TopicValues{Values: vals}
	}

	return &evm.FilterLogTriggerRequest{
		Addresses:  addresses,
		Topics:     topics,
		Confidence: evm.ConfidenceLevel(lt.Confidence),
	}
}

// buildDefaultReportPayload creates an ABI-encoded payload near 5 KB.
func buildDefaultReportPayload() []byte {
	typ, err := abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{Name: "FeedID", Type: "bytes32"},
		{Name: "Timestamp", Type: "uint32"},
		{Name: "Price", Type: "uint224"},
	})
	if err != nil {
		return []byte("fallback-payload")
	}

	type report struct {
		FeedID    [32]byte
		Timestamp uint32
		Price     *big.Int
	}

	// ~150 bytes per entry × 30 entries ≈ 4.5 KB
	reports := make([]report, 30)
	for i := range reports {
		var feedID [32]byte
		copy(feedID[:], common.LeftPadBytes(big.NewInt(int64(i+1)).Bytes(), 32))
		reports[i] = report{
			FeedID:    feedID,
			Timestamp: uint32(1700000000 + i),
			Price:     big.NewInt(int64((i + 1) * 100)),
		}
	}

	args := abi.Arguments{{Name: "Reports", Type: typ}}
	packed, err := args.Pack(reports)
	if err != nil {
		return []byte("fallback-payload")
	}
	return packed
}

// buildMaxResponse returns a string near 100 KB (ExecutionResponseLimit).
func buildMaxResponse() string {
	const targetLen = 99 * 1024 // 99 KB, safely under 100 KB
	var sb strings.Builder
	sb.Grow(targetLen)
	line := "[maxlimits] execution-complete "
	for sb.Len() < targetLen {
		sb.WriteString(line)
	}
	return sb.String()[:targetLen]
}
