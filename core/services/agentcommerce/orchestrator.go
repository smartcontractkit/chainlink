package agentcommerce

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// ServiceExecutor performs seller-side work after escrow has been locked.
type ServiceExecutor func(context.Context, SignedIntent) ([]byte, map[string]string, error)

// ProofVerifier performs buyer-side verification of the seller proof.
type ProofVerifier func(context.Context, SignedIntent, Proof) (VerificationResult, error)

// Orchestrator composes discovery, negotiation, signatures, escrow, verification,
// settlement, reputation, and audit logging into one ACP transaction loop.
type Orchestrator struct {
	Directory  *InMemoryDirectory
	Escrow     Escrow
	Reputation Reputation
	Audit      *AuditLog
	Wallet     Wallet
	Policy     Policy
	VerifyWith ProofVerifier
	ExecuteAs  ServiceExecutor
}

// RunTransaction executes one buyer-to-seller ACP transaction. The method is intentionally
// rail-agnostic: Amount.Rail can represent x402, stablecoins, EVM, Solana, bank rails, or a
// Chainlink-native adapter that releases via oracle-verified conditions.
func (o *Orchestrator) RunTransaction(ctx context.Context, query ServiceQuery, request ServiceRequest, sellerWallet Wallet) (SettlementReceipt, error) {
	if o.Directory == nil {
		return SettlementReceipt{}, fmt.Errorf("directory is required")
	}
	if o.Escrow == nil {
		return SettlementReceipt{}, fmt.Errorf("escrow is required")
	}
	if o.Reputation == nil {
		return SettlementReceipt{}, fmt.Errorf("reputation is required")
	}
	if o.Audit == nil {
		return SettlementReceipt{}, fmt.Errorf("audit log is required")
	}
	if o.Wallet == nil {
		return SettlementReceipt{}, fmt.Errorf("buyer wallet is required")
	}
	if o.ExecuteAs == nil {
		return SettlementReceipt{}, fmt.Errorf("service executor is required")
	}
	sellers := o.Directory.Discover(query)
	if len(sellers) == 0 {
		return SettlementReceipt{}, fmt.Errorf("no seller found for capability %q", query.Capability)
	}
	seller := sellers[0]
	price := seller.Pricing[query.Capability]
	terms := IntentTerms{ServiceRequest: request, Price: price, Buyer: o.Wallet.Address(), Seller: seller.ID, Timestamp: time.Now().UTC()}
	if err := EvaluatePolicy(o.Policy, seller, terms); err != nil {
		return SettlementReceipt{}, err
	}
	intent, err := SignIntent(terms, o.Wallet, sellerWallet)
	if err != nil {
		return SettlementReceipt{}, err
	}
	if err := VerifySignedIntent(intent, o.Wallet); err != nil {
		return SettlementReceipt{}, err
	}
	if err := o.Audit.Store("intent", intent.Hash, intent); err != nil {
		return SettlementReceipt{}, err
	}
	escrowReceipt, err := o.Escrow.Lock(ctx, intent)
	if err != nil {
		return SettlementReceipt{}, err
	}
	if err := o.Audit.Store("escrow", escrowReceipt.EscrowID, escrowReceipt); err != nil {
		return SettlementReceipt{}, err
	}
	proof, err := o.execute(ctx, intent)
	if err != nil {
		return SettlementReceipt{}, err
	}
	if err := o.Audit.Store("proof", proof.ExecutionID, proof); err != nil {
		return SettlementReceipt{}, err
	}
	result, err := o.verify(ctx, intent, proof)
	if err != nil {
		return SettlementReceipt{}, err
	}
	if err := o.Audit.Store("verification", intent.Hash, result); err != nil {
		return SettlementReceipt{}, err
	}
	var receipt SettlementReceipt
	if result.Verified {
		receipt, err = o.Escrow.Release(ctx, escrowReceipt.EscrowID)
	} else {
		receipt, err = o.Escrow.Refund(ctx, escrowReceipt.EscrowID)
	}
	if err != nil {
		return SettlementReceipt{}, err
	}
	delta := int64(1)
	if !result.Verified {
		delta = -1
	}
	if err := o.Reputation.Update(ctx, ReputationEvent{AgentID: seller.ID, IntentHash: intent.Hash, Delta: delta, Reason: result.Method}); err != nil {
		return SettlementReceipt{}, err
	}
	if err := o.Reputation.Update(ctx, ReputationEvent{AgentID: o.Wallet.Address(), IntentHash: intent.Hash, Delta: 1, Reason: "buyer-settled"}); err != nil {
		return SettlementReceipt{}, err
	}
	if err := o.Audit.Store("settlement", receipt.SettlementID, receipt); err != nil {
		return SettlementReceipt{}, err
	}
	return receipt, nil
}

func (o *Orchestrator) execute(ctx context.Context, intent SignedIntent) (Proof, error) {
	output, metadata, err := o.ExecuteAs(ctx, intent)
	if err != nil {
		return Proof{}, err
	}
	sum := sha256.Sum256(output)
	return Proof{OutputHash: hex.EncodeToString(sum[:]), Timestamp: time.Now().UTC(), ExecutionID: "exec-" + intent.Hash[:16], Metadata: metadata}, nil
}

func (o *Orchestrator) verify(ctx context.Context, intent SignedIntent, proof Proof) (VerificationResult, error) {
	if o.VerifyWith != nil {
		return o.VerifyWith(ctx, intent, proof)
	}
	if proof.OutputHash == "" {
		return VerificationResult{Verified: false, Method: "deterministic", Reason: "empty output hash", Time: time.Now().UTC()}, nil
	}
	return VerificationResult{Verified: true, Method: "deterministic", Time: time.Now().UTC()}, nil
}
