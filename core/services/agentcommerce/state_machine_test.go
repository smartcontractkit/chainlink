package agentcommerce

import "testing"

func TestExecutionAndSettlementStateMachinesFailClosed(t *testing.T) {
	if err := ValidateExecutionTransition(ExecutionPrepared, ExecutionSubmitted); err != nil {
		t.Fatalf("valid execution transition: %v", err)
	}
	if err := ValidateExecutionTransition(ExecutionPrepared, ExecutionVerified); err == nil {
		t.Fatal("skipped execution transition must fail")
	}
	if err := ValidateExecutionTransition(ExecutionExecuting, ExecutionUnknown); err != nil {
		t.Fatalf("UNKNOWN execution must remain recordable: %v", err)
	}
	if err := ValidateExecutionTransition(ExecutionUnknown, ExecutionVerified); err == nil {
		t.Fatal("UNKNOWN execution must not finalize directly")
	}
	if err := ValidateSettlementTransition(SettlementBroadcast, SettlementConfirmed); err != nil {
		t.Fatalf("valid settlement transition: %v", err)
	}
	if err := ValidateSettlementTransition(SettlementBroadcast, SettlementSettled); err == nil {
		t.Fatal("BROADCAST must not become SETTLED directly")
	}
	if err := ValidateSettlementTransition(SettlementConfirmed, SettlementUnknown); err != nil {
		t.Fatalf("UNKNOWN settlement must remain recordable: %v", err)
	}
	if err := ValidateSettlementTransition(SettlementUnknown, SettlementSettled); err == nil {
		t.Fatal("UNKNOWN settlement must not finalize directly")
	}
}
