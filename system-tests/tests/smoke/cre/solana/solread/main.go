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
	case config.TestCaseEVMReadAccountInfo:
		requireAccountInfo(t, runtime, cfg, client)
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

type T struct {
	*slog.Logger
}

func (t *T) Errorf(format string, args ...interface{}) {
	t.Logger.Error(fmt.Sprintf(format, args...))
}

func (t *T) FailNow() {
	panic("Test failed. Panic to stop execution")
}
