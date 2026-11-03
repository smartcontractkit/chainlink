package agentcommerce

import "fmt"

type ExecutionPhase string

const (
	ExecutionPrepared  ExecutionPhase = "PREPARED"
	ExecutionSubmitted ExecutionPhase = "SUBMITTED"
	ExecutionExecuting ExecutionPhase = "EXECUTING"
	ExecutionExecuted  ExecutionPhase = "EXECUTED"
	ExecutionVerified  ExecutionPhase = "VERIFIED"
	ExecutionFailed    ExecutionPhase = "FAILED"
	ExecutionUnknown   ExecutionPhase = "UNKNOWN"
)

func ValidateExecutionTransition(from, to ExecutionPhase) error {
	if from == to {
		return nil
	}
	allowed := map[ExecutionPhase]map[ExecutionPhase]bool{
		ExecutionPrepared:  {ExecutionSubmitted: true, ExecutionFailed: true, ExecutionUnknown: true},
		ExecutionSubmitted: {ExecutionExecuting: true, ExecutionFailed: true, ExecutionUnknown: true},
		ExecutionExecuting: {ExecutionExecuted: true, ExecutionFailed: true, ExecutionUnknown: true},
		ExecutionExecuted:  {ExecutionVerified: true, ExecutionFailed: true, ExecutionUnknown: true},
		ExecutionUnknown:   {ExecutionPrepared: true, ExecutionSubmitted: true, ExecutionExecuting: true, ExecutionExecuted: true, ExecutionFailed: true},
	}
	if !allowed[from][to] {
		return fmt.Errorf("invalid execution transition %s -> %s", from, to)
	}
	return nil
}

func ValidateSettlementTransition(from, to SettlementPhase) error {
	if from == to {
		return nil
	}
	allowed := map[SettlementPhase]map[SettlementPhase]bool{
		SettlementPrepared:  {SettlementSubmitted: true, SettlementRefunded: true, SettlementFailed: true, SettlementUnknown: true},
		SettlementSubmitted: {SettlementBroadcast: true, SettlementRefunded: true, SettlementFailed: true, SettlementUnknown: true},
		SettlementBroadcast: {SettlementConfirmed: true, SettlementRefunded: true, SettlementFailed: true, SettlementUnknown: true},
		SettlementConfirmed: {SettlementSettled: true, SettlementRefunded: true, SettlementFailed: true, SettlementUnknown: true},
		SettlementUnknown:   {SettlementPrepared: true, SettlementSubmitted: true, SettlementBroadcast: true, SettlementConfirmed: true, SettlementFailed: true},
	}
	if !allowed[from][to] {
		return fmt.Errorf("invalid settlement transition %s -> %s", from, to)
	}
	return nil
}
