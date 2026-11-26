package agentcommerce

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type nativeExecFunc func(context.Context, commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error)

func (f nativeExecFunc) RegisterToWorkflow(context.Context, commoncap.RegisterToWorkflowRequest) error {
	return nil
}
func (f nativeExecFunc) UnregisterFromWorkflow(context.Context, commoncap.UnregisterFromWorkflowRequest) error {
	return nil
}
func (f nativeExecFunc) Execute(ctx context.Context, req commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	return f(ctx, req)
}

func nativeFixture(t *testing.T) ([]byte, ocrtypes.ConfigDigest, ocrtypes.OnchainPublicKey) {
	t.Helper()
	seed := sha256.Sum256([]byte("native-capability-test-key"))
	key, err := crypto.ToECDSA(seed[:])
	if err != nil {
		t.Fatal(err)
	}
	digestSeed := sha256.Sum256([]byte("native-capability-test-digest"))
	var digest ocrtypes.ConfigDigest
	copy(digest[:], digestSeed[:])
	return crypto.FromECDSA(key), digest, crypto.PubkeyToAddress(key.PublicKey).Bytes()
}

func signedNativeResponse(t *testing.T, req commoncap.CapabilityRequest, keyBytes []byte, digest ocrtypes.ConfigDigest, output []byte) commoncap.CapabilityResponse {
	t.Helper()
	metadata := commoncap.ResponseMetadata{CapDON_N: 1, Metering: []commoncap.MeteringNodeDetail{{SpendUnit: "COMPUTE", SpendValue: "1"}}}
	report, err := commoncap.ResponseToReportData(req.Metadata.WorkflowExecutionID, req.Metadata.ReferenceID, output, metadata)
	if err != nil {
		t.Fatal(err)
	}
	key, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		t.Fatal(err)
	}
	sig, err := crypto.Sign(ocr2key.ReportToSigData3(digest, 1, report[:]), key)
	if err != nil {
		t.Fatal(err)
	}
	return commoncap.CapabilityResponse{Payload: &anypb.Any{Value: output}, Metadata: metadata, OCRAttestation: &commoncap.OCRAttestation{ConfigDigest: digest, SequenceNumber: 1, Sigs: []commoncap.AttributedSignature{{Signer: 0, Signature: sig}}}}
}

func validNativeBackend(t *testing.T, executable commoncap.Executable) (*NativeCapabilityBackend, ChainlinkExecutionRequest, NativeOCRPolicy) {
	t.Helper()
	key, digest, signer := nativeFixture(t)
	_ = key
	_ = digest
	intent := sha256.Sum256([]byte("native-capability-intent"))
	intentHash := hex.EncodeToString(intent[:])
	executionID, _ := ExecutionIDForIntent(intentHash)
	policy := NativeOCRPolicy{Signers: []ocrtypes.OnchainPublicKey{signer}, MinSignatures: 1, Issuer: "native-test-don", MaxLifetime: time.Minute}
	backend, err := NewNativeCapabilityBackend(NativeCapabilityBackendConfig{Executable: executable, Policy: policy, Network: "local-network", WorkflowID: "workflow", WorkflowOwner: "owner", CapabilityID: "capability", Method: "execute"})
	if err != nil {
		t.Fatal(err)
	}
	return backend, ChainlinkExecutionRequest{IntentHash: intentHash, ExecutionID: executionID, Network: "local-network", WorkflowID: "workflow", CapabilityID: "capability", RequestedAt: "2026-10-08T12:00:00Z"}, policy
}

func TestNativeCapabilityRejectsMalformedAndSubstitutedResponses(t *testing.T) {
	key, digest, _ := nativeFixture(t)
	tests := []struct {
		name   string
		mutate func(*commoncap.CapabilityResponse)
	}{
		{name: "empty output", mutate: func(r *commoncap.CapabilityResponse) { r.Payload.Value = nil }},
		{name: "missing OCR", mutate: func(r *commoncap.CapabilityResponse) { r.OCRAttestation = nil }},
		{name: "invalid signature", mutate: func(r *commoncap.CapabilityResponse) { r.OCRAttestation.Sigs[0].Signature[0] ^= 1 }},
		{name: "duplicate signer", mutate: func(r *commoncap.CapabilityResponse) {
			r.OCRAttestation.Sigs = append(r.OCRAttestation.Sigs, r.OCRAttestation.Sigs[0])
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			exec := nativeExecFunc(func(_ context.Context, req commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
				response := signedNativeResponse(t, req, key, digest, []byte("output"))
				test.mutate(&response)
				return response, nil
			})
			backend, request, _ := validNativeBackend(t, exec)
			if _, err := backend.Execute(context.Background(), request); err == nil {
				t.Fatal("expected rejection")
			}
		})
	}
	validExec := nativeExecFunc(func(_ context.Context, req commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
		return signedNativeResponse(t, req, key, digest, []byte("output")), nil
	})
	backend, request, _ := validNativeBackend(t, validExec)
	request.Network = "substituted"
	if _, err := backend.Execute(context.Background(), request); err == nil {
		t.Fatal("authority substitution accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := backend.Execute(ctx, request); err == nil {
		t.Fatal("cancellation ignored")
	}
}
