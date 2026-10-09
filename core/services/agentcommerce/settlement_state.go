package agentcommerce

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type SettlementPhase string

const (
	SettlementPrepared  SettlementPhase = "PREPARED"
	SettlementSubmitted SettlementPhase = "SUBMITTED"
	SettlementBroadcast SettlementPhase = "BROADCAST"
	SettlementConfirmed SettlementPhase = "CONFIRMED"
	SettlementSettled   SettlementPhase = "SETTLED"
	SettlementRefunded  SettlementPhase = "REFUNDED"
	SettlementFailed    SettlementPhase = "FAILED"
	SettlementUnknown   SettlementPhase = "UNKNOWN"
)

type SettlementScope struct {
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
	ScopeID  string `json:"scope_id"`
}

type RailSettlementRequest struct {
	Version           string
	Scope             SettlementScope
	IntentHash        string
	ExecutionID       string
	Rail              string
	Network           string
	Source            string
	Destination       string
	Asset             string
	Amount            string
	AuthorizationHash string
}

type RailSettlementReceipt struct {
	Version      string
	RequestHash  string
	Rail         string
	Phase        SettlementPhase
	ExternalID   string
	FinalityHash string
	ObservedAt   time.Time
	ReceiptHash  string
}

type SettlementEvidenceVerifier interface {
	VerifySettlementEvidence(context.Context, RailSettlementRequest, RailSettlementReceipt) error
}

func SettlementRequestHash(req RailSettlementRequest) (string, error) {
	if err := validateRailSettlementRequest(req); err != nil {
		return "", err
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func SealSettlementReceipt(receipt RailSettlementReceipt) (RailSettlementReceipt, error) {
	if receipt.Version == "" || receipt.RequestHash == "" || receipt.Rail == "" {
		return RailSettlementReceipt{}, errors.New("incomplete settlement receipt")
	}
	if receipt.ObservedAt.IsZero() {
		return RailSettlementReceipt{}, errors.New("settlement observation time is required")
	}
	switch receipt.Phase {
	case SettlementPrepared, SettlementSubmitted, SettlementBroadcast, SettlementConfirmed, SettlementSettled, SettlementRefunded, SettlementFailed, SettlementUnknown:
	default:
		return RailSettlementReceipt{}, errors.New("invalid settlement phase")
	}
	copy := receipt
	copy.ReceiptHash = ""
	payload, err := json.Marshal(copy)
	if err != nil {
		return RailSettlementReceipt{}, err
	}
	sum := sha256.Sum256(payload)
	receipt.ReceiptHash = hex.EncodeToString(sum[:])
	return receipt, nil
}

func VerifyFinalSettlement(ctx context.Context, req RailSettlementRequest, receipt RailSettlementReceipt, verifier SettlementEvidenceVerifier) error {
	if ctx == nil {
		return errors.New("context is required")
	}
	if verifier == nil {
		return errors.New("settlement evidence verifier is required")
	}
	wantRequestHash, err := SettlementRequestHash(req)
	if err != nil {
		return err
	}
	if receipt.RequestHash != wantRequestHash || receipt.Rail != req.Rail {
		return errors.New("settlement receipt identity mismatch")
	}
	sealed, err := SealSettlementReceipt(RailSettlementReceipt{
		Version: receipt.Version, RequestHash: receipt.RequestHash, Rail: receipt.Rail,
		Phase: receipt.Phase, ExternalID: receipt.ExternalID,
		FinalityHash: receipt.FinalityHash, ObservedAt: receipt.ObservedAt,
	})
	if err != nil || sealed.ReceiptHash != receipt.ReceiptHash {
		return errors.New("invalid settlement receipt seal")
	}
	if receipt.Phase != SettlementSettled {
		return fmt.Errorf("settlement is %s, not SETTLED", receipt.Phase)
	}
	if receipt.ExternalID == "" || receipt.FinalityHash == "" {
		return errors.New("final settlement evidence is incomplete")
	}
	return verifier.VerifySettlementEvidence(ctx, req, receipt)
}

func validateRailSettlementRequest(req RailSettlementRequest) error {
	if err := validateSettlementScope(req.Scope); err != nil {
		return err
	}
	values := map[string]string{
		"version": req.Version, "intent hash": req.IntentHash,
		"execution id": req.ExecutionID, "rail": req.Rail,
		"network": req.Network, "source": req.Source,
		"destination": req.Destination, "asset": req.Asset,
		"amount": req.Amount, "authorization hash": req.AuthorizationHash,
	}
	for name, value := range values {
		if err := requireCleanValue(name, value); err != nil {
			return err
		}
	}
	if err := ValidateSHA256Hex("intent hash", req.IntentHash); err != nil {
		return err
	}
	if err := ValidateSHA256Hex("authorization hash", req.AuthorizationHash); err != nil {
		return err
	}
	wantExecutionID, err := ExecutionIDForIntent(req.IntentHash)
	if err != nil {
		return err
	}
	if req.ExecutionID != wantExecutionID {
		return errors.New("settlement execution id mismatch")
	}
	if req.Rail != "x402" && req.Rail != "ccip" {
		return errors.New("unsupported settlement rail")
	}
	if strings.HasPrefix(req.Amount, "+") || (len(req.Amount) > 1 && req.Amount[0] == '0') {
		return errors.New("settlement amount must use canonical base-10 form")
	}
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok || amount.Sign() <= 0 {
		return errors.New("settlement amount must be a positive base-10 integer")
	}
	return nil
}

func validateSettlementScope(scope SettlementScope) error {
	for name, value := range map[string]string{"tenant id": scope.TenantID, "user id": scope.UserID, "scope id": scope.ScopeID} {
		if err := requireCleanValue(name, value); err != nil {
			return err
		}
	}
	return nil
}

func settlementScopeKey(scope SettlementScope) (string, error) {
	if err := validateSettlementScope(scope); err != nil {
		return "", err
	}
	payload, err := json.Marshal(scope)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}
