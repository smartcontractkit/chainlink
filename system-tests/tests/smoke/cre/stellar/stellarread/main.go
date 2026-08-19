//go:build wasip1

package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/stellar"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/stellar/stellarread/config"
)

// readContractBatchPassedLog is the single aggregate log the ReadContract test watches for;
// it is emitted only after every case in the batch asserted successfully.
const readContractBatchPassedLog = "Stellar ReadContract batch passed"

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		cfg := config.Config{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("unmarshal config: %w", err)
		}
		return cfg, nil
	}).Run(RunReadWorkflow)
}

func RunReadWorkflow(cfg config.Config, logger *slog.Logger, secretsProvider sdk.SecretsProvider) (sdk.Workflow[config.Config], error) {
	schedule := cfg.CronSchedule
	if schedule == "" {
		schedule = "*/30 * * * * *"
	}
	return sdk.Workflow[config.Config]{
		sdk.Handler(
			cron.Trigger(&cron.Config{Schedule: schedule}),
			func(cfg config.Config, runtime sdk.Runtime, _ *cron.Payload) (_ any, err error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("panic: %v", r)
					}
				}()
				client := stellar.Client{ChainSelector: cfg.ChainSelector}
				if cfg.ReadKind == config.TestCaseReadContract {
					return runReadContractBatch(cfg, runtime, client)
				}
				return runGetLatestLedger(cfg, runtime, client)
			},
		),
	}, nil
}

func runGetLatestLedger(cfg config.Config, runtime sdk.Runtime, client stellar.Client) (any, error) {
	runtime.Logger().Info("onStellarReadTrigger called", "workflow", cfg.WorkflowName)
	reply, err := client.GetLatestLedger(runtime, &stellar.GetLatestLedgerRequest{}).Await()
	if err != nil {
		msg := fmt.Sprintf("Stellar read failed: GetLatestLedger error: %v", err)
		runtime.Logger().Info(msg, "workflow", cfg.WorkflowName, "chainSelector", cfg.ChainSelector)
		return nil, fmt.Errorf("Stellar GetLatestLedger: %w", err)
	}
	if reply == nil {
		runtime.Logger().Info("Stellar read failed: GetLatestLedger reply is nil", "workflow", cfg.WorkflowName)
		return nil, errors.New("GetLatestLedger reply is nil")
	}
	if uint64(reply.Sequence) < cfg.MinLedgerSequence {
		msg := fmt.Sprintf("Stellar read failed: ledger sequence %d below expected minimum %d", reply.Sequence, cfg.MinLedgerSequence)
		runtime.Logger().Info(msg, "workflow", cfg.WorkflowName)
		return nil, fmt.Errorf("ledger sequence %d below expected minimum %d", reply.Sequence, cfg.MinLedgerSequence)
	}

	runtime.Logger().Info("Stellar read consensus succeeded", "sequence", reply.Sequence, "workflow", cfg.WorkflowName)
	return nil, nil
}

// pendingCase pairs a case with its in-flight ReadContract promise (or an arg-build error
// captured before dispatch). promise is nil when argErr is set.
type pendingCase struct {
	step    config.ReadContractAsserts
	promise sdk.Promise[*stellar.ReadContractResponse]
	argErr  error
}

// runReadContractBatch runs and asserts every case in cfg.Cases within this single trigger.
func runReadContractBatch(cfg config.Config, runtime sdk.Runtime, client stellar.Client) (any, error) {
	if len(cfg.Cases) == 0 {
		runtime.Logger().Info("Stellar ReadContract failed: no cases provided", "workflow", cfg.WorkflowName)
		return nil, errors.New("no ReadContract cases provided")
	}

	pending := make([]pendingCase, len(cfg.Cases))
	for i, step := range cfg.Cases {
		args, err := buildArgs(step)
		if err != nil {
			pending[i] = pendingCase{step: step, argErr: fmt.Errorf("bad args: %w", err)}
			continue
		}
		pending[i] = pendingCase{
			step: step,
			promise: client.ReadContract(runtime, &stellar.ReadContractRequest{
				ContractId:    step.ContractID,
				Function:      step.Function,
				Args:          args,
				SourceAccount: cfg.SourceAccount,
			}),
		}
	}

	for _, p := range pending {
		if err := assertReadContractCase(p); err != nil {
			runtime.Logger().Info("Stellar ReadContract case failed", "case", p.step.Name, "workflow", cfg.WorkflowName, "err", err.Error())
			return nil, fmt.Errorf("case %q: %w", p.step.Name, err)
		}
		runtime.Logger().Info("Stellar ReadContract case passed", "case", p.step.Name, "function", p.step.Function, "workflow", cfg.WorkflowName)
	}
	runtime.Logger().Info(readContractBatchPassedLog, "cases", len(cfg.Cases), "workflow", cfg.WorkflowName)
	return nil, nil
}

// assertReadContractCase awaits one dispatched ReadContract call and asserts it against the
// case's expectation. Returns nil on success, or a descriptive error on any deviation.
func assertReadContractCase(p pendingCase) error {
	step := p.step
	// A bad-args error is a hard failure even for ExpectError cases: it means the test
	// itself is malformed, not that the contract rejected the call.
	if p.argErr != nil {
		return p.argErr
	}
	reply, err := p.promise.Await()

	if step.ExpectError {
		switch {
		case err != nil:
			return nil // infra/consensus error is the expected outcome
		case reply != nil && reply.Error != "":
			return nil // contract-level error is the expected outcome
		default:
			return errors.New("expected an error, got success")
		}
	}
	if err != nil {
		return fmt.Errorf("consensus error: %w", err)
	}
	if reply == nil || reply.Error != "" {
		msg := "nil reply"
		if reply != nil {
			msg = reply.Error
		}
		return fmt.Errorf("contract error: %s", msg)
	}
	if step.ExpectedResult != "" && reply.Result != step.ExpectedResult {
		return fmt.Errorf("result mismatch: got %q want %q", reply.Result, step.ExpectedResult)
	}
	return nil
}

func buildArgs(step config.ReadContractAsserts) ([]*stellar.ScVal, error) {
	switch step.ArgMode {
	case "", "none":
		return nil, nil
	case "echo_bool":
		return []*stellar.ScVal{{Value: &stellar.ScVal_B{B: step.ArgBool}}}, nil
	case "echo_u32":
		return []*stellar.ScVal{{Value: &stellar.ScVal_U32{U32: step.ArgU32A}}}, nil
	case "sum_u32":
		return []*stellar.ScVal{
			{Value: &stellar.ScVal_U32{U32: step.ArgU32A}},
			{Value: &stellar.ScVal_U32{U32: step.ArgU32B}},
		}, nil
	case "echo_string":
		return []*stellar.ScVal{{Value: &stellar.ScVal_Str{Str: step.ArgString}}}, nil
	case "echo_bytes":
		b, err := hex.DecodeString(step.ArgBytesHex)
		if err != nil {
			return nil, fmt.Errorf("bad argBytesHex: %w", err)
		}
		return []*stellar.ScVal{{Value: &stellar.ScVal_BytesVal{BytesVal: b}}}, nil
	case "echo_i64":
		return []*stellar.ScVal{{Value: &stellar.ScVal_I64{I64: step.ArgI64}}}, nil
	case "echo_symbol":
		return []*stellar.ScVal{{Value: &stellar.ScVal_Sym{Sym: step.ArgSymbol}}}, nil
	case "echo_vec_u32":
		vals := make([]*stellar.ScVal, len(step.ArgVecU32))
		for i, u := range step.ArgVecU32 {
			vals[i] = &stellar.ScVal{Value: &stellar.ScVal_U32{U32: u}}
		}
		return []*stellar.ScVal{{Value: &stellar.ScVal_Vec{Vec: &stellar.ScVec{Values: vals}}}}, nil
	case "echo_address":
		recv, err := hex.DecodeString(step.ArgReceiverContractIDHex)
		if err != nil || len(recv) != 32 {
			return nil, errors.New("bad receiver contract id hex")
		}
		return []*stellar.ScVal{{Value: &stellar.ScVal_Address{Address: &stellar.ScAddress{Address: &stellar.ScAddress_ContractId{ContractId: recv}}}}}, nil
	case "echo_triple":
		return buildTripleArgs(step)
	default:
		return nil, fmt.Errorf("unknown argMode %q", step.ArgMode)
	}
}

func buildTripleArgs(step config.ReadContractAsserts) ([]*stellar.ScVal, error) {
	recv, err := hex.DecodeString(step.ArgReceiverContractIDHex)
	if err != nil || len(recv) != 32 {
		return nil, errors.New("bad receiver contract id hex")
	}
	execID, err := hex.DecodeString(step.ArgWorkflowExecutionIDHex)
	if err != nil || len(execID) != 32 {
		return nil, errors.New("bad workflow execution id hex")
	}
	rid, err := hex.DecodeString(step.ArgReportIDHex)
	if err != nil || len(rid) != 2 {
		return nil, errors.New("bad report id hex")
	}
	return []*stellar.ScVal{
		{Value: &stellar.ScVal_Address{Address: &stellar.ScAddress{Address: &stellar.ScAddress_ContractId{ContractId: recv}}}},
		{Value: &stellar.ScVal_BytesVal{BytesVal: execID}},
		{Value: &stellar.ScVal_BytesVal{BytesVal: rid}},
	}, nil
}
