//go:build wasip1

package main

import (
	"fmt"
	"log/slog"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/solana"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/solread/config"
)

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		cfg := config.Config{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
		}
		return cfg, nil
	}).Run(RunReadWorkflow)
}

func RunReadWorkflow(cfg config.Config, logger *slog.Logger, secretsProvider sdk.SecretsProvider) (sdk.Workflow[config.Config], error) {
	return sdk.Workflow[config.Config]{
		sdk.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onReadTrigger,
		),
	}, nil
}

func onReadTrigger(cfg config.Config, runtime sdk.Runtime, payload *cron.Payload) (_ any, err error) {
	runtime.Logger().Info("onReadTrigger called", "payload", payload)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	t := &T{Logger: runtime.Logger()}
	client := solana.Client{ChainSelector: chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector}
	switch cfg.TestCase {
	case config.TestCaseSolanaReadAccountInfo:
		requireAccountInfo(t, runtime, cfg, client)
	case config.TestCaseSolanaGetBalance:
		requireBalance(t, runtime, cfg, client)
	case config.TestCaseSolanaGetMultipleAccounts:
		requireMultipleAccounts(t, runtime, cfg, client)
	case config.TestCaseSolanaGetProgramAccounts:
		requireProgramAccounts(t, runtime, cfg, client)
	case config.TestCaseSolanaGetBlock:
		requireGetBlock(t, runtime, client)
	case config.TestCaseSolanaGetSlotHeight:
		requireSlotHeight(t, runtime, client)
	case config.TestCaseSolanaGetTransaction:
		requireTransaction(t, runtime, cfg, client)
	case config.TestCaseSolanaGetSignatureStatuses:
		requireSignatureStatuses(t, runtime, cfg, client)
	case config.TestCaseSolanaGetFeeForMessage:
		requireFeeForMessage(t, runtime, cfg, client)
	default:
		panic(fmt.Sprintf("unexpected test case: %s", cfg.TestCase))
	}

	runtime.Logger().Info("Read workflow test case passed for testcase "+cfg.TestCase.String(), "workflow", cfg.WorkflowName)
	return
}

func requireAccountInfo(t *T, runtime sdk.Runtime, cfg config.Config, client solana.Client) {
	accountInfoReply, err := client.GetAccountInfoWithOpts(runtime, &solana.GetAccountInfoWithOptsRequest{
		Account: cfg.AccountAddress,
		Opts: &solana.GetAccountInfoOpts{
			Encoding:       solana.EncodingType_ENCODING_TYPE_JSON_PARSED,
			Commitment:     solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
			DataSlice:      nil,
			MinContextSlot: 0,
		},
	}).Await()
	require.NoError(t, err, "failed to get account info")
	require.NotNil(t, accountInfoReply, "Account info should not be nil")
	require.NotNil(t, accountInfoReply.Value, "Account info value should not be nil")
	runtime.Logger().Info("Account info", "accountInfo", accountInfoReply.Value)
}

func requireBalance(t *T, runtime sdk.Runtime, cfg config.Config, client solana.Client) {
	balanceReply, err := client.GetBalance(runtime, &solana.GetBalanceRequest{
		Addr:       cfg.AccountAddress,
		Commitment: solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
	}).Await()
	require.NoError(t, err, "failed to get balance")
	require.NotNil(t, balanceReply, "balance reply should not be nil")
	require.Equal(t, cfg.ExpectedBalance.Uint64(), balanceReply.Value, "balance should match funded amount")
	runtime.Logger().Info("Balance", "lamports", balanceReply.Value)
}

func requireMultipleAccounts(t *T, runtime sdk.Runtime, cfg config.Config, client solana.Client) {
	multiReply, err := client.GetMultipleAccountsWithOpts(runtime, &solana.GetMultipleAccountsWithOptsRequest{
		Accounts: [][]byte{cfg.AccountAddress},
		Opts: &solana.GetMultipleAccountsOpts{
			Encoding:   solana.EncodingType_ENCODING_TYPE_JSON_PARSED,
			Commitment: solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
		},
	}).Await()
	require.NoError(t, err, "failed to get multiple accounts")
	require.NotNil(t, multiReply, "multiple accounts reply should not be nil")
	require.Len(t, multiReply.Value, 1, "should return exactly 1 account wrapper")
	require.NotNil(t, multiReply.Value[0], "account wrapper should not be nil")
	require.NotNil(t, multiReply.Value[0].Account, "account should not be nil")
	runtime.Logger().Info("Multiple accounts", "count", len(multiReply.Value))
}

func requireProgramAccounts(t *T, runtime sdk.Runtime, cfg config.Config, client solana.Client) {
	programAccountsReply, err := client.GetProgramAccounts(runtime, &solana.GetProgramAccountsRequest{
		Program: cfg.ProgramAddress,
		Opts: &solana.GetProgramAccountsOpts{
			Commitment: solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
		},
	}).Await()
	require.NoError(t, err, "failed to get program accounts")
	require.NotNil(t, programAccountsReply, "program accounts reply should not be nil")
	runtime.Logger().Info("Program accounts", "count", len(programAccountsReply.Value))
}

func requireGetBlock(t *T, runtime sdk.Runtime, client solana.Client) {
	slotReply, err := client.GetSlotHeight(runtime, &solana.GetSlotHeightRequest{
		Commitment: solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
	}).Await()
	require.NoError(t, err, "failed to get slot height")
	require.Greater(t, slotReply.Height, uint64(0), "slot height should be greater than 0")
	runtime.Logger().Info("Current slot height", "slot", slotReply.Height)

	blockReply, err := client.GetBlock(runtime, &solana.GetBlockRequest{
		Slot: slotReply.Height,
		Opts: &solana.GetBlockOpts{
			Commitment: solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
		},
	}).Await()
	require.NoError(t, err, "failed to get block")
	require.NotNil(t, blockReply, "block reply should not be nil")
	require.NotEmpty(t, blockReply.Blockhash, "block hash should not be empty")
	require.Greater(t, blockReply.BlockHeight, uint64(0), "block height should be greater than 0")
	runtime.Logger().Info("Block info",
		"blockHash", blockReply.Blockhash,
		"blockHeight", blockReply.BlockHeight,
		"parentSlot", blockReply.ParentSlot,
	)
}

func requireSlotHeight(t *T, runtime sdk.Runtime, client solana.Client) {
	slotReply, err := client.GetSlotHeight(runtime, &solana.GetSlotHeightRequest{
		Commitment: solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
	}).Await()
	require.NoError(t, err, "failed to get slot height")
	require.NotNil(t, slotReply, "slot height reply should not be nil")
	require.Greater(t, slotReply.Height, uint64(0), "slot height should be greater than 0")
	runtime.Logger().Info("Slot height", "height", slotReply.Height)
}

func requireTransaction(t *T, runtime sdk.Runtime, cfg config.Config, client solana.Client) {
	txReply, err := client.GetTransaction(runtime, &solana.GetTransactionRequest{
		Signature: cfg.TxSignature,
	}).Await()
	require.NoError(t, err, "failed to get transaction")
	require.NotNil(t, txReply, "transaction reply should not be nil")
	require.Greater(t, txReply.Slot, uint64(0), "transaction slot should be greater than 0")
	runtime.Logger().Info("Transaction", "slot", txReply.Slot)
}

func requireSignatureStatuses(t *T, runtime sdk.Runtime, cfg config.Config, client solana.Client) {
	statusReply, err := client.GetSignatureStatuses(runtime, &solana.GetSignatureStatusesRequest{
		Sigs: [][]byte{cfg.TxSignature},
	}).Await()
	require.NoError(t, err, "failed to get signature statuses")
	require.NotNil(t, statusReply, "signature statuses reply should not be nil")
	require.Len(t, statusReply.Results, 1, "should return exactly 1 status result")
	require.NotNil(t, statusReply.Results[0], "status result should not be nil")
	require.Empty(t, statusReply.Results[0].Err, "signature status should have no error")
	runtime.Logger().Info("Signature status",
		"slot", statusReply.Results[0].Slot,
		"confirmationStatus", statusReply.Results[0].ConfirmationStatus,
	)
}

func requireFeeForMessage(t *T, runtime sdk.Runtime, cfg config.Config, client solana.Client) {
	feeReply, err := client.GetFeeForMessage(runtime, &solana.GetFeeForMessageRequest{
		Message:    cfg.EncodedMessage,
		Commitment: solana.CommitmentType_COMMITMENT_TYPE_CONFIRMED,
	}).Await()
	require.NoError(t, err, "failed to get fee for message")
	require.NotNil(t, feeReply, "fee for message reply should not be nil")
	runtime.Logger().Info("Fee for message", "fee", feeReply.Fee)
}

type T struct {
	*slog.Logger
}

func (t *T) Errorf(format string, args ...interface{}) {
	t.Logger.Error(fmt.Sprintf(format, args...))
}

func (t *T) FailNow() {
	panic("Test failed. Panic to stop execution")
}
