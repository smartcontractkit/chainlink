package vaulttypes

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

var DefaultNamespace = "main"

const (
	// MethodSecretsCreate Note: additional methods should be reflected
	// in the `Methods` list below.
	MethodSecretsCreate = "vault.secrets.create"
	MethodSecretsGet    = "vault.secrets.get"
	MethodSecretsUpdate = "vault.secrets.update"
	MethodSecretsDelete = "vault.secrets.delete"
	MethodSecretsList   = "vault.secrets.list"

	MaxBatchSize = 10
)

var (
	Methods = []string{
		MethodSecretsCreate,
		// Uncomment the line below to expose the GetSecrets method
		// to the Gateway for testing purposes.
		// MethodSecretsGet,
		MethodSecretsUpdate,
		MethodSecretsDelete,
		MethodSecretsList,
	}
)

// SignedOCRResponse is the response format for OCR signed reports, as returned by the Vault DON.
// External clients should verify that the signatures match the payload and context, before trusting this response.
// Only after validating, clients should decode the payload for further processing.
// If however the Error field is non-empty, it indicates there was an error talking to the Vault DON.
type SignedOCRResponse struct {
	Error      string          `json:"error"`
	Payload    json.RawMessage `json:"payload"`
	Context    []byte          `json:"context"`
	Signatures [][]byte        `json:"signatures"`
}

func (r *SignedOCRResponse) String() string {
	return fmt.Sprintf("SignedOCRResponse { Error: %s, Payload: %s, Context: <[%d]byte blob>, Signatures: <[%d][]byte blob>}", r.Error, string(r.Payload), len(r.Context), len(r.Signatures))
}

type SecretsService interface {
	CreateSecrets(ctx context.Context, request *vaultcommon.CreateSecretsRequest) (*Response, error)
	UpdateSecrets(ctx context.Context, request *vaultcommon.UpdateSecretsRequest) (*Response, error)
	GetSecrets(ctx context.Context, requestID string, request *vaultcommon.GetSecretsRequest) (*Response, error)
	DeleteSecrets(ctx context.Context, request *vaultcommon.DeleteSecretsRequest) (*Response, error)
	ListSecretIdentifiers(ctx context.Context, request *vaultcommon.ListSecretIdentifiersRequest) (*Response, error)
}

func KeyFor(id *vaultcommon.SecretIdentifier) string {
	namespace := id.Namespace
	if namespace == "" {
		namespace = DefaultNamespace
	}
	return fmt.Sprintf("%s::%s::%s", id.Owner, namespace, id.Key)
}

type Request struct {
	Payload      proto.Message
	ResponseChan chan *Response

	IDVal         string
	ExpiryTimeVal time.Time
}

func (r *Request) ID() string {
	return r.IDVal
}

func (r *Request) Copy() *Request {
	newRequest := &Request{
		Payload: proto.Clone(r.Payload),

		// intentionally not copied as we want to keep the reference
		ResponseChan: r.ResponseChan,

		// copied by value
		IDVal:         r.IDVal,
		ExpiryTimeVal: r.ExpiryTimeVal,
	}
	return newRequest
}

func (r *Request) ExpiryTime() time.Time {
	return r.ExpiryTimeVal
}

func (r *Request) SendResponse(ctx context.Context, response *Response) {
	select {
	case <-ctx.Done():
		return
	case r.ResponseChan <- response:
	}
}

func (r *Request) SendTimeout(ctx context.Context) {
	r.SendResponse(ctx, &Response{
		ID:    r.IDVal,
		Error: fmt.Sprintf("timeout exceeded: could not process request %s before expiry", r.IDVal),
	})
}

type Response struct {
	ID         string
	Error      string
	Payload    []byte
	Format     string
	Context    []byte
	Signatures [][]byte
}

func (r *Response) ToJSONRPCResult() ([]byte, error) {
	return json.Marshal(SignedOCRResponse{
		Error:      r.Error,
		Payload:    r.Payload,
		Context:    r.Context,
		Signatures: r.Signatures,
	})
}

func (r *Response) RequestID() string {
	return r.ID
}

func (r *Response) String() string {
	return fmt.Sprintf("Response { ID: %s, Error: %s, Payload: %s, Format: %s }", r.ID, r.Error, string(r.Payload), r.Format)
}

func ValidateSignatures(resp *SignedOCRResponse, allowedSigners []common.Address, minRequired int) error {
	if len(resp.Context) < 64 {
		return fmt.Errorf("context too short: expected min 64 bytes, got %d bytes", len(resp.Context))
	}

	if len(resp.Signatures) < minRequired {
		return fmt.Errorf("not enough signatures: expected min %d, got %d", minRequired, len(resp.Signatures))
	}

	// The context contains:
	// 0:32 -> config digest
	// 32:64 -> epoch + round, namely:
	//   - 0:27 -> zero padding
	//   - 27:31 -> sequence number (big endian uint32)
	//   - 31:32 -> zero round value
	// 64:96 -> extra hash (not used by the vault plugin)
	cd, epochRound := resp.Context[:32], resp.Context[32:64]
	configDigest, err := ocr2types.BytesToConfigDigest(cd)
	if err != nil {
		return fmt.Errorf("invalid config digest in signature: %w", err)
	}

	epoch := binary.BigEndian.Uint32(epochRound[27:31])
	round := epochRound[31]

	fullHash := ocr2key.ReportToSigData(ocr2types.ReportContext{
		ReportTimestamp: ocr2types.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        epoch,
			Round:        round,
		},
	}, []byte(resp.Payload))

	validSigners := map[common.Address]bool{}
	for _, s := range resp.Signatures {
		signerPubkey, err := crypto.SigToPub(fullHash, s)
		if err != nil {
			return fmt.Errorf("invalid signature: %w", err)
		}
		signerAddr := crypto.PubkeyToAddress(*signerPubkey)

		for _, as := range allowedSigners {
			if as.Hex() == signerAddr.Hex() {
				validSigners[signerAddr] = true
				break
			}
		}

		if len(validSigners) >= minRequired {
			return nil
		}
	}

	return fmt.Errorf("only %d valid signatures, need at least %d", len(validSigners), minRequired)
}
