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

// OrchestratorConfig holds configuration for the orchestrator.
type OrchestratorConfig struct {
	Directory   *InMemoryDirectory
	Escrow      Escrow
	Reputation  Reputation
	Audit       *AuditLog
	Wallet      Wallet
	Policy      Policy
	VerifyWith  ProofVerifier
	ExecuteAs   ServiceExecutor
	Timeout     time.Duration
	MaxRetries  int
}

// Orchestrator composes discovery, negotiation, signatures, escrow, verification,
// settlement, reputation, and audit logging into one ACP transaction loop.
type Orchestrator struct {
	config OrchestratorConfig
}

// NewOrchestrator creates a new orchestrator with validation.
func NewOrchestrator(cfg OrchestratorConfig) (*Orchestrator, error) {
	if cfg.Directory == nil {
		return nil, fmt.Errorf("directory is required")
	}
	if cfg.Escrow == nil {
		return nil, fmt.Errorf("escrow is required")
	}
	if cfg.Reputation == nil {
		return nil, fmt.Errorf("reputation is required")
	}
	if cfg.Audit == nil {
		return nil, fmt.Errorf("audit log is required")
	}
	if cfg.Wallet == nil {
		return nil, fmt.Errorf("buyer wallet is required")
	}
	if cfg.ExecuteAs == nil {
		return nil, fmt.Errorf("service executor is required")
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 1
	}
	return &Orchestrator{config: cfg}, nil
}

// RunTransaction executes one buyer-to-seller ACP transaction.
func (o *Orchestrator) RunTransaction(ctx context.Context, query ServiceQuery, request ServiceRequest, sellerWallet Wallet) (receipt SettlementReceipt, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in RunTransaction: %v", r)
		}
	}()

	if err := o.validateInputs(query, request, sellerWallet); err != nil {
		return SettlementReceipt{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, o.config.Timeout)
	defer cancel()

	sellers := o.config.Directory.Discover(query)
	if len(sellers) == 0 {
		return SettlementReceipt{}, fmt.Errorf("no seller found for capability %q", query.Capability)
	}

	seller := sellers[0]
	price, ok := seller.Pricing[query.Capability]
	if !ok {
		return SettlementReceipt{}, fmt.Errorf("seller %q does not offer capability %q", seller.ID, query.Capability)
	}

	terms := IntentTerms{
		ServiceRequest: request,
		Price:          price,
		Buyer:          o.config.Wallet.Address(),
		Seller:         seller.ID,
		Timestamp:      time.Now().UTC(),
	}

	if err := EvaluatePolicy(o.config.Policy, seller, terms); err != nil {
		return SettlementReceipt{}, fmt.Errorf("policy evaluation failed: %w", err)
	}

	intent, err := SignIntent(terms, o.config.Wallet, sellerWallet)
	if err != nil {
		return SettlementReceipt{}, fmt.Errorf("signing failed: %w", err)
	}

	if err := VerifySignedIntent(intent, o.config.Wallet); err != nil {
		return SettlementReceipt{}, fmt.Errorf("intent verification failed: %w", err)
	}

	if err := o.config.Audit.Store("intent", intent.Hash, intent); err != nil {
		return SettlementReceipt{}, fmt.Errorf("audit store failed for intent: %w", err)
	}

	escrowReceipt, err := o.config.Escrow.Lock(ctx, intent)
	if err != nil {
		return SettlementReceipt{}, fmt.Errorf("escrow lock failed: %w", err)
	}

	if err := o.config.Audit.Store("escrow", escrowReceipt.EscrowID, escrowReceipt); err != nil {
		return SettlementReceipt{}, fmt.Errorf("audit store failed for escrow: %w", err)
	}

	proof, err := o.executeWithRetry(ctx, intent)
	if err != nil {
		return SettlementReceipt{}, fmt.Errorf("execution failed: %w", err)
	}

	if err := o.config.Audit.Store("proof", proof.ExecutionID, proof); err != nil {
		return SettlementReceipt{}, fmt.Errorf("audit store failed for proof: %w", err)
	}

	result, err := o.verifyWithRetry(ctx, intent, proof)
	if err != nil {
		return SettlementReceipt{}, fmt.Errorf("verification failed: %w", err)
	}

	if err := o.config.Audit.Store("verification", intent.Hash, result); err != nil {
		return SettlementReceipt{}, fmt.Errorf("audit store failed for verification: %w", err)
	}

	if result.Verified {
		receipt, err = o.config.Escrow.Release(ctx, escrowReceipt.EscrowID)
		if err != nil {
			return SettlementReceipt{}, fmt.Errorf("escrow release failed: %w", err)
		}
	} else {
		receipt, err = o.config.Escrow.Refund(ctx, escrowReceipt.EscrowID)
		if err != nil {
			return SettlementReceipt{}, fmt.Errorf("escrow refund failed: %w", err)
		}
	}

	if err := o.updateReputation(ctx, seller.ID, intent.Hash, result.Verified, result.Method); err != nil {
		return SettlementReceipt{}, fmt.Errorf("reputation update failed: %w", err)
	}

	if err := o.config.Audit.Store("settlement", receipt.SettlementID, receipt); err != nil {
		return SettlementReceipt{}, fmt.Errorf("audit store failed for settlement: %w", err)
	}

	return receipt, nil
}

func (o *Orchestrator) validateInputs(query ServiceQuery, request ServiceRequest, sellerWallet Wallet) error {
	if sellerWallet == nil {
		return fmt.Errorf("seller wallet is required")
	}
	if query.Capability == "" {
		return fmt.Errorf("capability is required")
	}
	if request.Service == "" {
		return fmt.Errorf("service is required")
	}
	if len(request.EscrowTerms.ReleaseConditions) == 0 {
		return fmt.Errorf("at least one release condition is required")
	}
	if request.EscrowTerms.Expiration <= 0 {
		return fmt.Errorf("escrow expiration must be positive")
	}
	return nil
}

func (o *Orchestrator) executeWithRetry(ctx context.Context, intent SignedIntent) (Proof, error) {
	var lastErr error
	for attempt := 0; attempt < o.config.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return Proof{}, fmt.Errorf("execution cancelled: %w", ctx.Err())
		default:
		}

		proof, err := o.execute(ctx, intent)
		if err == nil {
			return proof, nil
		}
		lastErr = err
		time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
	}
	return Proof{}, fmt.Errorf("execution failed after %d retries: %w", o.config.MaxRetries, lastErr)
}

func (o *Orchestrator) execute(ctx context.Context, intent SignedIntent) (Proof, error) {
	done := make(chan struct{})
	var output []byte
	var metadata map[string]string
	var err error

	go func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic in ExecuteAs: %v", r)
			}
			close(done)
		}()
		output, metadata, err = o.config.ExecuteAs(ctx, intent)
	}()

	select {
	case <-done:
		if err != nil {
			return Proof{}, err
		}
		if output == nil {
			return Proof{}, fmt.Errorf("execution returned nil output")
		}
		sum := sha256.Sum256(output)
		executionID := "exec-" + intent.Hash
		if len(intent.Hash) > 16 {
			executionID = "exec-" + intent.Hash[:16]
		}
		return Proof{
			OutputHash:  hex.EncodeToString(sum[:]),
			Timestamp:   time.Now().UTC(),
			ExecutionID: executionID,
			Metadata:    metadata,
		}, nil
	case <-ctx.Done():
		return Proof{}, fmt.Errorf("execution timed out: %w", ctx.Err())
	}
}

func (o *Orchestrator) verifyWithRetry(ctx context.Context, intent SignedIntent, proof Proof) (VerificationResult, error) {
	var lastErr error
	for attempt := 0; attempt < o.config.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return VerificationResult{}, fmt.Errorf("verification cancelled: %w", ctx.Err())
		default:
		}

		result, err := o.verify(ctx, intent, proof)
		if err == nil {
			return result, nil
		}
		lastErr = err
		time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
	}
	return VerificationResult{}, fmt.Errorf("verification failed after %d retries: %w", o.config.MaxRetries, lastErr)
}

func (o *Orchestrator) verify(ctx context.Context, intent SignedIntent, proof Proof) (VerificationResult, error) {
	if o.config.VerifyWith == nil {
		if proof.OutputHash == "" {
			return VerificationResult{
				Verified: false,
				Method:   "deterministic",
				Reason:   "empty output hash",
				Time:     time.Now().UTC(),
			}, nil
		}
		return VerificationResult{
			Verified: true,
			Method:   "deterministic",
			Time:     time.Now().UTC(),
		}, nil
	}

	done := make(chan struct{})
	var result VerificationResult
	var err error

	go func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic in VerifyWith: %v", r)
			}
			close(done)
		}()
		result, err = o.config.VerifyWith(ctx, intent, proof)
	}()

	select {
	case <-done:
		if err != nil {
			return VerificationResult{}, err
		}
		return result, nil
	case <-ctx.Done():
		return VerificationResult{}, fmt.Errorf("verification timed out: %w", ctx.Err())
	}
}

func (o *Orchestrator) updateReputation(ctx context.Context, sellerID, intentHash string, verified bool, reason string) error {
	delta := int64(-1)
	if verified {
		delta = 1
	}

	if err := o.config.Reputation.Update(ctx, ReputationEvent{
		AgentID:    sellerID,
		IntentHash: intentHash,
		Delta:      delta,
		Reason:     reason,
	}); err != nil {
		return fmt.Errorf("seller reputation update failed: %w", err)
	}

	if err := o.config.Reputation.Update(ctx, ReputationEvent{
		AgentID:    o.config.Wallet.Address(),
		IntentHash: intentHash,
		Delta:      1,
		Reason:     "buyer-settled",
	}); err != nil {
		return fmt.Errorf("buyer reputation update failed: %w", err)
	}

	return nil
}
