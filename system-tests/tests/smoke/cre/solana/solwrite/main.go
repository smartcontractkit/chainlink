//go:build wasip1

package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"math/big"

	solgo "github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/solwrite/config"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/solana"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"

	ag_binary "github.com/gagliardetto/binary"
)

func RunSolWriteWorkflow(cfg config.Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	return cre.Workflow[config.Config]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}), // every 30 seconds
			onTrigger,
		),
	}, nil

}
func onTrigger(config config.Config, runtime cre.Runtime, payload *cron.Payload) (string, error) {
	runtime.Logger().Info("Solana Write workflow started", "payload", payload)
	solClient := solana.Client{ChainSelector: chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector}
	runtime.Logger().Info("Got Solana client", "chainSelector", solClient.ChainSelector)
	// 1. Derive remaining
	remaining, err := deriveRemaining(runtime, config)
	if err != nil {
		return "", fmt.Errorf("failed to derive remaining: %w", err)
	}
	// 2. Get account ctx hash
	hash := calculateHash(remaining)
	// 3. Encode report payload
	encodedPayload, err := encodeReport(hash, config)
	if err != nil {
		return "", fmt.Errorf("failed to encode report payload")
	}
	// 4. Generate Report
	report, err := runtime.GenerateReport(&cre.ReportRequest{
		EncodedPayload: encodedPayload,
		EncoderName:    "solana",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await()
	if err != nil {
		return "", fmt.Errorf("failed to generate report: %w", err)
	}
	// 5. Execute WriteReport
	output, err := solClient.WriteReport(runtime, &solana.WriteCreReportRequest{
		Receiver:          config.Receiver.Bytes(),
		Report:            report,
		RemainingAccounts: remaining,
		ComputeConfig: &solana.ComputeConfig{
			ComputeLimit: 290_000,
		},
	}).Await()
	if err != nil {
		runtime.Logger().Error(fmt.Sprintf("[logger] failed to write report on-chain: %v", err))
		return "", fmt.Errorf("failed to write report on solana chain: %w", err)
	}

	runtime.Logger().With().Info("Submitted report on-chain")

	var message = "Solana Workflow successfully completed"
	if output.ErrorMessage != nil {
		message = *output.ErrorMessage
	}

	return message, nil
}

func main() {
	wasm.NewRunner(func(configBytes []byte) (config.Config, error) {
		cfg := config.Config{}
		if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		return cfg, nil
	}).Run(RunSolWriteWorkflow)
}

func deriveRemaining(runtime cre.Runtime, config config.Config) ([]*solana.AccountMeta, error) {
	authority, err := deriveForwarderAuthority(config.ForwarderState, config.Receiver, config.ForwarderProgramID)
	if err != nil {
		return nil, err
	}
	decimalReportSeeds := [][]byte{
		[]byte("decimal_report"),
		config.ReceiverState.Bytes(),
		config.FeedID[:],
	}
	decimalReportKey, _, err := solgo.FindProgramAddress(decimalReportSeeds, config.Receiver)
	if err != nil {
		return nil, err
	}

	hash := createReportHash(config.FeedID[:], authority.Bytes(), config.WFOwner[:], config.WFName[:])
	runtime.Logger().Info(fmt.Sprintf("repHash: %x feedID: %x sender: %x owner: %x name: %x", hash, config.FeedID, authority, config.WFOwner, config.WFName))
	writeFlagSeeds := [][]byte{
		[]byte("permission_flag"),
		config.ReceiverState.Bytes(),
		hash[:],
	}

	writeFlagKey, _, err := solgo.FindProgramAddress(writeFlagSeeds, config.Receiver)
	return []*solana.AccountMeta{
		{PublicKey: config.ForwarderState[:]},              // 0 state
		{PublicKey: authority[:]},                          // 1 authority
		{PublicKey: config.ReceiverState[:]},               // 2 cache state
		{PublicKey: config.Receiver[:]},                    // 3 dummy legacy store
		{PublicKey: config.Receiver[:]},                    // 4 dummy legacy feed config
		{PublicKey: config.Receiver[:]},                    // 5 dummy legacy writer
		{PublicKey: decimalReportKey[:], IsWritable: true}, // 6 decimal report pda
		{PublicKey: writeFlagKey[:]},                       // 7 write permission pda
	}, nil
}

func createReportHash(dataID []byte, forwarderAuthority []byte, workflowOwner []byte, workflowName []byte) [32]byte {
	var data []byte
	data = append(data, dataID...)
	data = append(data, forwarderAuthority...)
	data = append(data, workflowOwner...)
	data = append(data, workflowName...)

	return sha256.Sum256(data)
}

func calculateHash(accs []*solana.AccountMeta) [32]byte {
	var accounts = make([]byte, 0)
	for _, acc := range accs {
		accounts = append(accounts, acc.PublicKey[:]...)
	}
	return sha256.Sum256(accounts)
}

func deriveForwarderAuthority(forwarderState solgo.PublicKey, receiverProgram solgo.PublicKey, forwarderProgram solgo.PublicKey) (solgo.PublicKey, error) {
	seeds := [][]byte{
		[]byte("forwarder"),
		forwarderState[:],
		receiverProgram[:],
	}
	ret, _, err := solgo.FindProgramAddress(seeds, forwarderProgram)

	return ret, err
}

type ReceivedDecimalReport struct {
	Timestamp uint32
	Answer    [16]byte // u128 as 16 little-endian bytes
	DataID    [16]byte
}

type ForwarderReport struct {
	AccountHash [32]byte
	Payload     []byte
}

func encodeReport(accHash [32]byte, cfg config.Config) ([]byte, error) {
	var payloadBuf bytes.Buffer
	payloadEnc := ag_binary.NewBorshEncoder(&payloadBuf)
	var answer [16]byte
	copy(answer[:], big.NewInt(15).Bytes())

	reports := []ReceivedDecimalReport{
		{
			Timestamp: 1,
			Answer:    answer,
			DataID:    cfg.FeedID,
		},
	}
	if err := payloadEnc.Encode(reports); err != nil {
		return nil, fmt.Errorf("failed to borsh-encode ReceivedDecimalReport vec: %w", err)
	}
	payload := payloadBuf.Bytes()

	fr := ForwarderReport{
		AccountHash: accHash,
		Payload:     payload,
	}

	var outBuf bytes.Buffer
	outEnc := ag_binary.NewBorshEncoder(&outBuf)
	if err := outEnc.Encode(fr); err != nil {
		return nil, fmt.Errorf("failed to borsh-encode ForwarderReport: %w", err)
	}

	return outBuf.Bytes(), nil
}
